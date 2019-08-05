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

// CharacterCreateRequestData describes the data we're 
// expecting when a user attempts to create a character
type CharacterCreateRequestData struct {
  CharacterName       string   `json:"characterName"`
  CharacterStockImg   *string  `json:"characterStockImg,omitempty"`
  CharacterImg        *string  `json:"characterImg,omitmpty"`
  CharacterArchetype  *string  `json:"characterArchetype,omitempty"`
}


// CharacterUpdateRequestData describes the data we're 
// expecting when an admin user attempts to update a character
type CharacterUpdateRequestData struct {
  CharacterID         int      `json:"characterId"`
  CharacterName       *string  `json:"characterName,omitempty"`
  CharacterStockImg   *string  `json:"characterStockImg,omitempty"`
  CharacterImg        *string  `json:"characterImg,omitempty"`
  CharacterArchetype  *string  `json:"characterArchetype,omitempty"`
}


/*---------------------------------
          Response Data
----------------------------------*/

// CharacterGetAllResponseData is the data we send back
// after a successfully getting all character data
type CharacterGetAllResponseData struct {
  Characters  []*Character  `json:"characters"`
}


// CharacterCreateResponseData is the data we send
//  back after a successfully creating a new character
type CharacterCreateResponseData struct {
  Character  *Character  `json:"character"`
}


// CharacterUpdateResponseData is the data we send
//  back after a successfully udpating a new character
type CharacterUpdateResponseData struct {
  Character  *Character  `json:"character"`
}


/*--------------------------------
          API <--> SQL
----------------------------------*/

// Character is a translation from the SQL result
// which can have things like `sql.NullInt64`, so we
// need to translate that to regular JSON objects
type Character struct {
  CharacterID         int      `json:"characterId"`
  CharacterName       string   `json:"characterName"`
  CharacterStockImg   *string  `json:"characterStockImg,omitempty"`
  CharacterImg        *string  `json:"characterImg,omitempty"`
  CharacterArchetype  *string  `json:"characterArchetype,omitempty"`
}


// ToAPICharacter maps from a db.Character
// (which has things like sql.NullString) into
// an api.Character, which is JSON representable
func ToAPICharacter(dbChar *db.Character) *Character {
  character := new(Character)
  character.CharacterID = dbChar.CharacterID
  character.CharacterName = dbChar.CharacterName
  character.CharacterStockImg = FromNullString(dbChar.CharacterStockImg)
  character.CharacterImg = FromNullString(dbChar.CharacterImg)
  character.CharacterArchetype = FromNullString(dbChar.CharacterArchetype)

  return character
}


// ToDBCharacterCreate maps from an api.CharacterCreateRequestData
// to a db.CharacterCreate, which has fields like sql.NullString
func ToDBCharacterCreate(characterCreateRequestData *CharacterCreateRequestData) *db.CharacterCreate {
  dbCharacterCreate := new(db.CharacterCreate)
  dbCharacterCreate.CharacterName = characterCreateRequestData.CharacterName
  dbCharacterCreate.CharacterStockImg = ToNullString(characterCreateRequestData.CharacterStockImg)
  dbCharacterCreate.CharacterImg = ToNullString(characterCreateRequestData.CharacterImg)
  dbCharacterCreate.CharacterArchetype = ToNullString(characterCreateRequestData.CharacterArchetype)

  return dbCharacterCreate
}

// ToDBCharacterUpdate maps from an api.CharacterUpdateRequestData
// to a db.CharacterUpdate, which has fields like sql.NullString
func ToDBCharacterUpdate(characterUpdateRequestData *CharacterUpdateRequestData) *db.CharacterUpdate {
  dbCharacterUpdate := new(db.CharacterUpdate)
  dbCharacterUpdate.CharacterID = characterUpdateRequestData.CharacterID
  dbCharacterUpdate.CharacterName = ToNullString(characterUpdateRequestData.CharacterName)
  dbCharacterUpdate.CharacterStockImg = ToNullString(characterUpdateRequestData.CharacterStockImg)
  dbCharacterUpdate.CharacterImg = ToNullString(characterUpdateRequestData.CharacterImg)
  dbCharacterUpdate.CharacterArchetype = ToNullString(characterUpdateRequestData.CharacterArchetype)

  return dbCharacterUpdate
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


/*---------------------------------
             Handlers
----------------------------------*/

func (r *CharacterRouter) handleGetAll(res http.ResponseWriter, req *http.Request) {
  dbCharacters, err := r.SysUtils.Database.GetAllCharacters()
  if err != nil {
    http.Error(res, fmt.Sprintf("Error getting all characters from DB: %s", err.Error()), http.StatusInternalServerError)
    return
  }

  characters := make([]*Character, 0)
  for _, dbCharacter := range dbCharacters {
    character := ToAPICharacter(dbCharacter)
    characters = append(characters, character)
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
  createRequestData := new(CharacterCreateRequestData)

  err := decoder.Decode(createRequestData)
  if err != nil {
    http.Error(res, fmt.Sprintf("Invalid JSON request: %s", err.Error()), http.StatusBadRequest)
    return
  }

  dbCharacterCreate := ToDBCharacterCreate(createRequestData)
  dbCharacter, err := r.SysUtils.Database.CreateCharacter(dbCharacterCreate)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error creating new character in database: %s", err.Error()), http.StatusInternalServerError)
    return
  }
  character := ToAPICharacter(dbCharacter)

  response := Response{
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
  updateRequestData := new(CharacterUpdateRequestData)

  err := decoder.Decode(updateRequestData)
  if err != nil {
    http.Error(res, fmt.Sprintf("Invalid JSON request: %s", err.Error()), http.StatusBadRequest)
    return
  }

  dbCharacterUpdate := ToDBCharacterUpdate(updateRequestData)
  dbCharacter, err := r.SysUtils.Database.UpdateCharacter(dbCharacterUpdate)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error updating character in database: %s", err.Error()), http.StatusInternalServerError)
    return
  }
  character := ToAPICharacter(dbCharacter)

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
