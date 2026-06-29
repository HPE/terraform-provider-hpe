# UpdateRoleRequestRole

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Authority** | Pointer to **string** | Authority (Name) | [optional] 
**Description** | Pointer to **NullableString** | Description | [optional] 
**LandingUrl** | Pointer to **NullableString** | An optional override for the default landing page after login for a user. | [optional] 
**Multitenant** | Pointer to **bool** | Multitenant roles are copied to all tenant accounts and kept in sync until a sub-tenant user modifies their copy of the role. *Only available to master tenant, only applies to user roles* | [optional] 
**MultitenantLocked** | Pointer to **bool** | Multitenant Locked, prevents sub-tenant users from modifying their copy of multitenant roles. *Only available to master tenant, only applies to user roles* | [optional] 
**DefaultPersona** | Pointer to **NullableString** |  | [optional] 
**ResetPermissions** | Pointer to **bool** | Resets access levels for all feature permissions to their default values. | [optional] 
**ResetAllAccess** | Pointer to **bool** | Resets access levels for all permissions (including feature, non-feature, and default access levels) to their default values. | [optional] 
**FeaturePermissions** | Pointer to [**[]UpdateRoleRequestRoleFeaturePermissionsInner**](UpdateRoleRequestRoleFeaturePermissionsInner.md) | Set the access level for the specified permissions. | [optional] 
**GlobalSiteAccess** | Pointer to **string** | Set the default access level for for groups (sites). Only applies to user roles. | [optional] 
**Sites** | Pointer to [**[]UpdateRoleRequestRoleSitesInner**](UpdateRoleRequestRoleSitesInner.md) | Set the access level for the specified groups (sites). Only applies to user roles. | [optional] 
**GlobalZoneAccess** | Pointer to **string** | Set the default access level for for clouds (zones). Only applies to base account (tenant) roles. | [optional] 
**Zones** | Pointer to [**[]UpdateRoleRequestRoleZonesInner**](UpdateRoleRequestRoleZonesInner.md) | Set the access level for the specified clouds (zones). Only applies to base account (tenant) roles. | [optional] 
**GlobalInstanceTypeAccess** | Pointer to **string** | Set the default access level for for instance types | [optional] 
**InstanceTypePermissions** | Pointer to [**[]UpdateRoleRequestRoleInstanceTypePermissionsInner**](UpdateRoleRequestRoleInstanceTypePermissionsInner.md) | Set the access level for the specified instance types | [optional] 
**GlobalAppTemplateAccess** | Pointer to **string** | Set the default access level for blueprints | [optional] 
**AppTemplatePermissions** | Pointer to [**[]UpdateRoleRequestRoleAppTemplatePermissionsInner**](UpdateRoleRequestRoleAppTemplatePermissionsInner.md) | Set the access level for the specified blueprints (appTemplates) | [optional] 
**GlobalCatalogItemTypeAccess** | Pointer to **string** | Set the default access level for catalog item types | [optional] 
**CatalogItemTypePermissions** | Pointer to [**[]UpdateRoleRequestRoleCatalogItemTypePermissionsInner**](UpdateRoleRequestRoleCatalogItemTypePermissionsInner.md) | Set the access level for the specified catalog item types | [optional] 
**GlobalPersonaAccess** | Pointer to **string** | Set the default access level for personas | [optional] 
**PersonaPermissions** | Pointer to [**[]UpdateRoleRequestRolePersonaPermissionsInner**](UpdateRoleRequestRolePersonaPermissionsInner.md) | Set the access level for the specified personas | [optional] 
**GlobalVdiPoolAccess** | Pointer to **string** | Set the default access level for VDI pools | [optional] 
**VdiPoolPermissions** | Pointer to [**[]UpdateRoleRequestRoleVdiPoolPermissionsInner**](UpdateRoleRequestRoleVdiPoolPermissionsInner.md) | Set the access level for the specified VDI pools | [optional] 
**GlobalReportTypeAccess** | Pointer to **string** | Set the default access level for report types | [optional] 
**ReportTypePermissions** | Pointer to [**[]UpdateRoleRequestRoleReportTypePermissionsInner**](UpdateRoleRequestRoleReportTypePermissionsInner.md) | Set the access level for the specified report types | [optional] 
**GlobalTaskAccess** | Pointer to **string** | Set the default access level for tasks | [optional] 
**TaskPermissions** | Pointer to [**[]UpdateRoleRequestRoleTaskPermissionsInner**](UpdateRoleRequestRoleTaskPermissionsInner.md) | Set the access level for the specified tasks | [optional] 
**GlobalTaskSetAccess** | Pointer to **string** | Set the default access level for workflows (taskSets) | [optional] 
**TaskSetPermissions** | Pointer to [**[]UpdateRoleRequestRoleTaskSetPermissionsInner**](UpdateRoleRequestRoleTaskSetPermissionsInner.md) | Set the access level for the specified workflows (taskSets) | [optional] 
**GlobalClusterTypeAccess** | Pointer to **string** | Set the default access level for cluster types | [optional] 
**ClusterTypePermissions** | Pointer to [**[]UpdateRoleRequestRoleClusterTypePermissionsInner**](UpdateRoleRequestRoleClusterTypePermissionsInner.md) | Set the access level for the specified cluster types | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateRoleRequestRole{
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


