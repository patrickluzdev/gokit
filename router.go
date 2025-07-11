package gokit

import (
	"fmt"
	"net/http"
)

type router struct {
	mux        *http.ServeMux
	routes     []Route
	middleware []MiddlewareFunc
	prefix     string
}

type Route struct {
	method  string
	path    string
	handler HandlerFunc
}

func NewRouter() Router {
	return &router{
		mux:        http.NewServeMux(),
		routes:     make([]Route, 0),
		middleware: make([]MiddlewareFunc, 0),
		prefix:     "",
	}
}

func (r *router) DELETE(path string, handler HandlerFunc) {
	fullPath := r.prefix + path
	r.routes = append(r.routes, Route{method: "DELETE", path: fullPath, handler: handler})
	r.mux.HandleFunc("DELETE "+fullPath, func(w http.ResponseWriter, req *http.Request) {
		ctx := NewCtx(w, req)
		r.applyMiddleware(ctx)
		handler(ctx)
	})
}

func (r *router) GET(path string, handler HandlerFunc) {
	fullPath := r.prefix + path
	r.routes = append(r.routes, Route{method: "GET", path: fullPath, handler: handler})
	r.mux.HandleFunc("GET "+fullPath, func(w http.ResponseWriter, req *http.Request) {
		ctx := NewCtx(w, req)
		r.applyMiddleware(ctx)
		handler(ctx)
	})
}

func (r *router) Group(prefix string, fn func(Router)) {
	groupRouter := &router{
		mux:        r.mux,
		routes:     make([]Route, 0),
		middleware: make([]MiddlewareFunc, len(r.middleware)),
		prefix:     r.prefix + prefix,
	}
	copy(groupRouter.middleware, r.middleware)
	fn(groupRouter)
	r.routes = append(r.routes, groupRouter.routes...)
}

func (r *router) Listen(addr string) error {
	fmt.Printf("🚀 Server starting on %s\n", addr)
	return http.ListenAndServe(addr, r.mux)
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

func (r *router) applyMiddleware(ctx Context) {
	for _, mw := range r.middleware {
		mw(ctx)
	}
}

func (r *router) PATCH(path string, handler HandlerFunc) {
	fullPath := r.prefix + path
	r.routes = append(r.routes, Route{method: "PATCH", path: fullPath, handler: handler})
	r.mux.HandleFunc("PATCH "+fullPath, func(w http.ResponseWriter, req *http.Request) {
		ctx := NewCtx(w, req)
		r.applyMiddleware(ctx)
		handler(ctx)
	})
}

func (r *router) POST(path string, handler HandlerFunc) {
	fullPath := r.prefix + path
	r.routes = append(r.routes, Route{method: "POST", path: fullPath, handler: handler})
	r.mux.HandleFunc("POST "+fullPath, func(w http.ResponseWriter, req *http.Request) {
		ctx := NewCtx(w, req)
		r.applyMiddleware(ctx)
		handler(ctx)
	})
}

func (r *router) PUT(path string, handler HandlerFunc) {
	fullPath := r.prefix + path
	r.routes = append(r.routes, Route{method: "PUT", path: fullPath, handler: handler})
	r.mux.HandleFunc("PUT "+fullPath, func(w http.ResponseWriter, req *http.Request) {
		ctx := NewCtx(w, req)
		r.applyMiddleware(ctx)
		handler(ctx)
	})
}

func (r *router) Use(middleware ...MiddlewareFunc) {
	r.middleware = append(r.middleware, middleware...)
}

func (r *router) Routes() []Route {
	return r.routes
}

func (route *Route) Method() string {
	return route.method
}

func (route *Route) Path() string {
	return route.path
}

func (route *Route) Handler() HandlerFunc {
	return route.handler
}
