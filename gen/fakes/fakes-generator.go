package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/faetools/cgtools"
	"github.com/faetools/client"
	"golang.org/x/sync/errgroup"
)

type fakesGenerator struct {
	gen   *cgtools.Generator
	files client.HTTPRequestDoer
	cli   client.HTTPRequestDoer
	wg    *errgroup.Group
}

func (c *fakesGenerator) Do(req *http.Request) (*http.Response, error) {
	resp, err := c.files.Do(req)
	if err == nil && resp.StatusCode == http.StatusOK {
		return resp, nil
	}

	resp, err = c.cli.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s %s - got %s", req.Method, req.URL.Path, resp.Status)
	}

	r, w := io.Pipe()

	// write to w and close when done
	resp.Body = TeeReadCloser(resp.Body, w)

	c.wg.Go(func() error {
		return c.gen.Write(req.URL.Path+".json", r)
	})

	return resp, nil
}
