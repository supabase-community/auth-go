package integration_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/supabase-community/auth-go/types"
)

func TestToken(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	client := autoconfirmClient

	// Test login with email
	email := randomEmail()
	password := "password"

	_, err := client.Signup(types.SignupRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(err)

	token, err := client.Token(types.TokenRequest{
		GrantType: "password",
		Email:     email,
		Password:  password,
	})
	require.NoError(err)
	assert.Equal(email, token.User.Email)
	assert.NotEmpty(token.AccessToken)
	assert.NotEmpty(token.RefreshToken)
	assert.Equal("bearer", token.TokenType)
	assert.Equal(3600, token.ExpiresIn)
	assert.InDelta(time.Now().Add(3600*time.Second).Unix(), token.ExpiresAt, float64(time.Second))

	// Signin with email convenience method
	token, err = client.SignInWithEmailPassword(email, password)
	require.NoError(err)
	assert.Equal(email, token.User.Email)
	assert.NotEmpty(token.AccessToken)
	assert.NotEmpty(token.RefreshToken)
	assert.Equal("bearer", token.TokenType)
	assert.Equal(3600, token.ExpiresIn)
	assert.InDelta(time.Now().Add(3600*time.Second).Unix(), token.ExpiresAt, float64(time.Second))

	// Test login with phone
	phone := randomPhoneNumber()
	password = "password"

	_, err = client.Signup(types.SignupRequest{
		Phone:    phone,
		Password: password,
	})
	require.NoError(err)

	token, err = client.Token(types.TokenRequest{
		GrantType: "password",
		Phone:     phone,
		Password:  password,
	})
	require.NoError(err)
	assert.Equal(phone, token.User.Phone)
	assert.NotEmpty(token.AccessToken)
	assert.NotEmpty(token.RefreshToken)
	assert.Equal("bearer", token.TokenType)
	assert.Equal(3600, token.ExpiresIn)
	assert.InDelta(time.Now().Add(3600*time.Second).Unix(), token.ExpiresAt, float64(time.Second))

	// Signin with phone convenience method
	token, err = client.SignInWithPhonePassword(phone, password)
	require.NoError(err)
	assert.Equal(phone, token.User.Phone)
	assert.NotEmpty(token.AccessToken)
	assert.NotEmpty(token.RefreshToken)
	assert.Equal("bearer", token.TokenType)
	assert.Equal(3600, token.ExpiresIn)
	assert.InDelta(time.Now().Add(3600*time.Second).Unix(), token.ExpiresAt, float64(time.Second))

	// Incorrect password
	_, err = client.Token(types.TokenRequest{
		GrantType: "password",
		Email:     email,
		Password:  "wrong",
	})
	assert.Error(err)

	// Test login with refresh token
	email = randomEmail()
	user, err := client.Signup(types.SignupRequest{
		Email:    email,
		Password: "password",
	})
	require.NoError(err)
	require.NotEmpty(user.RefreshToken)

	token, err = client.Token(types.TokenRequest{
		GrantType:    "refresh_token",
		RefreshToken: user.RefreshToken,
	})
	require.NoError(err)
	assert.Equal(email, token.User.Email)
	assert.NotEmpty(token.AccessToken)
	assert.NotEmpty(token.RefreshToken)
	assert.Equal("bearer", token.TokenType)
	assert.Equal(3600, token.ExpiresIn)
	assert.InDelta(time.Now().Add(3600*time.Second).Unix(), token.ExpiresAt, float64(time.Second))

	// Refresh token convenience method
	token, err = client.RefreshToken(token.RefreshToken)
	require.NoError(err)
	assert.Equal(email, token.User.Email)
	assert.NotEmpty(token.AccessToken)
	assert.NotEmpty(token.RefreshToken)
	assert.Equal("bearer", token.TokenType)
	assert.Equal(3600, token.ExpiresIn)
	assert.InDelta(time.Now().Add(3600*time.Second).Unix(), token.ExpiresAt, float64(time.Second))

	// Invalid input tests
	tests := map[string]types.TokenRequest{
		"invalid_grant_type": {
			GrantType: "invalid",
		},
		"password/missing_email_phone": {
			GrantType: "password",
		},
		"password/missing_password": {
			GrantType: "password",
			Email:     email,
		},
		"password/refresh_token_provided": {
			GrantType:    "password",
			Email:        email,
			Password:     password,
			RefreshToken: "refresh_token",
		},
		"refresh_token/missing_refresh_token": {
			GrantType: "refresh_token",
		},
		"refresh_token/email_provided": {
			GrantType: "refresh_token",
			Email:     email,
		},
		"refresh_token/phone_provided": {
			GrantType: "refresh_token",
			Phone:     phone,
		},
		"refresh_token/password_provided": {
			GrantType: "refresh_token",
			Password:  password,
		},
		"pkce/missing_code": {
			GrantType: "pkce",
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := client.Token(test)
			require.Error(err)
		})
	}
}
