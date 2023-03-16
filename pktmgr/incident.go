// Package pktmgr handles a collection of messages from a single incident (i.e.,
// a single ICS-309 form).
package pktmgr

import (
	"errors"
	"fmt"
	"net/mail"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/rothskeller/packet/xscmsg"
	"github.com/rothskeller/packet/xscmsg/delivrcpt"
	"github.com/rothskeller/packet/xscmsg/readrcpt"
)

// maxMessageSize is the largest supported message.  Files larger than this are
// not read even if they have a name that looks like a message ID.
const maxMessageSize = 65536

// msgfileRE is the regular expression that a filename must match in order to be
// read as a message.  It matches a local message ID followed by ".txt",
// ".DR.txt" (for delivery receipts), or ".RR.txt" (for read receipts).  The
// local message ID must be three letters or digits with at least one letter,
// followed by a dash, a positive integer, and a suffix "P" or "M".
//
// Some senders might send us messages with IDs that are slightly out of spec.
// However, the message files are always named with the local message ID,
// presumed to have been generated by this program, so we can rely on its strict
// definition of message ID syntax.
var msgfileRE = regexp.MustCompile(`^(((?:[0-9][A-Z]{2}|[A-Z][A-Z0-9]{2})-)([0-9]*[1-9][0-9]*)([PM]))(?:\.([DR]R))?\.txt$`)

// An Incident represents a single incident and its related messages.  Each
// Incident is stored in a distinct directory in the file system, and
// corresponds to a single ICS-309 form.
type Incident struct {
	config     Config
	nextMsgNum int
	list       []*MAndR
	index      map[string]*MAndR
}

// NewIncident creates a new Incident with the specified configuration, using
// the current working directory for message storage.
func NewIncident(config Config) (i *Incident, err error) {
	i = new(Incident)
	i.config = config
	i.config.fillin()
	if err = i.Refresh(); err != nil {
		return nil, err
	}
	return i, nil
}

// Refresh re-loads all incident messages from disk.
func (i *Incident) Refresh() (err error) {
	var (
		dir   *os.File
		files []os.FileInfo
	)
	i.nextMsgNum, i.list, i.index = i.config.startMsgNum, i.list[:0], make(map[string]*MAndR)
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
		if err = i.readMessage(fi.Name()); err != nil {
			return err
		}
	}
	for _, mr := range i.list {
		if mr.M == nil {
			return fmt.Errorf("found receipt for non-existent message %s", mr.LMI)
		}
	}
	// Regenerate in case someone else changed message files.
	if i.config.BackgroundPDF != nil {
		go func() {
			i.config.BackgroundPDF.Lock()
			i.ics309()
			i.config.BackgroundPDF.Unlock()
		}()
	} else {
		i.ics309()
	}
	return nil
}

// readMessage reads the message in the specified file, and adds it into the
// list of messages (replacing any previously existing message with the same
// LMI).
func (i *Incident) readMessage(filename string) (err error) {
	var (
		lmi  string
		rcpt string
		mr   *MAndR
		m    *Message
	)
	if match := msgfileRE.FindStringSubmatch(filename); match != nil {
		lmi, rcpt = match[1], match[5]
		if match[2] == i.config.msgIDPrefix && match[4] == i.config.msgIDSuffix {
			num, _ := strconv.Atoi(match[3])
			if num >= i.nextMsgNum {
				i.nextMsgNum = num + 1
			}
		}
	} else {
		return nil
	}
	if m, err = readMessage(filename); err != nil {
		return fmt.Errorf("%s: %s", filename, err)
	}
	if mr = i.GetLMI(lmi); mr == nil {
		mr = i.addLMI(lmi)
	}
	m.mandr = mr
	switch rcpt {
	case "":
		mr.M = m
		if m.IsReceived() {
			if p := xscmsg.ParseSubject(m.RawMessage.Header.Get("Subject")); p != nil {
				mr.RMI = p.MessageNumber
			}
		}
	case "DR":
		mr.DR = m
		if m.IsReceived() {
			mr.RMI = m.Field("LocalMessageID").Value
		}
	case "RR":
		mr.RR = m
	}
	return nil
}

// Count returns the number of MAndRs in the Incident.
func (i *Incident) Count() int {
	return len(i.list)
}

// GetIndex returns the MAndR at the specified numbered position in the list of
// MAndRs.  It returns nil if no such MAndR exists.
func (i *Incident) GetIndex(idx int) *MAndR {
	if idx >= 0 && idx < len(i.list) {
		return i.list[idx]
	}
	return nil
}

// GetLMI returns the MandR with the specified local message ID.  It returns nil
// if no such MandR exists.
func (i *Incident) GetLMI(lmi string) *MAndR {
	return i.index[lmi]
}

// AddDraft creates a new draft message of the specified type with the next
// available local message ID, and returns it.  The local message ID, default
// body, and operator fields of the message are filled in.
func (i *Incident) AddDraft(tag string) (mr *MAndR) {
	xsc := xscmsg.Create(tag)
	if xsc == nil {
		return nil
	}
	lmi := fmt.Sprintf("%s%03d%s", i.config.msgIDPrefix, i.nextMsgNum, i.config.msgIDSuffix)
	mr = i.addLMI(lmi)
	i.nextMsgNum++
	mr.M = &Message{mandr: mr, Message: xsc}
	mr.M.SetOperatorFields()
	if f := mr.M.KeyField(xscmsg.FOriginMsgNo); f != nil {
		f.Value = mr.LMI
	}
	if f := mr.M.KeyField(xscmsg.FBody); f != nil {
		f.Value = i.config.DefBody
	}
	if tag == xscmsg.PlainTextTag {
		mr.M.KeyField(xscmsg.FSubject).Value = mr.LMI + "_"
	}
	mr.Save()
	return mr
}

// Remove removes a message from the Incident.
func (i *Incident) Remove(mr *MAndR) {
	for j, m := range i.list {
		if m == mr {
			i.list = append(i.list[:j], i.list[j+1:]...)
			break
		}
	}
	delete(i.index, mr.LMI)
	mr.remove()
}

// HasBulletin returns whether the Incident contains a bulletin from the
// specified area with the specified subject prefix.
func (i *Incident) HasBulletin(area, subjectPrefix string) bool {
	for j := len(i.list) - 1; j >= 0; j-- {
		if i.list[j].M.Bulletin == area &&
			strings.HasPrefix(i.list[j].M.RawMessage.Header.Get("Subject"), subjectPrefix) {
			return true
		}
	}
	return false
}

// Receive records a newly received message.  Raw is the raw text of the message
// as received from JNOS.  BBS is the name of the BBS from which the message was
// retrieved.  For bulletin messages, bulletin is the area to which the bulletin
// was sent; for all other messages, bulletin is an empty string.  When a
// delivery receipt should be sent for the message, it is queued to be sent and
// also returned by this function.
func (i *Incident) Receive(raw, bbs, bulletin string) (dr *Message, err error) {
	var (
		m  *Message
		mr *MAndR
	)
	if m, err = newMessage(raw); err != nil {
		return nil, err
	}
	m.BBS, m.Bulletin, m.Received = bbs, bulletin, time.Now()
	switch m.Type.Tag {
	case delivrcpt.Tag:
		mr = i.findWithSubject(m.Field("DeliveredSubject").Value)
		if mr == nil {
			return nil, errors.New("dropping received delivery receipt: no matching subject")
		}
		if !mr.M.IsSent() {
			return nil, errors.New("dropping received delivery receipt: matching message was received not sent")
		}
		if mr.DR != nil {
			return nil, errors.New("dropping received delivery receipt: already have one")
		}
		mr.DR = m
		mr.RMI = m.Field("LocalMessageID").Value
	case readrcpt.Tag:
		mr = i.findWithSubject(m.Field("ReadSubject").Value)
		if mr == nil {
			return nil, errors.New("dropping received read receipt: no matching subject")
		}
		if !mr.M.IsSent() {
			return nil, errors.New("dropping received read receipt: matching message was received not sent")
		}
		if mr.RR != nil {
			return nil, errors.New("dropping received read receipt: already have one")
		}
		mr.RR = m
	default:
		lmi := fmt.Sprintf("%s%03d%s", i.config.msgIDPrefix, i.nextMsgNum, i.config.msgIDSuffix)
		mr = i.addLMI(lmi)
		i.nextMsgNum++
		mr.Unread = true
		mr.M = m
		if p := xscmsg.ParseSubject(m.RawMessage.Header.Get("Subject")); p != nil {
			mr.RMI = p.MessageNumber
		}
		if f := mr.M.KeyField(xscmsg.FDestinationMsgNo); f != nil {
			f.Value = lmi
		}
		m.SetOperatorFields()
		if bulletin == "" {
			xdr := xscmsg.Create(delivrcpt.Tag)
			xdr.Field("DeliveredTo").Value = m.RawMessage.Header.Get("To")
			xdr.Field("DeliveredSubject").Value = m.RawMessage.Header.Get("Subject")
			xdr.Field("DeliveredTime").Value = m.Received.Format("1/2/2006 15:04")
			xdr.Field("LocalMessageID").Value = lmi
			mr.DR = &Message{
				Message: xdr,
				To:      []*mail.Address{m.From},
			}
			mr.DR.MarkReady()
			dr = mr.DR
		}
	}
	return dr, mr.Save()
}

// findWithSubject returns the message with the specified subject line.  It
// returns nil if no such message exists.  (This is used for matching delivery
// and read receipts against their primary messages.)
func (i *Incident) findWithSubject(subject string) *MAndR {
	for j := len(i.list) - 1; j >= 0; j-- {
		if i.list[j].M.RawMessage.Header.Get("Subject") == subject {
			return i.list[j]
		}
	}
	return nil
}

// addLMI adds a new MAndR to the Incident, with the specified LMI.
func (i *Incident) addLMI(lmi string) *MAndR {
	i.index[lmi] = &MAndR{LMI: lmi, incident: i}
	i.list = append(i.list, i.index[lmi])
	sort.Slice(i.list, func(a, b int) bool {
		return lmiLess(i.list[a].LMI, i.list[b].LMI)
	})
	return i.index[lmi]
}

// lmiLess returns whether a is less than b, where both are interpreted as
// local message IDs.
func lmiLess(a, b string) bool {
	if a[:3] != b[:3] {
		return a[:3] < b[:3]
	}
	if len(a) != len(b) {
		return len(a) < len(b)
	}
	return a < b
}
