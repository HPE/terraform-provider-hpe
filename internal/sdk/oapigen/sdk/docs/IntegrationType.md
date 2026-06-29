# IntegrationType

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Code** | Pointer to **string** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**Category** | Pointer to **string** |  | [optional] 
**Enabled** | Pointer to **bool** |  | [optional] 
**Creatable** | Pointer to **bool** |  | [optional] 
**HasCMDB** | Pointer to **bool** |  | [optional] 
**HasCMDBDiscovery** | Pointer to **bool** |  | [optional] 
**HasCM** | Pointer to **bool** |  | [optional] 
**HasDNS** | Pointer to **bool** |  | [optional] 
**HasApprovals** | Pointer to **bool** |  | [optional] 
**HasDeleteApprovals** | Pointer to **bool** |  | [optional] 
**HasReconfigureApprovals** | Pointer to **bool** |  | [optional] 
**IsPlugin** | Pointer to **bool** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &IntegrationType{
    // Set fields directly
}
```

### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


