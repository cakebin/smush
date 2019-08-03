package main

import (
  "log"
  "net/http"
  "os"

  "github.com/cakebin/smush/server/auth"
  "github.com/cakebin/smush/server/env"
  "github.com/cakebin/smush/server/routers/app"
  "github.com/cakebin/smush/server/services/database"
)


func main() {
  port := os.Getenv("PORT")
  if port == "" {
    log.Fatal("$PORT must be set")
  }

  dbURL := os.Getenv(("DATABASE_URL"))
  db, err := database.New(dbURL)
  if err != nil {
    log.Fatalf("Error opening database: %q", err)
  }

  authenticator := auth.NewAuthenticator()

  sysUtils := &env.SysUtils{
    Database: db,
    Authenticator: authenticator,
  }
  router := app.NewRouter(sysUtils)

  log.Printf("Listening on port %s", port) 
  http.ListenAndServe(":" + port, router)
}
