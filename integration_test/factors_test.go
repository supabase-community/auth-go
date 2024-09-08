package integration_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/supabase-community/auth-go/types"
)

func TestFactors(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	client := autoconfirmClient

	// Invalid request
	_, err := client.EnrollFactor(types.EnrollFactorRequest{})
	assert.Error(err)

	// No user token
	_, err = client.EnrollFactor(types.EnrollFactorRequest{
		FriendlyName: "test",
		FactorType:   types.FactorTypeTOTP,
		Issuer:       "test.com",
	})
	assert.Error(err)

	// Create a new user
	email := randomEmail()
	session, err := client.Signup(types.SignupRequest{
		Email:    email,
		Password: "password",
	})
	require.NoError(err)

	client = client.WithToken(session.AccessToken)

	// Enroll factor for client
	factorResp, err := client.EnrollFactor(types.EnrollFactorRequest{
		FriendlyName: "Test Factor",
		FactorType:   types.FactorTypeTOTP,
		Issuer:       "mine.com",
	})
	require.NoError(err)
	assert.Equal(types.FactorTypeTOTP, factorResp.Type)
	assert.NotEqual(uuid.Nil, factorResp.ID)
	assert.NotEmpty(factorResp.TOTP.Secret)
	assert.NotEmpty(factorResp.TOTP.QRCode)
	assert.NotEmpty(factorResp.TOTP.URI)

	// Create challenge with invalid request
	_, err = client.ChallengeFactor(types.ChallengeFactorRequest{})
	assert.Error(err)

	// Create challenge with invalid factor ID
	_, err = client.ChallengeFactor(types.ChallengeFactorRequest{
		FactorID: uuid.Nil,
	})
	assert.Error(err)

	// Create valid challenge
	challengeResp, err := client.ChallengeFactor(types.ChallengeFactorRequest{
		FactorID: factorResp.ID,
	})
	require.NoError(err)
	assert.NotEqual(uuid.Nil, challengeResp.ID)
	assert.Greater(challengeResp.ExpiresAt, time.Now())

	// Verify challenge with invalid request
	_, err = client.VerifyFactor(types.VerifyFactorRequest{})
	assert.Error(err)

	// Verify challenge with invalid factor ID
	_, err = client.VerifyFactor(types.VerifyFactorRequest{
		FactorID:    uuid.Nil,
		ChallengeID: challengeResp.ID,
		Code:        "blah",
	})
	assert.Error(err)

	// Verify challenge with invalid challenge ID
	_, err = client.VerifyFactor(types.VerifyFactorRequest{
		FactorID:    factorResp.ID,
		ChallengeID: uuid.Nil,
		Code:        "blah",
	})
	assert.Error(err)

	// Verify challenge with invalid code
	_, err = client.VerifyFactor(types.VerifyFactorRequest{
		FactorID:    factorResp.ID,
		ChallengeID: challengeResp.ID,
		Code:        "blah",
	})
	assert.Error(err)

	// Cannot test verify with a valid code without actually enrolling a TOTP
	// factor.

	// Delete factor with invalid request
	_, err = client.UnenrollFactor(types.UnenrollFactorRequest{})
	assert.Error(err)

	// Delete factor
	unenrollResp, err := client.UnenrollFactor(types.UnenrollFactorRequest{
		FactorID: factorResp.ID,
	})
	require.NoError(err)
	assert.Equal(factorResp.ID, unenrollResp.ID)
}
