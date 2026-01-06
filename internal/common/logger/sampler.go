package logger

import (
	"math/rand"
)

// Sampler determines whether an event should be kept or dropped.
// Philosophy: Always keep errors, slow requests, and VIP users.
// For normal requests, apply probabilistic sampling to reduce costs.
type Sampler struct {
	// BaseRate is the sampling rate for normal successful requests (e.g., 0.05 = 5%)
	BaseRate float64

	// P99ThresholdMs is the latency threshold above which requests are always kept
	P99ThresholdMs int64

	// Enabled controls whether sampling is active (if false, all events are kept)
	Enabled bool
}

// SamplerConfig holds configuration for creating a Sampler
type SamplerConfig struct {
	Enabled        bool
	BaseRate       float64
	P99ThresholdMs int64
}

// NewSampler creates a new Sampler with the given configuration.
func NewSampler(cfg SamplerConfig) *Sampler {
	return &Sampler{
		Enabled:        cfg.Enabled,
		BaseRate:       cfg.BaseRate,
		P99ThresholdMs: cfg.P99ThresholdMs,
	}
}

// DefaultSampler creates a sampler with sensible defaults:
// - 5% base sampling rate
// - 2000ms P99 threshold
func DefaultSampler() *Sampler {
	return &Sampler{
		Enabled:        true,
		BaseRate:       0.05, // 5%
		P99ThresholdMs: 2000, // 2 seconds
	}
}

// ShouldSample determines if an event should be kept based on tail sampling rules.
// Events are ALWAYS kept if:
// 1. Sampling is disabled
// 2. Status code >= 500 (server errors)
// 3. There's an error in the event
// 4. Request duration exceeds P99 threshold
//
// Otherwise, random sampling is applied at the configured BaseRate.
func (s *Sampler) ShouldSample(event *Event) bool {
	// If sampling is disabled, always keep
	if !s.Enabled {
		return true
	}

	// ALWAYS keep server errors (5xx)
	if event.StatusCode >= 500 {
		return true
	}

	// ALWAYS keep if there's an error context
	if event.Error != nil {
		return true
	}

	// ALWAYS keep slow requests (above P99 threshold)
	if event.DurationMs > s.P99ThresholdMs {
		return true
	}

	// ALWAYS keep client errors (4xx) - useful for debugging bad requests
	if event.StatusCode >= 400 && event.StatusCode < 500 {
		return true
	}

	// Random sampling for normal successful requests
	return rand.Float64() < s.BaseRate
}

// Deprecated: TailSampler is deprecated, use Sampler instead.
// This alias is kept for backward compatibility.
type TailSampler = Sampler

// Deprecated: NewTailSampler is deprecated, use NewSampler instead.
func NewTailSampler(cfg SamplerConfig) *Sampler {
	return NewSampler(cfg)
}
