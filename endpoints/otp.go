package endpoints

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/supabase-community/auth-go/types"
)

const otpPath = "/otp"

// POST /otp
// One-Time-Password. Will deliver a magiclink or SMS OTP to the user depending
// on whether the request contains an email or phone key.
//
// If CreateUser is true, the user will be automatically signed up if the user
// doesn't exist.
func (c *Client) OTP(req types.OTPRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	r, err := c.newRequest(otpPath, http.MethodPost, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	q := r.URL.Query()
	if req.RedirectTo != "" {
		q.Add("redirect_to", req.RedirectTo)
	}
	if req.EmailRedirectTo != "" {
		q.Add("redirect_to", req.EmailRedirectTo)
	}
	if len(q) > 0 {
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
