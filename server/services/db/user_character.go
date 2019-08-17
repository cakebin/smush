package db


/*---------------------------------
          Data Structures
----------------------------------*/

// UserCharacter describes the required and optional data
// needed to create a new "saved character" in our user_characters table
type UserCharacter struct {
  UserCharacterID  int64            `json:"userCharacterId"`
  UserID           int64            `json:"userId"`
  CharacterID      int64            `json:"characterId"`
  CharacterGsp     NullInt64JSON  `json:"characterGsp"`
  AltCostume       NullInt64JSON  `json:"altCostume"`
}


// UserCharacterCreate describes the data needed 
// to create a given "saved character" in our db
type UserCharacterCreate struct {
  UserID           int64            `json:"userId"`
  CharacterID      int64            `json:"characterId"`
  CharacterGsp     NullInt64JSON  `json:"characterGsp"`
  AltCostume       NullInt64JSON  `json:"altCostume"`
}


// UserCharacterUpdate describes the data needed 
// to update a given "saved character" in our db
type UserCharacterUpdate struct {
  UserCharacterID  int64            `json:"userCharacterId"`
  UserID           int64            `json:"userId"`
  CharacterID      NullInt64JSON  `json:"characterId"`
  CharacterGsp     NullInt64JSON  `json:"characterGsp"`
  AltCostume       NullInt64JSON  `json:"altCostume"`
}


// UserCharacterDelete describes the data needed 
// to delete a given "saved character" in our db
type UserCharacterDelete struct {
  UserID           int64          `json:"userId"`
  UserCharacterID  NullInt64JSON  `json:"userCharacterId"`
}


/*---------------------------------
            Interface
----------------------------------*/

// UserCharacterManager describes all of the methods used
// to int64eract with the user_characters table in our database
type UserCharacterManager interface {
  GetUserCharactersByUserID(userID int64) ([]*UserCharacter, error)

  CreateUserCharacter(userCharacterCreate *UserCharacterCreate) (int64, error)
  UpdateUserCharacter(userCharacterUpdate *UserCharacterUpdate) (int64, error)
  DeleteUserCharacterByID(userCharacterID int64) (int64, error)
}


/*---------------------------------
       Method Implementations
----------------------------------*/

// GetUserCharactersByUserID gets all of the "saved characters" for a given userID
func (db *DB) GetUserCharactersByUserID(userID int64) ([]*UserCharacter, error) {
  sqlStatement := `
    SELECT
      user_character_id,
      user_id,
      character_id,
      character_gsp,
      alt_costume
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
      &userCharacter.AltCostume,
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
func (db *DB) CreateUserCharacter(userCharacterCreate *UserCharacterCreate) (int64, error) {
  var userCharID int64
  sqlStatement := `
    INSERT INTO user_characters
      (user_id, character_id, character_gsp, alt_costume)
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
    userCharacterCreate.AltCostume,
  )

  err := row.Scan(&userCharID)
  if err != nil {
    return 0, err
  }

  return userCharID, nil
}


// UpdateUserCharacter updates an existing entry in the user_characters table
func (db *DB) UpdateUserCharacter(userCharacterUpdate *UserCharacterUpdate) (int64, error) {
  var userCharID int64
  sqlStatement := `
    UPDATE
      user_characters
    SET
      user_id = $1,
      character_id = $2,
      character_gsp = $3,
      alt_costume = $4
    WHERE
      user_character_id = $5
    RETURNING
      user_character_id
  `
  row := db.QueryRow(
    sqlStatement,
    userCharacterUpdate.UserID,
    userCharacterUpdate.CharacterID,
    userCharacterUpdate.CharacterGsp,
    userCharacterUpdate.AltCostume,
    userCharacterUpdate.UserCharacterID,
  )

  err := row.Scan(&userCharID)
  if err != nil {
    return 0, err
  }

  return userCharID, nil
}


// DeleteUserCharacterByID removes an existing entry in the user_characters table
func (db *DB) DeleteUserCharacterByID(userCharacterID int64) (int64, error) {
  var deletedUserCharID int64
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
