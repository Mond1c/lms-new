package domain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEmailNewValid(t *testing.T) {
	const validEmail = "admin@gmail.com"
	email, err := NewEmail(validEmail)
	require.NoError(t, err)
	require.Equal(t, validEmail, email.String())
}

func TestEmailNewInvalid(t *testing.T) {
	const invalidEmail = "admingmail.com"
	email, err := NewEmail(invalidEmail)
	require.Error(t, err)
	require.Empty(t, email.String())
}

func TestEmailFromTrusted(t *testing.T) {
	const trustedEmail = "admin@gmail.com"
	email := EmailFromTrusted(trustedEmail)
	require.Equal(t, trustedEmail, email.String())
}
