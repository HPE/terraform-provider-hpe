# SnapshotsInstance200ResponseSnapshotsInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**ExternalId** | Pointer to **NullableString** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**State** | Pointer to **NullableString** |  | [optional] 
**SnapshotType** | Pointer to **string** |  | [optional] 
**SnapshotCreated** | Pointer to **NullableTime** |  | [optional] 
**Zone** | Pointer to [**NullableSnapshotsInstance200ResponseSnapshotsInnerZone**](SnapshotsInstance200ResponseSnapshotsInnerZone.md) |  | [optional] 
**Datastore** | Pointer to **NullableString** |  | [optional] 
**ParentSnapshot** | Pointer to **NullableString** |  | [optional] 
**SnapshotFiles** | Pointer to [**[]SnapshotsInstance200ResponseSnapshotsInnerSnapshotFilesInner**](SnapshotsInstance200ResponseSnapshotsInnerSnapshotFilesInner.md) |  | [optional] 
**CurrentlyActive** | Pointer to **bool** |  | [optional] 
**MemorySnapshot** | Pointer to **bool** |  | [optional] 
**ForExport** | Pointer to **bool** |  | [optional] 
**ForBackup** | Pointer to **bool** |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &SnapshotsInstance200ResponseSnapshotsInner{
    // Set fields directly
}
```

### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### ExternalId (Nullable)

Use the Nullable wrapper methods:
- `obj.ExternalId.IsSet()` — check if set
- `obj.ExternalId.Get()` — get the inner value (returns pointer)
- `obj.ExternalId.Set(&val)` — set the value
- `obj.ExternalId.Unset()` — clear the value
### State (Nullable)

Use the Nullable wrapper methods:
- `obj.State.IsSet()` — check if set
- `obj.State.Get()` — get the inner value (returns pointer)
- `obj.State.Set(&val)` — set the value
- `obj.State.Unset()` — clear the value
### SnapshotCreated (Nullable)

Use the Nullable wrapper methods:
- `obj.SnapshotCreated.IsSet()` — check if set
- `obj.SnapshotCreated.Get()` — get the inner value (returns pointer)
- `obj.SnapshotCreated.Set(&val)` — set the value
- `obj.SnapshotCreated.Unset()` — clear the value
### Zone (Nullable)

Use the Nullable wrapper methods:
- `obj.Zone.IsSet()` — check if set
- `obj.Zone.Get()` — get the inner value (returns pointer)
- `obj.Zone.Set(&val)` — set the value
- `obj.Zone.Unset()` — clear the value
### Datastore (Nullable)

Use the Nullable wrapper methods:
- `obj.Datastore.IsSet()` — check if set
- `obj.Datastore.Get()` — get the inner value (returns pointer)
- `obj.Datastore.Set(&val)` — set the value
- `obj.Datastore.Unset()` — clear the value
### ParentSnapshot (Nullable)

Use the Nullable wrapper methods:
- `obj.ParentSnapshot.IsSet()` — check if set
- `obj.ParentSnapshot.Get()` — get the inner value (returns pointer)
- `obj.ParentSnapshot.Set(&val)` — set the value
- `obj.ParentSnapshot.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


