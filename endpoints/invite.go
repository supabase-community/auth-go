package endpoints

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/supabase-community/auth-go/types"
)

const invitePath = "/invite"

// POST /invite
//
// Invites a new user with an email.
// This endpoint requires the service_role or supabase_admin JWT set using WithToken.
func (c *Client) Invite(req types.InviteRequest) (*types.InviteResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	r, err := c.newRequest(invitePath, http.MethodPost, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	if req.RedirectTo != "" {
		q := r.URL.Query()
		q.Add("redirect_to", req.RedirectTo)
		r.URL.RawQuery = q.Encode()
	}

	resp, err := c.client.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, handleErrorResponse(resp)
	}

	var res types.InviteResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}
