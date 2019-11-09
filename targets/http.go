package target

import (
	"context"
	"io"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/mikemackintosh/twonicorn/config"
	"github.com/mikemackintosh/twonicorn/internal/log"
)

// HTTPTargetHandler defines a custom handler
type HTTPTargetHandler struct {
	method  string
	body    io.Reader
	url     string
	auth    string
	headers map[string]string
}

func (h *HTTPTargetHandler) SetURL(url string) {
	h.url = url
}

func (h *HTTPTargetHandler) SetMethod(method string) {
	h.method = method
}

func (h *HTTPTargetHandler) SetBody(body io.Reader) {
	h.body = body
}

func (h *HTTPTargetHandler) SetHeaders(headers map[string]string) {
	h.headers = headers
}

func (h *HTTPTargetHandler) MakeRequest() (*http.Request, error) {
	req, err := http.NewRequest(strings.ToUpper(h.method), h.url, h.body)
	if len(h.headers) > 0 {
		for hkey, hval := range h.headers {
			req.Header.Set(hkey, hval)
		}
	}

	return req, err
}

// NewHTTPTarget sets the default target settings
func NewHTTPTarget() *HTTPTargetHandler {
	var client = HTTPTargetHandler{}
	return &client
}

// NewHTTPTargetWithOptions configures an HTTP Handler with required options
func NewHTTPTargetWithOptions(method, url string, body io.Reader, headers map[string]string) *HTTPTargetHandler {
	var client = HTTPTargetHandler{
		method:  method,
		url:     url,
		body:    body,
		headers: headers,
	}
	return &client
}

// SetOptions will take an TargetHandler and apply all configurations
func (t *HTTPTargetHandler) SetOptions(method, url string, body io.Reader, headers map[string]string) *HTTPTargetHandler {
	t.method = method
	t.url = url
	t.body = body
	t.headers = headers

	return t
}

// Create performs a create action
func (t *HTTPTargetHandler) Send(ctx context.Context) (*http.Response, error) {
	req, err := t.MakeRequest()
	if err != nil {
		return nil, err
	}

	req.WithContext(ctx)

	if config.IsDebug() {
		x, _ := httputil.DumpRequest(req, true)
		log.Debugf(ctx, string(x))
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func init() {
	RegisterTarget("http", NewHTTPTarget())
}
