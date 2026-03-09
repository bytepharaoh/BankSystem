package toeken

import (
	"testing"
	"time"

	"github.com/bytepharoh/simplebank/util"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	secretKey := util.RandomString(32)
	username := util.RandomOwner()

	t.Run("new maker with short key", func(t *testing.T) {
		maker, err := NewJWTMaker("short-key")
		require.Error(t, err)
		require.Nil(t, maker)
	})

	t.Run("create and verify token", func(t *testing.T) {
		maker, err := NewJWTMaker(secretKey)
		require.NoError(t, err)

		token, payload, err := maker.CreateToken(username, time.Minute)
		require.NoError(t, err)
		require.NotEmpty(t, token)
		require.NotNil(t, payload)

		verifiedPayload, err := maker.VerifyToken(token)
		require.NoError(t, err)
		require.Equal(t, payload.ID, verifiedPayload.ID)
		require.Equal(t, payload.Username, verifiedPayload.Username)
	})

	t.Run("expired token", func(t *testing.T) {
		maker, err := NewJWTMaker(secretKey)
		require.NoError(t, err)

		jwtMaker := maker.(*JWTMaker)
		payload, err := NewPayload(username, time.Minute)
		require.NoError(t, err)
		past := time.Now().Add(-time.Minute)
		payload.ExpiredAt = past
		payload.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(past)

		jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
		token, err := jwtToken.SignedString([]byte(jwtMaker.secretKey))
		require.NoError(t, err)

		gotPayload, err := maker.VerifyToken(token)
		require.ErrorIs(t, err, ErrExpiredToken)
		require.Nil(t, gotPayload)
	})

	t.Run("invalid token", func(t *testing.T) {
		maker, err := NewJWTMaker(secretKey)
		require.NoError(t, err)

		payload, err := maker.VerifyToken("invalid-token")
		require.ErrorIs(t, err, ErrInvalidToken)
		require.Nil(t, payload)
	})
}
