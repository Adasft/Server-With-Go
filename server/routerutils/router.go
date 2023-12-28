package routerutils

import (
	"net/http"
)

type RouteHandlerFunc http.HandlerFunc
type Middleware func(http.Handler) http.Handler
type Routes []*Route

type HTTPMethod string

func (m *HTTPMethod) ToString() string {
	return string(*m)
}

type Route struct {
	method      HTTPMethod
	handlerFunc RouteHandlerFunc
	middleware  Middleware
}

func (r *Route) GetMethod() HTTPMethod {
	return r.method
}

func (r *Route) GetHandlerFunc() RouteHandlerFunc {
	return r.handlerFunc
}

func (r *Route) GetMiddleware() Middleware {
	return r.middleware
}

func (r *Route) ShouldApplyMiddleware() bool {
	return r.middleware != nil
}

type Router struct {
	pathRoutes map[string]Routes
}

func (r *Router) defineNewRoute(path string, handlerFunc RouteHandlerFunc, middleware Middleware, method HTTPMethod) {
	if r.pathRoutes == nil {
		return
	}

	r.pathRoutes[path] = append(r.pathRoutes[path], &Route{
		method:      method,
		handlerFunc: handlerFunc,
		middleware:  middleware,
	})
}

func (r *Router) GetRouteByMethod(path string, method HTTPMethod) *Route {
	routes := r.pathRoutes[path]

	for _, route := range routes {
		if route.method == method {
			return route
		}
	}
	return nil
}

func (r *Router) GetPathRoutes() *map[string]Routes {
	return &r.pathRoutes
}

func (r *Router) Get(path string, handlerFunc RouteHandlerFunc, middleware Middleware) {
	r.defineNewRoute(path, handlerFunc, middleware, http.MethodGet)
}

func (r *Router) Post(path string, handlerFunc RouteHandlerFunc, middleware Middleware) {
	r.defineNewRoute(path, handlerFunc, middleware, http.MethodPost)
}

func New() *Router {
	return &Router{
		pathRoutes: make(map[string]Routes),
	}
}
