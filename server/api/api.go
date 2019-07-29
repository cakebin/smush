package api


// Response is a generic response, returned
// when sending POST requests to our API routes
type Response struct {
  Success bool  `json:"success"`
  Error   error `json:"error"`
}
