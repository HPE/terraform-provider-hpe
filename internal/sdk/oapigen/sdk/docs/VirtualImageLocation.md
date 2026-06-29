# VirtualImageLocation

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Cloud** | Pointer to [**VirtualImageLocationCloud**](VirtualImageLocationCloud.md) |  | [optional] 
**Code** | Pointer to **string** |  | [optional] 
**InternalId** | Pointer to **NullableString** |  | [optional] 
**ExternalId** | Pointer to **string** |  | [optional] 
**ExternalDiskId** | Pointer to **string** |  | [optional] 
**RemotePath** | Pointer to **NullableString** |  | [optional] 
**ImagePath** | Pointer to **NullableString** |  | [optional] 
**ImageName** | Pointer to **string** |  | [optional] 
**ImageRegion** | Pointer to **string** |  | [optional] 
**ImageFolder** | Pointer to **NullableString** |  | [optional] 
**RefType** | Pointer to **string** |  | [optional] 
**RefId** | Pointer to **int64** |  | [optional] 
**NodeRefType** | Pointer to **NullableString** |  | [optional] 
**NodeRefId** | Pointer to **NullableString** |  | [optional] 
**SubRefType** | Pointer to **NullableString** |  | [optional] 
**SubRefId** | Pointer to **NullableString** |  | [optional] 
**IsPublic** | Pointer to **bool** |  | [optional] 
**SystemImage** | Pointer to **bool** |  | [optional] 
**DiskIndex** | Pointer to **int64** |  | [optional] 
**PricePlan** | Pointer to **NullableString** |  | [optional] 
**Volumes** | Pointer to **[]map[string]interface{}** |  | [optional] 
**StorageControllers** | Pointer to **[]map[string]interface{}** |  | [optional] 
**NetworkInterfaces** | Pointer to **[]map[string]interface{}** |  | [optional] 
**VirtualImage** | Pointer to [**VirtualImageLocationVirtualImage**](VirtualImageLocationVirtualImage.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &VirtualImageLocation{
    // Set fields directly
}
```

### InternalId (Nullable)

Use the Nullable wrapper methods:
- `obj.InternalId.IsSet()` — check if set
- `obj.InternalId.Get()` — get the inner value (returns pointer)
- `obj.InternalId.Set(&val)` — set the value
- `obj.InternalId.Unset()` — clear the value
### RemotePath (Nullable)

Use the Nullable wrapper methods:
- `obj.RemotePath.IsSet()` — check if set
- `obj.RemotePath.Get()` — get the inner value (returns pointer)
- `obj.RemotePath.Set(&val)` — set the value
- `obj.RemotePath.Unset()` — clear the value
### ImagePath (Nullable)

Use the Nullable wrapper methods:
- `obj.ImagePath.IsSet()` — check if set
- `obj.ImagePath.Get()` — get the inner value (returns pointer)
- `obj.ImagePath.Set(&val)` — set the value
- `obj.ImagePath.Unset()` — clear the value
### ImageFolder (Nullable)

Use the Nullable wrapper methods:
- `obj.ImageFolder.IsSet()` — check if set
- `obj.ImageFolder.Get()` — get the inner value (returns pointer)
- `obj.ImageFolder.Set(&val)` — set the value
- `obj.ImageFolder.Unset()` — clear the value
### NodeRefType (Nullable)

Use the Nullable wrapper methods:
- `obj.NodeRefType.IsSet()` — check if set
- `obj.NodeRefType.Get()` — get the inner value (returns pointer)
- `obj.NodeRefType.Set(&val)` — set the value
- `obj.NodeRefType.Unset()` — clear the value
### NodeRefId (Nullable)

Use the Nullable wrapper methods:
- `obj.NodeRefId.IsSet()` — check if set
- `obj.NodeRefId.Get()` — get the inner value (returns pointer)
- `obj.NodeRefId.Set(&val)` — set the value
- `obj.NodeRefId.Unset()` — clear the value
### SubRefType (Nullable)

Use the Nullable wrapper methods:
- `obj.SubRefType.IsSet()` — check if set
- `obj.SubRefType.Get()` — get the inner value (returns pointer)
- `obj.SubRefType.Set(&val)` — set the value
- `obj.SubRefType.Unset()` — clear the value
### SubRefId (Nullable)

Use the Nullable wrapper methods:
- `obj.SubRefId.IsSet()` — check if set
- `obj.SubRefId.Get()` — get the inner value (returns pointer)
- `obj.SubRefId.Set(&val)` — set the value
- `obj.SubRefId.Unset()` — clear the value
### PricePlan (Nullable)

Use the Nullable wrapper methods:
- `obj.PricePlan.IsSet()` — check if set
- `obj.PricePlan.Get()` — get the inner value (returns pointer)
- `obj.PricePlan.Set(&val)` — set the value
- `obj.PricePlan.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


