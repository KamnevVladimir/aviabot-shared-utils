package providers

import (
	"regexp"
	"strings"
	"testing"
)

func TestUUIDGenerator_Generate(t *testing.T) {
	generator := NewUUIDGenerator()
	
	id := generator.Generate()
	if id == "" {
		t.Error("UUIDGenerator.Generate() returned empty string")
	}
	
	// Check UUID-like format (8-4-4-4-12 hex characters)
	uuidRegex := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	if !uuidRegex.MatchString(id) {
		t.Errorf("UUIDGenerator.Generate() = %v, does not match UUID format", id)
	}
}

func TestUUIDGenerator_GenerateUnique(t *testing.T) {
	generator := NewUUIDGenerator()
	
	ids := make(map[string]bool)
	for i := 0; i < 1000; i++ {
		id := generator.Generate()
		if ids[id] {
			t.Errorf("UUIDGenerator.Generate() generated duplicate ID: %v", id)
		}
		ids[id] = true
	}
}

func TestPrefixedIDGenerator_Generate(t *testing.T) {
	prefix := "test"
	generator := NewPrefixedIDGenerator(prefix)
	
	id := generator.Generate()
	if !strings.HasPrefix(id, prefix+"-") {
		t.Errorf("PrefixedIDGenerator.Generate() = %v, expected prefix %v", id, prefix)
	}
	
	// Check that the part after prefix is UUID-like
	parts := strings.Split(id, "-")
	if len(parts) < 6 { // prefix + 5 UUID parts
		t.Errorf("PrefixedIDGenerator.Generate() = %v, invalid format", id)
	}
}

func TestSimpleIDGenerator_Generate(t *testing.T) {
	prefix := "simple"
	generator := NewSimpleIDGenerator(prefix).(*SimpleIDGenerator)
	
	// Test sequential generation
	id1 := generator.Generate()
	id2 := generator.Generate()
	id3 := generator.Generate()
	
	expected := []string{"simple-1", "simple-2", "simple-3"}
	actual := []string{id1, id2, id3}
	
	for i, exp := range expected {
		if actual[i] != exp {
			t.Errorf("SimpleIDGenerator.Generate() call %d = %v, want %v", i+1, actual[i], exp)
		}
	}
}

func TestHexIDGenerator_Generate(t *testing.T) {
	length := 16
	generator := NewHexIDGenerator(length)
	
	id := generator.Generate()
	if len(id) != length {
		t.Errorf("HexIDGenerator.Generate() length = %d, want %d", len(id), length)
	}
	
	// Check that it's hexadecimal
	hexRegex := regexp.MustCompile(`^[0-9a-f]+$`)
	if !hexRegex.MatchString(id) {
		t.Errorf("HexIDGenerator.Generate() = %v, not hexadecimal", id)
	}
}

func TestHexIDGenerator_DefaultLength(t *testing.T) {
	generator := NewHexIDGenerator(0) // Should use default length
	
	id := generator.Generate()
	if len(id) != 16 { // Default length
		t.Errorf("HexIDGenerator.Generate() with 0 length = %d chars, want 16", len(id))
	}
}

func TestHexIDGenerator_GenerateUnique(t *testing.T) {
	generator := NewHexIDGenerator(32)
	
	ids := make(map[string]bool)
	for i := 0; i < 100; i++ {
		id := generator.Generate()
		if ids[id] {
			t.Errorf("HexIDGenerator.Generate() generated duplicate ID: %v", id)
		}
		ids[id] = true
	}
}

func TestAllGenerators_Interface(t *testing.T) {
	generators := []interface{}{
		NewUUIDGenerator(),
		NewPrefixedIDGenerator("test"),
		NewSimpleIDGenerator("test"),
		NewHexIDGenerator(16),
	}
	
	for i, gen := range generators {
		if gen == nil {
			t.Errorf("Generator %d returned nil", i)
		}
	}
}