package integration_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/supabase-community/auth-go/types"
)

func TestSignup(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	// Signup with email
	email := randomEmail()
	userResp, err := client.Signup(types.SignupRequest{
		Email:    email,
		Password: "password",
	})
	require.NoError(err)
	assert.NotEqual(uuid.Nil, userResp.ID)
	assert.Equal(userResp.Email, email)
	assert.InDelta(time.Now().Unix(), userResp.ConfirmationSentAt.Unix(), float64(time.Second))
	assert.InDelta(time.Now().Unix(), userResp.CreatedAt.Unix(), float64(time.Second))
	assert.InDelta(time.Now().Unix(), userResp.UpdatedAt.Unix(), float64(time.Second))

	// Duplicate signup
	dupeUserResp, err := client.Signup(types.SignupRequest{
		Email:    email,
		Password: "password",
	})
	require.NoError(err)
	assert.NotEqual(uuid.Nil, dupeUserResp.ID)
	assert.Equal(dupeUserResp.Email, email)
	assert.InDelta(time.Now().Unix(), dupeUserResp.ConfirmationSentAt.Unix(), float64(time.Second))
	assert.InDelta(time.Now().Unix(), dupeUserResp.CreatedAt.Unix(), float64(time.Second))
	assert.InDelta(time.Now().Unix(), dupeUserResp.UpdatedAt.Unix(), float64(time.Second))
	assert.Equal(userResp.ID, dupeUserResp.ID)

	// Sign up with phone
	// Will error because SMS is not configured on the test server.
	user, err := client.Signup(types.SignupRequest{
		Phone:    "+15555555555",
		Password: "password",
	})
	assert.Error(err)
	assert.Nil(user)

	// Autoconfirmed signup
	// Should return a session
	email = randomEmail()
	session, err := autoconfirmClient.Signup(types.SignupRequest{
		Email:    email,
		Password: "password",
	})
	require.NoError(err)
	assert.NotEqual(uuid.Nil, session.ID)
	assert.Equal(session.Email, email)
	assert.InDelta(time.Now().Unix(), session.CreatedAt.Unix(), float64(time.Second))
	assert.InDelta(time.Now().Unix(), session.UpdatedAt.Unix(), float64(time.Second))
	assert.NotEmpty(session.AccessToken)
	assert.NotEmpty(session.RefreshToken)
	assert.Equal("bearer", session.TokenType)
	assert.Equal(3600, session.ExpiresIn)
	assert.InDelta(time.Now().Add(3600*time.Second).Unix(), session.ExpiresAt, float64(time.Second))

	// Sign up with signups disabled
	email = randomEmail()
	user, err = signupDisabledClient.Signup(types.SignupRequest{
		Email:    email,
		Password: "password",
	})
	assert.Error(err)
	assert.Nil(user)
}
