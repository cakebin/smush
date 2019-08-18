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

// TagCreateRequestData describes the data we're
// expecting when a user attemps to create a new tag
type TagCreateRequestData struct {
  TagName  string  `json:"tagName"`
}


// TagUpdateRequestData describes the data we're
// expecting when a user attemps update an existing tag
type TagUpdateRequestData struct {
  TagID    int     `json:"tagId"`
  TagName  string  `json:"tagName"`
}

/*---------------------------------
          Response Data
----------------------------------*/

// TagCreateResponseData describes the data we send
// back after successfully creating a new tag
type TagCreateResponseData struct {
  Tag  *Tag  `json:"tag"`
}


// TagUpdateResponseData describes the data we send
// back after successfully updating an existing tag
type TagUpdateResponseData struct {
  Tag  *Tag  `json:"tag"`
}


/*--------------------------------
          API <--> SQL
----------------------------------*/

// Tag is a translation from a db.Tag to an api.Tag
type Tag struct {
  TagID    int     `json:"tagId"`
  TagName  string  `json:"tagName"`
}


// ToAPITag maps from a db.Tag to an api.Tag
func ToAPITag(dbTag *db.Tag) *Tag {
  tag := new(Tag)
  tag.TagID = dbTag.TagID
  tag.TagName = dbTag.TagName

  return tag
}


// ToDBTagCreate maps from an api.TagCreateRequestData to a db.TagCreate
func ToDBTagCreate(tagCreateRequestData *TagCreateRequestData) *db.TagCreate {
  dbTagCreate := new(db.TagCreate)
  dbTagCreate.TagName = tagCreateRequestData.TagName

  return dbTagCreate
}


// ToDBTagUpdate maps from an api.TagUpdateRequestData to a db.TagUpdate
func ToDBTagUpdate(tagUpdateRequestData *TagUpdateRequestData) *db.TagUpdate {
  dbTagUpdate := new(db.TagUpdate)
  dbTagUpdate.TagID = tagUpdateRequestData.TagID
  dbTagUpdate.TagName = tagUpdateRequestData.TagName

  return dbTagUpdate
}

/*---------------------------------
             Router
----------------------------------*/

// TagRouter is responsible for serving /api/tag
type TagRouter struct {
  Services  *Services
}


func (r *TagRouter) ServeHTTP(res http.ResponseWriter, req *http.Request) {
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
    default:
      http.Error(res, fmt.Sprintf("Unsupported POST path %s", head), http.StatusBadRequest)
      return
    }
  // Unsupported Method Response
  default:
    http.Error(res, fmt.Sprintf("Unsupported Method type %s", req.Method), http.StatusBadRequest)
  }
}


// NewTagRouter makes a new api/tag router and hooks up its services
func NewTagRouter(routerServices *Services) *TagRouter {
  router := new(TagRouter)

  router.Services = routerServices

  return router
}


/*---------------------------------
             Handlers
----------------------------------*/


func (r *TagRouter) handleCreate(res http.ResponseWriter, req *http.Request) {
  decoder := json.NewDecoder(req.Body)
  createRequestData := new(TagCreateRequestData)

  err := decoder.Decode(createRequestData)
  if err != nil {
    http.Error(res, fmt.Sprintf("Invalid JSON request: %s", err.Error()), http.StatusBadRequest)
    return
  }

  dbTagCreate := ToDBTagCreate(createRequestData)
  tagID, err := r.Services.Database.CreateTag(dbTagCreate)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error creating new tag in database: %s", err.Error()), http.StatusInternalServerError)
    return
  }

  dbTag, err := r.Services.Database.GetTagByTagID(tagID)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error getting new tag in database: %s", err.Error()), http.StatusInternalServerError)
    return
  }
  tag := ToAPITag(dbTag)

  response := &Response{
    Success:  true,
    Error:    nil,
    Data:     TagCreateResponseData{
      Tag:  tag,
    },
  }

  res.Header().Set("Content-Type", "application/json")
  json.NewEncoder(res).Encode(response)
}



func (r *TagRouter) handleUpdate(res http.ResponseWriter, req *http.Request) {
  decoder := json.NewDecoder(req.Body)
  updateRequestData := new(TagUpdateRequestData)

  err := decoder.Decode(updateRequestData)
  if err != nil {
    http.Error(res, fmt.Sprintf("Invalid JSON request: %s", err.Error()), http.StatusBadRequest)
    return
  }

  dbTagUpdate := ToDBTagUpdate(updateRequestData)
  tagID, err := r.Services.Database.UpdateTag(dbTagUpdate)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error update tag in database: %s", err.Error()), http.StatusInternalServerError)
    return
  }

  dbTag, err := r.Services.Database.GetTagByTagID(tagID)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error getting updated tag in database: %s", err.Error()), http.StatusInternalServerError)
    return
  }
  tag := ToAPITag(dbTag)


  response := &Response{
    Success:  true,
    Error:    nil,
    Data:     TagUpdateResponseData{
      Tag:  tag,
    },
  }

  res.Header().Set("Content-Type", "application/json")
  json.NewEncoder(res).Encode(response)
}
