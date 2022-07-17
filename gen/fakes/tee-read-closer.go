package main

import "io"

type readCloser struct {
	io.Reader
	io.Closer
}

// TeeReadCloser returns a ReadCloser that writes to w what it reads from rc.
// When the ReadCloser is closed, rc is closed.
func TeeReadCloser(rc io.ReadCloser, w io.Writer) io.ReadCloser {
	return &readCloser{
		Reader: io.TeeReader(rc, w),
		Closer: rc,
	}
}
