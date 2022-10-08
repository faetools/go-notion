package clueo

import (
	"errors"
	"io/fs"
	"net/http"
)

type fallbacker struct {
	fallbacks []http.RoundTripper
	condition func(error) bool
}

func WithFallback(baseAndFallbacks ...http.RoundTripper) http.RoundTripper {
	return WithFallbackIf(func(error) bool { return true }, baseAndFallbacks...)
}

func WithFallbackForFSReader(baseAndFallbacks ...http.RoundTripper) http.RoundTripper {
	return WithFallbackIf(func(err error) bool {
		if errors.Is(err, fs.ErrNotExist) {
			return true
		}

		pathErr := &fs.PathError{}
		if errors.As(err, &pathErr) && pathErr.Err.Error() == "not a directory" {
			return true
		}

		return false
	}, baseAndFallbacks...)
}

func WithFallbackIf(condition func(error) bool, fallbacks ...http.RoundTripper) http.RoundTripper {
	if len(fallbacks) == 0 {
		panic("can't call WithFallbackIf without any roundtripper")
	}

	return &fallbacker{
		condition: condition,
		fallbacks: fallbacks,
	}
}

func (tr *fallbacker) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	for _, fallback := range tr.fallbacks {
		resp, err = fallback.RoundTrip(req)
		if err != nil && tr.condition(err) {
			continue
		}

		return resp, err
	}

	return resp, err // return the last error
}
