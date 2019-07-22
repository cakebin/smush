package db

import (
  "time"
)


// Match represents a recorded Smash Ultimate Online match outcome
type Match struct {
  ID                      int          `json:"id,omitempty"`
  OpponentCharacterName   string       `json:"opponentCharacterName"`
  OpponentCharacterGsp    int          `json:"opponentCharacterGsp,omitempty"`
  OpponentTeabag          bool         `json:"opponentTeabag,omitempty"`
  OpponentCamp            bool         `json:"opponentCamp,omitempty"`
  OpponentAwesome         bool         `json:"opponentAwesome,omitempty"`
  UserCharacterName       string       `json:"userCharacterName,omitempty"`
  UserCharacterGsp        int          `json:"userCharacterGsp,omitempty"`
  UserWin                 bool         `json:"userWin,omitempty"`
  Created                 time.Time    `json:"created"`
}


// MatchResponse represents an interaction with our database
// regarding operations related to the Matches table
type MatchResponse struct {
  Success  bool     `json:"success"`
  Error    error    `json:"error"`
}


// GetMatchByID gets a specific match by its id in the Matches table
func (db *DB) GetMatchByID(id int) (*Match, error) {
  row := db.QueryRow(`SELECT * FROM matches WHERE id = $1`, id)
  match := new(Match)
  err := row.Scan(
    &match.ID,
    &match.OpponentCharacterName,
    &match.OpponentCharacterGsp,
    &match.OpponentTeabag,
    &match.OpponentCamp,
    &match.OpponentAwesome,
    &match.UserCharacterName,
    &match.UserCharacterGsp,
    &match.UserWin,
    &match.Created,
  )

  if err != nil {
    return nil, err
  }

  return match, nil
} 


// GetAllMatches gets all of the matches from our database
func (db *DB) GetAllMatches() ([]*Match, error) {
  rows, err := db.Query("SELECT * FROM matches")
  if err != nil {
    return nil, err
  }
  defer rows.Close()

  matches := make([]*Match, 0)
  for rows.Next() {
    match := new(Match)
    err := rows.Scan(
      &match.ID,
      &match.OpponentCharacterName,
      &match.OpponentCharacterGsp,
      &match.OpponentTeabag,
      &match.OpponentCamp,
      &match.OpponentAwesome,
      &match.UserCharacterName,
      &match.UserCharacterGsp,
      &match.UserWin,
      &match.Created,
    )

    if err != nil {
      return nil, err
    }

    matches = append(matches, match)
  }

  if err = rows.Err(); err != nil {
    return nil, err
  }

  return matches, nil
}


// CreateMatch adds a new entry to the Matches table in our database
func (db *DB) CreateMatch(match Match) (bool, error) {
  sqlStatement := `
  INSERT INTO matches (
    opponent_character_name,
    opponent_character_gsp,
    opponent_teabag,
    opponent_camp,
    opponent_awesome,
    user_character_name,
    user_character_gsp,
    user_win
  )
  VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
  _, err := db.Exec(
    sqlStatement,
    match.OpponentCharacterName,
    match.OpponentCharacterGsp,
    match.OpponentTeabag,
    match.OpponentCamp,
    match.OpponentAwesome,
    match.UserCharacterName,
    match.UserCharacterGsp,
    match.UserWin,
  )

  if err != nil {
    return false, err
  }

  return true, nil
}
