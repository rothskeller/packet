//go:build windows

package cmd

import (
	"fmt"
	"log/slog"
	"net/url"
	"os/exec"
	"regexp"

	"github.com/rothskeller/packet/cmd/packet/cio"
	"github.com/rothskeller/packet/cmd/packet/osdep"
	"github.com/rothskeller/packet/cmd/packet/server"
	"github.com/rothskeller/packet/errors"
	"github.com/rothskeller/packet/form"
	"github.com/rothskeller/packet/message"
	"github.com/spf13/pflag"
)

const (
	outpostNewSlug = `Start an editor on a new draft message`
	outpostNewHelp = `
usage: packet outpost new ⇥addon type msgid opcall opname [taccall tacname actcall]

Command called by Outpost to start editing a new draft message.  The parameters are:

addon    ⇥Name of the Outpost "addon" for this message (e.g., "SCCoPIFO").
type     ⇥Name of the form to create ("form-something.html").
msgid    ⇥Message ID to assign to the new message.
opcall   ⇥FCC call sign of the operator creating the message.
opname   ⇥Name of the operator creating the message.
taccall  ⇥Tactical station call sign if any.
tacname  ⇥Tactical station name if any.
actcall  ⇥Active call sign (opcall or taccall).
`
)

var (
	addonRE    = regexp.MustCompile(`^[A-Z][A-Za-z0-9_]*$`)
	callSignRE = regexp.MustCompile(`(?i)^(?:A[A-L][0-9][A-Z]{1,3}|[KNW][A-Z][0-9][A-Z]{1,3}|[KNW][0-9][A-Z]{2,3})$`)
	tacCallRE  = regexp.MustCompile(`(?i)^[A-Z][A-Z0-9]{3,5}$`)
	msgIDRE    = regexp.MustCompile(`^[A-Za-z0-9][-A-Za-z0-9_]*$`)
	msgTypeRE  = regexp.MustCompile(`^form-[-a-z0-9]+\.html$`)
)

func cmdOutpostNew(args []string) (err error) {
	var (
		addon   string
		msgtype string
		msgID   string
		values  url.Values
		mtype   message.MType
	)
	flags := pflag.NewFlagSet("o-new", pflag.ContinueOnError)
	flags.Usage = func() {} // we do our own
	if err = flags.Parse(args); err == pflag.ErrHelp {
		return cmdOutpostHelp([]string{"new"})
	} else if err != nil {
		cio.Open().Error(err)
		return usage(outpostNewHelp)
	}
	if len(args) != 5 && len(args) != 8 {
		return usage(outpostNewHelp)
	}
	registerForms()
	values = make(url.Values)
	// Check addon name syntax.  Existence is checked below.
	if addon = decodeOutpostArg(args[0]); !addonRE.MatchString(addon) {
		slog.Error("invalid addon", "addon", addon)
		return fmt.Errorf("invalid addon name %q", addon)
	}
	// Check message type syntax.  Existence is checked below.
	if msgtype = decodeOutpostArg(args[1]); !msgTypeRE.MatchString(msgtype) {
		slog.Error("invalid msgtype", "msgtype", msgtype)
		return fmt.Errorf("invalid message type %q", msgtype)
	}
	// Check message ID syntax.
	if msgID = decodeOutpostArg(args[2]); !msgIDRE.MatchString(msgID) {
		slog.Error("invalid msgID", "msgID", msgID)
		return fmt.Errorf("invalid message ID %q", msgID)
	} else {
		values.Set("msgID", msgID)
	}
	// Check operator call syntax.
	if opCall := decodeOutpostArg(args[3]); !callSignRE.MatchString(opCall) {
		slog.Error("invalid opcall", "opcall", opCall)
		return fmt.Errorf("invalid operator call sign %q", opCall)
	} else {
		values.Set("opCall", opCall)
	}
	// Store operator name.
	values.Set("opName", decodeOutpostArg(args[4]))
	if len(args) == 8 && (args[5] != "" || args[6] != "") {
		// Check tactical call syntax.
		if tacCall, actCall := decodeOutpostArg(args[5]), decodeOutpostArg(args[7]); tacCall == actCall {
			if !tacCallRE.MatchString(tacCall) {
				slog.Error("invalid taccall", "taccall", tacCall)
				return fmt.Errorf("invalid tactical call sign %q", tacCall)
			} else {
				values.Set("tacCall", tacCall)
				values.Set("tacName", decodeOutpostArg(args[6]))
			}
		}
	}
	if mtype = message.FindType(func(mt message.MType) bool {
		if mt, ok := mt.(form.EditableFormType); ok {
			if mt.AddonName == addon && mt.HTMLName == msgtype && mt.CreateTag() != "" {
				values.Set("formtag", mt.CreateTag())
				return true
			}
		}
		return false
	}); mtype == nil {
		slog.Error("unknown form", "addon", addon, "msgtype", msgtype)
		return errors.Join(fmt.Errorf("unknown form %s/%s", values.Get("addon"), values.Get("msgtype")), err)
	}
	if err = outpostNECommon(values, "/outpost-new"); err != nil {
		return err
	}
	slog.Info("opened browser for new Outpost message", "addon", addon, "msgtype", msgtype, "msgID", msgID)
	return nil
}

func outpostNECommon(values url.Values, path string) (err error) {
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
