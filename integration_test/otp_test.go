package integration_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/supabase-community/auth-go/types"
)

func TestOTP(t *testing.T) {
	assert := assert.New(t)

	// Cannot create user from OTP if CreateUser is false
	email := randomEmail()
	err := client.OTP(types.OTPRequest{
		Email:      email,
		CreateUser: false,
	})
	assert.Error(err)

	// Create user from OTP
	err = client.OTP(types.OTPRequest{
		Email:      email,
		CreateUser: true,
	})
	assert.NoError(err)

	// Create user with SMS OTP, but SMS disabled
	phone := randomPhoneNumber()
	err = client.OTP(types.OTPRequest{
		Phone:      phone,
		CreateUser: true,
	})
	assert.Error(err)
}
