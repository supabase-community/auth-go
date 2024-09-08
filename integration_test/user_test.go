package integration_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/supabase-community/auth-go/types"
)

func TestUser(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	client := autoconfirmClient

	// Create user
	email := randomEmail()
	password := "password"
	session, err := client.Signup(types.SignupRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(err)

	// Get user
	user, err := client.WithToken(session.AccessToken).GetUser()
	require.NoError(err)
	assert.NotEqual(user.ID, uuid.Nil)
	assert.Equal(email, user.Email)
	assert.InDelta(time.Now().Unix(), user.CreatedAt.Unix(), float64(time.Second))
	assert.InDelta(time.Now().Unix(), user.UpdatedAt.Unix(), float64(time.Second))

	// Modify user metadata
	updateResp, err := client.WithToken(session.AccessToken).UpdateUser(types.UpdateUserRequest{
		Data: map[string]interface{}{
			"foo": "bar",
		},
	})
	require.NoError(err)
	assert.Equal("bar", updateResp.UserMetadata["foo"])
}
