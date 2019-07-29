package db

import (
  "time"
)

// Match represents a recorded Smash Ultimate Online match outcome
type Match struct {
  MatchID               *int      `json:"matchId,omitempty"`
  UserID                int       `json:"userId"`
  UserName              *string   `json:"userName,omitempty"`
  UserCharacterName     *string   `json:"userCharacterName,omitempty"`
  UserCharacterGsp      *int      `json:"userCharacterGsp,omitempty"`
  UserWin               *bool     `json:"userWin,omitempty"`
  OpponentCharacterName string    `json:"opponentCharacterName"`
  OpponentCharacterGsp  *int      `json:"opponentCharacterGsp,omitempty"`
  OpponentTeabag        *bool     `json:"opponentTeabag,omitempty"`
  OpponentCamp          *bool     `json:"opponentCamp,omitempty"`
  OpponentAwesome       *bool     `json:"opponentAwesome,omitempty"`
  Created               time.Time `json:"created"`
}


// GetMatchByID gets a specific match by its id in the Matches table
func (db *DB) GetMatchByID(id int) (*Match, error) {
  row := db.QueryRow(`SELECT * FROM matches WHERE id = $1`, id)
  match := new(Match)
  err := row.Scan(
    &match.MatchID,
    &match.UserID,
    &match.UserCharacterName,
    &match.UserCharacterGsp,
    &match.UserWin,
    &match.OpponentCharacterName,
    &match.OpponentCharacterGsp,
    &match.OpponentTeabag,
    &match.OpponentCamp,
    &match.OpponentAwesome,
    &match.Created,
  )

  if err != nil {
    return nil, err
  }

  return match, nil
}

// GetAllMatches gets all of the matches from our database
func (db *DB) GetAllMatches() ([]*Match, error) {
  sqlStatement := `
  SELECT
    users.user_id,
    users.user_name,
    matches.match_id,
    matches.user_character_name,
    matches.user_character_gsp,
    matches.user_win,
    matches.opponent_character_name,
    matches.opponent_character_gsp,
    matches.opponent_teabag,
    matches.opponent_camp,
    matches.opponent_awesome,
    matches.created
  FROM
    users, matches
  WHERE
    users.user_id = matches.user_id;`
  rows, err := db.Query(sqlStatement)
  if err != nil {
    return nil, err
  }
  defer rows.Close()

  matches := make([]*Match, 0)
  for rows.Next() {
    match := new(Match)
    err := rows.Scan(
      &match.UserID,
      &match.UserName,
      &match.MatchID,
      &match.UserCharacterName,
      &match.UserCharacterGsp,
      &match.UserWin,
      &match.OpponentCharacterName,
      &match.OpponentCharacterGsp,
      &match.OpponentTeabag,
      &match.OpponentCamp,
      &match.OpponentAwesome,
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
    user_id,
    user_character_name,
    user_character_gsp,
    user_win,
    opponent_character_name,
    opponent_character_gsp,
    opponent_teabag,
    opponent_camp,
    opponent_awesome
  )
  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
  _, err := db.Exec(
    sqlStatement,
    match.UserID,
    match.UserCharacterName,
    match.UserCharacterGsp,
    match.UserWin,
    match.OpponentCharacterName,
    match.OpponentCharacterGsp,
    match.OpponentTeabag,
    match.OpponentCamp,
    match.OpponentAwesome,
  )

  if err != nil {
    return false, err
  }

  return true, nil
}
