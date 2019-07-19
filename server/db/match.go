package db


// Match represents a recorded Smash Ultimate Online match outcome
type Match struct {
  OpponentCharacterName   string
  OpponentCharacterGsp    int
  OpponentTeabag          bool
  OpponentCamp            bool
  OpponentAwesome         bool
  UserCharacterName       string
  UserCharacterGsp        int
  UserWin                 bool
}

// AllMatches gets all of the matches from our database
func (db *DB) AllMatches() ([]*Match, error) {
  rows, err := db.Query("SELECT * FROM books")
  if err != nil {
      return nil, err
  }
  defer rows.Close()

  matches := make([]*Match, 0)
  for rows.Next() {
      match := new(Match)
      err := rows.Scan(
        &match.OpponentCharacterName,
        &match.OpponentCharacterGsp,
        &match.OpponentTeabag,
        &match.OpponentCamp,
        &match.OpponentAwesome,
        &match.UserCharacterName,
        &match.OpponentCharacterGsp,
        &match.UserWin,
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
