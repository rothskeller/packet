//go:build windows

package server

import (
	"context"
	_ "embed"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
	"unsafe"

	"github.com/rothskeller/packet/cmd/packet/osdep"
	"github.com/rothskeller/packet/errors"
	"github.com/rothskeller/packet/form/formdefs"
	"golang.org/x/sys/windows"
)

// OutpostGetAddress returns the URL of the currently running server.  If there
// is no server running, OutpostGetAddress invokes the server in a fashion
// appropriate for Outpost usage, waits for it to start up, and then returns its
// address.
func OutpostGetAddress() (address string, err error) {
	for attempts := range maxAttempts {
		if address, err = outpostAddressAttempt(); err != nil || address != "" {
			return address, err
		}
		if attempts == 0 {
			if err = outpostStartServer(); err != nil {
				return "", err
			}
		}
		time.Sleep(attemptDelay)
	}
	return "", errors.New("The packet server did not start successfully within the expected time.")
}

// outpostAddressAttempt tries to ping the server identified in the server
// address file, if any, and returns the server address if the ping was
// successful.  It returns an empty address and no error if there is no server
// address file or the addressed server is not pingable.  It returns an error if
// something went wrong that isn't just a "server not running".
func outpostAddressAttempt() (address string, err error) {
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

// outpostStartServer starts the server in a new, minimized console window.
func outpostStartServer() (err error) {
	// Before starting the server is the proper time to check for forms
	// updates.  Errors are logged but not returned.
	_ = formdefs.CheckForUpdates(false, true)
	// We can't use the Go standard library exec.Command.Start method,
	// because it doesn't allow access to the flags that would let us start
	// the process the way we want.  So we have to do it by calling the
	// Windows CreateProcess syscall directly.
	appName, _ := windows.UTF16PtrFromString(`C:\PackItForms\packet.exe`)
	commandLine, _ := windows.UTF16PtrFromString(`C:\PackItForms\packet.exe server start --outpost`)
	creationFlags := windows.CREATE_NEW_CONSOLE | windows.CREATE_NEW_PROCESS_GROUP
	startupInfo := new(windows.StartupInfo)
	startupInfo.Cb = uint32(unsafe.Sizeof(*startupInfo))
	startupInfo.Title, _ = windows.UTF16PtrFromString("Packet Server")
	startupInfo.XCountChars = 80
	startupInfo.YCountChars = 24
	startupInfo.Flags = 0x9 /* USECOUNTCHARS | USESHOWWINDOW */
	startupInfo.ShowWindow = windows.SW_SHOWMINIMIZED
	processInfo := new(windows.ProcessInformation)
	err = windows.CreateProcess(appName, commandLine, nil, nil, false, uint32(creationFlags), nil, nil, startupInfo, processInfo)
	if err != nil {
		return errors.NewF("The packet server could not be started: %s", err)
	}
	return nil
}
