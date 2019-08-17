package auth


import (
  "golang.org/x/crypto/bcrypt"
)


/*---------------------------------
            Interface
----------------------------------*/

// EncryptionManager describes all of the methods used
// for handling the encryption side of our auth layer
type EncryptionManager interface {
  HashPassword(password string) (string, error)
  CheckPassword(hashed string, password string) (bool, error) 
}


/*---------------------------------
       Method Implementations
----------------------------------*/


// HashPassword takes a given password and hashes it for storage in the user table
func (a *Auth) HashPassword(password string) (string, error) {
  bytePassword := []byte(password)

  hash, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.MinCost)
  if err != nil {
    return "", err
  }

  return string(hash), nil
}


// CheckPassword checks a given plain password against
// a hashed password to see if they're compatible
func (a *Auth) CheckPassword(hashed string, password string) (bool, error) {
  byteHash := []byte(hashed)
  bytePassword := []byte(password)

  err := bcrypt.CompareHashAndPassword(byteHash, bytePassword)
  if err != nil {
    return false, err
  }

  return  true, nil
}
