# Workflow

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Type** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**Labels** | Pointer to **[]string** |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 
**AccountId** | Pointer to **int64** |  | [optional] 
**Platform** | Pointer to **NullableString** |  | [optional] 
**Visibility** | Pointer to **string** |  | [optional] 
**AllowCustomConfig** | Pointer to **bool** |  | [optional] 
**Tasks** | Pointer to **[]int64** |  | [optional] 
**OptionTypes** | Pointer to [**[]AddWorkflows200ResponseAllOfTaskSetOptionTypesInner**](AddWorkflows200ResponseAllOfTaskSetOptionTypesInner.md) |  | [optional] 
**TaskSetTasks** | Pointer to [**[]AddWorkflows200ResponseAllOfTaskSetTaskSetTasksInner**](AddWorkflows200ResponseAllOfTaskSetTaskSetTasksInner.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &Workflow{
    // Set fields directly
}
```

### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### Labels (Nullable)

Use the Nullable wrapper methods:
- `obj.Labels.IsSet()` — check if set
- `obj.Labels.Get()` — get the inner value (returns pointer)
- `obj.Labels.Set(&val)` — set the value
- `obj.Labels.Unset()` — clear the value
### Platform (Nullable)

Use the Nullable wrapper methods:
- `obj.Platform.IsSet()` — check if set
- `obj.Platform.Get()` — get the inner value (returns pointer)
- `obj.Platform.Set(&val)` — set the value
- `obj.Platform.Unset()` — clear the value
### OptionTypes (Nullable)

Use the Nullable wrapper methods:
- `obj.OptionTypes.IsSet()` — check if set
- `obj.OptionTypes.Get()` — get the inner value (returns pointer)
- `obj.OptionTypes.Set(&val)` — set the value
- `obj.OptionTypes.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


