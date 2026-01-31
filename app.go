package cafe

import (
	"fmt"
	"net/http"
	"strings"
)

/*** Definitions ***/

type App struct {
	server      http.Server
	handler     *http.ServeMux
	routers     []mountedRouter
	routes      []route
	middlewares []middleware
}

type middleware func(next http.HandlerFunc) http.HandlerFunc

/*** Init ***/

func NewServer() App {
	return App{
		handler:     http.NewServeMux(),
		routers:     []mountedRouter{},
		routes:      []route{},
		middlewares: []middleware{},
	}
}

/*** Aggregation ***/

func (a *App) UseRouter(path string, ro *Router) {
	for _, mr := range a.routers {
		if mr.path == path {
			return
		}
	}
	a.routers = append(a.routers, mountedRouter{path: path, router: ro})
}

func (a *App) Use(mw middleware) {
	a.middlewares = append(a.middlewares, mw)
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

/*** Setup ***/

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
		h := setUpMiddlewares(r.handler, a.middlewares)
		a.handle(r.path, r.method, h)
	}

	for _, mr := range a.routers {
		routes := mr.router.getRoutes()
		for _, r := range routes {
			path := mr.path + r.path
			h := setUpMiddlewares(r.handler, a.middlewares)
			a.handle(path, r.method, h)
		}
	}
}

func setUpMiddlewares(f http.HandlerFunc, mws []middleware) http.HandlerFunc {
	if len(mws) == 0 {
		return f
	}

	for i := len(mws) - 1; i >= 0; i-- {
		f = mws[i](f)
	}
	return f
}

func (a *App) handle(path, method string, handler http.HandlerFunc) {
	patt := fmt.Sprintf("%s %s", method, path)
	if !strings.HasSuffix(patt, "/") {
		patt += "/"
	}
	patt += "{$}"
	a.handler.Handle(patt, handler)
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
