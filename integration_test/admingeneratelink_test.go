package integration_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/supabase-community/auth-go/types"
)

func TestAdminGenerateLink(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	admin := withAdmin(client)

	// Testing signup link
	email := randomEmail()
	resp, err := admin.AdminGenerateLink(types.AdminGenerateLinkRequest{
		Type:       "signup",
		Email:      email,
		Password:   "password",
		Data:       map[string]interface{}{},
		RedirectTo: "http://localhost:3000",
	})
	require.NoError(err)

	assert.EqualValues(resp.VerificationType, "signup")
	linkRegexp := regexp.MustCompile(`^http://localhost:9999/verify\?token=[a-zA-Z0-9_-]+&type=signup&redirect_to=http://localhost:3000$`)
	assert.Regexp(linkRegexp, resp.ActionLink)
	assert.NotEmpty(resp.HashedToken)
	assert.NotEmpty(resp.EmailOTP)
	assert.Equal("http://localhost:3000", resp.RedirectTo)
	assert.Contains(resp.ActionLink, resp.HashedToken)

	assert.NotEqual(uuid.Nil, resp.ID)
	assert.Equal(resp.Email, email)
	assert.InDelta(time.Now().Unix(), resp.ConfirmationSentAt.Unix(), float64(time.Second))
	assert.InDelta(time.Now().Unix(), resp.CreatedAt.Unix(), float64(time.Second))
	assert.InDelta(time.Now().Unix(), resp.UpdatedAt.Unix(), float64(time.Second))

	// Testing invalid requests
	tests := map[string]types.AdminGenerateLinkRequest{
		"signup/missing_email": {
			Type:     "signup",
			Password: "password",
		},
		"signup/missing_password": {
			Type:  "signup",
			Email: email,
		},
		"magiclink/missing_email": {
			Type: "magiclink",
		},
		"magiclink/password_provided": {
			Type:     "magiclink",
			Email:    email,
			Password: "password",
		},
		"invite/missing_email": {
			Type: "invite",
		},
		"invite/password_provided": {
			Type:     "invite",
			Email:    email,
			Password: "password",
		},
		"recovery/missing_email": {
			Type: "recovery",
		},
		"recovery/data_provided": {
			Type:  "recovery",
			Email: email,
			Data: map[string]interface{}{
				"foo": "bar",
			},
		},
		"recovery/password_provided": {
			Type:     "recovery",
			Email:    email,
			Password: "password",
		},
		"email_change_current/missing_email": {
			Type:     "email_change_current",
			NewEmail: email,
		},
		"email_change_current/missing_new_email": {
			Type:  "email_change_current",
			Email: email,
		},
		"email_change_current/data_provided": {
			Type:     "email_change_current",
			Email:    email,
			NewEmail: email,
			Data: map[string]interface{}{
				"foo": "bar",
			},
		},
		"email_change_current/password_provided": {
			Type:     "email_change_current",
			Email:    email,
			NewEmail: email,
			Password: "password",
		},
		"email_change_new/missing_email": {
			Type:     "email_change_new",
			NewEmail: email,
		},
		"email_change_new/missing new_email": {
			Type:  "email_change_new",
			Email: email,
		},
		"email_change_new/data_provided": {
			Type:     "email_change_new",
			Email:    email,
			NewEmail: email,
			Data: map[string]interface{}{
				"foo": "bar",
			},
		},
		"email_change_new/password_provided": {
			Type:     "email_change_new",
			Email:    email,
			NewEmail: email,
			Password: "password",
		},
	}
	for name, data := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := admin.AdminGenerateLink(data)
			assert.Error(err)
			assert.ErrorContains(err, "request is invalid")
		})
	}

	// Testing email change links
	newEmail := randomEmail()
	resp, err = admin.AdminGenerateLink(types.AdminGenerateLinkRequest{
		Type:       "email_change_current",
		Email:      email,
		NewEmail:   newEmail,
		RedirectTo: "http://localhost:3000",
	})
	require.NoError(err)
	assert.Equal("http://localhost:3000", resp.RedirectTo)
	assert.EqualValues(resp.VerificationType, "email_change_current")
	linkRegexp = regexp.MustCompile(`^http://localhost:9999/verify\?token=[a-zA-Z0-9_-]+&type=email_change&redirect_to=http://localhost:3000$`)
	assert.Regexp(linkRegexp, resp.ActionLink)
	assert.NotEmpty(resp.HashedToken)
	assert.NotEmpty(resp.EmailOTP)
	assert.Contains(resp.ActionLink, resp.HashedToken)

	assert.NotEqual(uuid.Nil, resp.ID)
	assert.Equal(resp.Email, email)
	assert.Equal(resp.EmailChange, newEmail)
	assert.InDelta(time.Now().Unix(), resp.ConfirmationSentAt.Unix(), float64(time.Second))
	assert.InDelta(time.Now().Unix(), resp.CreatedAt.Unix(), float64(time.Second))
	assert.InDelta(time.Now().Unix(), resp.UpdatedAt.Unix(), float64(time.Second))
}
