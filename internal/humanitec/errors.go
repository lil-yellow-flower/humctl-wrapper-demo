package humanitec

import "errors"

var (
	// ErrMissingAPIToken is returned when the API token is not set
	ErrMissingAPIToken = errors.New("Humanitec API token is not set in config.yaml")
) 