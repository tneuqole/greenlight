package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	r := httprouter.New()

	r.HandlerFunc(http.MethodGet, "/v1/health", app.healthHandler)
	r.HandlerFunc(http.MethodPost, "/v1/movies", app.postMovie)
	r.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.getMovie)

	return r
}
