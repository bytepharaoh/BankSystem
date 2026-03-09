package toeken

import "time"

// Maker defines behavior for creating and verifying authentication tokens.
type Maker interface {
	CreateToken(username string, duration time.Duration) (string, *Payload, error)
	VerifyToken(token string) (*Payload, error)
}
