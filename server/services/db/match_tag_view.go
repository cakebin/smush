package db


/*---------------------------------
          Data Structures
----------------------------------*/

// MatchTagView describes a JOIN between match_tags and tags tables
type MatchTagView struct {
  MatchTagID  int     `json:"matchTagId"`
  MatchID     int     `json:"matchId"`
  TagID       int     `json:"tagId"`
  TagName     string  `json:"tagName"`
}


/*---------------------------------
            Interface
----------------------------------*/

// MatchTagViewManager describes all of the methods
// used to interact with "match tag" views in our database
type MatchTagViewManager interface {
  GetAllMatchTagViews() ([]*MatchTagView, error)
  GetMatchTagViewsByMatchID(matchID int) ([]*MatchTagView, error)
}


/*---------------------------------
       Method Implementations
----------------------------------*/

// GetAllMatchTagViews gets all of the match tags
func (db *DB) GetAllMatchTagViews() ([]*MatchTagView, error) {
  sqlStatement := `
    SELECT
      match_tags.match_tag_id  AS  match_tag_id,
      match_tags.match_id      AS  match_id,
      tags.tag_id              AS  tag_id,
      tags.tag_name            AS  tag_name
    FROM
      match_tags
    LEFT JOIN tags ON tags.tag_id = match_tags.tag_ig
  `
  rows, err := db.Query(sqlStatement)
  if err != nil {
    return nil, err
  }
  defer rows.Close()

  matchTagViews := make([]*MatchTagView, 0)
  for rows.Next() {
    matchTagView := new(MatchTagView)
    err := rows.Scan(
      &matchTagView.MatchTagID,
      &matchTagView.MatchID,
      &matchTagView.TagID,
      &matchTagView.TagName,
    )
    if err != nil {
      return nil, err
    }

    matchTagViews = append(matchTagViews, matchTagView)
  }

  err = rows.Err()
  if err != nil {
    return nil, err
  }

  return matchTagViews, nil
}



// GetMatchTagViewsByMatchID gets all of the match tags for a given matchID
func (db *DB) GetMatchTagViewsByMatchID(matchID int) ([]*MatchTagView, error) {
  sqlStatement := `
    SELECT
      match_tags.match_tag_id  AS  match_tag_id,
      match_tags.match_id      AS  match_id,
      tags.tag_id              AS  tag_id,
      tags.tag_name            AS  tag_name
    FROM
      match_tags
    LEFT JOIN tags ON tags.tag_id = match_tags.tag_ig
    WHERE
      match_tags.match_id = $1
  `
  rows, err := db.Query(sqlStatement, matchID)
  if err != nil {
    return nil, err
  }
  defer rows.Close()

  matchTagViews := make([]*MatchTagView, 0)
  for rows.Next() {
    matchTagView := new(MatchTagView)
    err := rows.Scan(
      &matchTagView.MatchTagID,
      &matchTagView.MatchID,
      &matchTagView.TagID,
      &matchTagView.TagName,
    )
    if err != nil {
      return nil, err
    }

    matchTagViews = append(matchTagViews, matchTagView)
  }

  err = rows.Err()
  if err != nil {
    return nil, err
  }

  return matchTagViews, nil
}
