package role

const (
	RoleTypeUser = "user"
	// Provider-specific role type created to adopt new Morpheus naming conventions.
	RoleTypeTenant = "tenant"
	// API-specific legacy role type which we use to maintain API compatibility.
	// Account Roles and Tenant Roles are the same thing; Tenant Role is the newer name.
	RoleTypeAccountAPI = "account"

	// One of the access levels available for non-feature, fine-grained permissions.
	// When set to 'default', the resources for which access is being set
	// will inherit the access level specified by the "Default Access Level" field.
	DefaultPermissionAccessLevel = "default"
)
