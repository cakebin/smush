package db


/*---------------------------------
          Data Structures
----------------------------------*/

// MatchTag describes a "match tag" relationship
type MatchTag struct {
  MatchTagID  int64  `json:"matchTagId"`
  MatchID     int64  `json:"matchId"`
  TagID       int64  `json:"tagId"`
}

// MatchTagCreate describes the data needed
// to create a "match tag" relationship
type MatchTagCreate struct {
  MatchID  int64  `json:"matchId"`
  TagID    int64  `json:"tagId"`
}

/*---------------------------------
            Interface
----------------------------------*/

// MatchTagManager describes all of the methods used
// to interact with the match_tags table in our database
type MatchTagManager interface {
  CreateMatchTags(matchTagsCreate []*MatchTagCreate) ([]int64, error)
  DeleteMatchTagsByMatchID(matchID int64) ([]int64, error)
}


/*---------------------------------
       Method Implementations
----------------------------------*/

// CreateMatchTags creates adds multiple "match tag" relationships
func (db *DB) CreateMatchTags(matchTagsCreate []*MatchTagCreate) ([]int64, error) {
  table := "match_tags"
  columns := []string{"match_id", "tag_id"}
  numInserts := len(matchTagsCreate)
  returningCol := "match_tag_id"
  sqlStatement := MakeMultiInsertStatement(table, columns, numInserts, returningCol)
  expandedMatchTags := expandMatchTagsCreate(matchTagsCreate)
  rows, err := db.Query(sqlStatement, expandedMatchTags...)
  if err != nil {
    return nil, err
  }
  defer rows.Close()
  
  matchTagIDs := make([]int64, 0)
  for rows.Next() {
    var matchTagID int64
    err := rows.Scan(&matchTagID)
    if err != nil {
      return nil, err
    }

    matchTagIDs = append(matchTagIDs, matchTagID)
  }

  err = rows.Err()
  if err != nil {
    return nil, err
  }

  return matchTagIDs, nil
}


// DeleteMatchTagsByMatchID deletes a set of existing
// "match tag" relationships for a given matchID
func (db *DB) DeleteMatchTagsByMatchID(matchID int64) ([]int64, error) {
  sqlStatement := `
    DELETE FROM
      match_tags
    WHERE
      match_tag_id IN (SELECT match_tag_id FROM match_tags WHERE match_id = $1)
    RETURNING
      match_tag_id
  `
  rows, err := db.Query(sqlStatement, matchID)
  if err != nil {
    return nil, err
  }
  defer rows.Close()

  deletedMatchTagIDs := make([]int64, 0)
  for rows.Next() {
    var deletedMatchTagID int64
    err := rows.Scan(&deletedMatchTagID)
    if err != nil {
      return nil, err
    }

    deletedMatchTagIDs = append(deletedMatchTagIDs, deletedMatchTagID)
  }

  err = rows.Err()
  if err != nil {
    return nil, err
  }

  return deletedMatchTagIDs, nil
}


/*---------------------------------
            Helpers
----------------------------------*/

// ExpandMatchTags expands a slice of MatchTagCreate into a slice of
// sequential int, intended to be used in conjunction with a multi-insert
// statement (i.e. for use in CreateMatchTags)
func expandMatchTagsCreate(matchTagsCreate []*MatchTagCreate) ([]interface{}) {
  expandedMatchTags := make([]interface{}, 0)
  
  for _, matchTagCreate := range matchTagsCreate {
    expandedMatchTags = append(expandedMatchTags, matchTagCreate.MatchID)
    expandedMatchTags = append(expandedMatchTags, matchTagCreate.TagID)
  }
  
  return expandedMatchTags
}
