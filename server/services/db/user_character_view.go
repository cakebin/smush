package db


/*---------------------------------
            Interface
----------------------------------*/

// UserCharacterViewManager describes all of the methods used to int64eract with
// "saved characters" views in our database (data joined between users, characters, and user_characters)
type UserCharacterViewManager interface {
  GetUserCharacterViewsByUserID(userID int64) ([]*UserCharacterView, error)
  GetUserCharacterViewByUserCharacterID(userCharID int64) (*UserCharacterView, error)
}


/*---------------------------------
          Data Structures
----------------------------------*/

// UserCharacterView desribes a JOIN between the characters, users, and user_characters tables,
// containing all of the data necessary to show a "saved character" in the front end
type UserCharacterView struct {
  // Data from user_characters
  UserCharacterID  int64          `json:"userCharacterId"`
  CharacterGsp     NullInt64JSON  `json:"characterGsp"`
  AltCostume       NullInt64JSON  `json:"altCostume"`

  // Data from characters
  CharacterID      int64          `json:"characterId"`
  CharacterName    string         `json:"characterName"`

  // Data from users
  UserID           int64          `json:"userId"`
}


/*---------------------------------
       Method Implementations
----------------------------------*/

// GetUserCharacterViewsByUserID gets all of the data needed 
// to display a given user's "saved characters", whicn includes 
// joined data fromt the user_characters and characters table
func (db *DB) GetUserCharacterViewsByUserID(userID int64) ([]*UserCharacterView, error) {
  sqlStatement := `
    SELECT
      user_characters.user_character_id  AS  user_character_id,
      user_characters.character_gsp      AS  character_gsp,
      user_characters.alt_costume        AS  alt_costume,
      characters.character_id            AS  character_id,
      characters.character_name          AS  character_name,
      users.user_id                      AS  user_id
    FROM
      user_characters
    LEFT JOIN users ON users.user_id = user_characters.user_id
    LEFT JOIN characters ON characters.character_id = user_characters.character_id
    WHERE
      users.user_id = $1
  `
  rows, err := db.Query(sqlStatement, userID)
  if err != nil {
    return nil, err
  }
  defer rows.Close()

  userCharViews := make([]*UserCharacterView, 0)
  for rows.Next() {
    userCharView := new(UserCharacterView)
    err := rows.Scan(
      &userCharView.UserCharacterID,
      &userCharView.CharacterGsp,
      &userCharView.AltCostume,
      &userCharView.CharacterID,
      &userCharView.CharacterName,
      &userCharView.UserID,
    )
    if err != nil {
      return nil, err
    }

    userCharViews = append(userCharViews, userCharView)
  }

  err = rows.Err()
  if err != nil {
    return nil, err
  }

  return userCharViews, nil
}


// GetUserCharacterViewByUserCharacterID gets all of the data needed 
// to display a given "saved character", whicn includes 
// joined data fromt the user_characters and characters table
func (db *DB) GetUserCharacterViewByUserCharacterID(userCharID int64) (*UserCharacterView, error) {
  sqlStatement := `
    SELECT
      user_characters.user_character_id  AS  user_character_id,
      user_characters.character_gsp      AS  character_gsp,
      user_characters.alt_costume        AS  alt_costume,
      characters.character_id            AS  character_id,
      characters.character_name          AS  character_name,
      users.user_id                      AS  user_id
    FROM
      user_characters
    LEFT JOIN users on users.user_id = user_characters.user_id
    LEFT JOIN characters ON characters.character_id = user_characters.character_id
    WHERE
      user_character_id = $1
  `
  row := db.QueryRow(sqlStatement, userCharID)

  userCharView := new(UserCharacterView)
  err := row.Scan(
    &userCharView.UserCharacterID,
    &userCharView.CharacterGsp,
    &userCharView.AltCostume,
    &userCharView.CharacterID,
    &userCharView.CharacterName,
    &userCharView.UserID,
  )
  if err != nil {
    return nil, err
  }

  return userCharView, nil
}
