package app

import (
  "net/http"

  "github.com/cakebin/smush/server/util/routing"
)


// Router is the router responsible for serving "/"
// or delegating to the appropriate sub router
type Router struct {}

func (r *Router) ServeHTTP(res http.ResponseWriter, req *http.Request) {
  var head string
  head, req.URL.Path = routing.ShiftPath(req.URL.Path)

  switch head {
  case "", "static":
    http.FileServer((http.Dir("dist/static"))).ServeHTTP(res, req)
  default:
    http.Error(res, "404 Not Found", http.StatusNotFound)
  }
}
