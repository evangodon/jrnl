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
			group.GET("/", server.getDailyHandler())
			group.GET("/:date", server.getDailyHandler())
			group.GET("/template", server.getDailyTemplate())
			group.PATCH("/:id", server.updateDailyHandler())
			group.GET("/list", server.listDailyHandler())
			group.POST("/new", server.newDailyHandler())
		})

	return router
}
