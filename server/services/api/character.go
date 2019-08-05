package api

import (
  "encoding/json"
  "fmt"
  "net/http"

  "github.com/cakebin/smush/server/env"
  "github.com/cakebin/smush/server/services/db"
  "github.com/cakebin/smush/server/util/routing"
)


/*---------------------------------
          Request Data
----------------------------------*/

// CharacterCreateRequest describes the data we're 
// expecting when a user attempts to create a character
type CharacterCreateRequest struct {
  CharacterName      string   `json:"characterName"`
  CharacterStockImg  *string  `json:"characterStockImg,omitempty"`
}


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
  CharacterID  int  `json:"characterId"`
}


/*---------------------------------
             Router
----------------------------------*/

// CharacterRouter is responsible for serving "/api/character"
// Basically, connecting to our Postgres DB for all
// of the CRUD operations for our "Character" models
type CharacterRouter struct {
  SysUtils  *env.SysUtils
}


func (r *CharacterRouter) ServeHTTP(res http.ResponseWriter, req *http.Request) {
  var head string
  head, req.URL.Path = routing.ShiftPath(req.URL.Path)

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
    default:
      http.Error(res, fmt.Sprintf("Unsupported POST path %s", head), http.StatusBadRequest)
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

func (r *CharacterRouter) handleGetAll(res http.ResponseWriter, req *http.Request) {
  characters, err := r.SysUtils.Database.GetAllCharacters()
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
  var createRequest CharacterCreateRequest

  err := decoder.Decode(&createRequest)
  if err != nil {
    http.Error(res, fmt.Sprintf("Invalid JSON request: %s", err.Error()), http.StatusBadRequest)
    return
  }

  var character db.Character
  character.CharacterName = createRequest.CharacterName
  if createRequest.CharacterStockImg != nil {
    character.CharacterStockImg = *createRequest.CharacterStockImg
  }
  characterID, err := r.SysUtils.Database.CreateCharacter(character)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error creating new character in database: %s", err.Error()), http.StatusInternalServerError)
    return
  }

  response := Response{
    Success:  true,
    Error:    nil,
    Data:     CharacterCreateResponseData{
      CharacterID:  characterID,
    },
  }

  res.Header().Set("Content-Type", "application/json")
  json.NewEncoder(res).Encode(response)
}
