package auth

import (
  "database/sql"
  "encoding/json"
  "fmt"
  "net/http"
  "time"

  "github.com/cakebin/smush/server/api"
  "github.com/cakebin/smush/server/db"
  "github.com/cakebin/smush/server/env"
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
  default:
    http.Error(res, "404 Not found", http.StatusNotFound)
  }
}


func (r *Router) handleLogin(res http.ResponseWriter, req *http.Request) {
  decoder := json.NewDecoder(req.Body)
  var userLogin db.User

  err := decoder.Decode(&userLogin)
  if err != nil {
    http.Error(res, fmt.Sprintf("Invalid JSON request: %s", err.Error()), http.StatusBadRequest)
    return
  }

  user, err := r.SysUtils.Database.GetUserByEmail(userLogin.EmailAddress)
  if err == sql.ErrNoRows {
    http.Error(res, fmt.Sprintf("Invalid email address; user %s does not exist", userLogin.EmailAddress), http.StatusNotFound)
    return
  } else if err != nil {
    http.Error(res, fmt.Sprintf("Database error: %s", err.Error()), http.StatusInternalServerError)
    return
  }

  err = r.SysUtils.Authenticator.CheckPassword(
    *user.HashedPassword,
    *userLogin.Password,
  )
  if err != nil {
    http.Error(res, "Invalid email/password", http.StatusUnauthorized)
    return
  }

  // Short lifespan access token
  accessExpiration := time.Now().Add(5 * time.Minute)
  accessTokenStr, err := r.SysUtils.Authenticator.GetNewJWTToken(user.EmailAddress, accessExpiration)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error creating new access token: %s", err.Error()), http.StatusInternalServerError)
    return
  }

  // Longer lifespan refresh token
  refreshExpiration := time.Now().Add(time.Hour * 24)
  refreshTokenStr, err := r.SysUtils.Authenticator.GetNewJWTToken(user.EmailAddress, refreshExpiration)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error creating new refresh token: %s", err.Error()), http.StatusInternalServerError)
    return
  }
  // Also store this refresh token in the user table
  _, err = r.SysUtils.Database.UpdateUserRefreshTokenByID(refreshTokenStr, *user.UserID)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error adding new refresh token to database: %s", err.Error()), http.StatusInternalServerError)
    return
  }

  http.SetCookie(
    res,
    &http.Cookie{
      Name: "smush-access-token",
      Value: accessTokenStr,
      Expires: accessExpiration,
    },
  )
  http.SetCookie(
    res,
    &http.Cookie{
      Name: "smush-refresh-token",
      Value: refreshTokenStr,
      Expires: refreshExpiration,
    },
  )

  user.HashedPassword = nil
  user.RefreshToken = nil

  response := &api.AuthResponse{
    Success: true,
    Error: nil,
    User: user,
  }

  res.Header().Set("Content-Type", "application/json")
  json.NewEncoder(res).Encode(response)
}



func (r *Router) handleRegister(res http.ResponseWriter, req *http.Request) {
  decoder := json.NewDecoder(req.Body)
  var newUser db.User

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

  hashedPassword, err := r.SysUtils.Authenticator.HashPassword(*newUser.Password)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error when hashing password: %s", err.Error()), http.StatusInternalServerError)
    return
  }
  newUser.HashedPassword = &hashedPassword

  _, err = r.SysUtils.Database.CreateUser(newUser)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error creating new user in database: %s", err.Error()), http.StatusInternalServerError)
    return
  }

  response := &api.Response{
    Success: true,
    Error: nil,
  }

  res.Header().Set("Content-Type", "application/json")
  json.NewEncoder(res).Encode(response)
}



func (r *Router) handleLogout(res http.ResponseWriter, req *http.Request) {
  decoder := json.NewDecoder(req.Body)
  var userLogout db.User

  err := decoder.Decode(&userLogout)
  if err != nil {
    http.Error(res, fmt.Sprintf("Invalid JSON request: %s", err.Error()), http.StatusBadRequest)
    return
  }

  // We want to delete the refresh token when a user logs out
  _, err = r.SysUtils.Database.UpdateUserRefreshTokenByID("", *userLogout.UserID)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error removing refresh token from database: %s", err.Error()), http.StatusInternalServerError)
    return
  }

  // Delete existing Access/Refresh tokens in cookies
  http.SetCookie(
    res,
    &http.Cookie{
      Name: "smush-access-token",
      Value: "",
      MaxAge: 0,
    },
  )
  http.SetCookie(
    res,
    &http.Cookie{
      Name: "smush-refresh-token",
      Value: "",
      MaxAge: 0,
    },
  )

  response := &api.Response{
    Success: true,
    Error: nil,
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
