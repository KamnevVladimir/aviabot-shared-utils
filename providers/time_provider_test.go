package providers

import (
	"testing"
	"time"
)

func TestSystemTimeProvider_Now(t *testing.T) {
	provider := NewSystemTimeProvider()

	before := time.Now()
	now := provider.Now()
	after := time.Now()

	if now.Before(before) || now.After(after) {
		t.Errorf("SystemTimeProvider.Now() returned time outside expected range")
	}
}

func TestSystemTimeProvider_Interface(t *testing.T) {
	provider := NewSystemTimeProvider()
	if provider == nil {
		t.Fatal("NewSystemTimeProvider() returned nil")
	}
}

func TestFixedTimeProvider_Now(t *testing.T) {
	fixedTime := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	provider := NewFixedTimeProvider(fixedTime)

	now := provider.Now()
	if !now.Equal(fixedTime) {
		t.Errorf("FixedTimeProvider.Now() = %v, want %v", now, fixedTime)
	}
}

func TestFixedTimeProvider_SetTime(t *testing.T) {
	initialTime := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	newTime := time.Date(2025, 1, 2, 12, 0, 0, 0, time.UTC)

	provider := NewFixedTimeProvider(initialTime).(*FixedTimeProvider)

	// Check initial time
	if !provider.Now().Equal(initialTime) {
		t.Errorf("Initial time = %v, want %v", provider.Now(), initialTime)
	}

	// Set new time
	provider.SetTime(newTime)

	// Check new time
	if !provider.Now().Equal(newTime) {
		t.Errorf("After SetTime = %v, want %v", provider.Now(), newTime)
	}
}

func TestFixedTimeProvider_Interface(t *testing.T) {
	fixedTime := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	provider := NewFixedTimeProvider(fixedTime)
	if provider == nil {
		t.Fatal("NewFixedTimeProvider() returned nil")
	}
}
