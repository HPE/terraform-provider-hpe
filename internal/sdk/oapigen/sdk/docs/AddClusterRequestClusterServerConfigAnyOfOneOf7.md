# AddClusterRequestClusterServerConfigAnyOfOneOf7

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**NodeCount** | Pointer to **int64** | Number of workers or hosts | [optional] 
**PodCidr** | Pointer to **string** |  | [optional] 
**VmwareFolderId** | Pointer to **string** |  | [optional] 
**ServiceCidr** | Pointer to **string** |  | [optional] 
**CreateUser** | Pointer to **bool** |  | [optional] 
**DefaultRepoAccount** | Pointer to **NullableInt64** | Default Repo Account is the repository to be used when pulling images.  Default behavior is to be anonymous, which does have limits on allowed image pulls from public Docker Repos. | [optional] 
**ImageServer** | Pointer to **string** | Act as Image Server. Set to on to use the Default Repo Account to pull images. | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddClusterRequestClusterServerConfigAnyOfOneOf7{
    // Set fields directly
}
```

### DefaultRepoAccount (Nullable)

Use the Nullable wrapper methods:
- `obj.DefaultRepoAccount.IsSet()` — check if set
- `obj.DefaultRepoAccount.Get()` — get the inner value (returns pointer)
- `obj.DefaultRepoAccount.Set(&val)` — set the value
- `obj.DefaultRepoAccount.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


