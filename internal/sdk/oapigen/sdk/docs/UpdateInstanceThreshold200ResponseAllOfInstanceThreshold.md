# UpdateInstanceThreshold200ResponseAllOfInstanceThreshold

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**AutoUp** | Pointer to **bool** |  | [optional] 
**AutoDown** | Pointer to **bool** |  | [optional] 
**MinCount** | Pointer to **int64** |  | [optional] 
**MaxCount** | Pointer to **int64** |  | [optional] 
**ScaleIncrement** | Pointer to **int64** |  | [optional] 
**CpuEnabled** | Pointer to **bool** |  | [optional] 
**MinCpu** | Pointer to **int64** |  | [optional] 
**MaxCpu** | Pointer to **int64** |  | [optional] 
**MemoryEnabled** | Pointer to **bool** |  | [optional] 
**MinMemory** | Pointer to **int64** |  | [optional] 
**MaxMemory** | Pointer to **int64** |  | [optional] 
**DiskEnabled** | Pointer to **bool** |  | [optional] 
**MinDisk** | Pointer to **int64** |  | [optional] 
**MaxDisk** | Pointer to **int64** |  | [optional] 
**MinNetwork** | Pointer to **NullableString** |  | [optional] 
**NetworkEnabled** | Pointer to **bool** |  | [optional] 
**IopsEnabled** | Pointer to **bool** |  | [optional] 
**MinIops** | Pointer to **NullableString** |  | [optional] 
**MaxIops** | Pointer to **NullableString** |  | [optional] 
**Comment** | Pointer to **NullableString** |  | [optional] 
**ZoneId** | Pointer to **NullableInt64** |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateInstanceThreshold200ResponseAllOfInstanceThreshold{
    // Set fields directly
}
```

### MinNetwork (Nullable)

Use the Nullable wrapper methods:
- `obj.MinNetwork.IsSet()` — check if set
- `obj.MinNetwork.Get()` — get the inner value (returns pointer)
- `obj.MinNetwork.Set(&val)` — set the value
- `obj.MinNetwork.Unset()` — clear the value
### MinIops (Nullable)

Use the Nullable wrapper methods:
- `obj.MinIops.IsSet()` — check if set
- `obj.MinIops.Get()` — get the inner value (returns pointer)
- `obj.MinIops.Set(&val)` — set the value
- `obj.MinIops.Unset()` — clear the value
### MaxIops (Nullable)

Use the Nullable wrapper methods:
- `obj.MaxIops.IsSet()` — check if set
- `obj.MaxIops.Get()` — get the inner value (returns pointer)
- `obj.MaxIops.Set(&val)` — set the value
- `obj.MaxIops.Unset()` — clear the value
### Comment (Nullable)

Use the Nullable wrapper methods:
- `obj.Comment.IsSet()` — check if set
- `obj.Comment.Get()` — get the inner value (returns pointer)
- `obj.Comment.Set(&val)` — set the value
- `obj.Comment.Unset()` — clear the value
### ZoneId (Nullable)

Use the Nullable wrapper methods:
- `obj.ZoneId.IsSet()` — check if set
- `obj.ZoneId.Get()` — get the inner value (returns pointer)
- `obj.ZoneId.Set(&val)` — set the value
- `obj.ZoneId.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


