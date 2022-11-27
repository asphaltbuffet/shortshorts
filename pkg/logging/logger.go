// Package logging provides a logging singleton.
package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

// GetLogger returns the logger.
func GetLogger() *zap.Logger {
	if logger == nil {
		_ = setupLogger()
	}

	return logger
}

// Start starts the logger.
func Start() error {
	return setupLogger()
}

func setupLogger() error {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	l, _ := cfg.Build()

	logger = l

	return nil
}

// Shutdown shuts down the logger.
func Shutdown() error {
	err := logger.Sync()
	return err
}
