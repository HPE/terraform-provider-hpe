# ClusterPackageUpdate

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Cluster Package name | [optional] 
**Description** | Pointer to **NullableString** | Cluster Package description | [optional] 
**Code** | Pointer to **string** | Cluster Package code | [optional] 
**PackageVersion** | Pointer to **string** | Version of the cluster package | [optional] 
**PackageType** | Pointer to **string** | Package Type | [optional] 
**Type** | Pointer to **string** | type | [optional] 
**Enabled** | Pointer to **bool** | Can be used to enable / disable the cluster package. | [optional] [default to true]
**IconPath** | Pointer to **string** | Icon Path, relative location of an icon image, eg. /assets/containers-png/nginx.png. | [optional] 
**SpecTemplates** | Pointer to [**[]ClusterPackageUpdateSpecTemplatesInner**](ClusterPackageUpdateSpecTemplatesInner.md) | Array of resource spec templates | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ClusterPackageUpdate{
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


