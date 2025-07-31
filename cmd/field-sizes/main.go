package main

import (
	"fmt"
	"math"

	"github.com/rothskeller/packet/message"
	_ "github.com/rothskeller/packet/xscmsg"
)

func main() {
	for _, mtl := range message.RegisteredTypes {
		for _, mt := range mtl {
			m := message.Create(mt.Tag, mt.Version)
			if m == nil || mt.PDFBase == nil {
				fmt.Printf("Can't create %s %s.\n", mt.Tag, mt.Version)
				continue
			}
			fmt.Printf("%s %s:\n", mt.Tag, mt.Version)
			for _, f := range m.Base().Fields {
				showFieldSize(f, f.PDFRenderer)
			}
		}
	}
}

func showFieldSize(f *message.Field, r message.PDFRenderer) {
	switch r := r.(type) {
	case message.PDFMultiRenderer:
		for _, r2 := range r {
			showFieldSize(f, r2)
		}
	case *message.PDFTextRenderer:
		w, h := r.W, r.H
		if w == 0 {
			w = r.R - r.X
		}
		if h == 0 {
			h = r.B - r.Y
		}
		w = math.Ceil(w / (8 * 0.6))
		h = math.Ceil(h / (8 * 1.2))
		fmt.Printf("\t%s\tEditSize(%g, %g)\tEditWidth(%g)\n", f.Label, w, h, w)
	case *message.PDFMappedTextRenderer:
		h := r.H
		if h == 0 {
			h = r.B - r.Y
		}
		h = math.Ceil(h / 8 * 1.2)
		fmt.Printf("\t%s\tEditSize(?, %g)\n", f.Label, h)
	}
}
