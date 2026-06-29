# ApprovalItemApprovalItem

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**ExternalId** | Pointer to **NullableString** |  | [optional] 
**ExternalName** | Pointer to **NullableString** |  | [optional] 
**InternalId** | Pointer to **NullableString** |  | [optional] 
**ApprovedBy** | Pointer to **string** |  | [optional] 
**DeniedBy** | Pointer to **NullableString** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**ErrorMessage** | Pointer to **NullableString** |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 
**DateApproved** | Pointer to **time.Time** |  | [optional] 
**DateDenied** | Pointer to **NullableTime** |  | [optional] 
**Approval** | Pointer to [**ApprovalItemApprovalItemApproval**](ApprovalItemApprovalItemApproval.md) |  | [optional] 
**Reference** | Pointer to [**ApprovalItemApprovalItemReference**](ApprovalItemApprovalItemReference.md) |  | [optional] 
**Details** | Pointer to [**[]ApprovalItemApprovalItemDetailsInner**](ApprovalItemApprovalItemDetailsInner.md) | List of request change details associated with the approval item. Only present when the approval is for a reconfigure request.  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ApprovalItemApprovalItem{
    // Set fields directly
}
```

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
### InternalId (Nullable)

Use the Nullable wrapper methods:
- `obj.InternalId.IsSet()` — check if set
- `obj.InternalId.Get()` — get the inner value (returns pointer)
- `obj.InternalId.Set(&val)` — set the value
- `obj.InternalId.Unset()` — clear the value
### DeniedBy (Nullable)

Use the Nullable wrapper methods:
- `obj.DeniedBy.IsSet()` — check if set
- `obj.DeniedBy.Get()` — get the inner value (returns pointer)
- `obj.DeniedBy.Set(&val)` — set the value
- `obj.DeniedBy.Unset()` — clear the value
### ErrorMessage (Nullable)

Use the Nullable wrapper methods:
- `obj.ErrorMessage.IsSet()` — check if set
- `obj.ErrorMessage.Get()` — get the inner value (returns pointer)
- `obj.ErrorMessage.Set(&val)` — set the value
- `obj.ErrorMessage.Unset()` — clear the value
### DateDenied (Nullable)

Use the Nullable wrapper methods:
- `obj.DateDenied.IsSet()` — check if set
- `obj.DateDenied.Get()` — get the inner value (returns pointer)
- `obj.DateDenied.Set(&val)` — set the value
- `obj.DateDenied.Unset()` — clear the value
### Details (Nullable)

Use the Nullable wrapper methods:
- `obj.Details.IsSet()` — check if set
- `obj.Details.Get()` — get the inner value (returns pointer)
- `obj.Details.Set(&val)` — set the value
- `obj.Details.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


