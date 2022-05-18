// Package config handles reading and parsing the config.yaml file, which
// contains the configuration for the server.
package config

import (
	"errors"
	"log"
	"os"
	"sync"
	"time"

	"github.com/rothskeller/packet/wppsvr/interval"
	"gopkg.in/yaml.v3"
)

// PackItForms is the key into the Config.MinimumVersions map for retrieving the
// minimum version of the PackItForms encoding.
const PackItForms = "PackItForms"

// Config holds all of the configuration data.
type Config struct {
	BBSes           map[string]*BBSConfig     `yaml:"bbses"`
	SessionDefaults *SessionConfig            `yaml:"sessionDefaults"`
	Sessions        map[string]*SessionConfig `yaml:"sessions"`
	MinimumVersions map[string]string         `yaml:"minimumVersions"`
	ProblemActions  map[string]string         `yaml:"problemActions"`
	FormRouting     map[string]*FormRouting   `yaml:"formRouting"`
	ListenAddr      string                    `yaml:"listenAddr"`
	CanViewEveryone []string                  `yaml:"canViewEveryone"`
	CanEditSessions []string                  `yaml:"canEditSessions"`

	ProblemActionFlags map[string]Action `yaml:"-"`
}

// A FormRouting structure specifies the required form routing for a particular
// form type.
type FormRouting struct {
	HandlingOrder string   `yaml:"HandlingOrder"`
	ToICSPosition []string `yaml:"ToICSPosition"`
	ToLocation    []string `yaml:"ToLocation"`
}

// Action is a flag, or a bitmask of flags, describing the action(s) to take in
// response to a detected problem.
type Action uint8

// Values for Action:
const (
	ActionRespond Action = 1 << iota
	ActionReport
	ActionError
	ActionDontCount
	ActionDropMsg
)

// BBSConfig holds the configuration of a single BBS.
type BBSConfig struct {
	Transport string            `yaml:"transport"`
	AX25      string            `yaml:"ax25"`
	TCP       string            `yaml:"tcp"`
	Passwords map[string]string `yaml:"passwords"`
}

// SessionConfig holds the default configuration of a single session.
type SessionConfig struct {
	Name            string         `yaml:"name"`
	Prefix          string         `yaml:"prefix"`
	Start           string         `yaml:"start"`
	End             string         `yaml:"end"`
	ReportTo        []string       `yaml:"reportTo"`
	ExcludeFromWeek bool           `yaml:"excludeFromWeek"`
	ToBBSes         ScheduledValue `yaml:"toBBSes"`
	DownBBSes       ScheduledValue `yaml:"downBBSes"`
	Retrieve        []*Retrieval   `yaml:"retrieve"`
	MessageTypes    ScheduledValue `yaml:"messageTypes"`

	StartInterval interval.Interval `yaml:"-"`
	EndInterval   interval.Interval `yaml:"-"`
}

// Retrieval holds the configuration of a single retrieval for a session.
type Retrieval struct {
	When              string `yaml:"when"`
	BBS               string `yaml:"bbs"`
	Mailbox           string `yaml:"mailbox"`
	DontKillMessages  bool   `yaml:"dontKillMessages"`
	DontSendResponses bool   `yaml:"dontSendResponses"`

	Interval interval.Interval `yaml:"-"`
}

// ScheduledValue holds a value that changes on a set schedule.  The value may
// be either a single string or a list of strings, depending on context.
type ScheduledValue []struct {
	When string
	Then string

	Interval interval.Interval `yaml:"-"`
}

var (
	config *Config      // current configuration
	mutex  sync.RWMutex // mutex for access to configuration
)

// Get returns the current configuration.  Successive calls to Get may return
// different configurations if the config file has been changed in the interim.
func Get() *Config {
	mutex.RLock()
	defer mutex.RUnlock()
	return config
}

// Read reads the system configuration from the config.yaml file.  If an error
// occurs, the previous configuration is retained and the error is returned.
func Read() (err error) {
	var (
		newconfig Config
		configFH  *os.File
		decoder   *yaml.Decoder
	)
	if configFH, err = os.Open("config.yaml"); err != nil {
		log.Printf("ERROR: opening config file: %s", err)
		return err
	}
	defer configFH.Close()
	decoder = yaml.NewDecoder(configFH)
	decoder.KnownFields(true)
	if err = decoder.Decode(&newconfig); err != nil {
		log.Printf("ERROR: parsing config file: %s", err)
		return err
	}
	newconfig.applySessionDefaults()
	if !newconfig.Validate() {
		return errors.New("invalid configuration data")
	}
	SetConfig(&newconfig)
	return nil
}

// SetConfig sets the configuration that will be returned by subsequent calls to
// Get.  It should not be called by production code; it is intended for testing
// only.
func SetConfig(newconfig *Config) {
	mutex.Lock()
	config = newconfig
	mutex.Unlock()
}

// applySessionDefaults applies defaults from the SessionDefaults section to all
// sessions that don't have explicit settings for those parameters.
func (c *Config) applySessionDefaults() {
	if c.SessionDefaults == nil {
		return
	}
	for _, session := range c.Sessions {
		if session.ReportTo == nil {
			session.ReportTo = c.SessionDefaults.ReportTo
		}
		if session.ToBBSes == nil {
			session.ToBBSes = c.SessionDefaults.ToBBSes
		}
		if session.DownBBSes == nil {
			session.DownBBSes = c.SessionDefaults.DownBBSes
		}
		if session.MessageTypes == nil {
			session.MessageTypes = c.SessionDefaults.MessageTypes
		}
	}
}

// For finds the first item in the ScheduledValue whose "when" clause matches
// the specified time, and returns the corresponding value.  If no clauses
// match, the function returns an empty string.
func (sv ScheduledValue) For(t time.Time) string {
	for _, item := range sv {
		if item.Interval.Match(t) {
			return item.Then
		}
	}
	return ""
}

// AllFor finds all items in the ScheduledValue whose "when" clauses match the
// specified time, and returns the corresponding values.
func (sv ScheduledValue) AllFor(t time.Time) (list []string) {
	for _, item := range sv {
		if item.Interval.Match(t) {
			list = append(list, item.Then)
		}
	}
	return list
}
