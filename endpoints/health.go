package endpoints

import (
	"encoding/json"
	"net/http"

	"github.com/supabase-community/auth-go/types"
)

var healthPath = "/health"

// GET /health
//
// Check the health of the Auth server.
func (c *Client) HealthCheck() (*types.HealthCheckResponse, error) {
	r, err := c.newRequest(healthPath, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, handleErrorResponse(resp)
	}

	var res types.HealthCheckResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
