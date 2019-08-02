package database


/*---------------------------------
          Data Structures
----------------------------------*/

// Match describes the required and optional data
// needed to create a new match in our matches table
type Match struct {
  OpponentCharacterID   int       `json:"opponentCharacterId"`
  UserID                int       `json:"userId"`

  OpponentCharacterGsp  *int      `json:"opponentCharacterGsp,omitempty"`
  OpponentTeabag        *bool     `json:"opponentTeabag,omitempty"`
  OpponentCamp          *bool     `json:"opponentCamp,omitempty"`
  OpponentAwesome       *bool     `json:"opponentAwesome,omitempty"`
  UserCharacterID       *int      `json:"userCharacterId,omitempty"`
  UserCharacterGsp      *int      `json:"userCharacterGsp,omitempty"`
  UserWin               *bool     `json:"userWin,omitempty"`
}


/*---------------------------------
            Interface
----------------------------------*/

// MatchManager describes all of the methods used
// to interact with the matches table in our database
type MatchManager interface {
  CreateMatch(match Match) (int, error)
}


/*---------------------------------
       Method Implementations
----------------------------------*/

// CreateMatch adds a new entry to the matches table in our database
func (db *DB) CreateMatch(match Match) (int, error) {
  var matchID int
  sqlStatement := `
    INSERT INTO matches (
      opponent_character_id,
      user_id,
      opponent_character_gsp,
      opponent_teabag,
      opponent_camp,
      opponent_awesome,
      user_character_id
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
    match.OpponentCharacterGsp,
    match.OpponentTeabag,
    match.OpponentCamp,
    match.OpponentAwesome,
    match.UserCharacterID,
    match.UserCharacterGsp,
    match.UserWin,
  )
  err := row.Scan(matchID)

  if err != nil {
    return 0, err
  }

  return matchID, nil
}
