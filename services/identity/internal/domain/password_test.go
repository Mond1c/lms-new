package domain

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPasswordHashValid(t *testing.T) {
	const validPassword = "12345678"

	hash, err := HashPassword(validPassword)
	require.NoError(t, err)
	require.NotEmpty(t, hash)
}

func TestPasswordHashInvalid(t *testing.T) {
	shortPassword := strings.Repeat("1", MinPasswordLength-1)
	longPassowrd := strings.Repeat("1", MaxPasswordLength+1)

	hash, err := HashPassword(shortPassword)
	require.Error(t, err)
	require.Equal(t, err, ErrPasswordTooShort)
	require.Empty(t, hash)

	hash, err = HashPassword(longPassowrd)
	require.Error(t, err)
	require.Equal(t, err, ErrPasswordTooLong)
	require.Empty(t, hash)
}

func TestPasswordHashVerify(t *testing.T) {
	const validPassword = "12345678"
	hash, err := HashPassword(validPassword)
	require.NoError(t, err)
	require.NotEmpty(t, hash)

	err = hash.Verify(validPassword)
	require.NoError(t, err)
}
