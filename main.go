// Package main is the main package for shortshorts.
package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/asphaltbuffet/shortshorts/pkg/mqttwrapper"
)

func main() {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, _ := cfg.Build()

	defer logger.Sync() //nolint:errcheck // don't care about sync errors

	num := 2 // The number of messages to subscribe

	choke, _ := mqttwrapper.Start()

	receiveCount := 0

	for receiveCount < num {
		incoming := <-choke
		logger.Info("received message", zap.String("topic", incoming[0]), zap.String("message", incoming[1]))
		receiveCount++
	}

	mqttwrapper.Close()
}
