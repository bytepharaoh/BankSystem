package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	password := RandomString(8)
	wrongPassword := RandomString(8)

	hash1, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hash1)
	require.NotEqual(t, password, hash1)

	hash2, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hash2)

	// bcrypt adds a random salt, so the same password should hash differently each time.
	require.NotEqual(t, hash1, hash2)

	err = CheckPassword(password, hash1)
	require.NoError(t, err)
	require.NoError(t, CheckPassword(password, hash2))
	err = CheckPassword(wrongPassword, hash1)
	require.ErrorIs(t, err, bcrypt.ErrMismatchedHashAndPassword)
}
