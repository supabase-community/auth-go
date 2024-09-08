package integration_test

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"testing"

	backoff "github.com/cenkalti/backoff/v4"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/supabase-community/auth-go"
)

const (
	projectReference = "project_ref"
	apiKey           = "api_key"
	jwtSecret        = "secret"
)

var (
	// Global clients are used for all tests in this package.
	client               auth.Client
	autoconfirmClient    auth.Client
	signupDisabledClient auth.Client

	// Used to validate UUIDs.
	uuidRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{12}$`)
)

func randomString(n int) string {
	// Using all lower case because email addresses are lowercased by Auth.
	letterBytes := "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func randomEmail() string {
	return fmt.Sprintf("%s@test.com", randomString(10))
}

func randomNumberString(n int) string {
	numberBytes := "0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = numberBytes[rand.Intn(len(numberBytes))]
	}
	return string(b)
}

func randomPhoneNumber() string {
	return fmt.Sprintf("1%s", randomNumberString(10))
}

func adminToken() string {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"aud":  "admin",
		"sub":  "admin",
		"role": "supabase_admin",
		"exp":  9999999999,
	})
	token, err := t.SignedString([]byte(jwtSecret))
	if err != nil {
		panic(err)
	}
	return token
}

func withAdmin(c auth.Client) auth.Client {
	return c.WithToken(adminToken())
}

func TestMain(m *testing.M) {
	// Please refer to ./setup/docker-compose.yaml and ./README.md for more info
	// on this test set up.
	client = auth.New(projectReference, apiKey).WithCustomAuthURL("http://localhost:9999")
	autoconfirmClient = auth.New(projectReference, apiKey).WithCustomAuthURL("http://localhost:9998")
	signupDisabledClient = auth.New(projectReference, apiKey).WithCustomAuthURL("http://localhost:9997")

	// Ensure the server is ready before running tests.
	err := backoff.Retry(
		func() error {
			health, err := client.HealthCheck()
			if err != nil {
				return err
			}
			if health.Name != "GoTrue" {
				return fmt.Errorf("health check - unexpected server name: %s", health.Name)
			}

			health, err = autoconfirmClient.HealthCheck()
			if err != nil {
				return err
			}
			if health.Name != "GoTrue" {
				return fmt.Errorf("health check - unexpected server name: %s", health.Name)
			}

			health, err = signupDisabledClient.HealthCheck()
			if err != nil {
				return err
			}
			if health.Name != "GoTrue" {
				return fmt.Errorf("health check - unexpected server name: %s", health.Name)
			}
			return nil
		},
		backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 10),
	)
	if err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	code := m.Run()
	os.Exit(code)
}

func TestWithClient(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)
	require := require.New(t)

	c := auth.New(projectReference, apiKey).WithCustomAuthURL("http://localhost:9999")
	h, err := c.HealthCheck()
	require.NoError(err)
	assert.Equal("GoTrue", h.Name)

	roundTripper := &customRoundTripper{}
	c = c.WithClient(http.Client{
		Transport: roundTripper,
	})
	h, err = c.HealthCheck()
	require.NoError(err)
	assert.Equal("GoTrue", h.Name)
	assert.True(roundTripper.visited)
}

type customRoundTripper struct {
	visited bool
}

func (c *customRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	c.visited = true
	subC := http.Client{}
	return subC.Do(req)
}
