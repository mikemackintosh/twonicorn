package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mikemackintosh/twonicorn/config"
	"github.com/rs/xid"
)

var (
	flagConfigFile string
	flagPort       string
	flagOutput     string
	flagDebug      bool
)

func init() {
	flag.StringVar(&flagOutput, "o", "json", "Output Logging mechanism\nAvailable Outputs: "+strings.Join(supportedOutput, ", "))
	flag.StringVar(&flagPort, "p", "localhost:8080", "Local address:port to listen on")
	flag.StringVar(&flagConfigFile, "c", "", "Configuration file to execute")
	flag.BoolVar(&flagDebug, "d", false, "Enable debug or not")
}

func main() {
	flag.Parse()

	if len(flagConfigFile) == 0 {
		fmt.Println("Please provide a valid configuration file")
		os.Exit(0)
	}

	// Set debug status
	config.Registry.Debug = flagDebug

	rawConfig, err := os.Open(flagConfigFile)
	if err != nil {
		panic(fmt.Sprintf("Error opening configuration: %s", err))
	}

	hooks, err := config.ParseConfig(bufio.NewReader(rawConfig))
	if err != nil {
		panic(fmt.Sprintf("Error parsing configuration: %s", err))
	}

	http.HandleFunc("/hooks/", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := LoggingResponseWriter{ResponseWriter: w}
		reqId := xid.New().String()
		r.WithContext(context.WithValue(r.Context(), "reqid", reqId))
		HookHandler(&sw, r, hooks)
		duration := time.Now().Sub(start)

		logEntry := LogEntry{
			RequestID:  reqId,
			Host:       r.Host,
			RemoteAddr: r.RemoteAddr,
			Method:     r.Method,
			RequestURI: r.RequestURI,
			Proto:      r.Proto,
			Status:     sw.status,
			Time:       time.Now().Format("02/Jan/2006:15:04:05 -0700"),
			ContentLen: sw.length,
			Referer:    r.Referer(),
			UserAgent:  r.UserAgent(),
			Duration:   duration,
		}

		go func(flagOutput string, logEntry LogEntry) {
			switch flagOutput {
			case "yaml":
				logEntry.ToYaml()
			case "json":
				logEntry.ToJSON()
			case "kv":
				logEntry.ToKeyValue()
			default:
				logEntry.ToApache()
			}
		}(flagOutput, logEntry)

	})

	log.Printf("Listening on %s\n", flagPort)
	http.ListenAndServe(flagPort, http.DefaultServeMux)
}
