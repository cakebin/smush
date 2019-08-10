package routes

import (
  "encoding/json"
  "fmt"
  "net/http"
  "time"

  "github.com/cakebin/smush/server/services/db"
)


/*---------------------------------
          Request Data
----------------------------------*/

// MatchCreateRequestData describes the data we're 
// expecting when a user attempts to create a match
type MatchCreateRequestData struct {
  UserID                int     `json:"userId"`
  OpponentCharacterID   int     `json:"opponentCharacterId"`
  UserCharacterID       *int64  `json:"userCharacterId,omitempty"`
  UserCharacterGsp      *int64  `json:"userCharacterGsp,omitempty"`
  OpponentCharacterGsp  *int64  `json:"opponentCharacterGsp,omitempty"`
  UserWin               *bool   `json:"userWin,omitempty"`
  OpponentTeabag        *bool   `json:"opponentTeabag,omitempty"`
  OpponentCamp          *bool   `json:"opponentCamp,omitempty"`
  OpponentAwesome       *bool   `json:"opponentAwesome,omitempty"`
}


// MatchUpdateRequestData describes the data we're
// expecting when a user attempt to update a match
type MatchUpdateRequestData struct {
  MatchID               int         `json:"matchId"`
  OpponentCharacterID   *int64      `json:"opponentCharacterId"`
  OpponentCharacterGsp  *int64      `json:"opponentCharacterGsp"`
  OpponentTeabag        *bool       `json:"opponentTeabag"`
  OpponentCamp          *bool       `json:"opponentCamp"`
  OpponentAwesome       *bool       `json:"opponentAwesome"`
  UserCharacterID       *int64      `json:"userCharacterId"`
  UserCharacterGsp      *int64      `json:"userCharacterGsp"`
  UserWin               *bool       `json:"userWin"`
  Created               *time.Time  `json:"created"`
}


/*---------------------------------
          Response Data
----------------------------------*/

// MatchGetAllResponseData is the data we send back
// after a successfully get info for all matches in our db
type MatchGetAllResponseData struct {
  Matches  []*MatchView  `json:"matches"`
}


// MatchCreateResponseData is the data we send
// back after a successfully creating a new match
type MatchCreateResponseData struct {
  Match  *MatchView  `json:"match"`
}


// MatchUpdateResponseData is the data we send
// back after successfully updating a match
type MatchUpdateResponseData struct {
  Match  *MatchView  `json:"match"`
}


/*---------------------------------
          API <--> SQL
----------------------------------*/

// MatchView is a translation from the SQL result
// which can have things like `sql.NullInt64`, so we
// need to translate that to regular JSON objects
type MatchView struct {
  // Data from matches
  Created                time.Time  `json:"created"`
  UserID                 int        `json:"userId"`
  MatchID                int        `json:"matchId"`
  OpponentCharacterID    int        `json:"opponentCharacterId"`

  UserCharacterID        *int64     `json:"userCharacterId"`
  OpponentCharacterGsp   *int64     `json:"opponentCharacterGsp,omitempty"`
  OpponentTeabag         *bool      `json:"opponentTeabag,omitempty"`
  OpponentCamp           *bool      `json:"opponentCamp,omitempty"`
  OpponentAwesome        *bool      `json:"opponentAwesome,omitempty"`
  UserCharacterGsp       *int64     `json:"userCharacterGsp,omitempty"`
  UserWin                *bool      `json:"userWin,omitempty"`

  // Data from users
  UserName               string     `json:"userName"`

  // Data from characters
  OpponentCharacterName  string     `json:"opponentCharacterName"`
  UserCharacterName      *string    `json:"userCharacterName,omitempty"`
}


// ToAPIMatchView maps from a db.MatchView
// (which has things like sql.NullString) into
// an api.MatchView, which is JSON representable
func ToAPIMatchView(dbMatchView *db.MatchView) *MatchView {
  matchView := new(MatchView)
  matchView.Created = dbMatchView.Created
  matchView.UserID = dbMatchView.UserID
  matchView.MatchID = dbMatchView.MatchID
  matchView.OpponentCharacterID = dbMatchView.OpponentCharacterID
  matchView.UserCharacterID = FromNullInt64(dbMatchView.UserCharacterID)
  matchView.OpponentCharacterGsp = FromNullInt64(dbMatchView.OpponentCharacterGsp)
  matchView.OpponentTeabag = FromNullBool(dbMatchView.OpponentTeabag)
  matchView.OpponentCamp = FromNullBool(dbMatchView.OpponentCamp)
  matchView.OpponentAwesome = FromNullBool(dbMatchView.OpponentAwesome)
  matchView.UserCharacterGsp = FromNullInt64(dbMatchView.UserCharacterGsp)
  matchView.UserName = dbMatchView.UserName
  matchView.OpponentCharacterName = dbMatchView.OpponentCharacterName
  matchView.UserCharacterName = FromNullString(dbMatchView.UserCharacterName)
  matchView.UserWin = FromNullBool(dbMatchView.UserWin)

  return matchView
}


// ToDBMatchCreate maps from an api.MatchCreateRequestDat
// to a db.MatchCreate, which has fields like sql.NullBool
func ToDBMatchCreate(matchCreateRequestData *MatchCreateRequestData) *db.MatchCreate {
  dbMatchCreate := new(db.MatchCreate)
  dbMatchCreate.UserID = matchCreateRequestData.UserID
  dbMatchCreate.OpponentCharacterID = matchCreateRequestData.OpponentCharacterID
  dbMatchCreate.OpponentCharacterGsp = ToNullInt64(matchCreateRequestData.OpponentCharacterGsp)
  dbMatchCreate.OpponentTeabag = ToNullBool(matchCreateRequestData.OpponentTeabag)
  dbMatchCreate.OpponentCamp = ToNullBool(matchCreateRequestData.OpponentCamp)
  dbMatchCreate.OpponentAwesome = ToNullBool(matchCreateRequestData.OpponentAwesome)
  dbMatchCreate.UserCharacterID = ToNullInt64(matchCreateRequestData.UserCharacterID)
  dbMatchCreate.UserCharacterGsp = ToNullInt64(matchCreateRequestData.UserCharacterGsp)
  dbMatchCreate.UserWin = ToNullBool(matchCreateRequestData.UserWin)

  return dbMatchCreate
}


// ToDBMatchUpdate maps from an api.MatchUpdateRequestData
// to a db.MatchUpdate, which has fields like sql.NullBool
func ToDBMatchUpdate(matchUpdateRequestData *MatchUpdateRequestData) *db.MatchUpdate {
  dbMatchUpdate := new(db.MatchUpdate)
  dbMatchUpdate.MatchID = matchUpdateRequestData.MatchID
  dbMatchUpdate.OpponentCharacterID = ToNullInt64(matchUpdateRequestData.OpponentCharacterID)
  dbMatchUpdate.OpponentCharacterGsp = ToNullInt64(matchUpdateRequestData.OpponentCharacterGsp)
  dbMatchUpdate.OpponentTeabag = ToNullBool(matchUpdateRequestData.OpponentTeabag)
  dbMatchUpdate.OpponentCamp = ToNullBool(matchUpdateRequestData.OpponentCamp)
  dbMatchUpdate.OpponentAwesome = ToNullBool(matchUpdateRequestData.OpponentAwesome)
  dbMatchUpdate.UserCharacterID = ToNullInt64(matchUpdateRequestData.UserCharacterID)
  dbMatchUpdate.UserCharacterGsp = ToNullInt64(matchUpdateRequestData.UserCharacterGsp)
  dbMatchUpdate.UserWin = ToNullBool(matchUpdateRequestData.UserWin)
  dbMatchUpdate.Created = ToNullTime(matchUpdateRequestData.Created)

  return dbMatchUpdate
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


/*---------------------------------
             Handlers
----------------------------------*/

func (r *MatchRouter) handleGetAll(res http.ResponseWriter, req *http.Request) {
  dbMatchViews, err := r.Services.Database.GetAllMatchViews()
  if err != nil {
    http.Error(res, fmt.Sprintf("Error getting all matches from DB: %s", err.Error()), http.StatusInternalServerError)
    return
  }

  matchViews := make([]*MatchView, 0)
  for _, dbMatchView := range dbMatchViews {
    matchView := ToAPIMatchView(dbMatchView)
    matchViews = append(matchViews, matchView)
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
  createRequestData := new(MatchCreateRequestData)

  err := decoder.Decode(createRequestData)
  if err != nil {
    http.Error(res, fmt.Sprintf("Invalid JSON request: %s", err.Error()), http.StatusBadRequest)
    return
  }

  // First make a db.MatchCreate from the create request data
  dbMatchCreate := ToDBMatchCreate(createRequestData)

  // Then make the new match and fetch relevant match view data for it
  matchID, err := r.Services.Database.CreateMatch(dbMatchCreate)
  dbMatchView, err := r.Services.Database.GetMatchViewByMatchID(matchID)

  // Finally make a JSON representable version of the db.MatchView fetched results
  matchView := ToAPIMatchView(dbMatchView)

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
  updateRequestData := new(MatchUpdateRequestData)

  err := decoder.Decode(updateRequestData)
  if err != nil {
    http.Error(res, fmt.Sprintf("Invalid JSON request: %s", err.Error()), http.StatusBadRequest)
    return
  }

  dbMatchUpdate := ToDBMatchUpdate(updateRequestData)
  matchID, err := r.Services.Database.UpdateMatch(dbMatchUpdate)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error updating match in database: %s", err.Error()), http.StatusInternalServerError)
    return
  }
  dbMatchView, err := r.Services.Database.GetMatchViewByMatchID(matchID)
  matchView := ToAPIMatchView(dbMatchView)

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


// NewMatchRouter makes a new api/match router and hooks up its services
func NewMatchRouter(routerServices *Services) *MatchRouter {
  router := new(MatchRouter)

  router.Services = routerServices

  return router
}
