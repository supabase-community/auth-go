package endpoints

import (
	"encoding/json"
	"net/http"

	"github.com/supabase-community/auth-go/types"
)

var settingsPath = "/settings"

// GET /settings
//
// Returns the publicly available settings for this auth instance.
func (c *Client) GetSettings() (*types.SettingsResponse, error) {
	r, err := c.newRequest(settingsPath, http.MethodGet, nil)
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

	var res types.SettingsResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
