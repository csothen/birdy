package birdy

// GenericError is a generic error message returned by a server
type GenericError struct {
	Message string `json:"message"`
}
