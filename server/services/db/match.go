package db

import (
  "database/sql"
)


/*---------------------------------
          Data Structures
----------------------------------*/

// Match describes the required and optional data
// needed to create a new match in our matches table
type Match struct {
  MatchID               *int           `json:"matchId,omitempty"`
  OpponentCharacterID   int            `json:"opponentCharacterId"`
  UserID                int            `json:"userId"`
  OpponentCharacterGsp  sql.NullInt64  `json:"opponentCharacterGsp"`
  OpponentTeabag        sql.NullBool   `json:"opponentTeabag"`
  OpponentCamp          sql.NullBool   `json:"opponentCamp"`
  OpponentAwesome       sql.NullBool   `json:"opponentAwesome"`
  UserCharacterID       sql.NullInt64  `json:"userCharacterId"`
  UserCharacterGsp      sql.NullInt64  `json:"userCharacterGsp"`
  UserWin               sql.NullBool   `json:"userWin"`
}


// MatchCreate describes the data needed 
// to create a given match in our db
type MatchCreate struct {
  OpponentCharacterID   int            `json:"opponentCharacterId"`
  UserID                int            `json:"userId"`
  OpponentCharacterGsp  sql.NullInt64  `json:"opponentCharacterGsp"`
  OpponentTeabag        sql.NullBool   `json:"opponentTeabag"`
  OpponentCamp          sql.NullBool   `json:"opponentCamp"`
  OpponentAwesome       sql.NullBool   `json:"opponentAwesome"`
  UserCharacterID       sql.NullInt64  `json:"userCharacterId"`
  UserCharacterGsp      sql.NullInt64  `json:"userCharacterGsp"`
  UserWin               sql.NullBool   `json:"userWin"`
}

/*---------------------------------
            Interface
----------------------------------*/

// MatchManager describes all of the methods used
// to interact with the matches table in our database
type MatchManager interface {
  CreateMatch(matchCreate *MatchCreate) (int, error)
}

/*---------------------------------
       Method Implementations
----------------------------------*/

// CreateMatch adds a new entry to the matches table in our database
func (db *DB) CreateMatch(matchCreate *MatchCreate) (int, error) {
  var matchID int
  sqlStatement := `
    INSERT INTO matches (
      opponent_character_id,
      user_id,
      opponent_character_gsp,
      opponent_teabag,
      opponent_camp,
      opponent_awesome,
      user_character_id,
      user_character_gsp,
      user_win
    )
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    RETURNING
      match_id
  `
  row := db.QueryRow(
    sqlStatement,
    matchCreate.OpponentCharacterID,
    matchCreate.UserID,
    matchCreate.OpponentCharacterGsp,
    matchCreate.OpponentTeabag,
    matchCreate.OpponentCamp,
    matchCreate.OpponentAwesome,
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
