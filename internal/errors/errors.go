package errors

import "fmt"

// App-level error messages used across the codebase.
// Use these constants instead of hardcoded strings with fmt.Errorf.

const (
	// Command errors
	ErrNoCommandEntered = "no command entered. Please enter a command"
	ErrUnknownCommand   = "unknown command: %s. Type 'help' for available commands"
	ErrCommandNeedsArg  = "command %s requires an argument"

	// Data fetching errors
	ErrFetchData          = "failed to fetch data: %w"
	ErrMarshalData        = "failed to marshal data: %w"
	ErrUnmarshalData      = "failed to unmarshal data: %w"
	ErrDecodeResponseBody = "failed to decode response body: %w"
	ErrMakeRequest        = "failed to make GET request: %w"
	ErrNon200Response     = "received non-200 response code: %d"
	ErrGetBaseExperience  = "failed to get base experience: %w"

	// Cache errors
	ErrCacheAdd = "failed to add item to cache"
	ErrCacheGet = "failed to get item from cache"

	// Input/IO errors
	ErrScannerIO = "error reading input: %v"
)

// HTTP status code constants for structured error handling.

const (
	StatusBadRequest      = 400
	StatusUnauthorized    = 401
	StatusForbidden       = 403
	StatusNotFound        = 404
	StatusTooManyRequests = 429
	StatusInternalError   = 500
	StatusBadGateway      = 502
	StatusServiceUnavail  = 503
)

// UserMessage returns a human-readable message for the given HTTP status code.
func UserMessage(code int) string {
	switch code {
	case StatusBadRequest:
		return "bad request — check your input"
	case StatusNotFound:
		return "resource not found — the location or Pokemon does not exist"
	case StatusTooManyRequests:
		return "too many requests — slow down and try again"
	case StatusInternalError:
		return "internal server error — the API is having issues"
	case StatusBadGateway:
		return "bad gateway — the API might be down"
	case StatusServiceUnavail:
		return "service unavailable — try again later"
	default:
		return fmt.Sprintf("unexpected error (status %d)", code)
	}
}
