package clueo

import (
	"fmt"
	"io"
	"net/http"
)

// WithStatuses returns a roundtripper that returns an error if the response is not OK.
func WithOK(base http.RoundTripper) http.RoundTripper {
	return WithStatuses(base, http.StatusOK)
}

// WithStatuses returns a roundtripper that returns an error if the response does not have one of the desired status codes.
func WithStatuses(base http.RoundTripper, codes ...int) http.RoundTripper {
	return WithResponse(base, func(resp *http.Response) error {
		for _, code := range codes {
			if resp.StatusCode == code {
				return nil
			}
		}

		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err == nil && len(body) > 0 {
			return fmt.Errorf("got %s response and body: %s", resp.Status, string(body))
		}

		return fmt.Errorf("got %s response", resp.Status)
	})
}
