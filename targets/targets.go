package target

import (
	"context"
	"fmt"
	"net/http"
)

// targetRegistry
var targetRegistry = map[string]Target{}

// Target defines required target methods
type Target interface {
	Send(context.Context) (*http.Response, error)
}

// RegisterTarget registers a new target configuration
func RegisterTarget(name string, t Target) {
	targetRegistry[name] = t
}

// MustGetTarget ignores errors and returns nil for missing targets
func MustGetTarget(name string) Target {
	target, _ := GetTarget(name)
	return target
}

// GetTarget attempts to get a target, or returns an error
func GetTarget(name string) (Target, error) {
	if v, ok := targetRegistry[name]; ok {
		return v, nil
	}

	return nil, fmt.Errorf("target %q is not registered", name)
}
