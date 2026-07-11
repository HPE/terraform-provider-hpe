// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package client

// PermissionSet models

// CreatePermissionSet represents the request to create a new permission set
type CreatePermissionSet struct {
	Name        string       `json:"name"`
	Description string       `json:"description,omitempty"`
	Permissions []Permission `json:"permissions,omitempty"`
	Scope       string       `json:"scope,omitempty"` // "MSP" or "CLIENT"
}

// UpdatePermissionSet represents the request to update an existing permission set
type UpdatePermissionSet struct {
	Name        string       `json:"name,omitempty"`
	Description string       `json:"description,omitempty"`
	Permissions []Permission `json:"permissions,omitempty"`
}

// Permission represents a single permission
type Permission struct {
	Id   int    `json:"id,omitempty"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// PermissionSetResponse represents the API response for a permission set
type PermissionSetResponse struct {
	Id          int    `json:"id"`
	UniqueId    string `json:"uniqueId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Scope       string `json:"scope,omitempty"` // "MSP" or "CLIENT"

	Permissions []Permission `json:"permissions"`
}

// PermissionSetListResponse represents the response for listing permission sets
type PermissionSetListResponse struct {
	Results    []PermissionSetResponse `json:"results"`
	TotalCount int                     `json:"totalCount"`
	PageNo     int                     `json:"pageNo"`
	PageSize   int                     `json:"pageSize"`
}
