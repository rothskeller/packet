package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/rothskeller/pdf/pdfstruct"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: extract-field-rects pdf-file")
		os.Exit(2)
	}
	if err := extractFromFile(os.Args[1]); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s: %s\n", os.Args[1], err)
		os.Exit(1)
	}
}

func extractFromFile(fname string) (err error) {
	var (
		fh     *os.File
		pdf    *pdfstruct.PDF
		form   pdfstruct.Dict
		fields pdfstruct.Array
	)
	if fh, err = os.Open(fname); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
	if pdf, err = pdfstruct.Open(fh); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s: %s\n", fname, err)
		os.Exit(1)
	}
	switch ref := pdf.Catalog["AcroForm"].(type) {
	case nil:
		return errors.New("PDF does not contain a form")
	case pdfstruct.Dict:
		form = ref
	case pdfstruct.Reference:
		if form, err = pdf.GetDict(ref); err != nil {
			return fmt.Errorf("AcroForm: %s", err)
		}
	default:
		return errors.New("AcroForm: not a Dict")
	}
	switch a := form["Fields"].(type) {
	case nil:
		break
	case pdfstruct.Array:
		fields = a
	default:
		return errors.New("AcroForm/Fields: not an Array")
	}
	return extractFieldRects(pdf, fields, "AcroForm/Fields", "")
}

func extractFieldRects(pdf *pdfstruct.PDF, fields pdfstruct.Array, path, prefix string) (err error) {
	for i, f := range fields {
		var field pdfstruct.Dict

		switch ref := f.(type) {
		case pdfstruct.Reference:
			if field, err = pdf.GetDict(ref); err != nil {
				return fmt.Errorf("%s[%d]: %s", path, i, err)
			}
		case pdfstruct.Dict:
			field = ref
		default:
			return fmt.Errorf("%s[%d]: not a Dict", path, i)
		}
		if err = extractFieldRect(pdf, field, fmt.Sprintf("%s[%d]", path, i), prefix); err != nil {
			return err
		}
	}
	return nil
}

func extractFieldRect(pdf *pdfstruct.PDF, field pdfstruct.Dict, path, prefix string) (err error) {
	var (
		name string
		kids pdfstruct.Array
		rect pdfstruct.Array
		bbox [4]float64
	)
	switch ns := field["T"].(type) {
	case nil:
		return nil
	case string:
		name = ns
	default:
		return fmt.Errorf("%s.T: not a string", path)
	}
	if prefix != "" && name != "" {
		name = prefix + "." + name
	} else if name == "" {
		name = prefix
	}
	switch ref := field["Rect"].(type) {
	case nil:
		break
	case pdfstruct.Reference:
		if rect, err = pdf.GetArray(ref); err != nil {
			return fmt.Errorf("%s.Rect: %s", path, err)
		}
	case pdfstruct.Array:
		rect = ref
	default:
		return fmt.Errorf("%s.Rect: not an Array", path)
	}
	if rect != nil {
		if len(rect) != 4 {
			return fmt.Errorf("%s.Rect: wrong number of elements (%d)", path, len(rect))
		}
		for i := 0; i < 4; i++ {
			switch num := rect[i].(type) {
			case int:
				bbox[i] = float64(num)
			case float64:
				bbox[i] = num
			default:
				return fmt.Errorf("%s.Rect[%d]: not a number", path, i)
			}
		}
		fmt.Printf("%g, %g, %g, %g\t# %s\n", bbox[0], bbox[1], bbox[2], bbox[3], name)
	}
	switch ref := field["Kids"].(type) {
	case nil:
		break
	case pdfstruct.Reference:
		if kids, err = pdf.GetArray(ref); err != nil {
			return fmt.Errorf("%s.Kids: %s", path, err)
		}
	case pdfstruct.Array:
		kids = ref
	default:
		return fmt.Errorf("%s.Kids: not an Array", path)
	}
	if kids != nil {
		if err = extractFieldRects(pdf, kids, path+".Kids", name); err != nil {
			return err
		}
	}
	return nil
}
