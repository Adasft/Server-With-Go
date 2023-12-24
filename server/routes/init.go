package routes

import (
	"log"
	"net/http"

	"server/routermanager"
)

const (
	homePath    = "/"
	loginPath   = "/login"
	signupPath  = "/signup"
	recoverPath = "/login/recover"
)

func InitRouter() *routermanager.Router {
	router := routermanager.New()

	initHomeRouter(router)
	initLoginRouter(router)
	initSignupRouter(router)

	return router
}

func SetHandlerFunc(router *routermanager.Router) {
	ids := router.GetIds()

	for _, id := range ids {
		route, err := router.Get(id)

		if err != nil {
			log.Println(err)
		}

		http.HandleFunc(route.GetPath(), route.GetHandler())
	}
}
