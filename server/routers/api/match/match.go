package match

import (
  "encoding/json"
  "fmt"
  "net/http"
  "strconv"

  "github.com/cakebin/smush/server/db"
  "github.com/cakebin/smush/server/env"
  "github.com/cakebin/smush/server/util/routing"
)


// Router is responsible for serving "/api/matches"
// Basically, connecting to our Postgres DB for all
// of the CRUD operations for our "Match" models
type Router struct {
  SysUtils *env.SysUtils
}


func (r *Router) ServeHTTP(res http.ResponseWriter, req *http.Request) {
  var head string
  head, req.URL.Path = routing.ShiftPath(req.URL.Path)

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
    case "create":
      r.handleCreate(res, req)
    }
  // Unsupported Method Response
  default:
    http.Error(res, fmt.Sprintf("Unsupport Method type %s", req.Method), http.StatusBadRequest)
  }
}


func (r *Router) handleGetByID(res http.ResponseWriter, req *http.Request) {
  var head string
  head, req.URL.Path = routing.ShiftPath(req.URL.Path)

  id, err := strconv.Atoi(head)
  if err != nil {
    http.Error(res, fmt.Sprintf("Invalid match id :%s", head), http.StatusBadRequest)
    return
  }

  match, err := r.SysUtils.Database.GetMatchByID(id)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error getting match with id %q: :%s", id, err.Error()), http.StatusInternalServerError)
  }

  res.Header().Set("Content-Type", "application/json")
  json.NewEncoder(res).Encode(match)
}


func (r *Router) handleGetAll(res http.ResponseWriter, req *http.Request) {
  matches, err := r.SysUtils.Database.GetAllMatches()
  if err != nil {
    http.Error(res, fmt.Sprintf("Error getting all matches from DB: %s", err.Error()), http.StatusInternalServerError)
    return
  }

  res.Header().Set("Content-Type", "application/json")
  json.NewEncoder(res).Encode(matches)
}


func (r *Router) handleCreate(res http.ResponseWriter, req *http.Request) {
  decoder := json.NewDecoder(req.Body)
  var match db.Match

  err := decoder.Decode(&match)
  if err != nil {
    http.Error(res, fmt.Sprintf("Invalid JSON request: %s", err.Error()), http.StatusBadRequest)
    return
  }

  success, err := r.SysUtils.Database.CreateMatch(match)

  matchResponse := &db.MatchResponse{
    Success: success,
    Error: err,
  }

  res.Header().Set("Content-Type", "application/json")
  json.NewEncoder(res).Encode(matchResponse)
}


// NewRouter makes a new match router with access to the "SysUtils" environment object
func NewRouter(sysUtils *env.SysUtils) *Router {
  router := new(Router)
  router.SysUtils = sysUtils
  return router
}
