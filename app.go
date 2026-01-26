package nest

import (
	"fmt"
	"net/http"
	"strings"
)

type App struct {
	server  http.Server
	handler *http.ServeMux
	routers []mountedRouter
	routes  []route
}

func NewServer() App {
	return App{
		handler: http.NewServeMux(),
		routers: []mountedRouter{},
		routes:  []route{},
	}
}

func (a *App) UseRouter(path string, ro *router) {
	for _, mr := range a.routers {
		if mr.path == path {
			return
		}
	}
	a.routers = append(a.routers, mountedRouter{path: path, router: ro})
}

func (a *App) Listen(addr string) error {
	a.setUpRouters()
	a.server = http.Server{
		Addr:    addr,
		Handler: a.handler,
	}
	return a.server.ListenAndServe()
}

func (a *App) setUpRouters() {
	for _, r := range a.routes {
		a.handle(r.path, r.method, r.handler)
	}

	for _, mr := range a.routers {
		routes := mr.router.getRoutes()
		for _, r := range routes {
			path := mr.path + r.path
			a.handle(path, r.method, r.handler)
		}
	}
}

func (a *App) handle(path, method string, handler http.HandlerFunc) {
	patt := fmt.Sprintf("%s %s", method, path)
	if !strings.HasSuffix(patt, "/") {
		patt += "/"
	}
	patt += "{$}"
	a.handler.HandleFunc(patt, handler)
}

func addRoute(routes []route, path, method string, handler http.HandlerFunc) []route {
	for _, r := range routes {
		if r.path == path && r.method == method {
			return routes
		}
	}
	return append(routes, route{
		path:    path,
		method:  method,
		handler: handler,
	})
}

/*** Basic HTTP Methods ***/

func (a *App) Get(path string, handler http.HandlerFunc) {
	a.routes = addRoute(a.routes, path, "GET", handler)
}

func (a *App) Post(path string, handler http.HandlerFunc) {
	a.routes = addRoute(a.routes, path, "POST", handler)
}

func (a *App) Put(path string, handler http.HandlerFunc) {
	a.routes = addRoute(a.routes, path, "PUT", handler)
}

func (a *App) Delete(path string, handler http.HandlerFunc) {
	a.routes = addRoute(a.routes, path, "DELETE", handler)
}
