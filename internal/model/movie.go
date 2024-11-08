package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/lib/pq"
	"github.com/tneuqole/greenlight/internal/validator"
)

const Q_TIMEOUT = 3 * time.Second

var ErrInvalidRuntimeFormat = errors.New("invalid runtime format")

type Runtime int32

func (r Runtime) MarshalJSON() ([]byte, error) {
	s := fmt.Sprintf("%d mins", r)
	s = strconv.Quote(s)

	return []byte(s), nil
}

func (r *Runtime) UnmarshalJSON(data []byte) error {
	unquoted, err := strconv.Unquote(string(data))
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	parts := strings.Split(unquoted, " ")
	if len(parts) != 2 || parts[1] != "mins" {
		return ErrInvalidRuntimeFormat
	}

	i, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	*r = Runtime(i)

	return nil
}

type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`
	Runtime   Runtime   `json:"runtime,omitempty"`
	Genres    []string  `json:"genres,omitempty"`
	Version   int32     `json:"version"`
}

func ValidateMovie(v *validator.Validator, m *Movie) {
	v.Check(m.Title != "", "title", "must be provided")
	v.Check(len(m.Title) <= 500, "title", "must not be more than 500 bytes long")

	v.Check(m.Year >= 1888, "year", "must be greater than 1888")
	v.Check(m.Year <= int32(time.Now().Year()), "year", "must not be in the future")

	v.Check(m.Runtime > 0, "runtime", "must be a positive integer")

	v.Check(len(m.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(m.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(m.Genres), "genres", "values must be unique")
}

type MovieModel struct {
	DB *sql.DB
}

func (m MovieModel) Insert(movie *Movie) error {
	q := `INSERT INTO movie (title, year, runtime, genres)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at, version`

	ctx, cancel := context.WithTimeout(context.Background(), Q_TIMEOUT)
	defer cancel()

	args := []any{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)}
	return m.DB.QueryRowContext(ctx, q, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
}

func (m MovieModel) Get(id int64) (*Movie, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	q := `SELECT id, created_at, title, year, runtime, genres, version
	FROM movie WHERE id=$1`

	var movie Movie
	ctx, cancel := context.WithTimeout(context.Background(), Q_TIMEOUT)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, q, id).Scan(
		&movie.ID,
		&movie.CreatedAt,
		&movie.Title,
		&movie.Year,
		&movie.Runtime,
		pq.Array(&movie.Genres),
		&movie.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &movie, nil
}

func (m MovieModel) Update(movie *Movie) error {
	q := `UPDATE movie SET title=$1, year=$2, runtime=$3, genres=$4, version=version+1
	WHERE id=$5 AND version=$6
	RETURNING version`

	args := []any{
		movie.Title,
		movie.Year,
		movie.Runtime,
		pq.Array(&movie.Genres),
		movie.ID,
		movie.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), Q_TIMEOUT)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, q, args...).Scan(&movie.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (m MovieModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	q := `DELETE from movie where id=$1`

	ctx, cancel := context.WithTimeout(context.Background(), Q_TIMEOUT)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, q, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
