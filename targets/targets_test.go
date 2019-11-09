package target

import (
	"context"
	"net/http"
	"testing"
)

type testTargetHandler Handler

func (t *testTargetHandler) Send(ctx context.Context) (*http.Response, error) {
	return nil, nil
}
func resetRegistry() {
	targetRegistry = map[string]Target{}
}

func TestRegisterTarget(t *testing.T) {
	// Should be empty
	beforeCount := len(targetRegistry)
	expectedCount := beforeCount + 1

	// Register the test handler
	RegisterTarget("example", &testTargetHandler{})

	afterCount := len(targetRegistry)
	if expectedCount != afterCount {
		t.Errorf("want %d, got %d", expectedCount, afterCount)
	}

	resetRegistry()
}

func TestMustGetTarget(t *testing.T) {
	targetRegistry["example"] = &testTargetHandler{}
	module := MustGetTarget("example")
	if module == nil {
		t.Errorf("want type Handler, got %T", module)
	}

	resetRegistry()
}

func TestGetTarget(t *testing.T) {
	targetRegistry["example"] = &testTargetHandler{}
	_, err := GetTarget("example")
	if err != nil {
		t.Errorf("want nil, got %s", err)
	}

	resetRegistry()
}
