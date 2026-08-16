package endpoints

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/supabase-community/auth-go/types"
)

const adminCustomProvidersPath = "/admin/custom-providers"

// customProviderPath builds the path for a single provider. 
func customProviderPath(identifier string) string {
	// TODO: check if identifier contains a `:` and fail if not? 
	return fmt.Sprintf("%s/%s", adminCustomProvidersPath, url.PathEscape(identifier))
}

// GET /admin/custom-providers
//
// Get a list of all custom OAuth2/OIDC providers in the system.
//
// Set ProviderType on the request to filter by "oauth2" or "oidc".
func (c *Client) AdminListCustomOAuthProviders(req types.AdminListCustomOAuthProvidersRequest) (*types.AdminListCustomOAuthProvidersResponse, error) {
	path := adminCustomProvidersPath
	if req.ProviderType != "" {
		query := url.Values{}
		query.Set("type", string(req.ProviderType))
		path = fmt.Sprintf("%s?%s", path, query.Encode())
	}

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

	var res types.AdminListCustomOAuthProvidersResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// POST /admin/custom-providers
//
// Create a new custom OAuth2 or OIDC provider.
func (c *Client) AdminCreateCustomOAuthProvider(req types.AdminCreateCustomOAuthProviderRequest) (*types.AdminCreateCustomOAuthProviderResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	r, err := c.newRequest(adminCustomProvidersPath, http.MethodPost, bytes.NewBuffer(body))
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

	var res types.AdminCreateCustomOAuthProviderResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// GET /admin/custom-providers/{identifier}
//
// Get a custom OAuth2/OIDC provider by identifier.
func (c *Client) AdminGetCustomOAuthProvider(req types.AdminGetCustomOAuthProviderRequest) (*types.AdminGetCustomOAuthProviderResponse, error) {
	r, err := c.newRequest(customProviderPath(req.Identifier), http.MethodGet, nil)
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

	var res types.AdminGetCustomOAuthProviderResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// PUT /admin/custom-providers/{identifier}
//
// Update a custom OAuth2/OIDC provider by identifier. Omitted fields are left
// unchanged.
func (c *Client) AdminUpdateCustomOAuthProvider(req types.AdminUpdateCustomOAuthProviderRequest) (*types.AdminUpdateCustomOAuthProviderResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	r, err := c.newRequest(customProviderPath(req.Identifier), http.MethodPut, bytes.NewBuffer(body))
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

	var res types.AdminUpdateCustomOAuthProviderResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// DELETE /admin/custom-providers/{identifier}
//
// Delete a custom OAuth2/OIDC provider by identifier. The server returns no
// body on success.
func (c *Client) AdminDeleteCustomOAuthProvider(req types.AdminDeleteCustomOAuthProviderRequest) error {
	r, err := c.newRequest(customProviderPath(req.Identifier), http.MethodDelete, nil)
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
