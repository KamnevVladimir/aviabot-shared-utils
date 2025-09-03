package providers

import (
	"github.com/KamnevVladimir/aviabot-shared-core/domain/interfaces"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// UUIDGenerator generates UUID-like identifiers
type UUIDGenerator struct{}

// NewUUIDGenerator creates a new UUIDGenerator
func NewUUIDGenerator() interfaces.IDGenerator {
	return &UUIDGenerator{}
}

// Generate creates a new UUID-like identifier
func (g *UUIDGenerator) Generate() string {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		// Fallback to timestamp-based ID if random fails
		return fmt.Sprintf("id-%d", time.Now().UnixNano())
	}

	// Format as UUID-like string
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		bytes[0:4], bytes[4:6], bytes[6:8], bytes[8:10], bytes[10:16])
}

// PrefixedIDGenerator generates IDs with a specific prefix
type PrefixedIDGenerator struct {
	prefix string
	base   interfaces.IDGenerator
}

// NewPrefixedIDGenerator creates a new PrefixedIDGenerator
func NewPrefixedIDGenerator(prefix string) interfaces.IDGenerator {
	return &PrefixedIDGenerator{
		prefix: prefix,
		base:   NewUUIDGenerator(),
	}
}

// Generate creates a new ID with the specified prefix
func (g *PrefixedIDGenerator) Generate() string {
	baseID := g.base.Generate()
	return fmt.Sprintf("%s-%s", g.prefix, baseID)
}

// SimpleIDGenerator generates simple incremental IDs for testing
type SimpleIDGenerator struct {
	counter int
	prefix  string
}

// NewSimpleIDGenerator creates a new SimpleIDGenerator
func NewSimpleIDGenerator(prefix string) interfaces.IDGenerator {
	return &SimpleIDGenerator{
		counter: 0,
		prefix:  prefix,
	}
}

// Generate creates a simple incremental ID
func (g *SimpleIDGenerator) Generate() string {
	g.counter++
	return fmt.Sprintf("%s-%d", g.prefix, g.counter)
}

// HexIDGenerator generates hexadecimal IDs
type HexIDGenerator struct {
	length int
}

// NewHexIDGenerator creates a new HexIDGenerator with specified length
func NewHexIDGenerator(length int) interfaces.IDGenerator {
	if length <= 0 {
		length = 16 // Default length
	}
	return &HexIDGenerator{length: length}
}

// Generate creates a new hexadecimal ID
func (g *HexIDGenerator) Generate() string {
	bytes := make([]byte, g.length/2)
	_, err := rand.Read(bytes)
	if err != nil {
		// Fallback to timestamp-based hex ID
		return fmt.Sprintf("%x", time.Now().UnixNano())[:g.length]
	}
	return hex.EncodeToString(bytes)[:g.length]
}
