package gokit

import (
	"fmt"
	"reflect"
	"sync"
)

type Container struct {
	bindings   map[string]any
	instances  map[string]any
	singletons map[string]bool
	mu         sync.Mutex
}

func NewContainer() *Container {
	return &Container{
		bindings:   make(map[string]any),
		instances:  make(map[string]any),
		singletons: make(map[string]bool),
	}
}

func (c *Container) Make(key string) any {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.singletons[key] {
		if instance, exists := c.instances[key]; exists {
			return instance
		}
	}

	factory, exists := c.bindings[key]
	if !exists {
		panic(fmt.Sprintf("Service not found: %s", key))
	}

	instance := c.callFactory(factory)

	if c.singletons[key] {
		c.instances[key] = instance
	}

	return instance
}

func (c *Container) Bind(key string, factory any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.bindings[key] = factory
	c.singletons[key] = false
}

func (c *Container) Singleton(key string, factory any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.bindings[key] = factory
	c.singletons[key] = true
}

func (c *Container) callFactory(factory any) any {
	factoryValue := reflect.ValueOf(factory)
	if factoryValue.Kind() != reflect.Func {
		return factory
	}

	results := factoryValue.Call(nil)
	if len(results) == 0 {
		panic("Factory function must return a value")
	}

	return results[0].Interface()
}
