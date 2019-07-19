package main

import (
  "log"
  "net/http"
  "os"

  "github.com/cakebin/smush/server/routers/app"
)

func main() {
  port := os.Getenv("PORT")

  if port == "" {
    log.Fatal("$PORT must be set")
  }

  router := &app.Router{}

  log.Printf("Listening on port %s", port) 
  http.ListenAndServe(":" + port, router)
}
