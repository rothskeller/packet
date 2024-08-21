// pdfgrid overlays a grid on a PDF file, allowing easy determination of
// bounding boxes for form fields.
//
// usage: pdfgrid «pdf-file» [«grid-file»]
//
// «grid-file», if provided, is a list of "x y" coordinates, one pair per line.
// In its absence, lines are drawn every 10 points.
//
// output is written to basename(«pdf-file»).grid.pdf.
package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/phpdave11/gofpdf"
	"github.com/phpdave11/gofpdf/contrib/gofpdi"

	fpdi "github.com/phpdave11/gofpdi"
)

func main() {
	var (
		points [][]float64
		out    *gofpdf.Fpdf
		imp    *gofpdi.Importer
		sizes  map[int]map[string]map[string]float64
		tpl    int
		outfn  string
		err    error
	)
	if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Fprintln(os.Stderr, "usage: pdfgrid pdf-file [grid-file]")
		os.Exit(2)
	}
	// Read the points if any.
	if len(os.Args) == 3 {
		points = readPoints(os.Args[2])
	}
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
		drawGrid(out, points)
	}
	// Write the result.
	outfn = strings.TrimSuffix(os.Args[1], ".pdf") + ".grid.pdf"
	if err = out.OutputFileAndClose(outfn); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
}

func readPoints(fname string) (points [][]float64) {
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
		var x, y float64

		line := scan.Text()
		line, _, _ = strings.Cut(line, "#")
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 2 {
			goto SYNTAX
		}
		if x, err = strconv.ParseFloat(fields[0], 64); err != nil {
			goto SYNTAX
		}
		if y, err = strconv.ParseFloat(fields[1], 64); err != nil {
			goto SYNTAX
		}
		points = append(points, []float64{x, y})
	}
	if err = scan.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s: %s\n", fname, err)
		os.Exit(1)
	}
	return points
SYNTAX:
	fmt.Fprintf(os.Stderr, "ERROR: syntax error in %s\n", fname)
	os.Exit(1)
	return nil // not reachable
}

func drawGrid(out *gofpdf.Fpdf, points [][]float64) {
	var xs, ys []float64
	w, h := out.GetPageSize()
	out.SetLineWidth(0.5)
	out.SetFont("Helvetica", "", 8)
	out.SetMargins(0, 0, 0)
	if len(points) != 0 {
		for _, pt := range points {
			xs = append(xs, pt[0])
			ys = append(ys, pt[1])
		}
		slices.Sort(xs)
		slices.Sort(ys)
		xs = slices.Compact(xs)
		ys = slices.Compact(ys)
	} else {
		for x := 10.0; x < w; x += 10 {
			xs = append(xs, x)
		}
		for y := 10.0; y < h; y += 10 {
			ys = append(ys, y)
		}
	}
	for i, x := range xs {
		if i%3 == 2 {
			out.SetDrawColor(0, 255, 0)
			out.SetTextColor(0, 255, 0)
		} else {
			out.SetDrawColor(0, 0, 255)
			out.SetTextColor(0, 0, 255)
		}
		out.Line(x, 0, x, h)
		out.MoveTo(x, 0)
		out.TransformBegin()
		out.TransformRotate(270, x, 0)
		out.Write(8, strconv.FormatFloat(x, 'f', -1, 64))
		out.TransformEnd()
	}
	for i, y := range ys {
		if i%3 == 2 {
			out.SetDrawColor(0, 255, 0)
			out.SetTextColor(0, 255, 0)
		} else {
			out.SetDrawColor(0, 0, 255)
			out.SetTextColor(0, 0, 255)
		}
		out.Line(0, y, w, y)
		out.Text(0, y, strconv.FormatFloat(y, 'f', -1, 64))
	}
}
