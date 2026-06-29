# GetMigration200ResponseMigrationServersInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Status** | Pointer to **string** | Migration Server Status. The possible status values are: &#39;pending&#39;, &#39;precheck&#39;, &#39;running&#39;, &#39;failed&#39;, &#39;completed&#39; | [optional] 
**StatusMessage** | Pointer to **NullableString** | Status Message | [optional] 
**SourceServer** | Pointer to [**GetMigration200ResponseMigrationServersInnerSourceServer**](GetMigration200ResponseMigrationServersInnerSourceServer.md) |  | [optional] 
**DestinationServer** | Pointer to [**GetMigration200ResponseMigrationServersInnerDestinationServer**](GetMigration200ResponseMigrationServersInnerDestinationServer.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetMigration200ResponseMigrationServersInner{
    // Set fields directly
}
```

### StatusMessage (Nullable)

Use the Nullable wrapper methods:
- `obj.StatusMessage.IsSet()` — check if set
- `obj.StatusMessage.Get()` — get the inner value (returns pointer)
- `obj.StatusMessage.Set(&val)` — set the value
- `obj.StatusMessage.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


