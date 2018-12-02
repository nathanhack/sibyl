package rest

type ErrorState struct {
	Error         string `json:"Errors"`
	ErrorReturned bool   `json:"ErrorReturned"`
}
