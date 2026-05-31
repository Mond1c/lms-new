package jwt

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const userID = "01HX"
const secret = "secret"

func TestSignAndVerify(t *testing.T) {
	signer := NewSinger(secret, time.Hour)
	verifier := NewVerifier(secret)

	token, err := signer.Sign(userID)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	claims, err := verifier.Verify(token)
	require.NoError(t, err)
	require.Equal(t, userID, claims.UserID)
}

func TestVerify_WrongSecret(t *testing.T) {
	const wrongSecret = "wrongSecret"

	signer := NewSinger(secret, time.Hour)
	verifier := NewVerifier(wrongSecret)

	token, _ := signer.Sign(userID)
	_, err := verifier.Verify(token)
	require.Error(t, err)
}

func TestVerify_Expired(t *testing.T) {
	signer := NewSinger(secret, -time.Hour)
	verifier := NewVerifier(secret)

	token, _ := signer.Sign(userID)
	_, err := verifier.Verify(token)
	require.Error(t, err)
}

func TestVerify_Garbage(t *testing.T) {
	verifier := NewVerifier(secret)
	_, err := verifier.Verify("not-a-jwt")
	require.Error(t, err)

	_, err = verifier.Verify(strings.Repeat("x", 100))
	require.Error(t, err)
}
