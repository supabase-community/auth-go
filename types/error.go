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
// The error message format matches the previous fmt.Errorf format for backward compatibility.
func (e *AuthError) Error() string {
	if e.ErrorCode != "" {
		// Include error code if available, but keep the old format for compatibility
		return fmt.Sprintf("response status code %d (error_code: %s): %s",
			e.StatusCode, e.ErrorCode, e.Message)
	}
	// Match the old format: "response status code %d: %s"
	return fmt.Sprintf("response status code %d: %s", e.StatusCode, e.Message)
}
