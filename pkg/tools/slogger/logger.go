package slogger

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

// Default log level
var DefaultLogLevel slog.Level = slog.LevelDebug

// LogOption defines a function type that modifies the LogConfig
type LogOption func(*LogConfig)

// LogConfig stores the configuration for logging
type LogConfig struct {
	Level         slog.Level
	WithTimestamp bool
	WithContext   bool
}

// NewLogConfig creates a LogConfig with default values
func NewLogConfig() *LogConfig {
	return &LogConfig{
		Level:         DefaultLogLevel,
		WithTimestamp: true, // Default to include timestamp
		WithContext:   true, // Default to include context
	}
}

// WithLevel sets the log level
func WithLevel(level slog.Level) LogOption {
	return func(cfg *LogConfig) {
		cfg.Level = level
	}
}

// WithTimestamp enables or disables timestamp
func WithTimestamp(enable bool) LogOption {
	return func(cfg *LogConfig) {
		cfg.WithTimestamp = enable
	}
}

// WithContext enables or disables context
func WithContext(enable bool) LogOption {
	return func(cfg *LogConfig) {
		cfg.WithContext = enable
	}
}

// Log is a single function for logging at various levels
// The function expects a string message, a slice of slog.Attr (for attributes like errors),
// and a variadic number of LogOption (for customization such as level, timestamp, etc.)
func Log(ctx context.Context, msg string, attrs []slog.Attr, opts ...LogOption) {
	// Create a default log config
	cfg := NewLogConfig()

	// Apply all options to the log configuration
	for _, opt := range opts {
		opt(cfg)
	}

	// Collect log attributes
	var logAttrs []slog.Attr

	// Include timestamp if enabled
	if cfg.WithTimestamp {
		logAttrs = append(logAttrs, slog.String("timestamp", fmt.Sprintf("%v", time.Now())))
	}

	// Include context if enabled
	if cfg.WithContext && ctx != nil {
		logAttrs = append(logAttrs, slog.String("context", fmt.Sprintf("%v", ctx)))
	}

	// Append the provided attributes (e.g., error information)
	logAttrs = append(logAttrs, attrs...)

	// Log the message at the specified level
	slog.LogAttrs(ctx, cfg.Level, msg, logAttrs...)
}
