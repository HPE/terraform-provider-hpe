# ListHealth200ResponseAllOfHealthDatabase

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Success** | Pointer to **bool** |  | [optional] 
**ConnectionList** | Pointer to **[]map[string]interface{}** |  | [optional] 
**BusyConnections** | Pointer to **[]string** |  | [optional] 
**MaxConnections** | Pointer to **int64** |  | [optional] 
**MaxUsedConnections** | Pointer to **int64** |  | [optional] 
**UsedConnections** | Pointer to **int64** |  | [optional] 
**AbortedConnections** | Pointer to **int64** |  | [optional] 
**InnodbStatus** | Pointer to **NullableString** |  | [optional] 
**Stats** | Pointer to [**ListHealth200ResponseAllOfHealthDatabaseStats**](ListHealth200ResponseAllOfHealthDatabaseStats.md) |  | [optional] 
**Scans** | Pointer to [**ListHealth200ResponseAllOfHealthDatabaseScans**](ListHealth200ResponseAllOfHealthDatabaseScans.md) |  | [optional] 
**SlowQueries** | Pointer to [**[]ListHealth200ResponseAllOfHealthDatabaseSlowQueriesInner**](ListHealth200ResponseAllOfHealthDatabaseSlowQueriesInner.md) |  | [optional] 
**InnodbStats** | Pointer to [**ListHealth200ResponseAllOfHealthDatabaseInnodbStats**](ListHealth200ResponseAllOfHealthDatabaseInnodbStats.md) |  | [optional] 
**ScanPercent** | Pointer to **float32** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListHealth200ResponseAllOfHealthDatabase{
    // Set fields directly
}
```

### ConnectionList (Nullable)

Use the Nullable wrapper methods:
- `obj.ConnectionList.IsSet()` — check if set
- `obj.ConnectionList.Get()` — get the inner value (returns pointer)
- `obj.ConnectionList.Set(&val)` — set the value
- `obj.ConnectionList.Unset()` — clear the value
### InnodbStatus (Nullable)

Use the Nullable wrapper methods:
- `obj.InnodbStatus.IsSet()` — check if set
- `obj.InnodbStatus.Get()` — get the inner value (returns pointer)
- `obj.InnodbStatus.Set(&val)` — set the value
- `obj.InnodbStatus.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


