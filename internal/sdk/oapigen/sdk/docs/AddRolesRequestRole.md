# AddRolesRequestRole

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Authority** | **string** | Authority (Name) | 
**Description** | Pointer to **NullableString** | Description | [optional] 
**LandingUrl** | Pointer to **NullableString** | An optional override for the default landing page after login for a user. | [optional] 
**RoleType** | Pointer to **string** | Role type | [optional] [default to "user"]
**BaseRoleId** | Pointer to **int64** | Base Role ID. Create the new role with the same permissions and access levels that the specified base role has. If this is not passed, the role is create without any permissions. | [optional] 
**Multitenant** | Pointer to **bool** | Multitenant roles are copied to all tenant accounts and kept in sync until a sub-tenant user modifies their copy of the role. *Only available to master tenant, only applies to user roles* | [optional] [default to false]
**MultitenantLocked** | Pointer to **bool** | Multitenant Locked, prevents sub-tenant users from modifying their copy of multitenant roles. *Only available to master tenant, only applies to user roles* | [optional] [default to false]
**DefaultPersona** | Pointer to **NullableString** |  | [optional] 
**FeaturePermissions** | Pointer to [**[]AddRolesRequestRoleFeaturePermissionsInner**](AddRolesRequestRoleFeaturePermissionsInner.md) | Set the access level for the specified permissions. | [optional] 
**GlobalSiteAccess** | Pointer to **string** | Set the default access level for for groups (sites). Only applies to user roles. | [optional] 
**Sites** | Pointer to [**[]AddRolesRequestRoleSitesInner**](AddRolesRequestRoleSitesInner.md) | Set the access level for the specified groups (sites). Only applies to user roles. | [optional] 
**GlobalZoneAccess** | Pointer to **string** | Set the default access level for for clouds (zones). Only applies to base account (tenant) roles. | [optional] 
**Zones** | Pointer to [**[]AddRolesRequestRoleZonesInner**](AddRolesRequestRoleZonesInner.md) | Set the access level for the specified clouds (zones). Only applies to base account (tenant) roles. | [optional] 
**GlobalInstanceTypeAccess** | Pointer to **string** | Set the default access level for for instance types | [optional] 
**InstanceTypePermissions** | Pointer to [**[]AddRolesRequestRoleInstanceTypePermissionsInner**](AddRolesRequestRoleInstanceTypePermissionsInner.md) | Set the access level for the specified instance types | [optional] 
**GlobalAppTemplateAccess** | Pointer to **string** | Set the default access level for blueprints | [optional] 
**AppTemplatePermissions** | Pointer to [**[]AddRolesRequestRoleAppTemplatePermissionsInner**](AddRolesRequestRoleAppTemplatePermissionsInner.md) | Set the access level for the specified blueprints (appTemplates) | [optional] 
**GlobalCatalogItemTypeAccess** | Pointer to **string** | Set the default access level for catalog item types | [optional] 
**CatalogItemTypePermissions** | Pointer to [**[]AddRolesRequestRoleCatalogItemTypePermissionsInner**](AddRolesRequestRoleCatalogItemTypePermissionsInner.md) | Set the access level for the specified catalog item types | [optional] 
**GlobalPersonaAccess** | Pointer to **string** | Set the default access level for personas | [optional] 
**PersonaPermissions** | Pointer to [**[]AddRolesRequestRolePersonaPermissionsInner**](AddRolesRequestRolePersonaPermissionsInner.md) | Set the access level for the specified personas | [optional] 
**GlobalVdiPoolAccess** | Pointer to **string** | Set the default access level for VDI pools | [optional] 
**VdiPoolPermissions** | Pointer to [**[]AddRolesRequestRoleVdiPoolPermissionsInner**](AddRolesRequestRoleVdiPoolPermissionsInner.md) | Set the access level for the specified VDI pools | [optional] 
**GlobalReportTypeAccess** | Pointer to **string** | Set the default access level for report types | [optional] 
**ReportTypePermissions** | Pointer to [**[]AddRolesRequestRoleReportTypePermissionsInner**](AddRolesRequestRoleReportTypePermissionsInner.md) | Set the access level for the specified report types | [optional] 
**GlobalTaskAccess** | Pointer to **string** | Set the default access level for tasks | [optional] 
**TaskPermissions** | Pointer to [**[]AddRolesRequestRoleTaskPermissionsInner**](AddRolesRequestRoleTaskPermissionsInner.md) | Set the access level for the specified tasks | [optional] 
**GlobalTaskSetAccess** | Pointer to **string** | Set the default access level for workflows (taskSets) | [optional] 
**TaskSetPermissions** | Pointer to [**[]AddRolesRequestRoleTaskSetPermissionsInner**](AddRolesRequestRoleTaskSetPermissionsInner.md) | Set the access level for the specified workflows (taskSets) | [optional] 
**GlobalClusterTypeAccess** | Pointer to **string** | Set the default access level for cluster types | [optional] 
**ClusterTypePermissions** | Pointer to [**[]AddRolesRequestRoleClusterTypePermissionsInner**](AddRolesRequestRoleClusterTypePermissionsInner.md) | Set the access level for the specified cluster types | [optional] 
**Owner** | Pointer to **int64** | Set the role owner (tenant) by ID. *Only available to master tenant* | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddRolesRequestRole{
    // Set fields directly
}
```

### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### LandingUrl (Nullable)

Use the Nullable wrapper methods:
- `obj.LandingUrl.IsSet()` — check if set
- `obj.LandingUrl.Get()` — get the inner value (returns pointer)
- `obj.LandingUrl.Set(&val)` — set the value
- `obj.LandingUrl.Unset()` — clear the value
### DefaultPersona (Nullable)

Use the Nullable wrapper methods:
- `obj.DefaultPersona.IsSet()` — check if set
- `obj.DefaultPersona.Get()` — get the inner value (returns pointer)
- `obj.DefaultPersona.Set(&val)` — set the value
- `obj.DefaultPersona.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


