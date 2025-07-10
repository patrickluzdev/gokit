package gokit

import (
	"github.com/patrickluzdev/gokit/contracts"
	"github.com/patrickluzdev/gokit/facades"
	"github.com/patrickluzdev/gokit/providers"
)

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
	facades.SetApp(app)

	return app
}

func (a *Application) autoRegisterProviders() {
	a.providers = append(a.providers, &providers.ConfigProvider{})
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

func (a *Application) Config() contracts.Config {
	return a.Make("config").(contracts.Config)
}
