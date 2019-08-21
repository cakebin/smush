package routes

import (
  "encoding/json"
  "fmt"
  "net/http"
  "strconv"

  "github.com/cakebin/smush/server/services/db"
)


/*---------------------------------
          Response Data
----------------------------------*/

type UserGetAllResponseData struct {
  Users  []*db.User  `json:"users"`
}

// UserGetResponseData is the data we send back
// after a successfully getting all user's info
type UserGetResponseData struct {
  User            *db.UserProfileView      `json:"user"`
  UserCharacters  []*db.UserCharacterView  `json:"userCharacters"`
}


// UserUpdateResponseData is the data we send
// back after a successfully creating a new user
type UserUpdateResponseData struct {
  User            *db.UserProfileView      `json:"user"`
  UserCharacters  []*db.UserCharacterView  `json:"userCharacters"`
}


// UserUpdateDefaultUserCharacterResponseData is the data we send back
// after successfully updating a user's default user character
type UserUpdateDefaultUserCharacterResponseData struct {
  User            *db.UserProfileView      `json:"user"`
  UserCharacters  []*db.UserCharacterView  `json:"userCharacters"`
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
        case "getall":
          r.handleGetAll(res, req)
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

  userID, err := strconv.ParseInt(head, 10, 64)
  if err != nil {
    http.Error(res, fmt.Sprintf("Invalid user id: %s", head), http.StatusBadRequest)
    return
  }

  // Get the basic user profile information
  userProfileView, err := r.Services.Database.GetUserProfileViewByUserID(userID)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error getting user with userID %d: %s", userID, err.Error()), http.StatusInternalServerError)
    return
  }

  // Also get the user's saved characters
  userCharViews, err := r.Services.Database.GetUserCharacterViewsByUserID(userID)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error getting user's saved characters with userID %d: %s", userID, err.Error()), http.StatusInternalServerError)
    return
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


func (r *UserRouter) handleGetAll(res http.ResponseWriter, req *http.Request) {
  users, err := r.Services.Database.GetAllUsers()
  if err != nil {
    http.Error(res, fmt.Sprintf("Error getting all users: %s", err.Error()), http.StatusInternalServerError)
    return
  }

  response := &Response{
    Success:  true,
    Error:    nil,
    Data:     UserGetAllResponseData{
      Users: users,
    },
  }

  res.Header().Set("Content-Type", "application/json")
  json.NewEncoder(res).Encode(response)
}


func (r *UserRouter) handleUpdateProfile(res http.ResponseWriter, req *http.Request) {
  decoder := json.NewDecoder(req.Body)
  userProfileUpdate := new(db.UserProfileUpdate)

  err := decoder.Decode(userProfileUpdate)
  if err != nil {
    http.Error(res, fmt.Sprintf("Invalid JSON request: %s", err.Error()), http.StatusBadRequest)
    return
  }

  userID, err := r.Services.Database.UpdateUserProfile(userProfileUpdate)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error updating user in database: %s", err.Error()), http.StatusInternalServerError)
    return
  }
  userProfileView, err := r.Services.Database.GetUserProfileViewByUserID(userID)

  userCharViews, err := r.Services.Database.GetUserCharacterViewsByUserID(userID)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error fetching user character views in database after updating user_character: %s", err.Error()), http.StatusInternalServerError)
    return
  }

  response := &Response{
    Success:  true,
    Error:    nil,
    Data:     UserUpdateResponseData{
      User:           userProfileView,
      UserCharacters: userCharViews,
    },
  }

  res.Header().Set("Content-Type", "application/json")
  json.NewEncoder(res).Encode(response)
}


func (r *UserRouter) handleUpdateDefaultUserCharacter(res http.ResponseWriter, req *http.Request) {
  decoder := json.NewDecoder(req.Body)
  userDefaultUserCharUpdate := new(db.UserDefaultUserCharacterUpdate)

  err := decoder.Decode(userDefaultUserCharUpdate)
  if err != nil {
    http.Error(res, fmt.Sprintf("Invalid JSON request: %s", err.Error()), http.StatusBadRequest)
    return
  }

  userID, err := r.Services.Database.UpdateUserDefaultUserCharacter(userDefaultUserCharUpdate)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error updating user default character in database: %s", err.Error()), http.StatusInternalServerError)
    return
  }

  userProfileView, err := r.Services.Database.GetUserProfileViewByUserID(userID)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error getting user in database after updating default user character: %s", err.Error()), http.StatusInternalServerError)
    return
  }

  userCharViews, err := r.Services.Database.GetUserCharacterViewsByUserID(userID)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error fetching user character views in database after updating user_character: %s", err.Error()), http.StatusInternalServerError)
    return
  }

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
