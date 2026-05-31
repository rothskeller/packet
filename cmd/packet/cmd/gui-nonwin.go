//go:build !windows

package cmd

import (
	"sync"
)

// guiStartServer starts the packet server in an OS-appropriate way.
func guiStartServer() (address string, wg *sync.WaitGroup, err error) {
	return guiStartServerGoroutine()
}
