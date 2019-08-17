package db


/*---------------------------------
            Interface
----------------------------------*/

// MatchManager describes all of the methods used
// to interact with the matches table in our database
type MatchManager interface {
  CreateMatch(matchCreate *MatchCreate) (int64, error)
  UpdateMatch(matchUpdate *MatchUpdate) (int64, error)
}


/*---------------------------------
          Data Structures
----------------------------------*/

// Match describes the required and optional data
// needed to create a new match in our matches table
type Match struct {
  MatchID               *int64         `json:"matchId,omitempty"`
  OpponentCharacterID   int64          `json:"opponentCharacterId"`
  UserID                int64          `json:"userId"`
  OpponentCharacterGsp  NullInt64JSON  `json:"opponentCharacterGsp"`
  UserCharacterID       NullInt64JSON  `json:"userCharacterId"`
  UserCharacterGsp      NullInt64JSON  `json:"userCharacterGsp"`
  UserWin               NullBoolJSON   `json:"userWin"`
}


// MatchUpdate describes the data needed 
// to update a given user's profile information
type MatchUpdate struct {
  MatchID               int64               `json:"matchId"`
  OpponentCharacterID   NullInt64JSON       `json:"opponentCharacterId"`
  OpponentCharacterGsp  NullInt64JSON       `json:"opponentCharacterGsp"`
  UserCharacterID       NullInt64JSON       `json:"userCharacterId"`
  UserCharacterGsp      NullInt64JSON       `json:"userCharacterGsp"`
  UserWin               NullBoolJSON        `json:"userWin"`
  Created               NullTimeJSON        `json:"created"`
  MatchTags             *[]*MatchTagCreate  `json:"matchTags"`
}


// MatchCreate describes the data needed 
// to create a given match in our db
type MatchCreate struct {
  OpponentCharacterID   int64               `json:"opponentCharacterId"`
  UserID                int64               `json:"userId"`
  OpponentCharacterGsp  NullInt64JSON       `json:"opponentCharacterGsp"`
  UserCharacterID       NullInt64JSON       `json:"userCharacterId"`
  UserCharacterGsp      NullInt64JSON       `json:"userCharacterGsp"`
  UserWin               NullBoolJSON        `json:"userWin"`
  MatchTags             *[]*MatchTagCreate  `json:"matchTags"`
}


/*---------------------------------
       Method Implementations
----------------------------------*/

// CreateMatch adds a new entry to the matches table in our database
func (db *DB) CreateMatch(matchCreate *MatchCreate) (int64, error) {
  var matchID int64
  sqlStatement := `
    INSERT INTO matches (
      opponent_character_id,
      user_id,
      opponent_character_gsp,
      user_character_id,
      user_character_gsp,
      user_win
    )
    VALUES ($1, $2, $3, $4, $5, $6)
    RETURNING
      match_id
  `
  row := db.QueryRow(
    sqlStatement,
    matchCreate.OpponentCharacterID,
    matchCreate.UserID,
    matchCreate.OpponentCharacterGsp,
    matchCreate.UserCharacterID,
    matchCreate.UserCharacterGsp,
    matchCreate.UserWin,
  )

  err := row.Scan(&matchID)

  if err != nil {
    return 0, err
  }

  return matchID, nil
}


// UpdateMatch updates an entry in the matches table with the given data
func (db *DB) UpdateMatch(matchUpdate *MatchUpdate) (int64, error) {
  var matchID int64
  sqlStatement := `
    UPDATE
      matches
    SET
      opponent_character_id = $1,
      opponent_character_gsp = $2,
      user_character_id = $3,
      user_character_gsp = $4,
      user_win = $5,
      created = $6
    WHERE
      match_id = $7
    RETURNING
      match_id
  `
  row := db.QueryRow(
    sqlStatement,
    matchUpdate.OpponentCharacterID,
    matchUpdate.OpponentCharacterGsp,
    matchUpdate.UserCharacterID,
    matchUpdate.UserCharacterGsp,
    matchUpdate.UserWin,
    matchUpdate.Created,
    matchUpdate.MatchID,
  )
  err := row.Scan(&matchID)

  if err != nil {
    return 0, nil
  }

  return matchID, nil
}
