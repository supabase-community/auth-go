package integration_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/supabase-community/auth-go/types"
)

func TestInvite(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	admin := withAdmin(client)

	email := randomEmail()
	user, err := admin.Invite(types.InviteRequest{
		Email: email,
	})
	require.NoError(err)
	assert.NotEqual(user.ID, uuid.Nil)
	assert.Equal(email, user.Email)
	assert.InDelta(time.Now().Unix(), user.CreatedAt.Unix(), float64(time.Second))
	assert.InDelta(time.Now().Unix(), user.UpdatedAt.Unix(), float64(time.Second))
	assert.InDelta(time.Now().Unix(), user.InvitedAt.Unix(), float64(time.Second))
}
