package webserver

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/rothskeller/packet/wppsvr/config"
)

// serveLogin responds to POST /login requests.
func (ws *webserver) serveLogin(w http.ResponseWriter, r *http.Request) {
	callsign := r.FormValue("callsign")
	password := r.FormValue("password")
	if callsign == "" || password == "" || !validLogin(callsign, password) {
		http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
		return
	}
	token := randomToken()
	callsign = strings.ToUpper(callsign)
	ws.st.AddLogin(token, callsign, time.Now().Add(time.Hour))
	http.SetCookie(w, &http.Cookie{Name: "auth", Value: token, Path: "/", Secure: true})
	w.WriteHeader(http.StatusNoContent)
}

// validLogin determines whether a callsign/password combination is valid.  It
// does so by attempting to log into https://scc-ares-races.org with it.
func validLogin(callsign, password string) bool {
	var client http.Client
	client.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	response, err := client.PostForm("https://www.scc-ares-races.org/activities/login01.php", url.Values{
		"user_id":  {callsign},
		"password": {password},
		"Submit":   {"Log In"},
	})
	if err != nil {
		log.Printf("ERROR: checking login of %q: post to scc-ares-races.org: %s", callsign, err)
		return false
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusFound {
		log.Printf("ERROR: checking login of %q: scc-ares-races.org did not redirect", callsign)
		return false
	}
	if response.Header.Get("Location") != "events.php" {
		log.Printf("LOGIN FAIL: %s", callsign)
		return false
	}
	log.Printf("LOGIN: %s", callsign)
	return true
}

// checkLoggedIn verifies that the user is logged in, and returns their call
// sign.  If the user is not properly logged in, it emits a redirect to the
// login page and returns an empty string.
func (ws *webserver) checkLoggedIn(w http.ResponseWriter, r *http.Request) (callsign string) {
	if c, err := r.Cookie("auth"); err == nil {
		callsign = ws.st.GetLogin(c.Value)
	}
	if callsign == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	return callsign
}

// randomToken returns a random token string.
func randomToken() string {
	var (
		tokenb [24]byte
		err    error
	)
	if _, err = rand.Read(tokenb[:]); err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(tokenb[:])
}

// canEditSessions returns whether the viewer (identified by callsign) is
// allowed to edit session definitions.
func canEditSessions(callsign string) bool {
	for _, cs := range config.Get().CanEditSessions {
		if cs == callsign {
			return true
		}
	}
	return false
}

// canViewEveryone returns whether the viewer (identified by callsign) is
// allowed to view other people's messages.
func canViewEveryone(callsign string) bool {
	for _, cs := range config.Get().CanViewEveryone {
		if cs == callsign {
			return true
		}
	}
	return false
}
