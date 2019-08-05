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

  OpponentCharacterGsp  sql.NullInt64  `json:"opponentCharacterGsp,omitempty"`
  OpponentTeabag        sql.NullBool   `json:"opponentTeabag,omitempty"`
  OpponentCamp          sql.NullBool   `json:"opponentCamp,omitempty"`
  OpponentAwesome       sql.NullBool   `json:"opponentAwesome,omitempty"`
  UserCharacterID       sql.NullInt64  `json:"userCharacterId,omitempty"`
  UserCharacterGsp      sql.NullInt64  `json:"userCharacterGsp,omitempty"`
  UserWin               sql.NullBool   `json:"userWin,omitempty"`
}

/*---------------------------------
            Interface
----------------------------------*/

// MatchManager describes all of the methods used
// to interact with the matches table in our database
type MatchManager interface {
  CreateMatch(match Match) (*Match, error)
}

/*---------------------------------
       Method Implementations
----------------------------------*/

// CreateMatch adds a new entry to the matches table in our database
func (db *DB) CreateMatch(match Match) (*Match, error) {

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
    match.OpponentCharacterID,
    match.UserID,
    match.OpponentCharacterGsp.Int64,
    match.OpponentTeabag.Bool,
    match.OpponentCamp.Bool,
    match.OpponentAwesome.Bool,
    match.UserCharacterID.Int64,
    match.UserCharacterGsp.Int64,
    match.UserWin.Bool,
  )

  createdMatch := new(Match)
  err := row.Scan(&createdMatch.MatchID)

  if err != nil {
    return nil, err
  }

  return createdMatch, nil
}
