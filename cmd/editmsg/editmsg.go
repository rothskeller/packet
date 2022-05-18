// editmsg is a command for creating and editing packet messages.
//
// usage: editmsg [-new type] [-opname name] [-opcall call] msgnum
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/rothskeller/packet/editmsg"
	_ "github.com/rothskeller/packet/xscmsg/all"
)

func main() {
	var (
		msgnum  string
		err     error
		newtype = flag.String("new", "", "type for new message")
		opname  = flag.String("opname", "", "operator name for new message")
		opcall  = flag.String("opcall", "", "operator call sign for new message")
	)
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "usage: editmsg [-new type] [-opname name] [-opcall call] msgnum\n")
		os.Exit(2)
	}
	msgnum = flag.Arg(0)
	if *newtype != "" {
		err = editmsg.NewMessage(*newtype, msgnum, *opname, *opcall)
	} else {
		err = editmsg.EditMessage(msgnum)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
}
