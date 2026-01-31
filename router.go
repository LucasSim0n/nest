package cafe

import (
	"net/http"
)

/*** Definitions ***/

type route struct {
	path    string
	method  string
	handler http.HandlerFunc
}

type Router struct {
	routes      []route
	routers     []mountedRouter
	middlewares []middleware
}

type mountedRouter struct {
	path   string
	router *Router
}

/*** Init ***/

func NewRouter() *Router {
	return &Router{
		routes:  []route{},
		routers: []mountedRouter{},
	}
}

/*** Aggregation ***/

func (r *Router) UseRouter(path string, ro *Router) {
	for _, mr := range r.routers {
		if mr.path == path {
			return
		}
	}
	r.routers = append(r.routers, mountedRouter{path: path, router: ro})
}

func (r *Router) Use(mw middleware) {
	r.middlewares = append(r.middlewares, mw)
}

/*** Assembly ***/

func (r *Router) getRoutes() []route {
	mountedRoutes := []route{}
	for _, rt := range r.routes {
		rt.handler = setUpMiddlewares(rt.handler, r.middlewares)
		mountedRoutes = append(mountedRoutes, rt)
	}
	for _, mr := range r.routers {
		rtrRoutes := mr.router.getRoutes()
		for _, rt := range rtrRoutes {
			rt.path = mr.path + rt.path
			rt.handler = setUpMiddlewares(rt.handler, r.middlewares)
			mountedRoutes = append(mountedRoutes, rt)
		}
	}
	return mountedRoutes
}

/*** Basic HTTP Methods ***/

func (r *Router) Get(path string, handler http.HandlerFunc) {
	r.routes = addRoute(r.routes, path, "GET", handler)
}

func (r *Router) Post(path string, handler http.HandlerFunc) {
	r.routes = addRoute(r.routes, path, "POST", handler)
}

func (r *Router) Put(path string, handler http.HandlerFunc) {
	r.routes = addRoute(r.routes, path, "PUT", handler)
}

func (r *Router) Delete(path string, handler http.HandlerFunc) {
	r.routes = addRoute(r.routes, path, "DELETE", handler)
}
