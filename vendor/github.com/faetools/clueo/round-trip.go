package clueo

import "net/http"

type RoundTrip func(*http.Request) (*http.Response, error)

func (rt RoundTrip) RoundTrip(req *http.Request) (*http.Response, error) {
	return rt(req)
}
