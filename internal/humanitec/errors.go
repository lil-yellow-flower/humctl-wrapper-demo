package humanitec

import "errors"

var (
	// ErrMissingAPIToken is returned when the API token is not set
	ErrMissingAPIToken = errors.New("HUMANITEC_API_TOKEN environment variable is not set")
) 