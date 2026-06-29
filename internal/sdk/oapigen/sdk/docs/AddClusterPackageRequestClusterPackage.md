# AddClusterPackageRequestClusterPackage

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | Cluster Package name | 
**Description** | Pointer to **NullableString** | Cluster Package description | [optional] 
**Code** | **string** | Cluster Package code | 
**PackageVersion** | **string** | Version of the cluster package | 
**PackageType** | **string** | Package Type | 
**Type** | **string** | type | 
**Enabled** | **bool** | Can be used to enable / disable the cluster package. | [default to true]
**IconPath** | Pointer to **string** | Icon Path, relative location of an icon image, eg. /assets/containers-png/nginx.png. | [optional] 
**SpecTemplates** | **[]int64** | Array of resource spec templates | 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddClusterPackageRequestClusterPackage{
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


