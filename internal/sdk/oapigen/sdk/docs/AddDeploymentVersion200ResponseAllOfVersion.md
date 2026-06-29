# AddDeploymentVersion200ResponseAllOfVersion

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**DeployType** | Pointer to **string** |  | [optional] 
**DeploymentId** | Pointer to **int64** |  | [optional] 
**FetchUrl** | Pointer to **NullableString** |  | [optional] 
**GitUrl** | Pointer to **NullableString** |  | [optional] 
**GitRef** | Pointer to **NullableString** |  | [optional] 
**UserVersion** | Pointer to **string** |  | [optional] 
**Version** | Pointer to **string** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddDeploymentVersion200ResponseAllOfVersion{
    // Set fields directly
}
```

### FetchUrl (Nullable)

Use the Nullable wrapper methods:
- `obj.FetchUrl.IsSet()` — check if set
- `obj.FetchUrl.Get()` — get the inner value (returns pointer)
- `obj.FetchUrl.Set(&val)` — set the value
- `obj.FetchUrl.Unset()` — clear the value
### GitUrl (Nullable)

Use the Nullable wrapper methods:
- `obj.GitUrl.IsSet()` — check if set
- `obj.GitUrl.Get()` — get the inner value (returns pointer)
- `obj.GitUrl.Set(&val)` — set the value
- `obj.GitUrl.Unset()` — clear the value
### GitRef (Nullable)

Use the Nullable wrapper methods:
- `obj.GitRef.IsSet()` — check if set
- `obj.GitRef.Get()` — get the inner value (returns pointer)
- `obj.GitRef.Set(&val)` — set the value
- `obj.GitRef.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


