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
  CharacterID         int             `json:"characterId"`
  CharacterName       string          `json:"characterName"`
  CharacterStockImg   sql.NullString  `json:"characterStockImg"`
  CharacterImg        sql.NullString  `json:"characterImg"`
  CharacterArchetype  sql.NullString  `json:"characterArchetype"`
}


// CharacterUpdate describes the data needed 
// to update a given character in our db
type CharacterUpdate struct {
  CharacterID         int             `json:"characterId"`
  CharacterName       sql.NullString  `json:"characterName"`
  CharacterStockImg   sql.NullString  `json:"characterStockImg"`
  CharacterImg        sql.NullString  `json:"characterImg"`
  CharacterArchetype  sql.NullString  `json:"characterArchetype"`
}


// CharacterCreate describes the data needed 
// to create a given character in our db
type CharacterCreate struct {
  CharacterName       string          `json:"characterName"`
  CharacterStockImg   sql.NullString  `json:"characterStockImg"`
  CharacterImg        sql.NullString  `json:"characterImg"`
  CharacterArchetype  sql.NullString  `json:"characterArchetype"`
}


/*---------------------------------
            Interface
----------------------------------*/

// CharacterManager describes all of the methods used
// to interact with the characters table in our database
type CharacterManager interface {
  GetAllCharacters() ([]*Character, error)

  CreateCharacter(character CharacterCreate) (*Character, error)
  UpdateCharacter(update CharacterUpdate) (*Character, error)
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
      character_stock_img,
      character_img,
      character_archetype
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
      &character.CharacterImg,
      &character.CharacterArchetype,
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
func (db *DB) CreateCharacter(characterCreate CharacterCreate) (*Character, error) {
  sqlStatement := `
    INSERT INTO characters
      (character_name, character_stock_img, character_img, character_archetype)
    VALUES
      ($1, $2, $3, $4)
    RETURNING
      character_id,
      character_name,
      character_stock_img,
      character_img,
      character_archetype
  `
  row := db.QueryRow(
    sqlStatement,
    characterCreate.CharacterName,
    characterCreate.CharacterStockImg.String,
    characterCreate.CharacterImg.String,
    characterCreate.CharacterArchetype.String,
  )
  character := new(Character)
  err := row.Scan(
    &character.CharacterID,
    &character.CharacterName,
    &character.CharacterStockImg,
    &character.CharacterImg,
    &character.CharacterArchetype,
  )

  if err != nil {
    return nil, nil
  }

  return character, nil
}


// UpdateCharacter updates an existing entry in the characters table in our database
func (db *DB) UpdateCharacter(characterUpdate CharacterUpdate) (*Character, error) {
  sqlStatement := `
    UPDATE
      characters
    SET
      character_name = $1,
      character_stock_img = $2,
      character_img = $3,
      character_archetype = $4
    WHERE
      character_id = $5
    RETURNING
      character_id,
      character_name,
      character_stock_img,
      character_img,
      character_archetype
  `
  row := db.QueryRow(
    sqlStatement,
    characterUpdate.CharacterName.String,
    characterUpdate.CharacterStockImg.String,
    characterUpdate.CharacterImg.String,
    characterUpdate.CharacterArchetype.String,
    characterUpdate.CharacterID,
  )
  character := new(Character)
  err := row.Scan(
    &character.CharacterID,
    &character.CharacterName,
    &character.CharacterStockImg,
    &character.CharacterImg,
    &character.CharacterArchetype,
  )

  if err != nil {
    return nil, err
  }

  return character, nil
}
