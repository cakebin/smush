package api

import (
  "encoding/json"
  "fmt"
  "net/http"
  "time"

  "github.com/cakebin/smush/server/env"
  "github.com/cakebin/smush/server/services/db"
  "github.com/cakebin/smush/server/util/routing"
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

/*---------------------------------
          SQL --> API
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


/*---------------------------------
          Response Data
----------------------------------*/

// MatchGetAllResponseData is the data we send back
// after a successfully get info for all matches in our db
type MatchGetAllResponseData struct {
  Matches  []MatchView  `json:"matches"`
}


// MatchCreateResponseData is the data we send
//  back after a successfully creating a new match
type MatchCreateResponseData struct {
  Match  MatchView  `json:"match"`
}


/*---------------------------------
             Router
----------------------------------*/

// MatchRouter is responsible for serving "/api/matches"
// Basically, connecting to our Postgres DB for all
// of the CRUD operations for our "Match" models
type MatchRouter struct {
  SysUtils  *env.SysUtils
}


func (r *MatchRouter) ServeHTTP(res http.ResponseWriter, req *http.Request) {
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
  dbMatchViews, err := r.SysUtils.Database.GetAllMatchViews()
  if err != nil {
    http.Error(res, fmt.Sprintf("Error getting all matches from DB: %s", err.Error()), http.StatusInternalServerError)
    return
  }

  matchViews := make([]MatchView, 0)
  for _, dbMatchView := range dbMatchViews {
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
    matchView.UserName = dbMatchView.UserName
    matchView.OpponentCharacterName = dbMatchView.OpponentCharacterName
    matchView.UserCharacterName = FromNullString(dbMatchView.UserCharacterName)
    matchView.UserWin = FromNullBool(dbMatchView.UserWin)
    matchViews = append(matchViews, *matchView)
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
  var createRequestData MatchCreateRequestData

  err := decoder.Decode(&createRequestData)
  if err != nil {
    http.Error(res, fmt.Sprintf("Invalid JSON request: %s", err.Error()), http.StatusBadRequest)
    return
  }

  var match db.Match
  match.UserID = createRequestData.UserID
  match.OpponentCharacterID = createRequestData.OpponentCharacterID
  // Convert optional fields to sql.Null variants
  match.OpponentCharacterGsp = ToNullInt64(createRequestData.OpponentCharacterGsp)
  match.OpponentTeabag = ToNullBool(createRequestData.OpponentTeabag)
  match.OpponentCamp = ToNullBool(createRequestData.OpponentCamp)
  match.OpponentAwesome = ToNullBool(createRequestData.OpponentAwesome)
  match.UserCharacterID = ToNullInt64(createRequestData.UserCharacterID)
  match.UserCharacterGsp = ToNullInt64(createRequestData.UserCharacterGsp)
  match.UserWin = ToNullBool(createRequestData.UserWin)
  
  matchID, err := r.SysUtils.Database.CreateMatch(match)
  dbMatchView, err := r.SysUtils.Database.GetMatchViewByMatchID(matchID)

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
  matchView.UserName = dbMatchView.UserName
  matchView.OpponentCharacterName = dbMatchView.OpponentCharacterName
  matchView.UserCharacterName = FromNullString(dbMatchView.UserCharacterName)
  matchView.UserWin = FromNullBool(dbMatchView.UserWin)

  response := &Response{
    Success:  true,
    Error:    err,
    Data:     MatchCreateResponseData{
      Match:  *matchView,
    },
  }

  res.Header().Set("Content-Type", "application/json")
  json.NewEncoder(res).Encode(response)
}
