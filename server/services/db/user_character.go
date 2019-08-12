package db

import (
  "database/sql"
)


/*---------------------------------
          Data Structures
----------------------------------*/

// UserCharacter describes the required and optional data
// needed to create a new "saved character" in our user_characters table
type UserCharacter struct {
  UserCharacterID  int            `json:"userCharacterId"`
  UserID           int            `json:"userId"`
  CharacterID      int            `json:"characterId"`
  CharacterGsp     sql.NullInt64  `json:"characterGsp"`
}


// UserCharacterCreate describes the data needed 
// to create a given "saved character" in our db
type UserCharacterCreate struct {
  UserID           int            `json:"userId"`
  CharacterID      int            `json:"characterId"`
  CharacterGsp     sql.NullInt64  `json:"characterGsp"`
}


// UserCharacterUpdate describes the data needed 
// to update a given "saved character" in our db
type UserCharacterUpdate struct {
  UserCharacterID  int            `json:"userCharacterId"`
  UserID           sql.NullInt64  `json:"userId"`
  CharacterID      sql.NullInt64  `json:"characterId"`
  CharacterGsp     sql.NullInt64  `json:"characterGsp"`
}


/*---------------------------------
            Interface
----------------------------------*/

// UserCharacterManager describes all of the methods used
// to interact with the user_characters table in our database
type UserCharacterManager interface {
  GetUserCharactersByUserID(userID int) ([]*UserCharacter, error)

  CreateUserCharacter(userCharacterCreate *UserCharacterCreate) (int, error)
  UpdateUserCharacter(userCharacterUpdate *UserCharacterUpdate) (int, error)
  DeleteUserCharacterByID(userCharacterID int) (int, error)
}


/*---------------------------------
       Method Implementations
----------------------------------*/

// GetUserCharactersByUserID gets all of the "saved characters" for a given userID
func (db *DB) GetUserCharactersByUserID(userID int) ([]*UserCharacter, error) {
  sqlStatement := `
    SELECT
      user_character_id,
      user_id,
      character_id,
      character_gsp
    FROM
      user_characters
    WHERE
      user_id = $1
  `
  rows, err := db.Query(sqlStatement, userID)
  if err != nil {
    return nil, err
  }
  defer rows.Close()

  userCharacters := make([]*UserCharacter, 0)
  for rows.Next() {
    userCharacter := new(UserCharacter)
    err := rows.Scan(
      &userCharacter.UserCharacterID,
      &userCharacter.UserID,
      &userCharacter.CharacterID,
      &userCharacter.CharacterGsp,
    )

    if err != nil {
      return nil, err
    }

    userCharacters = append(userCharacters, userCharacter)
  }

  err = rows.Err()
  if err != nil {
    return nil, err
  }

  return userCharacters, nil
}


// CreateUserCharacter adds a new entry to the user_characters table
func (db *DB) CreateUserCharacter(userCharacterCreate *UserCharacterCreate) (int, error) {
  var userCharID int
  sqlStatement := `
    INSERT INTO user_characters
      (user_id, character_id, character_gsp)
    VALUES
      ($1, $2, $3, $4)
    RETURNING
      user_character_id
  `
  row := db.QueryRow(
    sqlStatement,
    userCharacterCreate.UserID,
    userCharacterCreate.CharacterID,
    userCharacterCreate.CharacterGsp,
  )

  err := row.Scan(&userCharID)
  if err != nil {
    return 0, err
  }

  return userCharID, nil
}


// UpdateUserCharacter updates an existing entry in the user_characters table
func (db *DB) UpdateUserCharacter(userCharacterUpdate *UserCharacterUpdate) (int, error) {
  var userCharID int
  sqlStatement := `
    UPDATE
      user_characters
    SET
      user_id = $1,
      character_id = $2,
      character_gsp = $3
    WHERE
      user_character_id = $4
    RETURNING
      user_character_id
  `
  row := db.QueryRow(
    sqlStatement,
    userCharacterUpdate.UserID,
    userCharacterUpdate.CharacterID,
    userCharacterUpdate.CharacterGsp,
    userCharacterUpdate.UserCharacterID,
  )

  err := row.Scan(&userCharID)
  if err != nil {
    return 0, err
  }

  return userCharID, nil
}


// DeleteUserCharacterByID removes an existing entry in the user_characters table
func (db *DB) DeleteUserCharacterByID(userCharacterID int) (int, error) {
  var deletedUserCharID int
  sqlStatement := `
    DELETE FROM
      user_characters
    WHERE
      user_character_id = $1
    RETURNING
      user_character_id
  `
  row := db.QueryRow(sqlStatement, userCharacterID)

  err := row.Scan(&deletedUserCharID)
  if err != nil {
    return 0, err
  }

  return deletedUserCharID, nil
}
