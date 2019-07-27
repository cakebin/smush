package db

import (
  "database/sql"
  _ "github.com/lib/pq"  // Needed for postgres
)


// Datastore defines all of the methods for interacting with our Database
// See specific files for their corresponding method implementations
// (i.e. for AllMatches, see 'server/db/match.go')
type Datastore interface {
  // Match API
  GetMatchByID(id int) (*Match, error)
  GetAllMatches() ([]*Match, error)
  CreateMatch(match Match) (bool, error)
  // UpdateMatch() (MatchResponse, error)
  // DeleteMatch() (MatchResponse, error)

  // User API
  GetUserByID(id int) (*User, error)
  CreateUser(user User) (bool, error)
}


// DB is the struct that we're going to use to implement our Datastore
// interface; All of the methods defined on Datastore will be implemented
// on this DB struct; DB will implement the Datastore interface
type DB struct {
  *sql.DB
}


// NewDB initializes a new postgres database connection and attaches
// said connection to our DB struct, which we can then call all of the
// methods described by the "Datastore" inferface
func NewDB(dataSourceName string) (*DB, error) {
  db, err := sql.Open("postgres", dataSourceName)
  if err != nil {
      return nil, err
  }
  if err = db.Ping(); err != nil {
      return nil, err
  }
  return &DB{db}, nil
}
