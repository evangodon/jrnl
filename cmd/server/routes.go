package main

import (
	"net/http"

	"github.com/uptrace/bunrouter"
	"github.com/uptrace/bunrouter/extra/reqlog"
)

func (server Server) routes() http.Handler {
	router := bunrouter.New(
		bunrouter.Use(reqlog.NewMiddleware()),
		bunrouter.Use(corsMiddleware),
		bunrouter.Use(errorHandler),
	)

	router.GET("/dbpath", server.dbPath)
	router.GET("/healthcheck", server.healthcheck)
	router.GET("/*path", server.sendWebApp())

	router.
		Use(server.validateAPIKeyMiddleware).
		WithGroup("/daily", func(group *bunrouter.Group) {
			group.GET("/", server.handleGetDaily())
			group.GET("/:date", server.handleGetDaily())
			group.GET("/template", server.handleGetDailyTemplate())
			group.PATCH("/:id", server.handleUpdateDaily())
			group.GET("/list", server.handleListDaily())
			group.POST("/new", server.handleNewDaily())
		})

	return router
}
