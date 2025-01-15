package endpoints

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/supabase-community/auth-go/types"
)

const adminUsersPath = "/admin/users"

// POST /admin/users
//
// Creates the user based on the user_id specified.
func (c *Client) AdminCreateUser(req types.AdminCreateUserRequest) (*types.AdminCreateUserResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	r, err := c.newRequest(adminUsersPath, http.MethodPost, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fullBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("response status code %d", resp.StatusCode)
		}
		return nil, fmt.Errorf("response status code %d: %s", resp.StatusCode, fullBody)
	}

	var res types.AdminCreateUserResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// GET /admin/users
//
// Get a list of users.
func (c *Client) AdminListUsers(req types.AdminListUsersRequest) (*types.AdminListUsersResponse, error) {
	r, err := c.newRequest(adminUsersPath, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	q := r.URL.Query()
	if req.Page != nil {
		q.Add("page", fmt.Sprintf("%d", *req.Page))
	}
	if req.PerPage != nil {
		q.Add("per_page", fmt.Sprintf("%d", *req.PerPage))
	}
	r.URL.RawQuery = q.Encode()

	resp, err := c.client.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fullBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("response status code %d", resp.StatusCode)
		}
		return nil, fmt.Errorf("response status code %d: %s", resp.StatusCode, fullBody)
	}

	var res types.AdminListUsersResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// GET /admin/users/{user_id}
//
// Get a user by their user_id.
func (c *Client) AdminGetUser(req types.AdminGetUserRequest) (*types.AdminGetUserResponse, error) {
	path := fmt.Sprintf("%s/%s", adminUsersPath, req.UserID)
	r, err := c.newRequest(path, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fullBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("response status code %d", resp.StatusCode)
		}
		return nil, fmt.Errorf("response status code %d: %s", resp.StatusCode, fullBody)
	}

	var res types.AdminGetUserResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// PUT /admin/users/{user_id}
//
// Update a user by their user_id.
func (c *Client) AdminUpdateUser(req types.AdminUpdateUserRequest) (*types.AdminUpdateUserResponse, error) {
	path := fmt.Sprintf("%s/%s", adminUsersPath, req.UserID)
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	r, err := c.newRequest(path, http.MethodPut, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fullBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("response status code %d", resp.StatusCode)
		}
		return nil, fmt.Errorf("response status code %d: %s", resp.StatusCode, fullBody)
	}

	var res types.AdminUpdateUserResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// DELETE /admin/users/{user_id}
//
// Delete a user by their user_id.
func (c *Client) AdminDeleteUser(req types.AdminDeleteUserRequest) error {
	path := fmt.Sprintf("%s/%s", adminUsersPath, req.UserID)
	r, err := c.newRequest(path, http.MethodDelete, nil)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fullBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("response status code %d", resp.StatusCode)
		}
		return fmt.Errorf("response status code %d: %s", resp.StatusCode, fullBody)
	}

	return nil
}
