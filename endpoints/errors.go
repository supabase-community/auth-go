package endpoints

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/supabase-community/auth-go/types"
)

// handleErrorResponse parses an HTTP error response and returns a structured AuthError.
// It attempts to parse the response body as JSON to extract error details.
// If JSON parsing fails, it falls back to using the raw response body as the error message.
func handleErrorResponse(resp *http.Response) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// If we can't read the body, return a basic error with just the status code
		return &types.AuthError{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("failed to read error response body: %v", err),
		}
	}

	// Try parsing as Supabase error response
	var apiErr struct {
		Error     string                 `json:"error"`
		ErrorCode string                 `json:"error_code"`
		Message   string                 `json:"message"`
		Details   map[string]interface{} `json:"details"`
	}

	if json.Unmarshal(body, &apiErr) == nil {
		// Use the parsed error fields
		message := apiErr.Error
		if message == "" {
			message = apiErr.Message
		}
		if message == "" {
			message = string(body)
		}

		return &types.AuthError{
			StatusCode: resp.StatusCode,
			Message:    message,
			ErrorCode:  apiErr.ErrorCode,
			Details:    apiErr.Details,
		}
	}

	// Fallback to text response
	return &types.AuthError{
		StatusCode: resp.StatusCode,
		Message:    string(body),
	}
}
