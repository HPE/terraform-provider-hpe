// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package client

// User models

// TimeZone represents a timezone object from the API
type TimeZone struct {
	Code  string `json:"code,omitempty"`
	Id    string `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Label string `json:"label,omitempty"`
}

// UserNotification represents a user notification preference
type UserNotification struct {
	NotifyType            string `json:"notifyType"`
	NotifyMethod          string `json:"notifyMethod"`
	NotifyInputType       string `json:"notifyInputType,omitempty"`
	NotifyRecurringReport bool   `json:"notifyRecurringReport"`
}

// TimezoneResponse represents the API response for timezones
type TimezoneResponse struct {
	Code  string `json:"code"`
	Id    string `json:"id"`
	Name  string `json:"name"`
	Label string `json:"label"`
}

// CreateUser represents the request to create a new user
type CreateUser struct {
	LoginName         string             `json:"loginName"`
	Password          string             `json:"password,omitempty"`
	FirstName         string             `json:"firstName"`
	LastName          string             `json:"lastName"`
	Designation       string             `json:"designation,omitempty"`
	Address           string             `json:"address,omitempty"`
	City              string             `json:"city,omitempty"`
	State             string             `json:"state,omitempty"`
	Zip               string             `json:"zip,omitempty"`
	Country           string             `json:"country,omitempty"`
	Email             string             `json:"email"`
	AltEmail          string             `json:"altEmail,omitempty"`
	PhoneNumber       string             `json:"phoneNumber,omitempty"`
	MobileNumber      string             `json:"mobileNumber,omitempty"`
	TimeZone          *TimeZone          `json:"timeZone,omitempty"`
	UserNotifications []UserNotification `json:"userNotifications,omitempty"`
	Roles             []UserRoleRef      `json:"roles,omitempty"`
	UserGroups        []UserGroupRef     `json:"userGroups,omitempty"`
	ChangePassword    bool               `json:"changePassword,omitempty"`
	AuthType          string             `json:"authType,omitempty"`
	OrganizationUnits []OrganizationUnit `json:"organizationUnits,omitempty"`
}

// UpdateUser represents the request to update an existing user
type UpdateUser struct {
	FirstName         string             `json:"firstName,omitempty"`
	LastName          string             `json:"lastName,omitempty"`
	Designation       string             `json:"designation,omitempty"`
	Address           string             `json:"address,omitempty"`
	City              string             `json:"city,omitempty"`
	State             string             `json:"state,omitempty"`
	Zip               string             `json:"zip,omitempty"`
	Country           string             `json:"country,omitempty"`
	Email             string             `json:"email,omitempty"`
	AltEmail          string             `json:"altEmail,omitempty"`
	PhoneNumber       string             `json:"phoneNumber,omitempty"`
	MobileNumber      string             `json:"mobileNumber,omitempty"`
	TimeZone          *TimeZone          `json:"timeZone,omitempty"`
	UserNotifications []UserNotification `json:"userNotifications,omitempty"`
	Roles             []UserRoleRef      `json:"roles,omitempty"`
	UserGroups        []UserGroupRef     `json:"userGroups,omitempty"`
	AuthType          string             `json:"authType,omitempty"`
	OrganizationUnits []OrganizationUnit `json:"organizationUnits,omitempty"`
}

// UserRoleRef represents a role reference for a user
type UserRoleRef struct {
	Id            int    `json:"id,omitempty"`
	UniqueId      string `json:"uniqueId,omitempty"`
	Name          string `json:"name,omitempty"`
	UserGroupRole bool   `json:"userGroupRole,omitempty"`
}

// UserGroupRef represents a user group reference
type UserGroupRef struct {
	Name     string `json:"name,omitempty"`
	UniqueId string `json:"uniqueId,omitempty"`
}

// OrganizationUnit represents an org unit for a user
type OrganizationUnit struct {
	Id       string `json:"id,omitempty"`
	UniqueId string `json:"uniqueId,omitempty"`
	Name     string `json:"name,omitempty"`
}

// UserResponse represents the API response for a user
type UserResponse struct {
	Id                string             `json:"id"`
	LoginName         string             `json:"loginName"`
	FirstName         string             `json:"firstName"`
	LastName          string             `json:"lastName"`
	Designation       string             `json:"designation"`
	Address           string             `json:"address"`
	City              string             `json:"city"`
	State             string             `json:"state"`
	Zip               string             `json:"zip"`
	Country           string             `json:"country"`
	Email             string             `json:"email"`
	AltEmail          string             `json:"altEmail"`
	PhoneNumber       string             `json:"phoneNumber"`
	MobileNumber      string             `json:"mobileNumber"`
	TimeZone          *TimeZone          `json:"timeZone"`
	Roles             []UserRoleRef      `json:"roles"`
	UserGroups        []UserGroupRef     `json:"userGroups"`
	UserNotifications []UserNotification `json:"userNotifications"`
	AuthType          string             `json:"authType"`
	Status            string             `json:"status"`
	CreatedTime       string             `json:"createdTime"`
	UpdatedTime       string             `json:"updatedTime"`
	OrganizationUnits []OrganizationUnit `json:"organizationUnits"`
}

// UserSearchResponse represents the response for user search
type UserSearchResponse struct {
	Results    []UserResponse `json:"results"`
	TotalCount int            `json:"totalCount"`
	PageNo     int            `json:"pageNo"`
	PageSize   int            `json:"pageSize"`
}

// UserMinimalResponse represents a minimal user response
type UserMinimalResponse struct {
	Id        string `json:"id"`
	LoginName string `json:"loginName"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

// ChangePasswordRequest represents the request to change a user's password
type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword,omitempty"`
	NewPassword     string `json:"newPassword"`
}

type DeleteUserRequest struct {
	TerminateReason  string `json:"terminateReason,omitempty"`
	DeactivateReason string `json:"deactivateReason,omitempty"`
	MaskType         string `json:"maskType,omitempty"`
}
