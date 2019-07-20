package api

import (
  "net/http"

  "github.com/cakebin/smush/server/env"
  "github.com/cakebin/smush/server/routers/api/match"
  "github.com/cakebin/smush/server/util/routing"
)


// Router is responsible for serving "/api"
// or delegating to the appropriate sub api-router
type Router struct {
  SysUtils *env.SysUtils
  MatchRouter *match.Router
}


func (r *Router) ServeHTTP(res http.ResponseWriter, req *http.Request) {
  var head string
  head, req.URL.Path = routing.ShiftPath(req.URL.Path)

  switch head {
  case "match":
    r.MatchRouter.ServeHTTP(res, req)
  default:
    http.Error(res, "404 Not Found", http.StatusNotFound)
  }
}


// NewRouter makes a new api router and sets up its children
// routers with access to the "SysUtils" environment object
func NewRouter(sysUtils *env.SysUtils) *Router {
  router := new(Router)
  router.SysUtils = sysUtils
  router.MatchRouter = match.NewRouter(sysUtils)
  return router
}
