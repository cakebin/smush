package env

import (
  "github.com/cakebin/smush/server/auth"
  "github.com/cakebin/smush/server/db"
)


// SysUtils holds application level references
// to be used by the rest of our Routers
type SysUtils struct {
  Database db.Datastore
  Authenticator auth.Authenticator
}
