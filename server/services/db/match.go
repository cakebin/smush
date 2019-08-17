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
  MatchID               *int64           `json:"matchId,omitempty"`
  OpponentCharacterID   int64            `json:"opponentCharacterId"`
  UserID                int64            `json:"userId"`
  OpponentCharacterGsp  NullInt64JSON  `json:"opponentCharacterGsp"`
  OpponentTeabag        NullBoolJSON   `json:"opponentTeabag"`
  OpponentCamp          NullBoolJSON   `json:"opponentCamp"`
  OpponentAwesome       NullBoolJSON   `json:"opponentAwesome"`
  UserCharacterID       NullInt64JSON  `json:"userCharacterId"`
  UserCharacterGsp      NullInt64JSON  `json:"userCharacterGsp"`
  UserWin               NullBoolJSON   `json:"userWin"`
}


// MatchUpdate describes the data needed 
// to update a given user's profile information
type MatchUpdate struct {
  MatchID               int64            `json:"matchId"`
  OpponentCharacterID   NullInt64JSON  `json:"opponentCharacterId"`
  OpponentCharacterGsp  NullInt64JSON  `json:"opponentCharacterGsp"`
  OpponentTeabag        NullBoolJSON   `json:"opponentTeabag"`
  OpponentCamp          NullBoolJSON   `json:"opponentCamp"`
  OpponentAwesome       NullBoolJSON   `json:"opponentAwesome"`
  UserCharacterID       NullInt64JSON  `json:"userCharacterId"`
  UserCharacterGsp      NullInt64JSON  `json:"userCharacterGsp"`
  UserWin               NullBoolJSON   `json:"userWin"`
  Created               NullTimeJSON   `json:"created"`
}


// MatchCreate describes the data needed 
// to create a given match in our db
type MatchCreate struct {
  OpponentCharacterID   int64            `json:"opponentCharacterId"`
  UserID                int64            `json:"userId"`
  OpponentCharacterGsp  NullInt64JSON  `json:"opponentCharacterGsp"`
  OpponentTeabag        NullBoolJSON   `json:"opponentTeabag"`
  OpponentCamp          NullBoolJSON   `json:"opponentCamp"`
  OpponentAwesome       NullBoolJSON   `json:"opponentAwesome"`
  UserCharacterID       NullInt64JSON  `json:"userCharacterId"`
  UserCharacterGsp      NullInt64JSON  `json:"userCharacterGsp"`
  UserWin               NullBoolJSON   `json:"userWin"`
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


// UpdateMatch updates an entry in the matches table with the given data
func (db *DB) UpdateMatch(matchUpdate *MatchUpdate) (int64, error) {
  var matchID int64
  sqlStatement := `
    UPDATE
      matches
    SET
      opponent_character_id = $1,
      opponent_character_gsp = $2,
      opponent_teabag = $3,
      opponent_camp = $4,
      opponent_awesome = $5,
      user_character_id = $6,
      user_character_gsp = $7,
      user_win = $8,
      created = $9
    WHERE
      match_id = $10
    RETURNING
      match_id
  `
  row := db.QueryRow(
    sqlStatement,
    matchUpdate.OpponentCharacterID,
    matchUpdate.OpponentCharacterGsp,
    matchUpdate.OpponentTeabag,
    matchUpdate.OpponentCamp,
    matchUpdate.OpponentAwesome,
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
