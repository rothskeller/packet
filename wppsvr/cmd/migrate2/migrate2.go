package main

import (
	"database/sql"
	"html"
	"log"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3" // database drive
	"github.com/rothskeller/packet/wppsvr/store"
)

var problemLabels = map[string]string{
	"BounceMessage":         "message has no return address (probably auto-response)",
	"CallSignConflict":      "call sign conflict",
	"DeliveryReceipt":       "DELIVERED receipt message",
	"FormCorrupt":           "incorrectly encoded form",
	"FormDestination":       "incorrect destination for form",
	"FormHandlingOrder":     "incorrect handling order for form",
	"FormInvalid":           "invalid form contents",
	"FormNoCallSign":        "no call sign in form message",
	"FormPracticeSubject":   "incorrect practice message details in form",
	"FormSubject":           "message subject doesn't agree with form contents",
	"FormToICSPosition":     `incorrect "To ICS Position" for form`,
	"FormToLocation":        `incorrect "To Location" for form`,
	"FormVersion":           "form version out of date",
	"FromBBSDown":           "message from incorrect BBS (simulated outage)",
	"HandlingOrderCode":     "unknown handling order code",
	"MessageCorrupt":        "message could not be parsed",
	"MessageFromWinlink":    "message sent from Winlink",
	"MessageNotASCII":       "message has non-ASCII characters",
	"MessageNotPlainText":   "not a plain text message",
	"MessageTooEarly":       "message before start of practice session",
	"MessageTooLate":        "message after end of practice session",
	"MessageTypeWrong":      "incorrect message type",
	"MsgNumFormat":          "incorrect message number format",
	"MsgNumPrefix":          "incorrect message number prefix",
	"NoCallSign":            "no call sign in message",
	"PIFOVersion":           "PackItForms version out of date",
	"PracticeAsFormName":    `incorrect subject line format (underline after "Practice")`,
	"PracticeSubjectFormat": "incorrect practice message details",
	"ReadReceipt":           "unexpected READ receipt message",
	"SessionDate":           "incorrect net date in subject",
	"SubjectFormat":         "incorrect subject line format",
	"SubjectHasSeverity":    "severity on subject line",
	"SubjectPlainForm":      "form name in subject of non-form message",
	"ToBBS":                 "message to incorrect BBS",
	"ToBBSDown":             "message to incorrect BBS (simulated outage)",
	"UnknownJurisdiction":   "unknown jurisdiction",
}

func main() {
	var (
		dbh *sql.DB
		tx  *sql.Tx
		err error
	)
	dbh, err = sql.Open("sqlite3", "file:wppsvr.db?mode=rw&cache=shared&_busy_timeout=1000&_txlock=immediate")
	if err != nil {
		log.Fatal("sqlite3.db:", err)
	}
	tx, err = dbh.Begin()
	if err != nil {
		log.Fatal("BEGIN:", err)
	}
	migrateMessages(tx)
	migrateSessions(tx)
	if err = tx.Commit(); err != nil {
		log.Fatal("COMMIT:", err)
	}
}

func migrateMessages(tx *sql.Tx) {
	var (
		rows *sql.Rows
		err  error
	)
	_, err = tx.Exec(`
CREATE TABLE new_message (
	id           text     PRIMARY KEY,
	hash         text     NOT NULL UNIQUE,
	deliverytime datetime NOT NULL,
	message      text     NOT NULL,
	session      integer  NOT NULL REFERENCES session ON DELETE CASCADE,
	fromaddress  text     NOT NULL,
	fromcallsign text     NOT NULL,
	frombbs      text     NOT NULL,
	tobbs        text     NOT NULL,
	jurisdiction text     NOT NULL,
	messagetype  text     NOT NULL,
	score        integer  NOT NULL,
	summary      text     NOT NULL,
	analysis     text     NOT NULL
);
`)
	if err != nil {
		log.Fatal("CREATE TABLE:", err)
	}
	rows, err = tx.Query(`SELECT id, hash, deliverytime, message, session, fromaddress, fromcallsign, frombbs, tobbs, jurisdiction, messagetype, subject, problems, actions FROM message`)
	if err != nil {
		log.Fatal("SELECT:", err)
	}
	for rows.Next() {
		var id, hash, message, fromAddress, fromCallSign, fromBBS, toBBS, jurisdiction, messageType, subject, problems, summary, analysis, response string
		var deliveryTime time.Time
		var session, actions, score int
		var importedProblems []string

		err = rows.Scan(&id, &hash, &deliveryTime, &message, &session, &fromAddress, &fromCallSign, &fromBBS, &toBBS, &jurisdiction, &messageType, &subject, &problems, &actions)
		if err != nil {
			log.Fatal("SELECT.scan:", err)
		}
		switch {
		case actions&24 != 0: // DontCount, DropMsg
			score = 0
		case actions&4 != 0: // Error
			score = 50
		case actions&3 != 0: // Respond, Report
			score = 80
		default:
			score = 100
		}
		if problems != "" {
			for _, problem := range strings.Split(problems, ";") {
				if summary != "" {
					summary = "multiple issues"
				} else {
					summary = problemLabels[problem]
					if summary == "" {
						if strings.IndexByte(problem, ' ') < 0 {
							log.Fatalf("unrecognized problem code %s", problem)
						}
						summary = problem
						importedProblems = append(importedProblems, problem)
					}
				}
			}
		}
		err = tx.QueryRow(`SELECT body FROM response WHERE responseto=? AND subject NOT LIKE 'DELIVERED:%'`, id).Scan(&response)
		switch err {
		case nil:
			analysis = "<h2>Analysis by pre-2023-08 packet practice software</h2><pre>" + html.EscapeString(response) + "</pre>"
		case sql.ErrNoRows:
			if len(importedProblems) != 0 {
				analysis = "<h2>Analysis by pre-2022-09 packet practice scripts</h2><pre>" + html.EscapeString(subject+"\n^^^ "+strings.Join(importedProblems, "\n^^^ ")) + "</pre>"
			}
		default:
			log.Fatal("SELECT.response:", err)
		}
		_, err = tx.Exec(`INSERT INTO new_message (id, hash, deliverytime, message, session, fromaddress, fromcallsign, frombbs, tobbs, jurisdiction, messagetype, score, summary, analysis) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
			id, hash, deliveryTime, message, session, fromAddress, fromCallSign, fromBBS, toBBS, jurisdiction, messageType, score, summary, analysis)
		if err != nil {
			log.Fatal("INSERT:", err)
		}
	}
	if err = rows.Err(); err != nil {
		log.Fatal("SELECT.end:", err)
	}
	if _, err = tx.Exec(`DROP TABLE message`); err != nil {
		log.Fatal("DROP TABLE:", err)
	}
	if _, err = tx.Exec(`ALTER TABLE new_message RENAME TO message`); err != nil {
		log.Fatal("ALTER TABLE:", err)
	}
	if _, err = tx.Exec(`CREATE INDEX message_session_idx ON message (session)`); err != nil {
		log.Fatal("CREATE INDEX:", err)
	}
}

func migrateSessions(tx *sql.Tx) {
	var (
		rows *sql.Rows
		err  error
	)
	_, err = tx.Exec(`
CREATE TABLE new_retrieval (
    session integer  REFERENCES session ON DELETE CASCADE,
    bbs     text     NOT NULL,
    lastrun datetime NOT NULL
)`)
	if err != nil {
		log.Fatal("CREATE TABLE:", err)
	}
	_, err = tx.Exec(`
CREATE TABLE new_session (
    id           integer  PRIMARY KEY,
    callsign     text     NOT NULL,
    name         text     NOT NULL,
    prefix       text     NOT NULL,
    start        datetime NOT NULL,
    end          datetime NOT NULL,
    reporttotext text     NOT NULL,
    reporttohtml text     NOT NULL,
    tobbses      text     NOT NULL,
    downbbses    text     NOT NULL,
    messagetypes text     NOT NULL,
    modelmessage text     NOT NULL,
    instructions text     NOT NULL,
    retrieveat   text     NOT NULL,
    report       text     NOT NULL,
    flags        integer  NOT NULL
)`)
	if err != nil {
		log.Fatal("CREATE TABLE:", err)
	}
	rows, err = tx.Query(`SELECT id, callsign, name, prefix, start, end, excludefromweek, reporttotext, reporttohtml, tobbses, downbbses, messagetypes, modified, running, imported, report FROM session`)
	if err != nil {
		log.Fatal("SELECT:", err)
	}
	for rows.Next() {
		var (
			id                                                                                                       int
			callsign, name, prefix, reporttotext, reporttohtml, tobbses, downbbses, messagetypes, report, retrieveat string
			start, end                                                                                               time.Time
			excludefromweek, modified, running, imported                                                             bool
			flags                                                                                                    store.SessionFlags
			retrievals                                                                                               map[string]time.Time
			rows2                                                                                                    *sql.Rows
		)
		err = rows.Scan(&id, &callsign, &name, &prefix, &start, &end, &excludefromweek, &reporttotext, &reporttohtml, &tobbses, &downbbses, &messagetypes, &modified, &running, &imported, &report)
		if err != nil {
			log.Fatal("SCAN:", err)
		}
		if strings.HasPrefix(reporttotext, "MESSAGE-SENDERS;") {
			flags |= store.ReportToSenders
			reporttotext = reporttotext[16:]
		} else if reporttotext == "MESSAGE-SENDERS" {
			flags |= store.ReportToSenders
			reporttotext = ""
		}
		if excludefromweek {
			flags |= store.ExcludeFromWeek
		}
		if modified {
			flags |= store.Modified
		}
		if running {
			flags |= store.Running
		}
		if imported {
			flags |= store.Imported
		}
		retrievals = make(map[string]time.Time)
		rows2, err = tx.Query(`SELECT bbs, interval, lastrun FROM retrieval WHERE session=?`, id)
		if err != nil {
			log.Fatal("SELECT:", err)
		}
		for rows2.Next() {
			var (
				bbs, when string
				lastrun   time.Time
			)
			err = rows2.Scan(&bbs, &when, &lastrun)
			if err != nil {
				log.Fatal("SCAN:", err)
			}
			if strings.Contains(when, "WEEKDAY") {
				if idx := strings.Index(when, "(DAY="); idx >= 0 {
					idx2 := strings.IndexByte(when[idx:], ' ')
					when = when[:idx+1] + when[idx+idx2+1:]
				}
				retrieveat = when
			}
			retrievals[bbs] = lastrun
		}
		if err = rows2.Err(); err != nil {
			log.Fatal("SELECT.end:", err)
		}
		_, err = tx.Exec(`INSERT INTO new_session (id, callsign, name, prefix, start, end, reporttotext, reporttohtml, tobbses, downbbses, messagetypes, modelmessage, instructions, retrieveat, report, flags) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
			id, callsign, name, prefix, start, end, reporttotext, reporttohtml, tobbses, downbbses, messagetypes, "", "", retrieveat, report, flags)
		if err != nil {
			log.Fatal("INSERT:", err)
		}
		for bbs, lastrun := range retrievals {
			_, err = tx.Exec(`INSERT INTO new_retrieval (session, bbs, lastrun) VALUES (?,?,?)`,
				id, bbs, lastrun)
			if err != nil {
				log.Fatal("INSERT:", err)
			}
		}
	}
	if err = rows.Err(); err != nil {
		log.Fatal("SELECT.end:", err)
	}
	if _, err = tx.Exec(`DROP TABLE retrieval`); err != nil {
		log.Fatal("DROP TABLE:", err)
	}
	if _, err = tx.Exec(`ALTER TABLE new_retrieval RENAME TO retrieval`); err != nil {
		log.Fatal("ALTER TABLE:", err)
	}
	if _, err = tx.Exec(`CREATE INDEX retrieval_session_idx ON retrieval (session)`); err != nil {
		log.Fatal("CREATE INDEX:", err)
	}
	if _, err = tx.Exec(`DROP TABLE session`); err != nil {
		log.Fatal("DROP TABLE:", err)
	}
	if _, err = tx.Exec(`ALTER TABLE new_session RENAME TO session`); err != nil {
		log.Fatal("ALTER TABLE:", err)
	}
	if _, err = tx.Exec(`CREATE UNIQUE INDEX session_call_end_idx ON session (callsign, end)`); err != nil {
		log.Fatal("CREATE INDEX:", err)
	}
	if _, err = tx.Exec(`CREATE INDEX session_end_idx ON session (end)`); err != nil {
		log.Fatal("CREATE INDEX:", err)
	}
	if _, err = tx.Exec(`CREATE INDEX session_running_idx ON session (flags) WHERE flags&1`); err != nil {
		log.Fatal("CREATE INDEX:", err)
	}
}
