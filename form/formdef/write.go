package formdef

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func Write(filename string, form *FormDef) (err error) {
	var fh *os.File

	if fh, err = os.Create(filename); err != nil {
		return err
	}
	defer fh.Close()
	fmt.Fprintf(fh, "addon    %s\n", form.AddonName)
	fmt.Fprintf(fh, "html     %s\n", form.HTMLName)
	fmt.Fprintf(fh, "version  %s\n", form.Version)
	fmt.Fprintf(fh, "title    %s\n", form.Title)
	if form.Sort != "" && form.Sort != form.Title {
		fmt.Fprintf(fh, "sort     %s\n", form.Title)
	}
	fmt.Fprintf(fh, "indef    %s\n", form.IndefName)
	if len(form.CreateTags) != 0 {
		fmt.Fprintf(fh, "create   %s\n", strings.Join(form.CreateTags, " "))
	}
	if form.SubjectTag != "" && (len(form.CreateTags) == 0 || form.SubjectTag != form.CreateTags[0]) {
		fmt.Fprintf(fh, "subject  %s\n", form.SubjectTag)
	}
	if form.HTMLFile != "" && filepath.Base(form.HTMLFile) != form.HTMLName {
		fmt.Fprintf(fh, "htmlFile %s\n", filepath.Base(form.HTMLFile))
	}
	if form.PDFFile != "" {
		fmt.Fprintf(fh, "pdfFile  %s\n", filepath.Base(form.PDFFile))
	}
	for _, fd := range form.Fields {
		emitField(fh, fd)
	}
	return nil
}

func emitField(fh *os.File, fd *FieldDef) {
	for _, c := range fd.Children {
		emitField(fh, c)
	}
	switch {
	case fd.Tag != "" && fd.Common != "":
		fmt.Fprintf(fh, "\nfield %s %s\n", fd.Tag, fd.Common)
	case fd.Tag != "":
		fmt.Fprintf(fh, "\nfield %s\n", fd.Tag)
	case fd.Common != "":
		fmt.Fprintf(fh, "\nfield - %s\n", fd.Common)
	default:
		fmt.Fprint(fh, "\nfield -\n")
	}
	if fd.Label != "" {
		fmt.Fprintf(fh, "  label    %s\n", maybeQuote(fd.Label))
	}
	if fd.ChildLabel != "" && fd.ChildLabel != fd.Label {
		fmt.Fprintf(fh, "  clabel   %s\n", maybeQuote(fd.ChildLabel))
	}
	if fd.Type != "text" {
		fmt.Fprintf(fh, "  type     %s\n", fd.Type)
	}
	if len(fd.Children) != 0 {
		ctags := make([]string, len(fd.Children))
		for i := range fd.Children {
			ctags[i] = fd.Children[i].Tag
		}
		fmt.Fprintf(fh, "  children %s\n", strings.Join(ctags, " "))
	}
	if fd.Value != "" {
		fmt.Fprintf(fh, "  value    %s\n", maybeQuote(fd.Value))
	}
	for _, c := range fd.Choices {
		if c.CondField != "" && c.CondValue != "" {
			fmt.Fprintf(fh, "  cchoice  if %s=%s then %s %s\n", c.CondField, maybeQuote(c.CondValue), c.Raw, c.Human)
		} else if c.CondField != "" {
			fmt.Fprintf(fh, "  cchoice  if %s then %s %s\n", c.CondField, c.Raw, c.Human)
		} else if c.Raw != c.Human {
			fmt.Fprintf(fh, "  mchoice  %s %s\n", c.Raw, c.Human)
		} else {
			fmt.Fprintf(fh, "  choice   %s\n", c.Raw)
		}
	}
	if (fd.Presence != "" && fd.Presence != Optional) || len(fd.PresenceCond) != 0 {
		fmt.Fprint(fh, "  presence ")
		for i, pc := range fd.PresenceCond {
			if i != 0 {
				fmt.Fprint(fh, " else if ")
			} else {
				fmt.Fprint(fh, "if ")
			}
			fmt.Fprint(fh, pc.OtherField)
			if pc.OtherValue != "" {
				fmt.Fprint(fh, "=")
				fmt.Fprint(fh, maybeQuote(pc.OtherValue))
			}
			fmt.Fprintf(fh, " then %s", pc.Presence)
		}
		if fd.Presence != "" && fd.Presence != Optional {
			if len(fd.PresenceCond) != 0 {
				fmt.Fprint(fh, " else ")
			}
			fmt.Fprint(fh, fd.Presence)
		}
		fmt.Fprintln(fh)
	}
	if fd.EditWidth != 0 {
		fmt.Fprintf(fh, "  width    %d\n", fd.EditWidth)
	}
	if fd.EditHelp != "" {
		fmt.Fprintf(fh, "  help     %q\n", fd.EditHelp)
	}
	if fd.CompareMethod != "" {
		fmt.Fprintf(fh, "  compare  %s\n", fd.CompareMethod)
	}
	for _, pr := range fd.PDF {
		emitPDFRenderer(fh, pr)
	}
}

func emitPDFRenderer(fh *os.File, pr PDFFieldRenderer) {
	fmt.Fprint(fh, "  pdf      ")
	for _, c := range pr.Conditions {
		if c.Set {
			fmt.Fprintf(fh, "if %s then ", c.Field)
		} else {
			fmt.Fprintf(fh, "if %s=%s then ", c.Field, maybeQuote(c.Value))
		}
	}
	switch r := pr.Renderer.(type) {
	case CircleRenderer:
		fmt.Fprint(fh, "circle ")
		if r.Page != 1 {
			fmt.Fprintf(fh, "P %d ", r.Page)
		}
		fmt.Fprintf(fh, "X %6.2f Y %6.2f R %.1f", r.Center.X, r.Center.Y, r.Radius)
		if r.Fill[0] != 0 || r.Fill[1] != 0 || r.Fill[2] != 153 || r.Fill[3] != 255 {
			fmt.Fprintf(fh, " C %02X%02X%02X", r.Fill[0], r.Fill[1], r.Fill[2])
			if r.Fill[3] != 255 {
				fmt.Fprintf(fh, "%02X", r.Fill[3])
			}
		}
	case CrossRenderer:
		fmt.Fprint(fh, "cross ")
		if r.Page != 1 {
			fmt.Fprintf(fh, "P %d ", r.Page)
		}
		fmt.Fprintf(fh, "L %6.2f R %6.2f B %6.2f T %6.2f", r.Rectangle.LLX, r.Rectangle.URX, r.Rectangle.LLY, r.Rectangle.URY)
		if r.Stroke[0] != 0 || r.Stroke[1] != 0 || r.Stroke[2] != 153 || r.Stroke[3] != 255 {
			fmt.Fprintf(fh, " C %02X%02X%02X", r.Stroke[0], r.Stroke[1], r.Stroke[2])
			if r.Stroke[3] != 255 {
				fmt.Fprintf(fh, "%02X", r.Stroke[3])
			}
		}
	case TextRenderer:
		fmt.Fprint(fh, "text ")
		if r.Page != 1 {
			fmt.Fprintf(fh, "P %d ", r.Page)
		}
		fmt.Fprintf(fh, "L %6.2f R %6.2f B %6.2f T %6.2f", r.Rectangle.LLX, r.Rectangle.URX, r.Rectangle.LLY, r.Rectangle.URY)
		if r.Baseline != 0 {
			fmt.Fprintf(fh, " BL %6.2f", r.Baseline)
		}
		if r.Font != "Times-Roman" {
			fmt.Fprintf(fh, " FT %s", r.Font)
		}
		if r.FontSize != 12 {
			fmt.Fprintf(fh, " FS %.1f", r.FontSize)
		}
		if r.MinFontSize != 8 {
			fmt.Fprintf(fh, " FS %.1f", r.MinFontSize)
		}
		if r.LineHeight != 1.15 {
			fmt.Fprintf(fh, " FS %.2f", r.LineHeight)
		}
		if r.Align != "" && r.Align != "lF" {
			fmt.Fprintf(fh, " A %s", r.Align)
		}
		if !r.Wrap {
			fmt.Fprint(fh, " WR f")
		}
		if r.Clip {
			fmt.Fprint(fh, " CL t")
		}
		if r.Color[0] != 0 || r.Color[1] != 0 || r.Color[2] != 153 {
			fmt.Fprintf(fh, " C %02X%02X%02X", r.Color[0], r.Color[1], r.Color[2])
		}
	}
	fmt.Fprintln(fh)

}
func maybeQuote(s string) string {
	if strings.ContainsAny(s, " #") {
		return strconv.Quote(s)
	}
	return s
}
