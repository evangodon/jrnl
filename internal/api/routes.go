package api

import (
	"net/http"

	"github.com/uptrace/bunrouter"
	"github.com/uptrace/bunrouter/extra/reqlog"
)

func (server Server) routes() http.Handler {
	router := bunrouter.New(
		bunrouter.Use(reqlog.NewMiddleware()),
	)

	router.GET("/dbpath", server.dbPath)
	router.GET("/healthcheck", server.healthcheck)

	router.Use(server.checkAPIKey).WithGroup("/daily", func(group *bunrouter.Group) {
		group.GET("/", server.getDailyHandler())
		group.GET("/:date", server.getDailyHandler())
		group.PATCH("/:id", server.updateDailyHandler())
		group.GET("/list", server.listDailyHandler())
		group.POST("/new", server.newDailyHandler())
	})

	return router
}
