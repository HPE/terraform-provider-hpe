# ClusterHistory

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**AccountId** | Pointer to **int64** |  | [optional] 
**UniqueId** | Pointer to **string** |  | [optional] 
**ProcessType** | Pointer to [**ClusterHistoryProcessType**](ClusterHistoryProcessType.md) |  | [optional] 
**DisplayName** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**SubType** | Pointer to **NullableString** |  | [optional] 
**SubId** | Pointer to **NullableString** |  | [optional] 
**ZoneId** | Pointer to **NullableInt64** |  | [optional] 
**IntegrationId** | Pointer to **NullableInt64** |  | [optional] 
**AppId** | Pointer to **NullableInt64** |  | [optional] 
**InstanceId** | Pointer to **NullableInt64** |  | [optional] 
**ContainerId** | Pointer to **NullableInt64** |  | [optional] 
**ServerId** | Pointer to **NullableInt64** |  | [optional] 
**ContainerName** | Pointer to **NullableString** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**Reason** | Pointer to **NullableString** |  | [optional] 
**Percent** | Pointer to **float64** |  | [optional] 
**StatusEta** | Pointer to **int64** |  | [optional] 
**Message** | Pointer to **NullableString** |  | [optional] 
**Output** | Pointer to **NullableString** |  | [optional] 
**Error** | Pointer to **NullableString** |  | [optional] 
**StartDate** | Pointer to **time.Time** |  | [optional] 
**EndDate** | Pointer to **time.Time** |  | [optional] 
**Duration** | Pointer to **int64** |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 
**CreatedBy** | Pointer to [**ClusterHistoryCreatedBy**](ClusterHistoryCreatedBy.md) |  | [optional] 
**UpdatedBy** | Pointer to [**ClusterHistoryUpdatedBy**](ClusterHistoryUpdatedBy.md) |  | [optional] 
**Events** | Pointer to [**[]ClusterHistoryEventsInner**](ClusterHistoryEventsInner.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ClusterHistory{
    // Set fields directly
}
```

### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### SubType (Nullable)

Use the Nullable wrapper methods:
- `obj.SubType.IsSet()` — check if set
- `obj.SubType.Get()` — get the inner value (returns pointer)
- `obj.SubType.Set(&val)` — set the value
- `obj.SubType.Unset()` — clear the value
### SubId (Nullable)

Use the Nullable wrapper methods:
- `obj.SubId.IsSet()` — check if set
- `obj.SubId.Get()` — get the inner value (returns pointer)
- `obj.SubId.Set(&val)` — set the value
- `obj.SubId.Unset()` — clear the value
### ZoneId (Nullable)

Use the Nullable wrapper methods:
- `obj.ZoneId.IsSet()` — check if set
- `obj.ZoneId.Get()` — get the inner value (returns pointer)
- `obj.ZoneId.Set(&val)` — set the value
- `obj.ZoneId.Unset()` — clear the value
### IntegrationId (Nullable)

Use the Nullable wrapper methods:
- `obj.IntegrationId.IsSet()` — check if set
- `obj.IntegrationId.Get()` — get the inner value (returns pointer)
- `obj.IntegrationId.Set(&val)` — set the value
- `obj.IntegrationId.Unset()` — clear the value
### AppId (Nullable)

Use the Nullable wrapper methods:
- `obj.AppId.IsSet()` — check if set
- `obj.AppId.Get()` — get the inner value (returns pointer)
- `obj.AppId.Set(&val)` — set the value
- `obj.AppId.Unset()` — clear the value
### InstanceId (Nullable)

Use the Nullable wrapper methods:
- `obj.InstanceId.IsSet()` — check if set
- `obj.InstanceId.Get()` — get the inner value (returns pointer)
- `obj.InstanceId.Set(&val)` — set the value
- `obj.InstanceId.Unset()` — clear the value
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
### ContainerName (Nullable)

Use the Nullable wrapper methods:
- `obj.ContainerName.IsSet()` — check if set
- `obj.ContainerName.Get()` — get the inner value (returns pointer)
- `obj.ContainerName.Set(&val)` — set the value
- `obj.ContainerName.Unset()` — clear the value
### Reason (Nullable)

Use the Nullable wrapper methods:
- `obj.Reason.IsSet()` — check if set
- `obj.Reason.Get()` — get the inner value (returns pointer)
- `obj.Reason.Set(&val)` — set the value
- `obj.Reason.Unset()` — clear the value
### Message (Nullable)

Use the Nullable wrapper methods:
- `obj.Message.IsSet()` — check if set
- `obj.Message.Get()` — get the inner value (returns pointer)
- `obj.Message.Set(&val)` — set the value
- `obj.Message.Unset()` — clear the value
### Output (Nullable)

Use the Nullable wrapper methods:
- `obj.Output.IsSet()` — check if set
- `obj.Output.Get()` — get the inner value (returns pointer)
- `obj.Output.Set(&val)` — set the value
- `obj.Output.Unset()` — clear the value
### Error (Nullable)

Use the Nullable wrapper methods:
- `obj.Error.IsSet()` — check if set
- `obj.Error.Get()` — get the inner value (returns pointer)
- `obj.Error.Set(&val)` — set the value
- `obj.Error.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


