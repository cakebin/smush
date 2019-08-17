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

// MatchGetAllResponseData is the data we send back
// after a successfully get info for all matches in our db
type MatchGetAllResponseData struct {
  Matches  []*db.MatchView  `json:"matches"`
}


// MatchCreateResponseData is the data we send
// back after a successfully creating a new match
type MatchCreateResponseData struct {
  Match  *db.MatchView  `json:"match"`
}


// MatchUpdateResponseData is the data we send
// back after successfully updating a match
type MatchUpdateResponseData struct {
  Match  *db.MatchView  `json:"match"`
}


/*---------------------------------
             Router
----------------------------------*/

// MatchRouter is responsible for serving "/api/matches"
// Basically, connecting to our Postgres DB for all
// of the CRUD operations for our "Match" models
type MatchRouter struct {
  Services  *Services
}


func (r *MatchRouter) ServeHTTP(res http.ResponseWriter, req *http.Request) {
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
      http.Error(res, fmt.Sprintf("Unsupport POST path %s", head), http.StatusBadRequest)
      return
    }
  // Unsupported Method Response
  default:
    http.Error(res, fmt.Sprintf("Unsupport Method type %s", req.Method), http.StatusBadRequest)
  }
}


// NewMatchRouter makes a new api/match router and hooks up its services
func NewMatchRouter(routerServices *Services) *MatchRouter {
  router := new(MatchRouter)

  router.Services = routerServices

  return router
}


/*---------------------------------
             Handlers
----------------------------------*/

func (r *MatchRouter) handleGetAll(res http.ResponseWriter, req *http.Request) {
  matchViews, err := r.Services.Database.GetAllMatchViews()
  if err != nil {
    http.Error(res, fmt.Sprintf("Error getting all matches from DB: %s", err.Error()), http.StatusInternalServerError)
    return
  }

  response := &Response{
    Success:  true,
    Error:    nil,
    Data:     MatchGetAllResponseData{
      Matches:  matchViews,
    },
  }

  res.Header().Set("Content-Type", "application/json")
  json.NewEncoder(res).Encode(response)
}


func (r *MatchRouter) handleCreate(res http.ResponseWriter, req *http.Request) {
  decoder := json.NewDecoder(req.Body)
  matchCreate := new(db.MatchCreate)

  err := decoder.Decode(matchCreate)
  if err != nil {
    http.Error(res, fmt.Sprintf("Invalid JSON request: %s", err.Error()), http.StatusBadRequest)
    return
  }

  // Make the new match and fetch relevant match view data for it
  matchID, err := r.Services.Database.CreateMatch(matchCreate)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error creating new match: %s", err.Error()), http.StatusInternalServerError)
    return
  }
  matchView, err := r.Services.Database.GetMatchViewByMatchID(matchID)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error getting new match view: %s", err.Error()), http.StatusInternalServerError)
    return
  }

  response := &Response{
    Success:  true,
    Error:    err,
    Data:     MatchCreateResponseData{
      Match:  matchView,
    },
  }

  res.Header().Set("Content-Type", "application/json")
  json.NewEncoder(res).Encode(response)
}


func (r *MatchRouter) handleUpdate(res http.ResponseWriter, req *http.Request) {
  decoder := json.NewDecoder(req.Body)
  matchUpdate := new(db.MatchUpdate)

  err := decoder.Decode(matchUpdate)
  if err != nil {
    http.Error(res, fmt.Sprintf("Invalid JSON request: %s", err.Error()), http.StatusBadRequest)
    return
  }

  matchID, err := r.Services.Database.UpdateMatch(matchUpdate)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error updating match in database: %s", err.Error()), http.StatusInternalServerError)
    return
  }
  matchView, err := r.Services.Database.GetMatchViewByMatchID(matchID)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error getting updated match view: %s", err.Error()), http.StatusInternalServerError)
    return
  }

  response := &Response{
    Success:   true,
    Error:     nil,
    Data:      MatchUpdateResponseData{
      Match:  matchView,
    },
  }

  res.Header().Set("Content-Type", "application/json")
  json.NewEncoder(res).Encode(response)
}
