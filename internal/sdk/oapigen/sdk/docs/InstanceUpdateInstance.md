# InstanceUpdateInstance

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Unique name scoped to your account for the instance. | [optional] 
**Description** | Pointer to **string** | Optional description field. | [optional] 
**InstanceContext** | Pointer to **string** | Environment | [optional] 
**Labels** | Pointer to **[]string** | Array of strings (keywords). | [optional] 
**Tags** | Pointer to [**[]InstanceUpdateInstanceTagsInner**](InstanceUpdateInstanceTagsInner.md) | Metadata tags, Array of objects having a name and value. | [optional] 
**AddTags** | Pointer to [**[]InstanceUpdateInstanceAddTagsInner**](InstanceUpdateInstanceAddTagsInner.md) | Add or update value of Metadata tags, Array of objects having a name and value. | [optional] 
**RemoveTags** | Pointer to [**[]InstanceUpdateInstanceRemoveTagsInner**](InstanceUpdateInstanceRemoveTagsInner.md) | Remove Metadata tags, Array of objects having a name and an optional value. If value is passed, it must match to be removed. | [optional] 
**PowerScheduleType** | Pointer to **int64** | Power schedule ID. | [optional] 
**Site** | Pointer to [**InstanceUpdateInstanceSite**](InstanceUpdateInstanceSite.md) |  | [optional] 
**OwnerId** | Pointer to **int64** | User ID, can be used to change instance owner. | [optional] 
**DisplayName** | Pointer to **string** | Name used in the UI for display | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &InstanceUpdateInstance{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


