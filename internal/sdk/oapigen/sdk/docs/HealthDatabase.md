# HealthDatabase

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
**Stats** | Pointer to [**HealthDatabaseStats**](HealthDatabaseStats.md) |  | [optional] 
**Scans** | Pointer to [**HealthDatabaseScans**](HealthDatabaseScans.md) |  | [optional] 
**SlowQueries** | Pointer to [**[]HealthDatabaseSlowQueriesInner**](HealthDatabaseSlowQueriesInner.md) |  | [optional] 
**InnodbStats** | Pointer to [**HealthDatabaseInnodbStats**](HealthDatabaseInnodbStats.md) |  | [optional] 
**ScanPercent** | Pointer to **float32** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &HealthDatabase{
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


