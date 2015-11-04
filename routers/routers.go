package routers

import (
	"github.com/aosen/kernel"
	"github.com/gorilla/mux"
	"net/http"
	"searchserver/controllers"
)

func Register(g *kernel.G) *mux.Router {
	// Routes consist of a path and a handler function.
	r := mux.NewRouter()
	// Bind to a port and pass our router in
	r.HandleFunc("/search/", g.Go(&controllers.SearchHandler{}))
	r.HandleFunc("/index/", g.Go(&controllers.IndexHandler{}))
	r.HandleFunc("/cut/", g.Go(&controllers.CutHandler{}))
	r.NotFoundHandler = http.HandlerFunc(g.Go(&controllers.NotFoundHandler{}))
	return r
}
