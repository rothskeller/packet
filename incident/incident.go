// Package incident manages collections of related messages.
//
// An incident is stored on disk as a directory of message files; each separate
// incident is a separate directory.  Specifically, package incident always
// works with the message files in the current working directory of the calling
// program.
//
// Within the directory, each non-receipt message is stored in a file called
// «LMI».txt, where «LMI» is the local message ID for the message.  If any
// remote message IDs for the message are known, symbolic links name «RMI».txt
// point to «LMI».txt.  (There may be multiple remote message IDs if the message
// was sent to multiple recipients.)
//
// Messages are automatically rendered in PDF format if the message type
// supports it and PDF rendering is built into the program; the PDF version is
// stored in «LMI».pdf, with possible symbolic link from «RMI».pdf.
//
// Delivery and read receipts are stored in «LMI».DR#.txt and «LMI».RR#.txt,
// respectively, where '#' is either absent or a serial number starting with 2.
// (Multiple receipts may be received for a message if it was sent to multiple
// recipients.)  There are no «RMI» symbolic links for those.
//
// On request, package incident can also generate an ICS-309 message log for the
// messages in the directory.  This is stored in CSV format in ics309.csv, and
// if PDF rendering is built into the program, it is rendered in PDF format in
// ics309.pdf as well.  Both files are automatically removed when any message is
// changed, so that the directory does not contain a stale log.
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
	"github.com/rothskeller/packet/xscmsg/delivrcpt"
	"github.com/rothskeller/packet/xscmsg/readrcpt"
)

// MsgIDRE is a regular expression matching a valid message ID.  Its substrings
// are the three-character prefix, the three-or-more-digit sequence number, and
// the optional suffix character.
var MsgIDRE = regexp.MustCompile(`^([0-9][A-Z]{2}|[A-Z][A-Z0-9]{2})-([1-9][0-9]{2,}|0[1-9][0-9]|00[1-9])([A-Z]?)$`)

// MessageExists returns true if a message exists with the specified LMI.
func MessageExists(lmi string) bool {
	if !MsgIDRE.MatchString(lmi) {
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
	if !MsgIDRE.MatchString(rmi) {
		return ""
	}
	if info, err = os.Stat(rmi + ".txt"); err != nil || info.Mode().Type() != os.ModeSymlink {
		return ""
	}
	if lmi, err = os.Readlink(rmi + ".txt"); err != nil || !strings.HasSuffix(lmi, ".txt") {
		return ""
	}
	lmi = lmi[:len(lmi)-4]
	if !MsgIDRE.MatchString(lmi) {
		return ""
	}
	return lmi
}

// ReadMessage reads a message from the incident directory and returns it.
func ReadMessage(lmi string) (env *envelope.Envelope, msg message.Message, err error) {
	var body string

	if env, body, err = readEnvelope(lmi, ""); err != nil {
		return env, nil, err
	}
	msg = message.Decode(env.SubjectLine, body)
	return env, msg, nil
}

// ReadReceipt reads a receipt for a message.  rcpt must be "DR#" or "RR#".
func ReadReceipt(lmi, rcpt string) (env *envelope.Envelope, msg message.Message, err error) {
	var body string

	if env, body, err = readEnvelope(lmi, rcpt); err != nil {
		return env, nil, err
	}
	msg = message.Decode(env.SubjectLine, body)
	return env, msg, nil
}

func readEnvelope(lmi, rcpt string) (env *envelope.Envelope, body string, err error) {
	var (
		fname    string
		contents []byte
	)
	if !MsgIDRE.MatchString(lmi) {
		return nil, "", errors.New("invalid LMI")
	}
	fname = lmi
	if rcpt != "" {
		fname += "." + rcpt
	}
	fname += ".txt"
	if contents, err = os.ReadFile(fname); err != nil {
		return nil, "", err
	}
	if env, body, err = envelope.ParseSaved(string(contents)); err != nil {
		return env, "", fmt.Errorf("stored message could not be parsed: %s", err)
	}
	return env, body, nil
}

// SaveMessage saves a (non-receipt) message to the incident directory,
// overwriting any previous message stored with the same LMI.  If rmi is not
// empty, an RMI symlink is created.  (Existing RMI symlinks are not disturbed.)
// If fast is true, PDFs are not generated even when possible; stale PDFs are
// removed.  If rawsubj is true, the envelope Subject: line is left unchanged
// rather than being regenerated based on the message contents.
func SaveMessage(lmi, rmi string, env *envelope.Envelope, msg message.Message, fast, rawsubj bool) (err error) {
	if !MsgIDRE.MatchString(lmi) {
		return errors.New("invalid LMI")
	}
	if rmi != "" && !MsgIDRE.MatchString(rmi) {
		rmi = "" // ignore ill-formed RMIs
	}
	switch msg.(type) {
	case *delivrcpt.DeliveryReceipt, *readrcpt.ReadReceipt:
		panic("cannot call SaveMessage for receipt message; call SaveReceipt instead")
	}
	if rmi != "" {
		return saveMessage(lmi+".txt", rmi+".txt", env, msg, fast, rawsubj)
	}
	return saveMessage(lmi+".txt", "", env, msg, fast, rawsubj)
}

// SaveReceipt saves a receipt message to the incident directory, with a unique
// sequence number to avoid overwriting other receipts for the same message.
func SaveReceipt(lmi string, env *envelope.Envelope, msg message.Message) (err error) {
	var (
		base     string
		filename string
		seq      = 1
	)
	if !MsgIDRE.MatchString(lmi) {
		return errors.New("invalid LMI")
	}
	switch msg.(type) {
	case *delivrcpt.DeliveryReceipt:
		base = lmi + ".DR"
	case *readrcpt.ReadReceipt:
		base = lmi + ".RR"
	default:
		panic("cannot call SaveReceipt on a non-receipt message")
	}
	filename = base + ".txt"
	for {
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			break
		} else {
			seq++
			filename = fmt.Sprintf("%s%d.txt", base, seq)
		}
	}
	return saveMessage(filename, "", env, msg, true, true)
}

// saveMessage is the common code between SaveMessage and SaveReceipt.
func saveMessage(filename, linkname string, env *envelope.Envelope, msg message.Message, fast, rawsubj bool) (err error) {
	var (
		content string
		modtime time.Time
	)
	// Encode the message.
	if !rawsubj {
		if env.Bulletin {
			env.SubjectLine = msg.EncodeBulletinSubject()
		} else {
			env.SubjectLine = msg.EncodeSubject()
		}
	}
	if b := msg.Base(); b.FHandling != nil && *b.FHandling == "IMMEDIATE" {
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
	filename = filename[:len(filename)-4] + ".pdf"
	if fast {
		os.Remove(filename)
		if linkname != "" {
			os.Remove(linkname[:len(linkname)-4] + ".pdf")
		}
		// This code could leave symlinks to nonexistent PDFs if there
		// are RMI links other than linkname.  TODO
	} else {
		// Render the PDF.  Ignore errors (can't allow them to prevent
		// us from saving a received message).
		if err = msg.RenderPDF(env, filename); err == nil && linkname != "" {
			linkname = linkname[:len(linkname)-4] + ".pdf"
			os.Symlink(filename, linkname)
		}
	}
	return nil
}

// RemoveMessage removes the message with the specified LMI.
func RemoveMessage(lmi string) {
	if !MsgIDRE.MatchString(lmi) {
		panic("invalid LMI")
	}
	os.Remove(lmi + ".txt")
	os.Remove(lmi + ".pdf")
	// This code could leave RMI symlinks to the message.  But client code
	// doesn't call this except for unsent messages, so it shouldn't be an
	// issue.
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
	if match := MsgIDRE.FindStringSubmatch(id); match != nil {
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
		if MsgIDRE.MatchString(lmi) {
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
		if match := MsgIDRE.FindStringSubmatch(mid); match == nil || match[2] != seqstr {
			continue
		}
		switch fi.Mode().Type() {
		case 0: // regular file
			lmis = append(lmis, mid)
		case os.ModeSymlink:
			var target string

			if !remote {
				break
			}
			if target, err = os.Readlink(fi.Name()); err != nil {
				break
			}
			if !strings.HasSuffix(target, ".txt") {
				break
			}
			mid = target[:len(target)-4]
			if MsgIDRE.MatchString(mid) {
				lmis = append(lmis, mid)
			}
		}
	}
	return lmis, nil
}

// A DeliveryInfo structure describes the delivery of a message to a recipient.
type DeliveryInfo struct {
	// Recipient is the address of the recipient to which the message was
	// addressed.  Display names are removed and domains are fleshed out.
	Recipient string
	// DeliveredTime is the date and time when the message was delivered to
	// the recipient, as described in the delivery receipt they sent back.
	// (It is a string because there is no standard time formatting for this
	// delivery receipt field.)  It is empty if no delivery receipt has been
	// received from this recipient.
	DeliveredTime string
	// RemoteMessageID is the message ID assigned to the message by the
	// recipient.  It is empty if no delivery receipt has been received from
	// this recipient.
	RemoteMessageID string
	// receivedBBS is the BBS where we retrieved the delivery receipt, which
	// presumably is also the BBS through which we sent the message.  We
	// need this to resolve addresses without a domain.
	receivedBBS string
}

// Deliveries returns the delivery information for an outgoing message.  One
// DeliveryInfo structure is returned for each distinct To/Cc/Bcc address in the
// message.  An error is returned only if files cannot be read or decoded.
func Deliveries(lmi string) (delivs []*DeliveryInfo, err error) {
	var (
		env   *envelope.Envelope
		addrs []*envelope.Address
		seq   = 2
	)
	if env, _, err = readEnvelope(lmi, ""); err != nil {
		return nil, err
	}
	if env.IsReceived() {
		return nil, fmt.Errorf("%s: not an outgoing message", lmi)
	}
	if addrs, err = envelope.ParseAddressList(env.To); err != nil {
		return nil, fmt.Errorf("%s: invalid To field: %s", lmi, err)
	}
	delivs = make([]*DeliveryInfo, len(addrs))
	for i, addr := range addrs {
		delivs[i] = &DeliveryInfo{Recipient: addr.Address}
	}
	if deliv, err := readDelivery(lmi + ".DR.txt"); err != nil && !os.IsNotExist(err) {
		return nil, err
	} else if deliv != nil {
		delivs = assignDelivery(delivs, deliv)
	}
	for {
		if deliv, err := readDelivery(fmt.Sprintf("%s.DR%d.txt", lmi, seq)); err != nil && !os.IsNotExist(err) {
			return nil, err
		} else if deliv != nil {
			delivs = assignDelivery(delivs, deliv)
		} else {
			break
		}
		seq++
	}
	return delivs, nil
}

// readDelivery reads a delivery receipt and creates a DeliveryInfo.
func readDelivery(fname string) (deliv *DeliveryInfo, err error) {
	var (
		contents []byte
		env      *envelope.Envelope
		body     string
		msg      message.Message
		dr       *delivrcpt.DeliveryReceipt
	)
	if contents, err = os.ReadFile(fname); err != nil {
		return nil, err
	}
	if env, body, err = envelope.ParseSaved(string(contents)); err != nil {
		return nil, fmt.Errorf("%s: stored message could not be parsed: %s", fname, err)
	}
	if addrs, err := envelope.ParseAddressList(env.From); err == nil {
		env.From = addrs[0].Address
	}
	msg = message.Decode(env.SubjectLine, body)
	if dr, _ = msg.(*delivrcpt.DeliveryReceipt); dr == nil || !env.IsReceived() {
		return nil, fmt.Errorf("%s: not a received delivery receipt", fname)
	}
	return &DeliveryInfo{env.From, dr.DeliveredTime, dr.LocalMessageID, env.ReceivedBBS}, nil
}

// assignDelivery assigns a DeliveryInfo to the correct recipient in the
// list of deliveries.  If no matching recipient is found, one is added.
func assignDelivery(delivs []*DeliveryInfo, deliv *DeliveryInfo) []*DeliveryInfo {
	for _, d := range delivs {
		if recipientMatch(d.Recipient, deliv) {
			d.DeliveredTime, d.RemoteMessageID = deliv.DeliveredTime, deliv.RemoteMessageID
			return delivs
		}
	}
	return append(delivs, deliv)
}

// recipientMatch returns whether the Recipient in candidate is a match for the
// target recipient.
func recipientMatch(tgt string, candidate *DeliveryInfo) bool {
	tlocal, tdomain, _ := strings.Cut(tgt, "@")
	clocal, cdomain, _ := strings.Cut(candidate.Recipient, "@")
	if !strings.EqualFold(tlocal, clocal) {
		return false
	}
	if tdomain == "" {
		tdomain = strings.ToLower(candidate.receivedBBS)
	}
	tfirst, trest, _ := strings.Cut(tdomain, ".")
	cfirst, crest, _ := strings.Cut(cdomain, ".")
	if !strings.EqualFold(tfirst, cfirst) {
		return false
	}
	tscco := trest == "" || !strings.EqualFold(trest, "ampr.org") || !strings.EqualFold(trest, "scc-ares-races.org")
	cscco := crest == "" || !strings.EqualFold(crest, "ampr.org") || !strings.EqualFold(crest, "scc-ares-races.org")
	if tscco != cscco {
		return false
	} else if tscco && cscco {
		return true
	}
	return strings.EqualFold(trest, crest)
}
