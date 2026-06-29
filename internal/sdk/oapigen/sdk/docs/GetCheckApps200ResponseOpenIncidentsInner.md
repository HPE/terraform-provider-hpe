# GetCheckApps200ResponseOpenIncidentsInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Account** | Pointer to [**GetCheckApps200ResponseOpenIncidentsInnerAccount**](GetCheckApps200ResponseOpenIncidentsInnerAccount.md) |  | [optional] 
**App** | Pointer to **NullableString** |  | [optional] 
**AutoClose** | Pointer to **bool** |  | [optional] 
**ChannelId** | Pointer to **string** |  | [optional] 
**CheckGroups** | Pointer to [**[]GetCheckApps200ResponseOpenIncidentsInnerCheckGroupsInner**](GetCheckApps200ResponseOpenIncidentsInnerCheckGroupsInner.md) |  | [optional] 
**Checks** | Pointer to [**[]GetCheckApps200ResponseOpenIncidentsInnerChecksInner**](GetCheckApps200ResponseOpenIncidentsInnerChecksInner.md) |  | [optional] 
**Comment** | Pointer to **NullableString** |  | [optional] 
**DisplayName** | Pointer to **string** |  | [optional] 
**Duration** | Pointer to **NullableString** |  | [optional] 
**EndDate** | Pointer to **NullableTime** |  | [optional] 
**InUptime** | Pointer to **bool** |  | [optional] 
**Muted** | Pointer to **bool** |  | [optional] 
**LastCheckTime** | Pointer to **time.Time** |  | [optional] 
**LastError** | Pointer to **string** |  | [optional] 
**LastMessage** | Pointer to **NullableString** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Resolution** | Pointer to **NullableString** |  | [optional] 
**Severity** | Pointer to **string** |  | [optional] 
**SeverityId** | Pointer to **int64** |  | [optional] 
**StartDate** | Pointer to **time.Time** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**Visibility** | Pointer to **string** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetCheckApps200ResponseOpenIncidentsInner{
    // Set fields directly
}
```

### App (Nullable)

Use the Nullable wrapper methods:
- `obj.App.IsSet()` — check if set
- `obj.App.Get()` — get the inner value (returns pointer)
- `obj.App.Set(&val)` — set the value
- `obj.App.Unset()` — clear the value
### Comment (Nullable)

Use the Nullable wrapper methods:
- `obj.Comment.IsSet()` — check if set
- `obj.Comment.Get()` — get the inner value (returns pointer)
- `obj.Comment.Set(&val)` — set the value
- `obj.Comment.Unset()` — clear the value
### Duration (Nullable)

Use the Nullable wrapper methods:
- `obj.Duration.IsSet()` — check if set
- `obj.Duration.Get()` — get the inner value (returns pointer)
- `obj.Duration.Set(&val)` — set the value
- `obj.Duration.Unset()` — clear the value
### EndDate (Nullable)

Use the Nullable wrapper methods:
- `obj.EndDate.IsSet()` — check if set
- `obj.EndDate.Get()` — get the inner value (returns pointer)
- `obj.EndDate.Set(&val)` — set the value
- `obj.EndDate.Unset()` — clear the value
### LastMessage (Nullable)

Use the Nullable wrapper methods:
- `obj.LastMessage.IsSet()` — check if set
- `obj.LastMessage.Get()` — get the inner value (returns pointer)
- `obj.LastMessage.Set(&val)` — set the value
- `obj.LastMessage.Unset()` — clear the value
### Resolution (Nullable)

Use the Nullable wrapper methods:
- `obj.Resolution.IsSet()` — check if set
- `obj.Resolution.Get()` — get the inner value (returns pointer)
- `obj.Resolution.Set(&val)` — set the value
- `obj.Resolution.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


