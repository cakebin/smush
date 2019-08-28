package db


/*---------------------------------
            Interface
----------------------------------*/

// UserRoleViewManager describes all of the methods used
// to interact with our user_roles table in our database
type UserRoleViewManager interface {
  GetUserRoleViewsByUserID(userID int64) ([]*UserRoleView, error)
}

/*---------------------------------
          Data Structures
----------------------------------*/

// UserRoleView describes the data relevant to a user's role
type UserRoleView struct {
  UserRoleID  int64   `json:"userRoleId"`
  UserID      int64   `json:"userId"`
  RoleID      int64   `json:"roleId"`
  RoleName    string  `json:"roleName"`
}


/*---------------------------------
       Method Implementations
----------------------------------*/


// GetUserRoleViewsByUserID gets a user's role information
func (db *DB) GetUserRoleViewsByUserID(userID int64) ([]*UserRoleView, error) {
  sqlStatement := `
    SELECT
      user_roles.user_role_id  AS  user_role_id,
      user_roles.user_id       AS  user_id,
      user_roles.role_id       AS  role_id,
      roles.role_name          AS  role_name
    FROM
      user_roles
    LEFT JOIN roles ON roles.role_id = user_roles.role_id
    WHERE
      user_id = $1
  `
  rows, err := db.Query(sqlStatement, userID)
  if err != nil {
    return nil, err
  }
  defer rows.Close()

  userRoleViews := make([]*UserRoleView, 0)
  for rows.Next() {
    userRoleView := new(UserRoleView)
    err := rows.Scan(
      &userRoleView.UserRoleID,
      &userRoleView.UserID,
      &userRoleView.RoleID,
      &userRoleView.RoleName,
    )
    if err != nil {
      return nil, err
    }

    userRoleViews = append(userRoleViews, userRoleView)
  }

  err = rows.Err()
  if err != nil {
    return nil, err
  }


  return userRoleViews, nil
}
