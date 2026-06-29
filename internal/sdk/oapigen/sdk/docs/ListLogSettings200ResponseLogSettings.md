# ListLogSettings200ResponseLogSettings

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Enabled** | Pointer to **bool** |  | [optional] 
**RetentionDays** | Pointer to **string** |  | [optional] 
**SyslogRules** | Pointer to **[]map[string]interface{}** |  | [optional] 
**Integrations** | Pointer to **[]map[string]interface{}** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListLogSettings200ResponseLogSettings{
    // Set fields directly
}
```

### SyslogRules (Nullable)

Use the Nullable wrapper methods:
- `obj.SyslogRules.IsSet()` — check if set
- `obj.SyslogRules.Get()` — get the inner value (returns pointer)
- `obj.SyslogRules.Set(&val)` — set the value
- `obj.SyslogRules.Unset()` — clear the value
### Integrations (Nullable)

Use the Nullable wrapper methods:
- `obj.Integrations.IsSet()` — check if set
- `obj.Integrations.Get()` — get the inner value (returns pointer)
- `obj.Integrations.Set(&val)` — set the value
- `obj.Integrations.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


