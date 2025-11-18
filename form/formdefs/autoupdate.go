package formdefs

import (
	"archive/zip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rothskeller/packet/errors"
)

const (
	updateFrequency = 24 * time.Hour
	requestTimeout  = 5 * time.Second
)

// updateInfo is the structure stored in update.json in the root of a bundle.
type updateInfo struct {
	URL             string    `json:"url,omitempty"`
	IfNoneMatch     string    `json:"ifNoneMatch,omitempty"`
	IfModifiedSince string    `json:"ifModifiedSince,omitempty"`
	LastCheck       time.Time `json:"lastCheck,omitempty"`
}

// CheckForUpdates fetches new versions of the form bundles from the Internet
// if appropriate and possible.  If force is true, checks are performed even if
// they've been done recently.  If readme is true, any README.html file in a
// new bundle is propagated to the root of the forms directory, which will
// cause it to be displayed next time there's an appropriate server request.
func CheckForUpdates(force, readme bool) (err error) {
	var (
		formsDir string
		ents     []os.DirEntry
	)
	if formsDir = FormsDir(); formsDir == "" {
		if force {
			return errors.New("can't check for forms update: no local forms directory")
		}
		return nil
	}
	if ents, err = os.ReadDir(formsDir); err != nil {
		return fmt.Errorf("forms update: %s", err)
	}
	for _, ent := range ents {
		if !ent.IsDir() || strings.HasSuffix(ent.Name(), ".new") {
			continue
		}
		if bundleErr := checkForBundleUpdate(formsDir, ent.Name(), force, readme); bundleErr != nil {
			bundleErr = fmt.Errorf("can't update bundle %q: %s", ent.Name(), bundleErr)
			err = errors.Join(err, bundleErr)
		}
	}
	return err
}

// checkForBundleUpdate conditionally fetches an update of the specified form
// bundle.  If force is true, the check is unconditional, otherwise it is
// checked only if it hasn't recently been checked.  If readme is true, and a
// new bundle is installed, any README.html file in the new bundle is merged
// into the README.html file in the root of the forms directory.
func checkForBundleUpdate(formsDir, bundle string, force, readme bool) (err error) {
	var (
		bundleDir string
		uiFile    string
		ui        *updateInfo
		zipFName  string
		zipFH     *os.File
		zipSize   int64
		bundleNew string
	)
	// Get the update info for the bundle and determine whether we should
	// do an update.
	bundleDir = filepath.Join(formsDir, bundle)
	uiFile = filepath.Join(bundleDir, "update.json")
	if ui, err = readUpdateInfo(uiFile); err != nil {
		return err
	} else if ui == nil || ui.URL == "" {
		slog.Debug("not updating bundle", "b", bundle, "r", "no update URL")
		return nil
	} else if !force && time.Since(ui.LastCheck) < updateFrequency {
		slog.Debug("not updating bundle", "b", bundle, "r", "checked recently")
		return nil
	}
	// Fetch the update, if there is one.
	zipFName = bundleDir + ".zip"
	if zipFH, zipSize, err = fetchUpdate(zipFName, ui); err != nil {
		return err
	} else if zipFH == nil {
		slog.Debug("not updating bundle", "b", bundle, "r", "no update available")
		goto CHECKED
	}
	defer os.Remove(zipFName)
	defer zipFH.Close()
	// Unpack the zip file.
	bundleNew = bundleDir + ".new"
	if err = unpackZip(zipFH, zipSize, bundleNew); err != nil {
		return err
	}
	// Remove the old bundle and move the new one into place.
	// move the new one into place.
	if err = os.RemoveAll(bundleDir); err != nil {
		slog.Error("os.RemoveAll", "d", bundleDir, "err", err)
		return fmt.Errorf("remove old bundle: %s", err)
	}
	if err = os.Rename(bundleNew, bundleDir); err != nil {
		slog.Error("os.Rename", "from", bundleNew, "to", bundleDir, "err", err)
		return fmt.Errorf("move new bundle into place: %s", err)
	}
	// The new bundle may have a new update URL.  Check for that.
	if ui2, err := readUpdateInfo(uiFile); err == nil && ui2 != nil {
		ui.URL = ui2.URL
	}
	// The new bundle may have a README.html.  Check for that.
	if readme {
		if err = propagateReadme(bundleDir, formsDir); err != nil {
			return err
		}
	}
	slog.Info("updated forms bundle", "b", bundle)
CHECKED:
	ui.LastCheck = time.Now()
	return writeUpdateInfo(ui, uiFile)
}

// propagateReadme takes the contents of the README.html in the bundle dir, if
// any, and appends them to the README.html in the formsDir, creating it if
// needed.
func propagateReadme(bundleDir, formsDir string) (err error) {
	var (
		filename string
		data     []byte
		fh       *os.File
	)
	// Read the README from the bundle dir.
	filename = filepath.Join(bundleDir, "README.html")
	if data, err = os.ReadFile(filename); os.IsNotExist(err) {
		return nil // no README file
	} else if err != nil {
		slog.Error("os.ReadFile", "f", filename, "err", err)
		return err
	}
	// Append to the README file in the forms dir.
	filename = filepath.Join(formsDir, "README.html")
	if fh, err = os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666); err != nil {
		slog.Error("os.OpenFile", "f", filename, "err", err)
		return err
	}
	if _, err = fh.Write(data); err != nil {
		fh.Close()
		slog.Error("fh.Write", "f", filename, "err", err)
		return fmt.Errorf("write %s: %s", filename, err)
	}
	if err = fh.Close(); err != nil {
		slog.Error("fh.Close", "f", filename, "err", err)
		return fmt.Errorf("close %s: %s", filename, err)
	}
	return nil
}

// readUpdateInfo reads the update.json file at the specified filename.
func readUpdateInfo(filename string) (ui *updateInfo, err error) {
	if data, err := os.ReadFile(filename); os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		ui = new(updateInfo)
		if err = json.Unmarshal(data, ui); err != nil {
			return nil, fmt.Errorf("json decode %s: %s", filename, err)
		}
		return ui, nil
	}
}

// writeUpdateInfo writes the update.json file.
func writeUpdateInfo(ui *updateInfo, filename string) (err error) {
	data, _ := json.Marshal(ui)
	if err = os.WriteFile(filename, data, 0666); err != nil {
		slog.Error("os.WriteFile", "f", filename, "err", err)
		return fmt.Errorf("can't write update info: %s", err)
	}
	return nil
}

// fetchUpdate retrieves the update zip file from the update server into the
// specified file.  If successful, it returns the handle to the open zip file
// and its size, and sets the IfNoneMatch and IfModifiedSince fields of the
// updateInfo to match the ETag and Last-Modified headers of the server
// response.  If there is no update available, it returns (nil, nil).  An error
// is returned if the server returns an error or the file cannot be saved.
func fetchUpdate(zipFName string, ui *updateInfo) (fh *os.File, size int64, err error) {
	var (
		ctx    context.Context
		cancel func()
		req    *http.Request
		resp   *http.Response
	)
	// Build the request.
	ctx, cancel = context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	req, _ = http.NewRequestWithContext(ctx, http.MethodGet, ui.URL, nil)
	if ui.IfModifiedSince != "" {
		req.Header.Set("If-Modified-Since", ui.IfModifiedSince)
	}
	if ui.IfNoneMatch != "" {
		req.Header.Set("If-None-Match", ui.IfNoneMatch)
	}
	// Issue the request and check the result.
	if resp, err = http.DefaultClient.Do(req); err != nil {
		slog.Error("http.Get", "u", ui.URL, "err", err)
		return nil, 0, fmt.Errorf("http GET %s: %s", ui.URL, err)
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case http.StatusNoContent, http.StatusNotModified, http.StatusNotFound:
		// These are "benign" responses meaning there is no update.
		return nil, 0, nil
	case http.StatusOK:
		// we have an update, continue below
	default:
		slog.Error("http.Get", "u", ui.URL, "code", resp.StatusCode, "status", resp.Status)
		return nil, 0, fmt.Errorf("http GET %s: %d %s", ui.URL, resp.StatusCode, resp.Status)
	}
	// Store the bundle into the specified file.
	if fh, err = os.OpenFile(zipFName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666); err != nil {
		slog.Error("os.Create", "f", zipFName, "err", err)
		return nil, 0, fmt.Errorf("create zip: %s", err)
	}
	if size, err = io.Copy(fh, resp.Body); err != nil {
		slog.Error("io.Copy", "f", zipFName, "err", err)
		return nil, 0, fmt.Errorf("copy to %s: %s", zipFName, err)
	}
	if _, err = fh.Seek(0, io.SeekStart); err != nil {
		slog.Error("fh.Seek", "f", zipFName, "err", err)
		return nil, 0, fmt.Errorf("seek in %s: %s", zipFName, err)
	}
	// Save the response information.
	ui.IfModifiedSince = resp.Header.Get("Last-Modified")
	ui.IfNoneMatch = resp.Header.Get("ETag")
	return fh, size, nil
}

// unpackZip unpacks the zip file opened as zipFH, which has size zipSize,
// into the directory bundleNew.
func unpackZip(fh *os.File, size int64, bundleNew string) (err error) {
	var (
		z *zip.Reader
	)
	// Remove the target directory if it is left over from a previous op.
	// Also remove it if this function returns with an error
	if err = os.RemoveAll(bundleNew); err != nil {
		slog.Error("os.RemoveAll", "d", bundleNew, "err", err)
		return err
	}
	defer func() {
		if err != nil {
			os.RemoveAll(bundleNew)
		}
	}()
	// Open the zip header.
	if z, err = zip.NewReader(fh, size); err != nil && err != zip.ErrInsecurePath {
		slog.Error("zip.NewReader", "err", err)
		return fmt.Errorf("open zip: %s", err)
	}
	// Unpack each file.
	for _, file := range z.File {
		if strings.HasSuffix(file.Name, "/") {
			continue // directory entry
		}
		fname := filepath.Join(bundleNew, strings.ReplaceAll(file.Name, "/", string(filepath.Separator)))
		dname := filepath.Dir(fname)
		if err = os.MkdirAll(dname, 0777); err != nil {
			slog.Error("os.MkdirAll", "d", dname, "err", err)
			return fmt.Errorf("mkdir %s: %s", dname, err)
		}
		if out, err := os.Create(fname); err != nil {
			slog.Error("os.Create", "f", fname, "err", err)
			return fmt.Errorf("create %s: %s", fname, err)
		} else if in, err := file.Open(); err != nil {
			out.Close()
			slog.Error("zip.file.Open", "f", file.Name, "err", err)
			return fmt.Errorf("open in zip %s: %s", file.Name, err)
		} else if _, err = io.Copy(out, in); err != nil {
			out.Close()
			slog.Error("io.Copy", "f", fname, "err", err)
			return fmt.Errorf("copy to %s: %s", fname, err)
		} else if err = out.Close(); err != nil {
			slog.Error("os.Close", "f", fname, "err", err)
			return fmt.Errorf("close %s: %s", fname, err)
		}
	}
	return nil
}
