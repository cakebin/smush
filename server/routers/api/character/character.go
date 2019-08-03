package character


import (
  "encoding/json"
  "fmt"
  "net/http"

  "github.com/cakebin/smush/server/api"
  "github.com/cakebin/smush/server/env"
  "github.com/cakebin/smush/server/services/db"
  "github.com/cakebin/smush/server/util/routing"
)


// Router is responsible for serving "/api/character"
// Basically, connecting to our Postgres DB for all
// of the CRUD operations for our "Character" models
type Router struct {
  SysUtils  *env.SysUtils
}


func (r *Router) ServeHTTP(res http.ResponseWriter, req *http.Request) {
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


func (r *Router) handleGetAll(res http.ResponseWriter, req *http.Request) {
  characters, err := r.SysUtils.Database.GetAllCharacters()
  if err != nil {
    http.Error(res, fmt.Sprintf("Error getting all characters from DB: %s", err.Error()), http.StatusInternalServerError)
    return
  }

  res.Header().Set("Content-Type", "application/json")
  json.NewEncoder(res).Encode(characters)
}


func (r *Router) handleCreate(res http.ResponseWriter, req *http.Request) {
  decoder := json.NewDecoder(req.Body)
  var character db.Character

  err := decoder.Decode(&character)
  if err != nil {
    http.Error(res, fmt.Sprintf("Invalid JSON request: %s", err.Error()), http.StatusBadRequest)
    return
  }

  characterID, err := r.SysUtils.Database.CreateCharacter(character)

  response := &api.Response{
    Success: true,
    Error: err,
    Data: characterID,
  }

  res.Header().Set("Content-Type", "application/json")
  json.NewEncoder(res).Encode(response)
}


// NewRouter makes a new match router with access to the "SysUtils" environment object
func NewRouter(sysUtils *env.SysUtils) *Router {
  router := new(Router)
  router.SysUtils = sysUtils
  return router
}
