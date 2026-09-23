package endpoints

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/supabase-community/auth-go/types"
)

const resendPath = "/resend"

// POST /resend
//
// Resends an existing signup confirmation email, email change email,
// SMS OTP or phone change OTP.
//
// You can specify a redirect url when you resend an email link using
// the emailRedirectTo option.
func (c *Client) Resend(req types.ResendRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	r, err := c.newRequest(resendPath, http.MethodPost, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	if req.EmailRedirectTo != "" {
		q := r.URL.Query()
		q.Add("redirect_to", req.EmailRedirectTo)
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
