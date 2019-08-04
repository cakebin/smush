package app

import (
  "net/http"

  "github.com/cakebin/smush/server/env"
  "github.com/cakebin/smush/server/services/api"
  "github.com/cakebin/smush/server/util/routing"
)

// Router is the router responsible for serving "/"
// or delegating to the appropriate sub router
type Router struct {
  SysUtils  *env.SysUtils
  APIRouter *api.Router
}

func (r *Router) ServeHTTP(res http.ResponseWriter, req *http.Request) {
  var head string
  head, req.URL.Path = routing.ShiftPath(req.URL.Path)

  switch head {
  case "static":
    http.FileServer(http.Dir("dist/static")).ServeHTTP(res, req)
  case "api":
    r.APIRouter.ServeHTTP(res, req)
  default:
    // Angular requires returning index.html if you are using routing in your app:
    // https://angular.io/guide/deployment#routed-apps-must-fallback-to-indexhtml
    http.ServeFile(res, req, "dist/static/index.html")
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
