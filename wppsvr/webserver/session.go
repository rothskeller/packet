package webserver

import (
	"fmt"
	"html"
	"mime/multipart"
	"net/http"
	"net/mail"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/common"
	"github.com/rothskeller/packet/message/plaintext"
	"github.com/rothskeller/packet/wppsvr/config"
	"github.com/rothskeller/packet/wppsvr/htmlb"
	"github.com/rothskeller/packet/wppsvr/interval"
	"github.com/rothskeller/packet/wppsvr/store"
)

var removeCR = strings.NewReplacer("\r\n", "\n", "\r", "\n")

// serveSessionEdit allows the definition of a session to be edited.
func (ws *webserver) serveSessionEdit(w http.ResponseWriter, r *http.Request) {
	var (
		callsign          string
		session           *store.Session
		startDate         string
		startTime         string
		startError        string
		endDate           string
		endTime           string
		endError          string
		nameError         string
		callSignError     string
		prefixError       string
		reportToTextError string
		reportToHTMLError string
		bbsError          string
		retrievalsError   string
		mtype             string
		msgTypesError     string
		plainSubjectError string
		plainBodyError    string
		formBodyError     string
		formImages        []*multipart.FileHeader
		formImageError    string
	)
	if callsign = checkLoggedIn(w, r); callsign == "" {
		return
	}
	if !canEditSessions(callsign) {
		http.Error(w, "403 Forbidden", http.StatusForbidden)
		return
	}
	if sid, err := strconv.Atoi(r.FormValue("id")); err == nil {
		session = ws.st.GetSession(sid)
	}
	if session == nil {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		if r.FormValue("delete") != "" && session.Flags&store.Running == 0 && session.Report == "" {
			ws.st.DeleteSession(session)
			http.Redirect(w, r, fmt.Sprintf("/sessions?year=%d", session.End.Year()), http.StatusSeeOther)
			return
		}
		startDate, startTime, startError = readStart(r, session)
		endDate, endTime, endError = ws.readEnd(r, session)
		nameError = readName(r, session)
		callSignError = readCallSign(r, session)
		prefixError = readPrefix(r, session)
		readExcludeFromWeek(r, session)
		reportToTextError = readReportToText(r, session)
		reportToHTMLError = readReportToHTML(r, session)
		bbsError = readBBSes(r, session)
		retrievalsError = readRetrievals(r, session)
		mtype = readMessage(r)
		msgTypesError = readMsgTypes(r, session, mtype == "any")
		plainSubjectError = readPlainSubject(r, session, mtype == "plain")
		plainBodyError = readPlainBody(session, mtype == "plain")
		formBodyError = readFormBody(r, session, mtype == "form")
		formImages, formImageError = ws.readFormImage(r, session, mtype == "form")
		readInstructions(r, session)
		if startError == "" && endError == "" && nameError == "" && callSignError == "" &&
			prefixError == "" && reportToTextError == "" && reportToHTMLError == "" && bbsError == "" &&
			retrievalsError == "" && msgTypesError == "" && plainSubjectError == "" && plainBodyError == "" &&
			formBodyError == "" && formImageError == "" {
			var copyImagesFromSession int
			if r.FormValue("copy") != "" {
				copyImagesFromSession = session.ID
				session.ID = 0
				session.Flags &^= store.Modified | store.Running | store.Imported
				session.Report = ""
				ws.st.CreateSession(session)
			} else {
				if session.Flags&store.Running != 0 {
					session.Flags |= store.Modified
				}
				ws.st.UpdateSession(session)
			}
			if mtype != "form" {
				ws.st.DeleteModelImages(session.ID)
			} else if len(formImages) != 0 {
				ws.st.DeleteModelImages(session.ID)
				for i, fh := range formImages {
					body, err := fh.Open()
					if err == nil {
						ws.st.SaveModelImage(session.ID, i+1, fh.Filename, body)
						body.Close()
					}
				}
			} else if r.FormValue("copy") != "" {
				count := ws.st.ModelImageCount(copyImagesFromSession)
				for pnum := 1; pnum <= count; pnum++ {
					if body := ws.st.ModelImage(copyImagesFromSession, pnum); body != nil {
						ws.st.SaveModelImage(session.ID, pnum, body.Name(), body)
						body.Close()
					}
				}
			}
			http.Redirect(w, r, fmt.Sprintf("/sessions?year=%d", session.End.Year()), http.StatusSeeOther)
			return
		}
	} else {
		if !session.Start.IsZero() {
			startDate, startTime = session.Start.Format("2006-01-02"), session.Start.Format("15:04")
		}
		if !session.End.IsZero() {
			endDate, endTime = session.End.Format("2006-01-02"), session.End.Format("15:04")
		}
		switch session.ModelMsg.(type) {
		case nil:
			mtype = "any"
		case *plaintext.PlainText:
			mtype = "plain"
		default:
			mtype = "form"
		}
	}
	// Start the HTML page.
	w.Header().Set("Cache-Control", "nostore")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := htmlb.HTML(w)
	defer html.Close()
	html.E("meta charset=utf-8")
	html.E("title>Weekly Packet Practice - Santa Clara County ARES/RACES")
	html.E("meta name=viewport content='width=device-width, initial-scale=1'")
	html.E("link rel=stylesheet href=/static/common.css")
	html.E("link rel=stylesheet href=/static/form.css")
	html.E("link rel=stylesheet href=/static/session.css")
	html.E("script src=/static/session.js")
	html.E("div id=org>Santa Clara County ARES<sup>®</sup>/RACES")
	html.E("div id=title>Weekly Packet Practice")
	html.E("div id=subtitle>Practice Session Editor")
	// Write the form.
	form := html.E("form class=form method=POST enctype=multipart/form-data")
	emitStart(form, startDate, startTime, startError != "" || r.Method == http.MethodGet, startError)
	emitEnd(form, endDate, endTime, endError != "", endError)
	emitName(form, session, nameError != "", nameError)
	emitCallSign(form, session, callSignError != "", callSignError)
	emitPrefix(form, session, prefixError != "", prefixError)
	emitExcludeFromWeek(form, session)
	emitReportToText(form, session, reportToTextError != "", reportToTextError)
	emitReportToHTML(form, session, reportToHTMLError != "", reportToHTMLError)
	emitBBSes(form, session, bbsError != "", bbsError)
	emitRetrievals(form, session, retrievalsError != "", retrievalsError)
	emitMessage(form, session, mtype)
	emitMsgTypes(form, session, mtype == "any", msgTypesError != "", msgTypesError)
	emitPlainSubject(form, session, mtype == "plain", plainSubjectError != "", plainSubjectError)
	emitPlainBody(form, session, mtype == "plain", plainBodyError != "", plainBodyError)
	emitFormBody(form, session, mtype == "form", formBodyError != "", formBodyError)
	ws.emitFormImage(form, session, mtype == "form", formImageError != "", formImageError)
	emitInstructions(form, session)
	emitButtons(form, session)
}

func readStart(r *http.Request, session *store.Session) (datestr, timestr, err string) {
	datestr, timestr = r.FormValue("startDate"), r.FormValue("startTime")
	session.Start = time.Time{}
	if datestr == "" {
		return datestr, timestr, "The start date is required."
	}
	if timestr == "" {
		return datestr, timestr, "The start time is required."
	}
	if _, err := time.Parse("2006-01-02", datestr); err != nil {
		return datestr, timestr, "The start date is not a valid YYYY-MM-DD date."
	}
	if _, err := time.Parse("15:04", timestr); err != nil {
		return datestr, timestr, "The start time is not a valid HH:MM time."
	}
	session.Start, _ = time.ParseInLocation("2006-01-02 15:04", datestr+" "+timestr, time.Local)
	return datestr, timestr, ""
}
func emitStart(form *htmlb.Element, date, time string, focus bool, err string) {
	row := form.E("div class='formRow sessionStart'")
	row.E("label for=startDate>Start at")
	dt := row.E("div class='formInput formRange'")
	dt.E("input type=date id=startDate name=startDate value=%s", date)
	dt.E("input type=time id=startTime name=startTime value=%s step=300", time)
	if err != "" {
		row.E("div class=formError>%s", err)
	}
	row.E("div class=formHelp>Date and time when we start accepting practice messages for this session.")
}

func (ws *webserver) readEnd(r *http.Request, session *store.Session) (datestr, timestr, err string) {
	datestr, timestr = r.FormValue("endDate"), r.FormValue("endTime")
	session.End = time.Time{}
	if datestr == "" {
		return datestr, timestr, "The end date is required."
	}
	if timestr == "" {
		return datestr, timestr, "The end time is required."
	}
	if _, err := time.Parse("2006-01-02", datestr); err != nil {
		return datestr, timestr, "The end date is not a valid YYYY-MM-DD date."
	}
	if _, err := time.Parse("15:04", timestr); err != nil {
		return datestr, timestr, "The end time is not a valid HH:MM time."
	}
	session.End, _ = time.ParseInLocation("2006-01-02 15:04", datestr+" "+timestr, time.Local)
	for _, existing := range ws.st.GetSessions(
		time.Date(session.End.Year(), session.End.Month(), session.End.Day(), 0, 0, 0, 0, time.Local),
		time.Date(session.End.Year(), session.End.Month(), session.End.Day()+1, 0, 0, 0, 0, time.Local),
	) {
		if existing.ID == session.ID && r.FormValue("copy") == "" {
			continue
		}
		return datestr, timestr, "Another session already ends on this date."
	}
	return datestr, timestr, ""
}
func emitEnd(form *htmlb.Element, date, time string, focus bool, err string) {
	row := form.E("div class='formRow sessionEnd'")
	row.E("label for=endDate>End at")
	dt := row.E("div class='formInput formRange'")
	dt.E("input type=date id=endDate name=endDate value=%s", date)
	dt.E("input type=time id=endTime name=endTime value=%s step=300", time)
	if err != "" {
		row.E("div class=formError>%s", err)
	}
	row.E("div class=formHelp>Date and time when we stop accepting practice messages for this session.  The session report will be sent at this time.")
}

func readName(r *http.Request, session *store.Session) string {
	if session.Name = strings.TrimSpace(r.FormValue("name")); session.Name == "" {
		return "The session name is required."
	}
	return ""
}
func emitName(form *htmlb.Element, session *store.Session, focus bool, err string) {
	row := form.E("div class='formRow sessionName'")
	row.E("label for=sessionName>Session name")
	row.E("input id=sessionName name=name value=%s", session.Name, focus, "autofocus")
	if err != "" {
		row.E("div class=formError>%s", err)
	}
	row.E("div class=formHelp>Name of the practice session.  Usually this is the name of the net that participants are checking into.  Do not include the date.")
}

var callSignRE = regexp.MustCompile(`(?i)^[A-Z][A-Z0-9]{2,5}$`)

func readCallSign(r *http.Request, session *store.Session) string {
	if session.CallSign = strings.ToUpper(strings.TrimSpace(r.FormValue("callsign"))); session.CallSign == "" {
		return "The call sign is required."
	} else if !callSignRE.MatchString(session.CallSign) {
		return "This is not a valid call sign.  The call sign is expected to contain between 3 and 6 letters and digits, starting with a letter."
	}
	return ""
}
func emitCallSign(form *htmlb.Element, session *store.Session, focus bool, err string) {
	row := form.E("div class='formRow callsign'")
	row.E("label for=callsign>Call sign")
	row.E("input id=callsign name=callsign value=%s", session.CallSign, focus, "autofocus")
	if err != "" {
		row.E("div class=formError>%s", err)
	}
	row.E("div class=formHelp>Call sign to which practice messages must be addressed, i.e., name of the JNOS mailbox from which to retrieve messages.")
}

var prefixRE = regexp.MustCompile(`(?i)^[A-Z][A-Z0-9]{2}$`)

func readPrefix(r *http.Request, session *store.Session) string {
	if session.Prefix = strings.ToUpper(strings.TrimSpace(r.FormValue("prefix"))); session.Prefix == "" {
		return "The message number prefix is required."
	} else if !prefixRE.MatchString(session.Prefix) {
		return "This is not a valid message number prefix.  The prefix must be three letters or digits, starting with a letter."
	}
	return ""
}
func emitPrefix(form *htmlb.Element, session *store.Session, focus bool, err string) {
	row := form.E("div class='formRow prefix'")
	row.E("label for=prefix>Message number prefix")
	row.E("input id=prefix name=prefix value=%s", session.Prefix, focus, "autofocus")
	if err != "" {
		row.E("div class=formError>%s", err)
	}
	row.E("div class=formHelp>Three-character message number prefix corresponding to call sign.")
}

func readExcludeFromWeek(r *http.Request, session *store.Session) {
	if r.FormValue("exclude") != "" {
		session.Flags |= store.ExcludeFromWeek
	} else {
		session.Flags &^= store.ExcludeFromWeek
	}
}
func emitExcludeFromWeek(form *htmlb.Element, session *store.Session) {
	row := form.E("div class='formRow exclude'")
	row.E("label for=exclude>Exclude from Week")
	row.E("div class=formInput").
		E("input type=checkbox id=exclude name=exclude", session.Flags&store.ExcludeFromWeek != 0, "checked")
	row.E("div class=formHelp>Exclude this session from weekly check-in counts.")
}

func readReportToText(r *http.Request, session *store.Session) (err string) {
	session.ReportToText = strings.Fields(r.FormValue("reportToText"))
	for _, addr := range session.ReportToText {
		if _, bad := mail.ParseAddress(addr); bad != nil {
			err = "“" + html.EscapeString(addr) + "” is not a valid packet address."
		}
	}
	if r.FormValue("reportToSenders") != "" {
		session.Flags |= store.ReportToSenders
	} else {
		session.Flags &^= store.ReportToSenders
	}
	return err
}
func emitReportToText(form *htmlb.Element, session *store.Session, focus bool, err string) {
	row := form.E("div class='formRow reportToText'")
	row.E("label for=reportToText class=textareaLabel>Packet Reports")
	in := row.E("div class='formInput reportTo'")
	in.E("textarea id=reportToText name=reportToText class=formInput", focus, "autofocus").R(strings.Join(session.ReportToText, "\n"))
	in.E("input type=checkbox id=reportToSenders name=reportToSenders", session.Flags&store.ReportToSenders != 0, "checked")
	in.E("label for=reportToSenders> Message senders")
	if err != "" {
		row.E("div class=formError>%s", err)
	}
	row.E("div class=formHelp>Packet addresses to which the plain text session report should be sent (one per line).")
}

func readReportToHTML(r *http.Request, session *store.Session) (err string) {
	session.ReportToHTML = strings.Fields(r.FormValue("reportToHTML"))
	for _, addr := range session.ReportToHTML {
		if _, bad := mail.ParseAddress(addr); bad != nil {
			err = "“" + html.EscapeString(addr) + "” is not a valid email address."
		}
	}
	return err
}
func emitReportToHTML(form *htmlb.Element, session *store.Session, focus bool, err string) {
	row := form.E("div class='formRow reportToHTML'")
	row.E("label for=reportToHTML>Email Reports")
	row.E("textarea id=reportToHTML name=reportToHTML class=formInput", focus, "autofocus").R(strings.Join(session.ReportToHTML, "\n"))
	if err != "" {
		row.E("div class=formError>%s", err)
	}
	row.E("div class=formHelp>Email addresses to which the HTML-formatted session report should be sent (one per line).")
}

func readBBSes(r *http.Request, session *store.Session) string {
	session.ToBBSes, session.DownBBSes = session.ToBBSes[:0], session.DownBBSes[:0]
	var bbsnames = make([]string, 0, len(config.Get().BBSes))
	var retrieve = make(map[string]*store.Retrieval)
	for name := range config.Get().BBSes {
		bbsnames = append(bbsnames, name)
	}
	sort.Strings(bbsnames)
	for _, r := range session.Retrieve {
		retrieve[r.BBS] = r
	}
	session.Retrieve = session.Retrieve[:0]
	for _, name := range bbsnames {
		ret := r.FormValue("retrieve."+name) != ""
		if ret {
			if retrieve[name] != nil {
				session.Retrieve = append(session.Retrieve, retrieve[name])
			} else {
				session.Retrieve = append(session.Retrieve, &store.Retrieval{BBS: name})
			}
		}
		isTo := ret && r.FormValue("destbbs."+name) != ""
		if isTo {
			session.ToBBSes = append(session.ToBBSes, name)
		}
		if !isTo && r.FormValue("downbbs."+name) != "" {
			session.DownBBSes = append(session.DownBBSes, name)
		}
	}
	if len(session.ToBBSes) == 0 {
		return "At least one BBS must be marked as a destination for practice messages."
	}
	return ""
}
func emitBBSes(form *htmlb.Element, session *store.Session, focus bool, err string) {
	var bbsnames = make([]string, 0, len(config.Get().BBSes))
	var retrieve = make(map[string]bool)
	for name := range config.Get().BBSes {
		bbsnames = append(bbsnames, name)
	}
	sort.Strings(bbsnames)
	for _, r := range session.Retrieve {
		retrieve[r.BBS] = true
	}
	for i, name := range bbsnames {
		row := form.E("div class=formRow")
		row.E("label for=retrieve.%s", name).TF("BBS %s", name)
		in := row.E("div class=formInput")
		in.E("input type=checkbox id=retrieve.%s name=retrieve.%s", name, name, retrieve[name], "checked", focus && i == 0, "autofocus")
		in.E("label for=retrieve.%s> Retrieve practice messages", name)
		in.E("br")
		in.E("input type=checkbox id=destbbs.%s name=destbbs.%s", name, name, inList(session.ToBBSes, name), "checked")
		in.E("label for=destbbs.%s> Accept retrieved practice messages", name)
		in.E("br")
		in.E("input type=checkbox id=downbbs.%s name=downbbs.%s", name, name, inList(session.DownBBSes, name), "checked")
		in.E("label for=downbbs.%s> Simulated outage", name)
		in.E("br")
		if i == 0 && err != "" {
			row.E("div class=formError>%s", err)
		}
	}
}

func readRetrievals(r *http.Request, session *store.Session) string {
	if r.FormValue("dontKillMessages") != "" {
		session.Flags |= store.DontKillMessages
	} else {
		session.Flags &^= store.DontKillMessages
	}
	if r.FormValue("dontSendResponses") != "" {
		session.Flags |= store.DontSendResponses
	} else {
		session.Flags &^= store.DontSendResponses
	}
	session.RetrieveAt = strings.TrimSpace(removeCR.Replace(r.FormValue("retrievals")))
	if interval.Parse(session.RetrieveAt) == nil {
		return "This is not a valid schedule string."
	}
	return ""
}
func emitRetrievals(form *htmlb.Element, session *store.Session, focus bool, err string) {
	row := form.E("div class=formRow")
	row.E("label for=reportToText class=textareaLabel>Retrieval Schedule")
	in := row.E("div class=formInput")
	in.E("textarea id=retrievals name=retrievals class=formInput", focus, "autofocus").R(session.RetrieveAt)
	in.E("input type=checkbox id=dontKillMessages name=dontKillMessages", session.Flags&store.DontKillMessages != 0, "checked")
	in.E("label for=dontKillMessages> Leave messages on BBS")
	in.E("br")
	in.E("input type=checkbox id=dontSendResponses name=dontSendResponses", session.Flags&store.DontSendResponses != 0, "checked")
	in.E("label for=dontSendResponses> Don’t send delivery receipts")
	if err != "" {
		row.E("div class=formError>%s", err)
	}
	row.E("div class=formHelp>Schedule for when to retrieve practice messages from BBSes.")
}

func readMessage(r *http.Request) string {
	switch mtype := r.FormValue("messageType"); mtype {
	case "any", "plain", "form":
		return mtype
	default:
		return "any"
	}
}
func emitMessage(form *htmlb.Element, session *store.Session, mtype string) {
	row := form.E("div class=formRow")
	row.E("label for=messageAny>Message to Send")
	in := row.E("div class=formInput")
	in.E("input type=radio id=anyMessage name=messageType value=any", mtype == "any", "checked")
	in.E("label for=dontKillMessages> Any message of specified type(s)")
	in.E("br")
	in.E("input type=radio id=plainMessage name=messageType value=plain", mtype == "plain", "checked")
	in.E("label for=dontSendResponses> Copy of provided plain text message")
	in.E("br")
	in.E("input type=radio id=formMessage name=messageType value=form", mtype == "form", "checked")
	in.E("label for=dontSendResponses> Copy of provided form")
}

func readMsgTypes(r *http.Request, session *store.Session, show bool) string {
	session.MessageTypes = session.MessageTypes[:0]
	session.ModelMessage = ""
	if !show {
		return ""
	}
	for _, tag := range r.Form["mtype"] {
		if tag == "plain" || config.Get().MessageTypes[tag] != nil {
			session.MessageTypes = append(session.MessageTypes, tag)
		}
	}
	if len(session.MessageTypes) == 0 {
		return "At least one message type must be accepted."
	}
	return ""
}
func emitMsgTypes(form *htmlb.Element, session *store.Session, show, focus bool, err string) {
	row := form.E("div id=mtypeRow class=formRow", !show, "style=display:none")
	row.E("label>Message Type(s)")
	in := row.E("div class=formInput")
	in.E("input type=checkbox id=mtype.plain name=mtype value=plain", inList(session.MessageTypes, plaintext.Type.Tag), "checked", focus, "autofocus")
	in.E("label for=mtype.plain> Plain text message")
	in.E("br")
	var tags = make([]string, 0, len(message.RegisteredTypes))
	for tag := range message.RegisteredTypes {
		if tag != plaintext.Type.Tag && config.Get().MessageTypes[tag] != nil {
			tags = append(tags, tag)
		}
	}
	sort.Slice(tags, func(i, j int) bool {
		return capitalize(message.RegisteredTypes[tags[i]].Name) < capitalize(message.RegisteredTypes[tags[j]].Name)
	})
	for _, tag := range tags {
		in.E("input type=checkbox id=mtype.%s name=mtype value=%s", tag, tag, inList(session.MessageTypes, tag), "checked")
		in.E("label for=mtype.%s> %s", tag, capitalize(message.RegisteredTypes[tag].Name))
		in.E("br")
	}
	if err != "" {
		row.E("div class=formError>%s", err)
	}
	row.E("div class=formHelp>Message type(s) that are accepted as practice messages.  Note that plain text messages are always accepted when sent from somewhere other than W*XSC.")
}

func readPlainSubject(r *http.Request, session *store.Session, show bool) string {
	// Actually reads both subject and body, but only returns errors for the
	// subject.
	if !show {
		return ""
	}
	subject := strings.TrimSpace(r.FormValue("plainSubject"))
	body := strings.TrimSpace(removeCR.Replace(r.FormValue("plainBody")))
	session.ModelMessage = "Subject: " + subject + "\n\n" + body
	if subject == "" {
		return "The subject line is required."
	}
	if _, _, _, formtag, realsubj := common.DecodeSubject(subject); realsubj == subject || formtag != "" {
		return "This is not a standard SCCo subject line for a plain text message."
	}
	return ""
}
func emitPlainSubject(form *htmlb.Element, session *store.Session, show, focus bool, err string) {
	var subject string

	if show {
		if strings.HasPrefix(session.ModelMessage, "Subject: ") {
			subject, _, _ = strings.Cut(session.ModelMessage[9:], "\n")
		}
		if subject == "" {
			subject = "XXX-###P_"
		}
	}
	row := form.E("div id=plainSubjectRow class=formRow", !show, "style=display:none")
	row.E("label for=plainSubject>Subject Line")
	row.E("input id=plainSubject name=plainSubject value=%s", subject, focus, "autofocus")
	if err != "" {
		row.E("div class=formError>%s", err)
	}
	row.E("div class=formHelp>Subject line of the expected plain text message, with message number (ignored), handling order, and subject.")
}

func readPlainBody(session *store.Session, show bool) string {
	// Body was actually read by readPlainSubject; we just return errors for
	// it here.
	if !show {
		return ""
	}
	_, body, _ := strings.Cut(session.ModelMessage, "\n\n")
	if body == "" {
		return "The message body is required."
	}
	return ""
}
func emitPlainBody(form *htmlb.Element, session *store.Session, show, focus bool, err string) {
	var body string

	if show && strings.HasPrefix(session.ModelMessage, "Subject: ") {
		_, body, _ = strings.Cut(session.ModelMessage[9:], "\n\n")
	}
	row := form.E("div id=plainBodyRow class=formRow", !show, "style=display:none")
	row.E("label for=plainBody>Message Body")
	row.E("textarea id=plainBody name=plainBody rows=8", focus, "autofocus").T(body)
	if err != "" {
		row.E("div class=formError>%s", err)
	}
	row.E("div class=formHelp>Body of the expected plain text message.")
}

func readFormBody(r *http.Request, session *store.Session, show bool) string {
	if !show {
		return ""
	}
	body := strings.TrimSpace(removeCR.Replace(r.FormValue("formBody")))
	session.ModelMessage = "Subject: \n\n" + body
	if body == "" {
		return "The encoded form body is required."
	}
	form := message.Decode("", body)
	if _, ok := form.(*plaintext.PlainText); ok {
		return "This is not a valid PackItForms-encoded form."
	}
	if _, ok := form.(message.ICompare); !ok {
		return "This is not a form type for which comparison logic has been implemented."
	}
	return ""
}
func emitFormBody(form *htmlb.Element, session *store.Session, show, focus bool, err string) {
	var body string

	if show && strings.HasPrefix(session.ModelMessage, "Subject: ") {
		_, body, _ = strings.Cut(session.ModelMessage[9:], "\n\n")
	}
	row := form.E("div id=formBodyRow class=formRow", !show, "style=display:none")
	row.E("label for=formBody>Encoded Form")
	row.E("textarea id=formBody name=formBody rows=8", focus, "autofocus").T(body)
	if err != "" {
		row.E("div class=formError>%s", err)
	}
	row.E("div class=formHelp>PackItForms encoding of the expected form.  Note that the message number and the operator-only fields will be ignored.  Handling and destination fields can be left blank to require the sender to provide a correct value from the recommended routing cheat sheet.")
}

func (ws *webserver) readFormImage(r *http.Request, session *store.Session, show bool) ([]*multipart.FileHeader, string) {
	if !show {
		return nil, ""
	}
	if ws.st.ModelImageCount(session.ID) == 0 && len(r.MultipartForm.File["formImage"]) == 0 {
		return nil, "At least one form image is required."
	}
	return r.MultipartForm.File["formImage"], ""
}
func (ws *webserver) emitFormImage(form *htmlb.Element, session *store.Session, show, focus bool, err string) {
	row := form.E("div id=formImageRow class=formRow", !show, "style=display:none")
	count := ws.st.ModelImageCount(session.ID)
	row.E("label for=formImage", count != 0, "class=imagelabel").R("Form Image(s)")
	in := row.E("div class=formInput")
	if count != 0 {
		for pnum := 1; pnum <= count; pnum++ {
			in.E("img class=modelimage src=/session/image?session=%d&page=%d", session.ID, pnum)
		}
		in.E("br")
	} else {
		in.E("div>No form images on file yet.")
	}
	in.E("input type=file id=formImage name=formImage accept=image/*,.jpg,.jpeg,.png capture=environment multiple", focus, "autofocus")
	if err != "" {
		row.E("div class=formError>%s", err)
	}
	row.E("div class=formHelp>Image(s) of the expected form.  This is what operators will be shown to tell them what to send as a practice message.")
}

func readInstructions(r *http.Request, session *store.Session) {
	session.Instructions = strings.TrimSpace(removeCR.Replace(r.FormValue("instructions")))
}
func emitInstructions(form *htmlb.Element, session *store.Session) {
	row := form.E("div class=formRow")
	row.E("label for=instructions>Extra Instructions")
	row.E("textarea id=instructions name=instructions class=formInput").R(session.Instructions)
	row.E("div class=formHelp>HTML-encoded additional instructions for practice message senders.")
}

func emitButtons(form *htmlb.Element, session *store.Session) {
	buttons := form.E("div class=formButtons")
	buttons.E("input type=submit name=save class='sbtn sbtn-primary' value=Save", session.Report != "", "class=sbtn-disabled disabled")
	buttons.E("input type=submit name=copy class='sbtn sbtn-primary' value='Save Copy'")
	buttons.E("button type=button class='sbtn sbtn-secondary' onclick=history.back()>Cancel")
	if session.Flags&store.Running == 0 && session.Report == "" {
		buttons.E("div class=formButtonSpace")
		buttons.E("input type=submit name=delete class='sbtn sbtn-danger' value='Delete Session'")
	}
}

func inList[T comparable](list []T, item T) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
