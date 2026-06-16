// Package incident handles collections of messages stored in the same
// directory.  Semantically, these are messages that belong to the same
// incident, i.e., that are recorded on the same ICS-309 form.
//
// Incidents are identified by the full absolute path of their directory.
// At a high level, there are only three operations one can perform:  Read,
// Write, and Watch.  These are the three exported functions in this package.
package incident

import (
	"context"
	"encoding/json/v2"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/rothskeller/packet/v4/cmd/packet/osdep"
)

const (
	lockFileName = ".packet.lock"
	jsonFileName = "packet.json"
)

var (
	ErrNotIncident = errors.New("not a valid incident directory")
)

var (
	watcher      *fsnotify.Watcher
	watching     = map[string]chan struct{}{}
	watcherMutex sync.Mutex
)

// An Incident is a collection of related messages stored in the same
// directory.
type Incident struct {
	// Dir is the absolute path to the directory for the incident.
	Dir string `json:"-"`
	// Seq is the sequence number, incremented with every change to the
	// incident.
	Seq int `json:"seq"`
	// Config is the configuration of the incident, i.e., the details that
	// (usually) don't change over time.
	Config *Config `json:"conf"`
	// BulletinChecks is a map from bulletin area name to the time at which
	// it was last checked.
	BulletinChecks map[string]time.Time `json:"bull"`
	// Log is an ordered list of log messages for the incident.  The first
	// element is always nil so that the slice index equals the item Index.
	Log []*LogEntry `json:"log"`
}

// IsIncident returns whether the named directory is an incident directory.
// (This is mostly used when presenting the user with a list of past incidents,
// to weed out past incident directories that no longer exist or no longer
// contain incident data.)
func IsIncident(dir string) bool {
	if !filepath.IsAbs(dir) {
		return false
	}
	if _, err := os.Stat(filepath.Join(dir, lockFileName)); err != nil {
		return false
	}
	if _, err := os.Stat(filepath.Join(dir, jsonFileName)); err != nil {
		return false
	}
	return true
}

// IsUnsafeIncidentDir returns whether the named directory is unsafe for use as
// an incident directory.  Unsafe directories are the root directory of any
// volume, the user's home directory, the Desktop or Documents subdirectory of
// the user's home directory, or (on Windows) C:\PackItForms.  These are
// considered unsafe because they are semantically too global to keep a single
// incident in.
func IsUnsafeIncidentDir(dir string) bool {
	if dir == "" || dir == "/" || dir[1:] == `:\` || strings.EqualFold(dir[1:], `:\PackItForms`) {
		return true
	}
	home := os.Getenv("HOME")
	if home == "" {
		return false
	}
	// I'm doing case insensitive comparisons because Windows and Mac both
	// have case insensitive pathnames.  If it causes a false positive on
	// Linux, too bad.
	home = strings.ToLower(home)
	dir = strings.ToLower(dir)
	if home == dir {
		return true
	}
	home += string(filepath.Separator)
	if !strings.HasPrefix(dir, home) {
		return false
	}
	dir = strings.TrimPrefix(dir, home)
	return dir == "desktop" || dir == "documents"
}

// Create creates a new incident in the named directory, which must not already
// be an incident directory.  It calls the supplied function while holding the
// write lock on the new incident.  If the supplied function returns nil, the
// new incident is saved; otherwise, no change is made to the named directory
// and the error is returned.
func Create(dir string, fn func(*Incident) error) (err error) {
	var (
		lockFH    *os.File
		inc       *Incident
		stateFile string
	)
	if !filepath.IsAbs(dir) {
		return ErrNotIncident
	}
	if lockFH, err = os.OpenFile(filepath.Join(dir, lockFileName), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0666); err != nil {
		return err
	}
	if err = osdep.WriteLock(lockFH); err != nil {
		lockFH.Close()
		os.Remove(lockFH.Name())
		return fmt.Errorf("write lock %s: %s", lockFH.Name(), err)
	}
	inc = &Incident{Dir: dir, Seq: 1, Config: configFromDefaults()}
	if err = fn(inc); err != nil {
		osdep.Unlock(lockFH)
		lockFH.Close()
		os.Remove(lockFH.Name())
		return err
	}
	if err = inc.writeStateFile(); err != nil {
		osdep.Unlock(lockFH)
		lockFH.Close()
		os.Remove(lockFH.Name())
		return err
	}
	stateFile = filepath.Join(dir, jsonFileName)
	if _, err = fmt.Fprintf(lockFH, "%08d", inc.Seq); err != nil {
		os.Remove(stateFile)
		osdep.Unlock(lockFH)
		lockFH.Close()
		os.Remove(lockFH.Name())
		return fmt.Errorf("write %s: %s", lockFH.Name(), err)
	}
	if err = osdep.Unlock(lockFH); err != nil {
		os.Remove(stateFile)
		osdep.Unlock(lockFH)
		lockFH.Close()
		os.Remove(lockFH.Name())
		return fmt.Errorf("unlock %s: %s", lockFH.Name(), err)
	}
	if err = lockFH.Close(); err != nil {
		os.Remove(stateFile)
		osdep.Unlock(lockFH)
		lockFH.Close()
		os.Remove(lockFH.Name())
		return fmt.Errorf("close %s: %s", lockFH.Name(), err)
	}
	return nil
}

// Read opens the incident in the named directory for read-only operations and
// invokes the supplied function while holding the read lock.  Read returns any
// error it encounters or any error returned by the supplied function.
func Read(dir string, fn func(*Incident) error) (err error) {
	var (
		lockFH    *os.File
		inc       *Incident
		stateFile string
	)
	if !filepath.IsAbs(dir) {
		slog.Error("not an incident directory", "dir", dir)
		return ErrNotIncident
	}
	if lockFH, err = os.Open(filepath.Join(dir, lockFileName)); err != nil {
		slog.Error("open lock file", "lf", filepath.Join(dir, lockFileName), "err", err)
		return err
	}
	defer lockFH.Close()
	if err = osdep.ReadLock(lockFH); err != nil {
		slog.Error("lock for read", "lf", lockFH.Name(), "err", err)
		return fmt.Errorf("read lock %s: %s", lockFH.Name(), err)
	}
	defer osdep.Unlock(lockFH)
	stateFile = filepath.Join(dir, jsonFileName)
	if inc, err = readStateFile(stateFile); err != nil {
		return err
	}
	inc.Dir = dir
	if err = fn(inc); err != nil {
		return err
	}
	return nil
}

// Write opens the incident in the named directory for write operations and
// invokes the supplied function while holding the write lock.  If Write
// encounters any error, it will return it.  If the supplied function returns
// nil, any changes it made to the Incident will be applied.  If it returns an
// error, no changes will be applied and Write will return the error.
func Write(dir string, fn func(*Incident) error) (err error) {
	var (
		lockFH    *os.File
		inc       *Incident
		stateFile string
	)
	if !filepath.IsAbs(dir) {
		slog.Error("not an incident directory", "dir", dir)
		return ErrNotIncident
	}
	if lockFH, err = os.OpenFile(filepath.Join(dir, lockFileName), os.O_RDWR, 0666); err != nil {
		slog.Error("open lock file", "lf", filepath.Join(dir, lockFileName), "err", err)
		return err
	}
	defer lockFH.Close()
	if err = osdep.WriteLock(lockFH); err != nil {
		slog.Error("lock for write", "lf", lockFH.Name(), "err", err)
		return fmt.Errorf("write lock %s: %s", lockFH.Name(), err)
	}
	defer osdep.Unlock(lockFH)
	stateFile = filepath.Join(dir, jsonFileName)
	if inc, err = readStateFile(stateFile); err != nil {
		return err
	}
	inc.Dir = dir
	inc.Seq++
	if err = fn(inc); err != nil {
		return err
	}
	// Remove the generated ICS-309 file after any successful incident
	// change.
	os.Remove(filepath.Join(inc.Dir, "ics309.pdf"))
	if err = inc.writeStateFile(); err != nil {
		return err
	}
	if _, err = fmt.Fprintf(lockFH, "%08d", inc.Seq); err != nil {
		slog.Error("write to lock file", "lf", lockFH.Name(), "err", err)
		return fmt.Errorf("write %s: %s", lockFH.Name(), err)
	}
	if err = osdep.Unlock(lockFH); err != nil {
		slog.Error("unlock", "lf", lockFH.Name(), "err", err)
		return fmt.Errorf("unlock %s: %s", lockFH.Name(), err)
	}
	if err = lockFH.Close(); err != nil {
		slog.Error("close lock file", "lf", lockFH.Name(), "err", err)
		return fmt.Errorf("close %s: %s", lockFH.Name(), err)
	}
	return nil
}

// Watch monitors the incident in the named directory for changes.  It returns
// nil when a change is observed or the stop channel (if any) is closed, and a
// context error when the supplied context is canceled.  If seq is less than the
// current sequence number of the incident, Watch returns nil immediately.
func Watch(ctx context.Context, dir string, seq int, stop <-chan struct{}) (err error) {
	var (
		lockFile string
		lockFH   *os.File
		curr     int
		ch       chan struct{}
	)
	if !IsIncident(dir) {
		slog.Error("not an incident directory", "dir", dir)
		return ErrNotIncident
	}
	// Open and lock the lock file.
	lockFile = filepath.Join(dir, lockFileName)
	if lockFH, err = os.Open(lockFile); err != nil {
		slog.Error("open lock file", "lf", lockFile, "err", err)
		return err
	}
	defer lockFH.Close()
	if err = osdep.ReadLock(lockFH); err != nil {
		slog.Error("read lock file", "lf", lockFile, "err", err)
		return fmt.Errorf("read lock %s: %s", lockFile, err)
	}
	defer osdep.Unlock(lockFH)
	// Read the sequence number from the lock file.  If it's old, the
	// incident has already changed, and we return immediately.
	if _, err = fmt.Fscanf(lockFH, "%d", &curr); err == nil && curr > seq {
		return nil
	}
	// Make sure the watcher is created and running.
	watcherMutex.Lock()
	if watcher == nil {
		if watcher, err = fsnotify.NewWatcher(); err != nil {
			slog.Error("fsnotify.NewWatcher", "err", err)
			watcherMutex.Unlock()
			return fmt.Errorf("can't create watcher: %s", err)
		}
		go listenToWatcher()
	}
	// If there's already a channel for writes to this file, we'll use it.
	if ch = watching[lockFile]; ch == nil {
		ch = make(chan struct{})
		watching[lockFile] = ch
		// If not, we'll add this file to the watcher list and create
		// a channel for it.
		if err = watcher.Add(lockFile); err != nil {
			slog.Error("watcher.Add", "f", lockFile, "err", err)
			watcherMutex.Unlock()
			return fmt.Errorf("can't watch %s: %s", lockFile, err)
		}
	}
	watcherMutex.Unlock()
	// Now that we're watching, we can release the lock on the lock file.
	osdep.Unlock(lockFH)
	lockFH.Close()
	// Wait for the channel to be closed, indicating a write, or for the
	// supplied context to be canceled.
	select {
	case <-stop:
		return nil
	case <-ch:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// listenToWatcher is the goroutine that handles events from the watcher.
// If it gets a WRITE event for a file we're watching, it closes the
// corresponding channel, which will wake up all other goroutines waiting on
// it.
func listenToWatcher() {
	for {
		select {
		case ev := <-watcher.Events:
			if ev.Has(fsnotify.Write) {
				watcherMutex.Lock()
				if ch := watching[ev.Name]; ch != nil {
					close(ch)
					delete(watching, ev.Name)
				}
				watcher.Remove(ev.Name)
				watcherMutex.Unlock()
			}
		case err := <-watcher.Errors:
			slog.Error("watcher.Errors", "err", err)
		}
	}
}

// readStateFile reads the incident state from the specified file.
func readStateFile(fname string) (inc *Incident, err error) {
	var fh *os.File

	if fh, err = os.Open(fname); err != nil {
		slog.Error("os.Open", "f", fname, "err", err)
		return nil, err
	}
	defer fh.Close()
	inc = new(Incident)
	if err = json.UnmarshalRead(fh, inc, json.RejectUnknownMembers(true)); err != nil {
		slog.Error("json.UnmarshalRead", "f", fname, "err", err)
		return nil, fmt.Errorf("parse %s: %s", fh.Name(), err)
	}
	if inc.Config == nil {
		inc.Config = new(Config)
	}
	return inc, nil
}

// writeStateFile writes out the new state file for an incident.
func (inc *Incident) writeStateFile() (err error) {
	var (
		fname string
		fh    *os.File
	)
	fname = filepath.Join(inc.Dir, jsonFileName)
	if fh, err = os.Create(fname); err != nil {
		slog.Error("os.Create", "f", fname, "err", err)
		return err
	}
	if err = json.MarshalWrite(fh, inc); err != nil {
		slog.Error("json.MarshalWrite", "f", fname, "err", err)
		fh.Close()
		os.Remove(fname) // better gone than corrupt?
		return fmt.Errorf("write %s: %s", fname, err)
	}
	if err = fh.Close(); err != nil {
		slog.Error("os.Close", "f", fname, "err", err)
		os.Remove(fname)
		return fmt.Errorf("close %s: %s", fname, err)
	}
	return nil
}
