package db

import (
  "database/sql"
)


/*---------------------------------
          Data Structures
----------------------------------*/

// UserProfileUpdate describes the data needed 
// to update a given user's profile information
type UserProfileUpdate struct {
  UserID               int            `json:"userId"`
  UserName             string         `json:"userName"`
  DefaultCharacterID   sql.NullInt64  `json:"defaultCharacterId"`
  DefaultCharacterGsp  sql.NullInt64  `json:"defaultCharacterGsp"`
}


// UserRefreshUpdate describes the data
// needed to update a given users refresh token
type UserRefreshUpdate struct {
  UserID        int     `json:"userId"`
  RefreshToken  string  `json:"refreshToken"`
}


// User describes the required and optional data
// needed to create a new user in our users table
type User struct {
  UserID          *int    `json:"userId,omitempty"`
  UserName        string  `json:"userName"`
  EmailAddress    string  `json:"emailAddress"`
  Password        string  `json:"password"`
  HashedPassword  string  `json:"hashedPassword"`
  RefreshToken    string  `json:"refreshToken"`
}


/*---------------------------------
            Interface
----------------------------------*/

// UserManager describes all of the methods used
// to interact with the users table in our database
type UserManager interface {
  GetUserIDByEmail(email string) (int, error)

  UpdateUserProfile(profileUpdate UserProfileUpdate) (int, error)
  UpdateUserRefreshToken(refreshUpdate UserRefreshUpdate) (int, error)

  CreateUser(user User) (int, error)
}


/*---------------------------------
       Method Implementations
----------------------------------*/

// GetUserIDByEmail gets a specific user's id from the users table by email
func (db *DB) GetUserIDByEmail(email string) (int, error) {
  var userID int
  sqlStatement := `
    SELECT
      user_id
    FROM
      users
    WHERE
      email_address = $1
  `
  row := db.QueryRow(sqlStatement, email)
  err := row.Scan(&userID)

  if err != nil {
    return 0, err
  }

  return userID, nil
}


// CreateUser adds a new entry to the users table in our database
func (db *DB) CreateUser(user User) (int, error) {
  var userID int
  sqlStatement := `
    INSERT INTO users
      (user_name, email_address, hashed_password)
    VALUES 
      ($1, $2, $3)
    RETURNING
      user_id
  `
  row := db.QueryRow(
    sqlStatement,
    user.UserName,
    user.EmailAddress,
    user.HashedPassword,
  )
  err := row.Scan(&userID)

  if err != nil {
    return 0, err
  }

  return userID, nil
}


// UpdateUserProfile updates an entry in the users table with the given data
func (db *DB) UpdateUserProfile(profileUpdate UserProfileUpdate) (int, error) {
  var userID int
  sqlStatement := `
    UPDATE
      users
    SET
      user_name = $1,
      default_character_id = $2,
      default_character_gsp = $3
    WHERE
      user_id = $4
    RETURNING
      user_id
  `
  row := db.QueryRow(
    sqlStatement,
    profileUpdate.UserName,
    profileUpdate.DefaultCharacterID,
    profileUpdate.DefaultCharacterGsp,
    profileUpdate.UserID,
  )
  err := row.Scan(&userID)
  
  if err != nil {
    return 0, err
  }

  return userID, nil

}


// UpdateUserRefreshToken updates an a user's refresh token; used for auth
func (db *DB) UpdateUserRefreshToken(refreshUpdate UserRefreshUpdate) (int, error) {
  var userID int
  sqlStatement := `
    UPDATE
      users
    SET
      refresh_token = $1
    WHERE
      user_id = $2
    RETURNING
      user_id
  `
  row := db.QueryRow(
    sqlStatement,
    refreshUpdate.RefreshToken,
    refreshUpdate.UserID,
  )
  err := row.Scan(&userID)

  if err != nil {
    return 0, err
  }

  return userID, nil
}
