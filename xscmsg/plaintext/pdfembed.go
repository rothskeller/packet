//go:build packetpdf

package plaintext

import (
	_ "embed" // .
	"fmt"
	"os"

	"github.com/go-pdf/fpdf"
	"github.com/rothskeller/packet/envelope"
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/xscmsg/delivrcpt"
)

//go:embed Go-Regular.ttf
var goRegular []byte

//go:embed Go-Bold.ttf
var goBold []byte

//go:embed Go-Mono.ttf
var goMono []byte

func init() {
	RenderPlainPDF = renderPDFActual
}

func renderPDFActual(env *envelope.Envelope, label, lmi, body, filename string) error {
	pdf := fpdf.New("P", "pt", "Letter", "")
	pdf.AddUTF8FontFromBytes("Go", "", goRegular)
	pdf.AddUTF8FontFromBytes("Go", "B", goBold)
	pdf.AddUTF8FontFromBytes("Go-Mono", "", goMono)
	pdf.SetMargins(36, 36, 36)
	pdf.SetAutoPageBreak(true, 48)
	pdf.AddPage()
	pdf.SetFont("Go", "B", 14)
	pdf.CellFormat(0, 21, label, "", 1, "L", false, 0, "")
	pdf.SetFontSize(12)
	if env.From != "" {
		pdf.Cell(63, 14.4, "From")
		pdf.SetFontStyle("")
		pdf.MultiCell(0, 14.4, env.From, "", "L", false)
		pdf.Ln(3.6)
	}
	if env.To != "" {
		pdf.SetFontStyle("B")
		pdf.Cell(63, 14.4, "To")
		pdf.SetFontStyle("")
		pdf.MultiCell(0, 14.4, env.To, "", "L", false)
		pdf.Ln(3.6)
	}
	pdf.SetFontStyle("B")
	pdf.Cell(63, 14.4, "Subject")
	pdf.SetFontStyle("")
	pdf.MultiCell(0, 14.4, env.SubjectLine, "", "L", false)
	pdf.Ln(3.6)
	if !env.Date.IsZero() {
		pdf.SetFontStyle("B")
		pdf.Cell(63, 14.4, "Sent")
		pdf.SetFontStyle("")
		pdf.CellFormat(0, 14.4, env.Date.Format("01/02/2006 15:04"), "", 1, "L", false, 0, "")
	}
	if received := getReceived(lmi, env); received != "" {
		pdf.Ln(3.6)
		pdf.SetFontStyle("B")
		pdf.Cell(63, 14.4, "Received")
		pdf.SetFontStyle("")
		pdf.CellFormat(0, 14.4, received, "", 1, "L", false, 0, "")
	}
	pdf.Ln(32.4)
	pdf.SetFont("Go-Mono", "", 11)
	pdf.MultiCell(0, 13.2, body, "", "L", false)
	return pdf.OutputFileAndClose(filename)
}

func getReceived(lmi string, env *envelope.Envelope) string {
	if env.IsReceived() {
		return fmt.Sprintf("%s as %s", env.ReceivedDate.Format("01/02/2006 15:04"), lmi)
	}
	// Not using incident.ReadReceipt here because that would create an
	// import cycle.
	contents, err := os.ReadFile(lmi + ".DR.txt")
	if err != nil {
		return ""
	}
	_, drbody, err := envelope.ParseSaved(string(contents))
	if err != nil {
		return ""
	}
	msg := message.Decode(env.SubjectLine, drbody)
	if dr, ok := msg.(*delivrcpt.DeliveryReceipt); ok {
		return fmt.Sprintf("%s as %s", dr.DeliveredTime, dr.LocalMessageID)
	}
	return ""
}
