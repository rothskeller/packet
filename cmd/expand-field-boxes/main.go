// usage: expand-field-boxes form-file
//
// expand-field-boxes reads the form file and its PDF file.  For each field
// definition, it expands the rectangle or circle to encompass all of the
// surrounding whitespace, minus an appropriate margin.  It writes the
// resulting form file to *.form.new.
package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/fs"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rothskeller/packet/v4/form/formdef"
	"github.com/rothskeller/pdf/v2"
)

func main() {
	var (
		formfilename string
		relformfile  string
		fsh          fs.FS
		basename     string
		tmp          *os.File
		dest         *pdf.PDF
		def          *formdef.FormDef
		magick       *exec.Cmd
		pages        []image.Image
		err          error
	)
	if len(os.Args) != 2 || !strings.HasSuffix(os.Args[1], ".form") {
		fmt.Fprintln(os.Stderr, "usage: expand-field-boxes form-file")
		os.Exit(2)
	}
	formfilename = os.Args[1]
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
		err = fmt.Errorf("%s: %s", formfilename, err)
		goto ERROR
	}
	// Copy the PDF in the same way that the packet code does, so the
	// result looks the same.  (magick doesn't handle form field default
	// appearances the way the packet code does.)
	{
		var in fs.File
		var src *pdf.PDF
		if in, err = fsh.Open(def.PDFFile); err != nil {
			goto ERROR
		}
		if src, err = pdf.Open(in.(pdf.Reader)); err != nil {
			goto ERROR
		}
		if tmp, err = os.CreateTemp("", "expandfieldboxes*.pdf"); err != nil {
			goto ERROR
		}
		dest = pdf.New(tmp)
		if err = dest.ImportPDF(src); err != nil {
			goto ERROR
		}
		if err = dest.Write(); err != nil {
			goto ERROR
		}
		tmp.Close()
	}
	// Convert the PDF to PNGs.
	println("Converting PDF to PNGs...")
	magick = exec.Command("magick", "-density", "600", tmp.Name(), "+adjoin", basename+"-%d.png")
	if err = magick.Run(); err != nil {
		goto ERROR
	}
	// Read the PNGs.
	for pagenum := 1; ; pagenum++ {
		var (
			pngfile string
			fh      *os.File
			im      image.Image
		)
		pngfile = fmt.Sprintf("%s-%d.png", basename, pagenum-1)
		if fh, err = os.Open(pngfile); errors.Is(err, os.ErrNotExist) {
			break
		} else if err != nil {
			goto ERROR
		}
		if im, err = png.Decode(fh); err != nil {
			err = fmt.Errorf("%s: %w", pngfile, err)
			goto ERROR
		}
		fh.Close()
		pages = append(pages, im)
	}
	// Expand each of the field boxes.
	for fd := range def.AllFields() {
		for _, pr := range fd.PDF {
			if err = expandFieldBox(pr, pages, dest); err != nil {
				goto ERROR
			}
		}
	}
	// Move the old fields file aside and save the new fields.
	if err = formdef.Write(formfilename+".new", def); err != nil {
		goto ERROR
	}
	cleanupPNGs(basename)
	os.Exit(0)

ERROR:
	fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
	cleanupPNGs(basename)
	os.Exit(1)
}

func cleanupPNGs(basename string) {
	files, _ := filepath.Glob(basename + "-*.png")
	for _, file := range files {
		os.Remove(file)
	}
}

func expandFieldBox(pr formdef.PDFFieldRenderer, pages []image.Image, dest *pdf.PDF) (err error) {
	var (
		pagenum int
		png     image.Image
		path    pdf.Path
		pdfBox  pdf.Rectangle
	)
	switch shape := pr.Renderer.(type) {
	case formdef.CircleRenderer:
		pagenum = shape.Page
	case formdef.CrossRenderer:
		pagenum = shape.Page
	case formdef.TextRenderer:
		pagenum = shape.Page
	default:
		panic("unknown shape type")
	}
	png = pages[pagenum-1]
	if path, err = dest.PagePath(pagenum); err != nil {
		return err
	}
	if pdfBox, err = dest.GetRectangle(path.K("MediaBox")); err != nil {
		return err
	}
	switch shape := pr.Renderer.(type) {
	case formdef.CrossRenderer:
		reduceRectangle(&shape.Rectangle)
		return expandRectangle(png, pdfBox, &shape.Rectangle, 0)
	case formdef.CircleRenderer:
		return expandCircle(png, pdfBox, &shape.Center, &shape.Radius, 1)
	case formdef.TextRenderer:
		return expandRectangle(png, pdfBox, &shape.Rectangle, 2)
	default:
		panic("unknown shape type")
	}
}

func reduceRectangle(rect *pdf.Rectangle) {
	rect.LLX = (rect.LLX + rect.URX) / 2
	rect.URX = rect.LLX + 1
	rect.LLY = (rect.LLY + rect.URY) / 2
	rect.URY = rect.LLY + 1
}

func expandRectangle(png image.Image, pdfBox pdf.Rectangle, rect *pdf.Rectangle, margin float64) (err error) {
	// Convert to PNG units.
	var l, r, b, t int
	l, b = pdfToPNG(pdfBox, rect.LLX, rect.LLY, png.Bounds())
	r, t = pdfToPNG(pdfBox, rect.URX, rect.URY, png.Bounds())
	// Grow the rectangle upward as far as the space remains blank.
	if isBlankRow(png, l, r, t) {
		for t > 0 && isBlankRow(png, l, r, t-1) {
			t--
		}
	}
	// Grow the rectangle downward as far as the space remains blank.
	if isBlankRow(png, l, r, b) {
		for b < 6599 && isBlankRow(png, l, r, b) {
			b++
		}
	}
	// Grow the rectangle leftward as far as the space remains blank.
	if isBlankCol(png, l, t, b) {
		for l > 0 && isBlankCol(png, l-1, t, b) {
			l--
		}
	}
	// Grow the rectangle rightward as far as the space remains blank.
	if isBlankCol(png, r, t, b) {
		for r < 5099 && isBlankCol(png, r, t, b) {
			r++
		}
	}
	// Convert back to PDF units.
	rect.LLX, rect.LLY = pngToPDF(png.Bounds(), l, b, pdfBox)
	rect.URX, rect.URY = pngToPDF(png.Bounds(), r, t, pdfBox)
	// Apply margin.
	rect.LLX += margin
	rect.LLY += margin
	rect.URX -= margin
	rect.URY -= margin
	return nil
}

func expandCircle(png image.Image, pdfBox pdf.Rectangle, center *pdf.Point, radius *float64, margin float64) (err error) {
	// Convert to PNG units.
	var x, y, r int
	x, y = pdfToPNG(pdfBox, center.X, center.Y, png.Bounds())
	r = 1
	// Expand the row at y to be as wide as possible in both directions.
	x1, x2 := x-r, x+r
	widenRow(png, &x1, &x2, y)
	// If the rows above or below are wider, move y.
	for y > 0 && ((x1 > 0 && isBlankRow(png, x1-1, x2, y-1)) || (x2 < png.Bounds().Dx() && isBlankRow(png, x1, x2+1, y-1))) {
		y--
		widenRow(png, &x1, &x2, y)
	}
	for y < png.Bounds().Dy() && ((x1 > 0 && isBlankRow(png, x1-1, x2, y+1)) || (x2 < png.Bounds().Dx() && isBlankRow(png, x1, x2+1, y+1))) {
		y++
		widenRow(png, &x1, &x2, y)
	}
	// Find out how many rows there are around y with the same width.
	y1, y2 := y, y
	for y1 > 0 && isBlankRow(png, x1, x2, y1-1) {
		y1--
	}
	for y2 < png.Bounds().Dy() && isBlankRow(png, x1, x2, y2+1) {
		y2++
	}
	// Convert x1, x2, y1, y2 to PDF units before calculating cx, cy, and r.
	x1p, y1p := pngToPDF(png.Bounds(), x1, y1, pdfBox)
	x2p, y2p := pngToPDF(png.Bounds(), x2, y2, pdfBox)
	center.X = (x1p + x2p) / 2
	center.Y = (y1p + y2p) / 2
	*radius = (x2p-x1p)/2 - margin
	return nil
}

func widenRow(png image.Image, l, r *int, y int) {
	for *l > 0 && isWhite(png.At(*l-1, y)) {
		*l--
	}
	for *r < 5099 && isWhite(png.At(*r, y)) {
		*r++
	}
}

func isBlankRow(png image.Image, l, r, y int) bool {
	for x := l; x < r; x++ {
		if !isWhite(png.At(x, y)) {
			return false
		}
	}
	return true
}
func isBlankCol(png image.Image, x, t, b int) bool {
	for y := t; y < b; y++ {
		if !isWhite(png.At(x, y)) {
			return false
		}
	}
	return true
}

func isWhite(c color.Color) bool {
	r, g, b, a := c.RGBA()
	return a == 0 || (r == 0xffff && g == 0xffff && b == 0xffff)
}

func pdfToPNG(pdfBox pdf.Rectangle, pdfX, pdfY float64, pngBox image.Rectangle) (pngX, pngY int) {
	// Adjust PDF coordinates so that bottom left is at (0,0).
	pdfX -= pdfBox.LLX
	pdfY -= pdfBox.LLY
	// Scale to range [0,1).
	pdfX /= (pdfBox.URX - pdfBox.LLX)
	pdfY /= (pdfBox.URY - pdfBox.LLY)
	// Adjust so that (0,0) is top left.
	pdfY = 1.0 - pdfY
	// Scale to PNG box size.
	pdfX *= float64(pngBox.Dx())
	pdfY *= float64(pngBox.Dy())
	// Adjust to PNG box origin.
	pdfX += float64(pngBox.Min.X)
	pdfY += float64(pngBox.Min.Y)
	return int(math.Round(pdfX)), int(math.Round(pdfY))
}

func pngToPDF(pngBox image.Rectangle, pngX, pngY int, pdfBox pdf.Rectangle) (pdfX, pdfY float64) {
	// Adjust PNG coordinates so that top left is at (0,0).
	pngX -= pngBox.Min.X
	pngY -= pngBox.Min.Y
	// Scale to range [0,1).
	pdfX = float64(pngX) / float64(pngBox.Dx())
	pdfY = float64(pngY) / float64(pngBox.Dy())
	// Adjust so that (0,0) is bottom left.
	pdfY = 1.0 - pdfY
	// Scale to PDF box size.
	pdfX *= (pdfBox.URX - pdfBox.LLX)
	pdfY *= (pdfBox.URY - pdfBox.LLY)
	// Adjust to PDF box origin.
	pdfX += pdfBox.LLX
	pdfY += pdfBox.LLY
	return pdfX, pdfY
}
