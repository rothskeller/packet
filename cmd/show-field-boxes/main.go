// usage: show-field-boxes form-file [output-file]
//
// show-field-boxes creates a PDF file using the template identified in the
// identified form file, and highlights the field drawing areas of all fields
// in the form.
//
// If an output-file is not given explicitly, the output is written to
// "${form-file%.form}.boxes.pdf".
package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/rothskeller/packet/cmd/packet/osdep"
	"github.com/rothskeller/packet/errors"
	"github.com/rothskeller/packet/form/formdef"
	"github.com/rothskeller/pdf/v2"
)

func main() {
	var (
		formfilename string
		relformfile  string
		outputname   string
		fsh          fs.FS
		def          *formdef.FormDef
		srcFH        fs.File
		src          *pdf.PDF
		destFH       *os.File
		dest         *pdf.PDF
		err          error
	)
	switch len(os.Args) {
	case 2:
		formfilename = os.Args[1]
		outputname = strings.TrimSuffix(formfilename, ".form") + ".boxes.pdf"
	case 3:
		formfilename, outputname = os.Args[1], os.Args[2]
	default:
		fmt.Fprintln(os.Stderr, "usage: show-field-boxes form-file [output-file]")
		os.Exit(2)
	}
	if filepath.IsAbs(formfilename) {
		fsh = os.DirFS(string(filepath.Separator))
		relformfile, err = filepath.Localize(formfilename[1:])
	} else {
		fsh = os.DirFS(".")
		relformfile, err = filepath.Localize(formfilename)
	}
	if err != nil {
		goto ERROR
	}
	if def, err = formdef.ReadFS(fsh, relformfile); err != nil {
		goto ERROR
	}
	if def.PDFFile == "" {
		err = fmt.Errorf("%s: no pdfFile set", formfilename)
		goto ERROR
	}
	if srcFH, err = fsh.Open(def.PDFFile); err != nil {
		goto ERROR
	}
	if src, err = pdf.Open(srcFH.(pdf.Reader)); err != nil {
		goto ERROR
	}
	if destFH, err = os.Create(outputname); err != nil {
		goto ERROR
	}
	dest = pdf.New(destFH)
	if err = dest.ImportPDF(src); err != nil {
		goto ERROR
	}
	for fd := range def.AllFields() {
		for _, pr := range fd.PDF {
			if err = markField(dest, pr); err != nil {
				err = fmt.Errorf("%s: %s", fd.Tag, err)
				goto ERROR
			}
		}
	}
	if err = dest.Write(); err != nil {
		goto ERROR
	}
	osdep.OpenFileCommand(outputname).Start()
	return
ERROR:
	fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
	os.Exit(1)

}

func markField(dest *pdf.PDF, pr formdef.PDFFieldRenderer) (err error) {
	switch pr := pr.Renderer.(type) {
	case formdef.CircleRenderer:
		err = markCircle(dest, pr.Page, pr.Center, pr.Radius)
	case formdef.CrossRenderer:
		err = markRect(dest, pr.Page, pr.Rectangle)
	case formdef.TextRenderer:
		err = errors.Join(
			markRect(dest, pr.Page, pr.Rectangle),
			markBaseline(dest, pr.Page, pr.Rectangle.LLX, pr.Rectangle.URX, pr.Baseline),
		)
	case formdef.BoxRenderer:
		err = markRect(dest, pr.Page, pr.Rectangle)
	}
	return err
}

func markRect(dest *pdf.PDF, pagenum int, rect pdf.Rectangle) (err error) {
	return (pdf.Box{
		Page:      pagenum,
		Rectangle: rect,
		Fill:      []byte{255, 0, 0, 128},
	}).Draw(dest)
}

func markBaseline(dest *pdf.PDF, pagenum int, x1, x2, y float64) (err error) {
	if y == 0 {
		return nil
	}
	return (pdf.Line{
		Page:   pagenum,
		P1:     pdf.PointXY(x1, y),
		P2:     pdf.PointXY(x2, y),
		Stroke: []byte{64, 64, 64},
		Width:  0.2,
	}).Draw(dest)
}

func markCircle(dest *pdf.PDF, pagenum int, center pdf.Point, radius float64) (err error) {
	return (pdf.Circle{
		Page:   pagenum,
		Center: center,
		Radius: radius,
		Fill:   []byte{255, 0, 0, 128},
	}).Draw(dest)
}
