package message

import (
	"fmt"
	"os"
	"strings"

	"github.com/rothskeller/packet/message/field"
	"github.com/rothskeller/pdf/v2"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/gofont/goregular"
)

// RenderPlainPDF creates a PDF representation of a message as plain text.
func RenderPlainPDF(m Message, filename, copyname string) (err error) {
	const (
		margin              = 48
		headingFont         = "Go-Bold"
		headingFontSize     = 14
		metadataFont        = "GoRegular"
		metadataFontSize    = 12
		metadataLineSpacing = 6
		metadataLabelFont   = "Go-Bold"
		metadataLabelWidth  = 60
		bodyFont            = "GoMono"
		bodyFontSize        = 10.5 // allows 80 columns to fit
		footerFont          = "GoRegular"
		footerFontSize      = 12
		timestampFormat     = "Monday, January 2, 2006 at 15:04:05"
	)
	var (
		fh       *os.File
		ph       *pdf.PDF
		avail    pdf.Rectangle
		from     string
		to       string
		date     string
		rcvd     string
		body     string
		msgID    string
		warnings pdf.ErrTextRendering
		page     = 1
	)
	if fh, err = os.Create(filename); err != nil {
		return err
	}
	ph = pdf.New(fh)
	pdf.AddTrueTypeFont(goregular.TTF)
	pdf.AddTrueTypeFont(gomono.TTF)
	pdf.AddTrueTypeFont(gobold.TTF)
	ph.Info["Title"] = m.Subject().EncodedSubject()
	ph.Info["Producer"] = "https://github.com/rothskeller/packet"
	ph.AddPage(pdf.USLetterPortrait)
	avail = pdf.USLetterPortrait
	// Set margins.
	avail.LLX, avail.LLY = avail.LLX+margin, avail.LLY+margin
	avail.URX, avail.URY = avail.URX-margin, avail.URY-margin
	// Room for footer.
	avail.LLY += 2 * footerFontSize
	// Add a banner.
	(pdf.Text{String: plainPDFBanner(m.Type().Name()), Rectangle: avail, Font: headingFont, FontSize: headingFontSize, Align: "lt"}).Draw(ph)
	avail.URY -= 2 * headingFontSize
	for f := range m.Fields() {
		switch f.Common() {
		case field.CHeaderFrom:
			from = f.Value(m)
		case field.CHeaderTo:
			to = f.Value(m)
		case field.CHeaderDate:
			date = f.Value(m)
		case field.CHeaderReceived:
			rcvd = f.Value(m)
		}
	}
	if from != "" {
		(pdf.Text{String: "From", Rectangle: avail, Font: metadataLabelFont, FontSize: metadataFontSize, Align: "lT"}).Draw(ph)
		avail.LLX += metadataLabelWidth // temporarily
		c, _ := (pdf.Text{String: from, Rectangle: avail, Font: metadataFont, FontSize: metadataFontSize, Align: "lT", Wrap: true}).DrawRect(ph)
		avail.LLX, avail.URY = avail.LLX-metadataLabelWidth, c.LLY-metadataLineSpacing
	}
	if to != "" {
		(pdf.Text{String: "To", Rectangle: avail, Font: metadataLabelFont, FontSize: metadataFontSize, Align: "lT"}).Draw(ph)
		avail.LLX += metadataLabelWidth // temporarily
		c, _ := (pdf.Text{String: m.To(), Rectangle: avail, Font: metadataFont, FontSize: metadataFontSize, Align: "lT", Wrap: true}).DrawRect(ph)
		avail.LLX, avail.URY = avail.LLX-metadataLabelWidth, c.LLY-metadataLineSpacing
	}
	if sub := m.Subject().EncodedSubject(); sub != "" {
		(pdf.Text{String: "Subject", Rectangle: avail, Font: metadataLabelFont, FontSize: metadataFontSize, Align: "lT"}).Draw(ph)
		avail.LLX += metadataLabelWidth // temporarily
		c, _ := (pdf.Text{String: sub, Rectangle: avail, Font: metadataFont, FontSize: metadataFontSize, Align: "lT", Wrap: true}).DrawRect(ph)
		avail.LLX, avail.URY = avail.LLX-metadataLabelWidth, c.LLY-metadataLineSpacing
	}
	if date != "" {
		(pdf.Text{String: "Date", Rectangle: avail, Font: metadataLabelFont, FontSize: metadataFontSize, Align: "lT"}).Draw(ph)
		avail.LLX += metadataLabelWidth // temporarily
		c, _ := (pdf.Text{String: date, Rectangle: avail, Font: metadataFont, FontSize: metadataFontSize, Align: "lT", Wrap: true}).DrawRect(ph)
		avail.LLX, avail.URY = avail.LLX-metadataLabelWidth, c.LLY-metadataLineSpacing
	}
	if rcvd != "" {
		(pdf.Text{String: "Received", Rectangle: avail, Font: metadataLabelFont, FontSize: metadataFontSize, Align: "lT"}).Draw(ph)
		avail.LLX += metadataLabelWidth // temporarily
		c, _ := (pdf.Text{String: rcvd, Rectangle: avail, Font: metadataFont, FontSize: metadataFontSize, Align: "lT", Wrap: true}).DrawRect(ph)
		avail.LLX, avail.URY = avail.LLX-metadataLabelWidth, c.LLY-metadataLineSpacing
	}
	avail.URY -= metadataFontSize + metadataLineSpacing
	body = m.Body().EncodedBody()
	body = strings.ReplaceAll(body, "\r", "")
	for body != "" {
		t := pdf.Text{String: body, Page: page, Rectangle: avail, Font: bodyFont, FontSize: bodyFontSize, Align: "lt", Wrap: true, Clip: true}
		fits, overflow, _, w := t.WrapText()
		switch w := w.(type) {
		case nil: // nothing
		case pdf.ErrTextRendering:
			w.DoesntFitY = false // we'll wrap to another page
			warnings.Merge(w)
		default:
			return w
		}
		t.String = fits
		t.Draw(ph)
		if body = overflow; body != "" {
			ph.AddPage(pdf.USLetterPortrait)
			page++
			avail = pdf.USLetterPortrait
			avail.LLX, avail.LLY = avail.LLX+margin, avail.LLY+margin
			avail.URX, avail.URY = avail.URX-margin, avail.URY-margin
			avail.LLY += 2 * footerFontSize
		}
	}
	msgID = m.Subject().SubjectMessageID()
	if err = RenderPDFFooters(ph, msgID, copyname); err != nil {
		return err
	}
	if err = ph.Write(); err != nil {
		return err
	}
	if err = fh.Close(); err != nil {
		return err
	}
	if w := warnings.AsError(); w != nil {
		return Warning{w}
	}
	return nil
}

// RenderPDFFooters writes a footer on every page of the PDF with a message ID,
// a copy name, and the page number.
func RenderPDFFooters(out *pdf.PDF, msgID, copyname string) (err error) {
	const (
		footerFont     = "GoRegular"
		footerFontSize = 12
	)
	var (
		pages int
		rect  = pdf.RectangleWH(36, 36, 540, 16)
		black = []byte{0, 0, 0}
	)
	if pages, err = out.NumPages(); err != nil {
		return err
	}
	for page := 1; page <= pages; page++ {
		if msgID != "" {
			(pdf.Text{String: msgID, Page: page, Rectangle: rect, Font: footerFont, FontSize: footerFontSize, Color: black, Align: "lB"}).Draw(out)
		}
		if copyname != "" {
			(pdf.Text{String: copyname, Page: page, Rectangle: rect, Font: footerFont, FontSize: footerFontSize, Color: black, Align: "cB"}).Draw(out)
		}
		(pdf.Text{String: fmt.Sprintf("Page %d of %d", page, pages), Page: page, Rectangle: rect, Font: footerFont, FontSize: footerFontSize, Color: black, Align: "rB"}).Draw(out)
	}
	return nil
}

func plainPDFBanner(typeName string) string {
	_, typeName, _ = strings.Cut(typeName, " ")
	return strings.ToUpper(typeName)
}
