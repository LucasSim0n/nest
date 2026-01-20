package quick

import "net/http"

type route struct {
	url     string
	method  string
	handler http.HandlerFunc
}

type Router struct {
	url    string
	routes []route
}

type App struct {
	server  *http.ServeMux
	routers map[string]*Router
}

func NewServer() App {
	return App{
		server:  http.NewServeMux(),
		routers: make(map[string]*Router, 0),
	}
}

func (r *Router) Route(url string)

func (r *Router) Get(url string, handler http.HandlerFunc) {
	r.routes = append(r.routes, route{
		method:  "GET",
		url:     url,
		handler: handler,
	})
}

func (r *Router) Post(url string, handler http.HandlerFunc) {
	r.routes = append(r.routes, route{
		method:  "POST",
		url:     url,
		handler: handler,
	})
}

func (r *Router) Put(url string, handler http.HandlerFunc) {
	r.routes = append(r.routes, route{
		method:  "PUT",
		url:     url,
		handler: handler,
	})
}

func (r *Router) Delete(url string, handler http.HandlerFunc) {
	r.routes = append(r.routes, route{
		method:  "DELETE",
		url:     url,
		handler: handler,
	})
}

func (a *App) UseRouter(url string, r *Router) {
	a.routers[url] = r
}

func (a *App) Listen()

func (a *App) setUpRouters()
