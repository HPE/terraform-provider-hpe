# CreateSubnetRequestSubnet

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Name for the subnet | [optional] 
**Description** | Pointer to **string** | Description | [optional] 
**Type** | Pointer to [**CreateSubnetRequestSubnetType**](CreateSubnetRequestSubnetType.md) |  | [optional] 
**NetworkId** | Pointer to **int64** | The ID of the Network this subnet belongs to. Required when not using the nested route &#x60;/api/networks/{networkId}/subnets&#x60;. | [optional] 
**Config** | Pointer to **map[string]interface{}** | Configuration object. Settings vary by type. | [optional] 
**Cidr** | Pointer to **string** | Subnet CIDR. Used directly by subnet types with &#x60;cidrEditable&#x60; and &#x60;cidrRequired&#x60; (e.g. Google). For Azure subnets, this is derived from &#x60;config.subnetCidr&#x60; and does not need to be set explicitly. | [optional] 
**Active** | Pointer to **bool** | Activate (true) or disable (false) the subnet | [optional] 
**DhcpServer** | Pointer to **bool** | DHCP Server enabled subnet | [optional] 
**AllowStaticOverride** | Pointer to **bool** | Allow IP Override | [optional] 
**Pool** | Pointer to [**CreateSubnetRequestSubnetPool**](CreateSubnetRequestSubnetPool.md) |  | [optional] 
**Tenants** | Pointer to [**[]CreateSubnetRequestSubnetTenantsInner**](CreateSubnetRequestSubnetTenantsInner.md) | Array of tenant account ID objects that are allowed access | [optional] 
**Visibility** | Pointer to **string** | private or public | [optional] [default to "private"]
**Labels** | Pointer to **[]string** | Array of label strings, can be used for filtering. | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &CreateSubnetRequestSubnet{
    // Set fields directly
}
```

### Labels (Nullable)

Use the Nullable wrapper methods:
- `obj.Labels.IsSet()` — check if set
- `obj.Labels.Get()` — get the inner value (returns pointer)
- `obj.Labels.Set(&val)` — set the value
- `obj.Labels.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


