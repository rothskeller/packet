// move-field-boxes adjusts the locations of field boxes in a form by a
// specified amount.
//
// usage: move-field-boxes [flags] form-file
//
//	-p NUM      affect fields on page NUM (default 1)
//	-x NUM      distance to move in horizontal direction
//	-y NUM      distance to move in vertical direction (+=up, -=down)
//	-xmin NUM   only affect fields at least this far from the left
//	-xmax NUM   only affect fields no more than this far to the right
//	-ymin NUM   only affect fields at least this far from the bottom
//	-ymax NUM   only affect fields no more than this far up
package main

import (
	"fmt"
	"os"

	"github.com/rothskeller/packet/form/formdef"
	"github.com/rothskeller/pdf/v2"
	"github.com/spf13/pflag"
)

func main() {
	var (
		flags                  *pflag.FlagSet
		page                   int
		x, y                   float64
		xmin, xmax, ymin, ymax float64
		fd                     *formdef.FormDef
		err                    error
	)
	flags = pflag.NewFlagSet("", pflag.ExitOnError)
	flags.IntVar(&page, "p", 1, "page number")
	flags.Float64Var(&x, "x", 0, "horizontal movement")
	flags.Float64Var(&y, "y", 0, "vertical movement")
	flags.Float64Var(&xmin, "xmin", 0, "minimum X of target fields")
	flags.Float64Var(&xmax, "xmax", 1000, "maximum X of target fields")
	flags.Float64Var(&ymin, "ymin", 0, "minimum Y of target fields")
	flags.Float64Var(&ymax, "ymax", 1000, "maximum Y of target fields")
	flags.Parse(os.Args[1:])
	if flags.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "usage: move-field-boxes [flags] form-file")
		os.Exit(2)
	}
	if fd, err = formdef.Read(flags.Arg(0)); err != nil {
		goto ERROR
	}
	for f := range fd.AllFields() {
		for _, p := range f.PDF {
			switch p := p.Renderer.(type) {
			case formdef.BoxRenderer:
				if p.Page == page {
					adjustRectangle(&p.Rectangle, x, y, xmin, xmax, ymin, ymax)
				}
			case formdef.CircleRenderer:
				if p.Page == page {
					adjustPoint(&p.Center, x, y, xmin, xmax, ymin, ymax)
				}
			case formdef.CrossRenderer:
				if p.Page == page {
					adjustRectangle(&p.Rectangle, x, y, xmin, xmax, ymin, ymax)
				}
			case formdef.TextRenderer:
				if p.Page == page {
					adjustRectangle(&p.Rectangle, x, y, xmin, xmax, ymin, ymax)
					if p.Baseline != 0 && p.Rectangle.LLX >= xmin && p.Rectangle.URX <= xmax && p.Rectangle.LLY >= ymin && p.Rectangle.URY <= ymax {
						p.Baseline += y
					}
				}
			default:
				panic("unknown PDF renderer type")
			}
		}
	}
	if err = formdef.Write(flags.Arg(0), fd); err != nil {
		goto ERROR
	}
	return
ERROR:
	fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
	os.Exit(1)
}

func adjustRectangle(r *pdf.Rectangle, x, y, xmin, xmax, ymin, ymax float64) {
	if r.LLX < xmin || r.URX > xmax || r.LLY < ymin || r.URY > ymax {
		return
	}
	r.LLX, r.URX = r.LLX+x, r.URX+x
	r.LLY, r.URY = r.LLY+y, r.URY+y
}

func adjustPoint(p *pdf.Point, x, y, xmin, xmax, ymin, ymax float64) {
	if p.X < xmin || p.X > xmax || p.Y < ymin || p.Y > ymax {
		return
	}
	p.X, p.Y = p.X+x, p.Y+y
}
