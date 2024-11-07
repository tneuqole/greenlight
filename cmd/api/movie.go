package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/tneuqole/greenlight/internal/model"
)

func (app *application) postMovie(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "post movie")
}

func (app *application) getMovie(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	movie := model.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Casablanca",
		Runtime:   102,
		Genres:    []string{"drama", "romance", "war"},
		Version:   1,
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
