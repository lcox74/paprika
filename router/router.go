package router

import (
	"context"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type RouterOption func(*Router)
type Router struct {
	ctx context.Context

	// The default page to display when there are no pages in the history.
	// If this is nil, then the router will not display anything.
	defaultPage Page

	// The number of pages to keep in the history. If this is 0, then the
	// router will not keep any history. If the history is full, then the
	// oldest page will be removed.
	history     uint
	pageHistory []Page
}

const (
	DefaultHistory uint = 10
)

// NewRouter creates a new router with the given options.
// The default page is nil, and the history is 10.
//
// Example:
//
//	r := router.NewRouter(
//		router.WithHistory(15),
//		router.WithDefaultPage(&ErrPage{}),
//	)
//
// This will create a router with a history of 15 pages, and the default
// page will be the ErrPage.
func NewRouter(options ...RouterOption) *Router {
	// Set Defaults
	r := &Router{
		ctx:         context.Background(),
		defaultPage: nil,
		history:     DefaultHistory,
	}

	// Apply options
	for _, option := range options {
		option(r)
	}

	return r
}

// WithHistory sets the number of pages to keep in the history. If this is 0,
// then the router will not keep any history. If the history is full, then the
// oldest page will be removed.
//
// Example:
//
//	r := router.NewRouter(
//		router.WithHistory(15),
//	)
//
// This will create a router with a history of 15 pages.
func WithHistory(recordedHistory uint) RouterOption {
	return func(r *Router) {
		r.history = recordedHistory
	}
}

// WithContext sets the context for the router. If this is not set, then the
// context will be context.Background().
//
// Example:
//
//	ctx, cancel := context.WithCancel(context.Background())
//	r := router.NewRouter(
//		router.WithContext(ctx),
//	)
//
// This will create a router with a context that can be cancelled outside of
// the router.
func WithContext(ctx context.Context) RouterOption {
	return func(r *Router) {
		r.ctx = ctx
	}
}

// WithDefaultPage sets the default page to display when there are no pages in
// the history. If this is nil, then the router will not display anything.
//
// Example:
//
//	r := router.NewRouter(
//		router.WithDefaultPage(&ErrPage{}),
//	)
//
// This will create a router with the default page set to the ErrPage.
func WithDefaultPage(page Page) RouterOption {
	return func(r *Router) {
		r.defaultPage = page
	}
}

// Run starts the router. This will loop until the window is closed or the
// context is cancelled. The loop will check the current page and if it should
// display the default page. If there is no default page we will continue to
// the next iteration to see if there will be a page to display. The loop will
// update the state of the current page and then draw the current page's state.
func (r *Router) Run() {
	var currentPage Page

	// Keep looping until the window is closed or the context is cancelled.
	for !rl.WindowShouldClose() {
		select {
		case <-r.ctx.Done():
			return
		default:
			// Check the current page and if it should display the default page.
			// If there is no default page we will continue to the next
			// iteration to see if there will be a page to display.
			currentPage = r.current()
			if currentPage == nil {
				if r.defaultPage != nil {
					r.Push(r.defaultPage)
					currentPage = r.defaultPage
				} else {
					continue
				}
			}

			// Update the state of the current page.
			updatePage(r, currentPage)

			// Draw the current page's state.
			rl.BeginDrawing()
			drawPage(r, currentPage)
			rl.EndDrawing()
		}
	}
}

// Sets a value on the router's context. This is useful for sharing data
// between pages.
func (r *Router) CtxSetValue(key string, value interface{}) {
	r.ctx = context.WithValue(r.ctx, key, value)
}

// Fetches a value from the router's context. This is useful for sharing data
// between pages.
func (r *Router) CtxValue(key string) interface{} {
	return r.ctx.Value(key)
}

// Push adds a page to the history. If the history is full, then the oldest
// page will be removed.
func (r *Router) Push(page Page) {
	// If the page history is full, then remove the oldest page.
	if len(r.pageHistory) >= int(r.history) {
		r.pageHistory = r.pageHistory[1:]
	}

	// Unmount the last page to clean up resources.
	if len(r.pageHistory) != 0 {
		lastPage := r.pageHistory[len(r.pageHistory)-1]
		unmountPage(r, lastPage)
	}

	// Add new page to history and mount it.
	r.pageHistory = append(r.pageHistory, page)
	mountPage(r, page)
}

// Pop removes the last page from the history and returns it. If there are no
// pages in the history, then nil is returned.
func (r *Router) Pop() Page {
	if len(r.pageHistory) == 0 {
		return nil
	}

	// Unmount the last page and mount the previous page.
	lastPage := r.pageHistory[len(r.pageHistory)-1]
	unmountPage(r, lastPage)

	// Remove last page from history and mount the previous page. If there is no
	// previous page, then the router will be empty. If a default page is set
	// then it will be mounted.
	r.pageHistory = r.pageHistory[:len(r.pageHistory)-1]
	if len(r.pageHistory) == 0 {
		return nil
	}
	mountPage(r, r.pageHistory[len(r.pageHistory)-1])

	return lastPage
}

func (r *Router) current() Page {
	if len(r.pageHistory) == 0 {
		return nil
	}
	return r.pageHistory[len(r.pageHistory)-1]
}
