package clueo

import "net/http"

type requestTransform struct {
	base         http.RoundTripper
	reqTransform func(*http.Request) error
}

func WithRequest(base http.RoundTripper, reqTransform func(*http.Request) error) http.RoundTripper {
	return &requestTransform{
		base:         base,
		reqTransform: reqTransform,
	}
}

// RoundTrip implements http.RoundTripper
func (tr *requestTransform) RoundTrip(req *http.Request) (*http.Response, error) {
	if err := tr.reqTransform(req); err != nil {
		return nil, err
	}

	return tr.base.RoundTrip(req)
}
