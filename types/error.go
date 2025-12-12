package types

import "fmt"

// AuthError represents a structured error from Supabase Auth API.
// It provides access to the HTTP status code, error message, error code,
// and additional details from the API response.
type AuthError struct {
	// StatusCode is the HTTP status code returned by the API.
	StatusCode int

	// Message is the error message from the API.
	Message string

	// ErrorCode is the Supabase-specific error code, if available.
	ErrorCode string

	// Details contains additional error information, if available.
	Details map[string]interface{}
}

// Error implements the error interface.
func (e *AuthError) Error() string {
	if e.ErrorCode != "" {
		return fmt.Sprintf("auth error (status %d, code %s): %s",
			e.StatusCode, e.ErrorCode, e.Message)
	}
	return fmt.Sprintf("auth error (status %d): %s", e.StatusCode, e.Message)
}
