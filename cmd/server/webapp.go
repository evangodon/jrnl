package main

import (
	"net/http"
	"path"

	"github.com/uptrace/bunrouter"

	"github.com/evangodon/jrnl/internal/cfg"
)

func (srv Server) sendWebApp() bunrouter.HandlerFunc {
	root := cfg.GetProjectRoot()
	buildpath := path.Join(root, "/app/build/")

	webappDir := http.Dir(buildpath)
	fileServer := http.FileServer(webappDir)

	return bunrouter.HTTPHandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/" {
			http.ServeFile(w, req, path.Join(buildpath, "200.html"))
			return
		}

		fileServer.ServeHTTP(w, req)
	})
}
