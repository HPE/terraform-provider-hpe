# SNMPCheck1

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **interface{}** |  | 
**Description** | Pointer to **NullableString** | Optional description field | [optional] 
**CheckInterval** | Pointer to **int32** | Number of milliseconds you want between check executions (minimum is 1 minute, depending on your subscription plan) | [optional] [default to 300000]
**InUptime** | Pointer to **bool** | Used to determine if check should affect account wide availability calculations | [optional] [default to true]
**Active** | Pointer to **bool** | Used to determine if check should be scheduled to execute | [optional] [default to true]
**Severity** | Pointer to **string** | Severity level threshold for sending notifications. | [optional] [default to "critical"]
**CheckType** | Pointer to [**SNMPCheck1AllOfCheckType**](SNMPCheck1AllOfCheckType.md) |  | [optional] 
**Config** | Pointer to [**SNMPCheck1AllOfConfig**](SNMPCheck1AllOfConfig.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &SNMPCheck1{
    // Set fields directly
}
```

### Name (Nullable)

Use the Nullable wrapper methods:
- `obj.Name.IsSet()` — check if set
- `obj.Name.Get()` — get the inner value (returns pointer)
- `obj.Name.Set(&val)` — set the value
- `obj.Name.Unset()` — clear the value
### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


