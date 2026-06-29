# GetCheckApps200ResponseOpenIncidentsInnerChecksInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Account** | Pointer to [**GetCheckApps200ResponseChecksInnerAccount**](GetCheckApps200ResponseChecksInnerAccount.md) |  | [optional] 
**Active** | Pointer to **bool** |  | [optional] 
**ApiKey** | Pointer to **string** |  | [optional] 
**Availability** | Pointer to **float32** |  | [optional] 
**CheckAgent** | Pointer to **NullableString** |  | [optional] 
**CheckInterval** | Pointer to **NullableInt64** |  | [optional] 
**CheckSpec** | Pointer to **NullableString** |  | [optional] 
**CheckType** | Pointer to [**GetCheckApps200ResponseChecksInnerCheckType**](GetCheckApps200ResponseChecksInnerCheckType.md) |  | [optional] 
**Config** | Pointer to [**GetCheckApps200ResponseChecksInnerConfig**](GetCheckApps200ResponseChecksInnerConfig.md) |  | [optional] 
**Container** | Pointer to [**GetCheckApps200ResponseChecksInnerContainer**](GetCheckApps200ResponseChecksInnerContainer.md) |  | [optional] 
**CreateIncident** | Pointer to **bool** |  | [optional] 
**Muted** | Pointer to **bool** |  | [optional] 
**CreatedBy** | Pointer to [**GetCheckApps200ResponseChecksInnerCreatedBy**](GetCheckApps200ResponseChecksInnerCreatedBy.md) |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**EndDate** | Pointer to **NullableTime** |  | [optional] 
**Health** | Pointer to **int64** |  | [optional] 
**InUptime** | Pointer to **bool** |  | [optional] 
**LastBoxStats** | Pointer to **NullableString** |  | [optional] 
**LastCheckStatus** | Pointer to **NullableString** |  | [optional] 
**LastError** | Pointer to **NullableString** |  | [optional] 
**LastErrorDate** | Pointer to **NullableTime** |  | [optional] 
**LastMessage** | Pointer to **NullableString** |  | [optional] 
**LastMetric** | Pointer to **NullableString** |  | [optional] 
**LastRunDate** | Pointer to **NullableTime** |  | [optional] 
**LastStats** | Pointer to **NullableString** |  | [optional] 
**LastSuccessDate** | Pointer to **NullableTime** |  | [optional] 
**LastTimer** | Pointer to **NullableInt64** |  | [optional] 
**LastUpdated** | Pointer to **NullableTime** |  | [optional] 
**LastWarningDate** | Pointer to **NullableTime** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**NextRunDate** | Pointer to **NullableTime** |  | [optional] 
**OutageTime** | Pointer to **int64** |  | [optional] 
**Severity** | Pointer to **string** |  | [optional] 
**StartDate** | Pointer to **NullableTime** |  | [optional] 
**Deleted** | Pointer to **bool** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetCheckApps200ResponseOpenIncidentsInnerChecksInner{
    // Set fields directly
}
```

### CheckAgent (Nullable)

Use the Nullable wrapper methods:
- `obj.CheckAgent.IsSet()` — check if set
- `obj.CheckAgent.Get()` — get the inner value (returns pointer)
- `obj.CheckAgent.Set(&val)` — set the value
- `obj.CheckAgent.Unset()` — clear the value
### CheckInterval (Nullable)

Use the Nullable wrapper methods:
- `obj.CheckInterval.IsSet()` — check if set
- `obj.CheckInterval.Get()` — get the inner value (returns pointer)
- `obj.CheckInterval.Set(&val)` — set the value
- `obj.CheckInterval.Unset()` — clear the value
### CheckSpec (Nullable)

Use the Nullable wrapper methods:
- `obj.CheckSpec.IsSet()` — check if set
- `obj.CheckSpec.Get()` — get the inner value (returns pointer)
- `obj.CheckSpec.Set(&val)` — set the value
- `obj.CheckSpec.Unset()` — clear the value
### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### EndDate (Nullable)

Use the Nullable wrapper methods:
- `obj.EndDate.IsSet()` — check if set
- `obj.EndDate.Get()` — get the inner value (returns pointer)
- `obj.EndDate.Set(&val)` — set the value
- `obj.EndDate.Unset()` — clear the value
### LastBoxStats (Nullable)

Use the Nullable wrapper methods:
- `obj.LastBoxStats.IsSet()` — check if set
- `obj.LastBoxStats.Get()` — get the inner value (returns pointer)
- `obj.LastBoxStats.Set(&val)` — set the value
- `obj.LastBoxStats.Unset()` — clear the value
### LastCheckStatus (Nullable)

Use the Nullable wrapper methods:
- `obj.LastCheckStatus.IsSet()` — check if set
- `obj.LastCheckStatus.Get()` — get the inner value (returns pointer)
- `obj.LastCheckStatus.Set(&val)` — set the value
- `obj.LastCheckStatus.Unset()` — clear the value
### LastError (Nullable)

Use the Nullable wrapper methods:
- `obj.LastError.IsSet()` — check if set
- `obj.LastError.Get()` — get the inner value (returns pointer)
- `obj.LastError.Set(&val)` — set the value
- `obj.LastError.Unset()` — clear the value
### LastErrorDate (Nullable)

Use the Nullable wrapper methods:
- `obj.LastErrorDate.IsSet()` — check if set
- `obj.LastErrorDate.Get()` — get the inner value (returns pointer)
- `obj.LastErrorDate.Set(&val)` — set the value
- `obj.LastErrorDate.Unset()` — clear the value
### LastMessage (Nullable)

Use the Nullable wrapper methods:
- `obj.LastMessage.IsSet()` — check if set
- `obj.LastMessage.Get()` — get the inner value (returns pointer)
- `obj.LastMessage.Set(&val)` — set the value
- `obj.LastMessage.Unset()` — clear the value
### LastMetric (Nullable)

Use the Nullable wrapper methods:
- `obj.LastMetric.IsSet()` — check if set
- `obj.LastMetric.Get()` — get the inner value (returns pointer)
- `obj.LastMetric.Set(&val)` — set the value
- `obj.LastMetric.Unset()` — clear the value
### LastRunDate (Nullable)

Use the Nullable wrapper methods:
- `obj.LastRunDate.IsSet()` — check if set
- `obj.LastRunDate.Get()` — get the inner value (returns pointer)
- `obj.LastRunDate.Set(&val)` — set the value
- `obj.LastRunDate.Unset()` — clear the value
### LastStats (Nullable)

Use the Nullable wrapper methods:
- `obj.LastStats.IsSet()` — check if set
- `obj.LastStats.Get()` — get the inner value (returns pointer)
- `obj.LastStats.Set(&val)` — set the value
- `obj.LastStats.Unset()` — clear the value
### LastSuccessDate (Nullable)

Use the Nullable wrapper methods:
- `obj.LastSuccessDate.IsSet()` — check if set
- `obj.LastSuccessDate.Get()` — get the inner value (returns pointer)
- `obj.LastSuccessDate.Set(&val)` — set the value
- `obj.LastSuccessDate.Unset()` — clear the value
### LastTimer (Nullable)

Use the Nullable wrapper methods:
- `obj.LastTimer.IsSet()` — check if set
- `obj.LastTimer.Get()` — get the inner value (returns pointer)
- `obj.LastTimer.Set(&val)` — set the value
- `obj.LastTimer.Unset()` — clear the value
### LastUpdated (Nullable)

Use the Nullable wrapper methods:
- `obj.LastUpdated.IsSet()` — check if set
- `obj.LastUpdated.Get()` — get the inner value (returns pointer)
- `obj.LastUpdated.Set(&val)` — set the value
- `obj.LastUpdated.Unset()` — clear the value
### LastWarningDate (Nullable)

Use the Nullable wrapper methods:
- `obj.LastWarningDate.IsSet()` — check if set
- `obj.LastWarningDate.Get()` — get the inner value (returns pointer)
- `obj.LastWarningDate.Set(&val)` — set the value
- `obj.LastWarningDate.Unset()` — clear the value
### NextRunDate (Nullable)

Use the Nullable wrapper methods:
- `obj.NextRunDate.IsSet()` — check if set
- `obj.NextRunDate.Get()` — get the inner value (returns pointer)
- `obj.NextRunDate.Set(&val)` — set the value
- `obj.NextRunDate.Unset()` — clear the value
### StartDate (Nullable)

Use the Nullable wrapper methods:
- `obj.StartDate.IsSet()` — check if set
- `obj.StartDate.Get()` — get the inner value (returns pointer)
- `obj.StartDate.Set(&val)` — set the value
- `obj.StartDate.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


