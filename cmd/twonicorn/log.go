package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"

	yaml "gopkg.in/yaml.v3"
)

var (
	// supportedOutput are currently supported
	supportedOutput = []string{
		"yaml",
		"kv",
		"json",
		"apache",
	}
)

const (
	// ApachePattern is a apache common-like access log pattern
	ApachePattern = "%s - - [%s] \"%s\" %d %d %f -- %s\n"
)

// LogEntry records received
type LogEntry struct {
	RequestID  string
	Host       string
	RemoteAddr string
	Method     string
	RequestURI string
	Proto      string
	Status     int
	ContentLen int
	UserAgent  string
	Time       string
	Referer    string
	Duration   time.Duration
}

// ToJSON will output the LogEntry payload in JSON format
func (l *LogEntry) ToJSON() {
	b, _ := json.Marshal(l)
	fmt.Fprintf(os.Stderr, string(b)+"\n")
}

// ToYAML will output the payload in YAML format
func (l *LogEntry) ToYaml() {
	b, _ := yaml.Marshal(l)
	fmt.Fprintf(os.Stderr, string(b)+"\n")
}

// ToApache will output LogEntry into an Apache-like access log entry
func (l *LogEntry) ToApache() {
	fmt.Fprintf(os.Stderr, ApachePattern,
		l.RemoteAddr,
		l.Time,
		fmt.Sprintf("%s %s %s", l.Method, l.RequestURI, l.Proto),
		l.Status,
		l.ContentLen,
		l.Duration.Seconds(),
		l.RequestID,
	)
}

// ToKeyValue will output LogEntry as a key=value collection
func (l *LogEntry) ToKeyValue() {
	v := reflect.ValueOf(*l)
	for i := 0; i < v.NumField(); i++ {
		val := fmt.Sprintf("%v", v.Field(i).Interface())
		fmt.Printf("%s=%v ", v.Type().Field(i).Name, strconv.QuoteToASCII(val))
	}
	fmt.Println()
}
