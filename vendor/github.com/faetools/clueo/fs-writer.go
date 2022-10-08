package clueo

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/afero"
)

var (
	// remove anything after a ';'
	reCleanMimeType = regexp.MustCompile(`;.*$`)

	// delete some queries so the filename is not too long
	queriesToDelete = []string{
		"X-Amz-Algorithm",
		"X-Amz-Credential",
		"X-Amz-Expires",
		"X-Amz-Content-Sha256",
		"X-Amz-Signature",
		"X-Amz-SignedHeaders",
		"X-Amz-Date",
		"x-id",
	}
)

type fsWriter struct {
	base http.RoundTripper
	fs   afero.Fs
}

// NewFSWriter returns a new roundtripper that writes responses to a filesystem.
func NewFSWriter(base http.RoundTripper, files afero.Fs) http.RoundTripper {
	return &fsWriter{
		base: base,
		fs:   files,
	}
}

func (tr *fsWriter) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := tr.base.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()

	path := requestToPath(req, resp.Header.Get(HeaderContentType))

	if err := tr.fs.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return nil, err
	}

	if err := afero.WriteFile(tr.fs, path, body, os.ModePerm); err != nil {
		return nil, err
	}

	resp.Body = io.NopCloser(bytes.NewReader(body))

	return resp, nil
}

func requestToPath(req *http.Request, mimeType string) (path string) {
	q := req.URL.Query()

	// delete some queries so the filename is not too long
	for _, key := range queriesToDelete {
		q.Del(key)
	}

	path = strings.TrimPrefix(filepath.Clean((&url.URL{
		// copy of URL without scheme
		Opaque:      req.URL.Opaque,
		User:        req.URL.User,
		Host:        req.URL.Host,
		Path:        req.URL.Path,
		RawPath:     req.URL.RawPath,
		OmitHost:    req.URL.OmitHost,
		ForceQuery:  req.URL.ForceQuery,
		RawQuery:    q.Encode(),
		Fragment:    req.URL.Fragment,
		RawFragment: req.URL.RawFragment,
	}).String()), "/")

	defer func() {
		// add file type extension
		mimeType = reCleanMimeType.ReplaceAllString(mimeType, "")
		if ext := mimeTypeExtensions[mimeType]; ext != "" {
			if filepath.Ext(path) != ext {
				path = filepath.Join(path, "index"+ext)
			}
		}
	}()

	if req.Body == nil {
		return
	}

	body, err := req.GetBody()
	if err != nil {
		return
	}
	defer body.Close()

	b, err := io.ReadAll(body)
	if err != nil {
		return
	}

	path = filepath.Join(path, string(b))

	return
}
