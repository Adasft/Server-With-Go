package routes

import (
	"net/http"

	"server/routerutils"
)

const (
	HomePath    = "/"
	LoginPath   = "/login"
	SignupPath  = "/signup"
	RecoverPath = "/login/recover"
)

func InitRouter() *routerutils.Router {
	router := routerutils.New()

	initHomeRouter(router)
	initLoginRouter(router)
	initSignupRouter(router)
	initRecoverRouter(router)

	return router
}

func processRoute(w http.ResponseWriter, r *http.Request, route *routerutils.Route) {
	if route.ShouldApplyMiddleware() {
		handler := http.HandlerFunc(route.GetHandlerFunc())
		middleware := route.GetMiddleware()
		middleware(handler).ServeHTTP(w, r)
	} else {
		route.GetHandlerFunc()(w, r)
	}
}

func configureRouteHandler(path string, router *routerutils.Router) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := router.GetRouteByMethod(path, routerutils.HTTPMethod(r.Method))
		processRoute(w, r, route)
	})
}

func setupRoutes(path string, router *routerutils.Router) {
	http.Handle(path, configureRouteHandler(path, router))
}

func SetHandlerFunc(router *routerutils.Router) {
	for path, _ := range *router.GetPathRoutes() {
		setupRoutes(path, router)
	}
}
