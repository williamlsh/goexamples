package main

import (
	"errors"
	"strings"
)

// StringService provides operations on strings.
type StringService interface {
	Uppercase(string) (string, error)
	Count(string) int
}

// stringService is a concrete implementation of StringService.
type stringService struct{}

func (stringService) Uppercase(s string) (string, error) {
	if s == "" {
		return "", ErrEmpty
	}
	return strings.ToUpper(s), nil
}

func (stringService) Count(s string) int {
	return len(s)
}

// ErrEmpty is returned when input string is empty.
var ErrEmpty = errors.New("Empty string")
