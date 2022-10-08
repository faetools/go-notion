package clueo

import (
	"net/http"

	"github.com/spf13/afero"
)

// NewReadWriter returns a roundtripper that both reads and writes to a filesystem,
// depending on if it finds the file associated with a request.
func NewReadWriter(base http.RoundTripper, fs afero.Fs) http.RoundTripper {
	return WithFallbackForFSReader(
		NewFSReader(afero.NewIOFS(fs)),
		NewFSWriter(WithOK(base), fs),
	)
}
