package routes

import (
  "net/http"
)


/*---------------------------------
             Responses
----------------------------------*/

// Response is the data 
type Response struct {
  Success  bool         `json:"success"`
  Error    error        `json:"error"`
  Data     interface{}  `json:"data"`
}


/*---------------------------------
             Router
----------------------------------*/

// APIRouter is responsible for serving "/api"
// or delegating to the appropriate sub api-router
type APIRouter struct {
  Services         *Services
  AuthRouter       *AuthRouter
  MatchRouter      *MatchRouter
  UserRouter       *UserRouter
  CharacterRouter  *CharacterRouter
}



func (r *APIRouter) ServeHTTP(res http.ResponseWriter, req *http.Request) {
  var head string
  head, req.URL.Path = ShiftPath(req.URL.Path)

  if head == "auth" {
    r.AuthRouter.ServeHTTP(res, req)
    return
  }

  /*-----------------------------------------
         API Route Authentication
  ------------------------------------------*/
  // Check the access token
  accessCookie, err := req.Cookie("smush-access-token")

  // If the access token exists, check the token value
  if err == nil {
    // If this is run, and there is an error (token is invalid),
    // err will not be nil in the next if block
    _, err = r.Services.Auth.CheckJWTToken(accessCookie.Value)
  }

  // Access token is invalid, or there wasn't an access cookie in the first place
  // I do not like reusing the same err here... the cascading logic is confusing. But it'll do for now.
  if err != nil {
    http.Error(res, "Session expired. Please log in again", http.StatusUnauthorized)
    return
  }

  // Once we're authorized, allow api requests
  switch head {
  case "match":
    r.MatchRouter.ServeHTTP(res, req)
  case "user":
    r.UserRouter.ServeHTTP(res, req)
  case "character":
    r.CharacterRouter.ServeHTTP(res, req)
  default:
    http.Error(res, "404 Not Found", http.StatusNotFound)
  }
}


// NewAPIRouter makes a new api router and sets up its children
// routers with access to our router services
func NewAPIRouter(routerServices *Services) *APIRouter {
  router := new(APIRouter)

  router.Services = routerServices
  router.AuthRouter = NewAuthRouter(routerServices)
  router.MatchRouter = NewMatchRouter(routerServices)
  router.UserRouter = NewUserRouter(routerServices)
  router.CharacterRouter = NewCharacterRouter(routerServices)

  return router
}
