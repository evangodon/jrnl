package api

import (
	"net/http"

	"github.com/uptrace/bunrouter"
	"github.com/uptrace/bunrouter/extra/reqlog"
)

func (app Application) routes() http.Handler {
	router := bunrouter.New(
		bunrouter.Use(reqlog.NewMiddleware()),
	)

	router.GET("/dbpath", app.dbPath)
	router.GET("/healthcheck", app.healthcheck)

	router.Use(app.checkAPIKey).WithGroup("/daily", func(group *bunrouter.Group) {
		group.GET("/", app.getDailyHandler())
		group.GET("/:date", app.getDailyHandler())
		group.PATCH("/:id", app.updateDailyHandler())
		group.GET("/list", app.listDailyHandler())
		group.POST("/new", app.newDailyHandler())
	})

	return router
}
