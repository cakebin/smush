package routes

import (
  "encoding/json"
  "fmt"
  "net/http"

  "github.com/cakebin/smush/server/services/db"
)


/*---------------------------------
          Response Data
----------------------------------*/

// UserCharacterCreateResponseData is the data we send back
// after a successfully creating a new "saved character" for a given user
type UserCharacterCreateResponseData struct {
  UserCharacters  []*db.UserCharacterView  `json:"userCharacters"`
  User            *db.UserProfileView      `json:"user"`
}


// UserCharacterUpdateResponseData is the data we send back
// after a successfully creating a new "saved character" for a given user
type UserCharacterUpdateResponseData struct {
  UserCharacters  []*db.UserCharacterView  `json:"userCharacters"`
  User            *db.UserProfileView      `json:"user"`
}


// UserCharacterDeleteResponseData is the data we send back
// after a successfully deleting a "saved character" for a given user
type UserCharacterDeleteResponseData struct {
  UserCharacters  []*db.UserCharacterView  `json:"userCharacters"`
  User            *db.UserProfileView      `json:"user"`
}


/*---------------------------------
             Router
----------------------------------*/

// UserCharacterRouter handles all of /api/user/character
type UserCharacterRouter struct {
  Services  *Services
}


func (r *UserCharacterRouter) ServeHTTP(res http.ResponseWriter, req *http.Request) {
  var head string
  head, req.URL.Path = ShiftPath(req.URL.Path)

  switch req.Method {
  // POST Request Handlers
  case http.MethodPost:
    switch head {
    case "create":
      r.handleCreate(res, req)
    case "update":
      r.handleUpdate(res, req)
    case "delete":
      r.handleDelete(res, req)
    default:
      http.Error(res, fmt.Sprintf("Unsupported POST path %s", head), http.StatusBadRequest)
      return
    }
  // Unsupported Method Response
  default:
    http.Error(res, fmt.Sprintf("Unsupported Method type %s", req.Method), http.StatusBadRequest)
  }
}


// NewUserCharacterRouter makes a new api/user/character router and hooks up its services
func NewUserCharacterRouter(routerServices *Services) *UserCharacterRouter {
  router := new(UserCharacterRouter)

  router.Services = routerServices

  return router
}


/*---------------------------------
             Handlers
----------------------------------*/

func (r *UserCharacterRouter) handleCreate(res http.ResponseWriter, req *http.Request) {
  decoder := json.NewDecoder(req.Body)
  userCharCreate := new(db.UserCharacterCreate)

  err := decoder.Decode(userCharCreate)
  if err != nil {
    http.Error(res, fmt.Sprintf("Invalid JSON request: %s", err.Error()), http.StatusBadRequest)
    return
  }

  _, err = r.Services.Database.CreateUserCharacter(userCharCreate)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error creating new user character in database: %s", err.Error()), http.StatusInternalServerError)
    return
  }

  userCharViews, err := r.Services.Database.GetUserCharacterViewsByUserID(userCharCreate.UserID)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error fetching user character views in database after creating new user_character: %s", err.Error()), http.StatusInternalServerError)
    return
  }

  userProfileView, err := r.Services.Database.GetUserProfileViewByUserID(userCharCreate.UserID)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error fetching user view in database after creating new user_character: %s", err.Error()), http.StatusInternalServerError)
    return
  }

  response := Response{
    Success:  true,
    Error:    nil,
    Data:    UserCharacterCreateResponseData{
      UserCharacters:  userCharViews,
      User:            userProfileView,
    },
  }

  res.Header().Set("Content-Type", "application/json")
  json.NewEncoder(res).Encode(response)
}


func (r *UserCharacterRouter) handleUpdate(res http.ResponseWriter, req *http.Request) {
  decoder := json.NewDecoder(req.Body)
  userCharUpdate := new(db.UserCharacterUpdate)

  err := decoder.Decode(userCharUpdate)
  if err != nil {
    http.Error(res, fmt.Sprintf("Invalid JSON request: %s", err.Error()), http.StatusBadRequest)
    return
  }

  _, err = r.Services.Database.UpdateUserCharacter(userCharUpdate)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error updating user character in database: %s", err.Error()), http.StatusInternalServerError)
    return
  }

  userCharViews, err := r.Services.Database.GetUserCharacterViewsByUserID(userCharUpdate.UserID)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error fetching user character views in database after updating user_character: %s", err.Error()), http.StatusInternalServerError)
    return
  }

  userProfileView, err := r.Services.Database.GetUserProfileViewByUserID(userCharUpdate.UserID)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error fetching user view in database after updating user_character: %s", err.Error()), http.StatusInternalServerError)
    return
  }

  response := &Response{
    Success:  true,
    Error:    nil,
    Data:     UserCharacterUpdateResponseData{
      UserCharacters:  userCharViews,
      User:            userProfileView,
    },
  }

  res.Header().Set("Content-Type", "application/json")
  json.NewEncoder(res).Encode(response)
}


func (r *UserCharacterRouter) handleDelete(res http.ResponseWriter, req *http.Request) {
  decoder := json.NewDecoder(req.Body)
  userCharDelete := new(db.UserCharacterDelete)

  err := decoder.Decode(userCharDelete)
  if err != nil {
    http.Error(res, fmt.Sprintf("Invalid JSON request: %s", err.Error()), http.StatusBadRequest)
    return
  }
  
  userProfileView, err := r.Services.Database.GetUserProfileViewByUserID(userCharDelete.UserID)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error fetching user view in database before deleting user_character: %s", err.Error()), http.StatusInternalServerError)
    return
  }
  // If we have a current default user character id, we may have to remove it before deleting the child row
  // Only remove the current default character if it's the userCharacter we are deleting
  if (userProfileView.DefaultUserCharacterID.Int64 == userCharDelete.UserCharacterID.Int64) {
    userDefaultUserCharUpdate := new(db.UserDefaultUserCharacterUpdate)
    userDefaultUserCharUpdate.UserID = userCharDelete.UserID
    userDefaultUserCharUpdate.UserCharacterID = userCharDelete.UserCharacterID

    _, err = r.Services.Database.UpdateUserDefaultUserCharacter(userDefaultUserCharUpdate)
    if err != nil {
      http.Error(res, fmt.Sprintf("Error updating user default character in database: %s", err.Error()), http.StatusInternalServerError)
      return
    }
  }

  _, err = r.Services.Database.DeleteUserCharacterByID(userCharDelete.UserCharacterID.Int64)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error deleting user character in database: %s", err.Error()), http.StatusInternalServerError)
    return
  }

  userCharViews, err := r.Services.Database.GetUserCharacterViewsByUserID(userCharDelete.UserID)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error fetching user character views in database after deleting user_character: %s", err.Error()), http.StatusInternalServerError)
    return
  }

  updatedUserProfileView, err := r.Services.Database.GetUserProfileViewByUserID(userCharDelete.UserID)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error fetching user view in database after deleting user_character: %s", err.Error()), http.StatusInternalServerError)
    return
  }

  response := &Response{
    Success:  true,
    Error:    nil,
    Data:     UserCharacterDeleteResponseData{
      UserCharacters:  userCharViews,
      User:            updatedUserProfileView,
    },
  }

  res.Header().Set("Content-Type", "application/json")
  json.NewEncoder(res).Encode(response)
}
