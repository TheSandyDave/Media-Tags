package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/TheSandyDave/Media-Tags/router"
	"github.com/sirupsen/logrus"
)

const shutdownTimeout = 30 * time.Second

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logrus.SetOutput(os.Stdout)

	logger := logrus.WithContext(ctx)

	serveAddress := "localhost:8080"

	API := router.TaggedMediaAPI{
		Spec: spec,
	}
	router := API.Configure(ctx)

	srv := &http.Server{
		Addr:    serveAddress,
		Handler: router,

		ReadHeaderTimeout: 20 * time.Second,
	}

	// run server in goroutine so graceful shutdowns can be handled

	go func() {
		logger.Infof("starting server on %s", serveAddress)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WithError(err).Error("failed starting server")
		}
	}()
	// listen for a shutdown signal
	<-ctx.Done()
	stop()

	logger.Infof("shutdown signal received, shutting down in: %v", shutdownTimeout)

	//initiate graceful shutdown timer on context to finish current request handling
	timeoutContext, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(timeoutContext); err != nil {
		logger.WithError(err).Error("shutdown forced")
	}

	<-timeoutContext.Done()
}
