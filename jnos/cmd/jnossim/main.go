package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/rothskeller/packet/jnos/simulator"
)

func main() {
	var (
		fh  *os.File
		sim *simulator.Simulator
		sig = make(chan os.Signal, 1)
		err error
	)
	if len(os.Args) > 1 {
		if fh, err = os.Open(os.Args[1]); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			os.Exit(1)
		}
		defer fh.Close()
	} else {
		fh = os.Stdin
	}
	signal.Notify(sig, os.Interrupt)
	sim, err = simulator.Start(fh)
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
