// renderpdf renders one or more messages into PDFs.
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/rothskeller/packet/envelope"
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/xscmsg"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: renderpdf message-file...\n")
		os.Exit(2)
	}
	xscmsg.Register()
	for _, mfile := range os.Args[1:] {
		mbytes, err := os.ReadFile(mfile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", mfile, err)
			continue
		}
		env, body, err := envelope.ParseSaved(string(mbytes))
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", mfile, err)
			continue
		}
		msg := message.Decode(env.SubjectLine, body)
		pfile := mfile
		if strings.HasSuffix(mfile, ".txt") {
			pfile = pfile[:len(mfile)-4]
		}
		pfile += ".pdf"
		if err = msg.RenderPDF(pfile); err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", pfile, err)
		}
	}
}
