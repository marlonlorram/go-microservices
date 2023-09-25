package config

import (
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"github.com/kelseyhightower/envconfig"
)

// Parse reads and maps environment variables to the provided type.
func Parse[T any](prefixes ...string) (*T, error) {
	var c T
	var prefix string

	if len(prefixes) > 0 {
		prefix = prefixes[0]
	}

	// Process the environment variables
	if err := envconfig.Process(prefix, &c); err != nil {
		return nil, fault.Wrap(err, fmsg.With("failed to process environment variables"))
	}

	return &c, nil
}
