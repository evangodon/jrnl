package api

import (
	"net/http"

	"github.com/uptrace/bunrouter"
)

func (app Application) routes() http.Handler {
	router := bunrouter.New()

	router.WithGroup("/daily", func(group *bunrouter.Group) {
		group.GET("/", app.getDailyHandler())
		group.GET("/list", app.listDailyHandler())
		group.POST("/new", app.newDailyHandler())
	})

	return router
}
