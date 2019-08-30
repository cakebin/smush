package auth

import (
  "github.com/cakebin/smush/server/services/db"
)


const adminRoleID int64 = 1


/*---------------------------------
            Interface
----------------------------------*/

// RoleManager describes all of the methods used
// for handling the permissions/roles auth layer
type RoleManager interface {
  HasRoleAdmin(userRoleViews []*db.UserRoleView) bool
}


/*---------------------------------
       Method Implementations
----------------------------------*/

// HasRoleAdmin checks to see whether or not user's
// fetched user roles has the "admin" role
func (a *Auth) HasRoleAdmin(userRoleViews []*db.UserRoleView) bool {
  for _, userRoleView := range userRoleViews {
    if userRoleView.RoleID == adminRoleID {
      return true
    }
  }

  return false
}
