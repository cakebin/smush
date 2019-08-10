package routes

import (
  "log"

  "github.com/cakebin/smush/server/services/auth"
  "github.com/cakebin/smush/server/services/db"
)


// Services describes what services are
// available to all routes in our application
type Services struct {
  Database  db.DatabaseManager
  Auth      auth.Authenticator
}


// NewRouterServices initializes all of the services available to our routers
func NewRouterServices() *Services {
  services := new(Services)

  database, err := db.New()
  if err != nil {
    log.Fatalf("Error opening database: %s", err.Error())
  }
  services.Database = database
  services.Auth = auth.New()

  return services
}
