package email


// Email is the struct to use to implement
// all of the Emailer interfaces
type Email struct {}


// Emailer combines all of the various
// aspecs of our email layer into one
type Emailer interface {
  PasswordEmailer
}


// New makes a new Email struct which
// implements all of the "Emailer" methods
func New() *Email {
  return &Email{}
}
