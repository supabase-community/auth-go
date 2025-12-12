package endpoints

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
		fullBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("response status code %d", resp.StatusCode)
		}
		return fmt.Errorf("response status code %d: %s", resp.StatusCode, fullBody)
	}

	return nil
}
