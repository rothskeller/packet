// jnoslist lists all messages held in the specified mailbox.  Run the command
// without arguments for a usage description.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/rothskeller/packet/jnos/cmd/jnosargs"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: jnoslist [transport-options] [area [to]]\n%s", jnosargs.Usage)
	}
	flag.Parse()
	conn, err := jnosargs.Connect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
	if area := flag.Arg(0); area != "" {
		if err := conn.SetArea(area); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			os.Exit(1)
		}
	}
	list, err := conn.List(flag.Arg(1))
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		conn.Close()
		os.Exit(1)
	}
	if list == nil {
		fmt.Println("No messages.")
		conn.Close()
		return
	}
	fmt.Printf("Area: %s - %d messages, %d new\n", list.Area, list.Count, list.CountNew)
	for _, m := range list.Messages {
		var flags string
		if m.Deleted {
			flags = "D"
		} else if m.Held {
			flags = "H"
		} else {
			flags = " "
		}
		if m.Read {
			flags += "R"
		} else {
			flags += " "
		}
		fmt.Printf("%3d %s %s %-8.8s %-13.13s %3d %s\n", m.Number, flags, m.Date, m.FromPrefix, m.ToPrefix, m.Size, m.SubjectPrefix)
	}
	if err = conn.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
}
