package web

import (
	"math/rand"
	"net/http"
	"strconv"

	"context"

	"flamingo/core/flamingo/event"
	"flamingo/core/flamingo/profiler"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type (
	// ContextKey is used for context.WithValue
	ContextKey string

	// Context defines what a controller sees
	Context interface {
		context.Context
		profiler.Profiler
		EventRouter() event.Router

		LoadVars(r *http.Request)
		Form(string) []string
		Form1(string) string
		FormAll() map[string][]string
		Param1(string) string
		ParamAll() map[string]string
		Query(string) []string
		QueryFirst(string) string
		QueryAll() map[string][]string
		Request() *http.Request

		ID() string

		Push(target string, opts *http.PushOptions) error

		Session() *sessions.Session
	}

	// ContextFactory creates a new context
	ContextFactory func(rw http.ResponseWriter, r *http.Request, session *sessions.Session) Context

	ctx struct {
		context.Context
		profiler.Profiler `inject:"private"`
		Eventrouter       event.Router `inject:"private"`

		vars    map[string]string
		request *http.Request
		id      string
		writer  http.ResponseWriter
		pusher  http.Pusher
		session *sessions.Session
	}
)

const (
	CONTEXT ContextKey = "context"
)

// ContextFromRequest returns a ctx enriched by Request Data
func ContextFromRequest(rw http.ResponseWriter, r *http.Request, session *sessions.Session) Context {
	c := new(ctx)
	c.vars = mux.Vars(r)
	c.request = r
	c.id = strconv.Itoa(rand.Int())
	c.writer = rw
	pusher, ok := rw.(http.Pusher)
	if ok {
		c.pusher = pusher
	}

	c.session = session
	c.Context = context.WithValue(r.Context(), "ID", c.id)
	return c
}

func (c *ctx) LoadVars(r *http.Request) {
	c.vars = mux.Vars(r)
}

func (c *ctx) EventRouter() event.Router {
	return c.Eventrouter
}

// Session returns the ctx Session
func (c *ctx) Session() *sessions.Session {
	return c.session
}

// Push pushes Assets via HTTP2 server push
func (c *ctx) Push(target string, opts *http.PushOptions) error {
	if c.pusher != nil {
		return c.pusher.Push(target, opts)
	}
	return nil
}

// ID returns the ctx Id (random Int)
func (c *ctx) ID() string {
	return c.id
}

// Form get POST value
func (c *ctx) Form(n string) []string {
	return c.FormAll()[n]
}

// Form1 get first POST value
func (c *ctx) Form1(n string) string {
	if len(c.Form(n)) < 1 {
		return ""
	}
	return c.Form(n)[0]
}

// FormAll get all POST values
func (c *ctx) FormAll() map[string][]string {
	c.Request().ParseForm()
	return c.Request().Form
}

// Param get querystring param
func (c *ctx) Param1(n string) string {
	return c.vars[n]
}

// Params get all querystring params
func (c *ctx) ParamAll() map[string]string {
	return c.vars
}

// Query looks up Raw Query map for Param
func (c *ctx) Query(n string) []string {
	return c.QueryAll()[n]
}

// QueryFirst  looks up Raw Query map for  First Param
func (c *ctx) QueryFirst(n string) string {
	return c.Query(n)[0]
}

// QueryAll returns a Map of the Raw Query
func (c *ctx) QueryAll() map[string][]string {
	return c.request.URL.Query()
}

// Request get the context's request
func (c *ctx) Request() *http.Request {
	return c.request
}
