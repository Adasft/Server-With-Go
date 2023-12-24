package routermanager

import (
	"errors"
	"net/http"

	"server/serrors"
)

type RouteHandler func(http.ResponseWriter, *http.Request)

type Route struct {
	path    string
	handler RouteHandler
}

func (r *Route) GetPath() string {
	return r.path
}

func (r *Route) GetHandler() RouteHandler {
	return r.handler
}

type Router struct {
	routes map[string]*Route
	ids    []string
}

func (r *Router) GetIds() []string {
	return r.ids
}

func (r *Router) Set(routeId string, path string, handler RouteHandler) {
	r.ids = append(r.ids, routeId)
	r.routes[routeId] = &Route{
		path:    path,
		handler: handler,
	}
}

func (r *Router) Get(routeId string) (*Route, error) {
	route, exists := r.routes[routeId]

	if !exists {
		return nil, errors.New(serrors.RouterIDNotExistError)
	}

	return route, nil
}

func New() *Router {
	return &Router{
		routes: make(map[string]*Route),
	}
}
