package msgifc

import (
	"io/fs"
	"iter"
	"net/http"
)

// MType is the interface satisfied by all messages types.
type MType interface {
	// Recognize detects whether the argument message is of the type
	// described by this MType, and if so, calls its SetType method to
	// assign this MType to it.
	Recognize(Message)
	// Name returns the name of the message type, as a phrase in lower case
	// (other than acronyms) starting with "a " or "an ".
	Name() string
	// Validate validates the contents of the message and returns any
	// problems.  Flags customize the validation.
	Validate(m Message, flags ValidateFlags) error
	// Fields returns an iterator on the set of message fields.
	Fields(m Message) iter.Seq[Field]
	// RenderPDF creates a PDF representation of the message in the
	// specified file.  If copyname is not empty, it is placed in the
	// footer of each page.  The returned error may be a Warning, showing a
	// non-fatal rendering issue.
	RenderPDF(m Message, filename, copyname string) error
}

// EditableMType is the interface satisfied by a message type that allows
// creation and editing of messages.
type EditableMType interface {
	MType
	// CreateTag returns the tag used to identify this message type on a
	// "new" command line.  The tag is case insensitive.  It is an error
	// for two message types to have the same tag.
	CreateTag() string
	// CreateKey returns the key (usually one or two letters) used to
	// identify this message type in a GUI dialog box for creating a new
	// message (or also on a "new" command line, as an alternative to
	// CreateTag).  The key is case insensitive.  It is an error for two
	// message types to have the same key, or for one to have a key that is
	// a prefix of another's, or for any key to be the same as any
	// CreateTag.  This method may return an empty string, in which case
	//there is no shortcut key in the dialog and the message type must be
	// selected with mouse or arrow keys.)
	CreateKey() string
	// NewDraft returns a new *message.DraftMessage of this type.  It has
	// default values filled in but is otherwise empty.
	NewDraft() Message
	// EditHTML returns the HTML for the edit page to edit the supplied
	// DraftMessage of this type.  vars customize the HTML based on the
	// the editing context.
	EditHTML(msg Message, vars EditHTMLVars) ([]byte, error)
	// EditAssets returns the file system containing assets used by the
	// HTML returned by EditHTML.
	EditAssets() fs.FS
	// FromPOST interprets the form POSTed by the HTML returned by EditHTML,
	// and translates it into a DraftMessage of this message type.  If the
	// POSTed form is invalid, FromPOST may return an error instead.
	FromPOST(r *http.Request) (Message, error)
}

// ValidateFlags are flags for the Validate method.
type ValidateFlags uint

// Values for ValidateFlags.
const (
	// VPIFOOnly indicates that validation should only complain about
	// problems that would block submission in the  PackItForms form editor.
	// If this flag is not set, more extensive checks are acceptable.
	VPIFOOnly ValidateFlags = 1 << iota
	// VPacket indicates that the message should be validated as a packet
	// message (i.e., it's known not to be a voice message).
	VPacket
)

type EditHTMLVars struct {
	// SubmitURL is the URL that the edit form should POST to.  The URL
	// should respond with either an error status with a text/plain body,
	// an http.StatusSeeOther with a redirect, or an http.StatusNoContent
	// (which means clase the window).
	SubmitURL string
	// SubmitLabel is the label of the submit button.  It is required.
	SubmitLabel string
	// SaveLabel is the label of the save button, if any.
	SaveLabel string
	// AssetBase is the URL base for any assets needed by the edit form.
	// It should correspond to the file system returned by EditAssets.
	AssetBase string
	// ShowAddressFields is a boolean indicating whether the To Address and
	// From Address fields should be shown and submitted.  (It's true for
	// manual/GUI editing and false for Outpost editing.)
	ShowAddressFields bool
	// FromAddress is the from address for the message (used only if
	// ShowAddressFields is true).
	FromAddress string
}
