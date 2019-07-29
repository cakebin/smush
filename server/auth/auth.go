package auth

import (
  "time"
)


// Authenticator blah
type Authenticator interface {
  // JWT Tokens
  GetNewJWTToken(emailAddress string, expiration time.Time) (string, error)
  RefreshJWTAccessToken(token string, expiration time.Time) (string, error)
  CheckJWTToken(token string) (bool, error)

  // Encryption
  HashPassword(password string) (string, error)
  CheckPassword(hashed string, password string) error 
}


// Auth blah
type Auth struct {}


// NewAuthenticator makes a new authenticator, duh
func NewAuthenticator() (*Auth) {
  return &Auth{}
}
