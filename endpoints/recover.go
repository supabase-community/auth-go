package endpoints

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/supabase-community/auth-go/types"
)

const recoverPath = "/recover"

// POST /recover
//
// Password recovery. Will deliver a password recovery mail to the user based
// on email address.
//
// By default recovery links can only be sent once every 60 seconds.
func (c *Client) Recover(req types.RecoverRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	r, err := c.newRequest(recoverPath, http.MethodPost, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	if req.RedirectTo != "" {
		q := r.URL.Query()
		q.Add("redirect_to", req.RedirectTo)
		r.URL.RawQuery = q.Encode()
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
