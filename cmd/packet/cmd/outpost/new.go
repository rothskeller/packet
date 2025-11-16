package outpost

import (
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"os/exec"
	"regexp"

	"github.com/spf13/cobra"

	"github.com/rothskeller/packet/cmd/packet/osdep"
	"github.com/rothskeller/packet/cmd/packet/server"
	"github.com/rothskeller/packet/form/formdef"
	"github.com/rothskeller/packet/form/formdefs"
	"github.com/rothskeller/packet/message"
)

var (
	addonRE    = regexp.MustCompile(`^[A-Z][A-Za-z0-9_]*$`)
	callSignRE = regexp.MustCompile(`(?i)^(?:A[A-L][0-9][A-Z]{1,3}|[KNW][A-Z][0-9][A-Z]{1,3}|[KNW][0-9][A-Z]{2,3})$`)
	msgIDRE    = regexp.MustCompile(`^[A-Za-z0-9][-A-Za-z0-9_]*$`)
	msgTypeRE  = regexp.MustCompile(`^form-[-a-z0-9]+\.html$`)
)

var newCmd = &cobra.Command{
	Use:   "new addon type msgid opcall opname [taccall tacname]",
	Short: "Start an editor on a new draft message",
	Long: `Command called by Outpost to start editing a new draft message.  The parameters are:

addon    Name of the Outpost "addon" for this message (e.g., "SCCoPIFO").
type     Name of the form to create ("form-something.html").
msgid    Message ID to assign to the new message.
opcall   FCC call sign of the operator creating the message.
opname   Name of the operator creating the message.
taccall  Tactical station call sign if any.
tacname  Tactical station name if any.`,
	DisableFlagsInUseLine: true,
	SilenceUsage:          true,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 5 && len(args) != 7 {
			return fmt.Errorf("accepts 5 or 7 arguments, received %d", len(args))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var (
			addon   string
			msgtype string
			msgID   string
			values  url.Values
			mtype   message.MType
		)
		values = make(url.Values)
		// Check addon name syntax.  Existence is checked below.
		if addon = decodeArg(args[0]); !addonRE.MatchString(addon) {
			slog.Error("invalid addon", "addon", addon)
			return fmt.Errorf("invalid addon name %q", addon)
		} else {
			values.Set("addon", addon)
		}
		// Check message type syntax.  Existence is checked below.
		if msgtype = decodeArg(args[1]); !msgTypeRE.MatchString(msgtype) {
			slog.Error("invalid msgtype", "msgtype", msgtype)
			return fmt.Errorf("invalid message type %q", msgtype)
		} else {
			values.Set("msgtype", msgtype)
		}
		// Check message ID syntax.
		if msgID = decodeArg(args[2]); !msgIDRE.MatchString(msgID) {
			slog.Error("invalid msgID", "msgID", msgID)
			return fmt.Errorf("invalid message ID %q", msgID)
		} else {
			values.Set("msgID", msgID)
		}
		// Check operator call syntax.
		if opCall := decodeArg(args[3]); !callSignRE.MatchString(opCall) {
			slog.Error("invalid opcall", "opcall", opCall)
			return fmt.Errorf("invalid operator call sign %q", opCall)
		} else {
			values.Set("opCall", opCall)
		}
		// Store operator name.
		values.Set("opName", decodeArg(args[4]))
		if len(args) == 7 && (args[5] != "" || args[6] != "") {
			// Check tactical call syntax.
			if tacCall := decodeArg(args[5]); !callSignRE.MatchString(tacCall) {
				slog.Error("invalid taccall", "taccall", tacCall)
				return fmt.Errorf("invalid tactical call sign %q", tacCall)
			} else {
				values.Set("tacCall", tacCall)
			}
			// Store tactical name.
			values.Set("tacName", decodeArg(args[6]))
		}
		if mtype = formdefs.Find(func(d *formdef.FormDef) bool {
			return d.AddonName == addon && d.HTMLName == msgtype && len(d.CreateTags) != 0
		}); mtype == nil {
			slog.Error("unknown form", "addon", addon, "msgtype", msgtype)
			return errors.Join(fmt.Errorf("unknown form %s/%s", values.Get("addon"), values.Get("msgtype")), err)
		}
		if err = neCommon(values, "/outpost-new"); err != nil {
			return err
		}
		slog.Info("opened browser for new message", "addon", addon, "msgtype", msgtype, "msgID", msgID)
		return nil
	},
}

func init() {
	Command.AddCommand(newCmd)
}

func neCommon(values url.Values, path string) (err error) {
	var (
		address string
		u       *url.URL
		cmd     *exec.Cmd
	)
	// Get the server address.  This also starts the server if not already
	// running.
	if address, err = server.GetAddress(true); err != nil {
		return fmt.Errorf("starting server: %s", err)
	}
	// Build the request URL.
	if u, err = url.Parse(address); err != nil {
		slog.Error("url.Parse", "url", address, "err", err)
		return fmt.Errorf("parse server address %q: %s", address, err)
	}
	u.Path = path
	u.RawQuery = values.Encode()
	// Open that URL in a browser.
	cmd = osdep.OpenURLCommand(u.String())
	if err = cmd.Run(); err != nil {
		slog.Error("cmd.Run", "url", u, "err", err)
		return fmt.Errorf("open form in browser: %s", err)
	}
	return nil
}
