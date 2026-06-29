# CheckApp

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Account** | Pointer to [**AddCheckApps200ResponseAllOfCheckAppAccount**](AddCheckApps200ResponseAllOfCheckAppAccount.md) |  | [optional] 
**Active** | Pointer to **bool** |  | [optional] 
**App** | Pointer to [**AddCheckApps200ResponseAllOfCheckAppApp**](AddCheckApps200ResponseAllOfCheckAppApp.md) |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**InUptime** | Pointer to **bool** |  | [optional] 
**LastCheckStatus** | Pointer to **NullableString** |  | [optional] 
**LastWarningDate** | Pointer to **NullableTime** |  | [optional] 
**LastErrorDate** | Pointer to **NullableTime** |  | [optional] 
**LastSuccessDate** | Pointer to **NullableTime** |  | [optional] 
**LastRunDate** | Pointer to **NullableTime** |  | [optional] 
**LastError** | Pointer to **NullableString** |  | [optional] 
**LastTimer** | Pointer to **int64** |  | [optional] 
**Health** | Pointer to **int64** |  | [optional] 
**History** | Pointer to **NullableString** |  | [optional] 
**Severity** | Pointer to **string** |  | [optional] 
**CreateIncident** | Pointer to **bool** |  | [optional] 
**Muted** | Pointer to **bool** |  | [optional] 
**CreatedBy** | Pointer to [**AddCheckApps200ResponseAllOfCheckAppCreatedBy**](AddCheckApps200ResponseAllOfCheckAppCreatedBy.md) |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 
**Availability** | Pointer to **NullableString** |  | [optional] 
**Checks** | Pointer to **[]int64** |  | [optional] 
**CheckGroups** | Pointer to **[]int64** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &CheckApp{
    // Set fields directly
}
```

### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### LastCheckStatus (Nullable)

Use the Nullable wrapper methods:
- `obj.LastCheckStatus.IsSet()` — check if set
- `obj.LastCheckStatus.Get()` — get the inner value (returns pointer)
- `obj.LastCheckStatus.Set(&val)` — set the value
- `obj.LastCheckStatus.Unset()` — clear the value
### LastWarningDate (Nullable)

Use the Nullable wrapper methods:
- `obj.LastWarningDate.IsSet()` — check if set
- `obj.LastWarningDate.Get()` — get the inner value (returns pointer)
- `obj.LastWarningDate.Set(&val)` — set the value
- `obj.LastWarningDate.Unset()` — clear the value
### LastErrorDate (Nullable)

Use the Nullable wrapper methods:
- `obj.LastErrorDate.IsSet()` — check if set
- `obj.LastErrorDate.Get()` — get the inner value (returns pointer)
- `obj.LastErrorDate.Set(&val)` — set the value
- `obj.LastErrorDate.Unset()` — clear the value
### LastSuccessDate (Nullable)

Use the Nullable wrapper methods:
- `obj.LastSuccessDate.IsSet()` — check if set
- `obj.LastSuccessDate.Get()` — get the inner value (returns pointer)
- `obj.LastSuccessDate.Set(&val)` — set the value
- `obj.LastSuccessDate.Unset()` — clear the value
### LastRunDate (Nullable)

Use the Nullable wrapper methods:
- `obj.LastRunDate.IsSet()` — check if set
- `obj.LastRunDate.Get()` — get the inner value (returns pointer)
- `obj.LastRunDate.Set(&val)` — set the value
- `obj.LastRunDate.Unset()` — clear the value
### LastError (Nullable)

Use the Nullable wrapper methods:
- `obj.LastError.IsSet()` — check if set
- `obj.LastError.Get()` — get the inner value (returns pointer)
- `obj.LastError.Set(&val)` — set the value
- `obj.LastError.Unset()` — clear the value
### History (Nullable)

Use the Nullable wrapper methods:
- `obj.History.IsSet()` — check if set
- `obj.History.Get()` — get the inner value (returns pointer)
- `obj.History.Set(&val)` — set the value
- `obj.History.Unset()` — clear the value
### Availability (Nullable)

Use the Nullable wrapper methods:
- `obj.Availability.IsSet()` — check if set
- `obj.Availability.Get()` — get the inner value (returns pointer)
- `obj.Availability.Set(&val)` — set the value
- `obj.Availability.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


