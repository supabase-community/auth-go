package endpoints

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/supabase-community/auth-go/types"
)

const magiclinkPath = "/magiclink"

// POST /magiclink
//
// DEPRECATED: Use /otp with Email and CreateUser=true instead of /magiclink.
//
// Magic Link. Will deliver a link (e.g.
// /verify?type=magiclink&token=fgtyuf68ddqdaDd) to the user based on email
// address which they can use to redeem an access_token.
//
// By default Magic Links can only be sent once every 60 seconds.
func (c *Client) Magiclink(req types.MagiclinkRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	r, err := c.newRequest(magiclinkPath, http.MethodPost, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	resp, err := c.client.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return handleErrorResponse(resp)
	}

	return nil
}
