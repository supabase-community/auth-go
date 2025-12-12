package endpoints

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/supabase-community/auth-go/types"
)

var userPath = "/user"

// GET /user
//
// Get the JSON object for the logged in user (requires authentication)
func (c *Client) GetUser() (*types.UserResponse, error) {
	r, err := c.newRequest(userPath, http.MethodGet, nil)
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

	var res types.UserResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// PUT /user
//
// Update a user (Requires authentication). Apart from changing email/password,
// this method can be used to set custom user data. Changing the email will
// result in a magiclink being sent out.
func (c *Client) UpdateUser(req types.UpdateUserRequest) (*types.UpdateUserResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	r, err := c.newRequest(userPath, http.MethodPut, bytes.NewBuffer(body))
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

	var res types.UpdateUserResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
