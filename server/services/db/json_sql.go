package db

import (
  "encoding/json"
  "database/sql"
  "time"

  "github.com/lib/pq"
)


// NullInt64JSON extends sql.NullInt64 to nicely (Un)Marshal JSON
type NullInt64JSON struct {
  sql.NullInt64
}


// MarshalJSON handles sql.NullInt64 to JSON
func (ni *NullInt64JSON) MarshalJSON() ([]byte, error) {
  if !ni.Valid {
    return []byte("null"), nil
  }

  return json.Marshal(ni.Int64)
}


// UnmarshalJSON handles JSON to sql.NullInt64
func (ni *NullInt64JSON) UnmarshalJSON(data []byte) error {
  // Unmarshalling into a pointer will let us detect null
  var integer *int64
  err := json.Unmarshal(data, &integer)
  if err != nil {
    return err
  }
 
  if integer != nil {
      ni.Valid = true
      ni.Int64 = *integer
  } else {
      ni.Valid = false
  }
  return nil
}


// NullStringJSON extends sql.NullString to nicely (Un)Marshal JSON
type NullStringJSON struct {
  sql.NullString
}


// MarshalJSON handles sql.NullString to JSON
func (ns *NullStringJSON) MarshalJSON() ([]byte, error) {
  if !ns.Valid {
    return []byte("null"), nil
  }

  return json.Marshal(ns.String)
}


// UnmarshalJSON handles JSON to sql.NullString
func (ns *NullStringJSON) UnmarshalJSON(data []byte) error {
  // Unmarshalling into a pointer will let us detect null
  var str *string
  err := json.Unmarshal(data, &str)
  if err != nil {
    return err
  }
 
  if str != nil {
      ns.Valid = true
      ns.String = *str
  } else {
      ns.Valid = false
  }
  return nil
}


// NullBoolJSON extends sql.NullBool to nicely (Un)Marshal JSON
type NullBoolJSON struct {
  sql.NullBool
}


// MarshalJSON handles sql.NullBool to JSON
func (nb *NullBoolJSON) MarshalJSON() ([]byte, error) {
  if !nb.Valid {
    return []byte("null"), nil
  }

  return json.Marshal(nb.Bool)
}


// UnmarshalJSON handles JSON to sql.NullBool
func (nb *NullBoolJSON) UnmarshalJSON(data []byte) error {
  // Unmarshalling into a point will let us detect null
  var boolean *bool
  err := json.Unmarshal(data, &boolean)
  if err != nil {
    return err
  }

  if boolean != nil {
    nb.Valid = true
    nb.Bool = *boolean
  } else {
    nb.Valid = false
  }

  return nil
}


// NullTimeJSON extends pq.NullTime to nicely (Un)Marshal JSON
type NullTimeJSON struct {
  pq.NullTime
}


// MarshalJSON handles pq.NullTime to JSON
func (nt *NullTimeJSON) MarshalJSON() ([]byte, error) {
  if !nt.Valid {
    return []byte("null"), nil
  }

  return json.Marshal(nt.Time)
}


// UnmarshalJSON handles JSON to pq.NullTime
func (nt *NullTimeJSON) UnmarshalJSON(data []byte) error {
  var t *time.Time
  err := json.Unmarshal(data, &t)
  if err != nil {
    return err
  }

  if t != nil {
    nt.Valid = true
    nt.Time = *t
  }

  return nil
}
