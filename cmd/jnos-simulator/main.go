package main

import (
	"fmt"
	"io"
	"os"
	"os/signal"

	"github.com/rothskeller/packet/jnos/simulator"
)

func main() {
	var (
		fh  *os.File
		sim *simulator.Simulator
		in  map[string]io.Reader
		err error
	)
	if len(os.Args) > 1 {
		if fh, err = os.Open(os.Args[1]); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			os.Exit(1)
		}
	}
	if fh != nil {
		in = map[string]io.Reader{"xndeoc": fh}
	}
	if sim, err = simulator.Start(in, "xndeoc"); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
	var ctrlc = make(chan os.Signal, 1)
	signal.Notify(ctrlc, os.Interrupt)
	<-ctrlc
	sim.Stop()
	for _, msg := range sim.Sent() {
		fmt.Println(msg)
	}
}
