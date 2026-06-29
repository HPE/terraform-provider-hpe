# Issue

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**AttachmentType** | Pointer to **string** |  | [optional] 
**App** | Pointer to **NullableString** |  | [optional] 
**Available** | Pointer to **bool** |  | [optional] 
**Check** | Pointer to **NullableString** |  | [optional] 
**CheckGroup** | Pointer to [**IssueCheckGroup**](IssueCheckGroup.md) |  | [optional] 
**CheckStatus** | Pointer to **NullableString** |  | [optional] 
**EndDate** | Pointer to **NullableTime** |  | [optional] 
**Health** | Pointer to **int64** |  | [optional] 
**InUptime** | Pointer to **bool** |  | [optional] 
**Incident** | Pointer to [**IssueIncident**](IssueIncident.md) |  | [optional] 
**LastCheckTime** | Pointer to **NullableTime** |  | [optional] 
**LastError** | Pointer to **NullableString** |  | [optional] 
**LastMessage** | Pointer to **NullableString** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Severity** | Pointer to **string** |  | [optional] 
**SeverityId** | Pointer to **int64** |  | [optional] 
**StartDate** | Pointer to **time.Time** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &Issue{
    // Set fields directly
}
```

### App (Nullable)

Use the Nullable wrapper methods:
- `obj.App.IsSet()` — check if set
- `obj.App.Get()` — get the inner value (returns pointer)
- `obj.App.Set(&val)` — set the value
- `obj.App.Unset()` — clear the value
### Check (Nullable)

Use the Nullable wrapper methods:
- `obj.Check.IsSet()` — check if set
- `obj.Check.Get()` — get the inner value (returns pointer)
- `obj.Check.Set(&val)` — set the value
- `obj.Check.Unset()` — clear the value
### CheckStatus (Nullable)

Use the Nullable wrapper methods:
- `obj.CheckStatus.IsSet()` — check if set
- `obj.CheckStatus.Get()` — get the inner value (returns pointer)
- `obj.CheckStatus.Set(&val)` — set the value
- `obj.CheckStatus.Unset()` — clear the value
### EndDate (Nullable)

Use the Nullable wrapper methods:
- `obj.EndDate.IsSet()` — check if set
- `obj.EndDate.Get()` — get the inner value (returns pointer)
- `obj.EndDate.Set(&val)` — set the value
- `obj.EndDate.Unset()` — clear the value
### LastCheckTime (Nullable)

Use the Nullable wrapper methods:
- `obj.LastCheckTime.IsSet()` — check if set
- `obj.LastCheckTime.Get()` — get the inner value (returns pointer)
- `obj.LastCheckTime.Set(&val)` — set the value
- `obj.LastCheckTime.Unset()` — clear the value
### LastError (Nullable)

Use the Nullable wrapper methods:
- `obj.LastError.IsSet()` — check if set
- `obj.LastError.Get()` — get the inner value (returns pointer)
- `obj.LastError.Set(&val)` — set the value
- `obj.LastError.Unset()` — clear the value
### LastMessage (Nullable)

Use the Nullable wrapper methods:
- `obj.LastMessage.IsSet()` — check if set
- `obj.LastMessage.Get()` — get the inner value (returns pointer)
- `obj.LastMessage.Set(&val)` — set the value
- `obj.LastMessage.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


