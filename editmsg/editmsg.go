// Package editmsg handles editing of packet messages in a text editor.
package editmsg

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/rothskeller/packet/pktmsg"
	"github.com/rothskeller/packet/xscmsg"
)

// ErrCancel is returned by NewMessage when the user cancels the creation of the
// new message (by saving the template message with no change, or saving an
// empty file).
var ErrCancel = errors.New("editing session canceled")

// ErrUnknownTag is returned by NewMessage when the supplied message type tag is
// not recognized.
var ErrUnknownTag = errors.New("unknown tag")

// NewMessage creates a new message and allows the user to edit it in a text
// editor.  tag identifies the type of message to be created; an empty string
// (or the word "plain") creates a plain text message.  msgnum is the message
// number to be injected into the message, if the message type permits.  opname
// and opcall are the operator name and call sign to be injected into the
// message, if the message type permits.  The new message is saved in
// «msgnum».txt after editing.
func NewMessage(tag, msgnum, opname, opcall string) (err error) {
	var (
		msg *pktmsg.Message
		out string
	)
	if _, err = os.Stat(msgnum + ".txt"); !errors.Is(err, fs.ErrNotExist) {
		return fs.ErrExist
	}
	if tag == "" || tag == "plain" {
		msg = pktmsg.New()
		msg.Header.Set("To", "")
		msg.Header.Set("Subject", msgnum+"_")
	} else {
		xsc := xscmsg.Create(tag)
		if xsc == nil {
			return ErrUnknownTag
		}
		if f := xsc.KeyField(xscmsg.FOriginMsgNo); f != nil {
			f.Value = msgnum
		}
		if f := xsc.KeyField(xscmsg.FOpName); f != nil {
			f.Value = opname
		}
		if f := xsc.KeyField(xscmsg.FOpCall); f != nil {
			f.Value = opcall
		}
		msg = pktmsg.New()
		msg.Body = xsc.Body(true)
		msg.Header.Set("To", "")
		msg.Header.Set("Subject", "(computed)")
	}
	if out, err = editMessage(msg.Encode(true)); err != nil {
		return err
	}
	if err = os.WriteFile(msgnum+".txt", []byte(out), 0666); err != nil {
		os.Remove(msgnum + ".txt")
	}
	return err
}

// EditMessage edits an existing message with the specified message number; the
// message is expected to exist in «msgnum».txt, and changes to it will be saved
// there.
func EditMessage(msgnum string) (err error) {
	var (
		orig []byte
		msg  *pktmsg.Message
		in   string
		out  string
	)
	if orig, err = os.ReadFile(msgnum + ".txt"); err != nil {
		return err
	}
	if msg, _ = pktmsg.ParseMessage(string(orig)); msg == nil {
		in = string(orig)
	} else if xsc := xscmsg.Recognize(msg, false); xsc != nil {
		msg.Header.Set("Subject", xsc.Subject())
		msg.Body = xsc.Body(true)
		in = msg.Encode(true)
	} else {
		in = msg.Encode(true)
	}
	if out, err = editMessage(in); err != nil {
		return err
	}
	if strings.TrimSpace(out) == strings.TrimSpace(string(orig)) {
		return nil
	}
	if err = os.WriteFile(msgnum+".txt", []byte(out), 0666); err != nil {
		os.Remove(msgnum + ".txt")
	}
	return err
}

func editMessage(in string) (out string, err error) {
	var (
		fh       *os.File
		editor   string
		cmd      *exec.Cmd
		buf      []byte
		problems []string
	)
	if fh, err = os.CreateTemp("", "editmsg*.txt"); err != nil {
		return "", err
	}
	defer os.Remove(fh.Name())
	io.WriteString(fh, in)
	fh.Close()
	if editor = os.Getenv("VISUAL"); editor == "" {
		if editor = os.Getenv("EDITOR"); editor == "" {
			if editor, _ = exec.LookPath("vim"); editor == "" {
				if editor, _ = exec.LookPath("vi"); editor == "" {
					if runtime.GOOS == "darwin" {
						editor = "open"
					}
				}
			}
		}
	}
RETRY:
	if editor == "" {
		cmd = exec.Command(fh.Name())
	} else {
		cmd = exec.Command(editor, fh.Name())
	}
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	if err = cmd.Run(); err != nil {
		return "", err
	}
	if buf, err = os.ReadFile(fh.Name()); err != nil {
		return "", err
	}
	out = string(buf)
	if strings.TrimSpace(out) == "" {
		return "", ErrCancel
	}
	if msg, err := pktmsg.ParseMessage(out); err != nil {
		problems = []string{err.Error()}
	} else if xsc := xscmsg.Recognize(msg, false); xsc != nil {
		problems = xsc.Validate(false)
		msg.Header.Set("Subject", xsc.Subject())
		msg.Body = xsc.Body(false)
		out = msg.Encode(false)
	} else if pktmsg.IsForm(msg.Body) {
		if form := pktmsg.ParseForm(msg.Body, false); form == nil {
			problems = []string{"form encoding is invalid"}
		} else {
			problems = []string{fmt.Sprintf("unknown form type %q", form.FormType)}
		}
	} else if xscmsg.ParseSubject(msg.Header.Get("Subject")) == nil {
		if msg.Header.Get("Subject") == "(computed)" {
			problems = []string{"message type was not recognized"}
		} else {
			problems = []string{"subject line encoding is incorrect"}
		}
	}
	if len(problems) != 0 {
		for _, p := range problems {
			fmt.Printf("ERROR: %s\n", p)
		}
		fmt.Print("Continue editing? [Y/n]  ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		if t := scanner.Text(); !strings.HasPrefix(t, "N") && !strings.HasPrefix(t, "n") {
			goto RETRY
		}
	}
	return out, nil
}
