package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/rothskeller/packet/pktmsg"
	"github.com/rothskeller/packet/xscmsg"
	_ "github.com/rothskeller/packet/xscmsg/all"
)

func main() {
	var raw string
	var msg *pktmsg.Message
	var xsc *xscmsg.Message
	var problems []string
	var err error

	if len(os.Args) > 1 {
		if mb, err := os.ReadFile(os.Args[1]); err != nil {
			log.Fatal(err)
		} else {
			raw = string(mb)
		}
	} else {
		if mb, err := io.ReadAll(os.Stdin); err != nil {
			log.Fatal(err)
		} else {
			raw = string(mb)
		}
	}
	if msg, err = pktmsg.ParseMessage(raw); err != nil {
		log.Fatal(err)
	}
	xsc = xscmsg.Recognize(msg, true)
	fmt.Fprintf(os.Stderr, "Message Type: %s\n", xsc.Type.Tag)
	problems = xsc.Validate(true)
	for _, problem := range problems {
		fmt.Println(problem)
	}
	if len(problems) != 0 {
		os.Exit(1)
	}
	msg.Header.Set("Subject", xsc.Subject())
	msg.Body = xsc.Body(false)
	out := msg.Encode(false)
	if out != raw {
		fmt.Printf("ERROR: round trip mismatch: input was:\n%s\n=== output was:\n%s\n", raw, out)
		os.Exit(1)
	}
}
