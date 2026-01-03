package server

import (
	_ "embed"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/rothskeller/packet/incident"
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/address"
)

// servePostManualReceive handles POST /manual-receive requests, which contain
// a received message that was manually provided.
func (s *Server) servePostManualReceive(w http.ResponseWriter, r *http.Request) {
	var (
		makedr bool
		mtext  string
		msg    *message.JustReceivedMessage
		err    error
	)
	// Check parameters.
	makedr = r.FormValue("mrdr") != ""
	mtext = strings.ReplaceAll(r.FormValue("mrmsg"), "\r\n", "\n")
	if msg, err = message.NewJustReceivedMessage(mtext, r.FormValue("mrbbs"), ""); err != nil {
		slog.Error("NewJustReceivedMessage", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = incident.Write(r.FormValue("dir"), func(i *incident.Incident) error {
		var dr *message.DraftMessage

		if dr, err = i.ReceiveMessage(msg); err != nil {
			slog.Error("ReceiveMessage", "err", err)
			return err
		}
		if dr != nil && makedr {
			if drid, err := i.AddDraftMessage(dr, false); err != nil {
				return fmt.Errorf("queueing delivery receipt: %s", err)
			} else {
				w.Header().Set("X-Packet-Action", fmt.Sprintf("manual-send:%d", drid))
			}
		}
		return nil
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// serveGetManualSendCommand handles GET /manual-send-command requests.  These
// have dir= and id= parameters, with the incident directory and log entry
// ident number for a draft message in the incident.  The handler responds with
// a text/plain body with the JNOS command to send the message.
func (s *Server) serveGetManualSendCommand(w http.ResponseWriter, r *http.Request) {
	var (
		dir   string
		id    int
		msg   message.Message
		ident string
		to    []string
		cmd   strings.Builder
		eb    string
		err   error
	)
	// Get the message from the incident and make sure it's proper.
	dir = r.FormValue("dir")
	if id, err = strconv.Atoi(r.FormValue("id")); err != nil || id <= 0 {
		http.Error(w, fmt.Sprintf("invalid message ident %q", r.FormValue("id")), http.StatusBadRequest)
		return
	}
	err = incident.Read(dir, func(i *incident.Incident) (err error) {
		if i.Config.TacCall != "" {
			ident = i.Config.OpCall
		}
		if le := i.GetLogEntryByIdent(id); le != nil {
			msg, err = i.GetMessageFromLogEntry(le)
		} else {
			err = fmt.Errorf("invalid message ident %q", r.FormValue("id"))
		}
		return err
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if dm, ok := msg.(*message.DraftMessage); !ok {
		http.Error(w, "not an unsent outgoing message", http.StatusBadRequest)
		return
	} else {
		if addrs, err := address.ParseList(dm.To()); err != nil || len(addrs) == 0 {
			http.Error(w, "invalid or empty To: address list", http.StatusBadRequest)
			return
		} else {
			for _, addr := range addrs {
				to = append(to, addr.Address)
			}
		}
	}
	// TODO: validate the message and include any errors in the HTML
	// Generate the actual send command.
	if msg.Bulletin() {
		cmd.WriteString("SB ")
	} else if len(to) > 1 {
		cmd.WriteString("SC ")
	} else {
		cmd.WriteString("SP ")
	}
	cmd.WriteString(to[0])
	cmd.WriteByte('\n')
	if len(to) > 1 {
		cmd.WriteString(strings.Join(to[1:], ","))
		cmd.WriteByte('\n')
	}
	cmd.WriteString(msg.Subject().EncodedSubject())
	cmd.WriteByte('\n')
	eb = msg.Payload().Encode()
	cmd.WriteString(eb)
	if !strings.HasSuffix(eb, "\n") {
		cmd.WriteByte('\n')
	}
	cmd.WriteString("/EX\n")
	if ident != "" {
		fmt.Fprintf(&cmd, "# DE %s\n", ident)
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	io.WriteString(w, cmd.String())
}

// servePostMarkSent handles POST /mark-sent requests, which mark a message as
// having been sent (manually).  They have dir= and id= parameters with the
// incident directory and log entry ident number.  They return 204 No Content
// unless an error occurs.
func (s *Server) servePostMarkSent(w http.ResponseWriter, r *http.Request) {
	var (
		dir string
		id  int
		msg message.Message
		err error
	)
	// Get the message from the incident and make sure it's proper.
	dir = r.FormValue("dir")
	if id, err = strconv.Atoi(r.FormValue("id")); err != nil || id <= 0 {
		http.Error(w, fmt.Sprintf("invalid message ident %q", r.FormValue("id")), http.StatusBadRequest)
		return
	}
	err = incident.Write(dir, func(i *incident.Incident) (err error) {
		if le := i.GetLogEntryByIdent(id); le == nil {
			return fmt.Errorf("invalid message ident %q", r.FormValue("id"))
		} else if msg, err = i.GetMessageFromLogEntry(le); err != nil {
			return err
		} else if dm, ok := msg.(*message.DraftMessage); !ok {
			return errors.New("not an unsent outgoing message")
		} else {
			return i.MarkMessageSent(dm, le)
		}
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// servePostMakeReceipt handles POST /make-receipt requests, which create a
// delivery receipt for a received message.  It has dir= and id= parameters.
// It responds with an error status and text/plain error message, or 204 No
// Content if successful.  In the latter case, the X-Packet-Action header is
// set to "manual-send:$IDENT" where $IDENT is the ident of the new receipt.
func (s *Server) servePostMakeReceipt(w http.ResponseWriter, r *http.Request) {
	var (
		dir  string
		id   int
		drid int
		err  error
	)
	// Get the message from the incident and make sure it's proper.
	dir = r.FormValue("dir")
	if id, err = strconv.Atoi(r.FormValue("id")); err != nil || id <= 0 {
		http.Error(w, fmt.Sprintf("invalid message ident %q", r.FormValue("id")), http.StatusBadRequest)
		return
	}
	err = incident.Write(dir, func(i *incident.Incident) (err error) {
		if le := i.GetLogEntryByIdent(id); le == nil {
			return fmt.Errorf("invalid message ident %q", r.FormValue("id"))
		} else if msg, err := i.GetMessageFromLogEntry(le); err != nil {
			return err
		} else if dr, err := i.MakeDeliveryReceipt(msg, le); err != nil {
			return err
		} else if drid, err = i.AddDraftMessage(dr, false); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("X-Packet-Action", fmt.Sprintf("manual-send:%d", drid))
	w.WriteHeader(http.StatusNoContent)
}

/*
import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/rothskeller/packet/cmd/packet/osdep"
	"github.com/rothskeller/packet/forms"
	"gopkg.in/ini.v1"
)

const (
	LATIN1_ENCODING  = "Latin-1"
	UTF8_ENCODING    = "UTF-8"
	WINDOWS_ENCODING = "Windows-1252"
)

type manualData struct {
	folder string
}

type manualSettings struct {
	Encoding          string `json:"encoding"`
	ArchiveFolder     string `json:"archiveFolder"`
	Call              string `json:"call"`
	Name              string `json:"name"`
	Prefix            string `json:"prefix"`
	UseTac            bool   `json:"useTac"`
	TacCall           string `json:"tacCall"`
	TacName           string `json:"tacName"`
	TacPrefix         string `json:"tacPrefix"`
	OpCall            string `json:"opCall"`
	OpName            string `json:"opName"`
	OpPrefix          string `json:"opPrefix"`
	NextMessageNumber int    `json:"nextMessageNumber"`
}

//go:embed pifo-manual/assets
var manualAssets embed.FS

func (s *Server) serveManualAsset(w http.ResponseWriter, r *http.Request) {
	http.ServeFileFS(w, r, manualAssets, filepath.Join("pifo-manual", "assets", r.PathValue("asset")))
}

//go:embed pifo-manual/manual-setup.html
var manualSetupHTML string

// serveGetManualSetup handles GET /manual-setup requests.
func (s *Server) serveGetManualSetup(w http.ResponseWriter, r *http.Request) {
	var (
		set      *manualSettings
		formHTML string
	)
	// Read or create the manual-mode settings.
	set = s.getManualSettings()
	set.ArchiveFolder = html.EscapeString(set.ArchiveFolder)
	formHTML = expandVariables(manualSetupHTML, map[string]string{
		"addon_version":     AddonVersion,
		"encoding":          set.Encoding,
		"archiveFolder":     set.ArchiveFolder,
		"useTac":            strconv.FormatBool(set.UseTac),
		"tacCall":           set.TacCall,
		"tacName":           set.TacName,
		"tacPrefix":         set.TacPrefix,
		"opCall":            set.OpCall,
		"opName":            set.OpName,
		"opPrefix":          set.OpPrefix,
		"nextMessageNumber": strconv.Itoa(set.NextMessageNumber),
	})
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, formHTML)
}

// servePostManualSetup handles POST /manual-setup requests.
func (s *Server) servePostManualSetup(w http.ResponseWriter, r *http.Request) {
	var (
		set *manualSettings
		enc []byte
	)
	if archiveFolder := r.FormValue("archiveFolder"); archiveFolder != "" {
		if stat, err := os.Stat(archiveFolder); err != nil || !stat.IsDir() {
			ErrorPage(w, http.StatusBadRequest, fmt.Errorf("%s is not a folder", archiveFolder), 0)
			return
		}
		testfname := filepath.Join(archiveFolder, "test")
		err := os.WriteFile(testfname, []byte("test"), 0666)
		os.Remove(testfname)
		if err != nil {
			ErrorPage(w, http.StatusBadRequest, fmt.Errorf("%s is not writable", archiveFolder), nil)
			return
		}
	}
	set = s.getManualSettings()
	set.OpCall = r.FormValue("opCall")
	set.OpName = r.FormValue("opName")
	set.OpPrefix = r.FormValue("opPrefix")
	set.UseTac = r.FormValue("useTac") != ""
	set.TacCall = r.FormValue("tacCall")
	set.TacName = r.FormValue("tacName")
	set.TacPrefix = r.FormValue("tacPrefix")
	set.ArchiveFolder = r.FormValue("archiveFolder")
	if set.Encoding = normalizeEncoding(r.FormValue("encoding")); set.Encoding == "" {
		set.Encoding = WINDOWS_ENCODING
	}
	set.NextMessageNumber, _ = strconv.Atoi(r.FormValue("nextMessageNumber"))
	s.setManualSettings(set)
	enc, _ = json.Marshal(set)
	log.Printf("servePostManualSetup %s", string(enc))
	if nextPage := r.FormValue("nextPage"); nextPage != "" {
		http.Redirect(w, r, nextPage, http.StatusSeeOther)
	} else {
		sendWindowClose(w)
	}
}

//go:embed pifo-manual/manual.html
var manualHTML string

// serveGetManual handles GET /manual requests.
func (s *Server) serveGetManual(w http.ResponseWriter, r *http.Request) {
	var (
		lfs []*forms.LaunchForm
		sb  strings.Builder
		err error
	)
	if lfs, err = forms.All(); err != nil {
		ErrorPage(w, http.StatusInternalServerError, err, nil)
		return
	}
	for _, lf := range lfs {
		fmt.Fprintf(&sb, `<option value="%s">%s</option>`,
			html.EscapeString(lf.AddonName+"/"+lf.FormHTML),
			html.EscapeString(lf.DisplayName))
	}
	var formHTML = expandVariables(manualHTML, map[string]string{
		"addon_version": AddonVersion,
		"form_options":  sb.String(),
	})
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, formHTML)
}

// servePostManualCreate handles POST /manual-create requests.
func (s *Server) servePostManualCreate(w http.ResponseWriter, r *http.Request) {
	var (
		addon   string
		msgtype string
		set     *manualSettings
		f       url.Values
	)
	addon, msgtype, _ = strings.Cut(r.FormValue("ADDON_MSG_TYPE"), "/")
	set = s.getManualSettings()
	f = make(url.Values)
	f.Set("message_status", "manual")
	f.Set("addon", addon)
	f.Set("addonver", AddonVersion)
	f.Set("msgtype", msgtype)
	f.Set("msgID", s.nextManualMessageNumber(set))
	f.Set("opName", set.OpName)
	f.Set("opCall", set.OpCall)
	if set.UseTac {
		f.Set("tacName", set.TacName)
		f.Set("tacCall", set.TacCall)
	}
	s.editCommon(w, f, addon, msgtype, "", s.manualCreate)
}

// nextManualMessageNumber increments the next message number in the settings
// and returns the formatted message ID.
func (s *Server) nextManualMessageNumber(set *manualSettings) (msgID string) {
	var prefix string

	if set.UseTac {
		prefix = set.TacPrefix
	} else {
		prefix = set.OpPrefix
	}
	msgID = fmt.Sprintf("%s-%03dM", prefix, set.NextMessageNumber)
	set.NextMessageNumber++
	s.setManualSettings(set)
	return msgID
}

//go:embed pifo-manual/message.html
var messageHTML string

func (s *Server) manualCreate(w http.ResponseWriter, r *http.Request) {
	var (
		parsed   *ParsedEmail
		message  string
		subject  string
		formHTML string
	)
	parsed, message = ParseEmail(r.FormValue("formtext"))
	subject = r.FormValue("subject")
	if h := parsed.Fields["5."]; h == "I" || h == "IMMEDIATE" {
		parsed.Urgent = true
	}
	formHTML = expandVariables(messageHTML, map[string]string{
		"readOnly": "true",
		"Subject":  subject,
		"Urgent":   strconv.FormatBool(parsed.Urgent),
		"Message":  message,
	})
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, formHTML)
}

//go:embed pifo-manual/manual-error.html
var manualErrorHTML string

var embeddedEXRE = regexp.MustCompile(`(^|[\r\n])\/EX([\r\n]|$)`)
var bodyContainsFormRE = regexp.MustCompile(`^[^\r\n]*![^!\r\n]+![\r\n]`)

// servePostManualCommand handles POST /manual-command requests.  These are the
// submission of the message.html form with the "Log and Send" button.
func (s *Server) servePostManualCommand(w http.ResponseWriter, r *http.Request) {
	var (
		urgent    bool
		bulletin  bool
		subject   string
		message   string
		prefix    string
		suffix    string
		addresses []string
		command   string
		set       *manualSettings
		problems  []string
	)
	// First, build the JNOS command to send the message.
	message = r.FormValue("message")
	urgent = r.FormValue("urgent") != ""
	bulletin = r.FormValue("bulletin") != ""
	subject = r.FormValue("subject")
	if strings.HasSuffix(message, "\n") {
		suffix = "/EX\r\n"
	} else {
		suffix = "\r\n/EX\r\n"
	}
	for _, alist := range strings.Split(r.FormValue("to"), ",") {
		for _, address := range strings.Split(alist, ";") {
			if address = strings.TrimSpace(address); address != "" {
				addresses = append(addresses, AsciifyHeader(address))
			}
		}
	}
	switch len(addresses) {
	case 0:
		problems = append(problems, "It does not have a To: address.")
	case 1:
		if bulletin {
			prefix = "SB " + addresses[0]
		} else {
			prefix = "SP " + addresses[0]
		}
	default:
		prefix = "SC " + addresses[0] + "\r\n" + strings.Join(addresses[1:], ",")
	}
	prefix += "\r\n" + AsciifyHeader(subject) + "\r\n"
	if urgent {
		prefix += "!URG!"
	}
	command = prefix + message + suffix
	// Does it contain an embedded /EX line?
	if embeddedEXRE.MatchString(message) {
		problems = append(problems, "It contains /EX on a line by itself (the end-of-message marker).")
	}
	// Does it contain invalid characters?
	if strings.ContainsFunc(command, invalidCharacterFunc) {
		problems = append(problems, "It contains character(s) outside the supported character set.")
	}
	// If there are problems, say so.
	if len(problems) != 0 {
		for i := range problems {
			problems[i] = html.EscapeString(problems[i])
		}
		out := strings.Replace(manualErrorHTML, "{{problems}}", strings.Join(problems, "<br>"), 1)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		io.WriteString(w, out)
		return
	}
	// Update the PIFO manual saved data.
	set = s.getManualSettings()
	if err := s.archiveManualMessage(set, subject, command); err != nil {
		log.Printf("ERROR: unable to archive manual message: %s", err)
	}
	s.logManualSend(set, message, subject, addresses)
	// Emit the command as plain text.
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	io.WriteString(w, command)
}

// invalidCharacterFunc returns whether the specified character is invalid for
// use in a packet message.  At present, it rejects any character that is not
// plain ASCII.  Over time this may be relaxed when character set transcoding is
// implemented.
func invalidCharacterFunc(r rune) bool {
	return (r < 0x20 || r > 0x7E) && r != 10 && r != 13
}

// archiveManualMessage writes the manual message to a unique filename in the
// archive folder, if there is one.
func (s *Server) archiveManualMessage(set *manualSettings, subject, command string) (err error) {
	if set.ArchiveFolder == "" {
		return nil
	}
	var basename = filepath.Join(set.ArchiveFolder, strings.Map(toFileNameMap, subject))
	var filename = basename + ".txt"
	var suffix = 1
	for {
		if _, err := os.Stat(filename + ".txt"); os.IsNotExist(err) {
			break
		}
		suffix++
		filename = fmt.Sprintf("%s (%d).txt", basename, suffix)
	}
	return os.WriteFile(filename, []byte(command), 0666)
}

// toFileNameMap is a function to be passed to strings.Map that replaces all
// unsafe filename characters with tildes.
func toFileNameMap(r rune) rune {
	if strings.ContainsRune(`<>:"/\|?*`, r) {
		return '~'
	} else {
		return r
	}
}

type manualLog struct {
	Messages []*manualLogEntry `json:"messages"`
}
type manualLogEntry struct {
	Date       string `json:"date"`
	Time       string `json:"time"`
	ToCall     string `json:"toCall"`
	FromCall   string `json:"fromCall"`
	FromNumber string `json:"fromNumber"`
	Subject    string `json:"subject"`
}

func (s *Server) logManualSend(set *manualSettings, message, subject string, addresses []string) {
	data := s.readManualLog()
	parsed, _ := ParseEmail(message)
	fromNumber := parsed.Fields["MsgNo"]
	if fromNumber == "" {
		fromNumber = getMessageNumberFromSubject(subject)
	}
	for i, a := range addresses {
		var item manualLogEntry

		item.ToCall, _, _ = strings.Cut(a, "@")
		item.ToCall = strings.TrimSpace(item.ToCall)
		if i != 0 {
			item.Time = `"`
			item.FromCall = `"`
			item.FromNumber = `"`
		} else {
			item.Date = time.Now().Format("01/02/2006")
			item.Time = time.Now().Format("15:04")
			if set.UseTac {
				item.FromCall = set.TacCall
			} else {
				item.FromCall = set.OpCall
			}
			item.FromNumber = fromNumber
			item.Subject = subject
		}
		data.Messages = append(data.Messages, &item)
	}
	s.writeManualLog(data)
}

var messageNumberFromSubjectRE = regexp.MustCompile(`(?i)^((?:[A-Z0-9]{1,3}-)?\d+[A-Z]?)`)

func getMessageNumberFromSubject(s string) string {
	if match := messageNumberFromSubjectRE.FindStringSubmatch(s); match != nil {
		return match[1]
	}
	return ""
}

func (s *Server) readManualLog() (ml *manualLog) {
	ml = new(manualLog)
	fname := s.findManualLogFile()
	if data, err := os.ReadFile(fname); err != nil {
		log.Printf("ERROR: read %s: %s", fname, err)
	} else if err = json.Unmarshal(data, ml); err != nil {
		log.Printf("ERROR: parse %s: %s", fname, err)
	}
	return ml
}

func (s *Server) writeManualLog(ml *manualLog) {
	data, _ := json.Marshal(ml)
	fname := s.findManualLogFile()
	if err := os.WriteFile(fname, data, 0666); err != nil {
		log.Printf("ERROR: write %s: %s", fname, err)
	}
}

func sendWindowClose(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte("<script>window.close();</script>"))
}

func (s *Server) getManualSettings() (settings *manualSettings) {
	var (
		data []byte
		set  manualSettings
		err  error
	)
	if data, err = os.ReadFile(s.findManualSettingsFile()); err != nil {
		log.Printf("ERROR: reading settings: %s", err)
	} else if err = json.Unmarshal(data, &set); err != nil {
		log.Printf("ERROR: parsing settings: %s", err)
	} else {
		settings = &set
	}
	if settings == nil {
		settings = getInitialManualSettings()
		s.setManualSettings(settings)
	}
	if settings.Encoding != "" {
		settings.Encoding = normalizeEncoding(settings.Encoding)
	} else {
		settings.Encoding = WINDOWS_ENCODING
	}
	if settings.UseTac {
		settings.Call = settings.TacCall
		settings.Name = settings.TacName
		settings.Prefix = settings.TacPrefix
	} else {
		settings.Call = settings.OpCall
		settings.Name = settings.OpName
		settings.Prefix = settings.OpPrefix
	}
	return settings
}

func (s *Server) setManualSettings(set *manualSettings) {
	var data, _ = json.Marshal(set)

	if err := os.WriteFile(s.findManualSettingsFile(), data, 0666); err != nil {
		log.Printf("ERROR: writing settinsg file: %s", err)
	}
}

func getInitialManualSettings() (set *manualSettings) {
	set = &manualSettings{Encoding: WINDOWS_ENCODING, NextMessageNumber: 1}
	if data, err := ini.Load(`c:\Program Files (x86)\SCCo Packet\Outpost.conf`); err != nil {
		return set
	} else if sect := data.Section("DataDirectory"); sect == nil {
		return set
	} else if key := sect.Key("DataDir"); key == nil {
		return set
	} else if d2, err := ini.Load(filepath.Join(key.String(), "Outpost.profile")); err != nil {
		return set
	} else if s2 := d2.Section("IDENTIFICATION"); s2 == nil {
		return set
	} else {
		if k2 := s2.Key("UsrName"); k2 != nil {
			set.OpName = k2.String()
		}
		if k2 := s2.Key("UsrCall"); k2 != nil {
			set.OpCall = k2.String()
		}
		if k2 := s2.Key("UsrID"); k2 != nil {
			set.OpPrefix = k2.String()
		}
		if k2 := s2.Key("TacName"); k2 != nil {
			set.TacName = k2.String()
		}
		if k2 := s2.Key("TacCall"); k2 != nil {
			set.TacCall = k2.String()
		}
		if k2 := s2.Key("TacID"); k2 != nil {
			set.TacPrefix = k2.String()
		}
		if k2, k3 := s2.Key("ActMyCall"), s2.Key("ActMyName"); k2 != nil && k2.String() == set.TacCall && k3 != nil && k3.String() == set.TacName {
			set.UseTac = true
		}
	}
	return set
}

func (s *Server) findManualLogFile() string {
	return filepath.Join(s.findManualDataFolder(), "manual-log.json")
}

func (s *Server) findManualSettingsFile() string {
	return filepath.Join(s.findManualDataFolder(), "manual-settings.json")
}

func (s *Server) findManualDataFolder() string {
	var (
		appData  string
		testFile string
		err      error
	)
	if s.manual.folder != "" {
		return s.manual.folder
	}
	if appData = os.Getenv("APPDATA"); appData == "" {
		err = errors.New("%APPDATA% is not set")
		goto FALLBACK
	}
	s.manual.folder = filepath.Join(appData, "PackItForms")
	if err = os.MkdirAll(s.manual.folder, 0777); err != nil {
		err = fmt.Errorf("create %s: %w", s.manual.folder, err)
		goto FALLBACK
	}
	testFile = filepath.Join(s.manual.folder, "test.txt")
	err = os.WriteFile(testFile, []byte("test"), 0666)
	os.Remove(testFile)
	if err == nil {
		return s.manual.folder
	}
FALLBACK:
	s.manual.folder = osdep.LogsDir
	log.Printf("manualDataFolder = %s; %s", s.manual.folder, err)
	return s.manual.folder
}

var variableRE = regexp.MustCompile(`\{\{[A-Za-z_]+\}\}`)

func expandVariables(buf string, vars map[string]string) string {
	return variableRE.ReplaceAllStringFunc(buf, func(s string) string {
		return vars[s[2:len(s)-2]]
	})
}

func normalizeEncoding(e string) string {
	switch strings.ToLower(e) {
	case "cp1252", "cp-1252", "win1252", "win-1252", "windows-1252":
		return WINDOWS_ENCODING
	case "cp819", "cp-819", "latin1", "latin-1", "iso8859-1", "iso 8859-1", "iso-8859-1", "iso_8859-1":
		return LATIN1_ENCODING
	case "utf8", "utf-8":
		return UTF8_ENCODING
	default:
		return e
	}
}
*/
