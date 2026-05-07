// Package config handles reading and parsing the config.yaml file, which
// contains the configuration for the server.
package config

import (
	"errors"
	"log"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

// Config holds all of the configuration data.
type Config struct {
	BBSes           map[string]*BBSConfig         `yaml:"bbses"`
	MinPIFOVersion  string                        `yaml:"minPIFOVersion"`
	MessageTypes    map[string]*MessageTypeConfig `yaml:"messageTypes"`
	ServerURL       string                        `yaml:"serverURL"`
	ListenAddr      string                        `yaml:"listenAddr"`
	SMTP            *SMTPConfig                   `yaml:"smtp"`
	CanViewEveryone []string                      `yaml:"canViewEveryone"`
	CanEditSessions []string                      `yaml:"canEditSessions"`
}

// An SMTPConfig describes how to send email via SMTP.
type SMTPConfig struct {
	From     string `yaml:"from"`
	Server   string `yaml:"server"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// A MessageTypeConfig structure specifies the configuration for a particular
// message type.
type MessageTypeConfig struct {
	MinimumVersion string   `yaml:"minimumVersion"`
	HandlingOrder  string   `yaml:"handlingOrder"`
	ToICSPosition  []string `yaml:"toICSPosition"`
	ToLocation     []string `yaml:"toLocation"`
}

// BBSConfig holds the configuration of a single BBS.
type BBSConfig struct {
	Transport string            `yaml:"transport"`
	AX25      string            `yaml:"ax25"`
	TCP       string            `yaml:"tcp"`
	Passwords map[string]string `yaml:"passwords"`
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
