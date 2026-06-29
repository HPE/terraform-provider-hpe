# HealthElastic

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Success** | Pointer to **bool** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**Master** | Pointer to [**HealthElasticMaster**](HealthElasticMaster.md) |  | [optional] 
**Nodes** | Pointer to [**[]HealthElasticNodesInner**](HealthElasticNodesInner.md) |  | [optional] 
**Stats** | Pointer to [**HealthElasticStats**](HealthElasticStats.md) |  | [optional] 
**Indices** | Pointer to **[]map[string]interface{}** |  | [optional] 
**BadIndices** | Pointer to **[]map[string]interface{}** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &HealthElastic{
    // Set fields directly
}
```

### Indices (Nullable)

Use the Nullable wrapper methods:
- `obj.Indices.IsSet()` — check if set
- `obj.Indices.Get()` — get the inner value (returns pointer)
- `obj.Indices.Set(&val)` — set the value
- `obj.Indices.Unset()` — clear the value
### BadIndices (Nullable)

Use the Nullable wrapper methods:
- `obj.BadIndices.IsSet()` — check if set
- `obj.BadIndices.Get()` — get the inner value (returns pointer)
- `obj.BadIndices.Set(&val)` — set the value
- `obj.BadIndices.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


