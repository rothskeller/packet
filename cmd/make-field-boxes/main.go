// usage: make-field-boxes [-p pagenum] png-file > form-file
//
// make-field-boxes is a tool to help create the PDF rendering specifications
// in a form file.  To use it, load the PDF form into an editor and make marks
// on it:  a red (#FF0000FF) mark in each text area, a green (#00FF00FF) mark
// in each checkbox, a blue (#0000FFFF) mark in each radio button, and magenta
// (#FF00FFFF) rectangles as needed.  Only those exact colors are handled.  The
// tool will write "pdf text", "pdf circle", "pdf cross", and "pdf box" lines
// to standard output for each region.  The red and green regions are enlarged
// to the surrounding edges with a margin; the blue regions are enlarged
// without a margin.  If a -p pagenum flag is specified, that page number is
// included in all generated lines.
package main

import (
	"cmp"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"slices"
	"strings"
)

type renderline struct {
	x    float64
	y    float64
	line string
}

var renderlines []renderline

func main() {
	var (
		pngfile string
		fh      *os.File
		im      image.Image
		bounds  image.Rectangle
		err     error
		pagenum = flag.Int("p", 1, "page number")
	)
	flag.Parse()
	if flag.NArg() != 1 || !strings.HasSuffix(flag.Arg(0), ".png") {
		fmt.Fprintln(os.Stderr, "usage: make-field-boxes [-p pagenum] png-file")
		os.Exit(2)
	}
	if fh, err = os.Open(flag.Arg(0)); err != nil {
		goto ERROR
	}
	if im, err = png.Decode(fh); err != nil {
		err = fmt.Errorf("%s: %w", pngfile, err)
		goto ERROR
	}
	fh.Close()
	bounds = im.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := im.At(x, y).RGBA()
			switch {
			case r == 0xffff && g == 0 && b == 0 && a == 0xffff:
				showRect(im, x, y, color.Black, 2, 2, "text", *pagenum)
			case r == 0 && g == 0xffff && b == 0 && a == 0xffff:
				showRect(im, x, y, color.Black, 0, 0, "cross", *pagenum)
			case r == 0 && g == 0 && b == 0xffff && a == 0xffff:
				showCircle(im, x, y, color.Black, 1, 3, *pagenum)
			case r == 0xffff && g == 0 && b == 0xffff && a == 0xffff:
				showRect(im, x, y, color.White, 0, 0, "box", *pagenum)
			}
		}
	}
	slices.SortFunc(renderlines, func(a, b renderline) int {
		if a.y != b.y {
			return -cmp.Compare(a.y, b.y)
		}
		return cmp.Compare(a.x, b.x)
	})
	for _, line := range renderlines {
		fmt.Println(line.line)
	}
	return

ERROR:
	fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
	os.Exit(1)
}

func showRect(im image.Image, x, y int, stop color.Color, margin, elm float64, keyword string, pagenum int) {
	var l, r, b, t = x, x + 1, y, y + 1
	// Grow the rectangle upward as far as the space remains blank.
	if isBlankRow(im, l, r, t, stop) {
		for t > im.Bounds().Min.Y && isBlankRow(im, l, r, t-1, stop) {
			t--
		}
	}
	// Grow the rectangle downward as far as the space remains blank.
	if isBlankRow(im, l, r, b, stop) {
		for b < im.Bounds().Max.Y && isBlankRow(im, l, r, b, stop) {
			b++
		}
	}
	// Grow the rectangle leftward as far as the space remains blank.
	if isBlankCol(im, l, t, b, stop) {
		for l > im.Bounds().Min.X && isBlankCol(im, l-1, t, b, stop) {
			l--
		}
	}
	// Grow the rectangle rightward as far as the space remains blank.
	if isBlankCol(im, r, t, b, stop) {
		for r < im.Bounds().Max.Y && isBlankCol(im, r, t, b, stop) {
			r++
		}
	}
	// Set the entire region to white, so it isn't picked up by another
	// area.
	draw.Draw(im.(draw.Image), image.Rect(l, t, r, b), &image.Uniform{color.White}, image.ZP, draw.Src)
	// Convert to PDF units.
	pl := float64(l-im.Bounds().Min.X)*612.0/float64(im.Bounds().Dx()) + margin + elm
	pr := float64(r-im.Bounds().Min.X)*612.0/float64(im.Bounds().Dx()) - margin
	pt := 792.0 - float64(t-im.Bounds().Min.Y)*792.0/float64(im.Bounds().Dy()) - margin
	pb := 792.0 - float64(b-im.Bounds().Min.Y)*792.0/float64(im.Bounds().Dy()) + margin
	// Write the line.
	line := fmt.Sprintf("  pdf      %s ", keyword)
	if pagenum != 1 {
		line += fmt.Sprintf("P %d ", pagenum)
	}
	line += fmt.Sprintf("L %6.2f R %6.2f B %6.2f T %6.2f", pl, pr, pb, pt)
	renderlines = append(renderlines, renderline{x: pl, y: pt, line: line})
}

func showCircle(im image.Image, x, y int, stop color.Color, margin, maxRadius float64, pagenum int) {
	// Expand the row at y to be as wide as possible in both directions.
	x1, x2 := x, x+1
	// debug = true
	widenRow(im, &x1, &x2, y, stop)
	// os.Exit(1)
	// If the rows above or below are wider, move y.
	for y > im.Bounds().Min.Y && ((x1 > im.Bounds().Min.X && isBlankRow(im, x1-1, x2, y-1, stop)) || (x2 < im.Bounds().Max.X && isBlankRow(im, x1, x2+1, y-1, stop))) {
		y--
		widenRow(im, &x1, &x2, y, stop)
	}
	for y < im.Bounds().Max.Y && ((x1 > im.Bounds().Min.X && isBlankRow(im, x1-1, x2, y+1, stop)) || (x2 < im.Bounds().Max.X && isBlankRow(im, x1, x2+1, y+1, stop))) {
		y++
		widenRow(im, &x1, &x2, y, stop)
	}
	// Find out how many rows there are around y with the same width.
	y1, y2 := y, y
	for y1 > im.Bounds().Min.Y && isBlankRow(im, x1, x2, y1-1, stop) {
		y1--
	}
	for y2 < im.Bounds().Max.Y && isBlankRow(im, x1, x2, y2+1, stop) {
		y2++
	}
	// Set the entire region to white, so it isn't picked up by another
	// area.
	my := (y1 + y2) / 2
	wy1, wy2 := my-(x2-x1)/2-1, my+(x2-x1)/2+1
	draw.Draw(im.(draw.Image), image.Rect(x1, wy1, x2, wy2), &image.Uniform{color.White}, image.ZP, draw.Src)
	// Convert x1, x2, y1, y2 to PDF units before calculating cx, cy, and r.
	px1 := float64(x1-im.Bounds().Min.X) * 612.0 / float64(im.Bounds().Dx())
	py1 := 792.0 - float64(y1-im.Bounds().Min.Y)*792.0/float64(im.Bounds().Dy())
	px2 := float64(x2-im.Bounds().Min.X) * 612.0 / float64(im.Bounds().Dx())
	py2 := 792.0 - float64(y2-im.Bounds().Min.Y)*792.0/float64(im.Bounds().Dy())
	cx := (px1 + px2) / 2
	cy := (py1 + py2) / 2
	r := (px2-px1)/2 - margin
	r = min(r, maxRadius)
	line := "  pdf      circle "
	if pagenum != 1 {
		line += fmt.Sprintf("P %d ", pagenum)
	}
	line += fmt.Sprintf("X %6.2f Y %6.2f R %4.2f", cx, cy, r)
	renderlines = append(renderlines, renderline{x: cx, y: cy, line: line})
}

func isBlankRow(png image.Image, l, r, y int, stop color.Color) bool {
	for x := l; x < r; x++ {
		if isColorMatch(png.At(x, y), stop) {
			return false
		}
	}
	return true
}

func isBlankCol(png image.Image, x, t, b int, stop color.Color) bool {
	for y := t; y < b; y++ {
		if isColorMatch(png.At(x, y), stop) {
			return false
		}
	}
	return true
}

func widenRow(im image.Image, l, r *int, y int, stop color.Color) {
	for *l > im.Bounds().Min.X && !isColorMatch(im.At(*l-1, y), stop) {
		*l--
	}
	for *r < im.Bounds().Max.X && !isColorMatch(im.At(*r, y), stop) {
		*r++
	}
}

func isColorMatch(t, c color.Color) bool {
	tr, tg, tb, ta := t.RGBA()
	cr, cg, cb, ca := c.RGBA()
	if tr == cr && tg == cg && tb == cb && ta == ca {
		return true
	}
	return ta == 0 && cr == 0xffff && cg == 0xffff && cb == 0xffff
}
