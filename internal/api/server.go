package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
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

func NewServer(srvConfig ServerConfig, dbClient db.DB) *Server {
	return &Server{
		cfg:      srvConfig,
		dbClient: dbClient,
		appCfg:   cfg.GetConfig(),
		logger:   logger.NewLogger(os.Stdout),
	}
}

// Get preferred outbound ip of this machine
func getOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func (server Server) Serve() error {
	addr := fmt.Sprintf("%s:%d", getOutboundIP().String(), server.cfg.Port)
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
			map[string]string{"signal": s.String()},
		)
		server.logger.PrintInfo(log)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		shutdownError <- srv.Shutdown(ctx)
	}()

	log := fmt.Sprintf("🚀 Starting server at %s", srv.Addr)
	server.logger.PrintInfo(log)

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
