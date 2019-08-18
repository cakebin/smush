package db


/*---------------------------------
          Data Structures
----------------------------------*/

// Tag describes the required data needed
// to create a new tag in our tags table
type Tag struct {
  TagID    int     `json:"tagId"`
  TagName  string  `json:"tagName"`
}


// TagCreate describes the data needed
// to create a new tag in our database
type TagCreate struct {
  TagName  string  `json:"tagName"`
}


// TagUpdate describes the data needed 
// to update a given tag in our database
type TagUpdate struct {
  TagID    int     `json:"tagId"`
  TagName  string  `json:"tagName"`
}


/*---------------------------------
            Interface
----------------------------------*/

// TagManager describes all of the methods used
// to interact with the tags table in our database
type TagManager interface {
  GetAllTags() ([]*Tag, error)
  GetTagByTagID(tagID int) (*Tag, error)

  CreateTag(tagCreate *TagCreate) (int, error)
  UpdateTag(tagUpdate *TagUpdate) (int, error)
}


/*---------------------------------
       Method Implementations
----------------------------------*/

// GetAllTags gets all of the tags we have in our database
func (db *DB) GetAllTags() ([]*Tag, error) {
  sqlStatement := `
    SELECT
      tag_id,
      tag_name
    FROM
      tags
  `
  rows, err := db.Query(sqlStatement)
  if err != nil {
    return nil, err
  }
  defer rows.Close()

  tags := make([]*Tag, 0)
  for rows.Next() {
    tag := new(Tag)

    err := rows.Scan(
      &tag.TagID,
      &tag.TagName,
    )
    if err != nil {
      return nil, err
    }

    tags = append(tags, tag)
  }

  err = rows.Err()
  if err != nil {
    return nil, err
  }

  return tags, nil
}


// GetTagByTagID gets a specific tag given a tagID
func (db *DB) GetTagByTagID(tagID int) (*Tag, error) {
  sqlStatement := `
    SELECT
      tag_id,
      tag_name
    FROM
      tags
    WHERE
      tag_id = $1
  `
  row := db.QueryRow(sqlStatement, tagID)

  tag := new(Tag)
  err := row.Scan(
    &tag.TagID,
    &tag.TagName,
  )
  if err != nil {
    return nil, err
  }

  return tag, nil
}


// CreateTag adds a new entry to the tags table
func (db *DB) CreateTag(tagCreate *TagCreate) (int, error) {
  sqlStatement := `
    INSERT INTO tags
      (tag_name)
    VALUES
      ($1)
    RETURNING
      tag_id
  `
  row := db.QueryRow(sqlStatement, tagCreate.TagName)

  var tagID int
  err := row.Scan(&tagID)
  if err != nil {
    return 0, err
  }

  return tagID, nil
}


// UpdateTag updates an existing entry in the tags table
func (db *DB) UpdateTag(tagUpdate *TagUpdate) (int, error) {
  sqlStatement := `
    UPDATE
      tags
    SET
      tag_name = $1
    WHERE
      tag_id = $2
    RETURNING
      tag_id
  `
  row := db.QueryRow(
    sqlStatement,
    tagUpdate.TagName,
    tagUpdate.TagID,
  )

  var tagID int
  err := row.Scan(&tagID)
  if err != nil {
    return 0, err
  }

  return tagID, nil
}
