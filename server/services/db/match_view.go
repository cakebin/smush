package db

import (
  "time"
)


/*---------------------------------
          Data Structures
----------------------------------*/

// MatchView desribes a JOIN between the matches, users, and characters tables,
// containing all of the data necessary to show a "match" in the front end
type MatchView struct {
  // Data from matches
  Created                time.Time  `json:"created"`
  MatchID                *int       `json:"matchId"`
  UserID                 int        `json:"userId"`
  UserCharacterID        *int       `json:"userCharacterId"`
  OpponentCharacterID    int        `json:"opponentCharacterId"`
  OpponentCharacterGsp   *int       `json:"opponentCharacterGsp,omitempty"`
  OpponentTeabag         *bool      `json:"opponentTeabag,omitempty"`
  OpponentCamp           *bool      `json:"opponentCamp,omitempty"`
  OpponentAwesome        *bool      `json:"opponentAwesome,omitempty"`
  UserCharacterGsp       *int       `json:"userCharacterGsp,omitempty"`
  UserWin                *bool      `json:"userWin,omitempty"`

  // Data from users
  UserName               *string    `json:"userName"`

  // Data from characters
  OpponentCharacterName  string     `json:"opponentCharacterName"`
  UserCharacterName      *string    `json:"userCharacterName"`
}


/*---------------------------------
            Interface
----------------------------------*/

// MatchViewManager describes all of the methods used to interact with
// match views in our database (data joined between match, character, user, etc)
type MatchViewManager interface {
  GetMatchViewByMatchID(matchID int) (*MatchView, error)
  GetAllMatchViews() ([]*MatchView, error)
}


/*---------------------------------
       Method Implementations
----------------------------------*/

// GetMatchViewByMatchID gets all of the data needed to display
// a recorded match, which includes joined data from the matches,
// users, and characters tables for a specific match_id
func (db *DB) GetMatchViewByMatchID(matchID int) (*MatchView, error) {
  sqlStatement := `
    SELECT
      matches.created                     AS created,
      matches.match_id                    AS match_id,
      users.user_id                       AS user_id,
      user_character.character_id         AS user_character_id,
      opponent_character.character_id     AS opponent_character_id,
      matches.opponent_character_gsp      AS opponent_character_gsp,
      matches.opponent_teabag             AS opponent_teabag,
      matches.opponent_camp               AS opponent_camp,
      matches.opponent_awesome            AS opponent_awesome,
      matches.user_character_gsp          AS user_character_gsp,
      matches.user_win                    AS user_win,
      users.user_name                     AS user_name,
      opponent_character.character_name   AS opponent_character_name,
      user_character.character_name       AS user_character_name
    FROM
      matches
    JOIN users ON users.user_id = matches.user_id
    JOIN characters opponent_character ON opponent_character.character_id = matches.opponent_character_id
    JOIN characters user_character ON user_character.character_id = matches.user_character_id
    WHERE
      matches.match_id = $1
  `

  row := db.QueryRow(sqlStatement, matchID)
  matchView := new(MatchView)
  err := row.Scan(
    &matchView.Created,
    &matchView.MatchID,
    &matchView.UserID,
    &matchView.UserCharacterID,
    &matchView.OpponentCharacterID,
    &matchView.OpponentCharacterGsp,
    &matchView.OpponentTeabag,
    &matchView.OpponentCamp,
    &matchView.OpponentAwesome,
    &matchView.UserCharacterGsp,
    &matchView.UserWin,
    &matchView.UserName,
    &matchView.OpponentCharacterName,
    &matchView.UserCharacterName,
  )

  if err != nil {
    return nil, err
  }

  return matchView, nil
}


// GetAllMatchViews gets all of the data needed to display all recorded matches,
// which includes joined data from the matches, users, and characters tables
func (db *DB) GetAllMatchViews() ([]*MatchView, error) {
  sqlStatement := `
    SELECT
      matches.created                     AS created,
      matches.match_id                    AS match_id,
      users.user_id                       AS user_id,
      user_character.character_id         AS user_character_id,
      opponent_character.character_id     AS opponent_character_id,
      matches.opponent_character_gsp      AS opponent_character_gsp,
      matches.opponent_teabag             AS opponent_teabag,
      matches.opponent_camp               AS opponent_camp,
      matches.opponent_awesome            AS opponent_awesome,
      matches.user_character_gsp          AS user_character_gsp,
      matches.user_win                    AS user_win,
      users.user_name                     AS user_name,
      opponent_character.character_name   AS opponent_character_name,
      user_character.character_name       AS user_character_name
    FROM
      matches
    JOIN users ON users.user_id = matches.user_id
    JOIN characters opponent_character ON opponent_character.character_id = matches.opponent_character_id
    JOIN characters user_character ON user_character.character_id = matches.user_character_id
  `

  rows, err := db.Query(sqlStatement)
  if err != nil {
    return nil, err
  }
  defer rows.Close()

  matchViews := make([]*MatchView, 0)
  for rows.Next() {
    matchView := new(MatchView)
    err := rows.Scan(
      &matchView.Created,
      &matchView.MatchID,
      &matchView.UserID,
      &matchView.UserCharacterID,
      &matchView.OpponentCharacterID,
      &matchView.OpponentCharacterGsp,
      &matchView.OpponentTeabag,
      &matchView.OpponentCamp,
      &matchView.OpponentAwesome,
      &matchView.UserCharacterGsp,
      &matchView.UserWin,
      &matchView.UserName,
      &matchView.OpponentCharacterName,
      &matchView.UserCharacterName,
    )

    if err != nil {
      return nil, err
    }

    matchViews = append(matchViews, matchView)
  }

  err = rows.Err()
  if err != nil {
    return nil, err
  }

  return matchViews, nil
}
