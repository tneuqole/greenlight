package model

import (
	"fmt"
	"strconv"
	"time"
)

type Runtime int32

func (r Runtime) MarshalJSON() ([]byte, error) {
	s := fmt.Sprintf("%d mins", r)
	s = strconv.Quote(s)

	return []byte(s), nil
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
