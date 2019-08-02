package auth

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/cakebin/smush/server/api"
	"github.com/cakebin/smush/server/env"
	"github.com/cakebin/smush/server/services/database"
	"github.com/cakebin/smush/server/util/routing"
)

// Router handles all of the authentication related routes
type Router struct {
	SysUtils *env.SysUtils
}

func (r *Router) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var head string
	head, req.URL.Path = routing.ShiftPath(req.URL.Path)

	switch head {
	case "login":
		r.handleLogin(res, req)
	case "logout":
		r.handleLogout(res, req)
	case "register":
		r.handleRegister(res, req)
	case "refresh":
		r.handleRefresh(res, req)
	default:
		http.Error(res, "404 Not found", http.StatusNotFound)
	}
}

func (r *Router) handleRefresh(res http.ResponseWriter, req *http.Request) {
	// We need the user to create a new token in case the access token is gone
	decoder := json.NewDecoder(req.Body)
	var userRefresh database.UserRefreshUpdate

	err := decoder.Decode(&userRefresh)
	if err != nil {
		http.Error(res, fmt.Sprintf("Invalid JSON request: %s", err.Error()), http.StatusBadRequest)
		return
	}

	// Check the refresh token. We don't care if there's an access token or not.
	// The refresh token should not be expired, because this endpoint is only being hit
	// BEFORE the refresh token expires. If it's already expired, the front end
	// will instead take care of the logout and prompt another login.
	refreshCookie, err := req.Cookie("smush-refresh-token")
	_, err = r.SysUtils.Authenticator.CheckJWTToken(refreshCookie.Value)
	if err != nil {
		// Refresh token is also invalid. This block should never be run in practice.
		http.Error(res, "Session expired. Please log in again", http.StatusUnauthorized)
		return
	}

	// Whether or not we have a cookie, the new/updated one will need an expiration time of five minutes from now
	newExpirationTime := time.Now().Add(5 * time.Minute)

	// Check the access token to see if we need a new cookie. We DON'T need to check the access token value.
	// It won't have a value if it's already expired (we won't be sent one to update)
	accessCookie, err := req.Cookie("smush-access-token")

	if err != nil {
		// We do NOT HAVE A COOKIE ANYMORE! So we need to make a new one.
		accessTokenStr, err := r.SysUtils.Authenticator.GetNewJWTToken(userRefresh.UserID, newExpirationTime)
		if err != nil {
			http.Error(res, fmt.Sprintf("Error creating new access token: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		http.SetCookie(
			res,
			&http.Cookie{
				Name:    "smush-access-token",
				Value:   accessTokenStr,
				Expires: newExpirationTime,
				Path:    "/api/",
			},
		)
	} else {
		// We DO HAVE A COOKIE! Update the existing one!
		newAccessToken, err := r.SysUtils.Authenticator.RefreshJWTAccessToken(accessCookie.Value, newExpirationTime)
		if err != nil {
			http.Error(res, fmt.Sprintf("Error updating existing access token: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		http.SetCookie(
			res,
			&http.Cookie{
				Name:    "smush-access-token",
				Value:   newAccessToken,
				Expires: newExpirationTime,
				Path:    "/api/",
			},
		)
	}

	// We are finally done! Send a new AuthResponse with the updated expiration time
	response := &api.AuthResponse{
		Success:           true,
		Error:             nil,
		User:              nil,
		AccessExpiration:  newExpirationTime,
		RefreshExpiration: refreshCookie.Expires,
	}

	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(response)
}

func (r *Router) handleLogin(res http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var credentialsRequest database.UserCredentialsView

	err := decoder.Decode(&credentialsRequest)
	if err != nil {
		http.Error(res, fmt.Sprintf("Invalid JSON request: %s", err.Error()), http.StatusBadRequest)
		return
	}

	userCredentialsView, err := r.SysUtils.Database.GetUserCredentialsViewByEmail(credentialsRequest.EmailAddress)
	if err == sql.ErrNoRows {
		http.Error(res, fmt.Sprintf("Invalid email address; user %s does not exist", credentialsRequest.EmailAddress), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(res, fmt.Sprintf("Database error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	err = r.SysUtils.Authenticator.CheckPassword(
		*userCredentialsView.HashedPassword,
		*credentialsRequest.Password,
	)
	if err != nil {
		http.Error(res, "Invalid email/password", http.StatusUnauthorized)
		return
	}

	// Short lifespan access token
	accessExpiration := time.Now().Add(5 * time.Minute)
	accessTokenStr, err := r.SysUtils.Authenticator.GetNewJWTToken(*userCredentialsView.UserID, accessExpiration)
	if err != nil {
		http.Error(res, fmt.Sprintf("Error creating new access token: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Longer lifespan refresh token
	refreshExpiration := time.Now().Add(time.Hour * 24)
	refreshTokenStr, err := r.SysUtils.Authenticator.GetNewJWTToken(*userCredentialsView.UserID, refreshExpiration)
	if err != nil {
		http.Error(res, fmt.Sprintf("Error creating new refresh token: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	// Also store this refresh token in the user table
	_, err = r.SysUtils.Database.UpdateUserRefreshToken(refreshTokenStr, *userCredentialsView.UserID)
	if err != nil {
		http.Error(res, fmt.Sprintf("Error adding new refresh token to database: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	http.SetCookie(
		res,
		&http.Cookie{
			Name:    "smush-access-token",
			Value:   accessTokenStr,
			Expires: accessExpiration,
			Path:    "/api/",
		},
	)
	http.SetCookie(
		res,
		&http.Cookie{
			Name:    "smush-refresh-token",
			Value:   refreshTokenStr,
			Expires: refreshExpiration,
			Path:    "/api/",
		},
	)

	userProfileView, err := r.SysUtils.Database.GetUserProfileViewByID(*userCredentialsView.UserID)
	if err != nil {
		http.Error(res, fmt.Sprintf("Could not get user data: %s", err.Error()), http.StatusBadRequest)
		return
	}

	response := &api.AuthResponse{
		Success:           true,
		Error:             nil,
		User:              userProfileView,
		AccessExpiration:  accessExpiration,
		RefreshExpiration: refreshExpiration,
	}

	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(response)
}

func (r *Router) handleRegister(res http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var newUser database.User

	err := decoder.Decode(&newUser)
	if err != nil {
		http.Error(res, fmt.Sprintf("Invalid JSON request: %s", err.Error()), http.StatusBadRequest)
		return
	}

	_, err = r.SysUtils.Database.GetUserByEmail(newUser.EmailAddress)
	if err != sql.ErrNoRows {
		http.Error(res, fmt.Sprintf("User already exists with email address %s", newUser.EmailAddress), http.StatusBadRequest)
		return
	}

	hashedPassword, err := r.SysUtils.Authenticator.HashPassword(newUser.Password)
	if err != nil {
		http.Error(res, fmt.Sprintf("Error when hashing password: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	newUser.HashedPassword = hashedPassword

	userID, err := r.SysUtils.Database.CreateUser(newUser)
	if err != nil {
		http.Error(res, fmt.Sprintf("Error creating new user in database: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	response := &api.Response{
		Success: true,
		Error:   nil,
		Data:    userID,
	}

	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(response)
}

func (r *Router) handleLogout(res http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var userLogout database.User

	err := decoder.Decode(&userLogout)
	if err != nil {
		http.Error(res, fmt.Sprintf("Invalid JSON request: %s", err.Error()), http.StatusBadRequest)
		return
	}

	// We want to delete the refresh token when a user logs out
	userID, err := r.SysUtils.Database.UpdateUserRefreshToken("", *userLogout.UserID)
	if err != nil {
		http.Error(res, fmt.Sprintf("Error removing refresh token from database: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Delete existing Access/Refresh tokens in cookies
	http.SetCookie(
		res,
		&http.Cookie{
			Name:   "smush-access-token",
			Value:  "",
			MaxAge: 0,
		},
	)
	http.SetCookie(
		res,
		&http.Cookie{
			Name:   "smush-refresh-token",
			Value:  "",
			MaxAge: 0,
		},
	)

	response := &api.Response{
		Success: true,
		Error:   nil,
		Data:    userID,
	}

	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(response)
}

// NewRouter makes a new match router with access to the "SysUtils" environment object
func NewRouter(sysUtils *env.SysUtils) *Router {
	router := new(Router)
	router.SysUtils = sysUtils
	return router
}
