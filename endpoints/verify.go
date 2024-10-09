package endpoints

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/supabase-community/auth-go/types"
)

const verifyPath = "/verify"

// GET /verify
//
// Verify a registration or a password recovery. Type can be signup or recovery
// or magiclink or invite and the token is a token returned from either /signup
// or /recover or /magiclink.
//
// The server returns a redirect response. This method will not follow the
// redirect, but instead returns the URL the client was told to redirect to,
// as well as parsing the parameters from the URL fragment.
//
// NOTE: This endpoint may return a nil error, but the Response can contain
// error details extracted from the returned URL. Please check that the Error,
// ErrorCode and/or ErrorDescription fields of the response are empty.
func (c *Client) Verify(req types.VerifyRequest) (*types.VerifyResponse, error) {
	if req.Type == "" {
		return nil, types.ErrInvalidVerifyRequest
	}
	if req.Token == "" {
		return nil, types.ErrInvalidVerifyRequest
	}
	if req.RedirectTo == "" {
		return nil, types.ErrInvalidVerifyRequest
	}

	r, err := c.newRequest(verifyPath, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	q := r.URL.Query()
	q.Add("type", string(req.Type))
	q.Add("token", req.Token)
	q.Add("redirect_to", req.RedirectTo)
	r.URL.RawQuery = q.Encode()

	// Set up a client that will not follow the redirect.
	noRedirClient := noRedirClient(c.client)

	resp, err := noRedirClient.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusSeeOther {
		fullBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("response status code %d", resp.StatusCode)
		}
		return nil, fmt.Errorf("response status code %d: %s", resp.StatusCode, fullBody)
	}

	redirURL := resp.Header.Get("Location")
	if redirURL == "" {
		return nil, fmt.Errorf("no redirect URL found in response")
	}
	u, err := url.Parse(redirURL)
	if err != nil {
		return nil, err
	}
	values, err := url.ParseQuery(u.Fragment)
	if err != nil {
		return nil, err
	}
	expiry := values.Get("expires_in")
	expiresIn := 0
	if expiry != "" {
		// e should remain 0 if this fails.
		expiresIn, _ = strconv.Atoi(expiry)
	}

	return &types.VerifyResponse{
		URL: redirURL,

		AccessToken:  values.Get("access_token"),
		TokenType:    values.Get("token_type"),
		ExpiresIn:    expiresIn,
		RefreshToken: values.Get("refresh_token"),
		Type:         (types.VerificationType)(values.Get("type")),

		Error:            values.Get("error"),
		ErrorCode:        values.Get("error_code"),
		ErrorDescription: values.Get("error_description"),
	}, nil
}

// POST /verify
//
// Verify a registration or a password recovery. Type can be signup or recovery
// or magiclink or invite and the token is a token returned from either /signup
// or /recover or /magiclink.
//
// This differs from GET /verify as it requires an email or phone to be given,
// which is used to verify the token associated to the user. It also returns a
// JSON response rather than a redirect.
func (c *Client) VerifyForUser(req types.VerifyForUserRequest) (*types.VerifyForUserResponse, error) {
	if req.Type == "" {
		return nil, types.ErrInvalidVerifyRequest
	}
	if req.Token == "" {
		return nil, types.ErrInvalidVerifyRequest
	}
	if req.RedirectTo == "" {
		return nil, types.ErrInvalidVerifyRequest
	}
	if req.Email == "" && req.Phone == "" {
		return nil, types.ErrInvalidVerifyRequest
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	r, err := c.newRequest(verifyPath, http.MethodPost, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fullBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("response status code %d", resp.StatusCode)
		}
		return nil, fmt.Errorf("response status code %d: %s", resp.StatusCode, fullBody)
	}

	var res types.VerifyForUserResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
