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

// CharacterGetAllResponseData is the data we send back
// after a successfully getting all character data
type CharacterGetAllResponseData struct {
  Characters  []*db.Character  `json:"characters"`
}


// CharacterCreateResponseData is the data we send
//  back after a successfully creating a new character
type CharacterCreateResponseData struct {
  Character  *db.Character  `json:"character"`
}


// CharacterUpdateResponseData is the data we send
//  back after a successfully udpating a new character
type CharacterUpdateResponseData struct {
  Character  *db.Character  `json:"character"`
}


/*---------------------------------
             Router
----------------------------------*/

// CharacterRouter is responsible for serving "/api/character"
// Basically, connecting to our Postgres DB for all
// of the CRUD operations for our "Character" models
type CharacterRouter struct {
  Services  *Services
}


func (r *CharacterRouter) ServeHTTP(res http.ResponseWriter, req *http.Request) {
  var head string
  head, req.URL.Path = ShiftPath(req.URL.Path)

  switch req.Method {
  // GET Request Handlers
  case http.MethodGet:
    switch head {
    case "getall":
      r.handleGetAll(res, req)
    default:
      http.Error(res, fmt.Sprintf("Unsupported GET path %s", head), http.StatusBadRequest)
      return
    }
  // POST Request Handlers
  case http.MethodPost:
    switch head {
    case "create":
      r.handleCreate(res, req)
    case "update":
      r.handleUpdate(res, req)
    default:
      http.Error(res, fmt.Sprintf("Unsupported POST path %s", head), http.StatusBadRequest)
      return
    }
  // Unsupported Method Response
  default:
    http.Error(res, fmt.Sprintf("Unsupported Method type %s", req.Method), http.StatusBadRequest)
  }
}


// NewCharacterRouter makes a new api/character router and hooks up its services
func NewCharacterRouter(routerServices *Services) *CharacterRouter {
  router := new(CharacterRouter)

  router.Services = routerServices

  return router
}


/*---------------------------------
             Handlers
----------------------------------*/

func (r *CharacterRouter) handleGetAll(res http.ResponseWriter, req *http.Request) {
  characters, err := r.Services.Database.GetAllCharacters()
  if err != nil {
    http.Error(res, fmt.Sprintf("Error getting all characters from DB: %s", err.Error()), http.StatusInternalServerError)
    return
  }

  response := &Response{
    Success:  true,
    Error:    nil,
    Data:     CharacterGetAllResponseData{
      Characters:  characters,
    },
  }

  res.Header().Set("Content-Type", "application/json")
  json.NewEncoder(res).Encode(response)
}


func (r *CharacterRouter) handleCreate(res http.ResponseWriter, req *http.Request) {
  decoder := json.NewDecoder(req.Body)
  characterCreate := new(db.CharacterCreate)

  err := decoder.Decode(characterCreate)
  if err != nil {
    http.Error(res, fmt.Sprintf("Invalid JSON request: %s", err.Error()), http.StatusBadRequest)
    return
  }

  character, err := r.Services.Database.CreateCharacter(characterCreate)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error creating new character in database: %s", err.Error()), http.StatusInternalServerError)
    return
  }

  response := &Response{
    Success:  true,
    Error:    nil,
    Data:     CharacterCreateResponseData{
      Character:  character,
    },
  }

  res.Header().Set("Content-Type", "application/json")
  json.NewEncoder(res).Encode(response)
}


func (r *CharacterRouter) handleUpdate(res http.ResponseWriter, req *http.Request) {
  decoder := json.NewDecoder(req.Body)
  characterUpdate := new(db.CharacterUpdate)

  err := decoder.Decode(characterUpdate)
  if err != nil {
    http.Error(res, fmt.Sprintf("Invalid JSON request: %s", err.Error()), http.StatusBadRequest)
    return
  }

  character, err := r.Services.Database.UpdateCharacter(characterUpdate)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error updating character in database: %s", err.Error()), http.StatusInternalServerError)
    return
  }

  response := &Response{
    Success:  true,
    Error:    nil,
    Data:     CharacterUpdateResponseData{
      Character:  character,
    },
  }

  res.Header().Set("Content-Type", "application/json")
  json.NewEncoder(res).Encode(response)
}
