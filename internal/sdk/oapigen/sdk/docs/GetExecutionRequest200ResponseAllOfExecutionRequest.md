# GetExecutionRequest200ResponseAllOfExecutionRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**UniqueId** | Pointer to **string** |  | [optional] 
**ContainerId** | Pointer to **NullableString** |  | [optional] 
**ServerId** | Pointer to **NullableString** |  | [optional] 
**InstanceId** | Pointer to **int64** |  | [optional] 
**ResourceId** | Pointer to **NullableString** |  | [optional] 
**AppId** | Pointer to **NullableString** |  | [optional] 
**StdOut** | Pointer to **string** |  | [optional] 
**StdErr** | Pointer to **string** |  | [optional] 
**ExitCode** | Pointer to **int64** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**ExpiresAt** | Pointer to **time.Time** |  | [optional] 
**CreatedById** | Pointer to **int64** |  | [optional] 
**StatusMessage** | Pointer to **NullableString** |  | [optional] 
**ErrorMessage** | Pointer to **NullableString** |  | [optional] 
**Config** | Pointer to **map[string]interface{}** |  | [optional] 
**RawData** | Pointer to **NullableString** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetExecutionRequest200ResponseAllOfExecutionRequest{
    // Set fields directly
}
```

### ContainerId (Nullable)

Use the Nullable wrapper methods:
- `obj.ContainerId.IsSet()` — check if set
- `obj.ContainerId.Get()` — get the inner value (returns pointer)
- `obj.ContainerId.Set(&val)` — set the value
- `obj.ContainerId.Unset()` — clear the value
### ServerId (Nullable)

Use the Nullable wrapper methods:
- `obj.ServerId.IsSet()` — check if set
- `obj.ServerId.Get()` — get the inner value (returns pointer)
- `obj.ServerId.Set(&val)` — set the value
- `obj.ServerId.Unset()` — clear the value
### ResourceId (Nullable)

Use the Nullable wrapper methods:
- `obj.ResourceId.IsSet()` — check if set
- `obj.ResourceId.Get()` — get the inner value (returns pointer)
- `obj.ResourceId.Set(&val)` — set the value
- `obj.ResourceId.Unset()` — clear the value
### AppId (Nullable)

Use the Nullable wrapper methods:
- `obj.AppId.IsSet()` — check if set
- `obj.AppId.Get()` — get the inner value (returns pointer)
- `obj.AppId.Set(&val)` — set the value
- `obj.AppId.Unset()` — clear the value
### StatusMessage (Nullable)

Use the Nullable wrapper methods:
- `obj.StatusMessage.IsSet()` — check if set
- `obj.StatusMessage.Get()` — get the inner value (returns pointer)
- `obj.StatusMessage.Set(&val)` — set the value
- `obj.StatusMessage.Unset()` — clear the value
### ErrorMessage (Nullable)

Use the Nullable wrapper methods:
- `obj.ErrorMessage.IsSet()` — check if set
- `obj.ErrorMessage.Get()` — get the inner value (returns pointer)
- `obj.ErrorMessage.Set(&val)` — set the value
- `obj.ErrorMessage.Unset()` — clear the value
### RawData (Nullable)

Use the Nullable wrapper methods:
- `obj.RawData.IsSet()` — check if set
- `obj.RawData.Get()` — get the inner value (returns pointer)
- `obj.RawData.Set(&val)` — set the value
- `obj.RawData.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


