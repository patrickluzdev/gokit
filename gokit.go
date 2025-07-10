package gokit

import "github.com/patrickluzdev/gokit/contracts"

type Application struct {
	*Container
	providers []contracts.ServiceProvider
	booted    bool
}

func (a *Application) AddProvider(provider contracts.ServiceProvider) {
	a.providers = append(a.providers, provider)
}

func New() contracts.Application {
	app := &Application{
		Container: NewContainer(),
		providers: make([]contracts.ServiceProvider, 0),
	}

	app.autoRegisterProviders()

	return app
}

func (a *Application) autoRegisterProviders() {
	temp := make([]contracts.ServiceProvider, 0)
	a.providers = append(a.providers, temp...)
}

func (a *Application) Boot() {
	if a.booted {
		return
	}

	for _, provider := range a.providers {
		provider.Register(a)
	}

	for _, provider := range a.providers {
		provider.Boot(a)
	}

	a.booted = true
}
