package main

import (
	"errors"
	"strings"
)

type StringService interface {
	Uppercase(str string) (string, error)
	Count(str string) int
}

// Implement interface
type stringService struct{}

func (stringService) Uppercase (str string) (string, error) {
	if str == "" {
		return "", ErrEmpty
	}
	return strings.ToUpper(str), nil
}

func (stringService) Count(str string) int {
	return len(str)
}

// ErrEmpty is returned when an input string is empty.
var ErrEmpty = errors.New("empty string")