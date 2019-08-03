package user


import (
  "encoding/json"
  "fmt"
  "net/http"
  "strconv"

  "github.com/cakebin/smush/server/api"
  "github.com/cakebin/smush/server/env"
  "github.com/cakebin/smush/server/services/db"
  "github.com/cakebin/smush/server/util/routing"
)


// Router is responsible for serving "/api/user"
// Basically, connecting to our Postgres DB for all
// of the CRUD operations for our "User" models
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
    }
  // Unsupported Method Response
  default:
    http.Error(res, fmt.Sprintf("Unsupported Method type %s", req.Method), http.StatusBadRequest)
  }
}


func (r *Router) handleGetByID(res http.ResponseWriter, req *http.Request) {
  var head string
  head, req.URL.Path = routing.ShiftPath(req.URL.Path)

  id, err := strconv.Atoi(head)
  if err != nil {
    http.Error(res, fmt.Sprintf("Invalid user id: %s", head), http.StatusBadRequest)
    return
  }

  userProfileView, err := r.SysUtils.Database.GetUserProfileViewByID(id)
  if err != nil {
    http.Error(res, fmt.Sprintf("Error getting user with id %q: %s", id, err.Error()), http.StatusInternalServerError)
    return
  }

  res.Header().Set("Content-Type", "application/json")
  json.NewEncoder(res).Encode(userProfileView)
}


func (r *Router) handleCreate(res http.ResponseWriter, req *http.Request) {
  decoder := json.NewDecoder(req.Body)
  var user db.User

  err := decoder.Decode(&user)
  if err != nil {
    http.Error(res, fmt.Sprintf("Invalid JSON request: %s", err.Error()), http.StatusBadRequest)
    return
  }
  
  hashedPassword, err := r.SysUtils.Authenticator.HashPassword(user.Password)
  user.HashedPassword = hashedPassword
  userID, err := r.SysUtils.Database.CreateUser(user)

  response := &api.Response{
    Success: true,
    Error: err,
    Data: userID,
  }

  res.Header().Set("Content-Type", "application/json")
  json.NewEncoder(res).Encode(response)
}

func (r *Router) handleUpdate(res http.ResponseWriter, req *http.Request) {
  decoder := json.NewDecoder(req.Body)
  var userProfileUpdate db.UserProfileUpdate

  err := decoder.Decode(&userProfileUpdate)
  if err != nil {
    http.Error(res, fmt.Sprintf("Invalid JSON request: %s", err.Error()), http.StatusBadRequest)
    return
  }

  userID, err := r.SysUtils.Database.UpdateUserProfile(userProfileUpdate)

  response := &api.Response{
    Success: true,
    Error: err,
    Data: userID,
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