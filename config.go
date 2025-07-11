package gokit

import (
	"bufio"
	"maps"
	"os"
	"strconv"
	"strings"
	"sync"
)

type config struct {
	values map[string]string
	mu     sync.RWMutex
}

func NewConfig(envPath string) Config {
	c := &config{
		values: make(map[string]string),
	}
	c.loadEnv(envPath)
	return c
}

func (c *config) loadEnv(path string) {
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

func (c *config) Get(key string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if value := os.Getenv(key); value != "" {
		return value
	}

	return c.values[key]
}

func (c *config) GetWithDefault(key, defaultValue string) string {
	if value := c.Get(key); value != "" {
		return value
	}
	return defaultValue
}

func (c *config) GetInt(key string) int {
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

func (c *config) GetBool(key string) bool {
	value := strings.ToLower(c.Get(key))
	return value == "true" || value == "1" || value == "yes"
}

func (c *config) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.values[key] = value
}

func (c *config) All() map[string]string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make(map[string]string)
	maps.Copy(result, c.values)
	return result
}
