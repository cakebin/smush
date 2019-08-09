package api

import (
  "database/sql"
  "time"

  "github.com/lib/pq"
)


// ToNullInt64 takes a pointer to an int64 and converts
// it to a sql.NullInt64 value; useful for converting
// optional API request int64 data for use in sql
func ToNullInt64(integer *int64) sql.NullInt64 {
  if integer == nil {
    return sql.NullInt64{}
  }

  return sql.NullInt64{Int64: *integer, Valid: true}
}


// ToNullBool takes a pointer to a bool and converts
// it to a sql.NullBool value; useful for converting
// optional API request bool data for use in sql
func ToNullBool(boolean *bool) sql.NullBool {
  if boolean == nil {
    return sql.NullBool{}
  }

  return sql.NullBool{Bool: *boolean, Valid: true}
}


// ToNullString takes a pointer to a string and converts
// it to a sql.NullString value; useful for converting
// optional API request string data for use in sql
func ToNullString(str *string) sql.NullString {
  if str == nil {
    return sql.NullString{}
  }

  return sql.NullString{String: *str, Valid: true}
}


func ToNullTime(t *time.Time) pq.NullTime {
  if t == nil {
    return pq.NullTime{}
  }

  return pq.NullTime{Time: *t, Valid: true}
}


// FromNullString converts a sql.NullString (usually scanned
// from a db.Query) and converts it to pointer to a normal
// string or nil if it's a "NULL" value in the query
func FromNullString(nullStr sql.NullString) *string {
  if !nullStr.Valid {
    return nil
  }

  return &nullStr.String
}


// FromNullBool converts a sql.NullBool (usually scanned
// from a db.Query) and converts it to a pointer to a normal
// bool or nil if it's a "NULL" value in the query
func FromNullBool(nullBool sql.NullBool) *bool {
  if !nullBool.Valid {
    return nil
  }

  return &nullBool.Bool
}


// FromNullInt64 converts a sql.NullBool (usually scanned
// from a db.Query) and converts it to a pointer to a normal bool
func FromNullInt64(nullInt sql.NullInt64) *int64 {
  if !nullInt.Valid {
    return nil
  }

  return &nullInt.Int64
}
