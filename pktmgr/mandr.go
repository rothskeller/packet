package pktmgr

import (
	"log"
	"os"

	"github.com/rothskeller/packet/xscpdf"
)

// MAndR stands for Message-and-Receipts.  It is the collection of a single
// human message along with the associated delivery receipt and/or read receipt.
type MAndR struct {
	// LMI is the local message ID for the human message.
	LMI string
	// RMI is the remote message ID for the human message, if known.
	RMI string
	// Message is the human message.
	M *Message
	// DR is the delivery receipt for the human message (either outbound,
	// for a received human message, or inbound, for a sent human message).
	DR *Message
	// RR is the read receipt for the human message (either outbound, for a
	// received human message, or inbound, for a sent human message).
	RR *Message
	// Unread is a flag indicating that the message has not been read by the
	// user.  Note that this flag is not persistent; on startup all existing
	// messages are considered read.
	Unread bool
	// incident is a pointer to the incident containing the MAndR.
	incident *Incident
}

// Save writes a newly received or modified message to disk.
func (m *MAndR) Save() (err error) {
	if err = m.M.save(m.LMI+".txt", m.LMI); err != nil {
		return err
	}
	if m.RMI != "" {
		_ = os.Symlink(m.LMI+".txt", m.RMI+".txt")
	}
	if m.DR != nil {
		if err = m.DR.save(m.LMI+".DR.txt", ""); err != nil {
			return err
		}
	}
	if m.RR != nil {
		if err = m.RR.save(m.LMI+".RR.txt", ""); err != nil {
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
func (m *MAndR) savePDF() {
	if err := xscpdf.MessageToPDF(m.M.Message, m.LMI+".pdf"); err == nil {
		if m.RMI != "" {
			_ = os.Symlink(m.LMI+".pdf", m.RMI+".pdf")
		}
	} else if err != xscpdf.ErrNoWriter {
		log.Printf("Saving %s.pdf: %s", m.LMI, err)
	}
}

// remove removes a message from disk.
func (m *MAndR) remove() {
	os.Remove(m.LMI + ".txt")
	os.Remove(m.LMI + ".pdf")
	os.Remove(m.LMI + ".DR.txt")
	os.Remove(m.LMI + ".RR.txt")
	if m.RMI != "" {
		os.Remove(m.RMI + ".txt")
		os.Remove(m.RMI + ".pdf")
	}
}
