package clueo

import (
	"fmt"
	"net/http"
	"strings"
)

const (
	HeaderContentType         = "Content-Type"
	MimeApplicationJSON       = "application/json"
	MIMETextHTML              = "text/html"
	MIMEAudioOgg              = "audio/ogg"
	MIMEApplicationJavaScript = "application/javascript"
	MIMEImageJPEG             = "image/jpeg"
	MIMEVideoMP4              = "video/mp4"
)

var mimeTypeExtensions = map[string]string{
	MimeApplicationJSON:       ".json",
	MIMETextHTML:              ".html",
	MIMEAudioOgg:              ".ogg",
	MIMEApplicationJavaScript: ".js",
	MIMEImageJPEG:             ".jpeg",
	MIMEVideoMP4:              ".mp4",
}

func WithResponseContainsContentType(base http.RoundTripper, mimeType string) http.RoundTripper {
	return WithResponse(base, func(resp *http.Response) error {
		contentType := resp.Header.Get(HeaderContentType)

		if strings.Contains(contentType, mimeType) {
			return nil
		}

		return fmt.Errorf("header content type is %q", contentType)
	})
}

func WithResponseContainsContentTypeApplicationJSON(base http.RoundTripper) http.RoundTripper {
	return WithResponseContainsContentType(base, MimeApplicationJSON)
}

func WithResponseContainsContentTypeTextHTML(base http.RoundTripper) http.RoundTripper {
	return WithResponseContainsContentType(base, MIMETextHTML)
}

func AddResponseContentType(base http.RoundTripper, value string) http.RoundTripper {
	return WithResponse(base, func(resp *http.Response) error {
		resp.Header.Add(HeaderContentType, value)
		return nil
	})
}

func AddResponseContentTypeApplicationJSON(base http.RoundTripper) http.RoundTripper {
	return AddResponseContentType(base, MimeApplicationJSON)
}

func AddResponseContentTypeTextHTML(base http.RoundTripper) http.RoundTripper {
	return AddResponseContentType(base, MIMETextHTML)
}
