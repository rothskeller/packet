// Package formdefs locates and reads all form definitions.
package formdefs

import (
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/rothskeller/packet/errors"
	"github.com/rothskeller/packet/form"
	"github.com/rothskeller/packet/form/formdef"
	"github.com/rothskeller/packet/forms"
	"github.com/rothskeller/packet/message"
)

type FormsFSI interface {
	fs.FS
	fs.ReadDirFS
	fs.ReadFileFS
}

var (
	ErrNoAppDir  = errors.New("The program could not locate the directory where forms should be installed.  It will use embedded forms, which may not be current.")
	ErrNoForms   = errors.New("No forms were found in the forms directory, so no form message types are defined.")
	ErrNoVersion = errors.New("The program could not determine its own version number.  It will use embedded forms, which may not be current.")
)

var once sync.Once
var FormsFS FormsFSI

// UseInternalForms can be set to true prior to calling RegisterForms, by
// clients that want to force usage of the embedded forms.
var UseInternalForms bool

// RegisterForms locates all form definitions, performs any necessary updates,
// and registers message types for all known forms.  The returned error gives
// any problems; they are always non-fatal.
func RegisterForms() (err error) {
	once.Do(func() { err = registerForms() })
	return err
}
func registerForms() (err error) {
	if UseInternalForms {
		FormsFS = forms.EmbeddedForms
	} else {
		FormsFS, err = getFormsFileSystem()
	}
	err = errors.Join(err, registerFSForms())
	return err
}

func getFormsFileSystem() (formsFS FormsFSI, err error) {
	var (
		dir  string
		ents []fs.DirEntry
	)
	// Figure out where the forms should go.
	if dir = FormsDir(); dir == "" {
		// Couldn't determine local cache location, so use embedded
		// forms.
		slog.Warn("couldn't determine FormsDir")
		return forms.EmbeddedForms, ErrNoAppDir
	}
	// Make sure that place exists.
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// The local cache directory doesn't exist.  Try to create it.
		// If we can't, use the embedded forms.
		if err = os.MkdirAll(dir, 0777); err != nil {
			slog.Warn("os.MkdirAll", "d", dir, "err", err)
			return forms.EmbeddedForms, errors.NewF("The program could not create the forms directory (%s).  It will use embedded forms, which may not be current.", err)
		} else {
			slog.Info("created Forms directory", "d", dir)
		}
	} else if err != nil {
		slog.Warn("os.Stat", "d", dir, "err", err)
		return forms.EmbeddedForms, errors.NewF("The program could not access the forms directory (%s).  It will use embedded forms, which may not be current.", err)
	}
	// Make sure that every bundle in the embedded forms exists in the
	// forms directory.
	ents, _ = forms.EmbeddedForms.ReadDir(".")
	for _, bundle := range ents {
		target := filepath.Join(dir, bundle.Name())
		if _, err := os.Stat(target); os.IsNotExist(err) {
			sub, _ := fs.Sub(forms.EmbeddedForms, bundle.Name())
			if err = os.CopyFS(target, sub); err != nil {
				slog.Warn("os.CopyFS", "bundle", bundle.Name(), "dest", target, "err", err)
				return forms.EmbeddedForms, fmt.Errorf("The program could not install the embedded %s forms in %s.  (Error: %s.)  It will use the embedded forms, which may not be current.", bundle.Name(), dir, err)
			} else {
				slog.Info("installed embedded forms", "bundle", bundle.Name())
			}
		} else {
			slog.Debug("bundle exists, not copying", "bundle", bundle.Name())
		}
	}
	return os.DirFS(dir).(FormsFSI), err
}

func registerFSForms() (err error) {
	var (
		forms   []string
		highest = make(map[string]map[string]map[int]*form.FormType)
	)
	fs.WalkDir(FormsFS, ".", func(path string, d fs.DirEntry, werr error) error {
		err = errors.Join(err, werr)
		if strings.HasSuffix(path, ".form") && !d.IsDir() {
			forms = append(forms, path)
		}
		return nil
	})
	if len(forms) == 0 {
		slog.Warn("no forms found")
		return ErrNoForms
	}
	for _, ff := range forms {
		if def, derr := formdef.ReadFS(FormsFS, ff); derr != nil {
			slog.Warn("form definition error", "f", ff, "err", derr)
			err = errors.Join(err, derr)
		} else {
			var ft = form.FormType{FormDef: def}
			var mt = message.MType(ft)
			if def.CreateTag != "" {
				mt = form.EditableFormType{FormType: &ft}
			}
			if derr = message.RegisterType(mt); derr != nil {
				slog.Warn("form registration error", "f", ff, "err", derr)
				err = errors.Join(err, derr)
			}
			// Keep track of the highest minor version number for
			// each major version of each form.
			if match := form.VersionRE.FindStringSubmatch(def.Version); match != nil {
				if highest[def.AddonName] == nil {
					highest[def.AddonName] = make(map[string]map[int]*form.FormType)
				}
				if highest[def.AddonName][def.HTMLName] == nil {
					highest[def.AddonName][def.HTMLName] = make(map[int]*form.FormType)
				}
				major, _ := strconv.Atoi(match[1])
				if highest[def.AddonName][def.HTMLName][major] == nil {
					highest[def.AddonName][def.HTMLName][major] = &ft
				} else if form.IsNewer(def.Version, highest[def.AddonName][def.HTMLName][major].Version) {
					highest[def.AddonName][def.HTMLName][major] = &ft
				}
			}
		}
	}
	// Mark the highest minor number of each major number of each form to
	// accept newer versions.
	for _, addons := range highest {
		for _, htmls := range addons {
			for _, major := range htmls {
				major.AcceptNewer = true
			}
		}
	}
	return err
}
