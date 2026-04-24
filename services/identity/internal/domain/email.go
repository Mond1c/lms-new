package domain

import (
	"errors"
	"regexp"
	"strings"
)

type Email struct {
	value string
}

var (
	ErrInvalidEmail = errors.New("invalid email")
	emailRegex      = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
)

func NewEmail(raw string) (Email, error) {
	v := strings.ToLower(strings.TrimSpace(raw))
	if !emailRegex.MatchString(v) {
		return Email{}, ErrInvalidEmail
	}
	return Email{value: v}, nil
}

func EmailFromTrusted(raw string) Email {
	return Email{value: raw}
}

func (e Email) String() string {
	return e.value
}
