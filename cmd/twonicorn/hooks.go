package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"path/filepath"

	"github.com/mikemackintosh/twonicorn/config"
	"github.com/mikemackintosh/twonicorn/internal/log"
	"github.com/mikemackintosh/twonicorn/internal/payloads"
	target "github.com/mikemackintosh/twonicorn/targets"
)

// HookHandler will get and parse an incoming request
func HookHandler(w http.ResponseWriter, r *http.Request, hooks *config.Entries) {
	var payload payloads.Payload
	var responseFromUpstream = respInvalidReq
	var requestError error
	var targetResponse *http.Response

	// Create a new context wrapper
	wrapper := NewContext(w, r)
	key := path.Base(r.URL.Path)
	prefix := path.Base(filepath.Dir(r.URL.Path))

	// Validate the key provided
	if len(key) == 0 {
		log.Infof(wrapper.Ctx, "Invalid key detected")
		wrapper.Write(respInvalidReq)
		return
	}

	// Piece together the key
	pathKey := fmt.Sprintf("%s%s", prefix, key)
	// POC Debugging
	log.Debugf(wrapper.Ctx, " --> Looking up key: %s", pathKey)
	hook, err := hooks.GetKey(pathKey)
	if err != nil {
		log.Infof(wrapper.Ctx, "Invalid key detected")
		wrapper.Write(respInvalidReq)
		return
	}

	// Looks like we found a hook
	log.Debugf(wrapper.Ctx, " --> Found hook: %s", hook.Name)

	// Parse the request paylaod into an interface, so we can handle it
	err = payloads.DecodeReceivedPayload(r, &payload)
	if err != nil {
		log.Infof(wrapper.Ctx, "Invalid payload decode: %s", err)
		wrapper.Write(respInvalidReq)
		return
	}

	// POC Debugging
	log.Debugf(wrapper.Ctx, " --> Payload: %+v", payload)
	log.Debugf(wrapper.Ctx, " --> Translate to: %+v", *hook.Translate)

	targetPayload, err := payloads.ParseTargetPayloadFormat(hook.Translate)
	if err != nil {
		log.Infof(wrapper.Ctx, "Invalid target payload marshal: %s", err)
		wrapper.Write(respInvalidReq)
		return
	}

	// Creates the target paload
	output, err := payloads.CreateTargetPayload(string(targetPayload), payload)
	if err != nil {
		log.Infof(wrapper.Ctx, "Invalid target payload created: %s", err)
		wrapper.Write(respInvalidReq)
		return
	}

	// POC Debugging
	log.Debugf(wrapper.Ctx, output.String())

	switch hook.Target {
	case "http":
		targetResponse, requestError = target.NewHTTPTargetWithOptions(
			hook.Method,
			hook.URL,
			ioutil.NopCloser(bytes.NewReader(output.Bytes())),
			hook.Headers,
		).Send(wrapper.Ctx)
	}

	// Check for an error from the target
	defer targetResponse.Body.Close()
	if requestError != nil {
		log.Infof(wrapper.Ctx, "Failed to send request: %s", err)
		wrapper.Write(respInvalidReq)
		return
	}

	// TODO: Put this in here, but don't do callbacks yet
	if hook.On != nil {
		// Let's check if the request failed or not
		switch targetResponse.StatusCode {
		case hook.On.SuccessCode:
		case hook.On.FailCode:
		default:
		}
	}

	// Get the response body and append it to the payload
	b, err := ioutil.ReadAll(targetResponse.Body)
	responseFromUpstream = ResponsePayload{StatusCode: targetResponse.StatusCode, Message: map[string]string{
		"upstream": string(b),
		"ok":       "Successful request",
	}}

	wrapper.Write(responseFromUpstream)
	return
}
