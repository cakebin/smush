package env

import (
  "github.com/cakebin/smush/server/services/auth"
  "github.com/cakebin/smush/server/services/db"
)


// SysUtils holds application level references
// to be used by the rest of our Routers
type SysUtils struct {
  Database       db.DatabaseManager
  Authenticator  auth.Authenticator
}
