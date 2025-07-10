package providers

import (
	"os"

	"github.com/patrickluzdev/gokit/config"
	"github.com/patrickluzdev/gokit/contracts"
)

type ConfigProvider struct{}

func (p *ConfigProvider) Register(app contracts.Application) {
	app.Singleton("config", func() any {
		envPath := ".env"
		if customPath := os.Getenv("CONFIG_PATH"); customPath != "" {
			envPath = customPath
		}
		return config.New(envPath)
	})
}

func (p *ConfigProvider) Boot(app contracts.Application) {}
