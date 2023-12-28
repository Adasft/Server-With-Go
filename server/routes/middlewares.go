package routes

import (
	"context"
	"net/http"

	"github.com/go-session/session"
	"server/errs"
)

func denyAccessToHomeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		store, err := session.Start(context.Background(), w, r)
		if err != nil {
			errs.InternalServerErrorHandler(w, err, LoginPath)
			return
		}

		_, ok := store.Get("user_id")

		if !ok {
			http.Redirect(w, r, LoginPath, http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func denyAccessIfAlreadyLoggedInMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		store, err := session.Start(context.Background(), w, r)
		if err != nil {
			errs.InternalServerErrorHandler(w, err, LoginPath)
			return
		}

		_, ok := store.Get("user_id")

		if ok {
			http.Redirect(w, r, HomePath, http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}
