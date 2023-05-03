package pktmgr

import (
	"fmt"
	"log"
	"net/mail"
	"os"
	"strings"

	"github.com/rothskeller/packet/pktmsg"
	"github.com/rothskeller/packet/typedmsg"
	"github.com/rothskeller/packet/xscmsg"
	"github.com/rothskeller/packet/xscmsg/delivrcpt"
	"github.com/rothskeller/packet/xscmsg/readrcpt"
	"github.com/rothskeller/packet/xscpdf"
)

// Message is a single human message along with the associated delivery receipt
// and/or read receipt.
type Message struct {
	// LMI is the local message ID for the human message.
	LMI string
	// RMI is the remote message ID for the human message, if known.
	RMI string
	// Message is the human message.
	M xscmsg.IMessage
	// DR is the delivery receipt for the human message (either outbound,
	// for a received human message, or inbound, for a sent human message).
	DR *delivrcpt.DeliveryReceipt
	// RR is the read receipt for the human message (either outbound, for a
	// received human message, or inbound, for a sent human message).
	RR *readrcpt.ReadReceipt
	// Unread is a flag indicating that the message has not been read by the
	// user.  Note that this flag is not persistent; on startup all existing
	// messages are considered read.
	Unread bool
	// incident is a pointer to the incident containing the MAndR.
	incident *Incident
}

// IsReceived returns whether the message was received (as opposed to being an
// outgoing message).
func (m *Message) IsReceived() bool { return m.M.GetRxBBS() != "" }

// IsSent returns whether the message was sent (as opposed to being a draft, or
// having been received).
func (m *Message) IsSent() bool { return !m.IsReceived() && m.M.GetSentDate() != "" }

// IsReady returns whether the message is ready to be sent.
func (m *Message) IsReady() bool {
	return !m.IsReceived() && !m.IsSent() && m.M.GetFrom() != ""
}

// IsDRReady returns whether the delivery receipt for the message is ready to be
// sent.
func (m *Message) IsDRReady() bool {
	return m.DR != nil && m.DR.GetRxBBS() == "" && m.DR.GetSentDate() == ""
}

// Save writes a newly received or modified message to disk.
func (m *Message) Save() (err error) {
	if !m.IsReceived() && !m.IsSent() {
		if lmi := m.M.GetOriginMsgID(); lmi != "" && lmi != m.LMI {
			// The LMI has been changed.
			m.remove()
			m.LMI = lmi
		}
	}
	if err = os.WriteFile(m.LMI+".txt", []byte(m.M.Save()), 0666); err != nil {
		return err
	}
	if m.RMI != "" {
		_ = os.Symlink(m.LMI+".txt", m.RMI+".txt")
	}
	if m.DR != nil {
		if err = os.WriteFile(m.LMI+".DR.txt", []byte(m.DR.Save()), 0666); err != nil {
			return err
		}
	}
	if m.RR != nil {
		if err = os.WriteFile(m.LMI+".RR.txt", []byte(m.RR.Save()), 0666); err != nil {
			return err
		}
	}
	if lock := m.incident.config.BackgroundPDF; lock != nil {
		go func() {
			lock.Lock()
			m.savePDF()
			m.incident.ics309()
			lock.Unlock()
		}()
	} else {
		m.savePDF()
		m.incident.ics309()
	}
	return nil
}
func (m *Message) savePDF() {
	if err := xscpdf.MessageToPDF(m.M, m.LMI+".pdf"); err == nil {
		if m.RMI != "" {
			_ = os.Symlink(m.LMI+".pdf", m.RMI+".pdf")
		}
	} else if err != xscpdf.ErrNoWriter {
		log.Printf("ERROR: Saving %s.pdf: %s", m.LMI, err)
	}
}

// remove removes a message from disk.
func (m *Message) remove() {
	os.Remove(m.LMI + ".txt")
	os.Remove(m.LMI + ".pdf")
	os.Remove(m.LMI + ".DR.txt")
	os.Remove(m.LMI + ".RR.txt")
	if m.RMI != "" {
		os.Remove(m.RMI + ".txt")
		os.Remove(m.RMI + ".pdf")
	}
}

// readMessage reads a physical message from a file.  It returns an error if the
// message cannot be read or parsed.
func readMessage(filename string) (tm typedmsg.Message, err error) {
	var (
		contents []byte
		pm       *pktmsg.Message
	)
	if contents, err = os.ReadFile(filename); err != nil {
		return nil, err
	}
	if pm, err = pktmsg.ParseMessage(string(contents)); err != nil {
		return nil, err
	}
	return typedmsg.Recognize(pm), nil
}

// MarkReady marks a message as queued to send.
func (m *Message) MarkReady() {
	m.incident.markReady(m.M)
}

func (i *Incident) markReady(m typedmsg.Message) {
	addr := fmt.Sprintf("%s@%s.ampr.org", strings.ToLower(i.config.callsign), strings.ToLower(i.config.BBS))
	maddr := &mail.Address{Name: i.config.name, Address: addr}
	m.SetFrom(maddr.String())
}

// MarkDraft marks a message as draft (not ready to send).
func (m *Message) MarkDraft() {
	m.M.SetFrom("")
}
