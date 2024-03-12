package userup

import (
	"crypto/rand"
	"fmt"
)

// NewSessionID creates a new random ID for a session ID.
func NewSessionID() string {
	// We're not handling the error b/c it is not something we can
	// deal with. It also isn't clear what the failure modes are.
	p, _ := rand.Prime(rand.Reader, 64)
	return fmt.Sprintf("%x", p)
}
