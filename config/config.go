package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config represents application configuration
type Config struct {
	values map[string]string
}

// NewConfig creates a new Config instance
func NewConfig() *Config {
	return &Config{
		values: make(map[string]string),
	}
}

// LoadFromEnv loads configuration from environment variables
func (c *Config) LoadFromEnv() {
	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		if len(pair) == 2 {
			c.values[pair[0]] = pair[1]
		}
	}
}

// Set sets a configuration value
func (c *Config) Set(key, value string) {
	c.values[key] = value
}

// Get gets a string configuration value
func (c *Config) Get(key string) string {
	return c.values[key]
}

// GetWithDefault gets a string configuration value with default
func (c *Config) GetWithDefault(key, defaultValue string) string {
	if value, exists := c.values[key]; exists && value != "" {
		return value
	}
	return defaultValue
}

// GetInt gets an integer configuration value
func (c *Config) GetInt(key string) (int, error) {
	value, exists := c.values[key]
	if !exists {
		return 0, fmt.Errorf("configuration key '%s' not found", key)
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("failed to parse '%s' as int: %w", key, err)
	}

	return intValue, nil
}

// GetIntWithDefault gets an integer configuration value with default
func (c *Config) GetIntWithDefault(key string, defaultValue int) int {
	intValue, err := c.GetInt(key)
	if err != nil {
		return defaultValue
	}
	return intValue
}

// GetBool gets a boolean configuration value
func (c *Config) GetBool(key string) (bool, error) {
	value, exists := c.values[key]
	if !exists {
		return false, fmt.Errorf("configuration key '%s' not found", key)
	}

	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf("failed to parse '%s' as bool: %w", key, err)
	}

	return boolValue, nil
}

// GetBoolWithDefault gets a boolean configuration value with default
func (c *Config) GetBoolWithDefault(key string, defaultValue bool) bool {
	boolValue, err := c.GetBool(key)
	if err != nil {
		return defaultValue
	}
	return boolValue
}

// GetDuration gets a duration configuration value
func (c *Config) GetDuration(key string) (time.Duration, error) {
	value, exists := c.values[key]
	if !exists {
		return 0, fmt.Errorf("configuration key '%s' not found", key)
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("failed to parse '%s' as duration: %w", key, err)
	}

	return duration, nil
}

// GetDurationWithDefault gets a duration configuration value with default
func (c *Config) GetDurationWithDefault(key string, defaultValue time.Duration) time.Duration {
	duration, err := c.GetDuration(key)
	if err != nil {
		return defaultValue
	}
	return duration
}

// GetStringSlice gets a string slice configuration value (comma-separated)
func (c *Config) GetStringSlice(key string) []string {
	value, exists := c.values[key]
	if !exists || value == "" {
		return []string{}
	}

	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}

// GetStringSliceWithDefault gets a string slice configuration value with default
func (c *Config) GetStringSliceWithDefault(key string, defaultValue []string) []string {
	slice := c.GetStringSlice(key)
	if len(slice) == 0 {
		return defaultValue
	}
	return slice
}

// GetRequired gets a required configuration value, panics if not found
func (c *Config) GetRequired(key string) string {
	value, exists := c.values[key]
	if !exists || value == "" {
		panic(fmt.Sprintf("required configuration key '%s' not found or empty", key))
	}
	return value
}

// GetRequiredInt gets a required integer configuration value, panics if not found or invalid
func (c *Config) GetRequiredInt(key string) int {
	intValue, err := c.GetInt(key)
	if err != nil {
		panic(fmt.Sprintf("required configuration key '%s': %v", key, err))
	}
	return intValue
}

// GetRequiredBool gets a required boolean configuration value, panics if not found or invalid
func (c *Config) GetRequiredBool(key string) bool {
	boolValue, err := c.GetBool(key)
	if err != nil {
		panic(fmt.Sprintf("required configuration key '%s': %v", key, err))
	}
	return boolValue
}

// GetRequiredDuration gets a required duration configuration value, panics if not found or invalid
func (c *Config) GetRequiredDuration(key string) time.Duration {
	duration, err := c.GetDuration(key)
	if err != nil {
		panic(fmt.Sprintf("required configuration key '%s': %v", key, err))
	}
	return duration
}

// Exists checks if a configuration key exists
func (c *Config) Exists(key string) bool {
	_, exists := c.values[key]
	return exists
}

// Keys returns all configuration keys
func (c *Config) Keys() []string {
	keys := make([]string, 0, len(c.values))
	for key := range c.values {
		keys = append(keys, key)
	}
	return keys
}

// Validate validates that all required keys are present
func (c *Config) Validate(requiredKeys []string) error {
	missing := []string{}
	for _, key := range requiredKeys {
		if !c.Exists(key) || c.Get(key) == "" {
			missing = append(missing, key)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required configuration keys: %s", strings.Join(missing, ", "))
	}

	return nil
}
