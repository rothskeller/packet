// usage: sign-forms zip-file
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
	if len(os.Args) != 2 || !strings.HasSuffix(strings.ToLower(os.Args[1]), ".zip") {
		fmt.Fprintln(os.Stderr, "usage: sign-forms zip-file")
		os.Exit(2)
	}
	outfn = os.Args[1][:len(os.Args[1])-4] + ".forms"
	h = sha512.New()
	if in, err = os.Open(os.Args[1]); err != nil {
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
