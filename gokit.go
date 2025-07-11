package gokit

type App struct {
	*Container
	providers []ServiceProvider
	booted    bool
}

func (a *App) AddProvider(provider ServiceProvider) {
	a.providers = append(a.providers, provider)
}

func New() Application {
	app := &App{
		Container: NewContainer(),
		providers: make([]ServiceProvider, 0),
	}

	app.autoRegisterProviders()
	app.boot()
	return app
}

func (a *App) autoRegisterProviders() {
	a.providers = append(a.providers, &ConfigProvider{}, &RouterProvider{})
}

func (a *App) boot() {
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

func (a *App) Config() Config {
	return a.Make(ConfigBinding).(Config)
}

func (a *App) Router() Router {
	return a.Make(RouterBinding).(Router)
}
