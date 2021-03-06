package db


/*---------------------------------
            Interface
----------------------------------*/

// UserManager describes all of the methods used
// to interact with the users table in our database
type UserManager interface {
  GetAllUsers() ([]*User, error)
  GetUserIDByEmail(email string) (int64, error)
  GetUserResetPasswordTokenByUserID(int64) (string, error)

  UpdateUserProfile(profileUpdate *UserProfileUpdate) (int64, error)
  UpdateUserRefreshToken(refreshUpdate *UserRefreshUpdate) (int64, error)
  UpdateUserResetPasswordToken(resetPasswordUpdate *UserResetPasswordUpdate) (int64, error)
  UpdateUserHashedPassword(hashedPasswordUpdate *UserHashedPasswordUpdate) (int64, error)
  UpdateUserDefaultUserCharacter(userCharUpdate *UserDefaultUserCharacterUpdate) (int64, error)

  CreateUser(userCreate *UserCreate) (int64, error)
}


/*---------------------------------
          Data Structures
----------------------------------*/

// User describes the data in the users table
type User struct {
  UserID   int64   `json:"userId"`
  UserName string  `json:"userName"`
}


// UserProfileUpdate describes the data needed
// to update a given user's profile information
type UserProfileUpdate struct {
  UserID    int64   `json:"userId"`
  UserName  string  `json:"userName"`
}


// UserDefaultUserCharacterUpdate describes the data needed
// to update a given user's default user character
type UserDefaultUserCharacterUpdate struct {
  UserID           int64          `json:"userId"`
  UserCharacterID  NullInt64JSON  `json:"userCharacterId"`
}


// UserRefreshUpdate describes the data
// needed to update a given users refresh token
type UserRefreshUpdate struct {
  UserID        int64   `json:"userId"`
  RefreshToken  string  `json:"refreshToken"`
}


// UserResetPasswordUpdate describes the data
// needed to update a given user's reset password token
type UserResetPasswordUpdate struct {
  UserID              int64   `json:"userId"`  
  ResetPasswordToken  string  `json:"resetPasswordToken"`
}


// UserHashedPasswordUpdate describes the data
// needed to update a given user's hashed password
type UserHashedPasswordUpdate struct {
  UserID          int64   `json:"userId"`
  HashedPassword  string  `json:"hashedPassword"`
}

// UserCreate describes the data needed
// to create a new user in our db
type UserCreate struct {
  UserName        string  `json:"userName"`
  EmailAddress    string  `json:"emailAddress"`
  HashedPassword  string  `json:"hashedPassword"`
}


/*---------------------------------
       Method Implementations
----------------------------------*/

// GetAllUsers fetches userId/userName for all users
func (db *DB) GetAllUsers() ([]*User, error) {
  sqlStatement := `
    SELECT
      user_id,
      user_name
    FROM
      users
  `
  rows, err := db.Query(sqlStatement)
  if err != nil {
    return nil, err
  }
  defer rows.Close()

  users := make([]*User, 0)
  for rows.Next() {
    user := new(User)
    err := rows.Scan(
      &user.UserID,
      &user.UserName,
    )
    if err != nil {
      return nil, err
    }

    users = append(users, user)
  }

  err = rows.Err()
  if err != nil {
    return nil, err
  }

  return users, nil
}


// GetUserIDByEmail gets a specific user's id from the users table by email
func (db *DB) GetUserIDByEmail(email string) (int64, error) {
  var userID int64
  sqlStatement := `
    SELECT
      user_id
    FROM
      users
    WHERE
      email_address = $1
  `
  row := db.QueryRow(sqlStatement, email)
  err := row.Scan(&userID)

  if err != nil {
    return 0, err
  }

  return userID, nil
}


// GetUserResetPasswordTokenByUserID gets a user's reset_password_token; used for "forgot password"
func (db *DB) GetUserResetPasswordTokenByUserID(userID int64) (string, error) {
  var resetPasswordToken string
  sqlStatement := `
    SELECT
      reset_password_token
    FROM
      users
    WHERE
      user_id = $1
  `
  row := db.QueryRow(
    sqlStatement,
    userID,
  )
  err := row.Scan(&resetPasswordToken)
  if err != nil {
    return "", err
  }

  return resetPasswordToken, nil
}


// CreateUser adds a new entry to the users table in our database
func (db *DB) CreateUser(userCreate *UserCreate) (int64, error) {
  var userID int64
  sqlStatement := `
    INSERT INTO users
      (user_name, email_address, hashed_password)
    VALUES
      ($1, $2, $3)
    RETURNING
      user_id
  `
  row := db.QueryRow(
    sqlStatement,
    userCreate.UserName,
    userCreate.EmailAddress,
    userCreate.HashedPassword,
  )
  err := row.Scan(&userID)

  if err != nil {
    return 0, err
  }

  return userID, nil
}


// UpdateUserProfile updates an entry in the users table with the given data
func (db *DB) UpdateUserProfile(profileUpdate *UserProfileUpdate) (int64, error) {
  var userID int64
  sqlStatement := `
    UPDATE
      users
    SET
      user_name = $1
    WHERE
      user_id = $2
    RETURNING
      user_id
  `
  row := db.QueryRow(
    sqlStatement,
    profileUpdate.UserName,
    profileUpdate.UserID,
  )
  err := row.Scan(&userID)

  if err != nil {
    return 0, err
  }

  return userID, nil

}


// UpdateUserRefreshToken updates a user's refresh token; used for auth
func (db *DB) UpdateUserRefreshToken(refreshUpdate *UserRefreshUpdate) (int64, error) {
  var userID int64
  sqlStatement := `
    UPDATE
      users
    SET
      refresh_token = $1
    WHERE
      user_id = $2
    RETURNING
      user_id
  `
  row := db.QueryRow(
    sqlStatement,
    refreshUpdate.RefreshToken,
    refreshUpdate.UserID,
  )
  err := row.Scan(&userID)

  if err != nil {
    return 0, err
  }

  return userID, nil
}


// UpdateUserResetPasswordToken updates a user's reset password token; used for "forgot password"
func (db *DB) UpdateUserResetPasswordToken(resetPasswordUpdate *UserResetPasswordUpdate) (int64, error) {
  var userID int64
  sqlStatement := `
    UPDATE
      users
    SET
      reset_password_token = $1
    WHERE
      user_id = $2
    RETURNING
      user_id
  `
  row := db.QueryRow(
    sqlStatement,
    resetPasswordUpdate.ResetPasswordToken,
    resetPasswordUpdate.UserID,
  )
  err := row.Scan(&userID)
  if err != nil {
    return 0, err
  }

  return userID, nil
}


// UpdateUserHashedPassword updates a user's hashed password
func (db *DB) UpdateUserHashedPassword(hashedPasswordUpdate *UserHashedPasswordUpdate) (int64, error) {
  var userID int64
  sqlStatement := `
    UPDATE
      users
    SET
      hashed_password = $1
    WHERE
      user_id = $2
    RETURNING
      user_id
  `
  row := db.QueryRow(
    sqlStatement,
    hashedPasswordUpdate.HashedPassword,
    hashedPasswordUpdate.UserID,
  )
  err := row.Scan(&userID)
  if err != nil {
    return 0, err
  }

  return userID, nil
}


// UpdateUserDefaultUserCharacter updates a user's default user character
func (db *DB) UpdateUserDefaultUserCharacter(userCharUpdate *UserDefaultUserCharacterUpdate) (int64, error) {
  var userID int64
  sqlStatement := `
    UPDATE
      users
    SET
      default_user_character_id = $1
    WHERE
      user_id = $2
    RETURNING
      user_id
  `
  row := db.QueryRow(
    sqlStatement,
    userCharUpdate.UserCharacterID,
    userCharUpdate.UserID,
  )
  err := row.Scan(&userID)

  if err != nil {
    return 0, err
  }

  return userID, nil
}
