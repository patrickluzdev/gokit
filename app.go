package gokit

import "github.com/patrickluzdev/gokit/contracts"

type Application struct {
	booted bool
}

func New() contracts.Application {
	app := &Application{}

	return app
}

func (a *Application) Boot() {
	a.booted = true
}
