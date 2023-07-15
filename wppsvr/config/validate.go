package config

import (
	"log"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"strings"

	"github.com/rothskeller/packet/message"
)

var fccCallRE = regexp.MustCompile(`^(?:A[A-L]|[KNW][A-Z]?)[0-9][A-Z]{1,3}$`)
var ax25RE = regexp.MustCompile(`^((?:A[A-L]|[KNW][A-Z]?)[0-9][A-Z]{1,3})-(?:[0-9]|1[0-5])$`)
var tacticalCallRE = regexp.MustCompile(`^[A-Z][A-Z0-9]{0,5}$`)

// Validate checks the configuration to make sure all fields have valid values.
// If there are any errors, they are logged, and the function returns false.
// If the configuration is valid, the function returns true.
func (c *Config) Validate() (valid bool) {
	var err error
	valid = true // assume valid until proven otherwise
	var haveHTMLReports bool

	// Check each of the BBS configurations.
	for bbsCall, bbsConf := range c.BBSes {
		if !fccCallRE.MatchString(bbsCall) {
			log.Printf("ERROR: config.bbses: %q is not a valid FCC call sign", bbsCall)
			valid = false
		}
		switch bbsConf.Transport {
		case "":
			log.Printf("ERROR: config.bbses[%q].transport is not specified", bbsCall)
			valid = false
		case "disable":
			break
		case "kpc3plus":
			if bbsConf.AX25 == "" {
				log.Printf("ERROR: config.bbses[%q].ax25 is not specified", bbsCall)
				valid = false
			} else if match := ax25RE.FindStringSubmatch(bbsConf.AX25); match == nil || match[1] != bbsCall {
				log.Printf("ERROR: config.bbses[%q].ax25 = %q is not a valid AX25 address for %s",
					bbsCall, bbsConf.AX25, bbsCall)
				valid = false
			}
		case "telnet":
			if bbsConf.TCP == "" {
				log.Printf("ERROR: config.bbses[%q].tcp is not specified", bbsCall)
				valid = false
			} else if _, _, err = net.SplitHostPort(bbsConf.TCP); err != nil {
				log.Printf("ERROR: config.bbses[%q].tcp = %q is not a host:port string", bbsCall, bbsConf.TCP)
				valid = false
			}
			for callsign := range bbsConf.Passwords {
				if !tacticalCallRE.MatchString(callsign) {
					log.Printf("ERROR: config.bbses[%q].passwords: %q is not a valid tactical call sign",
						bbsCall, callsign)
					valid = false
				}
			}
		default:
			log.Printf("ERROR: config.bbses[%q].transport = %q is not a known transport", bbsCall, bbsConf.Transport)
			valid = false
		}
	}

	// Check that we have a minimum PackItForms version.
	if c.MinPIFOVersion == "" {
		log.Printf("ERROR: config.minPIFOVersion is not specified")
		valid = false
	}

	// Check that we have configuration for every form.
	if c.MessageTypes == nil {
		log.Printf("ERROR: config.messageTypes is not specified")
		valid = false
	} else {
		for tag, mtc := range c.MessageTypes {
			if _, ok := message.RegisteredTypes[tag]; !ok {
				log.Printf("ERROR: config.messageTypes has entry for unknown message type %q", tag)
				valid = false
				continue
			}
			if tag == "plain" {
				continue
			}
			if mtc.MinimumVersion == "" {
				log.Printf("ERROR: config.messageTypes[%q].minimumVersion is not specified", tag)
				valid = false
			}
			switch mtc.HandlingOrder {
			case "", "IMMEDIATE", "PRIORITY", "ROUTINE":
				break
			case "computed":
				if tag != "ICS213" && tag != "EOC213RR" {
					log.Printf("ERROR: config.messageTypes[%q].handlingOrder = %q, but that form has no handling order computation defined", tag, mtc.HandlingOrder)
					valid = false
				}
			default:
				log.Printf("ERROR: config.messageTypes[%q].handlingOrder = %q is not a valid handling order", tag, mtc.HandlingOrder)
				valid = false
			}
		}
		for tag := range message.RegisteredTypes {
			if _, ok := c.MessageTypes[tag]; !ok {
				if tag == "plain" {
					c.MessageTypes["plain"] = new(MessageTypeConfig)
				} else if tag != "UNKNOWN" {
					log.Printf("ERROR: config.messageTypes[%q] is not specified", tag)
					valid = false
				}
			}
		}
	}

	// Check that we have a URL for the web server.
	if c.ServerURL == "" {
		log.Printf("ERROR: config.serverURL is not specified")
		valid = false
	} else if u, err := url.Parse(c.ServerURL); err != nil || u.Scheme != "https" {
		log.Printf("ERROR: config.serverURL has an invalid value, not an https:// URL")
		valid = false
	}

	// Check that we have a listen address for the web server.
	if c.ListenAddr == "" {
		log.Printf("ERROR: config.listenAddr is not specified")
		valid = false
	}

	// Check the SMTP configuration.
	if c.SMTP != nil {
		if c.SMTP.From == "" {
			log.Printf("ERROR: config.smtp.from is not specified")
			valid = false
		} else if _, err := mail.ParseAddress(c.SMTP.From); err != nil {
			log.Printf("ERROR: config.smtp.from is not a valid email address")
			valid = false
		}
		if c.SMTP.Server == "" {
			log.Printf("ERROR: config.smtp.server is not specified")
			valid = false
		} else if _, _, err := net.SplitHostPort(c.SMTP.Server); err != nil {
			log.Printf("ERROR: config.smtp.server is not a valid host:port")
			valid = false
		}
		if c.SMTP.Username == "" {
			log.Printf("ERROR: config.smtp.username is not specified")
			valid = false
		}
		if c.SMTP.Password == "" {
			log.Printf("ERROR: config.smtp.password is not specified")
			valid = false
		}
	} else if haveHTMLReports {
		log.Printf("ERROR: config.smtp is needed and not specified")
		valid = false
	}

	// Check that the permissions are granted to real call signs.
	for i, call := range c.CanViewEveryone {
		call = strings.ToUpper(call)
		if !fccCallRE.MatchString(call) {
			log.Printf("ERROR: config.canViewEveryone[%d] = %q: not a valid call sign", i, call)
			valid = false
		}
		c.CanViewEveryone[i] = call
	}
	for i, call := range c.CanEditSessions {
		call = strings.ToUpper(call)
		if !fccCallRE.MatchString(call) {
			log.Printf("ERROR: config.canEditSessions[%d] = %q: not a valid call sign", i, call)
			valid = false
		}
		c.CanEditSessions[i] = call
	}
	return valid
}
