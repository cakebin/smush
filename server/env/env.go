package env

import (
  "github.com/cakebin/smush/server/auth"
  "github.com/cakebin/smush/server/services/database"
)


// SysUtils holds application level references
// to be used by the rest of our Routers
type SysUtils struct {
  Database       database.DatabaseManager
  Authenticator  auth.Authenticator
}
