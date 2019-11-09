package main

import (
	"testing"
)

// TestResponsePayload checks that the correct status codes are used
// in respnose to different event types to ensure compatability with
// sending services
func TestResponsePayload(t *testing.T) {
	var testResponsePayload = []struct {
		got  ResponsePayload
		want int
	}{
		{
			got:  respSuccessfulReq,
			want: 200,
		},
		{
			got:  respInvalidReq,
			want: 400,
		},
		{
			got:  respQueuedReq,
			want: 201,
		},
	}

	for _, test := range testResponsePayload {
		if test.got.StatusCode != test.want {
			t.Errorf("want %d, got %d", test.got.StatusCode, test.want)
		}
	}
}
