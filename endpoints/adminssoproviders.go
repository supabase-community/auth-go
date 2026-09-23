package endpoints

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/supabase-community/auth-go/types"
)

const adminSSOPath = "/admin/sso/providers"

// GET /admin/sso/providers
//
// Get a list of all SAML SSO Identity Providers in the system.
func (c *Client) AdminListSSOProviders() (*types.AdminListSSOProvidersResponse, error) {
	r, err := c.newRequest(adminSSOPath, http.MethodGet, nil)
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

	var res types.AdminListSSOProvidersResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// POST /admin/sso/providers
//
// Create a new SAML SSO Identity Provider.
func (c *Client) AdminCreateSSOProvider(req types.AdminCreateSSOProviderRequest) (*types.AdminCreateSSOProviderResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	r, err := c.newRequest(adminSSOPath, http.MethodPost, bytes.NewBuffer(body))
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

	var res types.AdminCreateSSOProviderResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// GET /admin/sso/providers/{idp_id}
//
// Get a SAML SSO Identity Provider by ID.
func (c *Client) AdminGetSSOProvider(req types.AdminGetSSOProviderRequest) (*types.AdminGetSSOProviderResponse, error) {
	r, err := c.newRequest(fmt.Sprintf("%s/%s", adminSSOPath, req.ProviderID), http.MethodGet, nil)
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

	var res types.AdminGetSSOProviderResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// PUT /admin/sso/providers/{idp_id}
//
// Update a SAML SSO Identity Provider by ID.
func (c *Client) AdminUpdateSSOProvider(req types.AdminUpdateSSOProviderRequest) (*types.AdminUpdateSSOProviderResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	r, err := c.newRequest(fmt.Sprintf("%s/%s", adminSSOPath, req.ProviderID), http.MethodPut, bytes.NewBuffer(body))
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

	var res types.AdminUpdateSSOProviderResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// DELETE /admin/sso/providers/{idp_id}
//
// Delete a SAML SSO Identity Provider by ID.
func (c *Client) AdminDeleteSSOProvider(req types.AdminDeleteSSOProviderRequest) (*types.AdminDeleteSSOProviderResponse, error) {
	path := fmt.Sprintf("%s/%s", adminSSOPath, req.ProviderID)
	r, err := c.newRequest(path, http.MethodDelete, nil)
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

	var res types.AdminDeleteSSOProviderResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
