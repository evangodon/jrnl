package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/evangodon/jrnl/internal/logger"

	"github.com/evangodon/jrnl/internal/cfg"
	"github.com/evangodon/jrnl/internal/db"
)

type ServerConfig struct {
	Port int
	Env  string
}

type Server struct {
	cfg      ServerConfig
	dbClient db.DB
	appCfg   cfg.Config
	logger   *logger.Logger
}

func NewServer(srvConfig ServerConfig) *Server {
	return &Server{
		cfg:      srvConfig,
		dbClient: db.Connect(),
		appCfg:   cfg.GetConfig(),
		logger:   logger.NewLogger(os.Stdout),
	}
}

func main() {
	serverCfg := ServerConfig{
		Env: cfg.GetEnv(),
	}
	flag.IntVar(&serverCfg.Port, "port", 8090, "API server port")
	flag.Parse()

	server := NewServer(serverCfg)
	addr := fmt.Sprintf(":%d", serverCfg.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      server.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  8 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		log := fmt.Sprintf(
			"Shutting down due to signal: %s\n",
			strings.ToUpper(s.String()),
		)
		server.logger.PrintInfo(log)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		shutdownError <- srv.Shutdown(ctx)
	}()

	msg := fmt.Sprintf("🚀 Starting server at %s\n", srv.Addr)
	server.logger.Print(msg)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}

	err = <-shutdownError
	if err != nil {
		log.Fatal(err)
	}
}
