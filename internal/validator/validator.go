package validator

import (
	"regexp"
	"slices"
)

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) AddError(key, msg string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = msg
	}
}

func (v *Validator) Check(ok bool, key, msg string) {
	if !ok {
		v.AddError(key, msg)
	}
}

func PermittedValue[T comparable](v T, permitted ...T) bool {
	return slices.Contains(permitted, v)
}

func Matches(s string, rx *regexp.Regexp) bool {
	return rx.MatchString(s)
}

func Unique[T comparable](values []T) bool {
	unique := make(map[T]bool)
	for _, v := range values {
		unique[v] = true
	}

	return len(unique) == len(values)
}
