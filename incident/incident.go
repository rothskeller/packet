// Package incident manages collections of related messages.
//
// An incident is stored on disk as a directory of message files; each separate
// incident is a separate directory.  Specifically, package incident always
// works with the message files in the current working directory of the calling
// program.
//
// Within the directory, each non-receipt message is stored in a file called
// «LMI».txt, where «LMI» is the local message ID for the message.  If the
// remote message ID for the message is known, a symbolic link «RMI».txt points
// to «LMI».txt.
//
// Messages are automatically rendered in PDF format if the message type
// supports it; the PDF version is stored in «LMI».pdf, with possible symbolic
// link from «RMI».pdf.
//
// Delivery and read receipts are stored in «LMI».DR.txt and «LMI».RR.txt,
// respectively; there are no «RMI» symbolic links for those.
//
// On request, package incident can also generate an ICS-309 message log for the
// messages in the directory.  This is stored in CSV format in ics309.csv, and
// if the package is built with the "packetpdf" build tag, it is rendered in
// PDF format in ics309.pdf as well.  Both files are automatically removed when
// any message is changed, so that the directory does not contain a stale log.
package incident

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/rothskeller/packet/envelope"
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/delivrcpt"
	"github.com/rothskeller/packet/message/readrcpt"
)

// msgIDRE is a regular expression matching a valid message ID.  Its substrings
// are the three-character prefix, the three-or-more-digit sequence number, and
// the optional suffix character.
var msgIDRE = regexp.MustCompile(`^([0-9][A-Z]{2}|[A-Z][A-Z0-9]{2})-([1-9][0-9]{2,}|0[1-9][0-9]|00[1-9])([A-Z]?)$`)

// MessageExists returns true if a message exists with the specified LMI.
func MessageExists(lmi string) bool {
	if !msgIDRE.MatchString(lmi) {
		return false
	}
	if info, err := os.Stat(lmi + ".txt"); err == nil && info.Mode().IsRegular() {
		return true
	}
	return false
}

// LMIForRMI returns the LMI for the message with the given RMI, if any.  It
// returns "" if the RMI doesn't exist.
func LMIForRMI(rmi string) string {
	var (
		info os.FileInfo
		lmi  string
		err  error
	)
	if !msgIDRE.MatchString(rmi) {
		return ""
	}
	if info, err = os.Stat(rmi + ".txt"); err != nil || info.Mode().Type() != os.ModeSymlink {
		return ""
	}
	if lmi, err = os.Readlink(rmi + ".txt"); err != nil || !strings.HasSuffix(lmi, ".txt") {
		return ""
	}
	lmi = lmi[:len(lmi)-4]
	if !msgIDRE.MatchString(lmi) {
		return ""
	}
	return lmi
}

// ReadMessage reads a message from the incident directory and returns it.
func ReadMessage(lmi string) (env *envelope.Envelope, msg message.Message, err error) {
	var (
		contents []byte
		body     string
	)
	if !msgIDRE.MatchString(lmi) {
		return nil, nil, errors.New("invalid LMI")
	}
	if contents, err = os.ReadFile(lmi + ".txt"); err != nil {
		return nil, nil, err
	}
	if env, body, err = envelope.ParseSaved(string(contents)); err != nil {
		return env, nil, fmt.Errorf("stored message could not be parsed: %s", err)
	}
	msg = message.Decode(env.SubjectLine, body)
	return env, msg, nil
}

// SaveMessage saves a (non-receipt) message to the incident directory,
// overwriting any previous message stored with the same LMI.
func SaveMessage(lmi, rmi string, env *envelope.Envelope, msg message.Message) (err error) {
	if !msgIDRE.MatchString(lmi) {
		return errors.New("invalid LMI")
	}
	if rmi != "" && !msgIDRE.MatchString(rmi) {
		return errors.New("invalid RMI")
	}
	switch msg.(type) {
	case *delivrcpt.DeliveryReceipt, *readrcpt.ReadReceipt:
		panic("cannot call SaveMessage for receipt message; call SaveReceipt instead")
	}
	return saveMessage(lmi+".txt", rmi+".txt", env, msg)
}

// SaveReceipt saves a receipt message to the incident directory, overwriting
// any previous stored receipt of the same type with the same LMI.
func SaveReceipt(lmi string, env *envelope.Envelope, msg message.Message) (err error) {
	var (
		filename string
	)
	if !msgIDRE.MatchString(lmi) {
		return errors.New("invalid LMI")
	}
	switch msg.(type) {
	case *delivrcpt.DeliveryReceipt:
		filename = lmi + ".DR.txt"
	case *readrcpt.ReadReceipt:
		filename = lmi + ".RR.txt"
	default:
		panic("cannot call SaveReceipt on a non-receipt message")
	}
	return saveMessage(filename, "", env, msg)
}

// saveMessage is the common code between SaveMessage and SaveReceipt.
func saveMessage(filename, linkname string, env *envelope.Envelope, msg message.Message) (err error) {
	var (
		content string
		modtime time.Time
	)
	// Encode the message.
	env.SubjectLine = msg.EncodeSubject()
	if msg, ok := msg.(message.HumanMessage); ok && msg.GetHandling() == "IMMEDIATE" {
		env.OutpostUrgent = true
	} else {
		env.OutpostUrgent = false
	}
	content = env.RenderSaved(msg.EncodeBody())
	// Save the message to its text file.
	if err = os.WriteFile(filename, []byte(content), 0666); err != nil {
		return err
	}
	// Set the modification time of the text file based on the envelope.
	// (This allows AllLMIs to sort by file modification time.)
	if env.IsReceived() {
		modtime = env.ReceivedDate
	} else {
		modtime = env.Date
	}
	if !modtime.IsZero() {
		os.Chtimes(filename, modtime, modtime) // error ignored
	}
	// Create the RMI symlink for the text file if needed.
	if linkname != "" {
		os.Symlink(filename, linkname) // error ignored
	}
	// Remove any generated ICS-309 since it's now potentially out of date.
	RemoveICS309s()
	// If the message can be rendered as PDF, do that.
	if msg, ok := msg.(message.IRenderPDF); ok {
		filename = filename[:len(filename)-4] + ".pdf"
		if err = msg.RenderPDF(filename); err != nil {
			return err
		}
		if linkname != "" {
			linkname = linkname[:len(linkname)-4] + ".pdf"
			os.Symlink(filename, linkname) // error ignored
		}
	}
	return nil
}

// RemoveMessage removes the message with the specified LMI.
func RemoveMessage(lmi string) {
	if !msgIDRE.MatchString(lmi) {
		panic("invalid LMI")
	}
	os.Remove(lmi + ".txt")
	os.Remove(lmi + ".pdf")
}

// UniqueMessageID returns the provided message ID if there is no existing
// message in the directory with that ID (local or remote).  Otherwise, it
// increments the sequence number until the message ID is unused, and returns
// the result.
func UniqueMessageID(id string) string {
	var (
		prefix string
		seq    int
		suffix string
	)
	if match := msgIDRE.FindStringSubmatch(id); match != nil {
		prefix, suffix = match[1], match[3]
		seq, _ = strconv.Atoi(match[2])
	} else {
		panic("UniqueMessageID called for invalid ID")
	}
	for {
		if _, err := os.Stat(id + ".txt"); errors.Is(err, os.ErrNotExist) {
			return id
		}
		seq++
		id = fmt.Sprintf("%s-%03d%s", prefix, seq, suffix)
	}
}

// AllLMIs returns a list of local message IDs of all messages in the directory.
// The list is in chronological order.  An error is returned only if the
// directory cannot be read.
func AllLMIs() (lmis []string, err error) {
	var (
		dir   *os.File
		files []os.FileInfo
	)
	if dir, err = os.Open("."); err != nil {
		return nil, err
	}
	defer dir.Close()
	if files, err = dir.Readdir(0); err != nil {
		return nil, err
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime().Before(files[j].ModTime())
	})
	for _, fi := range files {
		var lmi string

		if !fi.Mode().IsRegular() {
			continue
		}
		if !strings.HasSuffix(fi.Name(), ".txt") {
			continue
		}
		lmi = fi.Name()[:len(fi.Name())-4]
		if msgIDRE.MatchString(lmi) {
			lmis = append(lmis, lmi)
		}
	}
	return lmis, nil
}

// SeqToLMI takes a sequence number and expands it to a list of existing message
// LMIs with that sequence number.  If remote is true, the LMIs of messages
// whose RMI has the requested sequence number are also included.  The LMIs are
// returned in unspecified order.  An error is returned only if the directory
// cannot be read.
func SeqToLMI(seq int, remote bool) (lmis []string, err error) {
	var (
		dir    *os.File
		files  []os.FileInfo
		seqstr = fmt.Sprintf("%03d", seq)
	)
	if seq <= 0 {
		panic("SeqToLMI sequence number must be positive")
	}
	if dir, err = os.Open("."); err != nil {
		return nil, err
	}
	defer dir.Close()
	if files, err = dir.Readdir(0); err != nil {
		return nil, err
	}
	for _, fi := range files {
		var mid string

		if !strings.HasSuffix(fi.Name(), ".txt") {
			continue
		}
		mid = fi.Name()[:len(fi.Name())-4]
		if match := msgIDRE.FindStringSubmatch(mid); match == nil || match[2] != seqstr {
			continue
		}
		switch fi.Mode().Type() {
		case 0: // regular file
			lmis = append(lmis, mid)
		case os.ModeSymlink:
			var target string

			if target, err = os.Readlink(fi.Name()); err != nil {
				break
			}
			if !strings.HasSuffix(target, ".txt") {
				break
			}
			mid = target[:len(target)-4]
			if match := msgIDRE.FindStringSubmatch(mid); match != nil && match[2] == seqstr {
				lmis = append(lmis, mid)
			}
		}
	}
	return lmis, nil
}

// RemoteMap returns a map from local message ID to remote message ID for those
// messages that have a remote message ID.  It returns an error only if the
// directory cannot be read.
func RemoteMap() (m map[string]string, err error) {
	var (
		dir   *os.File
		files []os.FileInfo
	)
	if dir, err = os.Open("."); err != nil {
		return nil, err
	}
	defer dir.Close()
	if files, err = dir.Readdir(0); err != nil {
		return nil, err
	}
	m = make(map[string]string)
	for _, fi := range files {
		var lmi, rmi string

		if fi.Mode().Type() != os.ModeSymlink {
			continue
		}
		if !strings.HasSuffix(fi.Name(), ".txt") {
			continue
		}
		if rmi = fi.Name()[:len(fi.Name())-4]; !msgIDRE.MatchString(rmi) {
			continue
		}
		if lmi, err = os.Readlink(fi.Name()); err != nil {
			continue
		}
		if !strings.HasSuffix(lmi, ".txt") {
			continue
		}
		if lmi = lmi[:len(lmi)-4]; !msgIDRE.MatchString(lmi) {
			continue
		}
		m[lmi] = rmi
	}
	return m, nil
}

// HasDeliveryReceipt returns whether the message with the specified LMI has a
// delivery receipt.
func HasDeliveryReceipt(lmi string) bool {
	if !msgIDRE.MatchString(lmi) {
		panic("HasDeliveryReceipt called for invalid LMI")
	}
	if _, err := os.Stat(lmi + ".DR.txt"); err == nil {
		return true
	}
	return false
}
