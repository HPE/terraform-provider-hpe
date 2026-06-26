// Defines types that are commonly used in api request and responses
package morpheus

import (
	_ "fmt"
)

// Common response types
// Common types used in all responses
// todo: moves these to shared_result_types.go or something

// StandardResult is a response format for most actions
type StandardResult struct {
	Success bool              `json:"success"`
	Message string            `json:"msg"`
	Errors  map[string]string `json:"errors"`
}

// StandardErrorResult is a format for request errors eg. http 400, 401, 500
type StandardErrorResult struct {
	StandardResult
}

// A standard format for Delete actions
type DeleteResult struct {
	StandardResult
}

// Common request types
type GetByIdRequest struct {
	Request
	ID int64 `json:"id"`
}

// tenantAbbrev is a response format for  describing a list of objects returned.
// Could maybe replace use of this with Account for this, it can unmarshal just id and name
type TenantAbbrev struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// MetaResult is a response format for  describing a list of objects returned.
// This is present in most list results.
type MetaResult struct {
	Total  int64       `json:"total"`
	Size   int64       `json:"size"`
	Max    interface{} `json:"max"`
	Offset int64       `json:"offset"`
}
