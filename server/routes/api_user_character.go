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
}


// UserCharacterUpdateRequestData describes the data we're 
// expecting when a user attempts to update one of their "saved characters"
type UserCharacterUpdateRequestData struct {
  UserCharacterID  int     `json:"userCharacterId"`
  UserID           *int64  `json:"userId"`
  CharacterID      *int64  `json:"characterId"`
  CharacterGsp     *int64  `json:"characterGsp"`
}


// UserCharacterDeleteRequestData describes the data we're 
// expecting when a user attempts to delete one of their "saved characters"
type UserCharacterDeleteRequestData struct {
  UserCharacterID  int  `json:"userCharacterId"`
}


/*---------------------------------
          Response Data
----------------------------------*/

// UserCharacterCreateResponseData is the data we send back
// after a successfully creating a new "saved character" for a given user
type UserCharacterCreateResponseData struct {
  UserCharacter *UserCharacterView  `json:"userCharacter"`
}


// UserCharacterUpdateResponseData is the data we send back
// after a successfully creating a new "saved character" for a given user
type UserCharacterUpdateResponseData struct {
  UserCharacter *UserCharacterView  `json:"userCharacter"`
}


// UserCharacterDeleteResponseData is the data we send back
// after a successfully deleting a "saved character" for a given user
type UserCharacterDeleteResponseData struct {
  UserCharacterID  int  `json:"userCharacterId"`
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

  return userCharView
}


// ToBDUserCharacterCreate maps from an api.UserCharacterCreateRequestData
// to a db.UserCharacterCreate, which has fields like sql.NullInt64
func ToBDUserCharacterCreate(userCharCreateRequestData *UserCharacterCreateRequestData) *db.UserCharacterCreate {
  dbUserCharCreate := new(db.UserCharacterCreate)
  dbUserCharCreate.UserID = userCharCreateRequestData.UserID
  dbUserCharCreate.CharacterID = userCharCreateRequestData.CharacterID
  dbUserCharCreate.CharacterGsp = ToNullInt64(userCharCreateRequestData.CharacterGsp)

  return dbUserCharCreate
}


// ToDBUserCharacterUpdate maps from an api.UserCharacterUpdateRequestData
// to a db.UserCharacterUpdate, which has fields like sql.NullInt64
func ToDBUserCharacterUpdate(userCharUpdateRequestData *UserCharacterUpdateRequestData) *db.UserCharacterUpdate {
  dbUserCharUpdate := new(db.UserCharacterUpdate)
  dbUserCharUpdate.UserCharacterID = userCharUpdateRequestData.UserCharacterID
  dbUserCharUpdate.UserID = ToNullInt64(userCharUpdateRequestData.UserID)
  dbUserCharUpdate.CharacterID = ToNullInt64(userCharUpdateRequestData.CharacterID)
  dbUserCharUpdate.CharacterGsp = ToNullInt64(userCharUpdateRequestData.CharacterGsp)

  return dbUserCharUpdate
}


/*---------------------------------
             Router
----------------------------------*/

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
  userCharID, err := r.Services.Database.CreateUserCharacter(dbUserCharCreate)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error creating new user character in database: %s", err.Error()), http.StatusInternalServerError)
    return
  }

  dbUserCharView, err := r.Services.Database.GetUserCharacterViewByUserCharacterID(userCharID)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error fetching new user character view in database: %s", err.Error()), http.StatusInternalServerError)
    return
  }
  userCharView := ToAPIUserCharacterView(dbUserCharView)

  response := Response{
    Success:  true,
    Error:    nil,
    Data:    UserCharacterCreateResponseData{
      UserCharacter:  userCharView,
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
  userCharID, err := r.Services.Database.UpdateUserCharacter(dbUserCharUpdate)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error updating user character in database: %s", err.Error()), http.StatusInternalServerError)
    return
  }

  dbUserCharView, err := r.Services.Database.GetUserCharacterViewByUserCharacterID(userCharID)
  userCharView := ToAPIUserCharacterView(dbUserCharView)

  response := &Response{
    Success:  true,
    Error:    nil,
    Data:     UserCharacterUpdateResponseData{
      UserCharacter:  userCharView,
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

  deletedUserCharID, err := r.Services.Database.DeleteUserCharacterByID(deleteRequestData.UserCharacterID)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error deleting user character in database: %s", err.Error()), http.StatusInternalServerError)
    return
  }

  response := &Response{
    Success:  true,
    Error:    nil,
    Data:     UserCharacterDeleteResponseData{
      UserCharacterID:  deletedUserCharID,
    },
  }

  res.Header().Set("Content-Type", "application/json")
  json.NewEncoder(res).Encode(response)
}
