// Package simulator provides a rudimentary simulation of a JNOS BBS, which can
// be used for testing JNOS-based services without connecting to a real BBS.
// This simulation implements only the features of JNOS that are used by package
// jnos.
//
// To use this simulator, call Start to start it, then create a jnos/telnet
// Transport connecting to simulator.ListenAddress.  Call the simulator's Stop
// method when it is no longer needed.
package simulator

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ListenAddress is the address to which the simulator listens.  Connect to it
// with a telnet transport.  Any login and password is accepted.
const ListenAddress = "localhost:63425"

// Start starts a new JNOS simulator.  It returnsa handle to the simulator,
// which can be examined after the interaction to see what messages were sent
// through the simulated BBS.  The messages available to read from the
// simulation are those in the supplied messages files, which can be mbox-style
// files or transcripts of an actual JNOS session that read messages.  home is
// the name of the default message area.  The simulator's Stop method should be
// called when it is no longer needed.
func Start(messages map[string]io.Reader, home string) (s *Simulator, err error) {
	s = new(Simulator)
	s.home = strings.ToLower(home)
	for area, reader := range messages {
		s.importMessages(strings.ToLower(area), reader)
	}
	if s.listener, err = net.Listen("tcp", ListenAddress); err != nil {
		return nil, err
	}
	go s.listen()
	return s, nil
}

// Simulator is a JNOS simulator.
type Simulator struct {
	area     string
	home     string
	messages map[string][]string
	sent     []string
	listener net.Listener
}

// Stop stops the simulator.
func (s *Simulator) Stop() {
	s.listener.Close()
}

// Sent returns the list of messages sent through the simulated BBS.
func (s *Simulator) Sent() []string { return s.sent }

var msgnum = regexp.MustCompile(`^Message #\d+$`)
var cmdprompt = regexp.MustCompile(`^\(#\d+\) >$`)
var fromline = regexp.MustCompile(`^From `)

// importMessages reads a file containing messages and makes them available to
// be "retrieved" from the simulated JNOS BBS.
func (s *Simulator) importMessages(area string, from io.Reader) {
	var (
		scan    *bufio.Scanner
		message string
		line    string
		inmsg   bool
	)
	if s.messages == nil {
		s.messages = make(map[string][]string)
	}
	s.messages[area] = nil
	scan = bufio.NewScanner(from)
	for scan.Scan() {
		line = scan.Text()
		if inmsg {
			if cmdprompt.MatchString(line) {
				inmsg = false
				s.messages[area] = append(s.messages[area], message)
				message = ""
				continue
			}
			if fromline.MatchString(line) {
				inmsg = false
				s.messages[area] = append(s.messages[area], message)
				message = ""
			}
		}
		if !inmsg {
			if msgnum.MatchString(line) {
				inmsg = true
				continue
			}
			if fromline.MatchString(line) {
				inmsg = true
			}
		}
		if inmsg {
			message = message + line + "\n"
		}
	}
	if inmsg {
		s.messages[area] = append(s.messages[area], message)
	}
}

// listen listens for incoming connections to the simulator, and hands them off
// to goroutines to be handled.
func (s *Simulator) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return
		}
		go s.handle(conn)
	}
}

// handle handles a single incoming connection to the simulator.
func (s *Simulator) handle(conn net.Conn) {
	var scan = bufio.NewScanner(conn)
	// First, lead the caller through a typical telnet login sequence, since
	// Start returns a telnet transport.  We don't care what login or
	// password they send.
	if _, err := conn.Write([]byte("login: ")); err != nil {
		conn.Close()
		return
	}
	scan.Scan()
	if _, err := conn.Write([]byte("Password [0] : ")); err != nil {
		conn.Close()
		return
	}
	scan.Scan()
	// JNOS command line loop.
	for {
		// Send a prompt and read a command.  The prompt has a bogus
		// message number in it; this simulator does not track a current
		// message number.
		if _, err := conn.Write([]byte("(#0) >\r\n")); err != nil {
			conn.Close()
			return
		}
		if !scan.Scan() {
			conn.Close()
			return
		}
		// Read a command from the command line and parcel it out to the
		// handler for it.  Only a small number of JNOS commands are
		// recognized; others are simply ignored.
		command := strings.ToLower(scan.Text())
		var err error
		if strings.HasPrefix(command, "l") {
			err = s.handleList(conn, command)
		} else if strings.HasPrefix(command, "a") {
			err = s.handleArea(conn, command)
		} else if strings.HasPrefix(command, "r") || strings.HasPrefix(command, "v") {
			err = s.handleRead(conn, command)
		} else if strings.HasPrefix(command, "s") {
			err = s.handleSend(conn, scan, command)
		} else if strings.HasPrefix(command, "b") {
			conn.Close()
			return
		}
		if err != nil {
			conn.Close()
			return
		}
	}
}

// handleArea handles the A command.
func (s *Simulator) handleArea(conn net.Conn, command string) (err error) {
	var (
		area string
	)
	time.Sleep(2 * time.Second)
	area = strings.ToLower(strings.TrimSpace(command[1:]))
	if _, ok := s.messages[area]; !ok {
		_, err = fmt.Fprintf(conn, "No such message area: %s\r\n", area)
		return err
	}
	s.area = area
	return nil
}

// handleList handles the LA and L> commands.
func (s *Simulator) handleList(conn net.Conn, command string) (err error) {
	var (
		wantto string
		buf    bytes.Buffer
		seen   bool
	)
	time.Sleep(2 * time.Second)
	if strings.HasPrefix(command, "l>") {
		wantto = strings.TrimSpace(command[2:])
	}
	fmt.Fprintf(&buf, "Mail area: ?\r\n%d messages  -  %d new\r\n\r\n", len(s.messages), len(s.messages))
	for i := 1; i <= len(s.messages[s.area]); i++ {
		if !seen {
			fmt.Fprint(&buf, "St.  #  TO            FROM     DATE   SIZE SUBJECT\r\n")
		}
		to, from, date, size, subject := parseMessage(s.messages[s.area][i-1])
		if to != "" && !strings.Contains(strings.ToLower(to), wantto) {
			continue
		}
		fmt.Fprintf(&buf, "  N %3d %-13.13s %-8.8s %s %4d %s\r\n", i, to, from, date, size, subject)
		seen = true
	}
	if !seen {
		fmt.Fprintf(&buf, "None to list.\r\n")
	}
	_, err = conn.Write(buf.Bytes())
	return err
}

// parseMessage is a *very* rudimentary mechanism to extract the message list
// fields from a raw message.
func parseMessage(m string) (to, from, date string, size int, subject string) {
	// Extract the interesting fields from the message headers.
	lines := strings.Split(m, "\n")
	for _, line := range lines {
		line = strings.TrimRight(line, "\r")
		if line == "" { // end of headers
			break
		}
		if to == "" && strings.HasPrefix(line, "To:") {
			to = strings.TrimSpace(line[3:])
		}
		if strings.HasPrefix(line, "From ") || strings.HasPrefix(line, "From:") {
			from = strings.TrimSpace(line[5:])
		}
		if date == "" && strings.HasPrefix(line, "Date:") {
			date = strings.TrimSpace(line[5:])
		}
		if subject == "" && strings.HasPrefix(line, "Subject:") {
			subject = strings.TrimSpace(line[8:])
		}
	}
	size = len(m)
	// Apply maximum lengths to the fields to match the list format.
	if len(to) > 13 {
		to = to[:13]
	}
	if len(from) > 8 {
		from = from[:8]
	}
	if len(subject) > 35 {
		subject = subject[:35]
	}
	// Convert the date format to that used in the list.
	if ts, err := time.ParseInLocation("Mon, 02 Jan", date[:11], time.Local); err == nil {
		date = ts.Format("Jan 02")
	} else if len(date) > 6 {
		date = date[:6]
	}
	return
}

// handleRead handles the R or V commands.  (It doesn't distinguish between
// them; it always returns the message in the form found in the input file.)
func (s *Simulator) handleRead(conn net.Conn, command string) (err error) {
	var msgnum int

	time.Sleep(2 * time.Second)
	if msgnum, err = strconv.Atoi(strings.TrimSpace(command[1:])); err != nil || msgnum < 1 || msgnum > len(s.messages[s.area]) {
		_, err = conn.Write([]byte("Invalid Message\r\n"))
		return err
	}
	if _, err = fmt.Fprintf(conn, "Message #%d \r\n", msgnum); err != nil {
		return err
	}
	_, err = conn.Write(bytes.ReplaceAll([]byte(s.messages[s.area][msgnum-1]), []byte{'\n'}, []byte{'\r', '\n'}))
	return err
}

// handleSend handles the SP and SC commands, by recording the outgoing message
// for later examination.
func (s *Simulator) handleSend(conn net.Conn, scan *bufio.Scanner, command string) (err error) {
	var (
		sb strings.Builder
	)
	time.Sleep(2 * time.Second)
	if strings.HasPrefix(command, "sp ") {
		// SP command.  Preserve the command line.
		to := strings.TrimSpace(command[3:])
		fmt.Fprintf(&sb, "SP %s\n", to)
	} else if strings.HasPrefix(command, "sc ") {
		// SC command.  Preserve the command line, then ask for Cc: and
		// preserve that too.
		to := strings.TrimSpace(command[3:])
		fmt.Fprintf(&sb, "SC %s\n", to)
		if _, err = conn.Write([]byte("Cc: ")); err != nil {
			return err
		}
		if !scan.Scan() {
			return scan.Err()
		}
		fmt.Fprintf(&sb, "Cc: %s\n", scan.Text())
	} else {
		// Some other command starting with S, not supported.
		_, err = conn.Write([]byte("Huh?\r\n"))
		return err
	}
	// Ask for the subject and preserve it.
	if _, err = conn.Write([]byte("Subject:\r\n")); err != nil {
		return err
	}
	if !scan.Scan() {
		return scan.Err()
	}
	fmt.Fprintf(&sb, "Subject: %s\n", scan.Text())
	// Ask for the body and preserve it.
	if _, err = conn.Write([]byte("Enter message.  End with /EX or ^Z in first column (^A aborts):\r\n")); err != nil {
		return err
	}
	for scan.Scan() {
		line := scan.Text()
		if line == "/EX" {
			s.sent = append(s.sent, sb.String())
			_, err = conn.Write([]byte("Msg queued\r\n"))
			return err
		}
		sb.WriteString(line)
		sb.WriteByte('\n')
	}
	return scan.Err()
}
