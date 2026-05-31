//go:build windows

package cmd

import (
	"sync"
	"unsafe"

	"github.com/rothskeller/packet/cmd/packet/server"
	"github.com/rothskeller/packet/errors"
	"golang.org/x/sys/windows"
)

var (
	kernel32             = windows.NewLazySystemDLL("kernel32.dll")
	procGetConsoleWindow = kernel32.NewProc("GetConsoleWindow")
)

// guiStartServer starts the packet server in an OS-appropriate way.
func guiStartServer() (address string, wg *sync.WaitGroup, err error) {
	if c, _, _ := procGetConsoleWindow.Call(); c != 0 {
		// We have a console window, so we'll run the server as a
		// goroutine in our own process.
		return guiStartServerGoroutine()
	}
	// We don't have a console window, so we'll start the server in its own
	// minimized console window.
	address, err = guiStartServerProcess()
	return address, nil, err
}

// guiStartServerProcess starts the server in a separate process with its own
// minimized console window.
func guiStartServerProcess() (address string, err error) {
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
		return "", errors.NewF("The packet server could not be started: %s", err)
	}
	return server.GetAddress(true)
}
