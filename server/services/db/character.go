package db


/*---------------------------------
          Data Structures
----------------------------------*/

// Character describes the required and optional data
// needed to create a new character in our characters table
type Character struct {
  CharacterID        *int     `json:"characterId,omitempty"`
  CharacterName      string   `json:"characterName"`
  CharacterStockImg  string   `json:"characterStockImg"`
}


/*---------------------------------
            Interface
----------------------------------*/

// CharacterManager describes all of the methods used
// to interact with the characters table in our database
type CharacterManager interface {
  GetAllCharacters() ([]*Character, error)

  CreateCharacter(character Character) (int, error)
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
func (db *DB) CreateCharacter(character Character) (int, error) {
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
    character.CharacterName,
    character.CharacterStockImg,
  )
  err := row.Scan(&characterID)

  if err != nil {
    return 0, nil
  }

  return characterID, nil
}
