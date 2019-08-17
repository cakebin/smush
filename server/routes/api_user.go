package routes

import (
  "encoding/json"
  "fmt"
  "net/http"
  "strconv"
  "time"

  "github.com/cakebin/smush/server/services/db"
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
}


// UserUpdateDefaultUserCharacterRequestData describes the data we're
// expecting when a user attempts to update their default user character
type UserUpdateDefaultUserCharacterRequestData struct {
  UserID           int     `json:"userId"`
  UserCharacterID  *int64  `json:"userCharacterId"`
}


/*---------------------------------
          Response Data
----------------------------------*/

// UserGetResponseData is the data we send back
// after a successfully getting all user's info
type UserGetResponseData struct {
  User            *UserProfileView      `json:"user"`
  UserCharacters  []*UserCharacterView  `json:"userCharacters"`
}


// UserUpdateResponseData is the data we send
// back after a successfully creating a new user
type UserUpdateResponseData struct {
  User  *UserProfileView  `json:"user"`
}


// UserUpdateDefaultUserCharacterResponseData is the data we send back
// after successfully updating a user's default user character
type UserUpdateDefaultUserCharacterResponseData struct {
  User            *UserProfileView      `json:"user"`
  UserCharacters  []*UserCharacterView  `json:"userCharacters"`
}


/*---------------------------------
          API <--> SQL
----------------------------------*/

// UserProfileView is a translation from the SQL result
// which can have things like `sql.NullInt64`, so we 
// need to translate that to regular JSON objects
type UserProfileView struct {
  UserID                        int        `json:"userId"`
  UserName                      string     `json:"userName"`
  EmailAddress                  string     `json:"emailAddress"`
  Created                       time.Time  `json:"created"`
  DefaultUserCharacterID        *int64     `json:"defaultUserCharacterId,omitempty"`
  DefaultUserCharacterGsp       *int64     `json:"defaultUserCharacterGsp,omitempty"`
  DefaultCharacterID            *int64     `json:"defaultCharacterId,omitempty"`
  DefaultCharacterName          *string    `json:"defaultCharacterName,omitempty"`
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
  userProfileView.DefaultUserCharacterID = FromNullInt64(dbUserProfileView.DefaultUserCharacterID)
  userProfileView.DefaultUserCharacterGsp = FromNullInt64(dbUserProfileView.DefaultUserCharacterGsp)
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

  return dbUserUpdate
}


// ToDBUserDefaultUserCharacterUpdate maps from an api.UserUpdateDefaultUserCharacterRequestData
// to a db.UserDefaultUserCharacterUpdate for updating said user's default_user_character_id
func ToDBUserDefaultUserCharacterUpdate(requestData *UserUpdateDefaultUserCharacterRequestData) *db.UserDefaultUserCharacterUpdate {
  dbUserUpdateDefaultUserChar := new(db.UserDefaultUserCharacterUpdate)
  dbUserUpdateDefaultUserChar.UserID = requestData.UserID
  dbUserUpdateDefaultUserChar.UserCharacterID = ToNullInt64(requestData.UserCharacterID)

  return dbUserUpdateDefaultUserChar
}

/*---------------------------------
             Router
----------------------------------*/

// UserRouter is responsible for serving "/api/user"
// Basically, connecting to our Postgres DB for all
// of the CRUD operations for our "User" models
type UserRouter struct {
  Services             *Services
  UserCharacterRouter  *UserCharacterRouter
}


func (r *UserRouter) ServeHTTP(res http.ResponseWriter, req *http.Request) {
  var head string
  head, req.URL.Path = ShiftPath(req.URL.Path)

  // Delegate to sub routers first
  switch head {
  case "character":
    r.UserCharacterRouter.ServeHTTP(res, req)

  // Otherwise, handle the user specific requests
  default:
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
        case "update_profile":
          r.handleUpdateProfile(res, req)
        case "update_default_user_character":
          r.handleUpdateDefaultUserCharacter(res, req)
        default:
          http.Error(res, fmt.Sprintf("Unsupport POST path %s", head), http.StatusBadRequest)
          return
        }

      // Unsupported Method Response
      default:
        http.Error(res, fmt.Sprintf("Unsupported Method type %s", req.Method), http.StatusBadRequest)
      }
  }
}


// NewUserRouter makes a new api/user router and hooks up its services
func NewUserRouter(routerServices *Services) *UserRouter {
  router := new(UserRouter)

  router.Services = routerServices
  router.UserCharacterRouter = NewUserCharacterRouter(routerServices)

  return router
}


/*---------------------------------
             Handlers
----------------------------------*/

func (r *UserRouter) handleGetByID(res http.ResponseWriter, req *http.Request) {
  var head string
  head, req.URL.Path = ShiftPath(req.URL.Path)

  userID, err := strconv.Atoi(head)
  if err != nil {
    http.Error(res, fmt.Sprintf("Invalid user id: %s", head), http.StatusBadRequest)
    return
  }

  // Get the basic user profile information
  dbUserProfileView, err := r.Services.Database.GetUserProfileViewByUserID(userID)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error getting user with userID %d: %s", userID, err.Error()), http.StatusInternalServerError)
    return
  }
  userProfileView := ToAPIUserProfileView(dbUserProfileView)

  // Also get the user's saved characters
  dbUserCharViews, err := r.Services.Database.GetUserCharacterViewsByUserID(userID)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error getting user's saved characters with userID %d: %s", userID, err.Error()), http.StatusInternalServerError)
    return
  }
  userCharViews := make([]*UserCharacterView, 0)
  for _, dbUserCharView := range dbUserCharViews {
    userCharView := ToAPIUserCharacterView(dbUserCharView)
    userCharViews = append(userCharViews, userCharView)
  }

  response := &Response{
    Success:  true,
    Error:    nil,
    Data:     UserGetResponseData{
      User:            userProfileView,
      UserCharacters:  userCharViews,
    },
  }

  res.Header().Set("Content-Type", "application/json")
  json.NewEncoder(res).Encode(response)
}


func (r *UserRouter) handleUpdateProfile(res http.ResponseWriter, req *http.Request) {
  decoder := json.NewDecoder(req.Body)
  updateRequestData := new(UserUpdateRequestData)

  err := decoder.Decode(updateRequestData)
  if err != nil {
    http.Error(res, fmt.Sprintf("Invalid JSON request: %s", err.Error()), http.StatusBadRequest)
    return
  }

  dbUserProfileUpdate := ToDBUserUpdate(updateRequestData)
  userID, err := r.Services.Database.UpdateUserProfile(dbUserProfileUpdate)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error updating user in database: %s", err.Error()), http.StatusInternalServerError)
    return
  }
  dbUserProfileView, err := r.Services.Database.GetUserProfileViewByUserID(userID)
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


func (r *UserRouter) handleUpdateDefaultUserCharacter(res http.ResponseWriter, req *http.Request) {
  decoder := json.NewDecoder(req.Body)
  updateRequestData := new(UserUpdateDefaultUserCharacterRequestData)

  err := decoder.Decode(updateRequestData)
  if err != nil {
    http.Error(res, fmt.Sprintf("Invalid JSON request: %s", err.Error()), http.StatusBadRequest)
    return
  }

  dbUserDefaultUserCharacterUpdate := ToDBUserDefaultUserCharacterUpdate(updateRequestData)
  userID, err := r.Services.Database.UpdateUserDefaultUserCharacter(dbUserDefaultUserCharacterUpdate)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error updating user default character in database: %s", err.Error()), http.StatusInternalServerError)
    return
  }

  dbUserProfileView, err := r.Services.Database.GetUserProfileViewByUserID(userID)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error getting user in database after updating default user character: %s", err.Error()), http.StatusInternalServerError)
    return
  }
  userProfileView := ToAPIUserProfileView(dbUserProfileView)

  dbUserCharViews, err := r.Services.Database.GetUserCharacterViewsByUserID(updateRequestData.UserID)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error fetching user character views in database after updating user_character: %s", err.Error()), http.StatusInternalServerError)
    return
  }
  userCharViews := ToAPIUserCharacterViews(dbUserCharViews)

  response := &Response{
    Success:  true,
    Error:    nil,
    Data:     UserUpdateDefaultUserCharacterResponseData{
      User:            userProfileView,
      UserCharacters:  userCharViews,
    },
  }

  res.Header().Set("Content-Type", "application/json")
  json.NewEncoder(res).Encode(response)
}
