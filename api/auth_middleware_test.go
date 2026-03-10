package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bytepharoh/simplebank/token"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type stubTokenMaker struct {
	verifyTokenFn func(token string) (*token.Payload, error)
}

func (s stubTokenMaker) CreateToken(username string, duration time.Duration) (string, *token.Payload, error) {
	return "", nil, errors.New("not implemented")
}

func (s stubTokenMaker) VerifyToken(token string) (*token.Payload, error) {
	return s.verifyTokenFn(token)
}

func TestAuthMiddleware(t *testing.T) {
	t.Parallel()

	makeServer := func(maker token.Maker) *gin.Engine {
		router := gin.New()
		router.GET("/auth", authMiddleware(maker), func(ctx *gin.Context) {
			payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
			ctx.JSON(http.StatusOK, gin.H{"username": payload.Username})
		})
		return router
	}

	testCases := []struct {
		name          string
		authHeader    string
		maker         token.Maker
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:       "MissingAuthorizationHeader",
			authHeader: "",
			maker: stubTokenMaker{
				verifyTokenFn: func(tok string) (*token.Payload, error) {
					return nil, nil
				},
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
				require.Contains(t, recorder.Body.String(), "authorization header is not provided")
			},
		},
		{
			name:       "InvalidAuthorizationHeaderFormat",
			authHeader: "Bearer",
			maker: stubTokenMaker{
				verifyTokenFn: func(tok string) (*token.Payload, error) {
					return nil, nil
				},
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
				require.Contains(t, recorder.Body.String(), "invalid authorization header format")
			},
		},
		{
			name:       "UnsupportedAuthorizationType",
			authHeader: "Basic token",
			maker: stubTokenMaker{
				verifyTokenFn: func(tok string) (*token.Payload, error) {
					return nil, nil
				},
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
				require.Contains(t, recorder.Body.String(), "unsupported authorization type")
			},
		},
		{
			name:       "InvalidToken",
			authHeader: "Bearer invalid-token",
			maker: stubTokenMaker{
				verifyTokenFn: func(tok string) (*token.Payload, error) {
					return nil, token.ErrInvalidToken
				},
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
				require.Contains(t, recorder.Body.String(), token.ErrInvalidToken.Error())
			},
		},
		{
			name:       "ExpiredToken",
			authHeader: "Bearer expired-token",
			maker: stubTokenMaker{
				verifyTokenFn: func(tok string) (*token.Payload, error) {
					return nil, token.ErrExpiredToken
				},
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
				require.Contains(t, recorder.Body.String(), token.ErrExpiredToken.Error())
			},
		},
		{
			name:       "OK",
			authHeader: "Bearer valid-token",
			maker: stubTokenMaker{
				verifyTokenFn: func(tok string) (*token.Payload, error) {
					return &token.Payload{
						ID:        uuid.New(),
						Username:  "valid_user",
						IssuedAt:  time.Now(),
						ExpiredAt: time.Now().Add(time.Minute),
					}, nil
				},
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				require.Contains(t, recorder.Body.String(), "valid_user")
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			server := makeServer(tc.maker)
			recorder := httptest.NewRecorder()

			request, err := http.NewRequest(http.MethodGet, "/auth", nil)
			require.NoError(t, err)
			if tc.authHeader != "" {
				request.Header.Set(authorizationHeaderKey, tc.authHeader)
			}

			server.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
