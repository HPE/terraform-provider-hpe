// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package models

// ErrBadRequest
type ErrBadRequest struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors,omitempty"`
}

// ErrInternalError
type ErrInternalError struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ErrNotFound
type ErrNotFound struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ErrUnauthorized
type ErrUnauthorized struct {
	Errors  string `json:"error"`
	Message string `json:"message"`
}

// Success Or Failure Message
type SuccessOrErrorMessage struct {
	IDModel
	Success bool        `json:"success,omitempty"`
	Msg     string      `json:"msg,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}
