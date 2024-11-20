package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	r := httprouter.New()

	r.NotFound = http.HandlerFunc(app.notFoundResponse)
	r.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	r.HandlerFunc(http.MethodGet, "/v1/health", app.healthHandler)

	r.HandlerFunc(http.MethodPost, "/v1/users", app.postUser)
	r.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUser)

	r.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.getAuthenticationToken)

	r.HandlerFunc(http.MethodPost, "/v1/movies", app.requirePermission("movies:write", app.postMovie))
	r.HandlerFunc(http.MethodGet, "/v1/movies", app.requirePermission("movies:read", app.getMovies))
	r.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.requirePermission("movies:read", app.getMovie))
	r.HandlerFunc(http.MethodPatch, "/v1/movies/:id", app.requirePermission("movies:write", app.putMovie))
	r.HandlerFunc(http.MethodDelete, "/v1/movies/:id", app.requirePermission("movies:write", app.deleteMovie))

	return app.recoverPanic(app.rateLimit(app.authenticate(r)))
}
