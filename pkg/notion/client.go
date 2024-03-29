package notion

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/faetools/client"
	"github.com/google/uuid"
)

const (
	versionHeader = "Notion-Version"
	version       = "2022-06-28"
)

var (
	maxPageSizeInt          = 100
	maxPageSize    PageSize = PageSize(maxPageSizeInt)

	// ErrBadRequest is returned when we get a 502 Bad Gateway or 504 Gateway Timeout response.
	ErrGatewayIssue = errors.New("gateway issue")
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

// CreateNotionPage creates a notion page or returns an error.
func (c Client) CreateNotionPage(ctx context.Context, p Page) (*Page, error) {
	// needs to be set
	p.Object = "page"

	if p.Id == "" {
		p.Id = UUID(uuid.NewString())
	}

	if p.Properties == nil {
		p.Properties = PropertyValueMap{}
	}

	resp, err := c.CreatePage(ctx, CreatePageJSONRequestBody(p))
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
	case http.StatusBadGateway, http.StatusGatewayTimeout:
		return nil, fmt.Errorf("%w (%s)", ErrGatewayIssue, resp.HTTPResponse.Status)
	default:
		return nil, fmt.Errorf("unknown %s response: %v",
			resp.HTTPResponse.Status, string(resp.Body))
	}
}

// GetNotionPage return the notion page or an error.
func (c Client) GetNotionPage(ctx context.Context, id Id) (*Page, error) {
	p, err := c.getNotionPage(ctx, id)
	if err != nil {
		// retry once if we get a Bad Gateway error
		if errors.Is(err, ErrGatewayIssue) {
			fmt.Printf("Got error %v getting notion page, retrying...\n", err)
			time.Sleep(time.Second)
			p, err = c.getNotionPage(ctx, id)
		}

		return p, err
	}

	return p, nil
}

func (c Client) getNotionPage(ctx context.Context, id Id) (*Page, error) {
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
	case http.StatusBadGateway, http.StatusGatewayTimeout:
		return nil, fmt.Errorf("%w (%s)", ErrGatewayIssue, resp.HTTPResponse.Status)
	default:
		return nil, fmt.Errorf("unknown %s response: %v",
			resp.HTTPResponse.Status, string(resp.Body))
	}
}

// UpdateNotionPage updates the notion page or returns an error.
func (c Client) UpdateNotionPage(ctx context.Context, p Page) (*Page, error) {
	// can't be present when updating
	p.CreatedTime = nil

	// needs to be set
	p.Object = "page"

	props := p.Properties
	for key, prop := range props {
		switch prop.Type {
		case PropertyTypeCreatedTime,
			PropertyTypeCreatedBy,
			PropertyTypeLastEditedTime,
			PropertyTypeLastEditedBy:
			// we can't update these
			delete(props, key)
		case PropertyTypeSelect:
			// if we want this to be empty, we have to set it to nil
			// if prop.Select != nil && prop.Select.Name == "" {
			// 	prop.Select = nil
			// 	props[key] = prop
			// }
		}
	}

	resp, err := c.UpdatePage(ctx, Id(p.Id), UpdatePageJSONRequestBody(p))
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
	case http.StatusBadGateway, http.StatusGatewayTimeout:
		return nil, fmt.Errorf("%w (%s)", ErrGatewayIssue, resp.HTTPResponse.Status)
	default:
		return nil, fmt.Errorf("unknown %s response: %v",
			resp.HTTPResponse.Status, string(resp.Body))
	}
}

// GetNotionBlock returns the notion block or an error.
func (c Client) GetNotionBlock(ctx context.Context, id Id) (*Block, error) {
	resp, err := c.GetBlock(ctx, id)
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
	case http.StatusBadGateway, http.StatusGatewayTimeout:
		return nil, fmt.Errorf("%w (%s)", ErrGatewayIssue, resp.HTTPResponse.Status)
	default:
		return nil, fmt.Errorf("unknown %s response: %v",
			resp.HTTPResponse.Status, string(resp.Body))
	}
}

// GetNotionDatabase returns the notion database or an error.
func (c Client) GetNotionDatabase(ctx context.Context, id Id) (*Database, error) {
	db, err := c.getNotionDatabase(ctx, id)
	if err != nil {
		// retry once if we get a Bad Gateway error
		if errors.Is(err, ErrGatewayIssue) {
			fmt.Printf("Got error %v updating notion database, retrying...\n", err)
			time.Sleep(time.Second)
			db, err = c.getNotionDatabase(ctx, id)
		}

		return db, err
	}

	return db, nil
}

func (c Client) getNotionDatabase(ctx context.Context, id Id) (*Database, error) {
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
	case http.StatusBadGateway, http.StatusGatewayTimeout:
		return nil, fmt.Errorf("%w (%s)", ErrGatewayIssue, resp.HTTPResponse.Status)
	default:
		return nil, fmt.Errorf("unknown %s response: %v",
			resp.HTTPResponse.Status, string(resp.Body))
	}
}

// GetAllDatabaseEntries returns all database entries or an error.
func (c Client) GetAllDatabaseEntries(ctx context.Context, id Id) (Pages, error) {
	return c.GetDatabaseEntries(ctx, id, nil, nil)
}

// GetDatabaseEntries return filtered and sorted database entries or an error.
func (c Client) GetDatabaseEntries(ctx context.Context, id Id, filter *Filter, sorts *Sorts) (Pages, error) {
	entries := Pages{}
	query := DatabaseQuery{
		Filter:   filter,
		PageSize: maxPageSizeInt,
		Sorts:    sorts,
	}

	for i := 0; ; i++ {
		results, next, err := c.QueryNotionDatabase(ctx, id, query)

		// retry once if we get a Bad Gateway error
		if errors.Is(err, ErrGatewayIssue) {
			fmt.Printf("Got error %v on page %d, retrying...\n", err, i)
			time.Sleep(time.Second)
			results, next, err = c.QueryNotionDatabase(ctx, id, query)
		}

		if err != nil {
			return nil, fmt.Errorf("on page %d: %w", i, err)
		}

		entries = append(entries, results...)

		if next == nil {
			return entries, nil
		}

		query.StartCursor = (*UUID)(next)
	}
}

func (c Client) QueryNotionDatabase(ctx context.Context, id Id, query DatabaseQuery) (Pages, *NextCursor, error) {
	resp, err := c.QueryDatabase(ctx, id, QueryDatabaseJSONRequestBody(query))
	if err != nil {
		return nil, nil, err
	}

	switch resp.StatusCode() {
	case http.StatusOK: // ok
		if resp.JSON200.HasMore {
			return resp.JSON200.Results, resp.JSON200.NextCursor, nil
		}

		return resp.JSON200.Results, nil, nil
	case http.StatusBadRequest:
		return nil, nil, resp.JSON400
	case http.StatusNotFound:
		return nil, nil, resp.JSON404
	case http.StatusTooManyRequests:
		return nil, nil, resp.JSON429
	case http.StatusBadGateway, http.StatusGatewayTimeout:
		return nil, nil, fmt.Errorf("%w (%s)", ErrGatewayIssue, resp.HTTPResponse.Status)
	default:
		return nil, nil, fmt.Errorf("unknown %s response: %v",
			resp.HTTPResponse.Status, string(resp.Body))
	}
}

func ensureDatabaseIsValid(db *Database) {
	// set mandatory values
	db.Object = "database"
	if db.Parent != nil && db.Parent.Type == "" {
		db.Parent.Type = "page_id"
	}

	if db.Description == nil {
		db.Description = RichTexts{}
	}
}

func ensureNewDatabaseIsValid(db *Database) {
	ensureDatabaseIsValid(db)

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
	ensureNewDatabaseIsValid(&db)

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
	case http.StatusBadGateway, http.StatusGatewayTimeout:
		return nil, fmt.Errorf("%w (%s)", ErrGatewayIssue, resp.HTTPResponse.Status)
	default:
		return nil, fmt.Errorf("unknown %s response: %v",
			resp.HTTPResponse.Status, string(resp.Body))
	}
}

// UpdateNotionDatabase updates a notion database or returns an error.
func (c Client) UpdateNotionDatabase(ctx context.Context, db Database) (*Database, error) {
	// can't be present when updating
	db.Parent = nil
	db.CreatedTime = nil

	ensureDatabaseIsValid(&db)

	updated, err := c.updateNotionDatabase(ctx, db)
	if err != nil {
		// retry once if we get a Bad Gateway error
		if errors.Is(err, ErrGatewayIssue) {
			fmt.Printf("Got error %v updating notion database, retrying...\n", err)
			time.Sleep(time.Second)
			updated, err = c.updateNotionDatabase(ctx, db)
		}

		return updated, err
	}

	return updated, nil
}

func (c Client) updateNotionDatabase(ctx context.Context, db Database) (*Database, error) {
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
	case http.StatusBadGateway, http.StatusGatewayTimeout:
		return nil, fmt.Errorf("%w (%s)", ErrGatewayIssue, resp.HTTPResponse.Status)
	default:
		return nil, fmt.Errorf("unknown %s response: %v",
			resp.HTTPResponse.Status, string(resp.Body))
	}
}

// ListAllUsers returns all users in the workspace.
func (c Client) ListAllUsers(ctx context.Context) (Users, error) {
	users := Users{}

	var cursor *StartCursor
	for {
		resp, err := c.ListUsers(ctx, &ListUsersParams{
			PageSize:    &maxPageSize,
			StartCursor: cursor,
		})
		if err != nil {
			return nil, err
		}

		switch resp.StatusCode() {
		case http.StatusOK: // ok
		case http.StatusBadRequest:
			return nil, resp.JSON400
		case http.StatusNotFound:
			return nil, resp.JSON404
		case http.StatusTooManyRequests:
			return nil, resp.JSON429
		case http.StatusBadGateway, http.StatusGatewayTimeout:
			return nil, fmt.Errorf("%w (%s)", ErrGatewayIssue, resp.HTTPResponse.Status)
		default:
			return nil, fmt.Errorf("unknown %s response: %v",
				resp.HTTPResponse.Status, string(resp.Body))
		}

		users = append(users, resp.JSON200.Results...)

		if !resp.JSON200.HasMore {
			return users, nil
		}

		cursor = (*StartCursor)(resp.JSON200.NextCursor)
	}
}

// GetAllBlocks returns all blocks of a given page or block.
func (c Client) GetAllBlocks(ctx context.Context, id Id) (Blocks, error) {
	var (
		all    Blocks
		blocks Blocks
		next   *StartCursor
		err    error
	)

	for i := 0; ; i++ {
		blocks, next, err = c.GetNextBlocks(ctx, id, next)
		if err != nil {
			return nil, fmt.Errorf("page %d of getting blocks for %s: %w", i, id, err)
		}

		all = append(all, blocks...)

		if next == nil {
			return all, nil
		}
	}
}

// GetNextBlocks gets the next blocks, starting at the cursor.
func (c Client) GetNextBlocks(ctx context.Context, id Id, cursor *StartCursor) (
	Blocks, *StartCursor, error,
) {
	blocks, next, err := c.getNextBlocks(ctx, id, cursor)

	// retry once if we get a Bad Gateway error
	if errors.Is(err, ErrGatewayIssue) {
		fmt.Printf("Got error %v, retrying...\n", err)
		time.Sleep(time.Second)
		blocks, next, err = c.getNextBlocks(ctx, id, cursor)
	}

	return blocks, next, err
}

func (c Client) getNextBlocks(ctx context.Context, id Id, cursor *StartCursor) (
	Blocks, *StartCursor, error,
) {
	resp, err := c.GetBlocks(ctx, id, &GetBlocksParams{
		PageSize:    &maxPageSize,
		StartCursor: cursor,
	})
	if err != nil {
		return nil, nil, err
	}

	switch resp.StatusCode() {
	case http.StatusOK: // ok
		var next *StartCursor
		if resp.JSON200.HasMore {
			next = (*StartCursor)(resp.JSON200.NextCursor)
		}

		return resp.JSON200.Results, next, nil
	case http.StatusBadRequest:
		return nil, nil, resp.JSON400
	case http.StatusNotFound:
		return nil, nil, resp.JSON404
	case http.StatusBadGateway, http.StatusGatewayTimeout:
		return nil, nil, fmt.Errorf("%w (%s)", ErrGatewayIssue, resp.HTTPResponse.Status)
	default:
		return nil, nil, fmt.Errorf("unknown %s response: %v",
			resp.HTTPResponse.Status, string(resp.Body))
	}
}

// PageWithinScope checks if an ancestor of the page has the stated UUID.
func (c Client) PageWithinScope(ctx context.Context, scope UUID, p *Page) (bool, error) {
	return c.parentWithinScope(ctx, scope, p.Parent)
}

// PageWithinScope checks if an ancestor of the database has the stated UUID.
func (c Client) DatabaseWithinScope(ctx context.Context, scope UUID, db *Database) (bool, error) {
	return c.parentWithinScope(ctx, scope, db.Parent)
}

// PageWithinScope checks if an ancestor of the block has the stated UUID.
func (c Client) BlockWithinScope(ctx context.Context, scope UUID, block *Block) (bool, error) {
	return c.parentWithinScope(ctx, scope, &block.Parent)
}

func (c Client) parentWithinScope(ctx context.Context, scope UUID, p *Parent) (bool, error) {
	switch {
	case p == nil:
		return false, nil
	case p.ID() == scope:
		return true, nil
	default:
		switch p.Type {
		case ParentTypeBlockId:
			parent, err := c.GetNotionBlock(ctx, Id(*p.BlockId))
			if err != nil {
				return false, err
			}

			return c.BlockWithinScope(ctx, scope, parent)
		case ParentTypeDatabaseId:
			parent, err := c.GetNotionDatabase(ctx, Id(*p.DatabaseId))
			if err != nil {
				return false, err
			}

			return c.DatabaseWithinScope(ctx, scope, parent)
		case ParentTypePageId:
			parent, err := c.GetNotionPage(ctx, Id(*p.PageId))
			if err != nil {
				return false, err
			}

			return c.PageWithinScope(ctx, scope, parent)
		case ParentTypeWorkspace:
			return false, nil
		default:
			return false, fmt.Errorf("invalid parent type %s", p.Type)
		}
	}

	return false, nil
}

func (c Client) GetNotionPagesByTitle(
	ctx context.Context, title string,
) (Pages, error) {
	if title == "" {
		return Pages{}, nil
	}

	resp, err := c.Search(ctx, SearchJSONRequestBody{
		Query: &title,
		Filter: &SearchFilter{
			Value:    SearchFilterValuePage,
			Property: SearchFilterPropertyObject,
		},
		// StartCursor: &"", // TODO pagination
		PageSize: &maxPageSizeInt,
	})
	if err != nil {
		return nil, fmt.Errorf("getting page by title for %s: %w", title, err)
	}

	switch resp.StatusCode() {
	case http.StatusOK: // ok
		pages := make(Pages, len(resp.JSON200.Results))
		for i, res := range resp.JSON200.Results {
			pages[i] = *res.Page
		}

		return pages, nil
	case http.StatusBadRequest:
		return nil, resp.JSON400
	case http.StatusNotFound:
		return nil, resp.JSON404
	case http.StatusBadGateway, http.StatusGatewayTimeout:
		return nil, fmt.Errorf("%w (%s)", ErrGatewayIssue, resp.HTTPResponse.Status)
	default:
		return nil, fmt.Errorf("unknown %s response: %v",
			resp.HTTPResponse.Status, string(resp.Body))
	}
}

func (c Client) GetNotionDatabasesByTitle(
	ctx context.Context, title string,
) (Databases, error) {
	if title == "" {
		return Databases{}, nil
	}

	resp, err := c.Search(ctx, SearchJSONRequestBody{
		Query: &title,
		Filter: &SearchFilter{
			Value:    SearchFilterValueDatabase,
			Property: SearchFilterPropertyObject,
		},
		// StartCursor: &"", // TODO pagination
		PageSize: &maxPageSizeInt,
	})
	if err != nil {
		return nil, fmt.Errorf("getting page by title for %s: %w", title, err)
	}

	switch resp.StatusCode() {
	case http.StatusOK: // ok
		dbs := make(Databases, len(resp.JSON200.Results))
		for i, res := range resp.JSON200.Results {
			dbs[i] = *res.Database
		}

		return dbs, nil
	case http.StatusBadRequest:
		return nil, resp.JSON400
	case http.StatusNotFound:
		return nil, resp.JSON404
	case http.StatusBadGateway, http.StatusGatewayTimeout:
		return nil, fmt.Errorf("%w (%s)", ErrGatewayIssue, resp.HTTPResponse.Status)
	default:
		return nil, fmt.Errorf("unknown %s response: %v",
			resp.HTTPResponse.Status, string(resp.Body))
	}
}

func (c Client) AppendBlocksToPage(ctx context.Context, pageID Id, blocks ...Block) (Blocks, error) {
	pageUUID := UUID(pageID)

	for i, b := range blocks {
		b.Parent = Parent{
			Type:   ParentTypePageId,
			PageId: &pageUUID,
		}

		if b.Id == "" {
			b.Id = UUID(uuid.NewString())
		}

		if err := b.Validate(); err != nil {
			return nil, fmt.Errorf("validating block %d to be appended: %w", i, err)
		}

		blocks[i] = b
	}

	resp, err := c.AppendBlocks(ctx, pageID, AppendBlocksJSONRequestBody{Children: blocks})
	if err != nil {
		return nil, fmt.Errorf("appending blocks to page %s: %w", pageID, err)
	}

	switch resp.StatusCode() {
	case http.StatusOK: // ok
		return resp.JSON200.Results, nil
	case http.StatusBadRequest:
		return nil, resp.JSON400
	case http.StatusNotFound:
		return nil, resp.JSON404
	case http.StatusTooManyRequests:
		return nil, resp.JSON429
	case http.StatusBadGateway, http.StatusGatewayTimeout:
		return nil, fmt.Errorf("%w (%s)", ErrGatewayIssue, resp.HTTPResponse.Status)
	default:
		return nil, fmt.Errorf("unknown %s response: %v",
			resp.HTTPResponse.Status, string(resp.Body))
	}
}
