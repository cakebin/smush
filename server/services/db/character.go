package db

import (
  "database/sql"
)


/*---------------------------------
          Data Structures
----------------------------------*/

// Character describes the required and optional data
// needed to create a new character in our characters table
type Character struct {
  CharacterID        int             `json:"characterId"`
  CharacterName      string          `json:"characterName"`
  CharacterStockImg  sql.NullString  `json:"characterStockImg"`
}


// CharacterUpdate describes the data needed 
// to update a given character in our db
type CharacterUpdate struct {
  CharacterID        int             `json:"characterId"`
  CharacterName      sql.NullString  `json:"characterName"`
  CharacterStockImg  sql.NullString  `json:"characterStockImg"`
}


// CharacterCreate describes the data needed 
// to create a given character in our db
type CharacterCreate struct {
  CharacterName      string          `json:"characterName"`
  CharacterStockImg  sql.NullString  `json:"characterStockImg"`
}


/*---------------------------------
            Interface
----------------------------------*/

// CharacterManager describes all of the methods used
// to interact with the characters table in our database
type CharacterManager interface {
  GetAllCharacters() ([]*Character, error)

  CreateCharacter(character CharacterCreate) (int, error)
  UpdateCharacter(update CharacterUpdate) (int, error)
}


/*---------------------------------
       Method Implementations
----------------------------------*/

// GetAllCharacters gets all of the characters we have in our database
func (db *DB) GetAllCharacters() ([]*Character, error) {
  sqlStatement := `
    SELECT
      character_id,
      character_name,
      character_stock_img
    FROM
      characters
  `
  rows, err := db.Query(sqlStatement)
  if err != nil {
    return nil, err
  }
  defer rows.Close()

  characters := make([]*Character, 0)
  for rows.Next() {
    character := new(Character)
    err := rows.Scan(
      &character.CharacterID,
      &character.CharacterName,
      &character.CharacterStockImg,
    )

    if err != nil {
      return nil, err
    }

    characters = append(characters, character)
  }

  err = rows.Err()
  if err != nil {
    return nil, err
  }

  return characters, nil
}


// CreateCharacter adds a new entry to the characters table in our database
func (db *DB) CreateCharacter(characterCreate CharacterCreate) (int, error) {
  var characterID int
  sqlStatement := `
    INSERT INTO characters
      (character_name, character_stock_img)
    VALUES
      ($1, $2)
    RETURNING
      character_id
  `
  row := db.QueryRow(
    sqlStatement,
    characterCreate.CharacterName,
    characterCreate.CharacterStockImg.String,
  )
  err := row.Scan(&characterID)

  if err != nil {
    return 0, nil
  }

  return characterID, nil
}


// UpdateCharacter updates an existing entry in the characters table in our database
func (db *DB) UpdateCharacter(characterUpdate CharacterUpdate) (int, error) {
  var characterID int
  sqlStatement := `
    UPDATE
      characters
    SET
      character_name = $1,
      character_stock_img = $2
    WHERE
      character_id = $3
    RETURNING
      character_id
  `
  row := db.QueryRow(
    sqlStatement,
    characterUpdate.CharacterName.String,
    characterUpdate.CharacterStockImg.String,
    characterUpdate.CharacterID,
  )
  err := row.Scan(&characterID)

  if err != nil {
    return 0, err
  }

  return characterID, nil
}
