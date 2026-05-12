package cmd

import (
	"cmp"
	"slices"

	"github.com/rothskeller/packet/cmd/packet/cio"
	"github.com/rothskeller/packet/form"
	"github.com/rothskeller/packet/form/formdef"
	"github.com/rothskeller/packet/form/formdefs"
	"github.com/rothskeller/packet/message"
	"github.com/spf13/pflag"
)

const (
	formsListSlug = `List installed forms`
	formsListHelp = `
usage: packet forms list [flags]
  -a, --all      ⇥include old and non-editable versions
  -v, --verbose  ⇥include Outpost form identifiers

The "packet forms list" command lists installed forms.  Unless the --all (or -a) flag is given, it lists only forms that can be used to create new outgoing messages.

The output is a table with the following columns:
  - ⇥Receive/send indicator.  This column is included only when the --all (or -a) flag is given.  It contains a "*" if the form can be used to create new outgoing messages, and is empty otherwise.  (In the CSV output, it contains "RS" or "R", respectively.)
  - ⇥Form tag.  This is the string used to identify the form type on the "new" command line and, usually, in the subject lines of the resulting messages.
  - ⇥Create key.  This is a one- or two-character string that can also be used to identify the form type on the "new" command line.
  - ⇥Version.  This is the version number of the form.
  - ⇥Description.  This is the title or description of the form.
Plus two additional columns if --verbose (or -v) is given:
  - ⇥Addon name.  The Outpost "addon" name for the form.  This appears surrounded by exclamation points on the first line of the encoded message.
  - ⇥HTML name.  The Outpost identifier for the form.  This appears on the "#T:" line in the encoded message.
`
)

func cmdFormsList(args []string) (err error) {
	var all, verbose bool

	flags := pflag.NewFlagSet("f-list", pflag.ContinueOnError)
	flags.BoolVarP(&all, "all", "a", false, "Include old and non-editable versions")
	flags.BoolVarP(&verbose, "verbose", "v", false, "Include Outpost form identifiers")
	flags.Usage = func() {} // we do our own
	if err = flags.Parse(args); err == pflag.ErrHelp {
		return cmdFormsHelp([]string{"list"})
	} else if err != nil {
		cio.Open().Error(err)
		return usage(formsListHelp)
	}
	if flags.NArg() != 0 {
		return usage(formsListHelp)
	}
	registerForms()
	var list []*formdef.FormDef
	if err := formdefs.RegisterForms(); err != nil {
		return err
	}
	for t := range message.AllTypes() {
		switch t := t.(type) {
		case *form.FormType:
			if !all {
				continue
			}
			list = append(list, t.FormDef)
		case form.EditableFormType:
			list = append(list, t.FormDef)
		default:
			continue
		}
	}
	slices.SortFunc(list, compareFormDefs)
	fl := cio.Open().NewFormsList(all, verbose)
	for _, fd := range list {
		var tag string
		if fd.CreateTag != "" {
			tag = fd.CreateTag
		} else {
			tag = fd.SubjectTag
		}
		fl.ShowForm(fd.CreateTag != "", tag, fd.CreateKey, fd.Version, fd.Title, fd.AddonName, fd.HTMLName)
	}
	fl.Close()
	return nil
}

func compareFormDefs(a, b *formdef.FormDef) int {
	ta, tb := a.CreateTag, b.CreateTag
	if ta == "" {
		ta = a.SubjectTag
	}
	if tb == "" {
		tb = b.SubjectTag
	}
	if ta != tb {
		return cmp.Compare(ta, tb)
	}
	if a.CreateTag != "" && b.CreateTag == "" {
		return -1
	}
	if a.CreateTag == "" && b.CreateTag != "" {
		return +1
	}
	return -cmp.Compare(a.Version, b.Version)
}
