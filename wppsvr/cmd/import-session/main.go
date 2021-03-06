// import-session reads a report generated by the old packet NCO scripts and
// imports it into the wppsvr database so that it is browsable through the
// wppsvr web UI.
package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/rothskeller/packet/wppsvr/config"
	"github.com/rothskeller/packet/wppsvr/store"
)

var session store.Session
var lines []string
var st *store.Store

func main() {
	var (
		fh  *os.File
		buf []byte
		err error
	)
	switch len(os.Args) {
	case 1:
		fh = os.Stdin
	case 2:
		if fh, err = os.Open(os.Args[1]); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "usage: import-session [report-filename]\n")
		os.Exit(2)
	}
	if buf, err = io.ReadAll(fh); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
	session.Report = string(buf)
	session.Imported = true
	if st, err = store.Open(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
	lines = strings.Split(session.Report, "\n")
	if len(lines) == 0 {
		log.Fatal("no input lines")
	}
	if !strings.HasPrefix(lines[0], "SCCo ARES/RACES Packet Practice Report For: ") {
		log.Fatal("missing header line")
	}
	if ts, _ := time.ParseInLocation("Monday, January 02, 2006.", lines[0][44:], time.Local); !ts.IsZero() {
		session.End = time.Date(ts.Year(), ts.Month(), ts.Day(), 20, 0, 0, 0, time.Local)
		session.Start = time.Date(ts.Year(), ts.Month(), ts.Day()-6, 0, 0, 0, 0, time.Local)
		switch ts.Weekday() {
		case time.Monday:
			session.CallSign = "PKTMON"
			session.Name = "SPECS Net"
			session.Prefix = "MON"
		case time.Tuesday:
			session.CallSign = "PKTTUE"
			session.Name = "SVECS Net"
			session.Prefix = "TUE"
		}
	}
	if session.CallSign == "" {
		log.Fatal("bad header line")
	}
	for len(lines) != 0 {
		if strings.HasPrefix(lines[0], "Simulated BBS outage: ") {
			session.DownBBSes = []string{lines[0][22:]}
			lines = lines[1:]
		}
		if strings.HasPrefix(lines[0], "##     FROM") {
			break
		}
		lines = lines[1:]
	}
	if len(lines) == 0 {
		log.Fatal("no message table header")
	}
	lines = lines[1:]
	st.CreateSession(&session)
	importMessages()
}

func importMessages() {
	var (
		msg *store.Message
	)
	for len(lines) != 0 {
		if lines[0] == "" || lines[0] == "Key:" {
			break
		}
		if lines[0][0] != ' ' { // new message line
			if msg != nil {
				parseErrors(msg)
				st.SaveMessage(msg)
			}
			msg = new(store.Message)
			msg.Session = session.ID
			sp := strings.IndexByte(lines[0][3:], ' ')
			if sp < 0 {
				sp = len(lines[0]) - 3
			}
			msg.FromAddress = lines[0][3 : 3+sp]
			msg.Subject = strings.TrimSpace(lines[0][3+sp:])
			msg.FromCallSign = strings.ToUpper(strings.Split(msg.FromAddress, "@")[0])
			sum := sha1.Sum([]byte(lines[0]))
			msg.Hash = hex.EncodeToString(sum[:])
			msg.LocalID = st.NextMessageID(session.Prefix)
		} else if strings.HasPrefix(lines[0], "     ^^^") {
			msg.Problems = append(msg.Problems, lines[0][9:])
			msg.Actions = config.ActionReport | config.ActionError
		} else if strings.HasPrefix(lines[0], "        ") {
			msg.Problems[len(msg.Problems)-1] += " " + lines[0][8:]
		}
		lines = lines[1:]
	}
	if msg != nil {
		parseErrors(msg)
		st.SaveMessage(msg)
	}
	importTooEarly()
}

func parseErrors(msg *store.Message) {
	var out []string

	for _, in := range msg.Problems {
		in = strings.TrimSpace(in)
		if strings.HasPrefix(in, "Duplicate check-in") {
			in = strings.TrimSpace(in[18:])
		}
		for in != "" {
			end := strings.IndexAny(in, ".(")
			if end < 0 {
				end = len(in) - 1
			}
			if in[end] != '(' {
				out = append(out, in[:end+1])
				in = strings.TrimSpace(in[end+1:])
				continue
			}
			close := strings.IndexByte(in, ')')
			parts := strings.Split(in[end+1:close], ";")
			for _, part := range parts {
				out = append(out, in[:end]+": "+strings.TrimSpace(part))
			}
			in = strings.TrimSpace(in[close+1:])
		}
	}
	msg.Problems = out
}

func importTooEarly() {
	for len(lines) != 0 {
		if strings.Contains(lines[0], "period which began") {
			break
		}
		lines = lines[1:]
	}
	if len(lines) == 0 {
		return
	}
	lines = lines[1:]
	for len(lines) != 0 {
		if len(lines[0]) != 0 && lines[0][0] >= '0' && lines[0][0] <= '9' {
			sum := sha1.Sum([]byte(lines[0]))
			parts := strings.SplitN(lines[0], " ", 3)
			msg := &store.Message{
				LocalID:      st.NextMessageID(session.Prefix),
				Hash:         hex.EncodeToString(sum[:]),
				Session:      session.ID,
				FromAddress:  parts[1],
				FromCallSign: strings.ToUpper(strings.Split(parts[1], "@")[0]),
				Subject:      strings.TrimSpace(parts[2]),
				Problems:     []string{"MessageTooEarly"},
				Actions:      config.ActionReport | config.ActionDontCount,
			}
			st.SaveMessage(msg)
		}
		lines = lines[1:]
	}
}
