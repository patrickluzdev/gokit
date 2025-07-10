package facades

import "github.com/patrickluzdev/gokit/contracts"

var app contracts.Application

func SetApp(a contracts.Application) {
	app = a
}

func App() contracts.Application {
	if app == nil {
		panic("Application not set in facades")
	}
	return app
}

func Config() contracts.Config {
	return App().Config()
}
