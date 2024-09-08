package endpoints

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/supabase-community/auth-go/types"
)

const adminGenerateLinkPath = "/admin/generate_link"

func validateAdminGenerateLinkRequest(req types.AdminGenerateLinkRequest) error {
	switch req.Type {
	case types.LinkTypeSignup:
		if req.Email == "" || req.Password == "" {
			return types.NewErrInvalidGenerateLinkRequest("email and password must be provided if Type is signup")
		}
	case types.LinkTypeMagicLink, types.LinkTypeInvite:
		if req.Email == "" {
			return types.NewErrInvalidGenerateLinkRequest("email must be provided if Type is magiclink or invite")
		}
		if req.Password != "" {
			return types.NewErrInvalidGenerateLinkRequest("password must not be provided if Type is magiclink or invite")
		}
	case types.LinkTypeRecovery:
		if req.Email == "" {
			return types.NewErrInvalidGenerateLinkRequest("email must be provided if Type is recovery")
		}
		if len(req.Data) > 0 {
			return types.NewErrInvalidGenerateLinkRequest("data must not be provided if Type is recovery")
		}
		if req.Password != "" {
			return types.NewErrInvalidGenerateLinkRequest("password must not be provided if Type is recovery")
		}
	case types.LinkTypeEmailChangeCurrent, types.LinkTypeEmailChangeNew:
		if req.Email == "" || req.NewEmail == "" {
			return types.NewErrInvalidGenerateLinkRequest("email and new_email must be provided if Type is email_change_current or email_change_new")
		}
		if len(req.Data) > 0 {
			return types.NewErrInvalidGenerateLinkRequest("data must not be provided if Type is email_change_current or email_change_new")
		}
		if req.Password != "" {
			return types.NewErrInvalidGenerateLinkRequest("password must not be provided if Type is email_change_current or email_change_new")
		}
	}

	return nil
}

// POST /admin/generate_link
//
// Returns the corresponding email action link based on the type specified.
// Among other things, the response also contains the query params of the action
// link as separate JSON fields for convenience (along with the email OTP from
// which the corresponding token is generated).
func (c *Client) AdminGenerateLink(req types.AdminGenerateLinkRequest) (*types.AdminGenerateLinkResponse, error) {
	err := validateAdminGenerateLinkRequest(req)
	if err != nil {
		return nil, err
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	r, err := c.newRequest(adminGenerateLinkPath, http.MethodPost, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fullBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("response status code %d", resp.StatusCode)
		}
		return nil, fmt.Errorf("response status code %d: %s", resp.StatusCode, fullBody)
	}

	var res types.AdminGenerateLinkResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
