package endpoints

import (
	"io"
	"net/http"
	"net/url"
)

const samlMetadataPath = "/sso/saml/metadata"
const samlACSPath = "/sso/saml/acs"

// GET /sso/saml/metadata
//
// Get the SAML metadata for the configured SAML provider.
//
// If successful, the server returns an XML response. Making sense of this is
// outside the scope of this client, so it is simply returned as []byte.
func (c *Client) SAMLMetadata() ([]byte, error) {
	r, err := c.newRequest(samlMetadataPath, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, handleErrorResponse(resp)
	}

	return io.ReadAll(resp.Body)
}

// POST /sso/saml/acs
//
// Implements the main Assertion Consumer Service endpoint behavior.
//
// This client does not provide a typed endpoint for SAML ACS. This method is
// provided for convenience and will simply POST your HTTP request to the
// endpoint and return the response.
//
// For required parameters, see the SAML spec or the Auth implementation
// of this endpoint.
//
// The server may issue redirects. Using the default HTTP client, this method
// will follow those redirects and return the final HTTP response. Should you
// prefer the client not to follow redirects, you can provide a custom HTTP
// client using WithClient(). See the example below.
//
// Example:
//
//	c := http.Client{
//		CheckRedirect: func(req *http.Request, via []*http.Request) error {
//			return http.ErrUseLastResponse
//		},
//	}
func (c *Client) SAMLACS(req *http.Request) (*http.Response, error) {
	acsURL := c.baseURL + samlACSPath
	u, err := url.Parse(acsURL)
	if err != nil {
		return nil, err
	}
	req.URL = u
	return c.client.Do(req)
}
