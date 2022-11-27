// Package servicemanager provides a way to manage services.
// Pulled from https://github.com/gbeletti/service-golang/tree/main/servicemanager.
package servicemanager

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// ErrServiceCanceled is the error returned when the service is canceled.
var ErrServiceCanceled = errors.New("service canceled")

// WaitShutdown waits until is going to die.
func WaitShutdown(shutdown func()) {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	s := <-sigc
	log.Printf("signal received [%v] canceling everything\n", s)
	shutdown()
}
