package config

import (
	"os"
	"testing"
	"time"
)

func TestNewConfig(t *testing.T) {
	config := NewConfig()
	if config == nil {
		t.Fatal("NewConfig() returned nil")
	}

	if config.values == nil {
		t.Error("NewConfig() values map is nil")
	}
}

func TestConfig_LoadFromEnv(t *testing.T) {
	// Set test environment variables
	os.Setenv("TEST_KEY", "test_value")
	os.Setenv("TEST_INT", "42")
	defer func() {
		os.Unsetenv("TEST_KEY")
		os.Unsetenv("TEST_INT")
	}()

	config := NewConfig()
	config.LoadFromEnv()

	if config.Get("TEST_KEY") != "test_value" {
		t.Errorf("LoadFromEnv() TEST_KEY = %v, want test_value", config.Get("TEST_KEY"))
	}

	if config.Get("TEST_INT") != "42" {
		t.Errorf("LoadFromEnv() TEST_INT = %v, want 42", config.Get("TEST_INT"))
	}
}

func TestConfig_SetAndGet(t *testing.T) {
	config := NewConfig()
	config.Set("key", "value")

	if config.Get("key") != "value" {
		t.Errorf("Get() = %v, want value", config.Get("key"))
	}

	if config.Get("nonexistent") != "" {
		t.Errorf("Get() nonexistent key = %v, want empty string", config.Get("nonexistent"))
	}
}

func TestConfig_GetWithDefault(t *testing.T) {
	config := NewConfig()
	config.Set("existing", "value")

	// Existing key
	if config.GetWithDefault("existing", "default") != "value" {
		t.Errorf("GetWithDefault() existing = %v, want value", config.GetWithDefault("existing", "default"))
	}

	// Non-existing key
	if config.GetWithDefault("nonexistent", "default") != "default" {
		t.Errorf("GetWithDefault() nonexistent = %v, want default", config.GetWithDefault("nonexistent", "default"))
	}

	// Empty value
	config.Set("empty", "")
	if config.GetWithDefault("empty", "default") != "default" {
		t.Errorf("GetWithDefault() empty = %v, want default", config.GetWithDefault("empty", "default"))
	}
}

func TestConfig_GetInt(t *testing.T) {
	config := NewConfig()
	config.Set("valid_int", "42")
	config.Set("invalid_int", "not_a_number")

	// Valid integer
	value, err := config.GetInt("valid_int")
	if err != nil {
		t.Errorf("GetInt() valid error = %v", err)
	}
	if value != 42 {
		t.Errorf("GetInt() valid = %v, want 42", value)
	}

	// Invalid integer
	_, err = config.GetInt("invalid_int")
	if err == nil {
		t.Error("GetInt() invalid should return error")
	}

	// Non-existing key
	_, err = config.GetInt("nonexistent")
	if err == nil {
		t.Error("GetInt() nonexistent should return error")
	}
}

func TestConfig_GetIntWithDefault(t *testing.T) {
	config := NewConfig()
	config.Set("valid_int", "42")
	config.Set("invalid_int", "not_a_number")

	// Valid integer
	if config.GetIntWithDefault("valid_int", 100) != 42 {
		t.Errorf("GetIntWithDefault() valid = %v, want 42", config.GetIntWithDefault("valid_int", 100))
	}

	// Invalid integer - should return default
	if config.GetIntWithDefault("invalid_int", 100) != 100 {
		t.Errorf("GetIntWithDefault() invalid = %v, want 100", config.GetIntWithDefault("invalid_int", 100))
	}

	// Non-existing key - should return default
	if config.GetIntWithDefault("nonexistent", 100) != 100 {
		t.Errorf("GetIntWithDefault() nonexistent = %v, want 100", config.GetIntWithDefault("nonexistent", 100))
	}
}

func TestConfig_GetBool(t *testing.T) {
	config := NewConfig()
	config.Set("true_bool", "true")
	config.Set("false_bool", "false")
	config.Set("invalid_bool", "not_a_bool")

	// Valid true
	value, err := config.GetBool("true_bool")
	if err != nil {
		t.Errorf("GetBool() true error = %v", err)
	}
	if !value {
		t.Errorf("GetBool() true = %v, want true", value)
	}

	// Valid false
	value, err = config.GetBool("false_bool")
	if err != nil {
		t.Errorf("GetBool() false error = %v", err)
	}
	if value {
		t.Errorf("GetBool() false = %v, want false", value)
	}

	// Invalid boolean
	_, err = config.GetBool("invalid_bool")
	if err == nil {
		t.Error("GetBool() invalid should return error")
	}
}

func TestConfig_GetDuration(t *testing.T) {
	config := NewConfig()
	config.Set("valid_duration", "30s")
	config.Set("invalid_duration", "not_a_duration")

	// Valid duration
	value, err := config.GetDuration("valid_duration")
	if err != nil {
		t.Errorf("GetDuration() valid error = %v", err)
	}
	if value != 30*time.Second {
		t.Errorf("GetDuration() valid = %v, want 30s", value)
	}

	// Invalid duration
	_, err = config.GetDuration("invalid_duration")
	if err == nil {
		t.Error("GetDuration() invalid should return error")
	}
}

func TestConfig_GetStringSlice(t *testing.T) {
	config := NewConfig()
	config.Set("comma_separated", "item1,item2,item3")
	config.Set("with_spaces", " item1 , item2 , item3 ")
	config.Set("empty_string", "")
	config.Set("single_item", "item1")

	// Comma separated
	slice := config.GetStringSlice("comma_separated")
	expected := []string{"item1", "item2", "item3"}
	if !stringSlicesEqual(slice, expected) {
		t.Errorf("GetStringSlice() comma_separated = %v, want %v", slice, expected)
	}

	// With spaces
	slice = config.GetStringSlice("with_spaces")
	if !stringSlicesEqual(slice, expected) {
		t.Errorf("GetStringSlice() with_spaces = %v, want %v", slice, expected)
	}

	// Empty string
	slice = config.GetStringSlice("empty_string")
	if len(slice) != 0 {
		t.Errorf("GetStringSlice() empty_string = %v, want empty slice", slice)
	}

	// Single item
	slice = config.GetStringSlice("single_item")
	if len(slice) != 1 || slice[0] != "item1" {
		t.Errorf("GetStringSlice() single_item = %v, want [item1]", slice)
	}

	// Non-existing key
	slice = config.GetStringSlice("nonexistent")
	if len(slice) != 0 {
		t.Errorf("GetStringSlice() nonexistent = %v, want empty slice", slice)
	}
}

func TestConfig_GetRequired(t *testing.T) {
	config := NewConfig()
	config.Set("existing", "value")
	config.Set("empty", "")

	// Existing key
	value := config.GetRequired("existing")
	if value != "value" {
		t.Errorf("GetRequired() existing = %v, want value", value)
	}

	// Test panic for non-existing key
	defer func() {
		if r := recover(); r == nil {
			t.Error("GetRequired() nonexistent should panic")
		}
	}()
	config.GetRequired("nonexistent")
}

func TestConfig_Exists(t *testing.T) {
	config := NewConfig()
	config.Set("existing", "value")

	if !config.Exists("existing") {
		t.Error("Exists() existing should return true")
	}

	if config.Exists("nonexistent") {
		t.Error("Exists() nonexistent should return false")
	}
}

func TestConfig_Keys(t *testing.T) {
	config := NewConfig()
	config.Set("key1", "value1")
	config.Set("key2", "value2")

	keys := config.Keys()
	if len(keys) != 2 {
		t.Errorf("Keys() length = %v, want 2", len(keys))
	}

	keyMap := make(map[string]bool)
	for _, key := range keys {
		keyMap[key] = true
	}

	if !keyMap["key1"] || !keyMap["key2"] {
		t.Errorf("Keys() = %v, want [key1, key2]", keys)
	}
}

func TestConfig_Validate(t *testing.T) {
	config := NewConfig()
	config.Set("key1", "value1")
	config.Set("key2", "value2")
	config.Set("empty", "")

	// Valid configuration
	err := config.Validate([]string{"key1", "key2"})
	if err != nil {
		t.Errorf("Validate() valid config error = %v", err)
	}

	// Missing required key
	err = config.Validate([]string{"key1", "nonexistent"})
	if err == nil {
		t.Error("Validate() missing key should return error")
	}

	// Empty required key
	err = config.Validate([]string{"key1", "empty"})
	if err == nil {
		t.Error("Validate() empty key should return error")
	}
}

// Helper function to compare string slices
func stringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
