package domain

import (
	"errors"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
)

type PasswordHash string

var (
	ErrPasswordTooShort = errors.New("password too short")
	ErrPasswordTooLong  = errors.New("password too long")
)

const (
	MinPasswordLength = 8
	MaxPasswordLength = 72
)

func HashPassword(raw string) (PasswordHash, error) {
	n := utf8.RuneCountInString(raw)
	if n < MinPasswordLength {
		return "", ErrPasswordTooShort
	}
	if n > MaxPasswordLength {
		return "", ErrPasswordTooLong
	}
	h, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost) // TODO: increase cost factor
	if err != nil {
		return "", err
	}
	return PasswordHash(h), nil
}

func (h PasswordHash) Verify(raw string) error {
	return bcrypt.CompareHashAndPassword([]byte(h), []byte(raw))
}
