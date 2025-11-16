// Package receipt defines the message, body, and subject types for Outpost
// receipt messages.
package receipt

import "github.com/rothskeller/packet/errors"

var (
	ErrNewlineInRMSubject = errors.New("The received-message subject must not contain newlines.")
	ErrNewlineInRMTo      = errors.New("The received-message To address must not contain newlines.")
	ErrNoRMSubject        = errors.New("The received-message subject must not be empty.")
	ErrNoRMTo             = errors.New("The received-message To address must not be empty.")
)
