package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/tneuqole/greenlight/internal/model"
	"github.com/tneuqole/greenlight/internal/validator"
)

func (app *application) postMovie(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string        `json:"title"`
		Year    int32         `json:"year"`
		Runtime model.Runtime `json:"runtime"`
		Genres  []string      `json:"genres"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	m := &model.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}

	if model.ValidateMovie(v, m); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Movies.Insert(m)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", m.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"movie": m}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getMovie(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	movie, err := app.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) putMovie(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	movie, err := app.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Title   *string        `json:"title"`
		Year    *int32         `json:"year"`
		Runtime *model.Runtime `json:"runtime"`
		Genres  []string       `json:"genres"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Title != nil {
		movie.Title = *input.Title
	}

	if input.Year != nil {
		movie.Year = *input.Year
	}

	if input.Runtime != nil {
		movie.Runtime = *input.Runtime
	}

	if input.Genres != nil {
		movie.Genres = input.Genres
	}

	v := validator.New()
	if model.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Movies.Update(movie)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteMovie(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Movies.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusNoContent, nil, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getMovies(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title  string
		Genres []string
		model.Filters
	}

	v := validator.New()
	qp := r.URL.Query()

	input.Title = app.readString(qp, "title", "")
	input.Genres = app.readList(qp, "genres", []string{})
	input.Page = app.readInt(qp, "page", 1, v)
	input.PageSize = app.readInt(qp, "page_size", 20, v)
	input.Sort = app.readString(qp, "sort", "id")
	input.SortAllowlist = []string{"id", "-id", "title", "-title", "year", "-year", "runtime", "-runtime"}

	if model.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	movies, metadata, err := app.models.Movies.GetAll(input.Title, input.Genres, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"metadata": metadata, "movies": movies}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
