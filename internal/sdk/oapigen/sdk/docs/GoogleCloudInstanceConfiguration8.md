# GoogleCloudInstanceConfiguration8

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**NoAgent** | Pointer to **NullableBool** | Skipping Agent installation will result in a lack of logging and guest operating system statistics. Automation scripts may also be adversely affected. | [optional] [default to false]
**GoogleZoneId** | Pointer to **int64** | External ID of the Google zone to use for instance. | [optional] 
**ExternalIpType** | Pointer to **int64** | External IP Type.  &#x60;-1&#x60; for ephemeral IP. | [optional] 
**NetworkTags** | Pointer to **string** | Network Tags | [optional] 
**ServiceAccount** | Pointer to **string** | Service Account | [optional] 
**AccessScope** | Pointer to **string** | Access Scope | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GoogleCloudInstanceConfiguration8{
    // Set fields directly
}
```

### NoAgent (Nullable)

Use the Nullable wrapper methods:
- `obj.NoAgent.IsSet()` — check if set
- `obj.NoAgent.Get()` — get the inner value (returns pointer)
- `obj.NoAgent.Set(&val)` — set the value
- `obj.NoAgent.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


