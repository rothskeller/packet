package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func (fd *formDefinition) emit(filename, pkgname, varname string) {
	fh, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer fh.Close()
	fmt.Fprintf(fh, "// Code generated by extract-pifo-fields. DO NOT EDIT.\n\npackage %s\n\nimport \"steve.rothskeller.net/packet/xscmsg/internal/xscform\"\n\nvar %s = &xscform.FormDefinition{\n", pkgname, varname)
	fmt.Fprintf(fh, "\tHTML:                   %q,\n", filepath.Base(fd.HTML))
	fmt.Fprintf(fh, "\tTag:                    %q,\n", fd.Tag)
	fmt.Fprintf(fh, "\tName:                   %q,\n", fd.Name)
	fmt.Fprintf(fh, "\tArticle:                %q,\n", fd.Article)
	fmt.Fprintf(fh, "\tVersion:                %q,\n", fd.Version)
	fmt.Fprintf(fh, "\tOriginNumberField:      %q,\n", fd.OriginNumberField)
	fmt.Fprintf(fh, "\tDestinationNumberField: %q,\n", fd.DestinationNumberField)
	fmt.Fprintf(fh, "\tHandlingOrderField:     %q,\n", fd.HandlingOrderField)
	fmt.Fprintf(fh, "\tSubjectField:           %q,\n", fd.SubjectField)
	fmt.Fprintf(fh, "\tOperatorNameField:      %q,\n", fd.OperatorNameField)
	fmt.Fprintf(fh, "\tOperatorCallField:      %q,\n", fd.OperatorCallField)
	fmt.Fprintf(fh, "\tActionDateField:        %q,\n", fd.ActionDateField)
	fmt.Fprintf(fh, "\tActionTimeField:        %q,\n", fd.ActionTimeField)
	fmt.Fprintf(fh, "\tFields:                 []*xscform.FieldDefinition{\n")
	for _, ff := range fd.Fields {
		ff.emit(fh)
	}
	fmt.Fprintf(fh, "\t},\n")
	if len(fd.Annotations) != 0 {
		var tags []string
		for tag := range fd.Annotations {
			tags = append(tags, tag)
		}
		sort.Strings(tags)
		fmt.Fprintf(fh, "\tAnnotations: map[string]string{\n")
		for _, tag := range tags {
			fmt.Fprintf(fh, "\t\t%q: %q,\n", tag, fd.Annotations[tag])
		}
		fmt.Fprintf(fh, "\t},\n")
	}
	if len(fd.Comments) != 0 {
		var tags []string
		for tag := range fd.Comments {
			tags = append(tags, tag)
		}
		sort.Strings(tags)
		fmt.Fprintf(fh, "\tComments: map[string]string{\n")
		for _, tag := range tags {
			fmt.Fprintf(fh, "\t\t%q: %q,\n", tag, fd.Comments[tag])
		}
		fmt.Fprintf(fh, "\t},\n")
	}
	fmt.Fprintf(fh, "}\n")
}
func (ff *fieldDefinition) emit(fh *os.File) {
	fmt.Fprintf(fh, "\t\t{\n")
	fmt.Fprintf(fh, "\t\t\tTag: %q,\n", ff.Tag)
	if len(ff.Values) != 0 {
		fmt.Fprintf(fh, "\t\t\tValues: []string{")
		for i, v := range ff.Values {
			if i != 0 {
				fmt.Fprint(fh, ", ")
			}
			fmt.Fprintf(fh, "%q", v)
		}
		fmt.Fprint(fh, "},\n")
	}
	if len(ff.Validations) != 0 {
		fmt.Fprintf(fh, "\t\t\tValidations: []xscform.ValidateFunc{")
		for i, v := range ff.Validations {
			if i != 0 {
				fmt.Fprint(fh, ", ")
			}
			fmt.Fprint(fh, validateFuncName(v))
		}
		fmt.Fprint(fh, "},\n")
	}
	if ff.Default != "" {
		fmt.Fprintf(fh, "\t\t\tDefault: %q,\n", ff.Default)
	}
	fmt.Fprintf(fh, "\t\t},\n")
}
func validateFuncName(v string) string {
	v = "xscform.Validate-" + v
	for {
		if idx := strings.IndexByte(v, '-'); idx >= 0 {
			v = v[:idx] + strings.ToUpper(v[idx+1:idx+2]) + v[idx+2:]
		} else {
			return v
		}
	}
}

func emitPackage(pkgdir string, files []string) {
	pkg := filepath.Base(pkgdir)
	fh, err := os.Create(filepath.Join(pkgdir, pkg+".def.go"))
	if err != nil {
		log.Fatal(err)
	}
	defer fh.Close()
	fmt.Fprintf(fh, "// Code generated by extract-pifo-fields. DO NOT EDIT.\n\npackage %s\n\nimport \"steve.rothskeller.net/packet/xscmsg/internal/xscform\"\n\nvar formDefinitions = []*xscform.FormDefinition{\n", pkg)
	for i := len(files) - 1; i >= 0; i-- {
		fmt.Fprintf(fh, "\t%s,\n", files[i][:len(files[i])-5])
	}
	fmt.Fprintf(fh, "}\n")
}
