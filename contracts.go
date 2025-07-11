package gokit

import "net/http"

const (
	RouterBinding = "gokit.router"
	ConfigBinding = "gokit.config"
)

type Config interface {
	Get(key string) string
	GetWithDefault(key, defaultValue string) string
	Set(key, value string)
	All() map[string]string
	GetInt(key string) int
	GetBool(key string) bool
}

type Application interface {
	Bind(key string, factory any)
	Singleton(key string, factory any)
	Make(key string) any

	AddProvider(provider ServiceProvider)

	Config() Config
	Router() Router
}

type ServiceProvider interface {
	Register(app Application)
	Boot(app Application)
}

type Router interface {
	GET(path string, handler HandlerFunc)
	POST(path string, handler HandlerFunc)
	PUT(path string, handler HandlerFunc)
	DELETE(path string, handler HandlerFunc)
	PATCH(path string, handler HandlerFunc)

	Group(prefix string, fn func(Router))

	Use(middleware ...MiddlewareFunc)

	Listen(addr string) error
}

type HandlerFunc func(Context)
type MiddlewareFunc func(Context)

type Context interface {
	Param(key string) string
	Query(key string) string
	Header(key string) string

	ParseJSON(v interface{}) error
	ParseString() string
	Body() []byte

	JSON(code int, obj any)
	String(code int, format string, values ...any)
	Data(code int, data []byte)

	Request() *http.Request
	Writer() http.ResponseWriter
}

type Response interface {
	Status() int
	Body() []byte
	Headers() map[string]string
}
