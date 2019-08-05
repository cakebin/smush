package api

import (
  "encoding/json"
  "fmt"
  "net/http"
  "strconv"
  "time"

  "github.com/cakebin/smush/server/env"
  "github.com/cakebin/smush/server/services/db"
  "github.com/cakebin/smush/server/util/routing"
)


/*---------------------------------
          Request Data
----------------------------------*/

// UserUpdateRequestData describes the data we're 
// expecting when a user attempts to update their profile
type UserUpdateRequestData struct {
  UserID               int       `json:"userId"`
  EmailAddress         string    `json:"emailAddress"`
  UserName             string    `json:"userName"`
  DefaultCharacterID   *int64    `json:"defaultCharacterId,omitempty"`
  DefaultCharacterGsp  *int64    `json:"defaultCharacterGsp,omitempty"`
}



/*---------------------------------
          SQL --> API
----------------------------------*/

// UserProfileView is a translation from the SQL result
// which can have things like `sql.NullInt64`, so we 
// need to translate that to regular JSON objects
type UserProfileView struct {
  UserID                int        `json:"userId"`
  UserName              string     `json:"userName"`
  EmailAddress          string     `json:"emailAddress"`
  Created               time.Time  `json:"created"`
  DefaultCharacterGsp   *int64     `json:"defaultCharacterGsp,omitempty"`
  DefaultCharacterID    *int64     `json:"defaultCharacterId,omitempty"`
  DefaultCharacterName  *string    `json:"defaultCharacterName,omitempty"`
}


/*---------------------------------
          Response Data
----------------------------------*/



// UserGetResponseData is the data we send back
// after a successfully getting all user's info
type UserGetResponseData struct {
  User  UserProfileView  `json:"user"`
}


// UserUpdateResponseData is the data we send
// back after a successfully creating a new user
type UserUpdateResponseData struct {
  UserID  int  `json:"userId"`
}


/*---------------------------------
             Router
----------------------------------*/

// UserRouter is responsible for serving "/api/user"
// Basically, connecting to our Postgres DB for all
// of the CRUD operations for our "User" models
type UserRouter struct {
  SysUtils  *env.SysUtils
}


func (r *UserRouter) ServeHTTP(res http.ResponseWriter, req *http.Request) {
  var head string
  head, req.URL.Path = routing.ShiftPath(req.URL.Path)

  switch req.Method {
  // GET Request Handlers
  case http.MethodGet:
    switch head {
    case "get":
      r.handleGetByID(res, req)
    default:
      http.Error(res, fmt.Sprintf("Unsupported GET path %s", head), http.StatusBadRequest)
      return
    }
  // POST Request Handlers
  case http.MethodPost:
    switch head {
    case "update":
      r.handleUpdate(res, req)
    default:
      http.Error(res, fmt.Sprintf("Unsupport POST path %s", head), http.StatusBadRequest)
      return
    }
  // Unsupported Method Response
  default:
    http.Error(res, fmt.Sprintf("Unsupported Method type %s", req.Method), http.StatusBadRequest)
  }
}


/*---------------------------------
             Handlers
----------------------------------*/

func (r *UserRouter) handleGetByID(res http.ResponseWriter, req *http.Request) {
  var head string
  head, req.URL.Path = routing.ShiftPath(req.URL.Path)

  id, err := strconv.Atoi(head)
  if err != nil {
    http.Error(res, fmt.Sprintf("Invalid user id: %s", head), http.StatusBadRequest)
    return
  }

  dbUserProfileView, err := r.SysUtils.Database.GetUserProfileViewByID(id)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error getting user with id %q: %s", id, err.Error()), http.StatusInternalServerError)
    return
  }

  var userProfileView UserProfileView
  userProfileView.UserID = dbUserProfileView.UserID
  userProfileView.UserName = dbUserProfileView.UserName
  userProfileView.EmailAddress = dbUserProfileView.EmailAddress
  userProfileView.Created = dbUserProfileView.Created
  userProfileView.DefaultCharacterID = FromNullInt64(dbUserProfileView.DefaultCharacterID)
  userProfileView.DefaultCharacterName = FromNullString(dbUserProfileView.DefaultCharacterName)

  response := &Response{
    Success:  true,
    Error:    nil,
    Data:     UserGetResponseData{
      User:  userProfileView,
    },
  }

  res.Header().Set("Content-Type", "application/json")
  json.NewEncoder(res).Encode(response)
}


func (r *UserRouter) handleUpdate(res http.ResponseWriter, req *http.Request) {
  decoder := json.NewDecoder(req.Body)
  var updateRequestData UserUpdateRequestData

  err := decoder.Decode(&updateRequestData)
  if err != nil {
    http.Error(res, fmt.Sprintf("Invalid JSON request: %s", err.Error()), http.StatusBadRequest)
    return
  }

  var userProfileUpdate db.UserProfileUpdate
  userProfileUpdate.UserID = updateRequestData.UserID
  userProfileUpdate.UserName = updateRequestData.UserName
  // Convert to sql.NullInt64 because a user can have no default character info
  userProfileUpdate.DefaultCharacterID = ToNullInt64(updateRequestData.DefaultCharacterID)
  userProfileUpdate.DefaultCharacterGsp = ToNullInt64(updateRequestData.DefaultCharacterGsp)


  userID, err := r.SysUtils.Database.UpdateUserProfile(userProfileUpdate)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error updating user in database: %s", err.Error()), http.StatusInternalServerError)
    return
  }

  response := &Response{
    Success:  true,
    Error:    nil,
    Data:     UserUpdateResponseData{
      UserID:  userID,
    },
  }

  res.Header().Set("Content-Type", "application/json")
  json.NewEncoder(res).Encode(response)
}
