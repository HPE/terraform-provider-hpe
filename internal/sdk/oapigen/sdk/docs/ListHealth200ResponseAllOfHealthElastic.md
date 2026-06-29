# ListHealth200ResponseAllOfHealthElastic

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Success** | Pointer to **bool** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**Master** | Pointer to [**ListHealth200ResponseAllOfHealthElasticMaster**](ListHealth200ResponseAllOfHealthElasticMaster.md) |  | [optional] 
**Nodes** | Pointer to [**[]ListHealth200ResponseAllOfHealthElasticNodesInner**](ListHealth200ResponseAllOfHealthElasticNodesInner.md) |  | [optional] 
**Stats** | Pointer to [**ListHealth200ResponseAllOfHealthElasticStats**](ListHealth200ResponseAllOfHealthElasticStats.md) |  | [optional] 
**Indices** | Pointer to **[]map[string]interface{}** |  | [optional] 
**BadIndices** | Pointer to **[]map[string]interface{}** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListHealth200ResponseAllOfHealthElastic{
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


