package config

// Registry contains the active configuration
var Registry = &Config{}

// Default config object
type Config struct {
	Debug bool
}

// IsDebug returns if debugging is enabled or not
// Default is false
func IsDebug() bool {
	if Registry == nil {
		return false
	}
	return Registry.Debug
}
