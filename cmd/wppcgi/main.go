// wppcgi is the CGI "script" for the weekly packet practice webserver.
package main

import (
	"log"
	"os"

	"github.com/rothskeller/packet/form/formdefs"
	"github.com/rothskeller/packet/wppsvr/config"
	"github.com/rothskeller/packet/wppsvr/store"
	"github.com/rothskeller/packet/wppsvr/webserver"
)

func main() {
	var (
		st  *store.Store
		err error
	)
	os.Setenv("TZ", "PST8PDT")
	// One of these chdirs should succeed; we don't care which one.  If
	// neither does, and there is no wppsvr.db in the current directory,
	// the store.Open will fail.
	os.Chdir("/Users/stever/src/packet/v2/wppsvr")
	os.Chdir("/home/sccares/private/wpp")
	openLog()
	formdefs.UseInternalForms = true
	formdefs.RegisterForms()
	if st, err = store.Open(); err != nil {
		log.Fatalf("ERROR: %s", err)
	}
	if err = config.Read(); err != nil {
		os.Exit(1)
	}
	webserver.HandleCGI(st)
}
