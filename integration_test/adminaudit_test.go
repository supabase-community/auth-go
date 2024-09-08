package integration_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/supabase-community/auth-go/types"
)

// This test creates some audit logs, then retrieves them with pagination.
// However, other tests running in parallel may also create logs, so we only
// check that _at least_ the logs we create here are returned.
//
// In order to test this functionality exactly, we need to run this test against
// a fresh database with an empty audit log table.
func TestAdminAudit(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	// User must be authenticated first
	_, err := client.AdminAudit(types.AdminAuditRequest{})
	assert.Error(err)

	client := client.WithToken(adminToken())

	// Get audit logs
	resp, err := client.AdminAudit(types.AdminAuditRequest{})
	require.NoError(err)
	assert.NotNil(resp.Logs)

	// Generate some log entries
	for i := 0; i < 10; i++ {
		pass := randomString(10)
		_, err = client.AdminCreateUser(types.AdminCreateUserRequest{
			Email:    randomEmail(),
			Password: &pass,
		})
		require.NoError(err)
	}

	// Get audit logs with pagination
	resp, err = client.AdminAudit(types.AdminAuditRequest{
		Page:    1,
		PerPage: 4,
	})
	require.NoError(err)
	assert.Len(resp.Logs, 4)
	assert.GreaterOrEqual(resp.TotalCount, 10)
	assert.EqualValues(2, resp.NextPage)
	assert.GreaterOrEqual(resp.TotalPages, uint(3))

	resp, err = client.AdminAudit(types.AdminAuditRequest{
		Page:    2,
		PerPage: 4,
	})
	require.NoError(err)
	assert.Len(resp.Logs, 4)
	assert.GreaterOrEqual(resp.TotalCount, 10)
	assert.EqualValues(3, resp.NextPage)
	assert.GreaterOrEqual(resp.TotalPages, uint(3))

	resp, err = client.AdminAudit(types.AdminAuditRequest{
		Page:    3,
		PerPage: 4,
	})
	require.NoError(err)
	assert.GreaterOrEqual(len(resp.Logs), 2)
	assert.GreaterOrEqual(resp.TotalCount, 10)
	assert.GreaterOrEqual(resp.TotalPages, uint(3))

	// Invalid - empty query
	_, err = client.AdminAudit(types.AdminAuditRequest{
		Query: &types.AuditQuery{},
	})
	assert.Error(err)

	// Invalid - invalid query column
	_, err = client.AdminAudit(types.AdminAuditRequest{
		Query: &types.AuditQuery{
			Column: "invalid",
			Value:  "valid",
		},
	})
	assert.Error(err)

	// Invalid - no query value
	_, err = client.AdminAudit(types.AdminAuditRequest{
		Query: &types.AuditQuery{
			Column: types.AuditQueryColumnAction,
		},
	})
	assert.Error(err)

	// Valid query
	resp, err = client.AdminAudit(types.AdminAuditRequest{
		Query: &types.AuditQuery{
			Column: types.AuditQueryColumnAction,
			Value:  "user_signedup",
		},
	})
	require.NoError(err)
	assert.GreaterOrEqual(len(resp.Logs), 10)

	resp, err = client.AdminAudit(types.AdminAuditRequest{
		Query: &types.AuditQuery{
			Column: types.AuditQueryColumnAction,
			Value:  "not user_signedup",
		},
	})
	require.NoError(err)
	assert.Len(resp.Logs, 0)
}
