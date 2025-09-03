package providers

import (
	"github.com/KamnevVladimir/aviabot-shared-core/domain/interfaces"
	"time"
)

// SystemTimeProvider provides current system time
type SystemTimeProvider struct{}

// NewSystemTimeProvider creates a new SystemTimeProvider
func NewSystemTimeProvider() interfaces.TimeProvider {
	return &SystemTimeProvider{}
}

// Now returns the current system time
func (p *SystemTimeProvider) Now() time.Time {
	return time.Now()
}

// FixedTimeProvider provides a fixed time for testing
type FixedTimeProvider struct {
	fixedTime time.Time
}

// NewFixedTimeProvider creates a new FixedTimeProvider with specified time
func NewFixedTimeProvider(fixedTime time.Time) interfaces.TimeProvider {
	return &FixedTimeProvider{
		fixedTime: fixedTime,
	}
}

// Now returns the fixed time
func (p *FixedTimeProvider) Now() time.Time {
	return p.fixedTime
}

// SetTime updates the fixed time
func (p *FixedTimeProvider) SetTime(newTime time.Time) {
	p.fixedTime = newTime
}
