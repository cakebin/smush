package db


import (
  "time"
)


// User represents a recorded user for our Smash Ultimate Online tracker app
type User struct {
  UserID                  *int          `json:"userId,omitempty"`
  UserName                string        `json:"userName"`
  EmailAddress            string        `json:"emailAddress"`
  DefaultCharacterName    *string       `json:"defaultCharacterName,omitempty"`
  DefaultCharacterGsp     *int          `json:"defaultCharacterGsp,omitempty"`
  Created                 time.Time     `json:"created"`
}


// UserResponse represents an interaction with our database
// regarding operations related to the Users table
type UserResponse struct {
  Success  bool   `json:"success"`
  Error    error  `json:"error"`
}


// GetUserByID gets a specific user by its id in the Users table
func (db *DB) GetUserByID(id int) (*User, error) {
  row := db.QueryRow(`SELECT * FROM users WHERE id = $1`, id)
  user := new(User)
  err := row.Scan(
    &user.UserID,
    &user.UserName,
    &user.EmailAddress,
    &user.DefaultCharacterName,
    &user.DefaultCharacterGsp,
  )

  if err != nil {
    return nil, err
  }

  return user, nil
}


// CreateUser adds a new entry to the Users table in our database
func (db *DB) CreateUser(user User) (bool, error) {
  sqlStatement := `
  INSERT INTO users (
    user_name,
    email_address,
    default_character_name,
    default_character_gsp,
  )
  VALUES ($1, $2, $3, $4)`
  _, err := db.Exec(
    sqlStatement,
    user.UserName,
    user.EmailAddress,
    user.DefaultCharacterName,
    user.DefaultCharacterGsp,
  )

  if err != nil {
    return false, err
  }

  return true,  nil
}
