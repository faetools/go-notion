package clueo

import (
	"net/http"
)

type responseTransform struct {
	base          http.RoundTripper
	respTransform func(*http.Response) error
}

func WithResponse(base http.RoundTripper, respTransform func(*http.Response) error) http.RoundTripper {
	return &responseTransform{
		base:          base,
		respTransform: respTransform,
	}
}

// RoundTrip implements http.RoundTripper
func (tr *responseTransform) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := tr.base.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	if err := tr.respTransform(resp); err != nil {
		return nil, err
	}

	return resp, nil
}
