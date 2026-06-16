// usage: extract-field-widths form-file
//
// extract-field-widths reads the specified form file and prints to standard
// output the CSS width specifier for each text area, assuming text at a font
// size of 1rem.
package main

import (
	"fmt"
	"os"

	"github.com/rothskeller/packet/v4/form/formdef"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: extract-field-widths form-file")
		os.Exit(2)
	}
	def, err := formdef.Read(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s: %s\n", os.Args[1], err)
		os.Exit(1)
	}
	for fd := range def.AllFields() {
		for _, pr := range fd.PDF {
			if pr, ok := pr.Renderer.(formdef.TextRenderer); ok {
				width := (pr.Rectangle.URX - pr.Rectangle.LLX) / (pr.FontSize / 12.0) / 12.0
				fmt.Printf("%-20.20s  style=\"width:%.2frem", fd.Label, width)
				if fd.Type == "multiline" {
					height := (pr.Rectangle.URY - pr.Rectangle.LLY) / (pr.FontSize / 12.0) / 12.0
					fmt.Printf(";height:%.2frem", height)
				}
				fmt.Println("\"")
			}
		}
	}
}
