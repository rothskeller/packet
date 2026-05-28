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
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/rothskeller/packet/cmd/packet/osdep"
	"github.com/rothskeller/packet/errors"
	"github.com/rothskeller/packet/form/formdefs"
	"github.com/rothskeller/packet/form/htmlop"
	"github.com/rothskeller/packet/incident"
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/packetver"
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
	// maxAttempts specifies how many attempts to make to reach a newly
	// started server before giving up on it.
	maxAttempts = 5
	// serverTimeout specifies how long the server can run without receiving
	// any requests.
	serverTimeout = time.Hour
)

// GetAddress returns the URL of the currently running server.  If there is no
// server running and start is true, GetAddress invokes the server, waits for it
// to start up, and then returns its address.
func GetAddress(start bool) (address string, err error) {
	var (
		attempts     int
		ctx          context.Context
		cancel       func()
		req          *http.Request
		resp         *http.Response
		cmd          *exec.Cmd
		attemptLimit = 1
	)
	for attempts < attemptLimit {
		attempts++
		// Read the address file.
		if fh, err := os.Open(osdep.AddressFile); os.IsNotExist(err) {
			err = nil
			goto START
		} else if err != nil {
			slog.Error("os.Open", "f", osdep.AddressFile, "err", err)
			return "", err
		} else if err = osdep.ReadLock(fh); err != nil {
			slog.Error("osdep.ReadLock", "f", osdep.AddressFile, "err", err)
			return "", fmt.Errorf("locking %s: %w", osdep.AddressFile, err)
		} else if buf, err := io.ReadAll(fh); err != nil {
			slog.Error("io.ReadAll", "f", osdep.AddressFile, "err", err)
			osdep.Unlock(fh)
			fh.Close()
			return "", fmt.Errorf("reading %s: %w", osdep.AddressFile, err)
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
			err = fmt.Errorf("%s/ping: %w", address, err)
		} else {
			resp.Body.Close()
			if resp.StatusCode == http.StatusNoContent {
				// We pinged the server successfully, so we have
				// the correct address.
				slog.Debug("server is running", "url", address)
				return address, nil
			} else {
				slog.Warn("ping server", "url", address, "code", resp.StatusCode, "status", resp.Status)
				err = fmt.Errorf("%s/ping responded with %d %s", address, resp.StatusCode, resp.Status)
			}
		}
	START:
		// We weren't able to contact a running server.
		if start {
			// We should try to start one.  But first, just before
			// starting the server is the proper time to check for
			// forms updates.  Errors will get logged but need not
			// be returned.
			_ = formdefs.CheckForUpdates(false, true)
			// Now that that's done, start the server.
			cmd = exec.Command(os.Args[0], "server", "start")
			cmd.SysProcAttr = osdep.DetachChild
			if err = cmd.Start(); err != nil {
				slog.Error("exec server start", "err", err)
				return "", fmt.Errorf("start %s serve: %w", os.Args[0], err)
			} else {
				slog.Debug("started server")
			}
			// Don't try again, but do wait a while for it to start up.
			start = false
			attemptLimit = maxAttempts
		}
		if attempts < attemptLimit {
			time.Sleep(attemptDelay)
		}
	}
	slog.Debug("server is not running")
	return "", err
}

// Start starts serving web requests on the specified port (or a random port if
// the specified port is zero).  If it finds that another server is already
// running (on any port), it writes that server's address to writeURL (if not
// nil) returns nil immediately.  Otherwise, it starts a new server and writes
// its address to writeURL (if not nil).  The new server will stop, and the
// function will return, after an hour of inactivity, on receipt of a POST /stop
// request, on receipt of an interrupt signal, or on creation/change of
// /tmp/packet-stop (Windows: C:\PackItForms\stop).  The function returns an
// error only if a new server fails to start.
func Start(port int, writeURL io.Writer) (err error) {
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
			if writeURL != nil {
				fmt.Fprintln(writeURL, address)
			}
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
	// Send a stop signal with the idle timer fires.
	server.idleTimer = time.NewTimer(serverTimeout)
	// Send a stop signal when the user hits Ctrl-C.
	go server.watchForInterrupt()
	// Set the server up to handle requests.
	server.registerHandlers()
	hserver.Handler = http.HandlerFunc(server.handleRequest)
	// Start a webserver on the port we're listening to.
	go hserver.Serve(listener)
	slog.Info("server listening", "url", server.address)
	if writeURL != nil {
		fmt.Fprintln(writeURL, server.address)
	}
	// Write the server address to the address file.
	addrFH.Seek(0, 0)
	addrFH.Truncate(0)
	fmt.Fprintln(addrFH, server.address)
	osdep.Unlock(addrFH)
	addrFH.Close()
	// Send a stop signal when we haven't received any server requests in a
	// long time.
	go server.watchForIdleTimeout()
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
	address   string
	stop      chan struct{}
	idleTimer *time.Timer
	mux       http.ServeMux
	// manual    manualData
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

// watchForInterrupt watches for an interrupt signal to the server, and stops
// the server gracefully.
func (s *Server) watchForInterrupt() {
	var ch = make(chan os.Signal, 1)

	signal.Notify(ch, os.Interrupt)
	<-ch
	slog.Info("stopping server: interrupt signal")
	close(s.stop)
}

// watchForIdleTimeout watches for an idle timeout, and stops the server when it
// happens.
func (s *Server) watchForIdleTimeout() {
	<-s.idleTimer.C
	slog.Info("stopping server: inactivity timeout")
	close(s.stop)
}

// handleRequest handles a web request to the server.
func (s *Server) handleRequest(w http.ResponseWriter, r *http.Request) {
	var attrs []slog.Attr

	// Reset the idle timer on any request.
	s.idleTimer.Reset(serverTimeout)
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
	if doc, err := htmlop.Parse(bytes.NewReader(manpageHTML)); err != nil {
		ErrPage(w, err.Error(), http.StatusInternalServerError)
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
func serveIncident(w http.ResponseWriter, r *http.Request, write bool, act func(*incident.Incident) error, finish func() error) {
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
			ErrPage(w, err.Error(), http.StatusBadRequest)
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
func serveLogIdent(w http.ResponseWriter, r *http.Request, write bool, act func(*incident.Incident, *incident.LogEntry) error, finish func() error) {
	serveIncident(w, r, write, func(i *incident.Incident) (err error) {
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
func serveMessage(w http.ResponseWriter, r *http.Request, write bool, act func(*incident.Incident, *incident.LogEntry, message.Message) error, finish func() error) {
	serveLogIdent(w, r, write, func(i *incident.Incident, le *incident.LogEntry) (err error) {
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
