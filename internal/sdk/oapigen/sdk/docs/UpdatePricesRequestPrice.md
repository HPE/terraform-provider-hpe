# UpdatePricesRequestPrice

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Price name | [optional] 
**Code** | Pointer to **string** | Price code, must be unique | [optional] 
**Account** | Pointer to [**UpdatePricesRequestPriceAccount**](UpdatePricesRequestPriceAccount.md) |  | [optional] 
**PriceType** | Pointer to **string** | Restricts query to only load only prices with specified priceType. * &#x60;fixed&#x60; - Everything * &#x60;compute&#x60; - Memory + CPU * &#x60;memory&#x60; - Memory * &#x60;cores&#x60; - Cores * &#x60;storage&#x60; - Storage * &#x60;datastore&#x60; - Datastore * &#x60;platform&#x60; - Platform * &#x60;software&#x60; - Software * &#x60;load_balancer&#x60; - Load Balancer * &#x60;load_balancer_virtual_server&#x60; - Load Balancer Virtual Server  | [optional] 
**PriceUnit** | Pointer to **string** | The unit of pricing | [optional] 
**IncurCharges** | Pointer to **string** | Indicates when to incur charge | [optional] 
**Currency** | Pointer to **string** | ISO Currency code | [optional] 
**Cost** | Pointer to **float32** | Cost | [optional] 
**MarkupType** | Pointer to **string** | Price adjustment type | [optional] 
**Markup** | Pointer to **float32** | Amount for &#x60;fixed&#x60; price adjustment type | [optional] 
**MarkupPercent** | Pointer to **float32** | Percent for &#x60;percent&#x60; price adjustment type | [optional] 
**CustomPrice** | Pointer to **float32** | Custom price for &#x60;custom&#x60; price adjustment type | [optional] 
**Platform** | Pointer to **string** | Platform.  Required for &#x60;platform&#x60; price type | [optional] 
**Software** | Pointer to **string** | Software.  Required for software price type | [optional] 
**VolumeType** | Pointer to [**UpdatePricesRequestPriceVolumeType**](UpdatePricesRequestPriceVolumeType.md) |  | [optional] 
**Datastore** | Pointer to [**UpdatePricesRequestPriceDatastore**](UpdatePricesRequestPriceDatastore.md) |  | [optional] 
**CrossCloudApply** | Pointer to **bool** | Apply price across clouds, optional true/false flag for datastore price type | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdatePricesRequestPrice{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


