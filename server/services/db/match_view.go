package db

import (
  "time"
)


/*---------------------------------
            Interface
----------------------------------*/

// MatchViewManager describes all of the methods used to interact with
// match views in our database (data joined between match, character, user, etc)
type MatchViewManager interface {
  GetMatchViewByMatchID(matchID int64) (*MatchView, error)
  GetAllMatchViews() ([]*MatchView, error)
}


/*---------------------------------
          Data Structures
----------------------------------*/

// MatchView desribes a JOIN between the matches, users, and characters tables,
// containing all of the data necessary to show a "match" in the front end
type MatchView struct {
  // Data from matches
  Created                time.Time        `json:"created"`
  UserID                 int64            `json:"userId"`
  MatchID                int64            `json:"matchId"`
  OpponentCharacterID    int64            `json:"opponentCharacterId"`

  UserCharacterID        NullInt64JSON    `json:"userCharacterId"`
  OpponentCharacterGsp   NullInt64JSON    `json:"opponentCharacterGsp,omitempty"`
  UserCharacterGsp       NullInt64JSON    `json:"userCharacterGsp,omitempty"`
  UserWin                NullBoolJSON     `json:"userWin,omitempty"`

  // Data from users
  UserName               string            `json:"userName"`

  // Data from characters
  OpponentCharacterName  string           `json:"opponentCharacterName"`
  UserCharacterName      NullStringJSON   `json:"userCharacterName,omitempty"`
  OpponentCharacterImg   string           `json:"opponentCharacterImage"`
  UserCharacterImg       NullStringJSON   `json:"userCharacterImage"`

  // Data from user characters
  AltCostume             NullInt64JSON    `json:"altCostume,omitempty"`

  // Data from match_tags; added seperately from the SQL Joins
  MatchTags              []*MatchTagView  `json:"matchTags"`
}


/*---------------------------------
       Method Implementations
----------------------------------*/

// GetMatchViewByMatchID gets all of the data needed to display
// an individual match, which includes joined data from the users and characters table
func (db *DB) GetMatchViewByMatchID(matchID int64) (*MatchView, error) {
  sqlStatement := `
    SELECT
      matches.created                         AS created,
      matches.match_id                        AS match_id,
      users.user_id                           AS user_id,
      player_character.character_id           AS player_character_id,
      opponent_character.character_id         AS opponent_character_id,
      matches.opponent_character_gsp          AS opponent_character_gsp,
      matches.user_character_gsp              AS player_character_gsp,
      matches.user_win                        AS user_win,
      users.user_name                         AS user_name,
      opponent_character.character_name       AS opponent_character_name,
      player_character.character_name         AS player_character_name,
      opponent_character.character_stock_img  AS opponent_character_img,
      player_character.character_stock_img    AS player_character_img,
      user_characters.alt_costume             AS alt_costume
    FROM
      matches
    LEFT JOIN users ON users.user_id = matches.user_id
    LEFT JOIN characters opponent_character ON opponent_character.character_id = matches.opponent_character_id
    LEFT JOIN characters player_character ON player_character.character_id = matches.user_character_id
    LEFT JOIN user_characters ON user_characters.character_id = matches.user_character_id AND user_characters.user_id = matches.user_id
    WHERE
     match_id = $1
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
    &matchView.UserCharacterGsp,
    &matchView.UserWin,
    &matchView.UserName,
    &matchView.OpponentCharacterName,
    &matchView.UserCharacterName,
    &matchView.OpponentCharacterImg,
    &matchView.UserCharacterImg,
    &matchView.AltCostume,
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
      matches.created                         AS created,
      matches.match_id                        AS match_id,
      users.user_id                           AS user_id,
      player_character.character_id           AS player_character_id,
      opponent_character.character_id         AS opponent_character_id,
      matches.opponent_character_gsp          AS opponent_character_gsp,
      matches.user_character_gsp              AS player_character_gsp,
      matches.user_win                        AS user_win,
      users.user_name                         AS user_name,
      opponent_character.character_name       AS opponent_character_name,
      player_character.character_name         AS player_character_name,
      opponent_character.character_stock_img  AS opponent_character_img,
      player_character.character_stock_img    AS player_character_img,
      user_characters.alt_costume             AS alt_costume
    FROM
      matches
    LEFT JOIN users ON users.user_id = matches.user_id
    LEFT JOIN characters opponent_character ON opponent_character.character_id = matches.opponent_character_id
    LEFT JOIN characters player_character ON player_character.character_id = matches.user_character_id
    LEFT JOIN user_characters ON user_characters.character_id = matches.user_character_id AND user_characters.user_id = matches.user_id
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
      &matchView.UserCharacterGsp,
      &matchView.UserWin,
      &matchView.UserName,
      &matchView.OpponentCharacterName,
      &matchView.UserCharacterName,
      &matchView.OpponentCharacterImg,
      &matchView.UserCharacterImg,
      &matchView.AltCostume,
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
