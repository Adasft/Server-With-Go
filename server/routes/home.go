package routes

import (
	"net/http"

	"server/errs"
	"server/routerutils"
	"server/template"
)

func homeHandlerGet(w http.ResponseWriter, r *http.Request) {
	_, err := template.Render(w, nil, template.GetView("index"), template.GetLayout("home"))

	if err != nil {
		errs.InternalServerErrorHandler(w, err, HomePath)
	}
}

func initHomeRouter(router *routerutils.Router) {
	router.Get(HomePath, homeHandlerGet, denyAccessToHomeMiddleware)
}
