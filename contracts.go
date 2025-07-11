package gokit

import (
	"database/sql"
	"net/http"
)

const (
	RouterBinding   = "gokit.router"
	ConfigBinding   = "gokit.config"
	DatabaseBinding = "gokit.database"
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
	DB() Database
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

	ParseJSON(v any) error
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

type Database interface {
	Connection() *sql.DB

	Table(name string) QueryBuilder
	Model(model any) QueryBuilder
	Raw(sql string, args ...any) QueryBuilder

	Begin() (Transaction, error)
	Transaction(fn func(tx Transaction) error) error

	Migrate() error
	Seed() error

	Exec(sql string, args ...any) error
}

type QueryBuilder interface {
	Select(columns ...string) QueryBuilder
	Where(query any, args ...any) QueryBuilder
	WhereIn(column string, values []any) QueryBuilder
	Join(query string, args ...any) QueryBuilder
	OrderBy(column string, direction ...string) QueryBuilder
	GroupBy(columns ...string) QueryBuilder
	Having(query any, args ...any) QueryBuilder
	Limit(limit int) QueryBuilder
	Offset(offset int) QueryBuilder

	Preload(query string, args ...any) QueryBuilder
	Joins(query string, args ...any) QueryBuilder

	Find(dest any) error
	First(dest any) error
	Create(value any) error
	Update(values any) error
	Delete() error
	Count() (int64, error)

	Begin() (QueryBuilder, error)
	Commit() error
	Rollback() error

	ToSQL() (string, []any, error)
}

type Transaction interface {
	QueryBuilder

	Commit() error
	Rollback() error
	SavePoint(name string) error
	RollbackTo(name string) error

	Table(name string) QueryBuilder
	Model(model any) QueryBuilder
}

type Migrator interface {
	Run() error
	Rollback() error
	Fresh() error
	Status() error
	AddMigration(migration Migration)
}

type Migration interface {
	ID() string
	Up() error
	Down() error
}

type SeederManager interface {
	Add(seeder Seeder)
	Run() error
}

type Seeder interface {
	Run() error
}
