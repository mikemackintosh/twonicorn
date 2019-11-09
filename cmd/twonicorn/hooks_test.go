package main

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/mikemackintosh/twonicorn/internal/payloads"
)

var testDebug = false

func TestHookReceive(t *testing.T) {
	var testHookReceive = []struct {
		got    string
		format string
		want   string
	}{
		// Tests splunk payloads
		{
			got: `{
	"result": {
		"sourcetype" : "mongod",
		"count" : "8"
	},
	"sid" : "scheduler_admin_search_W2_at_14232356_132",
	"results_link" : "http://web.example.local:8000/app/search/@go?sid=scheduler_admin_search_W2_at_14232356_132",
	"search_name" : null,
	"owner" : "admin",
	"app" : "search"
}`,
			format: "{\"summary\": {{index . \"result\" \"count\"}}}",
			want:   "{\"summary\": 8}",
		},
		// Tests stackdriver payloads
		{
			got: `{
	"incident": {
	  "incident_id": "f2e08c333dc64cb09f75eaab355393bz",
	  "resource_id": "i-4a266a2d",
	  "resource_name": "webserver-85",
	  "state": "open",
	  "started_at": 1385085727,
	  "ended_at": null,
	  "policy_name": "Webserver Health",
	  "condition_name": "CPU usage",
	  "url": "https://app.google.stackdriver.com/incidents/f333dc64z",
	  "summary": "CPU for webserver-85 is above the threshold of 1% with a value of 28.5%"
	},
  "version": 1.1
}`,
			format: "{\"summary\": \"{{index . \"incident\" \"state\"}} incident of {{index . \"incident\" \"summary\"}}\"}",
			want:   "{\"summary\": \"open incident of CPU for webserver-85 is above the threshold of 1% with a value of 28.5%\"}",
		},
	}

	for _, test := range testHookReceive {
		var payload payloads.Payload

		if testDebug {
			fmt.Printf("\033[38;5;154m---> Received payload:\033[0m \n%v\n", test.got)
			fmt.Printf("\033[38;5;154m---> Wanted Format:\033[0m \n%v\n", test.format)
		}

		r, err := http.NewRequest("GET", "/", strings.NewReader(test.got))
		if err != nil {
			t.Errorf("[err] Making test request: %s", err)
		}

		r.Header.Add("Content-type", "application/json")

		/*
			o, _ := httputil.DumpRequest(r, true)
			fmt.Printf("%s", string(o))
		*/

		err = payloads.DecodeReceivedPayload(r, &payload)
		if err != nil {
			t.Errorf("want %s, got %s", test.want, err)
		}

		if testDebug {
			fmt.Printf("\033[38;5;154m---> Parsed payload:\033[0m \n%#v\n", payload)
		}
		// Creates the target paload
		output, err := payloads.CreateTargetPayload(test.format, payload)
		if err != nil {
			t.Errorf("want %s, got %s", test.want, err)
		}

		if testDebug {
			fmt.Printf("\033[38;5;154m---> Parsed output:\033[0m \n%#v\n\n\n", output.String())
		}

		if output.String() != test.want {
			t.Errorf("want %s, got %s", test.want, output.String())
		}
	}
}
