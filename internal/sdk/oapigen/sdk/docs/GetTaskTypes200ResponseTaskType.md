# GetTaskTypes200ResponseTaskType

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Code** | Pointer to **string** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Category** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**Scriptable** | Pointer to **bool** |  | [optional] 
**Enabled** | Pointer to **bool** |  | [optional] 
**HasResults** | Pointer to **bool** |  | [optional] 
**AllowExecuteLocal** | Pointer to **NullableBool** |  | [optional] 
**AllowExecuteRemote** | Pointer to **NullableBool** |  | [optional] 
**AllowExecuteResource** | Pointer to **NullableBool** |  | [optional] 
**AllowLocalRepo** | Pointer to **NullableBool** |  | [optional] 
**AllowRemoteKeyAuth** | Pointer to **NullableBool** |  | [optional] 
**OptionTypes** | Pointer to [**[]GetTaskTypes200ResponseTaskTypeOptionTypesInner**](GetTaskTypes200ResponseTaskTypeOptionTypesInner.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetTaskTypes200ResponseTaskType{
    // Set fields directly
}
```

### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### AllowExecuteLocal (Nullable)

Use the Nullable wrapper methods:
- `obj.AllowExecuteLocal.IsSet()` — check if set
- `obj.AllowExecuteLocal.Get()` — get the inner value (returns pointer)
- `obj.AllowExecuteLocal.Set(&val)` — set the value
- `obj.AllowExecuteLocal.Unset()` — clear the value
### AllowExecuteRemote (Nullable)

Use the Nullable wrapper methods:
- `obj.AllowExecuteRemote.IsSet()` — check if set
- `obj.AllowExecuteRemote.Get()` — get the inner value (returns pointer)
- `obj.AllowExecuteRemote.Set(&val)` — set the value
- `obj.AllowExecuteRemote.Unset()` — clear the value
### AllowExecuteResource (Nullable)

Use the Nullable wrapper methods:
- `obj.AllowExecuteResource.IsSet()` — check if set
- `obj.AllowExecuteResource.Get()` — get the inner value (returns pointer)
- `obj.AllowExecuteResource.Set(&val)` — set the value
- `obj.AllowExecuteResource.Unset()` — clear the value
### AllowLocalRepo (Nullable)

Use the Nullable wrapper methods:
- `obj.AllowLocalRepo.IsSet()` — check if set
- `obj.AllowLocalRepo.Get()` — get the inner value (returns pointer)
- `obj.AllowLocalRepo.Set(&val)` — set the value
- `obj.AllowLocalRepo.Unset()` — clear the value
### AllowRemoteKeyAuth (Nullable)

Use the Nullable wrapper methods:
- `obj.AllowRemoteKeyAuth.IsSet()` — check if set
- `obj.AllowRemoteKeyAuth.Get()` — get the inner value (returns pointer)
- `obj.AllowRemoteKeyAuth.Set(&val)` — set the value
- `obj.AllowRemoteKeyAuth.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


