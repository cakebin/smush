package auth


// Auth is the struct that we're going to use 
// to implement all of out Authenticator interfaces
type Auth struct {}


// Authenticator combines all of the various
// aspects of our auth layer into one
type Authenticator interface {
  JWTManager
  EncryptionManager
  RoleManager
}


// New makes a new Auth struct which implements
// all of the "Authenticator" methods
func New() *Auth {
  return &Auth{}
}
