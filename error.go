package stormpath

type StormpathError struct {
	Status           int
	Code             int
	Message          string
	DeveloperMessage string
	MoreInfo         string
}
