package formdef

import (
	"bufio"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"unicode"

	"github.com/rothskeller/pdf/v2"
)

var markupColor = []byte{0, 0, 153, 255}

func Read(filename string) (form *FormDef, err error) {
	if filepath.IsAbs(filename) {
		if filename, err = filepath.Localize(filename[1:]); err != nil {
			return nil, err
		}
		return ReadFS(os.DirFS(string(filepath.Separator)), filename[1:])
	} else {
		if filename, err = filepath.Localize(filename); err != nil {
			return nil, err
		}
		return ReadFS(os.DirFS("."), filename)
	}
}

func ReadFS(formsFS fs.FS, filename string) (form *FormDef, err error) {
	var (
		fh      fs.File
		in      *bufio.Scanner
		fields  []string
		line    string
		linenum int
		fd      *FieldDef
	)
	if fh, err = formsFS.Open(filename); err != nil {
		return nil, err
	}
	defer fh.Close()
	in = bufio.NewScanner(fh)
	form = &FormDef{FormFS: formsFS}
TOPLEVEL:
	for in.Scan() {
		linenum++
		line = in.Text()
		if fields = tokenizeLine(line); len(fields) == 0 {
			line = ""
			continue
		}
		switch fields[0] {
		case "addon":
			form.AddonName = strings.Join(fields[1:], " ")
		case "html":
			form.HTMLName = strings.Join(fields[1:], " ")
			if form.HTMLFile == "" {
				form.HTMLFile = form.HTMLName
			}
		case "version":
			form.Version = strings.Join(fields[1:], " ")
		case "title":
			form.Title = strings.Join(fields[1:], " ")
			if form.Sort == "" {
				form.Sort = form.Title
			}
		case "sort":
			form.Sort = strings.Join(fields[1:], " ")
		case "indef":
			form.IndefName = strings.Join(fields[1:], " ")
		case "create":
			form.CreateTags = fields[1:]
			if len(form.CreateTags) != 0 && form.SubjectTag == "" {
				form.SubjectTag = form.CreateTags[0]
			}
		case "subject":
			form.SubjectTag = strings.Join(fields[1:], "_")
		case "htmlFile":
			form.HTMLFile = strings.Join(fields[1:], " ")
		case "pdfFile":
			form.PDFFile = strings.Join(fields[1:], " ")
		case "field":
			break TOPLEVEL
		default:
			return nil, fmt.Errorf("%s:%d: unknown keyword %q", filename, linenum, fields[0])
		}
		line = ""
	}
	if form.AddonName == "" {
		return nil, fmt.Errorf("%s: no 'addon' specified", filename)
	}
	if form.HTMLName == "" {
		return nil, fmt.Errorf("%s: no 'html' specified", filename)
	}
	if form.Version == "" {
		return nil, fmt.Errorf("%s: no 'version' specified", filename)
	}
	if form.Title == "" {
		return nil, fmt.Errorf("%s: no 'title' specified", filename)
	}
	if form.IndefName == "" {
		return nil, fmt.Errorf("%s: no 'indefName' specified", filename)
	}
	if form.SubjectTag == "" {
		return nil, fmt.Errorf("%s: no 'subject' (or 'create') specified", filename)
	}
	form.HTMLFile = path.Join(path.Dir(filename), form.HTMLFile)
	if form.PDFFile != "" {
		form.PDFFile = path.Join(path.Dir(filename), form.PDFFile)
	}
	for line != "" {
		fd = &FieldDef{Type: "text"}
		form.Fields = append(form.Fields, fd)
		switch len(fields) {
		case 3:
			if CommonTags.Has(fields[2]) {
				fd.Common = fields[2]
			} else {
				return nil, fmt.Errorf("%s:%d: second tag %q is not one of the common tags", filename, linenum, fields[2])
			}
			fallthrough
		case 2:
			if fields[1] != "-" {
				fd.Tag = fields[1]
			}
		default:
			return nil, fmt.Errorf("%s:%d: field line needs one or two arguments", filename, linenum)
		}
	FIELD:
		for in.Scan() {
			linenum++
			line = in.Text()
			if fields = tokenizeLine(line); len(fields) == 0 {
				line = ""
				continue
			}
			switch fields[0] {
			case "label":
				fd.Label = strings.Join(fields[1:], " ")
			case "clabel":
				fd.ChildLabel = strings.Join(fields[1:], " ")
			case "type":
				if len(fields) == 2 {
					if _, ok := typeHandlers[fields[1]]; ok {
						fd.Type = fields[1]
					} else {
						return nil, fmt.Errorf("%s:%d: unknown field type %q", filename, linenum, fields[1])
					}
				} else {
					return nil, fmt.Errorf("%s:%d: type needs a single argument", filename, linenum)
				}
			case "value":
				fd.Value = strings.Join(fields[1:], " ")
			case "choice":
				c := strings.Join(fields[1:], " ")
				fd.Choices = append(fd.Choices, Choice{Raw: c, Human: c})
			case "mchoice":
				if len(fields) < 3 {
					return nil, fmt.Errorf("%s:%d: mchoice needs two arguments", filename, linenum)
				}
				fd.Choices = append(fd.Choices, Choice{Raw: fields[1], Human: strings.Join(fields[2:], " ")})
			case "cchoice":
				then := slices.Index(fields, "then")
				if len(fields) < 6 || fields[1] != "if" || then < 3 || len(fields) < then+2 {
					return nil, fmt.Errorf("%s:%d: cchoice syntax error", filename, linenum)
				}
				c := Choice{CondField: strings.Join(fields[2:then], " "), Raw: fields[then+1], Human: strings.Join(fields[then+2:], " ")}
				c.CondField, c.CondValue, _ = strings.Cut(c.CondField, "=")
				c.CondValue = strings.Trim(c.CondValue, `"`)
				fd.Choices = append(fd.Choices, c)
			case "presence":
				if err = parsePresence(fd, fields[1:]); err != nil {
					return nil, fmt.Errorf("%s:%d: %s", filename, linenum, err)
				}
			case "width":
				if len(fields) != 2 {
					return nil, fmt.Errorf("%s:%d: invalid width value", filename, linenum)
				} else if fd.EditWidth, err = strconv.Atoi(fields[1]); err != nil || fd.EditWidth < 1 {
					return nil, fmt.Errorf("%s:%d: invalid width value", filename, linenum)
				}
			case "help":
				fd.EditHelp = strings.Join(fields[1:], " ")
			case "compare":
				fd.CompareMethod = strings.Join(fields[1:], " ")
			case "pdf":
				if err = parsePDFRender(fd, fields[1:]); err != nil {
					return nil, fmt.Errorf("%s:%d: %s", filename, linenum, err)
				}
			case "children":
				for _, ctag := range fields[1:] {
					idx := slices.IndexFunc(form.Fields, func(def *FieldDef) bool { return def.Tag == ctag })
					if idx < 0 {
						return nil, fmt.Errorf("%s:%d: no such field %q", filename, linenum, ctag)
					}
					child := form.Fields[idx]
					form.Fields = slices.Delete(form.Fields, idx, idx+1)
					child.Parent = fd
					fd.Children = append(fd.Children, child)
				}
			case "field":
				break FIELD
			default:
				return nil, fmt.Errorf("%s:%d: unknown keyword %q", filename, linenum, fields[0])
			}
			line = ""
		}
	}
	if err = in.Err(); err != nil {
		return nil, fmt.Errorf("close %s: %s", filename, err)
	}
	return form, nil
}

func parsePresence(fd *FieldDef, tokens []string) (err error) {
	fd.Presence = Optional // by default
	for len(tokens) >= 4 && tokens[0] == "if" && tokens[2] == "then" {
		pc := &PresenceCond{OtherField: tokens[1], Presence: Presence(tokens[3])}
		pc.OtherField, pc.OtherValue, _ = strings.Cut(pc.OtherField, "=")
		switch pc.Presence {
		case "required", "optional", "blocked":
			// ok
		default:
			return fmt.Errorf("invalid presence keyword %q", pc.Presence)
		}
		fd.PresenceCond = append(fd.PresenceCond, pc)
		tokens = tokens[4:] // remove conditional
		if len(tokens) != 0 {
			if len(tokens) < 2 || tokens[0] != "else" {
				return errors.New("invalid presence syntax")
			}
			tokens = tokens[1:] // remove "else"
		}
	}
	switch len(tokens) {
	case 0:
		fd.Presence = Optional
	case 1:
		switch tokens[0] {
		case "required", "optional", "blocked":
			fd.Presence = Presence(tokens[0])
		default:
			return fmt.Errorf("invalid presence keyword %q", tokens[0])
		}
	default:
		return errors.New("invalid presence syntax")
	}
	return nil
}

func parsePDFRender(fd *FieldDef, words []string) (err error) {
	var (
		pr    PDFFieldRenderer
		attrs map[string]string
	)
	for len(words) > 3 && words[0] == "if" && words[2] == "then" {
		var c PDFCondition
		c.Field, c.Value, c.Set = strings.Cut(words[1], "=")
		c.Set = !c.Set
		pr.Conditions = append(pr.Conditions, c)
		words = words[3:]
	}
	if len(words) == 0 {
		return errors.New("missing shape name")
	}
	if attrs, err = parsePDFAttrs(words[1:]); err != nil {
		return err
	}
	switch words[0] {
	case "box":
		pr.Renderer, err = parseBoxParams(attrs)
	case "circle":
		pr.Renderer, err = parseCircleParams(attrs)
	case "cross":
		pr.Renderer, err = parseCrossParams(attrs)
	case "text":
		pr.Renderer, err = parseTextParams(attrs)
	default:
		err = fmt.Errorf("unknown shape name %q", words[0])
	}
	if err != nil {
		return err
	}
	fd.PDF = append(fd.PDF, pr)
	return nil
}

func parsePDFAttrs(words []string) (attrs map[string]string, err error) {
	if len(words)%2 == 1 {
		return nil, errors.New("incomplete attribute pair")
	}
	attrs = make(map[string]string)
	for i := 0; i < len(words); i += 2 {
		if _, ok := attrs[words[i]]; ok {
			return nil, fmt.Errorf("duplicate attribute %q", words[i])
		}
		attrs[words[i]] = words[i+1]
	}
	return attrs, nil
}

func parseBoxParams(attrs map[string]string) (r BoxRenderer, err error) {
	var box pdf.Box

	if box.Page, err = getPageAttr(attrs); err != nil {
		return r, err
	}
	if box.Rectangle, err = getRectangleAttrs(attrs); err != nil {
		return r, err
	}
	if box.Fill, err = getColorAttr(attrs, "F", nil); err != nil {
		return r, err
	}
	if box.Stroke, err = getColorAttr(attrs, "S", nil); err != nil {
		return r, err
	}
	if box.Fill == nil && box.Stroke == nil {
		return r, errors.New("box must have F or S")
	}
	if s, ok := attrs["SW"]; ok {
		if box.Stroke == nil {
			return r, errors.New("box SW invalid without S")
		}
		if v, err := strconv.ParseFloat(s, 64); err != nil || v < 0 {
			return r, errors.New("invalid SW value")
		} else {
			delete(attrs, "SW")
			box.StrokeWidth = v
		}
	}
	if len(attrs) != 0 {
		return r, errors.New("excess attributes")
	}
	r.Box = &box
	return r, nil
}

func parseCircleParams(attrs map[string]string) (r CircleRenderer, err error) {
	var circle pdf.Circle

	if circle.Page, err = getPageAttr(attrs); err != nil {
		return r, err
	}
	if circle.Fill, err = getColorAttr(attrs, "C", markupColor); err != nil {
		return r, err
	}
	if s, ok := attrs["X"]; ok {
		if v, err := strconv.ParseFloat(s, 64); err != nil {
			return r, errors.New("non-numeric X value")
		} else {
			delete(attrs, "X")
			circle.Center.X = v
		}
	} else {
		return r, errors.New("no X value")
	}
	if s, ok := attrs["Y"]; ok {
		if v, err := strconv.ParseFloat(s, 64); err != nil {
			return r, errors.New("non-numeric Y value")
		} else {
			delete(attrs, "Y")
			circle.Center.Y = v
		}
	} else {
		return r, errors.New("no Y value")
	}
	if s, ok := attrs["R"]; ok {
		if v, err := strconv.ParseFloat(s, 64); err != nil || v < 0 {
			return r, errors.New("invalid R value")
		} else {
			delete(attrs, "R")
			circle.Radius = v
		}
	} else {
		return r, errors.New("no R value")
	}
	if len(attrs) != 0 {
		return r, errors.New("excess attributes")
	}
	r.Circle = &circle
	return r, nil
}

func parseCrossParams(attrs map[string]string) (r CrossRenderer, err error) {
	var cross pdf.Cross

	if cross.Rectangle, err = getRectangleAttrs(attrs); err != nil {
		return r, err
	}
	cross.LineWidth = (cross.Rectangle.URX - cross.Rectangle.LLX) / 10
	if cross.Page, err = getPageAttr(attrs); err != nil {
		return r, err
	}
	if cross.Stroke, err = getColorAttr(attrs, "C", markupColor); err != nil {
		return r, err
	}
	if len(attrs) != 0 {
		return r, errors.New("excess attributes")
	}
	r.Cross = &cross
	return r, nil
}

func parseTextParams(attrs map[string]string) (r TextRenderer, err error) {
	var text pdf.Text

	if text.Rectangle, err = getRectangleAttrs(attrs); err != nil {
		return r, err
	}
	if text.Page, err = getPageAttr(attrs); err != nil {
		return r, err
	}
	if text.Color, err = getColorAttr(attrs, "C", markupColor); err != nil {
		return r, err
	}
	text.Color = text.Color[:3] // remove alpha
	if s, ok := attrs["S"]; ok {
		delete(attrs, "S")
		text.String = s
	}
	if s, ok := attrs["BL"]; ok {
		if v, err := strconv.ParseFloat(s, 64); err != nil || v < text.Rectangle.LLY || v >= text.Rectangle.URY {
			return r, errors.New("invalid BL value")
		} else {
			delete(attrs, "BL")
			text.Baseline = v
		}
	}
	if s, ok := attrs["FT"]; ok {
		delete(attrs, "FT")
		text.Font = s
	} else {
		text.Font = "Times-Roman"
	}
	if s, ok := attrs["FS"]; ok {
		if v, err := strconv.ParseFloat(s, 64); err != nil || v <= 0 {
			return r, errors.New("invalid FS value")
		} else {
			delete(attrs, "FS")
			text.FontSize = v
		}
	} else {
		text.FontSize = 12
	}
	if s, ok := attrs["MS"]; ok {
		if v, err := strconv.ParseFloat(s, 64); err != nil || v <= 0 || v > text.FontSize {
			return r, errors.New("invalid MS value")
		} else {
			delete(attrs, "MS")
			text.MinFontSize = v
		}
	} else {
		text.MinFontSize = 8
	}
	if s, ok := attrs["LH"]; ok {
		if v, err := strconv.ParseFloat(s, 64); err != nil || v <= 0 {
			return r, errors.New("invalid LH value")
		} else {
			delete(attrs, "LH")
			text.LineHeight = v
		}
	} else {
		text.LineHeight = 1.15
	}
	if s, ok := attrs["A"]; ok {
		delete(attrs, "A")
		text.Align = s
		if _, ok := attrs["HA"]; ok {
			return r, errors.New("cannot specify both A and HA")
		}
		if _, ok := attrs["VA"]; ok {
			return r, errors.New("cannot specify both A and VA")
		}
	} else {
		if s, ok := attrs["HA"]; ok {
			delete(attrs, "HA")
			switch s {
			case "l", "left":
				text.Align = "l"
			case "c", "center":
				text.Align = "c"
			case "r", "right":
				text.Align = "r"
			default:
				return r, errors.New("invalid HA value")
			}
		} else {
			text.Align = "l"
		}
		if s, ok := attrs["VA"]; ok {
			delete(attrs, "VA")
			switch s {
			case "t", "top":
				text.Align += "t"
			case "c", "center":
				text.Align += "m"
			case "b", "bottom":
				text.Align += "b"
			case "baseline":
				text.Align += "F"
			default:
				return r, errors.New("invalid VA value")
			}
		} else {
			text.Align += "F"
		}
	}
	if s, ok := attrs["WR"]; ok {
		delete(attrs, "WR")
		switch s {
		case "t", "true":
			text.Wrap = true
		case "f", "false":
			text.Wrap = false
		default:
			return r, errors.New("invalid WR value")
		}
	} else {
		text.Wrap = true
	}
	if s, ok := attrs["CL"]; ok {
		delete(attrs, "CL")
		switch s {
		case "t", "true":
			text.Clip = true
		case "f", "false":
			text.Clip = false
		default:
			return r, errors.New("invalid CL value")
		}
	}
	if len(attrs) != 0 {
		return r, errors.New("excess attributes")
	}
	r.Text = &text
	return r, nil
}

func getRectangleAttrs(attrs map[string]string) (rect pdf.Rectangle, err error) {
	var x, y, w, h, t, b, l, r float64
	var hx, hy, hw, hh, ht, hb, hl, hr bool

	if s, ok := attrs["X"]; ok {
		if v, err := strconv.ParseFloat(s, 64); err != nil {
			return rect, errors.New("non-numeric X value")
		} else {
			delete(attrs, "X")
			x, hx = v, true
		}
	}
	if s, ok := attrs["Y"]; ok {
		if v, err := strconv.ParseFloat(s, 64); err != nil {
			return rect, errors.New("non-numeric Y value")
		} else {
			delete(attrs, "Y")
			y, hy = v, true
		}
	}
	if s, ok := attrs["W"]; ok {
		if v, err := strconv.ParseFloat(s, 64); err != nil {
			return rect, errors.New("non-numeric W value")
		} else {
			delete(attrs, "W")
			w, hw = v, true
		}
	}
	if s, ok := attrs["H"]; ok {
		if v, err := strconv.ParseFloat(s, 64); err != nil {
			return rect, errors.New("non-numeric H value")
		} else {
			delete(attrs, "H")
			h, hh = v, true
		}
	}
	if s, ok := attrs["T"]; ok {
		if v, err := strconv.ParseFloat(s, 64); err != nil {
			return rect, errors.New("non-numeric T value")
		} else {
			delete(attrs, "T")
			t, ht = v, true
		}
	}
	if s, ok := attrs["B"]; ok {
		if v, err := strconv.ParseFloat(s, 64); err != nil {
			return rect, errors.New("non-numeric B value")
		} else {
			delete(attrs, "B")
			b, hb = v, true
		}
	}
	if s, ok := attrs["L"]; ok {
		if v, err := strconv.ParseFloat(s, 64); err != nil {
			return rect, errors.New("non-numeric L value")
		} else {
			delete(attrs, "L")
			l, hl = v, true
		}
	}
	if s, ok := attrs["R"]; ok {
		if v, err := strconv.ParseFloat(s, 64); err != nil {
			return rect, errors.New("non-numeric R value")
		} else {
			delete(attrs, "R")
			r, hr = v, true
		}
	}
	if hx && hl && x != l {
		return rect, errors.New("X!=L")
	} else if hl {
		x, hx = l, true
	}
	if hy && hb && y != b {
		return rect, errors.New("X!=L")
	} else if hb {
		y, hy = b, true
	}
	if !hx {
		if hw && hr {
			x, hx, hw = r-w, true, false
		} else {
			return rect, errors.New("incomplete rect (X)")
		}
	}
	if !hy {
		if hh && ht {
			y, hy, hh = t-h, true, false
		} else {
			return rect, errors.New("incomplete rect (Y)")
		}
	}
	if hw && hr {
		if x+w != r {
			return rect, errors.New("inconsistent X/W/R")
		}
	} else if hw {
		r = x + w
	} else if !hr {
		return rect, errors.New("incomplete rect (R)")
	}
	if hh && ht {
		if y+h != t {
			return rect, errors.New("inconsistent Y/H/T")
		}
	} else if hh {
		t = y + h
	} else if !ht {
		return rect, errors.New("incomplete rect (T)")
	}
	return pdf.Rectangle{LLX: x, LLY: y, URX: r, URY: t}, nil
}

func getPageAttr(attrs map[string]string) (page int, err error) {
	if s, ok := attrs["P"]; ok {
		if v, err := strconv.Atoi(s); err != nil || v < 1 {
			return 0, errors.New("invalid P")
		} else {
			delete(attrs, "P")
			return v, nil
		}
	}
	return 1, nil
}

func getColorAttr(attrs map[string]string, key string, def []byte) (color []byte, err error) {
	if s, ok := attrs[key]; ok {
		if len(s) != 6 && len(s) != 8 {
			return nil, errors.New("invalid " + key)
		} else if color, err = hex.DecodeString(s); err != nil {
			return nil, errors.New("invalid " + key)
		} else {
			delete(attrs, key)
			return color, nil
		}
	}
	return def, nil
}

// tokenizeLine returns the set of unquoted-whitespace-delimited tokens on the
// line.  Quotes are removed.
func tokenizeLine(line string) (tokens []string) {
	var (
		curtoken string
		escaped  bool
		quoted   bool
	)
LOOP:
	for _, r := range line {
		switch {
		case escaped:
			curtoken += string(r)
			escaped = false
		case r == '\\':
			escaped = true
		case r == '"':
			quoted = !quoted
		case unicode.IsSpace(r) && !quoted:
			if curtoken != "" {
				tokens = append(tokens, curtoken)
				curtoken = ""
			}
		case r == '#' && !quoted:
			break LOOP
		default:
			curtoken += string(r)
		}
	}
	if curtoken != "" {
		tokens = append(tokens, curtoken)
	}
	return tokens
}
