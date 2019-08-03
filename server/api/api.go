package api

import (
	"time"

	"github.com/cakebin/smush/server/services/database"
)

// Response is a generic response, returned
// when sending POST requests to our API routes
type Response struct {
	Success bool        `json:"success"`
	Error   error       `json:"error"`
	Data    interface{} `json:"data"`
}

// AuthResponse is an api response related to auth endpoints (i.e login)
type AuthResponse struct {
	Success           bool                      `json:"success"`
	Error             error                     `json:"error"`
	User              *database.UserProfileView `json:"user"`
	AccessExpiration  time.Time                 `json:"accessExpiration"`
	RefreshExpiration time.Time                 `json:"refreshExpiration"`
}
