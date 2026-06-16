// Package server contains the web server.
package server

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/rothskeller/packet/v4/cmd/packet/osdep"
	"github.com/rothskeller/packet/v4/errors"
	"github.com/rothskeller/packet/v4/form/htmlop"
	"github.com/rothskeller/packet/v4/incident"
	"github.com/rothskeller/packet/v4/message"
	"github.com/rothskeller/packet/v4/packetver"
)

const (
	// AddonVersion is the addon version number used in all messages.
	AddonVersion = "4.0"
	// pingTimeout specifies how long to wait for the server to respond to a
	// /ping request.
	pingTimeout = time.Second
	// attemptDelay specifies how long to wait between attempts to reach a
	// newly started server.
	attemptDelay = time.Second
	// waitAttempts specifies how many attempts to make to reach a newly
	// started server before giving up on it.
	waitAttempts = 5
)

// Logpath is the path to the log file, set by main.
var Logpath string

// GetAddress returns the URL of the currently running server.  If wait is true
// (i.e., we just started a server and are waiting for it to start responding to
// pings), GetAddress will try repeatedly for a while before returning.
func GetAddress(wait bool) (address string, err error) {
	maxAttempts := 1
	if wait {
		maxAttempts = waitAttempts
	}
	for range maxAttempts {
		if address, err = addressAttempt(); err != nil || address != "" {
			return address, err
		}
		time.Sleep(attemptDelay)
	}
	return "", nil
}

// addressAttempt tries to ping the server identified in the server
// address file, if any, and returns the server address if the ping was
// successful.  It returns an empty address and no error if there is no server
// address file or the addressed server is not pingable.  It returns an error if
// something went wrong that isn't just a "server not running".
func addressAttempt() (address string, err error) {
	var (
		ctx    context.Context
		cancel func()
		req    *http.Request
		resp   *http.Response
	)
	// Read the address file.
	if fh, err := os.Open(osdep.AddressFile); os.IsNotExist(err) {
		return "", nil
	} else if err != nil {
		slog.Error("os.Open", "f", osdep.AddressFile, "err", err)
		return "", errors.NewF("The packet server address file could not be opened: %s", err)
	} else if err = osdep.ReadLock(fh); err != nil {
		slog.Error("osdep.ReadLock", "f", osdep.AddressFile, "err", err)
		return "", errors.NewF("The packet server address file could not be locked: %s", err)
	} else if buf, err := io.ReadAll(fh); err != nil {
		slog.Error("io.ReadAll", "f", osdep.AddressFile, "err", err)
		osdep.Unlock(fh)
		fh.Close()
		return "", errors.NewF("The packet server address file could not be read: %s", err)
	} else {
		osdep.Unlock(fh)
		fh.Close()
		address = strings.TrimSpace(string(buf))
	}
	// Ping the server to see if it is still running.
	ctx, cancel = context.WithTimeout(context.Background(), pingTimeout)
	req, _ = http.NewRequestWithContext(ctx, http.MethodGet, address+"/ping", nil)
	resp, err = http.DefaultClient.Do(req)
	cancel()
	if err != nil {
		slog.Debug("ping server", "url", address, "err", err)
		return "", nil
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		slog.Error("ping server", "url", address, "code", resp.StatusCode, "status", resp.Status)
		return "", errors.NewF("The packet server responded to a ping request with an unexpected response %d.", resp.StatusCode)
	}
	slog.Debug("server is running", "url", address)
	return address, nil
}

// Start starts serving web requests on the specified port (or the default port
// or a random port if the specified port is zero).  If it finds that another
// server is already running (on any port), it returns nil immediately.
// Otherwise, it starts a new server.  The new server will stop, and the
// function will return, on receipt of a POST /stop request, on receipt of an
// interrupt signal, or on creation/change of /tmp/packet-stop (Windows:
// C:\PackItForms\stop).  The function returns an error only if a new server
// fails to start.
func Start(port int, outpost bool) (err error) {
	var (
		addressDir string
		addrFH     *os.File
		address    string
		listener   net.Listener
		server     Server
		hserver    http.Server
		defPort    bool
	)
	// Open and read the address file with a write lock.
	addressDir = filepath.Dir(osdep.AddressFile)
	if err = os.MkdirAll(addressDir, 0777); err != nil {
		slog.Error("os.MkdirAll", "d", addressDir, "err", err)
		return errors.New("unable to create server address file")
	} else if addrFH, err = os.OpenFile(osdep.AddressFile, os.O_RDWR|os.O_CREATE, 0644); err != nil {
		slog.Error("os.OpenFile", "f", osdep.AddressFile, "err", err)
		return errors.New("unable to create server address file")
	} else if err = osdep.WriteLock(addrFH); err != nil {
		slog.Error("WriteLock", "f", osdep.AddressFile, "err", err)
		addrFH.Close()
		return errors.New("unable to create server address file")
	} else if buf, err := io.ReadAll(addrFH); err != nil {
		slog.Error("io.ReadAll", "f", osdep.AddressFile, "err", err)
		osdep.Unlock(addrFH)
		addrFH.Close()
		return errors.New("unable to create server address file")
	} else {
		address = strings.TrimSpace(string(buf))
	}
	// Ping the address we read to see if there's a server running there.
	// (i.e., we were in a race with it to start up, and we lost.)
	if address != "" {
		ctx, cancel := context.WithTimeout(context.Background(), pingTimeout)
		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, address+"/ping", nil)
		resp, err := http.DefaultClient.Do(req)
		cancel()
		if err == nil {
			resp.Body.Close()
		}
		if err == nil && resp.StatusCode == http.StatusNoContent {
			// Yes, it's running happily.  Exit.
			slog.Debug("server already running", "url", address)
			fmt.Printf("A packet server is already running at %s .\n", address)
			osdep.Unlock(addrFH)
			addrFH.Close()
			return nil
		}
	}
	// If the caller didn't specify a port, try our default port first.
	if port == 0 {
		port, defPort = 45674, true
	}
	// Open a listening port.
	server.address = fmt.Sprintf("127.0.0.1:%d", port)
	if listener, err = net.Listen("tcp4", server.address); defPort && isBusyPortError(err) {
		// If that was an attempt at our default port and it failed
		// because the port was in use, use a random port instead.
		server.address = "127.0.0.1:0"
		listener, err = net.Listen("tcp4", server.address)
	}
	if err != nil {
		slog.Error("net.Listen", "a", server.address, "err", err)
		return err
	}
	server.address = "http://" + listener.Addr().String()
	// Send a stop signal when the packet-stop file is touched.
	server.stop = make(chan struct{}, 1)
	go server.watchForStopFile()
	// Send a stop signal when the user hits Ctrl-C or the window is closed.
	go server.watchForInterrupt()
	// Send a stop signal if Outpost is closed.
	if outpost {
		server.outpost = true
		go server.watchForOutpostClose()
	}
	// Set the server up to handle requests.
	server.registerHandlers()
	hserver.Handler = http.HandlerFunc(server.handleRequest)
	// Start a webserver on the port we're listening to.
	go hserver.Serve(listener)
	slog.Info("server listening", "url", server.address)
	fmt.Printf(`
PACKET SERVER v%s at %s
Keep this window open until all packet-related browser tabs are closed.

Activity is being logged to %s.
`, packetver.Version, server.address, Logpath)
	// Write the server address to the address file.
	addrFH.Seek(0, 0)
	addrFH.Truncate(0)
	fmt.Fprintln(addrFH, server.address)
	osdep.Unlock(addrFH)
	addrFH.Close()
	// Wait until (a) we've been idle for a long time, (b) the stop file is
	// touched, or (c) we've received a POST /stop request.
	<-server.stop
	// Remove the address file.  That saves anyone trying to send pings to
	// us after we're gone.
	os.Remove(osdep.AddressFile)
	// Shut down the server and exit.
	hserver.Shutdown(context.Background())
	slog.Info("server exited", "url", server.address)
	return nil
}

// Server represents the packet HTTP server.
type Server struct {
	address string
	stop    chan struct{}
	mux     http.ServeMux
	outpost bool
}

// watchForStopFile watches for the creation or update of the "packet-stop"
// file.  When that happens, it writes on the server's stop channel.
func (s *Server) watchForStopFile() {
	var (
		watcher *fsnotify.Watcher
		err     error
	)
	if watcher, err = fsnotify.NewWatcher(); err != nil {
		slog.Error("fsnotify.NewWatcher", "err", err)
		close(s.stop)
		return
	}
	if err = watcher.Add(filepath.Dir(osdep.ServerStopFile)); err != nil {
		slog.Error("watcher.Add", "f", osdep.ServerStopFile, "err", err)
		close(s.stop)
		return
	}
	for {
		select {
		case event := <-watcher.Events:
			if event.Name != osdep.ServerStopFile {
				break
			}
			slog.Info("stopping server: stop file modtime has changed")
			close(s.stop)
			return
		case err := <-watcher.Errors:
			slog.Error("watcher.Error", "err", err)
			close(s.stop)
			return
		}
	}
}

// watchForOutpostClose opens a connection to opdirect.  When the connection is
// closed (meaning Outpost is closed), it stops the server, unless the server
// has been used for non-Outpost requests.
func (s *Server) watchForOutpostClose() {
	var (
		conn net.Conn
		err  error
		buf  = make([]byte, 1)
	)
	if conn, err = net.Dial("tcp4", "127.0.0.1:9334"); err != nil {
		slog.Error("net.Dial (opdirect)", "err", err)
		close(s.stop)
		return
	}
	conn.Read(buf)
	if s.outpost {
		slog.Info("stopping server: connection to opdirect has been closed")
		close(s.stop)
	} else {
		slog.Info("opdirect connection closed; not stopping server due to non-Outpost usage")
	}
}

// watchForInterrupt watches for an interrupt signal to the server, and stops
// the server gracefully.
func (s *Server) watchForInterrupt() {
	var ch = make(chan os.Signal, 1)

	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch
	slog.Info("stopping server: interrupt signal")
	close(s.stop)
}

// handleRequest handles a web request to the server.
func (s *Server) handleRequest(w http.ResponseWriter, r *http.Request) {
	var attrs []slog.Attr

	// Log the request.
	r.FormValue("x") // force the form to be parsed
	for k, vs := range r.Form {
		for _, v := range vs {
			attrs = append(attrs, slog.String(k, v))
		}
	}
	slog.LogAttrs(context.Background(), slog.LevelDebug, r.Method+" "+r.URL.Path, attrs...)
	// Never cache any response.
	w.Header().Set("Cache-Control", "no-store, private")
	// Handle the request.
	s.mux.ServeHTTP(w, r)
}

// handleStop handles the POST /stop request by stopping the server.
func (s *Server) handleStop(w http.ResponseWriter, r *http.Request) {
	slog.Info("stopping server: received POST /stop")
	close(s.stop)
}

// registerHandlers registers the server handlers.
func (s *Server) registerHandlers() {
	s.mux.HandleFunc("POST /stop", s.handleStop)
	s.mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNoContent) })
	s.mux.HandleFunc("GET /outpost-new", s.outpostNewRequest)
	s.mux.HandleFunc("GET /outpost-edit", s.outpostEditRequest)
	s.mux.HandleFunc("POST /outpost-submit", s.outpostSubmit)
	s.mux.HandleFunc("GET /pdf/{filename}", s.serveRenderedPDF)
	s.mux.HandleFunc("GET /incident-config", s.serveGetIncidentConfig)
	s.mux.HandleFunc("POST /incident-config", s.servePostIncidentConfig)
	s.mux.HandleFunc("GET /incident", s.serveGetIncident)
	s.mux.HandleFunc("GET /incident/log", s.serveGetIncidentLog)
	s.mux.HandleFunc("POST /manual-receive", s.servePostManualReceive)
	s.mux.HandleFunc("GET /new-message", s.serveGetNewMessage)
	s.mux.HandleFunc("GET /edit-message", s.serveGetEditMessage)
	s.mux.HandleFunc("GET /assets/{tag}/{asset...}", s.serveGetAsset)
	s.mux.HandleFunc("POST /send-message", s.servePostSendMessage)
	s.mux.HandleFunc("GET /manual-send-command", s.serveGetManualSendCommand)
	s.mux.HandleFunc("POST /mark-sent", s.servePostMarkSent)
	s.mux.HandleFunc("GET /view-message", s.serveGetViewMessage)
	s.mux.HandleFunc("POST /print-message", s.servePostPrintMessage)
	s.mux.HandleFunc("POST /make-receipt", s.servePostMakeReceipt)
	s.mux.HandleFunc("POST /edit-log-entry", s.servePostEditLogEntry)
	s.mux.HandleFunc("POST /reset-log-entry", s.servePostResetLogEntry)
	s.mux.HandleFunc("POST /delete-log-entry", s.servePostDeleteLogEntry)
	s.mux.HandleFunc("POST /toggle-flag", s.servePostToggleFlag)
	s.mux.HandleFunc("POST /new-message-from", s.servePostNewMessageFrom)
	s.mux.HandleFunc("POST /delete-message", s.servePostDeleteMessage)
	s.mux.HandleFunc("GET /view-encoded", s.serveGetViewEncoded)
	s.mux.HandleFunc("POST /view-ics309", s.servePostViewICS309)
	s.mux.HandleFunc("GET /manpage.html", s.serveGetManPage)
	s.mux.HandleFunc("/incident-open", s.serveIncidentOpen)
	s.mux.HandleFunc("GET /county-seal.svg", s.serveGetCountySeal)
	s.mux.HandleFunc("GET /Go-Regular.woff2", s.serveGetGoRegular)
	s.mux.HandleFunc("GET /Go-Bold.woff2", s.serveGetGoBold)
	s.mux.HandleFunc("GET /Go-Italic.woff2", s.serveGetGoItalic)
	s.mux.HandleFunc("GET /Go-Mono.woff2", s.serveGetGoMono)
	s.mux.HandleFunc("GET /favicon.ico", s.serveGetFavicon)
	s.mux.HandleFunc("GET /favicon-16x16.png", s.serveGetFavicon16)
	s.mux.HandleFunc("GET /favicon-32x32.png", s.serveGetFavicon32)
	s.mux.HandleFunc("GET /apple-touch-icon.png", s.serveGetAppleTouchIcon)
	s.mux.HandleFunc("POST /connect-bbs", s.servePostConnectBBS)
	s.mux.HandleFunc("POST /connect-abort", s.servePostConnectAbort)
	s.mux.HandleFunc("GET /connect-progress", s.serveGetConnectProgress)
	s.mux.HandleFunc("POST /set-view-flag", s.servePostSetViewFlag)
	s.mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/incident-open", http.StatusSeeOther)
	})
}

//go:embed assets/manpage.html
var manpageHTML []byte

func (s *Server) serveGetManPage(w http.ResponseWriter, r *http.Request) {
	s.outpost = false
	if doc, err := htmlop.Parse(bytes.NewReader(manpageHTML)); err != nil {
		s.ErrPage(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		htmlop.Expand(doc, map[string]string{"VERSION": packetver.Version})
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		htmlop.Minify(w, doc)
	}
}

//go:embed assets/county-seal.svg
var countySealSVG []byte

func (s *Server) serveGetCountySeal(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml")
	w.Write(countySealSVG)
}

//go:embed assets/Go-Regular.woff2
var goRegularWOFF2 []byte

func (s *Server) serveGetGoRegular(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "font/woff2")
	w.Write(goRegularWOFF2)
}

//go:embed assets/Go-Bold.woff2
var goBoldWOFF2 []byte

func (s *Server) serveGetGoBold(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "font/woff2")
	w.Write(goBoldWOFF2)
}

//go:embed assets/Go-Italic.woff2
var goItalicWOFF2 []byte

func (s *Server) serveGetGoItalic(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "font/woff2")
	w.Write(goItalicWOFF2)
}

//go:embed assets/Go-Mono.woff2
var goMonoWOFF2 []byte

func (s *Server) serveGetGoMono(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "font/woff2")
	w.Write(goMonoWOFF2)
}

//go:embed assets/apple-touch-icon.png
var appleTouchIcon []byte

func (s *Server) serveGetAppleTouchIcon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/png")
	w.Write(appleTouchIcon)
}

//go:embed assets/favicon.ico
var favicon []byte

func (s *Server) serveGetFavicon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/vnd.microsoft.icon")
	w.Write(favicon)
}

//go:embed assets/favicon-16x16.png
var favicon16 []byte

func (s *Server) serveGetFavicon16(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/png")
	w.Write(favicon16)
}

//go:embed assets/favicon-32x32.png
var favicon32 []byte

func (s *Server) serveGetFavicon32(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/png")
	w.Write(favicon32)
}

// serveIncident is a helper function for handlers that take a dir= parameter
// identifying an incident directory.  If the supplied dir is a valid incident
// directory, act will be invoked with that incident locked (for read or write
// depending on the write flag).  If act is successful and finish is non-nil,
// finish will be invoked.  If finish is nil, it defaults to
// {w.WriteHeader(http.StatusNoContent); return nil}.  If any error occurs
// (invalid dir, act or finish return non-nil), serveIncident issues an error
// response to the client.  This will be either an HTML error page or a
// text/plain body, depending on the request's Accept header.
func (s *Server) serveIncident(w http.ResponseWriter, r *http.Request, write bool, act func(*incident.Incident) error, finish func() error) {
	var err error

	operate := func(i *incident.Incident) (err error) {
		if err = act(i); err != nil {
			return err
		}
		if finish == nil {
			w.WriteHeader(http.StatusNoContent)
			return nil
		}
		return finish()
	}
	if write {
		err = incident.Write(r.FormValue("dir"), operate)
	} else {
		err = incident.Read(r.FormValue("dir"), operate)
	}
	if err != nil {
		if strings.Contains(r.Header.Get("Accept"), "html") {
			s.ErrPage(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
}

// serveLogIdent is a helper function for handlers that take dir= and id=
// parameters identifying an incident directory and a log entry in the
// incident.  If the supplied dir is a valid incident directory, act will be
// invoked with that incident locked (for read or write depending on the write
// flag) and that log entry fetched.  If act is successful and finish is
// non-nil, finish will be invoked.  If finish is nil, it defaults to
// {w.WriteHeader(http.StatusNoContent); return nil}.  If any error occurs
// (invalid dir or id, or act or finish return non-nil), serveLogIdent issues
// an error response to the client.
func (s *Server) serveLogIdent(w http.ResponseWriter, r *http.Request, write bool, act func(*incident.Incident, *incident.LogEntry) error, finish func() error) {
	s.serveIncident(w, r, write, func(i *incident.Incident) (err error) {
		ident, _ := strconv.Atoi(r.FormValue("id"))
		if le := i.GetLogEntryByIdent(ident); le == nil {
			return errors.NewF("There is no message with LEID %q in this incident.", r.FormValue("id"))
		} else {
			return act(i, le)
		}
	}, finish)
}

// serveMessage is a helper function for handlers that take dir= and id=
// parameters identifying an incident directory and a log entry in the
// incident.  If the supplied dir is a valid incident directory, act will be
// invoked with that incident locked (for read or write depending on the write
// flag), that log entry fetched, and the message referred to by that log entry
// fetched.  If act is successful and finish is non-nil, finish will be
// invoked.  If finish is nil, it defaults to
// {w.WriteHeader(http.StatusNoContent); return nil}.  If any error occurs
// (invalid dir or id, or act or finish return non-nil), serveLogIdent issues
// an error response to the client.
func (s *Server) serveMessage(w http.ResponseWriter, r *http.Request, write bool, act func(*incident.Incident, *incident.LogEntry, message.Message) error, finish func() error) {
	s.serveLogIdent(w, r, write, func(i *incident.Incident, le *incident.LogEntry) (err error) {
		if msg, err := i.GetMessageFromLogEntry(le); msg == nil && err != nil {
			return err
		} else if msg == nil {
			return errors.New("There is no message associated with this log entry.")
		} else {
			return act(i, le, msg)
		}
	}, finish)
}

// Defining this here allows us to avoid putting the isBusyPortError function is
// OS-dependent code.
const WSAEADDRINUSE syscall.Errno = 10048

// isBusyPortError returns whether the error indicates an attempt to bind to a
// port already in use.
func isBusyPortError(err error) bool {
	switch err := err.(type) {
	case *net.OpError:
		switch e2 := err.Err.(type) {
		case *os.SyscallError:
			switch e3 := e2.Err.(type) {
			case syscall.Errno:
				return e3 == syscall.EADDRINUSE || e3 == /*windows.*/ WSAEADDRINUSE
			}
		}
	}
	return false
}
