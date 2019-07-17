package main

import (
  "log"
  "net/http"
  "os"
)

func main() {
  port := os.Getenv("PORT")

  if port == "" {
    log.Fatal("$PORT must be set")
  }

  log.Printf("Listening on port %s", port) 
  http.Handle("/", http.FileServer((http.Dir("dist"))))
  http.ListenAndServe(":" + port, nil)
}
