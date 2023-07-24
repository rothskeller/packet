package shell

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"

	"github.com/rothskeller/packet/envelope"
	"github.com/rothskeller/packet/message"
)

// maxMessageSize is the largest supported message.  Files larger than this are
// not read even if they have a name that looks like a message ID.
const maxMessageSize = 65536

// msgfileRE is the regular expression that a filename must match in order to be
// read as a message.  It matches a local message ID followed by ".txt".  The
// local message ID must be three letters or digits with at least one letter,
// followed by a dash, a positive integer, and an optional suffix letter.
var msgfileRE = regexp.MustCompile(`^((?:[0-9][A-Z]{2}|[A-Z][A-Z0-9]{2})-[0-9]*[1-9][0-9]*[A-Z]?)\.txt$`)

// delivReceiptRE is the regular expression that a filename must match in order
// to be read as a delivery receipt.  It matches a local message ID followed by
// ".DR.txt".
var delivReceiptRE = regexp.MustCompile(`^((?:[0-9][A-Z]{2}|[A-Z][A-Z0-9]{2})-[0-9]*[1-9][0-9]*[A-Z]?)\.DR\.txt$`)

// forEachMessageFile calls the supplied function for each message file in the
// directory, in chronological order by modification time.  The function is
// passed the local message ID from the filename.
func forEachMessageFile(fn func(lmi string)) (err error) {
	var (
		dir   *os.File
		files []os.FileInfo
	)
	if dir, err = os.Open("."); err != nil {
		return err
	}
	defer dir.Close()
	if files, err = dir.Readdir(0); err != nil {
		return err
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime().Before(files[j].ModTime())
	})
	for _, fi := range files {
		if !fi.Mode().IsRegular() {
			continue
		}
		if fi.Size() > maxMessageSize {
			continue
		}
		if match := msgfileRE.FindStringSubmatch(fi.Name()); match != nil {
			fn(match[1])
		}
	}
	return nil
}

// forEachDeliveryReceipt calls the supplied function for each delivery receipt
// file in the directory, in unspecified order.  The function is passed the
// local message ID from the filename.
func forEachDeliveryReceipt(fn func(lmi string)) (err error) {
	var (
		dir   *os.File
		files []os.FileInfo
	)
	if dir, err = os.Open("."); err != nil {
		return err
	}
	defer dir.Close()
	if files, err = dir.Readdir(0); err != nil {
		return err
	}
	for _, fi := range files {
		if !fi.Mode().IsRegular() {
			continue
		}
		if fi.Size() > maxMessageSize {
			continue
		}
		if match := delivReceiptRE.FindStringSubmatch(fi.Name()); match != nil {
			fn(match[1])
		}
	}
	return nil
}

// forEachMessageSymlink calls the supplied function for each symbolic link to a
// message file in the directory, in unspecified order.  The function is passed
// the local and remote message IDs from the filename.
func forEachMessageSymlink(fn func(local, remote string)) (err error) {
	var (
		dir   *os.File
		files []os.FileInfo
	)
	if dir, err = os.Open("."); err != nil {
		return err
	}
	defer dir.Close()
	if files, err = dir.Readdir(0); err != nil {
		return err
	}
	for _, fi := range files {
		var (
			remote string
			dest   string
			local  string
			err    error
		)
		if fi.Mode().Type() != os.ModeSymlink {
			continue
		}
		if match := msgfileRE.FindStringSubmatch(fi.Name()); match != nil {
			remote = match[1]
		} else {
			continue
		}
		if dest, err = os.Readlink(fi.Name()); err != nil {
			continue
		}
		if match := msgfileRE.FindStringSubmatch(dest); match != nil {
			local = match[1]
		} else {
			continue
		}
		for _, fi2 := range files {
			if fi2.Name() == dest {
				if fi2.Mode().IsRegular() && fi2.Size() <= maxMessageSize {
					fn(local, remote)
				}
				break
			}
		}
	}
	return nil
}

// readMessage reads and parses the message with the specified LMI.
func readMessage(lmi string) (env *envelope.Envelope, msg message.Message, err error) {
	var body string

	contents, err := os.ReadFile(lmi + ".txt")
	if err != nil {
		return nil, nil, err
	}
	env, body, err = envelope.ParseSaved(string(contents))
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s.txt: %s\n", lmi, err)
		return env, nil, err
	}
	msg = message.Decode(env.SubjectLine, body)
	return env, msg, nil
}

// saveMessage saves a message.
func saveMessage(lmi string, env *envelope.Envelope, msg message.Message) (err error) {
	env.SubjectLine = msg.EncodeSubject()
	if msg, ok := msg.(message.HumanMessage); ok && msg.GetHandling() == "IMMEDIATE" {
		env.OutpostUrgent = true
	} else {
		env.OutpostUrgent = false
	}
	content := env.RenderSaved(msg.EncodeBody())
	if err = os.WriteFile(lmi+".txt", []byte(content), 0666); err != nil {
		return err
	}
	if msg, ok := msg.(message.IRenderPDF); ok {
		if err = msg.RenderPDF(lmi + ".pdf"); err != nil {
			return err
		}
	}
	return nil
}

// getNextMessageID returns the next message ID with the same prefix and suffix
// as the supplied model.  If any messages exist with that prefix and suffix,
// the sequence number of the returned ID is one higher than the highest
// existing sequence number of messages with that prefix and suffix, but no less
// than the sequence number in the model.  If no such messages exist, the model
// itself is returned.
func getNextMessageID(model string) (msgid string) {
	mparts := msgIDRE.FindStringSubmatch(model)
	seqnum, _ := strconv.Atoi(mparts[2])
	forEachMessageFile(func(lmi string) {
		parts := msgIDRE.FindStringSubmatch(lmi)
		if parts[1] != mparts[1] || parts[3] != mparts[3] {
			return
		}
		var seq, _ = strconv.Atoi(parts[2])
		if seq >= seqnum {
			seqnum = seq + 1
		}
	})
	forEachMessageSymlink(func(_, remote string) {
		parts := msgIDRE.FindStringSubmatch(remote)
		if parts[1] != mparts[1] || parts[3] != mparts[3] {
			return
		}
		var seq, _ = strconv.Atoi(parts[2])
		if seq >= seqnum {
			seqnum = seq + 1
		}
	})
	return fmt.Sprintf("%s-%03d%s", mparts[1], seqnum, mparts[3])
}
