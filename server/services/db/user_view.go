package db

import (
  "database/sql"
  "time"
)

/*---------------------------------
          Data Structures
----------------------------------*/

// UserProfileView describes all of the required
// and optional data needed for a user's public
type UserProfileView struct {
  // Data from Users
  UserID                int             `json:"userId"`
  UserName              string          `json:"userName"`
  EmailAddress          string          `json:"emailAddress"`
  Created               time.Time       `json:"created"`
  DefaultCharacterGsp   sql.NullInt64   `json:"defaultCharacterGsp"`

  // Data from characters
  DefaultCharacterID    sql.NullInt64   `json:"defaultCharacterId"`
  DefaultCharacterName  sql.NullString  `json:"defaultCharacterName"`
}

// UserCredentialsView describes all of the data
// needed for a user's authentication credentials
type UserCredentialsView struct {
  EmailAddress    string  `json:"email"`
  UserID          int     `json:"userId"`
  UserName        string  `json:"userName"`
  HashedPassword  string  `json:"hashedPassword"`
}

/*---------------------------------
            Interface
----------------------------------*/

// UserViewManager describes all of the methods used to interact with
// user views in our database (data joined between match, character, user, etc)
type UserViewManager interface {
  GetUserProfileViewByID(userID int) (*UserProfileView, error)
  GetUserCredentialsViewByEmail(email string) (*UserCredentialsView, error)
}

/*---------------------------------
       Method Implementations
----------------------------------*/

// GetUserProfileViewByID gets all of the data needed to display
// a user's profile, which includes joined data from the characters table
func (db *DB) GetUserProfileViewByID(userID int) (*UserProfileView, error) {
  sqlStatement := `
    SELECT
      users.user_id                     AS user_id,
      users.user_name                   AS user_name,
      users.email_address               AS email_address,
      users.created                     AS created,
      users.default_character_gsp       AS default_character_gsp,
      characters.character_id           AS default_character_id,
      characters.character_name         AS default_character_name
    FROM
      users
    LEFT JOIN characters ON characters.character_id = users.default_character_id
    WHERE
      user_id = $1
  `
  row := db.QueryRow(sqlStatement, userID)
  userProfileView := new(UserProfileView)
  err := row.Scan(
    &userProfileView.UserID,
    &userProfileView.UserName,
    &userProfileView.EmailAddress,
    &userProfileView.Created,
    &userProfileView.DefaultCharacterGsp,
    &userProfileView.DefaultCharacterID,
    &userProfileView.DefaultCharacterName,
  )

  if err != nil {
    return nil, err
  }

  return userProfileView, nil
}

// GetUserCredentialsViewByEmail gets a user's auth related
// information by their email; used for user authentication
func (db *DB) GetUserCredentialsViewByEmail(email string) (*UserCredentialsView, error) {
  sqlStatement := `
    SELECT
      user_id,
      user_name,
      email_address,
      hashed_password
    FROM
      users
    WHERE
      email_address = $1
  `
  row := db.QueryRow(sqlStatement, email)
  userCredentialsView := new(UserCredentialsView)
  err := row.Scan(
    &userCredentialsView.UserID,
    &userCredentialsView.UserName,
    &userCredentialsView.EmailAddress,
    &userCredentialsView.HashedPassword,
  )

  if err != nil {
    return nil, err
  }

  return userCredentialsView, nil
}
