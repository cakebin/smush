package api

import (
	"time"

	"github.com/cakebin/smush/server/db"
)

// Response is a generic response, returned
// when sending POST requests to our API routes
type Response struct {
	Success bool  `json:"success"`
	Error   error `json:"error"`
}

// AuthResponse is an api response related to auth endpoints (i.e login)
type AuthResponse struct {
	Success           bool      `json:"success"`
	Error             error     `json:"error"`
	User              *db.User  `json:"user"`
	AccessExpiration  time.Time `json:"accessExpiration"`
	RefreshExpiration time.Time `json:"refreshExpiration"`
}
