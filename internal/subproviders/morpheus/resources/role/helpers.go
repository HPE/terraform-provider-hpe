package role

import (
	"fmt"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

// For now, keeping this as a separate struct as if we try to unmarshal the permission set
// from the TF config into an sdk.AddRolesRequestRole struct it will fail and throw an error
// due to the missing required Authority field

type rolePermissionSetPost struct {
	// Set the access level for the specified permissions.
	Permissions []sdk.AddRolesRequestRolePermissionsInner `json:"permissions,omitempty"`
	// Set the default access level for for groups (sites). Only applies to user roles.
	GlobalSiteAccess string `json:"globalSiteAccess,omitempty"`
	// Set the access level for the specified groups (sites). Only applies to user roles.
	Sites []sdk.AddRolesRequestRoleSitesInner `json:"sites,omitempty"`
	// Set the default access level for for clouds (zones). Only applies to base account (tenant) roles.
	GlobalZoneAccess string `json:"globalZoneAccess,omitempty"`
	// Set the access level for the specified clouds (zones). Only applies to base account (tenant) roles.
	Zones []sdk.AddRolesRequestRoleZonesInner `json:"zones,omitempty"`
	// Set the default access level for for instance types
	GlobalInstanceTypeAccess string `json:"globalInstanceTypeAccess,omitempty"`
	// Set the access level for the specified instance types
	InstanceTypes []sdk.AddRolesRequestRoleInstanceTypesInner `json:"instanceTypes,omitempty"`
	// Set the default access level for blueprints
	GlobalAppTemplateAccess string `json:"globalAppTemplateAccess,omitempty"`
	// Set the access level for the specified blueprints (appTemplates)
	AppTemplates []sdk.AddRolesRequestRoleAppTemplatesInner `json:"appTemplates,omitempty"`
	// Set the default access level for catalog item types
	GlobalCatalogItemTypeAccess string `json:"globalCatalogItemTypeAccess,omitempty"`
	// Set the access level for the specified catalog item types
	CatalogItemTypes []sdk.AddRolesRequestRoleCatalogItemTypesInner `json:"catalogItemTypes,omitempty"`
	// Set the default access level for personas
	GlobalPersonaAccess string `json:"globalPersonaAccess,omitempty"`
	// Set the access level for the specified personas
	Personas []sdk.AddRolesRequestRolePersonasInner `json:"personas,omitempty"`
	// Set the default access level for VDI pools
	GlobalVdiPoolAccess string `json:"globalVdiPoolAccess,omitempty"`
	// Set the access level for the specified VDI pools
	VdiPools []sdk.AddRolesRequestRoleVdiPoolsInner `json:"vdiPools,omitempty"`
	// Set the default access level for report types
	GlobalReportTypeAccess string `json:"globalReportTypeAccess,omitempty"`
	// Set the access level for the specified report types
	ReportTypes []sdk.AddRolesRequestRoleReportTypesInner `json:"reportTypes,omitempty"`
	// Set the default access level for tasks
	GlobalTaskAccess string `json:"globalTaskAccess,omitempty"`
	// Set the access level for the specified tasks
	Tasks []sdk.AddRolesRequestRoleTasksInner `json:"tasks,omitempty"`
	// Set the default access level for workflows (taskSets)
	GlobalTaskSetAccess string `json:"globalTaskSetAccess,omitempty"`
	// Set the access level for the specified workflows (taskSets)
	TaskSets []sdk.AddRolesRequestRoleTaskSetsInner `json:"taskSets,omitempty"`
}

type rolePermissionSetGet struct {
	// support unmarshalling on "permissions" key for deep compare
	Permissions []sdk.AddRoles200ResponseAllOfFeaturePermissionsInner `json:"permissions,omitempty"`

	// the original fields from get
	FeaturePermissions          []sdk.AddRoles200ResponseAllOfFeaturePermissionsInner      `json:"featurePermissions,omitempty"`
	GlobalSiteAccess            *string                                                    `json:"globalSiteAccess,omitempty"`
	Sites                       []sdk.AddRoles200ResponseAllOfSitesInner                   `json:"sites,omitempty"`
	GlobalZoneAccess            *string                                                    `json:"globalZoneAccess,omitempty"`
	Zones                       []sdk.AddRoles200ResponseAllOfSitesInner                   `json:"zones,omitempty"`
	GlobalInstanceTypeAccess    *string                                                    `json:"globalInstanceTypeAccess,omitempty"`
	InstanceTypePermissions     []sdk.AddRoles200ResponseAllOfInstanceTypePermissionsInner `json:"instanceTypePermissions,omitempty"`
	GlobalAppTemplateAccess     *string                                                    `json:"globalAppTemplateAccess,omitempty"`
	AppTemplatePermissions      []sdk.AddRoles200ResponseAllOfAppTemplatePermissionsInner  `json:"appTemplatePermissions,omitempty"`
	GlobalCatalogItemTypeAccess *string                                                    `json:"globalCatalogItemTypeAccess,omitempty"`
	CatalogItemTypePermissions  []sdk.AddRoles200ResponseAllOfSitesInner                   `json:"catalogItemTypePermissions,omitempty"`
	GlobalPersonaAccess         *string                                                    `json:"globalPersonaAccess,omitempty"`
	PersonaPermissions          []sdk.AddRoles200ResponseAllOfInstanceTypePermissionsInner `json:"personaPermissions,omitempty"`
	GlobalVdiPoolAccess         *string                                                    `json:"globalVdiPoolAccess,omitempty"`
	VdiPoolPermissions          []sdk.AddRoles200ResponseAllOfSitesInner                   `json:"vdiPoolPermissions,omitempty"`
	GlobalReportTypeAccess      *string                                                    `json:"globalReportTypeAccess,omitempty"`
	ReportTypePermissions       []sdk.AddRoles200ResponseAllOfInstanceTypePermissionsInner `json:"reportTypePermissions,omitempty"`
	GlobalTaskAccess            *string                                                    `json:"globalTaskAccess,omitempty"`
	TaskPermissions             []sdk.AddRoles200ResponseAllOfAppTemplatePermissionsInner  `json:"taskPermissions,omitempty"`
	GlobalTaskSetAccess         *string                                                    `json:"globalTaskSetAccess,omitempty"`
	TaskSetPermissions          []sdk.AddRoles200ResponseAllOfAppTemplatePermissionsInner  `json:"taskSetPermissions,omitempty"`
}

// A struct used to unmarshal the equivalent types from the generated SDK models
// Contains naming of fields in line with the Morpheus UI and documentation,
// not the internal API naming used by the OpenAPI spec and generated SDK.
// type rolePermissionSet struct {
// 	// Feature Permissions
// 	Features sdk.AddRolesRequestRolePermissionsInner `json:"featurePermissions,omitempty"`
// 	// Sets the Groups Default Access option
// 	GroupsDefaultAccess *string `json:"globalSiteAccess,omitempty"`
//
// 	// Groups Permissions
// 	Groups sdk.AddRolesRequestRoleSitesInner
// 	// Sets the Clouds Default Access option
// 	// CloudsDefaultAccess
// 	// // Clouds Permissions
// 	// Clouds
// 	//
// 	// // Sets the Instance Types Default Access option
// 	// InstanceTypesDefaultAccess
// 	//
// 	// // Instance Types Permissions
// 	// InstanceTypes
// 	// // Sets the Blueprints Default Access option
// 	// BlueprintsDefaultAccess
// 	// // Blueprints Permissions
// 	// Blueprints
// 	// // Sets the Catalog Item Types Default Access option
// 	// CatalogItemTypesDefaultAccess
// 	// // Catalog Item Types Permissions
// 	// CatalogItemTypes
// 	//
// 	// // Sets the Personas Default Access option
// 	// PersonasDefaultAccess
// 	// // Personas Permissions
// 	// Personas
// 	// // Sets the Vdi Pools Default Access option
// 	// VdiPoolsDefaultAccess
// 	// // VdiPools Permissions
// 	// VdiPools
// 	// // Sets the Report Types Default Access option
// 	// ReportTypesDefaultAccess
// 	// // Report Types Permissions
// 	// ReportTypes
// 	//
// 	// // Sets the Tasks Default Access option
// 	// TasksDefaultAccess
// 	// // Tasks Permissions
// 	// Tasks
// 	// // Sets the Workflows Default Access option
// 	// WorkflowsDefaultAccess
// 	// // Workflows Permissions
// 	// Workflows
// 	//
// 	// // TODO: Implement cluster types support into generated SDK
// 	// // (requires update of OpenAPI spec)
// 	// // // Sets the Cluster Types Default Access option
// 	// // ClusterTypesDefaultAccess
// 	// // // Cluster Types Permissions
// 	// // ClusterTypes
// }

func permissionSetConfigPermissions() string {
	return `
{
  "permissions": [
    {
      "code": "integrations-ansible",
      "access": "full"
    }
  ]
}
`
}

func permissionSetConfigSites(siteId int64) string {
	return fmt.Sprintf(`
"sites": [
  {
    "id": %d,
    "access": "full"
  }
]
`, siteId)
}
