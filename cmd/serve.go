package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/uptrace/bunrouter"
	"github.com/urfave/cli/v2"
)

// ~~~~~ Config ~~~~~ //
type config struct {
	port int
	env  string
}

// ~~~~~ Router ~~~~~ //
func (app *application) routes() http.Handler {
	router := bunrouter.New()

	router.GET("/", func(_ http.ResponseWriter, req bunrouter.Request) error {
		fmt.Println(req.Method, req.Route(), req.Params().Map())
		return nil
	})

	router.GET("/list", func(_ http.ResponseWriter, req bunrouter.Request) error {
		fmt.Println(req.Method, req.Route(), req.Params().Map())
		return nil
	})

	router.POST("/new", func(_ http.ResponseWriter, req bunrouter.Request) error {
		return nil
	})

	return router
}

// ~~~~~ App ~~~~~ //
type application struct {
	cfg config
}

func (app application) serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.cfg.port),
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
			"Shutting down due to signal: %s",
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

	err = <-shutdownError
	if err != nil {
		return err
	}

	return nil
}

var ServeCmd = &cli.Command{
	Name:    "serve",
	Aliases: []string{"s"},
	Usage:   "Start the server",
	Action: func(_ *cli.Context) error {
		cfg := config{
			port: 8080,
		}

		app := &application{
			cfg: cfg,
		}

		app.serve()
		return nil
	},
}
