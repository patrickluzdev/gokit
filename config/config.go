package config

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/patrickluzdev/gokit/contracts"
)

type Config struct {
	values map[string]string
	mu     sync.RWMutex
}

func New(envPath string) contracts.Config {
	c := &Config{
		values: make(map[string]string),
	}
	c.loadEnv(envPath)
	return c
}

func (c *Config) loadEnv(path string) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if parts := strings.SplitN(line, "=", 2); len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			value = strings.Trim(value, `"'`)
			c.values[key] = value
		}
	}
}

func (c *Config) Get(key string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if value := os.Getenv(key); value != "" {
		return value
	}

	return c.values[key]
}

func (c *Config) GetWithDefault(key, defaultValue string) string {
	if value := c.Get(key); value != "" {
		return value
	}
	return defaultValue
}

func (c *Config) GetInt(key string) int {
	value := c.Get(key)
	if value == "" {
		return 0
	}

	intVal, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return intVal
}

func (c *Config) GetBool(key string) bool {
	value := strings.ToLower(c.Get(key))
	return value == "true" || value == "1" || value == "yes"
}

func (c *Config) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.values[key] = value
}

func (c *Config) All() map[string]string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make(map[string]string)
	for k, v := range c.values {
		result[k] = v
	}
	return result
}
