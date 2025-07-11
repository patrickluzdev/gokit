package gokit

import (
	"os"
)

type ConfigProvider struct{}

func (p *ConfigProvider) Register(app Application) {
	app.Singleton(ConfigBinding, func() any {
		envPath := ".env"
		if customPath := os.Getenv("CONFIG_PATH"); customPath != "" {
			envPath = customPath
		}
		return NewConfig(envPath)
	})
}

func (p *ConfigProvider) Boot(app Application) {}

type RouterProvider struct{}

func (p *RouterProvider) Register(app Application) {
	app.Singleton(RouterBinding, func() any {
		return NewRouter()
	})
}

func (p *RouterProvider) Boot(app Application) {}

type DatabaseProvider struct{}

func (p *DatabaseProvider) Register(app Application) {
	app.Singleton(DatabaseBinding, func() any {
		config := app.Config()
		db := NewDB(config)
		return db
	})
}

func (p *DatabaseProvider) Boot(app Application) {
	app.Make(DatabaseBinding)
}
