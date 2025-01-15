package integration_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/supabase-community/auth-go/types"
)

func TestAdminCreateUser(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	admin := withAdmin(client)

	pass := "password"
	email := randomEmail()
	req := types.AdminCreateUserRequest{
		Email:    email,
		Role:     "admin",
		Password: &pass,
	}
	resp, err := admin.AdminCreateUser(req)
	require.NoError(err)
	require.Regexp(uuidRegex, resp.ID)
	assert.Equal(resp.Email, email)
	assert.Equal(resp.Role, "admin")

	_, err = admin.AdminCreateUser(types.AdminCreateUserRequest{})
	assert.Error(err)
}

func TestAdminListUsers(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	admin := withAdmin(client)

	// Create a user that we know should be returned
	pass := "password"
	email := randomEmail()
	req := types.AdminCreateUserRequest{
		Email:    email,
		Role:     "test",
		Password: &pass,
	}
	createResp, err := admin.AdminCreateUser(req)
	require.NoError(err)
	require.Regexp(uuidRegex, createResp.ID)

	// Then list and look up the user we just created
	page := 0
	perPage := 1
	resp, err := admin.AdminListUsers(types.AdminListUsersRequest{
		Page:    &page,
		PerPage: &perPage,
	})
	require.NoError(err)
	assert.NotEmpty(resp)
	for _, u := range resp.Users {
		assert.NotEqual(uuid.Nil, u.ID)
		assert.Equal(u.Email, createResp.Email)
	}
}

func TestAdminGetUser(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	// Need admin credential
	_, err := client.AdminGetUser(types.AdminGetUserRequest{
		UserID: uuid.Nil,
	})
	assert.Error(err)

	admin := withAdmin(client)

	// Create a user
	pass := "password"
	email := randomEmail()
	req := types.AdminCreateUserRequest{
		Email:    email,
		Role:     "test",
		Password: &pass,
	}
	createResp, err := admin.AdminCreateUser(req)
	require.NoError(err)
	require.Regexp(uuidRegex, createResp.ID)

	// Get that user
	resp, err := admin.AdminGetUser(types.AdminGetUserRequest{
		UserID: createResp.ID,
	})
	require.NoError(err)
	assert.Equal(resp.Email, email)
	assert.Equal(resp.Role, "test")

	// Cannot get a user that doesn't exist
	_, err = admin.AdminGetUser(types.AdminGetUserRequest{
		UserID: uuid.New(),
	})
	assert.Error(err)
}

func TestAdminUpdateUser(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	// Need admin credential
	_, err := client.AdminUpdateUser(types.AdminUpdateUserRequest{
		UserID: uuid.Nil,
	})
	assert.Error(err)

	admin := withAdmin(client)

	// Create a user
	pass := "password"
	email := randomEmail()
	req := types.AdminCreateUserRequest{
		Email:    email,
		Role:     "test",
		Password: &pass,
	}
	createResp, err := admin.AdminCreateUser(req)
	require.NoError(err)
	require.Regexp(uuidRegex, createResp.ID)

	// Update the user
	resp, err := admin.AdminUpdateUser(types.AdminUpdateUserRequest{
		UserID: createResp.ID,
		Role:   "admin",
	})
	require.NoError(err)
	assert.Equal(resp.Email, email)
	assert.Equal(resp.Role, "admin")

	// Cannot update a user that doesn't exist
	_, err = admin.AdminUpdateUser(types.AdminUpdateUserRequest{
		UserID: uuid.New(),
	})
	assert.Error(err)
}

func TestAdminDeleteUser(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	// Need admin credential
	err := client.AdminDeleteUser(types.AdminDeleteUserRequest{
		UserID: uuid.Nil,
	})
	assert.Error(err)

	admin := withAdmin(client)

	// Create a user
	pass := "password"
	email := randomEmail()
	resp, err := admin.AdminCreateUser(types.AdminCreateUserRequest{
		Email:    email,
		Role:     "test",
		Password: &pass,
	})
	require.NoError(err)
	assert.NotEqual(uuid.Nil, resp.ID)

	// Delete the user
	err = admin.AdminDeleteUser(types.AdminDeleteUserRequest{
		UserID: resp.ID,
	})
	assert.NoError(err)

	// Cannot delete a user that doesn't exist
	err = admin.AdminDeleteUser(types.AdminDeleteUserRequest{
		UserID: uuid.New(),
	})
	assert.Error(err)
}
