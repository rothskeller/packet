// pdfrects overlays rectangles on a PDF file, allowing easy determination of
// bounding boxes for form fields.
//
// usage: pdfrects «pdf-file» rects-file»
//
// «rects-file» has lines of the form llx, lly, urx, ury
// Comments start with #, blank lines are ignored.
//
// output is written to basename(«pdf-file»).rects.pdf.
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/rothskeller/gofpdf"
	"github.com/rothskeller/gofpdf/contrib/gofpdi"

	fpdi "github.com/rothskeller/gofpdi"
)

func main() {
	var (
		rects [][]float64
		out   *gofpdf.Fpdf
		imp   *gofpdi.Importer
		sizes map[int]map[string]map[string]float64
		tpl   int
		outfn string
		err   error
	)
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: pdfrects pdf-file rects-file")
		os.Exit(2)
	}
	// Read the rectangles.
	rects = readRectangles(os.Args[2])
	// ImportPage will panic if the file doesn't exist or can't be read.
	// I'd rather have a clean error, so open it with fpdi to get the error
	// if any.
	if _, err = fpdi.NewPdfReader(os.Args[1]); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
	out = gofpdf.New("P", "pt", "Letter", "")
	imp = gofpdi.NewImporter()
	// Read the first page, just so the importer has the page details.
	// Ignore the result.
	imp.ImportPage(out, os.Args[1], 1, "/MediaBox")
	sizes = imp.GetPageSizes()
	// Handle each page.
	for pnum := 1; sizes[pnum] != nil; pnum++ {
		orient := "P"
		w, h := sizes[pnum]["/MediaBox"]["w"], sizes[pnum]["/MediaBox"]["h"]
		if w > h {
			orient = "L"
		}
		out.AddPageFormat(orient, gofpdf.SizeType{Wd: w, Ht: h})
		tpl = imp.ImportPage(out, os.Args[1], pnum, "/MediaBox")
		imp.UseImportedTemplate(out, tpl, 0, 0, w, h)
		drawRects(out, rects)
	}
	// Write the result.
	outfn = strings.TrimSuffix(os.Args[1], ".pdf") + ".rects.pdf"
	if err = out.OutputFileAndClose(outfn); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
}

func readRectangles(fname string) (rects [][]float64) {
	var (
		fh   *os.File
		scan *bufio.Scanner
		err  error
	)
	if fh, err = os.Open(fname); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
	defer fh.Close()
	scan = bufio.NewScanner(fh)
	for scan.Scan() {
		var llx, lly, urx, ury float64

		line := scan.Text()
		line, _, _ = strings.Cut(line, "#")
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 4 {
			goto SYNTAX
		}
		for i := range fields {
			fields[i] = strings.TrimSuffix(fields[i], ",")
		}
		if llx, err = strconv.ParseFloat(fields[0], 64); err != nil {
			goto SYNTAX
		}
		if lly, err = strconv.ParseFloat(fields[1], 64); err != nil {
			goto SYNTAX
		}
		if urx, err = strconv.ParseFloat(fields[2], 64); err != nil {
			goto SYNTAX
		}
		if ury, err = strconv.ParseFloat(fields[3], 64); err != nil {
			goto SYNTAX
		}
		if llx > urx || lly > ury {
			goto SYNTAX
		}
		rects = append(rects, []float64{llx, lly, urx, ury})
	}
	if err = scan.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s: %s\n", fname, err)
		os.Exit(1)
	}
	return rects
SYNTAX:
	fmt.Fprintf(os.Stderr, "ERROR: syntax error in %s\n", fname)
	os.Exit(1)
	return nil // not reachable
}

func drawRects(out *gofpdf.Fpdf, rects [][]float64) {
	_, h := out.GetPageSize()
	out.SetLineWidth(0.5)
	out.SetMargins(0, 0, 0)
	out.SetDrawColor(255, 0, 0)
	out.SetFillColor(255, 128, 128)
	for _, rect := range rects {
		out.Rect(rect[0], h-rect[3], rect[2]-rect[0], rect[3]-rect[1], "DF")
	}
}
