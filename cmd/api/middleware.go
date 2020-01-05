package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")

		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthentication(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		token := r.Header.Get("Authorization")
		app.infoLog.Printf("%s", token)
		id, err := app.verifyToken(token)
		if err != nil {
			app.generateErrorResponse(w, err, http.StatusUnauthorized)
			return
		}

		if id == 0 || id < 0 {
			app.generateErrorResponse(w, err, http.StatusUnauthorized)
			return
		}

		user, err := app.user.UserGetByID(id)
		if err != nil {
			app.generateErrorResponse(w, err, http.StatusUnauthorized)
		}
		ctx := context.WithValue(r.Context(), contextRequestUser, user)
		r = r.WithContext(ctx)

		h(w, r, ps)
	}
}
