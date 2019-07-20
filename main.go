package main

import (
  "log"
  "net/http"
  "os"

  "github.com/cakebin/smush/server/db"
  "github.com/cakebin/smush/server/env"
  "github.com/cakebin/smush/server/routers/app"
)


func main() {
  port := os.Getenv("PORT")
  if port == "" {
    log.Fatal("$PORT must be set")
  }

  dbURL := os.Getenv(("DATABASE_URL"))
  database, err := db.NewDB(dbURL)
  if err != nil {
    log.Fatalf("Error opening database: %q", err)
  }

  sysUtils := &env.SysUtils{Database: database}
  router := app.NewRouter(sysUtils)

  log.Printf("Listening on port %s", port) 
  http.ListenAndServe(":" + port, router)
}
