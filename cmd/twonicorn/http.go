package main

import (
	"context"
	"encoding/json"
	"net/http"
)

// RequestWrapper does things
type RequestWrapper struct {
	Ctx     context.Context
	Request *http.Request
	Writer  http.ResponseWriter
}

// WriteResponse users a responserwriter and response struct to
// respond to the request
func (r *RequestWrapper) Write(resp ResponsePayload) {
	payload, err := json.Marshal(resp.Message)
	if err != nil {
		panic(err)
	}

	r.Writer.WriteHeader(resp.StatusCode)
	_, _ = r.Writer.Write(payload)
}

// ResponsePayload is a struct for responses
type ResponsePayload struct {
	StatusCode int
	Message    map[string]string
}

// respSuccessfulReq is a default invalid request
var respSuccessfulReq = ResponsePayload{
	StatusCode: http.StatusOK,
	Message:    map[string]string{"ok": "Successful request"},
}

// respQueuedReq indicates the request payload was queued
var respQueuedReq = ResponsePayload{
	StatusCode: http.StatusCreated,
	Message:    map[string]string{"ok": "Request queued"},
}

// respInvalidReq is a default invalid request
var respInvalidReq = ResponsePayload{
	StatusCode: http.StatusBadRequest,
	Message:    map[string]string{"err": "Invalid request"},
}

// NewContext creates a new context
func NewContext(w http.ResponseWriter, r *http.Request) *RequestWrapper {
	request := RequestWrapper{
		Writer:  w,
		Request: r,
		Ctx:     r.Context(),
	}

	return &request
}

type LoggingResponseWriter struct {
	http.ResponseWriter
	status int
	length int
}

func (w *LoggingResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *LoggingResponseWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}
	n, err := w.ResponseWriter.Write(b)
	w.length += n
	return n, err
}
