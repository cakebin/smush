package routes

import (
  "net/http"
)


/*---------------------------------
             Router
----------------------------------*/

// StaticRouter is responsible for serving "/static"
type StaticRouter struct {}


func (r *StaticRouter) ServeHTTP(res http.ResponseWriter, req *http.Request) {
  // For now just serve the files; but we can make 
  // optimizations for serving static files here later
  http.FileServer(http.Dir("dist/static")).ServeHTTP(res, req)
}


// NewStaticRouter makes a new api/user router for serving our static assets
func NewStaticRouter() *StaticRouter {
  router := new(StaticRouter)

  return router
}
