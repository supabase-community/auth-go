package endpoints

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/tomnomnom/linkheader"

	"github.com/supabase-community/auth-go/types"
)

const adminAuditPath = "/admin/audit"

// GET /admin/audit
//
// Get audit logs.
//
// May optionally specify a query to use for filtering the audit logs. The
// column and value must be specified if using a query.
//
// The result may also be paginated. By default, 50 results will be returned
// per request. This can be configured with PerPage in the request. The response
// will include the total number of results, as well as the total number of pages
// and, if not already on the last page, the next page number.
func (c *Client) AdminAudit(req types.AdminAuditRequest) (*types.AdminAuditResponse, error) {
	if req.Query != nil {
		if req.Query.Column != types.AuditQueryColumnAuthor && req.Query.Column != types.AuditQueryColumnAction && req.Query.Column != types.AuditQueryColumnType {
			return nil, types.ErrInvalidAdminAuditRequest
		}
		if req.Query.Value == "" {
			return nil, types.ErrInvalidAdminAuditRequest
		}
	}

	r, err := c.newRequest(adminAuditPath, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	q := r.URL.Query()
	if req.Query != nil {
		q.Add("query", fmt.Sprintf("%s:%s", req.Query.Column, req.Query.Value))
	}
	if req.Page != 0 && req.PerPage != 0 {
		q.Add("page", fmt.Sprintf("%d", req.Page))
		q.Add("per_page", fmt.Sprintf("%d", req.PerPage))
	}
	r.URL.RawQuery = q.Encode()

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

	var logs []types.AuditLogEntry
	if err := json.NewDecoder(resp.Body).Decode(&logs); err != nil {
		return nil, err
	}

	// Result count should be given in X-Total-Count header.
	count := resp.Header.Get("X-Total-Count")
	resultCount := 0
	if count != "" {
		resultCount, _ = strconv.Atoi(count)
	}

	// Parse Link header from response to get total pages
	links := linkheader.Parse(resp.Header.Get("Link"))

	// Header should only contain one 'last' link
	var totPages uint = req.Page
	l := links.FilterByRel("last")
	if len(l) == 1 {
		// Parse it's URL as a URL to get the query params
		lastURL, err := url.Parse(l[0].URL)
		if err == nil {
			// Look for the ?page=X query param
			last := lastURL.Query().Get("page")
			lastPage, err := strconv.Atoi(last)
			if err == nil {
				totPages = uint(lastPage)
			}
		}
	}

	// Header may contain one 'next' link
	var nextPage uint = 0
	n := links.FilterByRel("next")
	if len(n) == 1 {
		nextURL, err := url.Parse(n[0].URL)
		if err == nil {
			next := nextURL.Query().Get("page")
			nPage, err := strconv.Atoi(next)
			if err == nil {
				nextPage = uint(nPage)
			}
		}
	}

	return &types.AdminAuditResponse{
		Logs: logs,

		TotalCount: resultCount,
		NextPage:   nextPage,
		TotalPages: totPages,
	}, nil
}
