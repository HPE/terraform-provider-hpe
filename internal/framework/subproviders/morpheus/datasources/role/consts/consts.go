// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package consts

const (
	ErrorNoRoleFound        = `no role found`
	ErrorNoValidSearchTerms = `no valid search terms - an id or name is required`
	ErrorRunningPreApply    = `Error running pre-apply plan: exit status 1`
	ErrorMultipleRoles      = `multiple roles were returned`
	RoleTypeUser            = "user"
	// Provider-specific role type created to adopt new Morpheus naming conventions.
	RoleTypeTenant = "tenant"
	// API-specific legacy role type which we use to maintain API compatibility.
	// Account Roles and Tenant Roles are the same thing; Tenant Role is the newer name.
	RoleTypeAccountAPI = "account"
)
