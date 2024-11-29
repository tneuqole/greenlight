package main

import (
	"expvar"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	r := httprouter.New()

	r.NotFound = http.HandlerFunc(app.notFoundResponse)
	r.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	r.HandlerFunc(http.MethodGet, "/v1/health", app.healthHandler)
	r.Handler(http.MethodGet, "/debug/vars", expvar.Handler())

	r.HandlerFunc(http.MethodPost, "/v1/users", app.postUser)
	r.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUser)
	r.HandlerFunc(http.MethodPut, "/v1/users/password", app.putPassword)

	r.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.postAuthenticationToken)
	r.HandlerFunc(http.MethodPost, "/v1/tokens/password-reset", app.postPasswordResetToken)

	r.HandlerFunc(http.MethodPost, "/v1/movies", app.requirePermission("movies:write", app.postMovie))
	r.HandlerFunc(http.MethodGet, "/v1/movies", app.requirePermission("movies:read", app.getMovies))
	r.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.requirePermission("movies:read", app.getMovie))
	r.HandlerFunc(http.MethodPatch, "/v1/movies/:id", app.requirePermission("movies:write", app.putMovie))
	r.HandlerFunc(http.MethodDelete, "/v1/movies/:id", app.requirePermission("movies:write", app.deleteMovie))

	return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(r)))))
}
