package main

import (
	"fmt"
	"net/http"
)

func (app *application) postMovie(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "post movie")
}

func (app *application) getMovie(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
	}

	fmt.Fprintf(w, "get movie %d\n", id)
}
