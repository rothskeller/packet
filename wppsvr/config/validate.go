package config

import (
	"fmt"
	"log"
	"net"
	"regexp"
)

var fccCallRE = regexp.MustCompile(`^[AKNW][A-Z]?[0-9][A-Z]{1,3}$`)
var ax25RE = regexp.MustCompile(`^([AKNW][A-Z]?[0-9][A-Z]{1,3})-(?:[0-9]|1[0-5])$`)
var tacticalCallRE = regexp.MustCompile(`^[A-Z][A-Z0-9]{0,5}$`)
var prefixRE = regexp.MustCompile(`^(?:[A-Z]{3}|[A-Z][0-9][A-Z]|[0-9][A-Z]{2})$`)

// Validate checks the configuration to make sure all fields have valid values.
// If there are any errors, they are logged, and the function returns false.
// If the configuration is valid, the function returns true.
func (c *Config) Validate() (valid bool) {
	var err error

	valid = true // assume valid until proven otherwise

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

	// Check each of the session configurations.
	for toCallSign, session := range c.Sessions {
		if !tacticalCallRE.MatchString(toCallSign) {
			log.Printf("ERROR: config.sessions: %q is not a valid tactical call sign", toCallSign)
			valid = false
		}
		if session.Name == "" {
			log.Printf("ERROR: config.sessions[%q].name is not specified", toCallSign)
			valid = false
		}
		if session.Prefix == "" {
			log.Printf("ERROR: config.sessions[%q].prefix is not specified", toCallSign)
			valid = false
		} else if !prefixRE.MatchString(session.Prefix) {
			log.Printf("ERROR: config.sessions[%q].prefix = %q is not a valid prefix", toCallSign, session.Prefix)
			valid = false
		}
		if session.Start == "" {
			log.Printf("ERROR: config.sessions[%q].start is not specified", toCallSign)
			valid = false
		} else if session.StartInterval, err = ParseInterval(session.Start); err != nil {
			log.Printf("ERROR: config.sessions[%q].start = %q: %s", toCallSign, session.Start, err)
			valid = false
		}
		if session.End == "" {
			log.Printf("ERROR: config.sessions[%q].end is not specified", toCallSign)
			valid = false
		} else if session.EndInterval, err = ParseInterval(session.End); err != nil {
			log.Printf("ERROR: config.sessions[%q].end = %q: %s", toCallSign, session.End, err)
			valid = false
		}
		if !session.ToBBSes.validate(fmt.Sprintf("config.sessions[%q].toBBS", toCallSign)) {
			valid = false
		}
		for i, item := range session.ToBBSes {
			if _, ok := c.BBSes[item.Then]; !ok {
				log.Printf("ERROR: config.sessions[%q].toBBS[%d] = %q is not a recognized BBS",
					toCallSign, i, item.Then)
				valid = false
			}
		}
		if !session.DownBBSes.validate(fmt.Sprintf("config.sessions[%q].downBBSes", toCallSign)) {
			valid = false
		}
		for i, item := range session.DownBBSes {
			if _, ok := c.BBSes[item.Then]; !ok {
				log.Printf("ERROR: config.sessions[%q].downBBSes[%d] = %q is not a recognized BBS",
					toCallSign, i, item.Then)
				valid = false
			}
		}
		for _, bbs := range session.RetrieveFromBBSes {
			if _, ok := c.BBSes[bbs]; !ok {
				log.Printf("ERROR: config.sessions[%q].retrieveFromBBSes: %q is not a recognized BBS",
					toCallSign, bbs)
				valid = false
			}
		}
		if session.RetrieveAt == "" {
			log.Printf("ERROR: config.sessions[%q].retrieveAt is not specified", toCallSign)
			valid = false
		} else if session.RetrieveAtInterval, err = ParseInterval(session.RetrieveAt); err != nil {
			log.Printf("ERROR: config.sessions[%q].retrieveAt = %q: %s", toCallSign, session.RetrieveAt, err)
			valid = false
		}
		if !session.MessageTypes.validate(fmt.Sprintf("config.sessions[%q].messageTypes", toCallSign)) {
			valid = false
		}
		for i, item := range session.MessageTypes {
			if LookupMessageType(item.Then) == nil {
				log.Printf("ERROR: config.sessions[%q].messageTypes[%d] = %q is not a recognized message type",
					toCallSign, i, item.Then)
				valid = false
			}
		}
	}

	// Check that we have minimum versions for every form.
	if c.MinimumVersions == nil {
		log.Printf("ERROR: config.minimumVersions is not specified")
		valid = false
	} else {
		if c.MinimumVersions[PackItForms] == "" {
			log.Printf("ERROR: config.minimumVersions[%q] is not specified", PackItForms)
			valid = false
		}
		for _, mtype := range ValidMessageTypes {
			if form := mtype.Form(); form != nil {
				if c.MinimumVersions[form.TypeCode()] == "" {
					log.Printf("ERROR: config.minimumVersions[%q] is not specified", form.TypeCode())
					valid = false
				}
			}
		}
	}
	return valid
}

// validate makes sure that the interval specified in a scheduled value is
// parseable.
func (sv ScheduledValue) validate(label string) (valid bool) {
	var err error

	valid = true
	for i, item := range sv {
		if item.When == "" {
			log.Printf("ERROR: %s[%d].when is not specified", label, i)
			valid = false
		} else if item.Interval, err = ParseInterval(item.When); err != nil {
			log.Printf("ERROR: %s[%d].when = %q: %s", label, i, item.When, err)
			valid = false
		}
	}
	return valid
}
