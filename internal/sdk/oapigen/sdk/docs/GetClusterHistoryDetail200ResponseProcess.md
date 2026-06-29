# GetClusterHistoryDetail200ResponseProcess

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**AccountId** | Pointer to **int64** |  | [optional] 
**UniqueId** | Pointer to **string** |  | [optional] 
**ProcessType** | Pointer to [**GetClusterHistoryDetail200ResponseProcessProcessType**](GetClusterHistoryDetail200ResponseProcessProcessType.md) |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**SubType** | Pointer to **NullableString** |  | [optional] 
**SubId** | Pointer to **NullableString** |  | [optional] 
**ZoneId** | Pointer to **NullableString** |  | [optional] 
**IntegrationId** | Pointer to **NullableString** |  | [optional] 
**AppId** | Pointer to **NullableInt64** |  | [optional] 
**InstanceId** | Pointer to **NullableString** |  | [optional] 
**ContainerId** | Pointer to **NullableString** |  | [optional] 
**ServerId** | Pointer to **int64** |  | [optional] 
**ContainerName** | Pointer to **NullableString** |  | [optional] 
**DisplayName** | Pointer to **string** |  | [optional] 
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
**CreatedBy** | Pointer to [**GetClusterHistoryDetail200ResponseProcessCreatedBy**](GetClusterHistoryDetail200ResponseProcessCreatedBy.md) |  | [optional] 
**UpdatedBy** | Pointer to [**GetClusterHistoryDetail200ResponseProcessUpdatedBy**](GetClusterHistoryDetail200ResponseProcessUpdatedBy.md) |  | [optional] 
**Events** | Pointer to [**[]GetClusterHistoryDetail200ResponseProcessEventsInner**](GetClusterHistoryDetail200ResponseProcessEventsInner.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetClusterHistoryDetail200ResponseProcess{
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


