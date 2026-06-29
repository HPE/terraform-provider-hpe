# AddIntegrations200ResponseAllOfIntegrationOneOfConfig

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Inventory** | Pointer to **NullableString** |  | [optional] 
**DefaultBranch** | Pointer to **string** |  | [optional] 
**CacheEnabled** | Pointer to **NullableString** |  | [optional] 
**AnsiblePlaybooks** | Pointer to **string** |  | [optional] 
**AnsibleRoles** | Pointer to **string** |  | [optional] 
**AnsibleGroupVars** | Pointer to **string** |  | [optional] 
**AnsibleHostVars** | Pointer to **string** |  | [optional] 
**AnsibleCommandBus** | Pointer to **bool** |  | [optional] 
**AnsibleVerbose** | Pointer to **bool** |  | [optional] 
**AnsibleGalaxyEnabled** | Pointer to **bool** |  | [optional] 
**AnsibleDefaultBranch** | Pointer to **string** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddIntegrations200ResponseAllOfIntegrationOneOfConfig{
    // Set fields directly
}
```

### Inventory (Nullable)

Use the Nullable wrapper methods:
- `obj.Inventory.IsSet()` — check if set
- `obj.Inventory.Get()` — get the inner value (returns pointer)
- `obj.Inventory.Set(&val)` — set the value
- `obj.Inventory.Unset()` — clear the value
### CacheEnabled (Nullable)

Use the Nullable wrapper methods:
- `obj.CacheEnabled.IsSet()` — check if set
- `obj.CacheEnabled.Get()` — get the inner value (returns pointer)
- `obj.CacheEnabled.Set(&val)` — set the value
- `obj.CacheEnabled.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


