// Package jnosargs contains the command-line processing to get JNOS transport
// details.  It is common code used by all of the JNOS package commands.
package jnosargs

import (
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/rothskeller/packet/jnos"
	"github.com/rothskeller/packet/jnos/kpc3plus"
	"github.com/rothskeller/packet/jnos/simulator"
	"github.com/rothskeller/packet/jnos/telnet"
)

// Usage is string describing the transport-options.  It can be incorporated
// into the usage messages of the commands that use this package.
const Usage = `transport-options are:
    -bbs address
    -mbox mailbox
    -call callsign
    -port serialport
    -pwd password
For TNC-based RF connections:
    address must be an AX.25 address (A1AAA-1),
    mailbox must be set,
    callsign must be set if it's not the same as mailbox, and
    serialport must be set.
For network connections,
    address must be a host:port network address,
    mailbox must be set, and
    password must be set to either the password or the name of a file
    containing the password.
For JNOS simulations,
    address must be the name of a file containing the simulated messages
    to be read.
`

var (
	bbs  = flag.String("bbs", "", "BBS address")
	mbox = flag.String("mbox", "", "BBS mailbox")
	call = flag.String("call", "", "FCC callsign")
	port = flag.String("port", "", "serial port")
	pwd  = flag.String("pwd", "", "password or file containing password")
)

var (
	ax25RE = regexp.MustCompile(`(?i)^(?:A[A-L]|[KNW][A-Z]?)[0-9][A-Z]{1,3}-(?:[0-9]|1[0-5])$`)
	fccRE  = regexp.MustCompile(`(?i)^(?:A[A-L]|[KNW][A-Z]?)[0-9][A-Z]{1,3}$`)
)

// Connect connects to a BBS as specified through the command line parameters.
func Connect() (c *jnos.Conn, err error) {
	if *bbs == "" {
		fmt.Fprintf(os.Stderr, "ERROR: -bbs not specified\n")
		os.Exit(2)
	}
	if fh, err := os.Open(*bbs); err == nil {
		return connectSimulator(fh)
	}
	if ax25RE.MatchString(*bbs) {
		return connectKPC3Plus()
	}
	return connectTelnet()
}

func connectSimulator(fh *os.File) (*jnos.Conn, error) {
	if _, err := simulator.Start(map[string]io.Reader{"": fh}, "x"); err != nil {
		return nil, err
	}
	*bbs = simulator.ListenAddress
	*mbox = "x"
	*pwd = "x"
	return connectTelnet()
}

func connectKPC3Plus() (*jnos.Conn, error) {
	if *port == "" {
		fmt.Fprintf(os.Stderr, "ERROR: -port required for TNC connection\n")
		os.Exit(2)
	}
	if *mbox == "" {
		fmt.Fprintf(os.Stderr, "ERROR: -mbox required for BBS connection\n")
		os.Exit(2)
	}
	if !fccRE.MatchString(*mbox) && *call == "" {
		fmt.Fprintf(os.Stderr, "ERROR: -mbox doesn't look like an FCC call sign, so -call required\n")
		os.Exit(2)
	}
	if *call != "" && !fccRE.MatchString(*call) {
		fmt.Fprintf(os.Stderr, "ERROR: -call doesn't look like an FCC call sign\n")
		os.Exit(2)
	}
	return kpc3plus.Connect(*port, strings.ToUpper(*bbs), strings.ToUpper(*mbox), strings.ToUpper(*call), nil)
}

func connectTelnet() (*jnos.Conn, error) {
	if *mbox == "" {
		fmt.Fprintf(os.Stderr, "ERROR: -mbox required for BBS connection\n")
		os.Exit(2)
	}
	if *pwd == "" {
		fmt.Fprintf(os.Stderr, "ERROR: -pwd required for BBS connection\n")
		os.Exit(2)
	}
	if buf, err := os.ReadFile(*pwd); err == nil {
		*pwd = strings.TrimSpace(string(buf))
	}
	return telnet.Connect(*bbs, *mbox, *pwd, nil)
}
