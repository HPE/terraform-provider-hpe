// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package client

// Role models

// Role represents a role for creation, update, and retrieval
type Role struct {
	Id       int    `json:"id,omitempty"`
	UniqueId string `json:"uniqueId,omitempty"`

	Name        string `json:"name"`
	Description string `json:"description"`

	DefaultRole bool `json:"defaultRole"`

	Clients    []RoleClientRef `json:"clients,omitempty"`
	AllClients bool            `json:"allClients,omitempty"`
	Scope      string          `json:"scope,omitempty"` // ej. MSP, CLIENT

	AllDevices     bool `json:"allDevices"`
	AllCredentials bool `json:"allCredentials"`
	AllAuthzTags   bool `json:"allAuthzTags"`

	Users          []RoleUserRef        `json:"users"`
	UserGroups     []RoleUserGroupRef   `json:"userGroups"`
	Devices        []RoleDeviceRef      `json:"devices"`
	DeviceGroups   []RoleDeviceGroupRef `json:"deviceGroups"`
	Resources      []RoleDeviceRef      `json:"resources"`
	CredentialSets []RoleCredentialRef  `json:"credentialSets"`
	Permissions    []RolePermissionRef  `json:"permissions"`

	CreatedTime string `json:"createdTime,omitempty"`
	UpdatedTime string `json:"updatedTime,omitempty"`
}

// PermissionSetRef represents a permission set reference
type PermissionSetRef struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// RoleClientRef represents a client reference in a role
type RoleClientRef struct {
	UniqueId string `json:"uniqueId"`
	Name     string `json:"name,omitempty"`
}

// RoleUserRef represents a user reference in a role
type RoleUserRef struct {
	Id        string `json:"id,omitempty"`
	LoginName string `json:"loginName,omitempty"`
}

// RoleUserGroupRef represents a user group reference in a role
type RoleUserGroupRef struct {
	UniqueId string `json:"uniqueId,omitempty"`
	Name     string `json:"name,omitempty"`
}

// RoleDeviceRef represents a device reference in a role
type RoleDeviceRef struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// RoleDeviceGroupRef represents a device group reference in a role
type RoleDeviceGroupRef struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// RoleCredentialRef represents a credential set reference in a role
type RoleCredentialRef struct {
	UniqueId string `json:"uniqueId,omitempty"`
	Name     string `json:"name,omitempty"`
}

// RolePermissionRef represents a permission reference in a role
type RolePermissionRef struct {
	Id   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// RoleSearchResponse represents the response for role search
type RoleSearchResponse struct {
	Results      []Role `json:"results"`
	TotalResults int    `json:"totalResults"`
	PageNo       int    `json:"pageNo"`
	PageSize     int    `json:"pageSize"`
	TotalPages   int    `json:"totalPages"`
	NextPage     bool   `json:"nextPage"`
}
