package config

import (
	"fmt"
	"log"
	"net"
	"regexp"
	"strings"

	"github.com/rothskeller/packet/wppsvr/interval"
	"github.com/rothskeller/packet/xscmsg"
)

var fccCallRE = regexp.MustCompile(`^[AKNW][A-Z]?[0-9][A-Z]{1,3}$`)
var ax25RE = regexp.MustCompile(`^([AKNW][A-Z]?[0-9][A-Z]{1,3})-(?:[0-9]|1[0-5])$`)
var tacticalCallRE = regexp.MustCompile(`^[A-Z][A-Z0-9]{0,5}$`)
var prefixRE = regexp.MustCompile(`^(?:[A-Z]{3}|[A-Z][0-9][A-Z]|[0-9][A-Z]{2})$`)
var variableRE = regexp.MustCompile(`\{[^}]*\}`)

// Validate checks the configuration to make sure all fields have valid values.
// If there are any errors, they are logged, and the function returns false.
// If the configuration is valid, the function returns true.
func (c *Config) Validate(knownProbs map[string]map[string]struct{}, knownVars map[string]struct{}) (valid bool) {
	var (
		err          error
		sawRetrieval = make(map[string]bool)
	)
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
		} else if session.StartInterval = interval.Parse(session.Start); session.StartInterval == nil {
			log.Printf("ERROR: config.sessions[%q].start = %q is not a valid interval", toCallSign, session.Start)
			valid = false
		}
		if session.End == "" {
			log.Printf("ERROR: config.sessions[%q].end is not specified", toCallSign)
			valid = false
		} else if session.EndInterval = interval.Parse(session.End); session.EndInterval == nil {
			log.Printf("ERROR: config.sessions[%q].end = %q is not a valid interval", toCallSign, session.End)
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
		for i, r := range session.Retrieve {
			key := r.Mailbox + "@" + r.BBS
			if sawRetrieval[key] {
				log.Printf("ERROR: config.sessions[%q].retrieve: multiple retrievals from %s", toCallSign, key)
				valid = false
			}
			sawRetrieval[key] = true
			if r.Interval = interval.Parse(r.When); r.Interval == nil {
				log.Printf("ERROR: config.sessions[%q].retrieve[%d].when = %q is not a valid interval", toCallSign, i, r.When)
				valid = false
			}
			if _, ok := c.BBSes[r.BBS]; !ok {
				log.Printf("ERROR: config.sessions[%q].retrieve[%d].bbs = %q is not a recognized BBS", toCallSign, i, r.BBS)
				valid = false
			}
			if r.Mailbox == "" {
				r.Mailbox = toCallSign
			} else if !tacticalCallRE.MatchString(r.Mailbox) {
				log.Printf("ERROR: config.sessions[%q].retrieve[%d].mailbox = %q is not a valid mailbox name", toCallSign, i, r.Mailbox)
				valid = false
			}
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
		for _, mtype := range ValidMessageTypes() {
			if mtype.TypeTag() == "plain" {
				continue
			}
			if c.MinimumVersions[mtype.TypeTag()] == "" {
				log.Printf("ERROR: config.minimumVersions[%q] is not specified", mtype.TypeTag())
				valid = false
			}
		}
	}

	// Parse the problems.
	if c.Problems == nil {
		log.Printf("ERROR: config.problems is not specified")
		valid = false
	} else {
		for code, problem := range c.Problems {
			probVars, ok := knownProbs[code]
			if !ok {
				log.Printf("ERROR: config.problems has entry for unknown problem %q", code)
				valid = false
				continue
			}
			for _, word := range strings.Fields(problem.Actions) {
				switch word {
				case "report":
					problem.ActionFlags |= ActionReport
				case "error":
					problem.ActionFlags |= ActionError
				case "dontcount":
					problem.ActionFlags |= ActionDontCount
				case "dropmsg":
					problem.ActionFlags |= ActionDropMsg
				default:
					log.Printf("ERROR: config.problems[%q].actions contains %q, which is not recognized", code, word)
					valid = false
				}
			}
			if problem.Response != "" {
				if problem.Responses == nil {
					problem.Responses = make(map[string]string)
				}
				problem.Responses[""] = problem.Response
			}
			if len(problem.Responses) != 0 {
				problem.ActionFlags |= ActionRespond
			}
			for resptag, response := range problem.Responses {
				for _, variable := range variableRE.FindAllString(response, -1) {
					if probVars != nil {
						if _, ok := probVars[variable[1:len(variable)-1]]; ok {
							continue
						}
					}
					if _, ok := knownVars[variable[1:len(variable)-1]]; ok {
						continue
					}
					log.Printf("ERROR: config.problems[%q].responses[%q] refers to unknown variable %q", code, resptag, variable)
					valid = false
				}
			}
		}
		for code := range knownProbs {
			if _, ok := c.Problems[code]; !ok {
				log.Printf("ERROR: config.problems has no entry for problem %q", code)
				valid = false
			}
		}
	}

	// Verify the jurisdiction abbreviations and fill out the map.
	if c.Jurisdictions == nil {
		log.Printf("ERROR: config.jurisdictions is not specified")
		valid = false
	} else {
		for name, abbr := range c.Jurisdictions {
			if len(abbr) < 1 || len(abbr) > 3 ||
				strings.IndexFunc(abbr, func(r rune) bool { return r < 'A' || r > 'Z' }) >= 0 {
				log.Printf("ERROR: config.jurisdictions[%q]: %q is not a valid abbreviation", name, abbr)
				valid = false
			}
			if upcase := strings.ToUpper(name); upcase != name {
				delete(c.Jurisdictions, name)
				c.Jurisdictions[upcase] = abbr
			}
			c.Jurisdictions[abbr] = abbr
		}
	}

	// Check that we have form routing for every form.
	if c.FormRouting == nil {
		log.Printf("ERROR: config.formRouting is not specified")
		valid = false
	} else {
		for _, mtype := range ValidMessageTypes() {
			if mtype.TypeTag() == "plain" {
				continue
			}
			if fr := c.FormRouting[mtype.TypeTag()]; fr == nil {
				log.Printf("ERROR: config.formRouting[%q] is not specified", mtype.TypeTag())
				valid = false
			} else {
				switch fr.HandlingOrder {
				case "":
					break
				case "computed":
					if _, ok := ComputedRecommendedHandlingOrder[mtype.TypeTag()]; !ok {
						log.Printf("ERROR: config.formRouting[%q].HandlingOrder = %q, but that form has no handling order computation defined", mtype.TypeTag(), fr.HandlingOrder)
						valid = false
					}
				default:
					if _, ok := xscmsg.ParseHandlingOrder(fr.HandlingOrder); !ok {
						log.Printf("ERROR: config.formRouting[%q].HandlingOrder = %q is not a valid handling order", mtype.TypeTag(), fr.HandlingOrder)
						valid = false
					}
				}
			}
		}
	}

	// Check that we have a listen address for the web server.
	if c.ListenAddr == "" {
		log.Printf("ERROR: config.listenAddr is not specified")
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

// validate makes sure that the interval specified in a scheduled value is
// parseable.
func (sv ScheduledValue) validate(label string) (valid bool) {
	valid = true
	for i, item := range sv {
		if item.When == "" {
			log.Printf("ERROR: %s[%d].when is not specified", label, i)
			valid = false
		} else if sv[i].Interval = interval.Parse(item.When); sv[i].Interval == nil {
			log.Printf("ERROR: %s[%d].when = %q is not a valid interval", label, i, item.When)
			valid = false
		}
	}
	return valid
}
