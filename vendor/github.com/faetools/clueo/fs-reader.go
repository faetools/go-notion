package clueo

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"net/http"
)

var numMimeTypes = len(mimeTypeExtensions) + 1 // plus other

type fsReader struct {
	fs       fs.FS
	mimeType string
}

func newFSReader(files fs.FS, mimeType string) http.RoundTripper {
	return &fsReader{
		fs:       files,
		mimeType: mimeType,
	}
}

// NewFSReader returns a new roundtripper that returns responses using files in a filesystem.
func NewFSReader(fs fs.FS) http.RoundTripper {
	fallbacks := make([]http.RoundTripper, 0, numMimeTypes)

	for mimeType := range mimeTypeExtensions {
		// try to read from the files of the mime type
		fallbacks = append(fallbacks,
			AddResponseContentType(&fsReader{fs: fs, mimeType: mimeType}, mimeType))
	}

	// else try to read from other files
	fallbacks = append(fallbacks, &fsReader{fs: fs})

	return WithFallbackForFSReader(fallbacks...)
}

func (c *fsReader) RoundTrip(req *http.Request) (*http.Response, error) {
	f, err := c.fs.Open(requestToPath(req, c.mimeType))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	body, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return newResponse(req, http.StatusOK, body), nil
}

func newResponse(req *http.Request, status int, body []byte) *http.Response {
	return &http.Response{
		StatusCode:    status,
		Status:        fmt.Sprintf("%d %s", status, http.StatusText(status)),
		Body:          io.NopCloser(bytes.NewReader(body)),
		Request:       req,
		Header:        http.Header{},
		ContentLength: int64(len(body)),
	}
}
