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
          Response Data
----------------------------------*/

// UserGetResponseData is the data we send back
// after a successfully getting all user's info
type UserGetResponseData struct {
  User  *UserProfileView  `json:"user"`
}


// UserUpdateResponseData is the data we send
// back after a successfully creating a new user
type UserUpdateResponseData struct {
  User  *UserProfileView  `json:"user"`
}


/*---------------------------------
          API <--> SQL
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


// ToAPIUserProfileView maps from a db.UserProfileView
// (which has things like sql.NullString) into
// an api.UserProfileView, which is JSON representable
func ToAPIUserProfileView(dbUserProfileView *db.UserProfileView) *UserProfileView {
  userProfileView := new(UserProfileView)
  userProfileView.UserID = dbUserProfileView.UserID
  userProfileView.UserName = dbUserProfileView.UserName
  userProfileView.EmailAddress = dbUserProfileView.EmailAddress
  userProfileView.Created = dbUserProfileView.Created
  userProfileView.DefaultCharacterGsp = FromNullInt64(dbUserProfileView.DefaultCharacterGsp)
  userProfileView.DefaultCharacterID = FromNullInt64(dbUserProfileView.DefaultCharacterID)
  userProfileView.DefaultCharacterName = FromNullString(dbUserProfileView.DefaultCharacterName)

  return userProfileView
}


// ToDBUserUpdate maps from an api.UserUpdateRequestDat 
// into a sb.UserProfileUpdate, which has things like sql.NullInt64
func ToDBUserUpdate(userUpdateRequestData *UserUpdateRequestData) *db.UserProfileUpdate {
  dbUserUpdate := new(db.UserProfileUpdate)
  dbUserUpdate.UserID = userUpdateRequestData.UserID
  dbUserUpdate.UserName = userUpdateRequestData.UserName
  dbUserUpdate.DefaultCharacterID = ToNullInt64(userUpdateRequestData.DefaultCharacterID)
  dbUserUpdate.DefaultCharacterGsp = ToNullInt64(userUpdateRequestData.DefaultCharacterGsp)

  return dbUserUpdate
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
  userProfileView := ToAPIUserProfileView(dbUserProfileView)

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
  updateRequestData := new(UserUpdateRequestData)

  err := decoder.Decode(updateRequestData)
  if err != nil {
    http.Error(res, fmt.Sprintf("Invalid JSON request: %s", err.Error()), http.StatusBadRequest)
    return
  }

  dbUserProfileUpdate := ToDBUserUpdate(updateRequestData)
  userID, err := r.SysUtils.Database.UpdateUserProfile(dbUserProfileUpdate)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error updating user in database: %s", err.Error()), http.StatusInternalServerError)
    return
  }
  dbUserProfileView, err := r.SysUtils.Database.GetUserProfileViewByID(userID)
  userProfileView := ToAPIUserProfileView(dbUserProfileView)
  

  response := &Response{
    Success:  true,
    Error:    nil,
    Data:     UserUpdateResponseData{
      User:  userProfileView,
    },
  }

  res.Header().Set("Content-Type", "application/json")
  json.NewEncoder(res).Encode(response)
}
