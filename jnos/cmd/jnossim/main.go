package main

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"

	"github.com/rothskeller/packet/jnos/simulator"
)

func main() {
	var (
		messages = map[string]io.Reader{}
		home     string
		sim      *simulator.Simulator
		sig      = make(chan os.Signal, 1)
		err      error
	)
	for _, arg := range os.Args[1:] {
		var area, filename string
		var found bool
		if area, filename, found = strings.Cut(arg, ":"); found {
			if len(messages) == 0 {
				home = area
			}
			if messages[area], err = os.Open(filename); err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
				os.Exit(1)
			}
		} else if len(messages) == 0 {
			if messages[""], err = os.Open(filename); err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
				os.Exit(1)
			}
		} else {
			fmt.Fprintln(os.Stderr, "with multiple areas, every message file must be prefixed with area:")
			os.Exit(1)
		}
	}
	if len(messages) > 1 && messages[""] != nil {
		fmt.Fprintln(os.Stderr, "with multiple areas, every message file must be prefixed with area:")
		os.Exit(1)
	}
	if len(messages) == 0 {
		messages[""] = os.Stdin
	}
	signal.Notify(sig, os.Interrupt)
	sim, err = simulator.Start(messages, home)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
	}
	<-sig
	sim.Stop()
	fmt.Printf("%d messages sent:\n", len(sim.Sent()))
	for i, m := range sim.Sent() {
		fmt.Printf("==== SENT MESSAGE %d ====\n%s\n", i, m)
	}
}
