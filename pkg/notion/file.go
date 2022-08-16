package notion

import (
	"fmt"
	"net/url"
	"path/filepath"
	"time"
)

// URL return the URL of the file.
func (f File) URL() string {
	switch f.Type {
	case FileTypeExternal:
		return f.External.Url
	case FileTypeFile:
		return f.File.Url
	default:
		panic(fmt.Errorf("invalid File of type %q", f.Type))
	}
}

func (f File) GetName() string {
	if f.Name != nil {
		return *f.Name
	}

	u, err := url.Parse(f.URL())
	if err != nil {
		panic(fmt.Errorf("invalid File with unparsable URL: %w", err))
	}

	return filepath.Base(u.Path)
}

func (f File) ExpiryTime() *time.Time {
	if f.File != nil {
		return &f.File.ExpiryTime
	}

	return nil
}

// URL return the URL of the file.
func (f FileWithCaption) URL() string {
	switch f.Type {
	case FileWithCaptionTypeExternal:
		return f.External.Url
	case FileWithCaptionTypeFile:
		return f.File.Url
	default:
		panic(fmt.Errorf("invalid FileWithCaption of type %q", f.Type))
	}
}

func (f FileWithCaption) ExpiryTime() *time.Time {
	if f.File != nil {
		return &f.File.ExpiryTime
	}

	return nil
}

func (f FileWithCaption) GetName() string {
	u, err := url.Parse(f.URL())
	if err != nil {
		panic(fmt.Errorf("invalid FileWithCaption with unparsable URL: %w", err))
	}

	return filepath.Base(u.Path)
}

func (f FileWithCaption) GetFile() File {
	return File{
		External: f.External,
		File:     f.File,
		Type:     FileType(f.Type),
	}
}
