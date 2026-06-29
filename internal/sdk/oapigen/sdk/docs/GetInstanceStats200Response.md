# GetInstanceStats200Response

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**InstanceStats** | Pointer to [**GetInstanceStats200ResponseInstanceStats**](GetInstanceStats200ResponseInstanceStats.md) |  | [optional] 
**ZoneIds** | Pointer to **[]int64** | Array of Cloud IDs that are included in the stats. By default all the clouds the user has access to are returned. | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetInstanceStats200Response{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


