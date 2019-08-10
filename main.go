package main

import (
  "log"
  "net/http"
  "os"

  "github.com/cakebin/smush/server/routes"
)


func main() {
  port := os.Getenv("PORT")
  if port == "" {
    log.Fatal("$PORT must be set")
  }

  router := routes.NewRouter()

  log.Printf("Listening on port %s", port) 
  http.ListenAndServe(":" + port, router)
}
