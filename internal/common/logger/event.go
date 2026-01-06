package logger

import (
	"context"
	"time"
)

// Event is the core structure for Canonical Log Lines.
// Instead of many scattered log lines, we emit ONE comprehensive event per request
// containing all context needed for debugging.
type Event struct {
	// Request metadata
	Timestamp time.Time `json:"timestamp"`
	RequestID string    `json:"request_id"`
	TraceID   string    `json:"trace_id,omitempty"`
	SpanID    string    `json:"span_id,omitempty"`

	// Service metadata
	Service     string `json:"service"`
	Version     string `json:"version"`
	Environment string `json:"environment,omitempty"`

	// HTTP metadata
	Method     string            `json:"method"`
	Path       string            `json:"path"`
	Query      string            `json:"query,omitempty"`
	StatusCode int               `json:"status_code"`
	DurationMs int64             `json:"duration_ms"`
	BytesSent  int64             `json:"bytes_sent,omitempty"`
	ClientIP   string            `json:"client_ip,omitempty"`
	UserAgent  string            `json:"user_agent,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`

	// User context - added by auth middleware when authenticated
	User *UserContext `json:"user,omitempty"`

	// Error context - added when errors occur
	Error *ErrorContext `json:"error,omitempty"`

	// Database and external call metrics
	DBQueries     int   `json:"db_queries,omitempty"`
	DBDurationMs  int64 `json:"db_duration_ms,omitempty"`
	ExternalCalls int   `json:"external_calls,omitempty"`
	CacheHit      *bool `json:"cache_hit,omitempty"`

	// Feature flags for debugging rollouts
	FeatureFlags map[string]bool `json:"feature_flags,omitempty"`

	// Business-specific context - add domain-specific data here
	Custom map[string]any `json:"custom,omitempty"`

	// Final outcome
	Outcome string `json:"outcome"` // "success" or "error"
}

// UserContext contains authenticated user information
type UserContext struct {
	ID       string `json:"id"`
	Email    string `json:"email,omitempty"`
	Role     string `json:"role,omitempty"`
	TenantID string `json:"tenant_id,omitempty"`
}

// ErrorContext contains detailed error information
type ErrorContext struct {
	Type      string `json:"type"`
	Code      string `json:"code"`
	Message   string `json:"message"`
	Retriable bool   `json:"retriable"`
	Stack     string `json:"stack,omitempty"`
}

// eventKey is a custom type to avoid context key collisions
type eventKey struct{}

// GetEvent retrieves the Event from context.
// Returns nil if no event is present.
func GetEvent(ctx context.Context) *Event {
	event, ok := ctx.Value(eventKey{}).(*Event)
	if !ok {
		return nil
	}
	return event
}

// WithEvent stores an Event in the context.
func WithEvent(ctx context.Context, event *Event) context.Context {
	return context.WithValue(ctx, eventKey{}, event)
}

// AddUserContext adds user information to the event.
// Call this from auth middleware after successful authentication.
func AddUserContext(ctx context.Context, userID, email string) {
	if event := GetEvent(ctx); event != nil {
		event.User = &UserContext{
			ID:    userID,
			Email: email,
		}
	}
}

// AddUserContextFull adds complete user information to the event.
func AddUserContextFull(ctx context.Context, user UserContext) {
	if event := GetEvent(ctx); event != nil {
		event.User = &user
	}
}

// AddError adds error context to the event.
func AddError(ctx context.Context, errType, code, message string, retriable bool) {
	if event := GetEvent(ctx); event != nil {
		event.Error = &ErrorContext{
			Type:      errType,
			Code:      code,
			Message:   message,
			Retriable: retriable,
		}
		event.Outcome = "error"
	}
}

// AddErrorWithStack adds error context with stack trace.
func AddErrorWithStack(ctx context.Context, errType, code, message, stack string, retriable bool) {
	if event := GetEvent(ctx); event != nil {
		event.Error = &ErrorContext{
			Type:      errType,
			Code:      code,
			Message:   message,
			Retriable: retriable,
			Stack:     stack,
		}
		event.Outcome = "error"
	}
}

// IncrementDBQueries increments the database query counter.
func IncrementDBQueries(ctx context.Context) {
	if event := GetEvent(ctx); event != nil {
		event.DBQueries++
	}
}

// AddDBDuration adds to the total database query duration.
func AddDBDuration(ctx context.Context, durationMs int64) {
	if event := GetEvent(ctx); event != nil {
		event.DBDurationMs += durationMs
	}
}

// IncrementExternalCalls increments the external API call counter.
func IncrementExternalCalls(ctx context.Context) {
	if event := GetEvent(ctx); event != nil {
		event.ExternalCalls++
	}
}

// SetCacheHit sets whether the request used cached data.
func SetCacheHit(ctx context.Context, hit bool) {
	if event := GetEvent(ctx); event != nil {
		event.CacheHit = &hit
	}
}

// SetFeatureFlag sets a feature flag value.
func SetFeatureFlag(ctx context.Context, flag string, enabled bool) {
	if event := GetEvent(ctx); event != nil {
		if event.FeatureFlags == nil {
			event.FeatureFlags = make(map[string]bool)
		}
		event.FeatureFlags[flag] = enabled
	}
}

// SetCustom sets a custom field in the event.
func SetCustom(ctx context.Context, key string, value any) {
	if event := GetEvent(ctx); event != nil {
		if event.Custom == nil {
			event.Custom = make(map[string]any)
		}
		event.Custom[key] = value
	}
}
