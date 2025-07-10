package contracts

type Application interface {
	Bind(key string, factory any)
	Singleton(key string, factory any)
	Make(key string) any

	AddProvider(provider ServiceProvider)
	Boot()
}

type ServiceProvider interface {
	Register(app Application)
	Boot(app Application)
}
