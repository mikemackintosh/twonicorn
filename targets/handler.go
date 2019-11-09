package target

import (
	"io"
	"net/http"
)

// Handler holds the connection meta for a target
type Handler struct {
	method  string
	body    io.Reader
	url     string
	auth    string
	headers map[string]string
}

func (h *Handler) SetURL(url string) {
	h.url = url
}

func (h *Handler) SetMethod(method string) {
	h.method = method
}

func (h *Handler) SetBody(body io.Reader) {
	h.body = body
}

func (h *Handler) SetHeaders(headers map[string]string) {
	h.headers = headers
}

func (h *Handler) MakeRequest() (*http.Request, error) {
	req, err := http.NewRequest(h.method, h.url, h.body)
	if len(h.headers) > 0 {
		for hkey, hval := range h.headers {
			req.Header.Set(hkey, hval)
		}
	}

	return req, err
}
