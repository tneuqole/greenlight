package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/tneuqole/greenlight/internal/validator"
)

type envelope map[string]any

func (app *application) readIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id param")
	}

	return id, nil
}

func (app *application) writeJSON(w http.ResponseWriter, status int, data any, headers http.Header) error {
	// note: MarhsalIndent() is slower than Marhsal()
	j, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	for k, v := range headers {
		w.Header()[k] = v
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(j)

	return nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dest any) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	err := d.Decode(dest)
	if err != nil {
		var (
			syntaxError           *json.SyntaxError
			unmarshalTypeError    *json.UnmarshalTypeError
			invalidUnmarshalError *json.InvalidUnmarshalError
			maxBytesError         *http.MaxBytesError
		)

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains malformed JSON at character %d", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains malformed JSON")
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("wrong type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("wrong type at character %d", unmarshalTypeError.Offset)
		case errors.Is(err, io.EOF):
			return errors.New("body cannot be empty")
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			field := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("unknown field %s", field)
		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body cannot be larger than %d bytes", maxBytesError.Limit)
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			return err
		}
	}

	err = d.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return errors.New("body must be a single json object")
	}

	return nil
}

func (app *application) readString(qp url.Values, key string, defaultVal string) string {
	s := qp.Get(key)
	if s == "" {
		return defaultVal
	}

	return s
}

func (app *application) readList(qp url.Values, key string, defaultVal []string) []string {
	vals := qp.Get(key)
	if vals == "" {
		return defaultVal
	}

	return strings.Split(vals, ",")
}

func (app *application) readInt(qp url.Values, key string, defaultVal int, v *validator.Validator) int {
	s := qp.Get(key)
	if s == "" {
		return defaultVal
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, "must be an integer value")
		return defaultVal
	}

	return i
}

func (app *application) background(fn func()) {
	app.wg.Add(1)
	go func() {
		defer app.wg.Done()
		defer func() {
			if err := recover(); err != nil {
				app.logger.Error(fmt.Sprintf("%v", err))
			}
		}()

		fn()
	}()
}
