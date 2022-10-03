package api

import (
	"net/http"

	"github.com/uptrace/bunrouter"
)

func (app Application) routes() http.Handler {
	router := bunrouter.New()

	router.Use(app.checkAPIKey).WithGroup("/daily", func(group *bunrouter.Group) {

		group.GET("/", app.getDailyHandler())
		group.GET("/:date", app.getDailyHandler())
		group.PATCH("/:id", app.updateDailyHandler())
		group.GET("/list", app.listDailyHandler())
		group.POST("/new", app.newDailyHandler())

	})

	return router
}
