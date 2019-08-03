package db

import (
  "database/sql"

  _ "github.com/lib/pq" // Needed for the postgres driver
)


// DB is the struct that we're going to use to implement all of out
// Datasbase interfaces; All of the methods defined on each of our
// interfaces will be implemented on this DB struct
type DB struct {
  *sql.DB
}


type DatabaseManager interface {
  MatchManager
  MatchViewManager
  UserManager
  UserViewManager
  CharacterManager
}


// New initializes a new postgres database connection and attaches
// said connection to our DB struct, which we can then call all of
// the methods described by the our varies Database interfaces
func New(dataSourceName string) (*DB, error) {
  db, err := sql.Open("postgres", dataSourceName)
  if err != nil {
    return nil, err
  }
  if err = db.Ping(); err != nil {
    return nil, err
  }
  return &DB{db}, nil
}


