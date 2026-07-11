// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package client

// UserGroup models

// CreateUserGroup represents the request to create a new user group
type CreateUserGroup struct {
	Name        string        `json:"name"`
	UniqueId    string        `json:"uniqueId,omitempty"`
	Description string        `json:"description,omitempty"`
	Roles       []UserRoleRef `json:"roles,omitempty"`
}

// UpdateUserGroup represents the request to update an existing user group
type UpdateUserGroup struct {
	Name        string        `json:"name,omitempty"`
	Description string        `json:"description,omitempty"`
	Roles       []UserRoleRef `json:"roles,omitempty"`
}

// UserRef represents a user reference
type UserRef struct {
	Id          string `json:"id"`
	LoginName   string `json:"loginName,omitempty"`
	LastName    string `json:"lastName,omitempty"`
	FirstName   string `json:"firstName,omitempty"`
	Email       string `json:"email,omitempty"`
	PhoneNumber string `json:"phoneNumber,omitempty"`
}

// UserGroupResponse represents the API response for a user group
type UserGroupResponse struct {
	UniqueId    string        `json:"uniqueId"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Roles       []UserRoleRef `json:"roles"`
	CreatedTime string        `json:"createdTime"`
	UpdatedTime string        `json:"updatedTime"`
}

type GetUsersFromUserGroupResponse struct {
	Results []UserRef `json:"results"`
}
