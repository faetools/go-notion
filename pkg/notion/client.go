package notion

import (
	"context"
	"fmt"
	"net/http"

	"github.com/faetools/client"
	"github.com/google/uuid"
)

const (
	versionHeader           = "Notion-Version"
	version                 = "2022-02-22"
	maxPageSize    PageSize = 100
	maxPageSizeInt          = 100
)

// NewDefaultClient returns a new client with the default options.
func NewDefaultClient(bearer string, opts ...client.Option) (*Client, error) {
	opts = append([]client.Option{
		client.WithBearer(bearer),
		client.WithRequestEditorFn(func(_ context.Context, req *http.Request) error {
			req.Header.Set(versionHeader, version)
			return nil
		}),
	}, opts...)

	return NewClient(opts...)
}

// Error ensures responses with an error fulfill the error interface.
func (e *Error) Error() string {
	return fmt.Sprintf("%d %s: %s - %s", e.Status, http.StatusText(e.Status), e.Code, e.Message)
}

// GetNotionPage return the notion page or an error.
func (c Client) GetNotionPage(ctx context.Context, id Id) (*Page, error) {
	resp, err := c.GetPage(ctx, id)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode() {
	case http.StatusOK: // ok
		return resp.JSON200, nil
	case http.StatusBadRequest:
		return nil, resp.JSON400
	case http.StatusNotFound:
		return nil, resp.JSON404
	case http.StatusTooManyRequests:
		return nil, resp.JSON429
	default:
		return nil, fmt.Errorf("unknown error response: %v", string(resp.Body))
	}
}

// GetNotionDatabase returns the notion database or an error.
func (c Client) GetNotionDatabase(ctx context.Context, id Id) (*Database, error) {
	resp, err := c.GetDatabase(ctx, id)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode() {
	case http.StatusOK: // ok
		return resp.JSON200, nil
	case http.StatusBadRequest:
		return nil, resp.JSON400
	case http.StatusNotFound:
		return nil, resp.JSON404
	case http.StatusTooManyRequests:
		return nil, resp.JSON429
	default:
		return nil, fmt.Errorf("unknown error response: %v", string(resp.Body))
	}
}

func ensureDatabaseIsValid(db *Database) {
	// set mandatory values
	db.Object = "database"
	if db.Parent != nil && db.Parent.Type == "" {
		db.Parent.Type = "page_id"
	}

	// initialize properties
	if db.Properties == nil {
		db.Properties = PropertyMetaMap{"Title": TitleProperty}
		return
	}

	// make sure a title property is present
	for _, prop := range db.Properties {
		if prop.Title != nil {
			return
		}
	}

	db.Properties["Title"] = TitleProperty
}

// CreateNotionDatabase creates a notion database or returns an error.
func (c Client) CreateNotionDatabase(ctx context.Context, db Database) (*Database, error) {
	ensureDatabaseIsValid(&db)

	// create a UUID for the new database
	db.Id = UUID(uuid.NewString())

	resp, err := c.CreateDatabase(ctx, CreateDatabaseJSONRequestBody(db))
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode() {
	case http.StatusOK: // ok
		return resp.JSON200, nil
	case http.StatusBadRequest:
		return nil, resp.JSON400
	case http.StatusNotFound:
		return nil, resp.JSON404
	case http.StatusTooManyRequests:
		return nil, resp.JSON429
	default:
		return nil, fmt.Errorf("unknown error response: %v", string(resp.Body))
	}
}

// UpdateNotionDatabase updates a notion database or returns an error.
func (c Client) UpdateNotionDatabase(ctx context.Context, db Database) (*Database, error) {
	// can't be present when updating
	db.Parent = nil
	db.CreatedTime = nil

	ensureDatabaseIsValid(&db)

	resp, err := c.UpdateDatabase(ctx, Id(db.Id), UpdateDatabaseJSONRequestBody(db))
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode() {
	case http.StatusOK: // ok
		return resp.JSON200, nil
	case http.StatusBadRequest:
		return nil, resp.JSON400
	case http.StatusNotFound:
		return nil, resp.JSON404
	case http.StatusTooManyRequests:
		return nil, resp.JSON429
	default:
		return nil, fmt.Errorf("unknown error response: %v", string(resp.Body))
	}
}
