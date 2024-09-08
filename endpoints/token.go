package endpoints

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/supabase-community/auth-go/types"
)

const tokenPath = "/token"

// Sign in with email and password
//
// This is a convenience method that calls Token with the password grant type
func (c *Client) SignInWithEmailPassword(email, password string) (*types.TokenResponse, error) {
	return c.Token(types.TokenRequest{
		GrantType: "password",
		Email:     email,
		Password:  password,
	})
}

// Sign in with phone and password
//
// This is a convenience method that calls Token with the password grant type
func (c *Client) SignInWithPhonePassword(phone, password string) (*types.TokenResponse, error) {
	return c.Token(types.TokenRequest{
		GrantType: "password",
		Phone:     phone,
		Password:  password,
	})
}

// Sign in with refresh token
//
// This is a convenience method that calls Token with the refresh_token grant type
func (c *Client) RefreshToken(refreshToken string) (*types.TokenResponse, error) {
	return c.Token(types.TokenRequest{
		GrantType:    "refresh_token",
		RefreshToken: refreshToken,
	})
}

// POST /token
//
// This is an OAuth2 endpoint that currently implements the password,
// refresh_token, and PKCE grant types
func (c *Client) Token(req types.TokenRequest) (*types.TokenResponse, error) {
	switch req.GrantType {
	case "password":
		if (req.Email == "" && req.Phone == "") || req.Password == "" || req.RefreshToken != "" {
			return nil, types.ErrInvalidTokenRequest
		}
	case "refresh_token":
		if req.RefreshToken == "" || req.Email != "" || req.Phone != "" || req.Password != "" {
			return nil, types.ErrInvalidTokenRequest
		}
	case "pkce":
		if req.Code == "" || req.CodeVerifier == "" {
			return nil, types.ErrInvalidTokenRequest
		}
	default:
		return nil, types.ErrInvalidTokenRequest
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	r, err := c.newRequest(tokenPath+"?grant_type="+req.GrantType, http.MethodPost, bytes.NewBuffer(body))
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

	var res types.TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}
