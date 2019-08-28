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

// TagGetAllResponseData is the data we send back
// after a successfully getiing all tags in our db
type TagGetAllResponseData struct {
  Tags  []*db.Tag  `json:"tags"`
}


// TagCreateResponseData describes the data we send
// back after successfully creating a new tag
type TagCreateResponseData struct {
  Tag  *db.Tag  `json:"tag"`
}


// TagUpdateResponseData describes the data we send
// back after successfully updating an existing tag
type TagUpdateResponseData struct {
  Tag  *db.Tag  `json:"tag"`
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
    // Check for admin role
    accessCookie, err := req.Cookie("smush-access-token")
    if err != nil {
      http.Error(res, fmt.Sprintf("Access token expired; can't get user ID from cookie"), http.StatusUnauthorized)
      return
    }
    userID, err := r.Services.Auth.GetUserIDFromJWTToken(accessCookie.Value)
    userRoleViews, err := r.Services.Database.GetUserRoleViewsByUserID(userID)
    if err != nil {
      http.Error(res, fmt.Sprintf("Error fetching user role from db: %s", err.Error()), http.StatusInternalServerError)
      return
    }
    hasRoleAdmin := r.Services.Auth.HasRoleAdmin(userRoleViews)
    if !hasRoleAdmin {
      http.Error(res, fmt.Sprintf("User not authorized to POST to api/tag"), http.StatusUnauthorized)
      return
    }

    switch head {
    case "create":
      r.handleCreate(res, req)
    case "update":
      r.handleUpdate(res, req)
    case "delete":
      r.handleDelete(res, req)
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


func (r *TagRouter) handleGetAll(res http.ResponseWriter, req *http.Request) {
  tags, err := r.Services.Database.GetAllTags()
  if err != nil {
    http.Error(res, fmt.Sprintf("Error getting all tags from DB: %s", err.Error()), http.StatusInternalServerError)
    return
  }

  response := &Response{
    Success:  true,
    Error:    nil,
    Data:     TagGetAllResponseData{
      Tags: tags,
    },
  }

  res.Header().Set("Content-Type", "application/json")
  json.NewEncoder(res).Encode(response)
}


func (r *TagRouter) handleCreate(res http.ResponseWriter, req *http.Request) {
  decoder := json.NewDecoder(req.Body)
  tagCreate := new(db.TagCreate)

  err := decoder.Decode(tagCreate)
  if err != nil {
    http.Error(res, fmt.Sprintf("Invalid JSON request: %s", err.Error()), http.StatusBadRequest)
    return
  }

  tagID, err := r.Services.Database.CreateTag(tagCreate)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error creating new tag in database: %s", err.Error()), http.StatusInternalServerError)
    return
  }

  tag, err := r.Services.Database.GetTagByTagID(tagID)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error getting new tag in database: %s", err.Error()), http.StatusInternalServerError)
    return
  }

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
  tagUpdate := new(db.TagUpdate)

  err := decoder.Decode(tagUpdate)
  if err != nil {
    http.Error(res, fmt.Sprintf("Invalid JSON request: %s", err.Error()), http.StatusBadRequest)
    return
  }

  tagID, err := r.Services.Database.UpdateTag(tagUpdate)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error update tag in database: %s", err.Error()), http.StatusInternalServerError)
    return
  }

  tag, err := r.Services.Database.GetTagByTagID(tagID)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error getting updated tag in database: %s", err.Error()), http.StatusInternalServerError)
    return
  }

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


func (r *TagRouter) handleDelete(res http.ResponseWriter, req *http.Request) {
  decoder := json.NewDecoder(req.Body)
  tagDelete := new(db.TagDelete)

  err := decoder.Decode(tagDelete)
  if err != nil {
    http.Error(res, fmt.Sprintf("Invalid JSON request: %s", err.Error()), http.StatusBadRequest)
    return
  }

  _, err = r.Services.Database.DeleteTagByTagID(tagDelete.TagID)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error deleting tag in database: %s", err.Error()), http.StatusInternalServerError)
    return
  }

  response := &Response{
    Success:  true,
    Error:    nil,
  }

  res.Header().Set("Content-Type", "application/json")
  json.NewEncoder(res).Encode(response)
}
