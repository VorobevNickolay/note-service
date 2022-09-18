package app

import "encoding/json"

type ErrorModel struct {
	Error string `json:"error"`
}

type TokenModel struct {
	Token string `json:"token"`
}

type ValidationErrors struct {
	Errors map[string]string `json:"errors"`
}

func NewValidationErrors() ValidationErrors {
	return ValidationErrors{
		Errors: make(map[string]string, 0),
	}
}

func (ve ValidationErrors) Error() string {
	res, _ := json.Marshal(ve)
	return string(res)
}

var UnknownError = ErrorModel{"Unknown error"}
