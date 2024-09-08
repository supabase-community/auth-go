package endpoints

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/supabase-community/auth-go/types"
)

const ssoPath = "/sso"

// POST /sso
//
// Initiate an SSO session with the given provider.
//
// If successful, the server returns a redirect to the provider's authorization
// URL. The client will follow it and return the final HTTP response.
//
// Auth allows you to skip following the redirect by setting SkipHTTPRedirect
// on the request struct. In this case, the URL to redirect to will be returned
// in the response.
func (c *Client) SSO(req types.SSORequest) (*types.SSOResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	r, err := c.newRequest(http.MethodPost, ssoPath, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(r)
	if err != nil {
		return nil, err
	}

	if !req.SkipHTTPRedirect {
		// If the client is following redirects, we can return the response
		// directly.
		return &types.SSOResponse{
			HTTPResponse: resp,
		}, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusSeeOther {
		fullBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("response status code %d", resp.StatusCode)
		}
		return nil, fmt.Errorf("response status code %d: %s", resp.StatusCode, fullBody)
	}

	// If the client is not following redirects, we can unmarshal the response from
	// the server to get the URL.
	var res types.SSOResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
