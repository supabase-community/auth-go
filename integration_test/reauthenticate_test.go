package integration_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/supabase-community/auth-go/types"
)

func TestReauthenticate(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	// User must be authenticated first
	err := client.Reauthenticate()
	assert.Error(err)

	client := autoconfirmClient

	// Create a new user
	email := randomEmail()
	session, err := client.Signup(types.SignupRequest{
		Email:    email,
		Password: "password",
	})
	require.NoError(err)

	client = client.WithToken(session.AccessToken)
	err = client.Reauthenticate()
	assert.NoError(err)
}
