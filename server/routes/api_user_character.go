package routes

import (
  "encoding/json"
  "fmt"
  "net/http"

  "github.com/cakebin/smush/server/services/db"
)


/*---------------------------------
          Request Data
----------------------------------*/

// UserCharacterCreateRequestData describes the data we're 
// expecting when a user attempts to create a "saved character"
type UserCharacterCreateRequestData struct {
  UserID           int     `json:"userId"`
  CharacterID      int     `json:"characterId"`
  CharacterGsp     *int64  `json:"characterGsp"`
  AltCostume       *int64  `json:"altCostume"`
}


// UserCharacterUpdateRequestData describes the data we're 
// expecting when a user attempts to update one of their "saved characters"
type UserCharacterUpdateRequestData struct {
  UserCharacterID  int     `json:"userCharacterId"`
  UserID           int     `json:"userId"`
  CharacterID      *int64  `json:"characterId"`
  CharacterGsp     *int64  `json:"characterGsp"`
  AltCostume       *int64  `json:"altCostume"`
}


// UserCharacterDeleteRequestData describes the data we're 
// expecting when a user attempts to delete one of their "saved characters"
type UserCharacterDeleteRequestData struct {
  UserID           int  `json:"userId"`
  UserCharacterID  int  `json:"userCharacterId"`
}


/*---------------------------------
          Response Data
----------------------------------*/

// UserCharacterCreateResponseData is the data we send back
// after a successfully creating a new "saved character" for a given user
type UserCharacterCreateResponseData struct {
  UserCharacters  []*UserCharacterView  `json:"userCharacters"`
  User            *UserProfileView      `json:"user"`
}


// UserCharacterUpdateResponseData is the data we send back
// after a successfully creating a new "saved character" for a given user
type UserCharacterUpdateResponseData struct {
  UserCharacters  []*UserCharacterView  `json:"userCharacters"`
  User            *UserProfileView      `json:"user"`
}


// UserCharacterDeleteResponseData is the data we send back
// after a successfully deleting a "saved character" for a given user
type UserCharacterDeleteResponseData struct {
  UserCharacters  []*UserCharacterView  `json:"userCharacters"`
  User            *UserProfileView      `json:"user"`
}


/*--------------------------------
          API <--> SQL
----------------------------------*/

// UserCharacterView is a translation from the SQL result
// which can have things like `sql.NullInt64`, so we
// need to translate that to regular JSON objects
type UserCharacterView struct {
  UserCharacterID  int     `json:"userCharacterId"`
  UserID           int     `json:"userId"`
  CharacterID      int     `json:"characterId"`
  CharacterName    string  `json:"characterName"`
  CharacterGsp     *int64  `json:"characterGsp"`
  AltCostume       *int64  `json:"altCostume"`
}


// ToAPIUserCharacterView maps from a db.UserCharacterView
// (which has things like sql.NullInt64) into
// an api.UserCharacter, which is JSON representable
func ToAPIUserCharacterView(dbUserCharView *db.UserCharacterView) *UserCharacterView {
  userCharView := new(UserCharacterView)
  userCharView.UserCharacterID = dbUserCharView.UserCharacterID
  userCharView.UserID = dbUserCharView.UserID
  userCharView.CharacterID = dbUserCharView.CharacterID
  userCharView.CharacterName = dbUserCharView.CharacterName
  userCharView.CharacterGsp = FromNullInt64(dbUserCharView.CharacterGsp)
  userCharView.AltCostume = FromNullInt64(dbUserCharView.AltCostume)

  return userCharView
}


// ToAPIUserCharacterViews maps from a []*db.UserCharacterView to a []*UserCharacterView
func ToAPIUserCharacterViews(dbUserCharViews []*db.UserCharacterView) []*UserCharacterView {
  userCharViews := make([]*UserCharacterView, 0)
  for _, dbUserCharView := range dbUserCharViews {
    userCharView := ToAPIUserCharacterView(dbUserCharView)
    userCharViews = append(userCharViews, userCharView)
  }

  return userCharViews
}


// ToBDUserCharacterCreate maps from an api.UserCharacterCreateRequestData
// to a db.UserCharacterCreate, which has fields like sql.NullInt64
func ToBDUserCharacterCreate(userCharCreateRequestData *UserCharacterCreateRequestData) *db.UserCharacterCreate {
  dbUserCharCreate := new(db.UserCharacterCreate)
  dbUserCharCreate.UserID = userCharCreateRequestData.UserID
  dbUserCharCreate.CharacterID = userCharCreateRequestData.CharacterID
  dbUserCharCreate.CharacterGsp = ToNullInt64(userCharCreateRequestData.CharacterGsp)
  dbUserCharCreate.AltCostume = ToNullInt64(userCharCreateRequestData.AltCostume)

  return dbUserCharCreate
}


// ToDBUserCharacterUpdate maps from an api.UserCharacterUpdateRequestData
// to a db.UserCharacterUpdate, which has fields like sql.NullInt64
func ToDBUserCharacterUpdate(userCharUpdateRequestData *UserCharacterUpdateRequestData) *db.UserCharacterUpdate {
  dbUserCharUpdate := new(db.UserCharacterUpdate)
  dbUserCharUpdate.UserCharacterID = userCharUpdateRequestData.UserCharacterID
  dbUserCharUpdate.UserID = userCharUpdateRequestData.UserID
  dbUserCharUpdate.CharacterID = ToNullInt64(userCharUpdateRequestData.CharacterID)
  dbUserCharUpdate.CharacterGsp = ToNullInt64(userCharUpdateRequestData.CharacterGsp)
  dbUserCharUpdate.AltCostume = ToNullInt64(userCharUpdateRequestData.AltCostume)

  return dbUserCharUpdate
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
  createRequestData := new(UserCharacterCreateRequestData)

  err := decoder.Decode(createRequestData)
  if err != nil {
    http.Error(res, fmt.Sprintf("Invalid JSON request: %s", err.Error()), http.StatusBadRequest)
    return
  }

  dbUserCharCreate := ToBDUserCharacterCreate(createRequestData)
  _, err = r.Services.Database.CreateUserCharacter(dbUserCharCreate)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error creating new user character in database: %s", err.Error()), http.StatusInternalServerError)
    return
  }

  dbUserCharViews, err := r.Services.Database.GetUserCharacterViewsByUserID(createRequestData.UserID)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error fetching user character views in database after creating new user_character: %s", err.Error()), http.StatusInternalServerError)
    return
  }
  userCharViews := ToAPIUserCharacterViews(dbUserCharViews)


  dbUserProfileView, err := r.Services.Database.GetUserProfileViewByUserID(createRequestData.UserID)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error fetching user view in database after creating new user_character: %s", err.Error()), http.StatusInternalServerError)
    return
  }
  userProfileView := ToAPIUserProfileView(dbUserProfileView) 

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
  updateRequestData := new(UserCharacterUpdateRequestData)

  err := decoder.Decode(updateRequestData)
  if err != nil {
    http.Error(res, fmt.Sprintf("Invalid JSON request: %s", err.Error()), http.StatusBadRequest)
    return
  }

  dbUserCharUpdate := ToDBUserCharacterUpdate(updateRequestData)
  _, err = r.Services.Database.UpdateUserCharacter(dbUserCharUpdate)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error updating user character in database: %s", err.Error()), http.StatusInternalServerError)
    return
  }

  dbUserCharViews, err := r.Services.Database.GetUserCharacterViewsByUserID(updateRequestData.UserID)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error fetching user character views in database after updating user_character: %s", err.Error()), http.StatusInternalServerError)
    return
  }
  userCharViews := ToAPIUserCharacterViews(dbUserCharViews)

  dbUserProfileView, err := r.Services.Database.GetUserProfileViewByUserID(updateRequestData.UserID)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error fetching user view in database after updating user_character: %s", err.Error()), http.StatusInternalServerError)
    return
  }
  userProfileView := ToAPIUserProfileView(dbUserProfileView)

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
  deleteRequestData := new(UserCharacterDeleteRequestData)

  err := decoder.Decode(deleteRequestData)
  if err != nil {
    http.Error(res, fmt.Sprintf("Invalid JSON request: %s", err.Error()), http.StatusBadRequest)
    return
  }
  
  dbUserProfileView, err := r.Services.Database.GetUserProfileViewByUserID(deleteRequestData.UserID)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error fetching user view in database before deleting user_character: %s", err.Error()), http.StatusInternalServerError)
    return
  }
  // If we have a current default user character id, we may have to remove it before deleting the child row
  userProfileView := ToAPIUserProfileView(dbUserProfileView)
  // Only remove the current default character if it's the userCharacter we are deleting
  if (*userProfileView.DefaultUserCharacterID == int64(deleteRequestData.UserCharacterID)) {
    updateRequestData := new(db.UserDefaultUserCharacterUpdate)
    updateRequestData.UserID = deleteRequestData.UserID
    updateRequestData.UserCharacterID = ToNullInt64(nil)
    _, err = r.Services.Database.UpdateUserDefaultUserCharacter(updateRequestData)
    if err != nil {
      http.Error(res, fmt.Sprintf("Error updating user default character in database: %s", err.Error()), http.StatusInternalServerError)
      return
    }
  }

  _, err = r.Services.Database.DeleteUserCharacterByID(deleteRequestData.UserCharacterID)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error deleting user character in database: %s", err.Error()), http.StatusInternalServerError)
    return
  }

  dbUserCharViews, err := r.Services.Database.GetUserCharacterViewsByUserID(deleteRequestData.UserID)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error fetching user character views in database after deleting user_character: %s", err.Error()), http.StatusInternalServerError)
    return
  }
  userCharViews := ToAPIUserCharacterViews(dbUserCharViews)

  dbUpdatedUserProfileView, err := r.Services.Database.GetUserProfileViewByUserID(deleteRequestData.UserID)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error fetching user view in database after deleting user_character: %s", err.Error()), http.StatusInternalServerError)
    return
  }
  updatedUserProfileView := ToAPIUserProfileView(dbUpdatedUserProfileView)

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
