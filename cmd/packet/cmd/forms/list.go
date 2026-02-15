package forms

import (
	"cmp"
	"fmt"
	"slices"

	"github.com/rothskeller/packet/form"
	"github.com/rothskeller/packet/form/formdef"
	"github.com/rothskeller/packet/form/formdefs"
	"github.com/rothskeller/packet/message"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists installed forms",
	Long: `Prints a list of the installed forms.  By default, only forms that can be
created and edited are listed; to list other forms that can be received, add a
--all (-a) flag.  The --verbose (-v) flag adds information about how the forms
are identified in PackItForms encoding.`,
	Args:         cobra.NoArgs,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var list []*formdef.FormDef
		all, _ := cmd.Flags().GetBool("all")
		verbose, _ := cmd.Flags().GetBool("verbose")
		if err := formdefs.RegisterForms(); err != nil {
			return err
		}
		for t := range message.AllTypes() {
			switch t := t.(type) {
			case form.FormType:
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
		fmt.Println("* TAG          CT  VER  DESCRIPTION")
		if verbose {
			fmt.Println("                        PACKITFORMS")
		}
		for _, fd := range list {
			var creatable, tag string
			if fd.CreateTag != "" {
				creatable = "*"
				tag = fd.CreateTag
			} else {
				tag = fd.SubjectTag
			}
			fmt.Printf("%-1.1s %-12.12s %-2.2s  %-4.4s %s\n", creatable, tag, fd.CreateKey, fd.Version, fd.Title)
			if verbose {
				fmt.Printf("                        %-12.12s #T: %s\n", "!"+fd.AddonName+"!", fd.HTMLName)
			}
		}
		return nil
	},
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

func init() {
	listCmd.Flags().BoolP("all", "a", false, "Include old and non-editable versions")
	listCmd.Flags().BoolP("verbose", "v", false, "Include Outpost form identifiers")
	Command.AddCommand(listCmd)
}
