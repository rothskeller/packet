// usage: sign-forms bundle-name zip-file
//
// sign-forms computes a digital signature for the specified zip file and
// creates the corresponding ".forms" file with the same basename as the
// supplied ZIP file.
package main

import (
	"crypto"
	"crypto/sha512"
	"fmt"
	"hash"
	"io"
	"os"
	"strings"

	"github.com/rothskeller/packet/form/formdefs"
)

func main() {
	var (
		h     hash.Hash
		outfn string
		in    *os.File
		out   *os.File
		sig   []byte
		err   error
	)
	if len(os.Args) != 3 ||
		!formdefs.ValidBundleNameRE.MatchString(os.Args[1]) ||
		!strings.HasSuffix(strings.ToLower(os.Args[2]), ".zip") {
		fmt.Fprintln(os.Stderr, "usage: sign-forms bundle-name zip-file")
		os.Exit(2)
	}
	outfn = os.Args[1][:len(os.Args[2])-4] + ".forms"
	h = sha512.New()
	io.WriteString(h, os.Args[1])
	if in, err = os.Open(os.Args[2]); err != nil {
		goto ERROR
	}
	defer in.Close()
	if _, err = io.Copy(h, in); err != nil {
		goto ERROR
	}
	if _, err = in.Seek(0, 0); err != nil {
		goto ERROR
	}
	if sig, err = formsBundlePrivateKey.Sign(nil, h.Sum(nil), crypto.SHA512); err != nil {
		goto ERROR
	}
	if out, err = os.Create(outfn); err != nil {
		goto ERROR
	}
	if _, err = fmt.Fprintf(out, "PackItFormBundle%-15s\n", os.Args[1]); err != nil {
		goto ERROR
	}
	if _, err = out.Write(sig); err != nil {
		goto ERROR
	}
	if _, err = io.Copy(out, in); err != nil {
		goto ERROR
	}
	if err = out.Close(); err != nil {
		goto ERROR
	}
	return

ERROR:
	fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
	os.Remove(outfn)
	os.Exit(1)
}
