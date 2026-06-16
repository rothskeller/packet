package incident

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/rothskeller/packet/v4/message/messageid"
)

// nextMessageID returns the next local message ID in sequence (and marks it
// used).  If tx is true, it uses the transmit message ID sequence; otherwise
// the receive message ID sequence.  It returns an error only if the config
// contains an invalid message ID template, which shouldn't be possible.
func (i *Incident) nextMessageID(tx bool) (lid string, err error) {
	var next string

	// The configuration has the next message ID to use.
	if tx {
		lid = i.Config.TxMessageID
	} else {
		lid = i.Config.RxMessageID
	}
	// But, we need to be sure the corresponding filename isn't already in
	// use.
	for {
		fname := filepath.Join(i.Dir, lid+".txt")
		if _, err := os.Stat(fname); err != nil { // assume ENOENT
			break
		}
		if lid, err = messageid.Increment(lid); err != nil {
			slog.Error("increment message ID", "err", err)
			return "", err
		}
	}
	// Update the configuration to use the next message ID next time.  If
	// Rx and Tx are sharing the same ID stream, update both.
	if next, err = messageid.Increment(lid); err != nil {
		slog.Error("increment message ID", "err", err)
		return "", err
	}
	if i.Config.TxMessageID == i.Config.RxMessageID {
		i.Config.TxMessageID, i.Config.RxMessageID = next, next
	} else if tx {
		i.Config.TxMessageID = next
	} else {
		i.Config.RxMessageID = next
	}
	return lid, nil
}
