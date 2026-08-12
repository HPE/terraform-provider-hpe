// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

//nolint:lll
package errors

import "strings"

// BaseError is an error type that all other error types embed.
type BaseError struct {
	ErrorResponse ErrorResponse
	Info          string
	OriginalError error
}

// ErrorResponse should be used to return details of a problem
//
//nolint:maligned
type ErrorResponse struct {
	// message: clear and concise description of the error condition, suitable for display to an end user
	Message string `json:"message,omitempty"`
	// details: optional verbose description of the error condition, suitable for display to an end user
	Details string `json:"details,omitempty"`
	// recommendedActions: steps that an end user can perform to correct the error condition (list of strings)
	RecommendedActions []string `json:"recommendedActions,omitempty"`
	// nestedErrors: list of subsidiary errors that led to this error condition, each of which is a nested set of these same attributes
	NestedErrors string `json:"nestedErrors,omitempty"`
	// errorSource: identifies some element of the request that caused the error, e.g. a specific form field or a specific resource
	ErrorSource bool `json:"errorSource,omitempty"`
	// errorCode: an opaque string uniquely identifying the error for programmatic use
	ErrorCode string `json:"errorCode,omitempty"`
	// data: arbitrary data associated with the error condition, for programmatic use
	Data string `json:"data,omitempty"`
	// canForce: a Boolean indicating whether this error condition can be ignored by re-sending the request with the force=true query parameter as described below.
	CanForce bool `json:"canForce,omitempty"`
}

//nolint:gocritic
func (e *BaseError) Error() string {
	if e.ErrorResponse.Message != "" {
		return e.ErrorResponse.Message
	} else if e.Info != "" {
		return e.Info
	} else if e.OriginalError != nil {
		return e.OriginalError.Error()
	}

	return "An error occurred."
}

type ErrBadRequest struct {
	BaseError
}

// MakeErrBadRequest helper to create ErrBadRequest
func MakeErrBadRequest(errorResponse ErrorResponse) *ErrBadRequest {
	return &ErrBadRequest{BaseError{ErrorResponse: errorResponse}}
}

// ErrForbidden is error type that can be returned from a function and propogated to be handled appropriately
// Used to indicate insufficient access
type ErrForbidden struct {
	BaseError
	ForbiddenThings []string
}

// MakeErrForbidden helper to create ErrForbidden
func MakeErrForbidden(forbiddenThings ...string) *ErrForbidden {
	return &ErrForbidden{
		ForbiddenThings: forbiddenThings,
		BaseError:       BaseError{Info: "Forbidden: " + strings.Join(forbiddenThings, ", ")},
	}
}

// ErrUnauthorized is a error type that can be returned from a function and propagated up to be handled appropriately.
// Used to indicate that there was a conflict.
// see SetResponseIfError
type ErrUnauthorized struct {
	BaseError
	UnauthorizedReason string
}

// MakeErrUnauthorized helper to create ErrUnauthorized
func MakeErrUnauthorized(reason string) *ErrUnauthorized {
	return &ErrUnauthorized{
		UnauthorizedReason: reason,
		BaseError:          BaseError{Info: "Unauthorized access: " + reason},
	}
}

// ErrInternalError is an error type that can be returned from a
// function and propagated up to be handled appropriately. Used to indicate
// that something went wrong. See SetInternalErrorWithErrorResponse.
type ErrInternalError struct {
	BaseError
}

// MakeErrInternalError  helper to create ErrInternalErrorDetails
func MakeErrInternalError(errorResponse ErrorResponse) *ErrInternalError {
	return &ErrInternalError{BaseError{ErrorResponse: errorResponse}}
}
