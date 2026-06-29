# ListHealth200ResponseAllOfHealthRabbit

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Success** | Pointer to **bool** |  | [optional] 
**BusyQueues** | Pointer to **[]map[string]interface{}** |  | [optional] 
**ErrorQueues** | Pointer to **[]map[string]interface{}** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**Queues** | Pointer to [**[]ListHealth200ResponseAllOfHealthRabbitQueuesInner**](ListHealth200ResponseAllOfHealthRabbitQueuesInner.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListHealth200ResponseAllOfHealthRabbit{
    // Set fields directly
}
```

### BusyQueues (Nullable)

Use the Nullable wrapper methods:
- `obj.BusyQueues.IsSet()` — check if set
- `obj.BusyQueues.Get()` — get the inner value (returns pointer)
- `obj.BusyQueues.Set(&val)` — set the value
- `obj.BusyQueues.Unset()` — clear the value
### ErrorQueues (Nullable)

Use the Nullable wrapper methods:
- `obj.ErrorQueues.IsSet()` — check if set
- `obj.ErrorQueues.Get()` — get the inner value (returns pointer)
- `obj.ErrorQueues.Set(&val)` — set the value
- `obj.ErrorQueues.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


