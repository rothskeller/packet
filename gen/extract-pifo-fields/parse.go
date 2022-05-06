package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// formDefinition and fieldDefinition are parallel to the same types in the
// xscmsg package, except that fieldDefinition.Validations is []string rather
// than []func.
type formDefinition struct {
	HTML                   string
	Tag                    string
	Name                   string
	Article                string
	Version                string
	OriginNumberField      string
	DestinationNumberField string
	HandlingOrderField     string
	SubjectField           string
	OperatorNameField      string
	OperatorCallField      string
	ActionDateField        string
	ActionTimeField        string
	Fields                 []*fieldDefinition
	Annotations            map[string]string
	Comments               map[string]string
}
type fieldDefinition struct {
	Tag         string
	Values      []string
	Validations []string
	Default     string
}

var filenameRE = regexp.MustCompile(`(form-[^.]*)\.v([0-9.]+?)\.html$`)

func extract(filename string) string {
	fd, err := parseFormDefinition(filename)
	if err != nil {
		log.Fatal(err)
	}
	pkgname := filepath.Base(filepath.Dir(filename))
	varname := pkgname
	if d := varname[len(varname)-1]; d >= '0' && d <= '9' {
		varname += "v"
	}
	match := filenameRE.FindStringSubmatch(filename)
	varname += strings.ReplaceAll(match[2], ".", "")
	fd.HTML = match[1] + ".html"
	fd.emit(filepath.Join(filepath.Dir(filename), varname+".def.go"), pkgname, varname)
	return varname
}

func parseFormDefinition(filename string) (fd *formDefinition, err error) {
	fd = new(formDefinition)
	fd.Annotations = make(map[string]string)
	fd.Comments = make(map[string]string)
	fd.HTML = filename
	if err = fd.parseHTML(filename); err != nil {
		return nil, err
	}
	if fd.Tag == "" {
		return nil, errors.New("missing form tag")
	}
	if fd.Name == "" {
		return nil, errors.New("missing form name")
	}
	if fd.Article == "" {
		return nil, errors.New("missing form article")
	}
	if fd.Version == "" {
		return nil, errors.New("missing form version")
	}
	if fd.OriginNumberField == "" {
		return nil, errors.New("missing form OriginNumber field")
	} else if !fd.hasField(fd.OriginNumberField) {
		return nil, fmt.Errorf("form OriginNumber field %q not found in form", fd.OriginNumberField)
	}
	// Note, destination number field is not required since PackItForms HTML
	// forms generally don't have a field name for it.
	if fd.DestinationNumberField != "" && !fd.hasField(fd.DestinationNumberField) {
		return nil, fmt.Errorf("form DestinationNumber field %q not found in form", fd.DestinationNumberField)
	}
	for _, ff := range fd.Fields {
		if len(ff.Values) == 3 && ff.Values[0] == "IMMEDIATE" && ff.Values[1] == "PRIORITY" && ff.Values[2] == "ROUTINE" {
			fd.HandlingOrderField = ff.Tag
			break
		}
	}
	if fd.HandlingOrderField == "" {
		return nil, errors.New("missing form HandlingOrder field")
	}
	if fd.SubjectField == "" {
		return nil, errors.New("missing form Subject field")
	} else if !fd.hasField(fd.SubjectField) {
		return nil, fmt.Errorf("form Subject field %q not found in form", fd.SubjectField)
	}
	if fd.OperatorNameField == "" {
		return nil, errors.New("missing form OperatorName field")
	} else if !fd.hasField(fd.OperatorNameField) {
		return nil, fmt.Errorf("form OperatorName field %q not found in form", fd.OperatorNameField)
	}
	if fd.OperatorCallField == "" {
		return nil, errors.New("missing form OperatorCall field")
	} else if !fd.hasField(fd.OperatorCallField) {
		return nil, fmt.Errorf("form OperatorCall field %q not found in form", fd.OperatorCallField)
	}
	if fd.ActionDateField == "" {
		return nil, errors.New("missing form ActionDate field")
	} else if !fd.hasField(fd.ActionDateField) {
		return nil, fmt.Errorf("form ActionDate field %q not found in form", fd.ActionDateField)
	}
	if fd.ActionTimeField == "" {
		return nil, errors.New("missing form ActionTime field")
	} else if !fd.hasField(fd.ActionTimeField) {
		return nil, fmt.Errorf("form ActionTime field %q not found in form", fd.ActionTimeField)
	}
	if len(fd.Annotations) == 0 {
		fd.Annotations = nil
	}
	for _, field := range fd.Fields {
		var comment = strings.Join(field.Validations, " ")
		if len(field.Values) != 0 {
			if comment != "" {
				comment += ": "
			}
			comment += strings.Join(field.Values, ", ")
			field.Validations = append(field.Validations, "select")
		}
		if comment != "" {
			fd.Comments[field.Tag] = comment
		}
	}
	if len(fd.Comments) == 0 {
		fd.Comments = nil
	}
	return fd, nil
}

func (fd *formDefinition) parseHTML(filename string) (err error) {
	var (
		fh   *os.File
		root *html.Node
	)
	if fh, err = os.Open(filename); err != nil {
		return err
	}
	defer fh.Close()
	if root, err = html.Parse(fh); err != nil {
		return fmt.Errorf("%s: %s", filename, err)
	}
	return fd.parseNode(root)
}

func (fd *formDefinition) parseNode(node *html.Node) (err error) {
	if node.Type == html.ElementNode {
		switch node.DataAtom {
		case atom.Meta:
			err = fd.parseMeta(node)
		case atom.Title:
			err = fd.parseTitle(node)
		case atom.Div:
			err = fd.parseDiv(node)
		case atom.Input:
			err = fd.parseInput(node)
		case atom.Select:
			err = fd.parseSelect(node)
		case atom.Textarea:
			err = fd.parseTextArea(node)
		}
		if err != nil {
			return err
		}
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if err = fd.parseNode(child); err != nil {
			return err
		}
	}
	return nil
}

func (fd *formDefinition) parseMeta(node *html.Node) (err error) {
	var name, content string
	for _, a := range node.Attr {
		switch a.Key {
		case "name":
			name = a.Val
		case "content":
			content = a.Val
		}
	}
	if name == "pack-it-forms-subject-suffix" {
		re := regexp.MustCompile(`^_([A-Z][-A-Za-z0-9]*)_\{\{field:(.*)\}\}$`)
		if match := re.FindStringSubmatch(content); match != nil {
			fd.Tag = match[1]
			if dot := strings.LastIndexByte(match[2], '.'); dot >= 0 {
				fd.SubjectField = match[2][:dot+1]
			} else {
				fd.SubjectField = match[2]
			}
		} else {
			return errors.New("unexpected format for meta pack-it-form-subject-suffix")
		}
	}
	return nil
}

var proseCaseRE = regexp.MustCompile(`\b[A-Z][a-z]*\b`)

func (fd *formDefinition) parseTitle(node *html.Node) (err error) {
	if node.FirstChild == nil || node.FirstChild.Type != html.TextNode || node.FirstChild.NextSibling != nil {
		return errors.New("unexpected format for title")
	}
	fd.Name = node.FirstChild.Data
	fd.Name = proseCaseRE.ReplaceAllStringFunc(fd.Name, strings.ToLower)
	if idx := strings.IndexByte("aeioAEFHILMNORSX", fd.Name[0]); idx >= 0 {
		fd.Article = "an"
	} else {
		fd.Article = "a"
	}
	return nil
}

func (fd *formDefinition) parseDiv(node *html.Node) (err error) {
	for _, a := range node.Attr {
		if a.Key == "class" && a.Val == "version" && fd.Version == "" {
			if node.FirstChild == nil || node.FirstChild.Type != html.TextNode || node.FirstChild.NextSibling != nil {
				return errors.New("unexpected format for div version")
			}
			fd.Version = node.FirstChild.Data
		}
		if a.Key == "data-include-html" {
			fname := filepath.Join(filepath.Dir(fd.HTML), "resources", "html", a.Val+".html")
			if err = fd.parseHTML(fname); err != nil {
				return err
			}
		}
	}
	return nil
}

func (fd *formDefinition) parseInput(node *html.Node) (err error) {
	var field *fieldDefinition
	var itype, name, note, class, value string
	var required, checked, added bool
	for _, a := range node.Attr {
		switch a.Key {
		case "id":
			if a.Val == "form-subject" {
				return nil
			}
		case "type":
			itype = a.Val
		case "name":
			name = a.Val
		case "class":
			class = a.Val
		case "value":
			value = a.Val
		case "required":
			required = true
		case "checked":
			checked = true
		}
	}
	if name == "" {
		return nil
	}
	if dot := strings.LastIndexByte(name, '.'); dot >= 0 {
		name, note = name[:dot+1], name[dot+1:]
	}
	for _, existingField := range fd.Fields {
		if existingField.Tag == name {
			if itype == "radio" {
				field, added = existingField, true
			} else {
				return fmt.Errorf("multiple inputs with name %q", name)
			}
		}
	}
	if field == nil {
		field = new(fieldDefinition)
		field.Tag = name
		if note != "" {
			fd.Annotations[field.Tag] = note
		}
		if required {
			field.Validations = append(field.Validations, "required")
		}
		for _, classpart := range strings.Fields(class) {
			switch classpart {
			case "hidden":
				return nil
			case "time", "date", "call-sign", "message-number", "required-for-complete", "phone-number", "cardinal-number", "real-number", "frequency", "frequency-offset":
				field.Validations = append(field.Validations, classpart)
			case "no-msg-init", "init-on-submit", "include-in-submit", "clearable":
				break
			default:
				return fmt.Errorf("don't know how to handle input class=%q", classpart)
			}
		}
	}
	switch itype {
	case "", "text":
		if strings.HasPrefix(value, "{{") {
			switch value {
			case "{{envelope:sender_message_number}}",
				"{{envelope:sender_message_number|view-by:receiver}}":
				fd.OriginNumberField = field.Tag
			case "{{envelope:receiver_message_number}}",
				"{{envelope:receiver_message_number|view-by:sender}}":
				fd.DestinationNumberField = field.Tag
			case "{{envelope:viewer_operator_call_sign}}":
				fd.OperatorCallField = field.Tag
			case "{{envelope:viewer_operator_name}}":
				fd.OperatorNameField = field.Tag
			case "{{expand-while-null:{{envelope:viewer_date}},{{open-tmpl}}date{{close-tmpl}}}}":
				fd.ActionDateField = field.Tag
			case "{{expand-while-null:{{envelope:viewer_time}},{{open-tmpl}}time{{close-tmpl}}}}":
				fd.ActionTimeField = field.Tag
			case "{{date}}":
				field.Default = "«date»"
			case "{{time}}", "{{envelope:viewer_message_number}}":
				break
			default:
				return fmt.Errorf("don't know how to handle value %q", value)
			}
		} else if value != "" {
			field.Default = value
		}
	case "radio":
		field.Values = append(field.Values, value)
		if checked {
			field.Default = value
		}
	case "checkbox":
		field.Validations = append(field.Validations, "boolean")
	case "button":
		return nil
	default:
		return fmt.Errorf("don't know how to handle input type=%q", itype)
	}
	if !added {
		fd.Fields = append(fd.Fields, field)
	}
	return nil
}

func (fd *formDefinition) parseSelect(node *html.Node) (err error) {
	var field fieldDefinition
	var name, class string
	var required bool
	for _, a := range node.Attr {
		switch a.Key {
		case "name":
			name = a.Val
		case "class":
			class = a.Val
		case "required":
			required = true
		}
	}
	if name == "" {
		return nil
	}
	if dot := strings.LastIndexByte(name, '.'); dot >= 0 {
		field.Tag = name[:dot+1]
		fd.Annotations[field.Tag] = name[dot+1:]
	} else {
		field.Tag = name
	}
	for _, existing := range fd.Fields {
		if existing.Tag == field.Tag {
			return fmt.Errorf("multiple inputs with name %q", field.Tag)
		}
	}
	if required {
		field.Validations = append(field.Validations, "required")
	}
	for _, classpart := range strings.Fields(class) {
		switch classpart {
		case "hidden":
			return nil
		case "required-for-complete":
			field.Validations = append(field.Validations, classpart)
		default:
			return fmt.Errorf("don't know how to handle select class=%q", classpart)
		}
	}
OPTION:
	for option := node.FirstChild; option != nil; option = option.NextSibling {
		if option.Type == html.CommentNode || (option.Type == html.TextNode && strings.TrimSpace(option.Data) == "") {
			continue
		}
		if option.Type != html.ElementNode || option.DataAtom != atom.Option {
			return errors.New("select has non-option child")
		}
		for _, a := range option.Attr {
			switch a.Key {
			case "disabled":
				continue OPTION
			case "value":
				field.Values = append(field.Values, a.Val)
			}
		}
	}
	fd.Fields = append(fd.Fields, &field)
	return nil
}

func (fd *formDefinition) parseTextArea(node *html.Node) (err error) {
	var field fieldDefinition
	var name, class string
	var required bool
	for _, a := range node.Attr {
		switch a.Key {
		case "name":
			name = a.Val
		case "class":
			class = a.Val
		case "required":
			required = true
		}
	}
	if name == "" {
		return nil
	}
	if dot := strings.LastIndexByte(name, '.'); dot >= 0 {
		field.Tag = name[:dot+1]
		fd.Annotations[field.Tag] = name[dot+1:]
	} else {
		field.Tag = name
	}
	for _, existing := range fd.Fields {
		if existing.Tag == field.Tag {
			return fmt.Errorf("multiple inputs with name %q", field.Tag)
		}
	}
	if required {
		field.Validations = append(field.Validations, "required")
	}
	for _, classpart := range strings.Fields(class) {
		switch classpart {
		case "hidden":
			return nil
		default:
			return fmt.Errorf("don't know how to handle textarea class=%q", classpart)
		}
	}
	fd.Fields = append(fd.Fields, &field)
	return nil
}

func (fd *formDefinition) hasField(tag string) bool {
	for _, f := range fd.Fields {
		if f.Tag == tag {
			return true
		}
	}
	return false
}
