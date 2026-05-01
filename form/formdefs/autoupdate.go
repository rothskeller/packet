package formdefs

import (
	"archive/zip"
	"context"
	"crypto"
	"crypto/ed25519"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/rothskeller/packet/errors"
)

const (
	updateFrequency = 24 * time.Hour
	requestTimeout  = 5 * time.Second
)

// ValidBundleNameRE matches a valid bundle name.
var ValidBundleNameRE = regexp.MustCompile(`^[A-Z][-A-Za-z0-9_]{0,14}$`)

// updateInfo is the structure stored in update.json in the root of a bundle.
type updateInfo struct {
	URL             string    `json:"url,omitempty"`
	IfNoneMatch     string    `json:"ifNoneMatch,omitempty"`
	IfModifiedSince string    `json:"ifModifiedSince,omitempty"`
	LastCheck       time.Time `json:"lastCheck,omitempty"`
}

// CheckForUpdates fetches new versions of the form bundles from the Internet if
// appropriate and possible.  If force is true, checks are performed even if
// they've been done recently.  If readme is true, any README.txt file in a new
// bundle is propagated to the root of the forms directory, which will cause it
// to be displayed next time there's an appropriate server request.
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
// new bundle is installed, any README.txt file in the new bundle is merged into
// the README.txt file in the root of the forms directory.
func checkForBundleUpdate(formsDir, bundle string, force, readme bool) (err error) {
	var (
		bundleDir   string
		uiFile      string
		ui          *updateInfo
		bundleFName string
		bundleFH    *os.File
		bundleSize  int64
		readmeText  string
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
	bundleFName = bundleDir + ".forms"
	if bundleFH, bundleSize, err = fetchUpdate(bundleFName, ui); err != nil {
		return err
	} else if bundleFH == nil {
		slog.Debug("not updating bundle", "b", bundle, "r", "no update available")
		goto CHECKED
	}
	defer os.Remove(bundleFName)
	defer bundleFH.Close()
	// Install the bundle.
	if _, readmeText, err = installBundle(bundleFName, bundle, bundleFH, bundleSize, formsDir); err != nil {
		return err
	}
	// Save the README text if we're supposed to.
	if readme {
		if err = AppendReadme(readmeText); err != nil {
			return err
		}
	}
	// Save the update information.
	if ui2, err := readUpdateInfo(uiFile); err != nil {
		return err
	} else {
		slog.Info("updated forms bundle", "b", bundle)
		if ui2 == nil { // new bundle doesn't have auto-update
			return nil
		} else {
			ui.URL = ui2.URL // keep the URL from the new bundle
		}
	}
CHECKED:
	ui.LastCheck = time.Now()
	return writeUpdateInfo(ui, uiFile)
}

// InstallBundle installs a bundle from a specified source, which can be either
// a filename or an https:// URL.  If successful, it returns the bundle name and
// any README text from the bundle.  (The caller can either display it or pass
// it to AppendReadme for later display.)
func InstallBundle(source string) (bundle, readme string, err error) {
	var (
		formsDir string
		ui       updateInfo
		fh       *os.File
		size     int64
	)
	if formsDir = FormsDir(); formsDir == "" {
		slog.Error("no local forms dir")
		return "", "", errors.New("Forms cannot be installed on this system because there is no local forms directory.")
	}
	if err = os.MkdirAll(formsDir, 0777); err != nil {
		slog.Error("os.MkdirAll", "d", formsDir, "err", err)
		return "", "", errors.NewF("The forms directory %s could not be created.", formsDir)
	}
	if strings.HasPrefix(source, "https://") {
		var tempfilename = filepath.Join(formsDir, "install-temp.forms")
		defer os.Remove(tempfilename)
		ui.URL = source
		if fh, size, err = fetchUpdate(tempfilename, &ui); err != nil {
			return "", "", err
		}
		source = tempfilename
	} else if fh, err = os.Open(source); err != nil {
		slog.Error("os.Open", "f", source, "err", err)
		return "", "", errors.NewF("The forms bundle file %s could not be opened.", source)
	} else if stat, err := fh.Stat(); err != nil {
		slog.Error("fh.Stat", "f", source, "err", err)
		return "", "", errors.NewF("The forms bundle file %s could not be read.", source)
	} else {
		size = stat.Size()
	}
	defer fh.Close()
	// Install the bundle.
	if bundle, readme, err = installBundle(source, "", fh, size, formsDir); err != nil {
		return "", "", err
	}
	// If we read from a URL, save the update information.
	if ui.URL != "" {
		uiFile := filepath.Join(formsDir, bundle, "update.json")
		if ui2, err := readUpdateInfo(uiFile); err != nil {
			return "", "", err
		} else if ui2 != nil && ui.URL == ui2.URL {
			// Only store the update info if the bundle has auto
			// update from the same source where we got it.
			ui.LastCheck = time.Now()
			if err = writeUpdateInfo(&ui, uiFile); err != nil {
				return "", "", err
			}
		}
	}
	slog.Info("Installed bundle", "b", bundle, "source", source)
	return bundle, readme, nil
}

// fetchUpdate retrieves the update bundle from the update server into the
// specified file.  If successful, it returns the handle to the open bundle and
// its size, and sets the IfNoneMatch and IfModifiedSince fields of the
// updateInfo to match the ETag and Last-Modified headers of the server
// response.  If there is no update available, it returns (nil, nil).  An error
// is returned if the server returns an error or the file cannot be saved.
func fetchUpdate(bundleFName string, ui *updateInfo) (fh *os.File, size int64, err error) {
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
	if fh, err = os.OpenFile(bundleFName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666); err != nil {
		slog.Error("os.Create", "f", bundleFName, "err", err)
		return nil, 0, fmt.Errorf("create zip: %s", err)
	}
	if size, err = io.Copy(fh, resp.Body); err != nil {
		slog.Error("io.Copy", "f", bundleFName, "err", err)
		return nil, 0, fmt.Errorf("copy to %s: %s", bundleFName, err)
	}
	if _, err = fh.Seek(0, io.SeekStart); err != nil {
		slog.Error("fh.Seek", "f", bundleFName, "err", err)
		return nil, 0, fmt.Errorf("seek in %s: %s", bundleFName, err)
	}
	// Save the response information.
	ui.IfModifiedSince = resp.Header.Get("Last-Modified")
	ui.IfNoneMatch = resp.Header.Get("ETag")
	return fh, size, nil
}

// installBundle installs a bundle file with the specified source (for error
// messages), bundle name (if known, otherwise read from the file), open file
// handle (rewound to the beginning), and size, into formsDir.  It returns the
// bundle name and the contents of the README.txt file in the bundle, if any.
func installBundle(source, expect string, fh *os.File, size int64, formsDir string) (bundle, readme string, err error) {
	var (
		bundleDir  string
		bundleNew  string
		launchFile string
		rmfile     string
	)
	if bundle, err = unpackBundle(source, expect, fh, size, formsDir); err != nil {
		return "", "", err
	}
	bundleDir = filepath.Join(formsDir, bundle)
	bundleNew = bundleDir + ".new"
	// There should be a $bundle.launch file in that directory.  Check for
	// that.
	launchFile = filepath.Join(bundleNew, bundle+".launch")
	if _, err = os.Stat(launchFile); err != nil {
		return "", "", errors.NewF("The new forms bundle does not contain a %s.launch file.", bundle)
	}
	// Generate the associated $bundle.ini file.
	if err = writeAddonINI(launchFile); err != nil {
		return "", "", errors.NewF("%s.ini could not be added to the forms bundle: %s", bundle, err)
	}
	// Remove the old bundle and move the new one into place.
	// move the new one into place.
	if err = os.RemoveAll(bundleDir); err != nil {
		slog.Error("os.RemoveAll", "d", bundleDir, "err", err)
		return "", "", errors.NewF("The old bundle directory %s could not be removed.", bundleDir)
	}
	if err = os.Rename(bundleNew, bundleDir); err != nil {
		slog.Error("os.Rename", "from", bundleNew, "to", bundleDir, "err", err)
		return "", "", errors.NewF("The new bundle directory could not be moved to %s.", bundleDir)
	}
	// The new bundle may have a README.txt.  Check for that.
	rmfile = filepath.Join(bundleDir, "README.txt")
	if rm, err := os.ReadFile(rmfile); err == nil || os.IsNotExist(err) {
		readme = string(rm)
	} else {
		slog.Error("os.ReadFile", "f", rmfile, "err", err)
		return "", "", errors.New("The README.txt file in the new bundle could not be read.")
	}
	return bundle, readme, nil
}

// unpackBundle unpacks the forms bundle opened from source as fh, which has the
// specified size, into the formsDir/bundleName.new directory, where bundleName
// is the bundle name read from the file, and returns the bundle name.  If an
// expect string is supplied, the bundle name must match it.
func unpackBundle(source, expect string, fh *os.File, size int64, formsDir string) (bundleName string, err error) {
	var (
		dir string
		zr  io.ReaderAt
		z   *zip.Reader
	)
	// Verify the digital signature of the bundle.
	if zr, bundleName, err = verifySignature(source, fh, size, expect); err != nil {
		return "", err
	}
	dir = filepath.Join(formsDir, bundleName+".new")
	// Remove the target directory if it is left over from a previous op.
	// Also remove it if this function returns with an error
	if err = os.RemoveAll(dir); err != nil {
		slog.Error("os.RemoveAll", "d", dir, "err", err)
		return "", errors.NewF("The previously existing directory %s could not be removed.", dir)
	}
	defer func() {
		if err != nil {
			os.RemoveAll(dir)
		}
	}()
	// Open the zip header.
	if z, err = zip.NewReader(zr, size); err != nil && err != zip.ErrInsecurePath {
		slog.Error("zip.NewReader", "err", err)
		return "", fmt.Errorf("open zip: %s", err)
	}
	// Unpack each file.
	for _, file := range z.File {
		if strings.HasSuffix(file.Name, "/") {
			continue // directory entry
		}
		fname := filepath.Join(dir, strings.ReplaceAll(file.Name, "/", string(filepath.Separator)))
		dname := filepath.Dir(fname)
		if err = os.MkdirAll(dname, 0777); err != nil {
			slog.Error("os.MkdirAll", "d", dname, "err", err)
			return "", errors.NewF("The directory %s could not be created.", dname)
		}
		if out, err := os.Create(fname); err != nil {
			slog.Error("os.Create", "f", fname, "err", err)
			return "", errors.NewF("The file %s could not be created.", fname)
		} else if in, err := file.Open(); err != nil {
			out.Close()
			slog.Error("zip.file.Open", "f", file.Name, "err", err)
			return "", errors.NewF("The bundle file %s could not be decoded.", source)
		} else if _, err = io.Copy(out, in); err != nil {
			out.Close()
			slog.Error("io.Copy", "f", fname, "err", err)
			return "", errors.NewF("The file %s could not be written.", fname)
		} else if err = out.Close(); err != nil {
			slog.Error("os.Close", "f", fname, "err", err)
			return "", errors.NewF("The file %s could not be written.", fname)
		}
	}
	return bundleName, nil
}

// verifySignature verifies that the forms bundle was digitally signed by the
// key for this version of the packet software.  It assumes the bundle file is
// opened and rewound.  If an expect string is given, it verifies that it
// matches the bundle name in the file.  If the signature is verified, it
// returns an io.ReaderAt that addresses the ZIP part of the file and the actual
// bundle name in the file.
func verifySignature(source string, fh *os.File, size int64, expect string) (zr io.ReaderAt, bundle string, err error) {
	var (
		h   hash.Hash
		sig = make([]byte, ed25519.SignatureSize)
	)
	if _, err = fh.Read(sig[:32]); err != nil {
		slog.Error("fh.Read 1", "src", source, "err", err)
		return nil, "", errors.NewF("The bundle file %s could not be read.", source)
	}
	if string(sig[:16]) != "PackItFormBundle" {
		slog.Error("not a form bundle header", "src", source)
		return nil, "", errors.NewF("The file %s is not a forms bundle file.", source)
	}
	if bundle = strings.TrimRight(string(sig[16:32]), " \n"); !ValidBundleNameRE.MatchString(bundle) {
		slog.Error("contains invalid bundle name", "src", source)
		return nil, "", errors.NewF("The file %s contains an invalid forms bundle name %q.", source, bundle)
	} else if expect != "" && bundle != expect {
		slog.Error("contains wrong bundle name", "src", source, "exp", expect, "act", bundle)
		return nil, "", errors.NewF("The file %s contains forms bundle %q, not %q.", source, bundle, expect)
	}
	if _, err = fh.Read(sig); err != nil {
		slog.Error("fh.Read 2", "src", source, "err", err)
		return nil, "", errors.NewF("The bundle file %s could not be read.", source)
	}
	h = sha512.New()
	io.WriteString(h, bundle)
	if _, err = io.Copy(h, fh); err != nil {
		slog.Error("io.Copy", "src", source, "err", err)
		return nil, "", errors.NewF("The bundle file %s could not be read.", source)
	}
	if err = ed25519.VerifyWithOptions(formsBundlePublicKey, h.Sum(nil), sig, &ed25519.Options{Hash: crypto.SHA512}); err != nil {
		slog.Error("ed25519.VerifyWithOptions", "src", source, "err", err)
		return nil, "", errors.NewF("The signature of the bundle file %s is not correct.  The bundle file has been corrupted, was improperly signed, or is intended for a different version of SCCo Packet software.", source)
	}
	zr = io.NewSectionReader(fh, ed25519.SignatureSize+32, size-ed25519.SignatureSize-32)
	return zr, bundle, nil
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

// AppendReadme appends the supplied text to the README.txt in the forms
// directory, creating it if needed.
func AppendReadme(text string) (err error) {
	var (
		filename string
		fh       *os.File
	)
	if text == "" {
		return nil
	}
	filename = filepath.Join(FormsDir(), "README.txt")
	if fh, err = os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666); err != nil {
		slog.Error("os.OpenFile", "f", filename, "err", err)
		return errors.NewF("The file %s could not be opened or created.", filename)
	}
	if _, err = io.WriteString(fh, text); err != nil {
		fh.Close()
		slog.Error("io.WriteString", "f", filename, "err", err)
		return errors.NewF("The file %s could not be written.", filename)
	}
	if err = fh.Close(); err != nil {
		slog.Error("fh.Close", "f", filename, "err", err)
		return errors.NewF("The file %s could not be written.", filename)
	}
	return nil
}
