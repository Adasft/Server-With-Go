package routes

import (
	"context"
	"net/http"

	"github.com/go-session/session"
	"server/routermanager"
	"server/serrors"
	"server/template"
)

func getHomeHandler(w http.ResponseWriter, r *http.Request) {
	store, err := session.Start(context.Background(), w, r)
	if err != nil {
		serrors.InternalServerErrorHandler(w, err, loginPath)
		return
	}

	_, ok := store.Get("user_id")

	if !ok {
		http.Redirect(w, r, loginPath, http.StatusSeeOther)
		return
	}

	_, err = template.Render(w, nil, template.GetView("index"), template.GetLayout("home"))

	if err != nil {
		serrors.InternalServerErrorHandler(w, err, homePath)
	}
}

func initHomeRouter(router *routermanager.Router) {
	router.Set("home", homePath, getHomeHandler)
}
