package endpoints

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"

	"github.com/supabase-community/auth-go/types"
)

const authorizePath = "/authorize"

func generatePKCEParams() (*types.PKCEParams, error) {
	data := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, data); err != nil {
		return nil, err
	}

	// RawURLEncoding since "code challenge can only contain alphanumeric characters, hyphens, periods, underscores and tildes"
	verifier := base64.RawURLEncoding.EncodeToString(data)
	sha := sha256.Sum256([]byte(verifier))
	challenge := base64.RawURLEncoding.EncodeToString(sha[:])
	return &types.PKCEParams{
		Challenge:       challenge,
		ChallengeMethod: "S256",
		Verifier:        verifier,
	}, nil
}

// GET /authorize
//
// Get access_token from external oauth provider.
//
// Scopes are optional additional scopes depending on the provider (email and
// name are requested by default).
//
// If successful, the server returns a redirect response. This method will not
// follow the redirect, but instead returns the URL the client was told to
// redirect to.
func (c *Client) Authorize(req types.AuthorizeRequest) (*types.AuthorizeResponse, error) {
	r, err := c.newRequest(authorizePath, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	q := r.URL.Query()
	q.Add("scopes", req.Scopes)
	q.Add("provider", string(req.Provider))
	q.Add("redirect_to", req.RedirectTo)

	verifier := ""

	if string(req.FlowType) == string(types.FlowPKCE) {
		pkce, err := generatePKCEParams()
		if err != nil {
			return nil, err
		}
		q.Add("code_challenge", pkce.Challenge)
		q.Add("code_challenge_method", pkce.ChallengeMethod)
		verifier = pkce.Verifier
	}

	r.URL.RawQuery = q.Encode()

	// Set up a client that will not follow the redirect.
	noRedirClient := noRedirClient(c.client)

	resp, err := noRedirClient.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusFound {
		return nil, handleErrorResponse(resp)
	}

	url := resp.Header.Get("Location")
	if url == "" {
		return nil, fmt.Errorf("no redirect URL found in response")
	}
	return &types.AuthorizeResponse{
		AuthorizationURL: url,
		Verifier:         verifier,
	}, nil
}
