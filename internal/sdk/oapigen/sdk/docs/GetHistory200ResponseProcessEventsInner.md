# GetHistory200ResponseProcessEventsInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**ProcessId** | Pointer to **int64** |  | [optional] 
**AccountId** | Pointer to **int64** |  | [optional] 
**UniqueId** | Pointer to **string** |  | [optional] 
**ProcessType** | Pointer to [**GetHistory200ResponseProcessEventsInnerProcessType**](GetHistory200ResponseProcessEventsInnerProcessType.md) |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**RefType** | Pointer to **string** |  | [optional] 
**RefId** | Pointer to **int64** |  | [optional] 
**SubType** | Pointer to **NullableString** |  | [optional] 
**SubId** | Pointer to **NullableString** |  | [optional] 
**ZoneId** | Pointer to **int64** |  | [optional] 
**IntegrationId** | Pointer to **NullableString** |  | [optional] 
**InstanceId** | Pointer to **int64** |  | [optional] 
**ContainerId** | Pointer to **int64** |  | [optional] 
**ServerId** | Pointer to **int64** |  | [optional] 
**ContainerName** | Pointer to **string** |  | [optional] 
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
**CreatedBy** | Pointer to [**GetHistory200ResponseProcessEventsInnerCreatedBy**](GetHistory200ResponseProcessEventsInnerCreatedBy.md) |  | [optional] 
**UpdatedBy** | Pointer to [**GetHistory200ResponseProcessEventsInnerUpdatedBy**](GetHistory200ResponseProcessEventsInnerUpdatedBy.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetHistory200ResponseProcessEventsInner{
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
### IntegrationId (Nullable)

Use the Nullable wrapper methods:
- `obj.IntegrationId.IsSet()` — check if set
- `obj.IntegrationId.Get()` — get the inner value (returns pointer)
- `obj.IntegrationId.Set(&val)` — set the value
- `obj.IntegrationId.Unset()` — clear the value
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


