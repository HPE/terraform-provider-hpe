# GetApprovals200ResponseApproval

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**InternalId** | Pointer to **NullableString** |  | [optional] 
**ExternalId** | Pointer to **NullableString** |  | [optional] 
**ExternalName** | Pointer to **NullableString** |  | [optional] 
**RequestType** | Pointer to **string** |  | [optional] 
**Account** | Pointer to [**NullableGetApprovals200ResponseApprovalAccount**](GetApprovals200ResponseApprovalAccount.md) |  | [optional] 
**Approver** | Pointer to [**NullableGetApprovals200ResponseApprovalApprover**](GetApprovals200ResponseApprovalApprover.md) |  | [optional] 
**AccountIntegration** | Pointer to **NullableString** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**ErrorMessage** | Pointer to **NullableString** |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 
**RequestBy** | Pointer to **string** |  | [optional] 
**ApprovalItems** | Pointer to [**[]GetApprovals200ResponseApprovalApprovalItemsInner**](GetApprovals200ResponseApprovalApprovalItemsInner.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetApprovals200ResponseApproval{
    // Set fields directly
}
```

### InternalId (Nullable)

Use the Nullable wrapper methods:
- `obj.InternalId.IsSet()` — check if set
- `obj.InternalId.Get()` — get the inner value (returns pointer)
- `obj.InternalId.Set(&val)` — set the value
- `obj.InternalId.Unset()` — clear the value
### ExternalId (Nullable)

Use the Nullable wrapper methods:
- `obj.ExternalId.IsSet()` — check if set
- `obj.ExternalId.Get()` — get the inner value (returns pointer)
- `obj.ExternalId.Set(&val)` — set the value
- `obj.ExternalId.Unset()` — clear the value
### ExternalName (Nullable)

Use the Nullable wrapper methods:
- `obj.ExternalName.IsSet()` — check if set
- `obj.ExternalName.Get()` — get the inner value (returns pointer)
- `obj.ExternalName.Set(&val)` — set the value
- `obj.ExternalName.Unset()` — clear the value
### Account (Nullable)

Use the Nullable wrapper methods:
- `obj.Account.IsSet()` — check if set
- `obj.Account.Get()` — get the inner value (returns pointer)
- `obj.Account.Set(&val)` — set the value
- `obj.Account.Unset()` — clear the value
### Approver (Nullable)

Use the Nullable wrapper methods:
- `obj.Approver.IsSet()` — check if set
- `obj.Approver.Get()` — get the inner value (returns pointer)
- `obj.Approver.Set(&val)` — set the value
- `obj.Approver.Unset()` — clear the value
### AccountIntegration (Nullable)

Use the Nullable wrapper methods:
- `obj.AccountIntegration.IsSet()` — check if set
- `obj.AccountIntegration.Get()` — get the inner value (returns pointer)
- `obj.AccountIntegration.Set(&val)` — set the value
- `obj.AccountIntegration.Unset()` — clear the value
### ErrorMessage (Nullable)

Use the Nullable wrapper methods:
- `obj.ErrorMessage.IsSet()` — check if set
- `obj.ErrorMessage.Get()` — get the inner value (returns pointer)
- `obj.ErrorMessage.Set(&val)` — set the value
- `obj.ErrorMessage.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


