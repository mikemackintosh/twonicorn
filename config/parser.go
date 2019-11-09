package config

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"net/http"

	yaml "gopkg.in/yaml.v3"
)

var validMethods = []string{
	http.MethodGet,
	http.MethodHead,
	http.MethodPost,
	http.MethodPut,
	http.MethodPatch,
	http.MethodDelete,
	http.MethodConnect,
	http.MethodOptions,
	http.MethodTrace,
}

// Entries is the master collection of Entry items
type Entries map[string]Entry

// Entry is a configuration entry
type Entry struct {
	PathKey string `yaml:"-"`
	Name    string `yaml:"-"`
	// Source    string `yaml:"source"`
	Target    string            `yaml:"target"`
	Method    string            `yaml:"method"`
	URL       string            `yaml:"url"`
	Owner     string            `yaml:"owner"`
	Headers   map[string]string `yaml:"headers"`
	Translate *map[string]interface{}
	On        *struct {
		SuccessCode int        `yaml:"success_code"`
		Success     OnCallback `yaml:"success"`
		FailCode    int        `yaml:"fail_code"`
		Fail        OnCallback `yaml:"fail"`
		Other       OnCallback `yaml:"other"`
	} `yaml:"on,omitempty"`
}

// OnCallback handles payloads for result handlers
type OnCallback struct {
	Target  string            `yaml:"target"`
	Payload OnCallbackPayload `yaml:"payload"`
}

// OnCallbackPayload
type OnCallbackPayload map[string]interface{}

// Validate will validate a config entry
func (e *Entry) Validate() error {
	if len(e.Target) == 0 {
		return errors.New("invalid target provided")
	} else if len(e.URL) == 0 {
		return errors.New("invalid url provided")
	} else if len(e.Method) == 0 {
		return errors.New("invalid method provided")
	} else if len(e.Owner) == 0 {
		return errors.New("invalid owner provided")
	} else if e.Translate == nil {
		return errors.New("invalid translate provided")
	}

	// check that the supplied method in the config is valid
	if !Contains(validMethods, e.Method) {
		return errors.New("Invalid method")
	}

	return nil
}

// GetKey returns the Entry for the corresponding key
func (e *Entries) GetKey(key string) (*Entry, error) {
	if entry, ok := (*e)[key]; ok {
		return &entry, nil
	}
	return nil, errors.New("invalid key")
}

// ParseConfig will parse a configuration
func ParseConfig(c io.Reader) (*Entries, error) {
	var parsedEntries Entries
	err := yaml.NewDecoder(c).Decode(&parsedEntries)
	if err != nil {
		return nil, err
	}

	var entries = Entries{}
	for n, e := range parsedEntries {
		hash := ComputeKey(n)
		e.Name = n
		e.PathKey = hash
		entries[hash] = e
	}

	return &entries, nil
}

// ComputeKey will compute the key for the config entry
func ComputeKey(s string) string {
	hasher := sha256.New()
	_, _ = hasher.Write([]byte(s))
	return hex.EncodeToString(hasher.Sum(nil))[:22]
}

func Contains(collection interface{}, x interface{}) bool {
	switch collection.(type) {
	case []string:
		for _, n := range collection.([]string) {
			if x.(string) == n {
				return true
			}
		}
	}

	return false
}
