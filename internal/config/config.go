package config

// Config represents the application configuration.
type Config struct {
	URL              string
	Token            string
	OutputPath       string
	IncludeComments  bool
	IncludeMeta      bool
	Timeout          int
}

// DefaultConfig returns the default configuration.
func DefaultConfig() *Config {
	return &Config{
		IncludeComments: true,
		IncludeMeta:     true,
		Timeout:         30,
	}
}
