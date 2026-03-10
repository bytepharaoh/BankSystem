package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var ErrExpiredToken = errors.New("token has expired")

// Payload contains the token data.
type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
	jwt.RegisteredClaims
}

// NewPayload creates a new token payload for a specific username and duration.
func NewPayload(username string, duration time.Duration) (*Payload, error) {
	if username == "" {
		return nil, errors.New("username is required")
	}
	if duration <= 0 {
		return nil, errors.New("duration must be greater than zero")
	}

	now := time.Now()
	payload := &Payload{
		ID:        uuid.New(),
		Username:  username,
		IssuedAt:  now,
		ExpiredAt: now.Add(duration),
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
			Subject:   username,
		},
	}

	return payload, nil
}

// Valid checks whether the token payload is still valid.
func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
