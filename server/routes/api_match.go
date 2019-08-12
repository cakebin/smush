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
  UserID                int                  `json:"userId"`
  OpponentCharacterID   int                  `json:"opponentCharacterId"`
  UserCharacterID       *int64               `json:"userCharacterId,omitempty"`
  UserCharacterGsp      *int64               `json:"userCharacterGsp,omitempty"`
  OpponentCharacterGsp  *int64               `json:"opponentCharacterGsp,omitempty"`
  UserWin               *bool                `json:"userWin,omitempty"`
  MatchTags             *[]*MatchTagRequest  `json:"matchTags,omitempty"`
}


// MatchTagRequest describes the data we're expecting
// to create a new "match tag" relationship
type MatchTagRequest struct {
  MatchID  int  `json:"matchId"`
  TagID    int  `json:"tagId"`
}


// MatchUpdateRequestData describes the data we're
// expecting when a user attempt to update a match
type MatchUpdateRequestData struct {
  MatchID               int                  `json:"matchId"`
  OpponentCharacterID   *int64               `json:"opponentCharacterId,omitempty"`
  OpponentCharacterGsp  *int64               `json:"opponentCharacterGsp,omitempty"`
  MatchTags             *[]*MatchTagRequest  `json:"matchTags,omitempty"`
  UserCharacterID       *int64               `json:"userCharacterId,omitempty"`
  UserCharacterGsp      *int64               `json:"userCharacterGsp,omitempty"`
  UserWin               *bool                `json:"userWin,omitempty"`
  Created               *time.Time           `json:"created,omitempty"`
}


/*---------------------------------
          Response Data
----------------------------------*/

// MatchGetAllResponseData is the data we send back
// after a successfully get info for all matches in our db
type MatchGetAllResponseData struct {
  Matches  []*MatchMetaData  `json:"matches"`
}


// MatchMetaData combines the data from MatchView and MatchTagView
type MatchMetaData struct {
  Match      *MatchView       `json:"match"`
  MatchTags  []*MatchTagView  `json:"matchTags"`
}


// MatchCreateResponseData is the data we send
// back after a successfully creating a new match
type MatchCreateResponseData struct {
  Match      *MatchView       `json:"match"`
  MatchTags  []*MatchTagView  `json:"matchTags"`
}


// MatchUpdateResponseData is the data we send
// back after successfully updating a match
type MatchUpdateResponseData struct {
  Match      *MatchView       `json:"match"`
  MatchTags  []*MatchTagView  `json:"matchTags"`
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
  UserCharacterGsp       *int64     `json:"userCharacterGsp,omitempty"`
  UserWin                *bool      `json:"userWin,omitempty"`

  // Data from users
  UserName               string     `json:"userName"`

  // Data from characters
  OpponentCharacterName  string     `json:"opponentCharacterName"`
  UserCharacterName      *string    `json:"userCharacterName,omitempty"`
}


// MatchTagView is a translation from the SQL result
// which can have things like `sql.NullInt64`, so we
// need to translate that to regular JSON objects
type MatchTagView struct {
  MatchTagID  int     `json:"matchTagId"`
  MatchID     int     `json:"matchId"`
  TagID       int     `json:"tagId"`
  TagName     string  `json:"tagName"`
}


// ToAPIMatchView maps from a db.MatchView
// (which can have things like sql.NullString) into
// an api.MatchView, which is JSON representable
func ToAPIMatchView(dbMatchView *db.MatchView) *MatchView {
  matchView := new(MatchView)
  matchView.Created = dbMatchView.Created
  matchView.UserID = dbMatchView.UserID
  matchView.MatchID = dbMatchView.MatchID
  matchView.OpponentCharacterID = dbMatchView.OpponentCharacterID
  matchView.UserCharacterID = FromNullInt64(dbMatchView.UserCharacterID)
  matchView.OpponentCharacterGsp = FromNullInt64(dbMatchView.OpponentCharacterGsp)
  matchView.UserCharacterGsp = FromNullInt64(dbMatchView.UserCharacterGsp)
  matchView.UserName = dbMatchView.UserName
  matchView.OpponentCharacterName = dbMatchView.OpponentCharacterName
  matchView.UserCharacterName = FromNullString(dbMatchView.UserCharacterName)
  matchView.UserWin = FromNullBool(dbMatchView.UserWin)

  return matchView
}


// ToAPIMatchViews maps from a []*db.MatchView to a []*api.MatchView
func ToAPIMatchViews(dbMatchViews []*db.MatchView) []*MatchView {
  matchViews := make([]*MatchView, 0)

  for _, dbMatchView := range dbMatchViews {
    matchView := ToAPIMatchView(dbMatchView)
    matchViews = append(matchViews, matchView)
  }

  return matchViews
}


// ToAPIMatchMetadata prepares match metadata
func ToAPIMatchMetadata(matchView *MatchView, matchTagViews []*MatchTagView) *MatchMetaData {
  matchMetaData := new(MatchMetaData)
  matchMetaData.Match = matchView
  matchMetaData.MatchTags = matchTagViews

  return matchMetaData
}


// ToAllAPIMatchMetadata prepares all match metadata given a full list of matchViews and matchTagViews
func ToAllAPIMatchMetadata(allMatchViews []*MatchView, allMatchTagViews []*MatchTagView) []*MatchMetaData {
  allMatchMetadata := make([]*MatchMetaData, 0)

  for _, matchView := range allMatchViews {
    filteredMatchTagViews := FilterMatchTagViewsByMatchID(allMatchTagViews, matchView.MatchID)
    matchMetadata := ToAPIMatchMetadata(matchView, filteredMatchTagViews)
    allMatchMetadata = append(allMatchMetadata, matchMetadata)
  }

  return allMatchMetadata
}


// FilterMatchTagViewsByMatchID finds the MatchTagViews associated with a given matchID
func FilterMatchTagViewsByMatchID(matchTagViews []*MatchTagView, matchID int) []*MatchTagView {
  filteredMatchTagViews := make([]*MatchTagView, 0)

  for _, matchTagView := range matchTagViews {
    if matchTagView.MatchID == matchID {
      filteredMatchTagViews = append(filteredMatchTagViews, matchTagView)
    }
  }

  return filteredMatchTagViews
}


// ToAPIMatchTagView maps from a db.MatchTagView
// (which can have things like sql.NullString) into
// an api.MatchTagView, which is JSON representable
func ToAPIMatchTagView(dbMatchTagView *db.MatchTagView) *MatchTagView {
  matchTagView := new(MatchTagView)
  matchTagView.MatchTagID = dbMatchTagView.MatchTagID
  matchTagView.MatchID = dbMatchTagView.MatchID
  matchTagView.TagID = dbMatchTagView.TagID
  matchTagView.TagName = dbMatchTagView.TagName

  return matchTagView
}


// ToAPIMatchTagViews maps from a []*db.MatchTagView to a []*api.MatchTagView
func ToAPIMatchTagViews(dbMatchTagViews []*db.MatchTagView) []*MatchTagView {
  matchTagViews := make([]*MatchTagView, 0)

  for _, dbMatchTagView := range dbMatchTagViews {
    matchTagView := ToAPIMatchTagView(dbMatchTagView)
    matchTagViews = append(matchTagViews, matchTagView)
  }

  return matchTagViews
}

// ToDBMatchCreate maps from an api.MatchCreateRequestDat
// to a db.MatchCreate, which has fields like sql.NullBool
func ToDBMatchCreate(matchCreateRequestData *MatchCreateRequestData) *db.MatchCreate {
  dbMatchCreate := new(db.MatchCreate)
  dbMatchCreate.UserID = matchCreateRequestData.UserID
  dbMatchCreate.OpponentCharacterID = matchCreateRequestData.OpponentCharacterID
  dbMatchCreate.OpponentCharacterGsp = ToNullInt64(matchCreateRequestData.OpponentCharacterGsp)
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
  dbMatchUpdate.UserCharacterID = ToNullInt64(matchUpdateRequestData.UserCharacterID)
  dbMatchUpdate.UserCharacterGsp = ToNullInt64(matchUpdateRequestData.UserCharacterGsp)
  dbMatchUpdate.UserWin = ToNullBool(matchUpdateRequestData.UserWin)
  dbMatchUpdate.Created = ToNullTime(matchUpdateRequestData.Created)

  return dbMatchUpdate
}


// ToDBMatchTagCreate maps from an api.MatchTag to a db.MatchTagCreate
func ToDBMatchTagCreate(matchTagRequest *MatchTagRequest) *db.MatchTagCreate {
  dbMatchTagCreate := new(db.MatchTagCreate)
  dbMatchTagCreate.MatchID = matchTagRequest.MatchID
  dbMatchTagCreate.TagID = matchTagRequest.TagID

  return dbMatchTagCreate
}


// ToDBMatchTagsCreate maps a slice of api.MatchTag to a slice of db.MatchTagCreate
func ToDBMatchTagsCreate(matchTagRequests *[]*MatchTagRequest) []*db.MatchTagCreate {
  dbMatchTagsCreate := make([]*db.MatchTagCreate, 0)

  for _, matchTagRequest := range *matchTagRequests {
    dbMatchTagCreate := ToDBMatchTagCreate(matchTagRequest)
    dbMatchTagsCreate = append(dbMatchTagsCreate, dbMatchTagCreate)
  }

  return dbMatchTagsCreate
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
  dbMatchViews, err := r.Services.Database.GetAllMatchViews()
  if err != nil {
    http.Error(res, fmt.Sprintf("Error getting all matches from DB: %s", err.Error()), http.StatusInternalServerError)
    return
  }

  dbMatchTagViews, err := r.Services.Database.GetAllMatchTagViews()
  if err != nil {
    http.Error(res, fmt.Sprintf("Error getting all matches tags from DB: %s", err.Error()), http.StatusInternalServerError)
    return
  }

  matchViews := ToAPIMatchViews(dbMatchViews)
  matchTagViews := ToAPIMatchTagViews(dbMatchTagViews)
  matchMetadata := ToAllAPIMatchMetadata(matchViews, matchTagViews)

  response := &Response{
    Success:  true,
    Error:    nil,
    Data:     MatchGetAllResponseData{
      Matches:  matchMetadata,
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
  if err != nil {
    http.Error(res, fmt.Sprintf("Error creating new match: %s", err.Error()), http.StatusInternalServerError)
    return
  }

  // Then make any match tag relationships
  if len(*createRequestData.MatchTags) == 0 {
    dbMatchTagsCreate := ToDBMatchTagsCreate(createRequestData.MatchTags)
    _, err := r.Services.Database.CreateMatchTags(dbMatchTagsCreate)
    if err != nil {
      http.Error(res, fmt.Sprintf("Error creating new match tags: %s", err.Error()), http.StatusInternalServerError)
      return
    }
  }

  dbMatchView, err := r.Services.Database.GetMatchViewByMatchID(matchID)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error getting match view: %s", err.Error()), http.StatusInternalServerError)
    return
  }
  dbMatchTagViews, err := r.Services.Database.GetMatchTagViewsByMatchID(matchID)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error getting match tag view: %s", err.Error()), http.StatusInternalServerError)
    return
  }

  // Finally make a JSON representable version of the db.MatchView fetched results
  matchView := ToAPIMatchView(dbMatchView)
  matchTagViews := ToAPIMatchTagViews(dbMatchTagViews)

  response := &Response{
    Success:  true,
    Error:    nil,
    Data:     MatchCreateResponseData{
      Match:      matchView,
      MatchTags:  matchTagViews,
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

  // First update match specific info
  dbMatchUpdate := ToDBMatchUpdate(updateRequestData)
  matchID, err := r.Services.Database.UpdateMatch(dbMatchUpdate)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error updating match in database: %s", err.Error()), http.StatusInternalServerError)
    return
  }

  // Then make update match tag relationships
  if len(*updateRequestData.MatchTags) == 0 {
    // Delete older match tag relationships 
    _, err := r.Services.Database.DeleteMatchTagsByMatchID(matchID)
    if err != nil {
      http.Error(res, fmt.Sprintf("Error deleting match tags in database: %s", err.Error()), http.StatusInternalServerError)
      return
    }

    // Then nake new match tag relatioships
    dbMatchTagsCreate := ToDBMatchTagsCreate(updateRequestData.MatchTags)
    _, err = r.Services.Database.CreateMatchTags(dbMatchTagsCreate)
    if err != nil {
      http.Error(res, fmt.Sprintf("Error creating new match tags: %s", err.Error()), http.StatusInternalServerError)
      return
    }
  }

  dbMatchView, err := r.Services.Database.GetMatchViewByMatchID(matchID)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error getting match view: %s", err.Error()), http.StatusInternalServerError)
    return
  }
  dbMatchTagViews, err := r.Services.Database.GetMatchTagViewsByMatchID(matchID)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error getting match tag view: %s", err.Error()), http.StatusInternalServerError)
    return
  }

  matchView := ToAPIMatchView(dbMatchView)
  matchTagViews := ToAPIMatchTagViews(dbMatchTagViews)

  response := &Response{
    Success:   true,
    Error:     nil,
    Data:      MatchUpdateResponseData{
      Match:      matchView,
      MatchTags:  matchTagViews,
    },
  }

  res.Header().Set("Content-Type", "application/json")
  json.NewEncoder(res).Encode(response)
}
