# InstancesCloneImage

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**TemplateName** | Pointer to **string** | Image Template Name | [optional] [default to "{server.name}-{timestamp}"]
**ZoneFolder** | Pointer to **string** | Zone Folder externalId. This is required for VMware | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &InstancesCloneImage{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


