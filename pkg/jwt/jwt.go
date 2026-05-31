package jwt

import (
	"errors"
	"fmt"
	"time"

	jwtv5 "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string `json:"sub"`
	jwtv5.RegisteredClaims
}

type Signer struct {
	secret []byte
	ttl    time.Duration
}

func NewSinger(secret string, ttl time.Duration) *Signer {
	return &Signer{
		secret: []byte(secret),
		ttl:    ttl,
	}
}

func (s *Signer) Sign(userID string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwtv5.RegisteredClaims{
			ExpiresAt: jwtv5.NewNumericDate(now.Add(s.ttl)),
			IssuedAt:  jwtv5.NewNumericDate(now),
			Issuer:    "lms",
		},
	}

	token := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.secret)
	if err != nil {
		return "", fmt.Errorf("sign: %w", err)
	}
	return signed, nil
}

func (s *Signer) TTL() time.Duration { return s.ttl }

var ErrInvalidToken = errors.New("invalid token")

type Verifier struct {
	secret []byte
}

func NewVerifier(secret string) *Verifier {
	return &Verifier{secret: []byte(secret)}
}

func (v *Verifier) Verify(raw string) (*Claims, error) {
	token, err := jwtv5.ParseWithClaims(raw, &Claims{}, func(token *jwtv5.Token) (any, error) {
		if _, ok := token.Method.(*jwtv5.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return v.secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}
	return claims, nil
}
