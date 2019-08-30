package auth


import (
  "errors"
  "os"
  "time"

  "github.com/dgrijalva/jwt-go"
)


/*---------------------------------
          Data Structures
----------------------------------*/

var jwtKey = []byte(os.Getenv("JWT_TOKEN_SECRET"))


// Claims is a custom extended jwt.StandardClaims to include
// a user's email address as part of the claims
type Claims struct {
  UserID               int
  jwt.StandardClaims
}


/*---------------------------------
            Interface
----------------------------------*/

// JWTManager describes all of the methods used
// for handling the JSON web token side of our auth layer
type JWTManager interface {
  GetNewJWTToken(id int64, expiration time.Time) (string, error)
  RefreshJWTAccessToken(token string, expiration time.Time) (string, error)
  CheckJWTToken(token string) (bool, error)
  GetUserIDFromJWTToken(token string) (int64, error)
}


/*---------------------------------
       Method Implementations
----------------------------------*/

// GetNewJWTToken generates a new jwt access token
// for a given user with a given expiration date
func (a *Auth) GetNewJWTToken(id int64, expirationTime time.Time) (string, error) {
  claims := &Claims{
    UserID: int(id),
    StandardClaims: jwt.StandardClaims{
      // In JWT, the expiry time is expressed as unix milliseconds
      ExpiresAt: expirationTime.Unix(),
    },
  }

  token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
  tokenStr, err := token.SignedString(jwtKey)
  if err != nil {
    return "", err
  }

  return tokenStr, nil
}


// CheckJWTToken takes a jwt token string and checks
// if is still valid (i.e. didn't expire, of the right sign format, etc)
func (a *Auth) CheckJWTToken(token string) (bool, error) {
  claims := &Claims{}
  parsedToken, err := jwt.ParseWithClaims(
    token,
    claims,
    func(token *jwt.Token) (interface{}, error) { return jwtKey, nil },
  )
  if err != nil {
    return false, err
  }
  
  if !parsedToken.Valid {
    return false, errors.New("Token Expired")
  }

  return true, nil
}


// RefreshJWTAccessToken takes an existing access jwt token (likely expired)
// And updates its expiration time so that the token is still valid for future use
func (a *Auth) RefreshJWTAccessToken(token string, newExpiration time.Time) (string, error) {
  claims := &Claims{}
  _, err := jwt.ParseWithClaims(
    token,
    claims,
    func(token *jwt.Token) (interface{}, error) { return jwtKey, nil },
  )
  if err != nil {
    return "", err
  }
  claims.ExpiresAt = newExpiration.Unix()
  newAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
  newAccessTokenStr, err := newAccessToken.SignedString(jwtKey)
  if err != nil {
    return "", err
  }

  return newAccessTokenStr, nil
}


// GetUserIDFromJWTToken extracts the stored userID fromt the claims in a JWT token 
func (a *Auth) GetUserIDFromJWTToken(token string) (int64, error) {
  claims := new(Claims)
  parsedToken, err := jwt.ParseWithClaims(
    token,
    claims,
    func(token *jwt.Token) (interface{}, error) { return jwtKey, nil},
  )
  if err != nil {
    return 0, err
  }

  if !parsedToken.Valid {
    return 0, errors.New("Token Expired")
  }

  return int64(claims.UserID), nil
}
