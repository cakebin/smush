package api

import (
  "database/sql"
  "encoding/json"
  "fmt"
  "net/http"
  "time"

  "github.com/cakebin/smush/server/env"
  "github.com/cakebin/smush/server/services/db"
  "github.com/cakebin/smush/server/util/routing"
)


/*---------------------------------
          Request Data
----------------------------------*/

// LoginRequestData describes the data we're 
// expecting when a user attempts to log in
type LoginRequestData struct {
  EmailAddress  string  `json:"emailAddress"`
  Password      string  `json:"password"`
}


// LogoutRequestData describes the data we're 
// expecting when a user attempts to log out
type LogoutRequestData struct {
  UserID  int  `json:"userId"`
}


// RegisterRequestData describes the data we're 
// expecting when a user attempts register
type RegisterRequestData struct {
  UserName      string  `json:"userName"`
  EmailAddress  string  `json:"emailAddress"`
  Password      string  `json:"password"`
}


// RefreshRequestData describes the data we're expecting
// when a user attempts refresh their access token
type RefreshRequestData struct {
  UserID  int  `json:"userId"`
}


/*---------------------------------
          Response Data
----------------------------------*/

// LoginResponseData is the data we
// send back after a successful log in
type LoginResponseData struct {
  User               UserProfileView     `json:"user"`
  AccessExpiration   time.Time           `json:"accessExpiration"`
  RefreshExpiration  time.Time           `json:"refreshExpiration"`
}


// LogoutResponseData is the data we
// send back after a successful log out
type LogoutResponseData struct {
  UserID  int  `json:"userId"`
}


// RegisterResponseData is the data we
// send back after a successful register
type RegisterResponseData struct {
  UserID  int  `json:"userId"`
}


// RefreshResponseData is the data we
// send back after a successful refresh
type RefreshResponseData struct {
  AccessExpiration   time.Time  `json:"accessExpiration"`
}


/*---------------------------------
             Router
----------------------------------*/

// AuthRouter handles all of the authentication related routes
type AuthRouter struct {
  SysUtils  *env.SysUtils
}


func (r *AuthRouter) ServeHTTP(res http.ResponseWriter, req *http.Request) {
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


/*---------------------------------
             Handlers
----------------------------------*/
func (r *AuthRouter) handleRefresh(res http.ResponseWriter, req *http.Request) {
  // We need the user to create a new token in case the access token is gone
  var refreshRequestData RefreshRequestData
  decoder := json.NewDecoder(req.Body)
  err := decoder.Decode(&refreshRequestData)
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
    // We DO NOT HAVE A COOKIE ANYMORE! So we need to make a new one.
    accessTokenStr, err := r.SysUtils.Authenticator.GetNewJWTToken(refreshRequestData.UserID, newExpirationTime)
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

  // We are finally done! Send a new Response with the updated expiration time
  response := &Response{
    Success:  true,
    Error:    nil,
    Data:     RefreshResponseData{
      AccessExpiration:  newExpirationTime,
    },
  }

  res.Header().Set("Content-Type", "application/json")
  json.NewEncoder(res).Encode(response)
}


func (r *AuthRouter) handleLogin(res http.ResponseWriter, req *http.Request) {
  var loginRequestData LoginRequestData
  decoder := json.NewDecoder(req.Body)
  err := decoder.Decode(&loginRequestData)
  if err != nil {
    http.Error(res, fmt.Sprintf("Invalid JSON request: %s", err.Error()), http.StatusBadRequest)
    return
  }

  userCredentialsView, err := r.SysUtils.Database.GetUserCredentialsViewByEmail(loginRequestData.EmailAddress)
  if err == sql.ErrNoRows {
    http.Error(res, fmt.Sprintf("Invalid email address; user %s does not exist", loginRequestData.EmailAddress), http.StatusNotFound)
    return
  } else if err != nil {
    http.Error(res, fmt.Sprintf("Database error: %s", err.Error()), http.StatusInternalServerError)
    return
  }

  err = r.SysUtils.Authenticator.CheckPassword(
    userCredentialsView.HashedPassword,
    loginRequestData.Password,
  )
  if err != nil {
    http.Error(res, "Invalid email/password", http.StatusUnauthorized)
    return
  }

  // Short lifespan access token
  accessExpiration := time.Now().Add(5 * time.Minute)
  accessTokenStr, err := r.SysUtils.Authenticator.GetNewJWTToken(userCredentialsView.UserID, accessExpiration)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error creating new access token: %s", err.Error()), http.StatusInternalServerError)
    return
  }

  // Longer lifespan refresh token
  refreshExpiration := time.Now().Add(time.Hour * 24)
  refreshTokenStr, err := r.SysUtils.Authenticator.GetNewJWTToken(userCredentialsView.UserID, refreshExpiration)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error creating new refresh token: %s", err.Error()), http.StatusInternalServerError)
    return
  }

  // Also store this refresh token in the user table
  var userRefreshUpdate db.UserRefreshUpdate
  userRefreshUpdate.UserID = userCredentialsView.UserID
  userRefreshUpdate.RefreshToken = refreshTokenStr
  _, err = r.SysUtils.Database.UpdateUserRefreshToken(userRefreshUpdate)
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

  dbUserProfileView, err := r.SysUtils.Database.GetUserProfileViewByID(userCredentialsView.UserID)
  if err != nil {
    http.Error(res, fmt.Sprintf("Could not get user data for id %d: %s", userCredentialsView.UserID, err.Error()), http.StatusBadRequest)
    return
  }
  
  var userProfileView UserProfileView
  userProfileView.UserID = dbUserProfileView.UserID
  userProfileView.UserName = dbUserProfileView.UserName
  userProfileView.EmailAddress = dbUserProfileView.EmailAddress
  userProfileView.Created = dbUserProfileView.Created
  userProfileView.DefaultCharacterGsp = FromNullInt64(dbUserProfileView.DefaultCharacterGsp)
  userProfileView.DefaultCharacterID = FromNullInt64(dbUserProfileView.DefaultCharacterID)
  userProfileView.DefaultCharacterName = FromNullString(dbUserProfileView.DefaultCharacterName)

  response := &Response{
    Success:           true,
    Error:             nil,
    Data:  LoginResponseData{
      User:               userProfileView,
      AccessExpiration:   accessExpiration,
      RefreshExpiration:  refreshExpiration,
    },
  }

  res.Header().Set("Content-Type", "application/json")
  json.NewEncoder(res).Encode(response)
}


func (r *AuthRouter) handleRegister(res http.ResponseWriter, req *http.Request) {
  var registerRequestData RegisterRequestData
  decoder := json.NewDecoder(req.Body)
  err := decoder.Decode(&registerRequestData)
  if err != nil {
    http.Error(res, fmt.Sprintf("Invalid JSON request: %s", err.Error()), http.StatusBadRequest)
    return
  }

  _, err = r.SysUtils.Database.GetUserIDByEmail(registerRequestData.EmailAddress)
  if err != sql.ErrNoRows {
    http.Error(res, fmt.Sprintf("User already exists with email address %s", registerRequestData.EmailAddress), http.StatusBadRequest)
    return
  }

  hashedPassword, err := r.SysUtils.Authenticator.HashPassword(registerRequestData.Password)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error when hashing password: %s", err.Error()), http.StatusInternalServerError)
    return
  }

  var newUser db.User
  newUser.UserName = registerRequestData.UserName
  newUser.EmailAddress = registerRequestData.EmailAddress
  newUser.HashedPassword = hashedPassword
  userID, err := r.SysUtils.Database.CreateUser(newUser)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error creating new user in database: %s", err.Error()), http.StatusInternalServerError)
    return
  }

  response := &Response{
    Success: true,
    Error:   nil,
    Data:    RegisterResponseData{
      UserID: userID,
    },
  }

  res.Header().Set("Content-Type", "application/json")
  json.NewEncoder(res).Encode(response)
}


func (r *AuthRouter) handleLogout(res http.ResponseWriter, req *http.Request) {
  decoder := json.NewDecoder(req.Body)
  var logoutRequestData LogoutRequestData

  err := decoder.Decode(&logoutRequestData)
  if err != nil {
    http.Error(res, fmt.Sprintf("Invalid JSON request: %s", err.Error()), http.StatusBadRequest)
    return
  }

  // We want to delete the refresh token when a user logs out
  var userRefreshUpdate db.UserRefreshUpdate
  userRefreshUpdate.UserID = logoutRequestData.UserID
  userRefreshUpdate.RefreshToken = ""
  userID, err := r.SysUtils.Database.UpdateUserRefreshToken(userRefreshUpdate)
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
      MaxAge: -1,
    },
  )
  http.SetCookie(
    res,
    &http.Cookie{
      Name:   "smush-refresh-token",
      Value:  "",
      MaxAge: -1,
    },
  )

  response := &Response{
    Success: true,
    Error:   nil,
    Data:    LogoutResponseData{
      UserID: userID,
    },
  }

  res.Header().Set("Content-Type", "application/json")
  json.NewEncoder(res).Encode(response)
}
