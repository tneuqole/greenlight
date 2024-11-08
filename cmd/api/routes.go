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
	r.HandlerFunc(http.MethodPost, "/v1/movies", app.postMovie)
	r.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.getMovie)
	r.HandlerFunc(http.MethodPut, "/v1/movies/:id", app.putMovie)
	r.HandlerFunc(http.MethodPatch, "/v1/movies/:id", app.putMovie)
	r.HandlerFunc(http.MethodDelete, "/v1/movies/:id", app.deleteMovie)

	return app.recoverPanic(r)
}
