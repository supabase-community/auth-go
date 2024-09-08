package integration_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/supabase-community/auth-go/types"
)

func TestRecover(t *testing.T) {
	assert := assert.New(t)

	email := randomEmail()
	err := client.Recover(types.RecoverRequest{
		Email: email,
	})
	assert.NoError(err)

	err = client.Recover(types.RecoverRequest{})
	assert.Error(err)
}
