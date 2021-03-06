package routes

import (
  "database/sql"
  "encoding/json"
  "fmt"
  "net/http"
  "strconv"
  "time"

  "github.com/cakebin/smush/server/services/db"
  "github.com/cakebin/smush/server/services/email"
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
  UserID  int64  `json:"userId"`
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
  UserID  int64  `json:"userId"`
}


// ForgotPasswordRequestData describes the data we're expecting
// when a user requests to send an email to reset their password
type ForgotPasswordRequestData struct {
  UserEmail  string  `json:"userEmail"`
}

// ResetPasswordRequestData describes the data we're expecting
// when a user requests to reset their password
type ResetPasswordRequestData struct {
  Token        string  `json:"token"`
  NewPassword  string  `json:"newPassword"`
}

/*---------------------------------
          Response Data
----------------------------------*/

// LoginResponseData is the data we
// send back after a successful log in
type LoginResponseData struct {
  User               *db.UserProfileView      `json:"user"`
  UserCharacters     []*db.UserCharacterView  `json:"userCharacters"`
  AccessExpiration   time.Time             `json:"accessExpiration"`
  RefreshExpiration  time.Time             `json:"refreshExpiration"`
}


// LogoutResponseData is the data we
// send back after a successful log out
type LogoutResponseData struct {
  UserID  int64  `json:"userId"`
}


// RegisterResponseData is the data we
// send back after a successful register
type RegisterResponseData struct {
  UserID  int64  `json:"userId"`
}


// RefreshResponseData is the data we
// send back after a successful refresh
type RefreshResponseData struct {
  AccessExpiration  time.Time  `json:"accessExpiration"`
}


/*---------------------------------
             Router
----------------------------------*/

// AuthRouter handles all of the authentication related routes
type AuthRouter struct {
  Services  *Services
}


func (r *AuthRouter) ServeHTTP(res http.ResponseWriter, req *http.Request) {
  var head string
  head, req.URL.Path = ShiftPath(req.URL.Path)

  switch head {
  case "login":
    r.handleLogin(res, req)
  case "logout":
    r.handleLogout(res, req)
  case "register":
    r.handleRegister(res, req)
  case "refresh":
    r.handleRefresh(res, req)
  case "forgot-password":
    r.handleForgotPassword(res, req)
  case "reset-password":
    r.handleResetPassword(res, req)
  default:
    http.Error(res, "404 Not found", http.StatusNotFound)
  }
}


// NewAuthRouter makes a new api/auth router and hooks up its services
func NewAuthRouter(routerServices *Services) *AuthRouter {
  router := new(AuthRouter)

  router.Services = routerServices

  return router
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
  _, err = r.Services.Auth.CheckJWTToken(refreshCookie.Value)
  if err != nil {
    // Refresh token is also invalid. This block should never be run in practice.
    http.Error(res, fmt.Sprintf("Session expired. Please log in again. Token error: %s", err.Error()), http.StatusUnauthorized)
    return
  }

  // Whether or not we have a cookie, the new/updated one will need an expiration time of five minutes from now
  newExpirationTime := time.Now().Add(5 * time.Minute)

  // Check the access token to see if we need a new cookie. We DON'T need to check the access token value.
  // It won't have a value if it's already expired (we won't be sent one to update)
  accessCookie, err := req.Cookie("smush-access-token")

  if err != nil {
    // We DO NOT HAVE A COOKIE ANYMORE! So we need to make a new one.
    accessTokenStr, err := r.Services.Auth.GetNewJWTToken(refreshRequestData.UserID, newExpirationTime)
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
    newAccessToken, err := r.Services.Auth.RefreshJWTAccessToken(accessCookie.Value, newExpirationTime)
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

  userCredentialsView, err := r.Services.Database.GetUserCredentialsViewByEmail(loginRequestData.EmailAddress)
  if err == sql.ErrNoRows {
    http.Error(res, fmt.Sprintf("Invalid email address; user %s does not exist", loginRequestData.EmailAddress), http.StatusNotFound)
    return
  } else if err != nil {
    http.Error(res, fmt.Sprintf("Database error: %s", err.Error()), http.StatusInternalServerError)
    return
  }

  _, err = r.Services.Auth.CheckPassword(
    userCredentialsView.HashedPassword,
    loginRequestData.Password,
  )
  if err != nil {
    http.Error(res, "Invalid email/password", http.StatusUnauthorized)
    return
  }

  // Short lifespan access token
  accessExpiration := time.Now().Add(5 * time.Minute)
  accessTokenStr, err := r.Services.Auth.GetNewJWTToken(userCredentialsView.UserID, accessExpiration)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error creating new access token: %s", err.Error()), http.StatusInternalServerError)
    return
  }

  // Longer lifespan refresh token
  refreshExpiration := time.Now().Add(time.Hour * 24)
  refreshTokenStr, err := r.Services.Auth.GetNewJWTToken(userCredentialsView.UserID, refreshExpiration)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error creating new refresh token: %s", err.Error()), http.StatusInternalServerError)
    return
  }

  // Also store this refresh token in the user table
  userRefreshUpdate := new(db.UserRefreshUpdate)
  userRefreshUpdate.UserID = userCredentialsView.UserID
  userRefreshUpdate.RefreshToken = refreshTokenStr

  _, err = r.Services.Database.UpdateUserRefreshToken(userRefreshUpdate)
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

  // Get the basic user profile information
  userProfileView, err := r.Services.Database.GetUserProfileViewByUserID(userCredentialsView.UserID)
  if err != nil {
    http.Error(res, fmt.Sprintf("Could not get user data for id %d: %s", userCredentialsView.UserID, err.Error()), http.StatusBadRequest)
    return
  }

  // Also get the user's saved characters
  userCharViews, err := r.Services.Database.GetUserCharacterViewsByUserID(userCredentialsView.UserID)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error getting user's saved characters with userID %d: %s", userCredentialsView.UserID, err.Error()), http.StatusInternalServerError)
    return
  }

  // Finally get the user roles after authentication
  userRoleViews, err := r.Services.Database.GetUserRoleViewsByUserID(userCredentialsView.UserID)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error fetching user roles from database: %s", err.Error()), http.StatusInternalServerError)
    return
  }
  userProfileView.UserRoles = userRoleViews

  response := &Response{
    Success:           true,
    Error:             nil,
    Data:  LoginResponseData{
      User:               userProfileView,
      UserCharacters:     userCharViews,
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

  _, err = r.Services.Database.GetUserIDByEmail(registerRequestData.EmailAddress)
  if err != sql.ErrNoRows {
    http.Error(res, fmt.Sprintf("User already exists with email address %s", registerRequestData.EmailAddress), http.StatusBadRequest)
    return
  }

  hashedPassword, err := r.Services.Auth.HashPassword(registerRequestData.Password)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error when hashing password: %s", err.Error()), http.StatusInternalServerError)
    return
  }


  userCreate := new(db.UserCreate)
  userCreate.UserName = registerRequestData.UserName
  userCreate.EmailAddress = registerRequestData.EmailAddress
  userCreate.HashedPassword = hashedPassword
  userID, err := r.Services.Database.CreateUser(userCreate)
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
  userRefreshUpdate :=  new(db.UserRefreshUpdate)
  userRefreshUpdate.UserID = logoutRequestData.UserID
  userRefreshUpdate.RefreshToken = ""
  userID, err := r.Services.Database.UpdateUserRefreshToken(userRefreshUpdate)
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


func (r *AuthRouter) handleForgotPassword(res http.ResponseWriter, req *http.Request) {
  var forgotPasswordRequestData ForgotPasswordRequestData
  decoder := json.NewDecoder(req.Body)
  err := decoder.Decode(&forgotPasswordRequestData)
  if err != nil {
    http.Error(res, fmt.Sprintf("Invalid JSON request: %s", err.Error()), http.StatusBadRequest)
    return
  }

  // Check if we have a user with that email address before sending an email
  userID, err := r.Services.Database.GetUserIDByEmail(forgotPasswordRequestData.UserEmail)
  if err == sql.ErrNoRows {
    http.Error(res, fmt.Sprintf("No such user exists with email address %s", err.Error()), http.StatusInternalServerError)
    return
  }

  resetPasswordRequest, err := http.NewRequest("GET", "https://smush-tracker.herokuapp.com/reset-password/token", nil)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error when attempting to build reset password url: %s", err.Error()), http.StatusInternalServerError)
    return
  }

  queryParam := resetPasswordRequest.URL.Query()
  resetExpirationTime := time.Now().Add(15 * time.Minute)
  resetPasswordToken, err := r.Services.Auth.GetNewJWTToken(userID, resetExpirationTime)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error creating new access token: %s", err.Error()), http.StatusInternalServerError)
    return
  }

  // Update the user's reset_password_token for validation later
  resetPasswordUpdate := new(db.UserResetPasswordUpdate)
  resetPasswordUpdate.UserID = userID
  resetPasswordUpdate.ResetPasswordToken = resetPasswordToken
  userID, err = r.Services.Database.UpdateUserResetPasswordToken(resetPasswordUpdate)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error saving reset_password_token :%s for userID: %d", err.Error(), userID), http.StatusInternalServerError)
    return
  }

  queryParam.Add("t", resetPasswordToken)
  queryParam.Add("e", strconv.FormatInt(resetExpirationTime.Unix() * 1000, 10))
  resetPasswordRequest.URL.RawQuery = queryParam.Encode()
  resetURL := resetPasswordRequest.URL.String()

  resetPWInfo := new(email.ResetPWInfo)
  resetPWInfo.UserEmail = forgotPasswordRequestData.UserEmail
  resetPWInfo.ResetURL = resetURL
  success, err := r.Services.Email.SendResetPWEmail(resetPWInfo)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error when attempting to send email to: %s", err.Error()), http.StatusInternalServerError)
    return
  }

  response := &Response{
    Success:  success,
    Error:    nil,
  }

  res.Header().Set("Content-Type", "application/json")
  json.NewEncoder(res).Encode(response)
}


func (r *AuthRouter) handleResetPassword(res http.ResponseWriter, req *http.Request) {
  resetPasswordRequest := new(ResetPasswordRequestData)
  decoder := json.NewDecoder(req.Body)
  err := decoder.Decode(&resetPasswordRequest)
  if err != nil {
    http.Error(res, fmt.Sprintf("Invalid JSON request: %s", err.Error()), http.StatusBadRequest)
    return
  }

  // Check if the reset token has expired
  _, err = r.Services.Auth.CheckJWTToken(resetPasswordRequest.Token)
  if err != nil {
    http.Error(res, "Reset Password token has expired", http.StatusBadRequest)
    return
  }

  // Get the userID from the reset token
  userID, err := r.Services.Auth.GetUserIDFromJWTToken(resetPasswordRequest.Token)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error getting userID from reset token: %s", err.Error()), http.StatusInternalServerError)
    return
  }

  // Check if the reset token matches the one we stored for the user's database
  currentResetPasswordToken, err := r.Services.Database.GetUserResetPasswordTokenByUserID(userID)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error getting user's current reset token: %s", err.Error()), http.StatusInternalServerError)
    return
  }

  if currentResetPasswordToken != resetPasswordRequest.Token {
    http.Error(res, "Provided reset token does not match user's current reset token", http.StatusBadRequest)
    return
  }

  // If we're good, update the new password for the user with the usual hashing
  newHashedPassword, err := r.Services.Auth.HashPassword(resetPasswordRequest.NewPassword)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error when hashing new password: %s", err.Error()), http.StatusInternalServerError)
    return
  }
  hashedPasswordUpdate := new(db.UserHashedPasswordUpdate)
  hashedPasswordUpdate.UserID = userID
  hashedPasswordUpdate.HashedPassword = newHashedPassword
  userID, err = r.Services.Database.UpdateUserHashedPassword(hashedPasswordUpdate)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error when updating user's hashed password: %s", err.Error()), http.StatusInternalServerError)
    return
  }

  // Clear out the reset password token
  resetPasswordUpdate := new(db.UserResetPasswordUpdate)
  resetPasswordUpdate.UserID = userID
  resetPasswordUpdate.ResetPasswordToken = ""
  userID, err = r.Services.Database.UpdateUserResetPasswordToken(resetPasswordUpdate)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error clearing reset_password_token :%s for userID: %d", err.Error(), userID), http.StatusInternalServerError)
    return
  }

  // Send a success response, which would have the front end prompt them to log in
  response := &Response{
    Success:  true,
    Error:    nil,
  }

  res.Header().Set("Content-Type", "application/json")
  json.NewEncoder(res).Encode(response)
}
