package main

import (
	"context"
	"database/sql"
	"errors"
	"expvar"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/Luqqk/goexrates/internal/jsonlog"
)

// Define an application struct to hold the dependencies for our HTTP handlers, helpers,
// and middleware.
type application struct {
	config config
	logger *jsonlog.Logger
	db     *sql.DB
}

func (app *application) serve() error {
	// Declare a HTTP server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.config.port),
		Handler: app.routes(),
		// Create a new Go log.Logger instance with the log.New() function, passing in
		// our custom Logger as the first parameter. The "" and 0 indicate that the
		// log.Logger instance should not use a prefix or any flags.
		ErrorLog:     log.New(app.logger, "", 0),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	// Create a shutdownError channel. Use this to receive any errors returned
	// by the graceful Shutdown() function.
	shutdownError := make(chan error)

	// Start a background goroutine that catches SIGINT and SIGTERM
	// to allow graceful shutdowns.
	go func() {
		// Create a quit channel which carries os.Signal values.
		quit := make(chan os.Signal, 1)
		// Listen for incoming SIGINT and SIGTERM signals and
		// relay them to the quit channel.
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		// Read the signal from the quit channel. It blocks until a signal is received.
		s := <-quit
		// Log the signal that has been caught.
		app.logger.PrintInfo("caught signal", map[string]string{
			"signal": s.String(),
		})
		// Create a context with a 5-second timeout.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		// Call Shutdown() on the server, pass in the context. Shutdown() will return
		// nil if the graceful shutdown was successful, or an error (which may happen
		// because of a problem closing the listeners, or because the shutdown didn't
		// complete before the 5-second context deadline is hit). Relay this return
		// value to the shutdownError channel.
		shutdownError <- srv.Shutdown(ctx)
	}()

	app.logger.PrintInfo("starting server", map[string]string{
		"addr": srv.Addr,
		"env":  app.config.env,
	})
	// Calling Shutdown() on our server will cause ListenAndServe() to immediately
	// return a http.ErrServerClosed error. So if we see this error, it is actually a
	// good thing and an indication that the graceful shutdown has started. So we check
	// specifically for this, only returning the error if it is NOT http.ErrServerClosed.
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	// Otherwise, we wait to receive the return value from Shutdown() on the
	// shutdownError channel. If return value is an error, we know that there was a
	// problem with the graceful shutdown and we return the error.
	err = <-shutdownError
	if err != nil {
		return err
	}
	// At this point we know that the graceful shutdown completed successfully and we
	// log a proper message.
	app.logger.PrintInfo("stopped server", map[string]string{
		"addr": srv.Addr,
		"env":  app.config.env,
	})
	return nil
}

func (app *application) exposeMetrics() {
	// Expose additional variables/stats that will be displayed via metrics endpoint.
	// Publish current version of the API.
	expvar.NewString("version").Set(version)
	// Publish the number of active goroutines.
	expvar.Publish("goroutines", expvar.Func(func() interface{} {
		return runtime.NumGoroutine()
	}))
	// Publish the database connection pool statistics.
	expvar.Publish("database", expvar.Func(func() interface{} {
		return app.db.Stats()
	}))
	// Publish the current Unix timestamp.
	expvar.Publish("timestamp", expvar.Func(func() interface{} {
		return time.Now().Unix()
	}))
}
