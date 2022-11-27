// Package servicemanager provides a way to manage services.
// Pulled from https://github.com/gbeletti/service-golang/tree/main/servicemanager.
package servicemanager

import (
	"errors"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/asphaltbuffet/shortshorts/pkg/logging"
)

var (
	// ErrServiceCanceled is the error returned when the service is canceled.
	ErrServiceCanceled = errors.New("service canceled")
	logger             *zap.Logger
)

// WaitShutdown waits until is going to die.
func WaitShutdown(shutdown func()) {
	logger = logging.GetLogger()

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	s := <-sigc
	logger.Info("shutdown signal received, canceling everything", zap.String("signal", s.String()))
	shutdown()
}
