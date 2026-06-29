# GetCypherKey200Response

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Success** | Pointer to **bool** |  | [optional] 
**Data** | Pointer to [**GetCypherKey200ResponseAllOfData**](GetCypherKey200ResponseAllOfData.md) |  | [optional] 
**Type** | Pointer to **string** | Type of data that was written to the key | [optional] 
**LeaseDuration** | Pointer to **int32** | Lease duration in seconds, 0 means no expiry. | [optional] 
**Cypher** | Pointer to [**GetCypherKey200ResponseAllOfCypher**](GetCypherKey200ResponseAllOfCypher.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetCypherKey200Response{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


