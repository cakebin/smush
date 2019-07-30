package db

import (
  "time"
)

// User represents a recorded user for our Smash Ultimate Online tracker app
type User struct {
  UserID               *int      `json:"userId,omitempty"`
  UserName             string    `json:"userName"`
  EmailAddress         string    `json:"emailAddress"`
  DefaultCharacterName *string   `json:"defaultCharacterName,omitempty"`
  DefaultCharacterGsp  *int      `json:"defaultCharacterGsp,omitempty"`
  Created              time.Time `json:"created"`
  Password             *string   `json:"password,omitempty"`
  HashedPassword       *string   `json:"hashedPassword,omitempty"`
  RefreshToken         *string   `json:"refreshToken,omitempty"`
}

// GetUserByID gets a specific user by its id in the Users table
func (db *DB) GetUserByID(id int) (*User, error) {
  sqlStatement := `
  SELECT
    user_id,
    user_name,
    email_address,
    default_character_name,
    default_character_gsp,
    created
  FROM
    users
  WHERE
    user_id = $1
  `
  row := db.QueryRow(sqlStatement, id)
  user := new(User)
  err := row.Scan(
    &user.UserID,
    &user.UserName,
    &user.EmailAddress,
    &user.DefaultCharacterName,
    &user.DefaultCharacterGsp,
    &user.Created,
  )

  if err != nil {
    return nil, err
  }

  return user, nil
}


// GetUserByEmail gets a user's auth related information
// by their email; used for user authentication
func (db *DB) GetUserByEmail(email string) (*User, error) {
  sqlStatement := `
  SELECT
    user_id,
    user_name,
    email_address,
    default_character_name,
    default_character_gsp,
    created,
    hashed_password,
    refresh_token
  FROM
   users
  WHERE
    email_address = $1`
  row := db.QueryRow(sqlStatement, email)
  user := new(User)
  err := row.Scan(
    &user.UserID,
    &user.UserName,
    &user.EmailAddress,
    &user.DefaultCharacterName,
    &user.DefaultCharacterGsp,
    &user.Created,
    &user.HashedPassword,
    &user.RefreshToken,
  )

  if err != nil {
    return nil, err
  }

  return user, nil
}

// UpdateUserRefreshTokenByID updates a user's stored jwt
// refresh token by their email; used for user authentication
func (db *DB) UpdateUserRefreshTokenByID(token string, id int) (bool, error) {
  sqlStatement := `
  UPDATE users SET
    refresh_token = $1
  WHERE
    user_id = $2`
  _, err := db.Exec(
    sqlStatement,
    token,
    id,
  )

  if err != nil {
    return false, err
  }

  return true, nil
}

// UpdateUser updates an existing entry in the Users table
func (db *DB) UpdateUser(user User) (bool, error) {
  sqlStatement := `
  UPDATE users SET 
    user_name = $1,
    default_character_name = $2,
    default_character_gsp = $3
  WHERE user_id = $4`
  _, err := db.Exec(
    sqlStatement,
    user.UserName,
    user.DefaultCharacterName,
    user.DefaultCharacterGsp,
    user.UserID,
  )

  if err != nil {
    return false, err
  }

  return true, nil
}

// CreateUser adds a new entry to the Users table in our database
func (db *DB) CreateUser(user User) (bool, error) {
  sqlStatement := `
  INSERT INTO users (
    user_name,
    email_address,
    hashed_password,
    refresh_token
  )
  VALUES ($1, $2, $3, $4)`
  _, err := db.Exec(
    sqlStatement,
    user.UserName,
    user.EmailAddress,
    user.HashedPassword,
    user.RefreshToken,
  )

  if err != nil {
    return false, err
  }

  return true, nil
}
