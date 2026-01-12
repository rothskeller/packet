package message

import (
	"fmt"
	"os"
	"strings"

	"github.com/rothskeller/packet/message/field"
	"github.com/rothskeller/pdf/v2"
)

// RenderPlainPDF creates a PDF representation of a message as plain text.
func RenderPlainPDF(m Message, filename, copyname string) (err error) {
	const (
		margin              = 48
		headingFont         = "Helvetica-Bold"
		headingFontSize     = 14
		metadataFont        = "Helvetica"
		metadataFontSize    = 12
		metadataLineSpacing = 14
		metadataLabelFont   = "Helvetica-Bold"
		metadataLabelWidth  = 60
		bodyFont            = "Courier"
		bodyFontSize        = 12
		footerFont          = "Helvetica"
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
		(pdf.Text{String: from, Rectangle: avail, Font: metadataFont, FontSize: metadataFontSize, Align: "lT"}).Draw(ph)
		avail.LLX, avail.URY = avail.LLX-metadataLabelWidth, avail.URY-metadataLineSpacing
	}
	if to != "" {
		(pdf.Text{String: "To", Rectangle: avail, Font: metadataLabelFont, FontSize: metadataFontSize, Align: "lT"}).Draw(ph)
		avail.LLX += metadataLabelWidth // temporarily
		(pdf.Text{String: m.To(), Rectangle: avail, Font: metadataFont, FontSize: metadataFontSize, Align: "lT"}).Draw(ph)
		avail.LLX, avail.URY = avail.LLX-metadataLabelWidth, avail.URY-metadataLineSpacing
	}
	if sub := m.Subject().EncodedSubject(); sub != "" {
		(pdf.Text{String: "Subject", Rectangle: avail, Font: metadataLabelFont, FontSize: metadataFontSize, Align: "lT"}).Draw(ph)
		avail.LLX += metadataLabelWidth // temporarily
		(pdf.Text{String: sub, Rectangle: avail, Font: metadataFont, FontSize: metadataFontSize, Align: "lT"}).Draw(ph)
		avail.LLX, avail.URY = avail.LLX-metadataLabelWidth, avail.URY-metadataLineSpacing
	}
	if date != "" {
		(pdf.Text{String: "Date", Rectangle: avail, Font: metadataLabelFont, FontSize: metadataFontSize, Align: "lT"}).Draw(ph)
		avail.LLX += metadataLabelWidth // temporarily
		(pdf.Text{String: date, Rectangle: avail, Font: metadataFont, FontSize: metadataFontSize, Align: "lT"}).Draw(ph)
		avail.LLX, avail.URY = avail.LLX-metadataLabelWidth, avail.URY-metadataLineSpacing
	}
	if rcvd != "" {
		(pdf.Text{String: "Received", Rectangle: avail, Font: metadataLabelFont, FontSize: metadataFontSize, Align: "lT"}).Draw(ph)
		avail.LLX += metadataLabelWidth // temporarily
		(pdf.Text{String: rcvd, Rectangle: avail, Font: metadataFont, FontSize: metadataFontSize, Align: "lT"}).Draw(ph)
		avail.LLX, avail.URY = avail.LLX-metadataLabelWidth, avail.URY-metadataLineSpacing
	}
	avail.URY -= metadataLineSpacing
	body = m.Body().EncodedBody()
	for body != "" {
		t := pdf.Text{String: body, Page: page, Rectangle: avail, Font: bodyFont, FontSize: bodyFontSize, Align: "lt", Wrap: true, Clip: true}
		fits, overflow, _, w := t.WrapText()
		if w != nil {
			etr := w.(pdf.ErrTextRendering)
			warnings.Merge(etr)
		}
		t.String = fits
		t.Draw(ph)
		if body = overflow; body != "" {
			ph.AddPage(pdf.USLetterPortrait)
			page++
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
	return warnings.AsError()
}

// RenderPDFFooters writes a footer on every page of the PDF with a message ID,
// a copy name, and the page number.
func RenderPDFFooters(out *pdf.PDF, msgID, copyname string) (err error) {
	const (
		footerFont     = "Helvetica"
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
