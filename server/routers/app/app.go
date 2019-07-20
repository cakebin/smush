package app

import (
  "net/http"

  "github.com/cakebin/smush/server/env"
  "github.com/cakebin/smush/server/routers/api"
  "github.com/cakebin/smush/server/util/routing"
)


// Router is the router responsible for serving "/"
// or delegating to the appropriate sub router
type Router struct {
  SysUtils *env.SysUtils
  APIRouter *api.Router
}


func (r *Router) ServeHTTP(res http.ResponseWriter, req *http.Request) {
  var head string
  head, req.URL.Path = routing.ShiftPath(req.URL.Path)

  switch head {
  case "", "static":
    http.FileServer((http.Dir("dist/static"))).ServeHTTP(res, req)
  case "api":
    r.APIRouter.ServeHTTP(res, req)
  default:
    http.Error(res, "404 Not Found", http.StatusNotFound)
  }
}


// NewRouter makes a new app router and sets up its children
// routers with access to the "SysUtils" environment object
func NewRouter(sysUtils *env.SysUtils) *Router {
  router := new(Router)
  router.SysUtils = sysUtils
  router.APIRouter = api.NewRouter(sysUtils)
  return router
}
