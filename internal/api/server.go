package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/evangodon/jrnl/internal/db"
)

// ~~~~~ Config ~~~~~ //
type Config struct {
	Port int
	Env  string
}

// ~~~~~ App ~~~~~ //
type Application struct {
	Cfg      Config
	DBClient db.DB
}

func (app Application) Serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.Cfg.Port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  8 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		fmt.Printf(
			"Shutting down due to signal: %s\n",
			map[string]string{"signal": s.String()},
		)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		shutdownError <- srv.Shutdown(ctx)
	}()

	fmt.Println("Starting server", map[string]string{
		"addr": srv.Addr,
	})

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	// err = <-shutdownError
	if err != nil {
		return err
	}

	return nil
}
