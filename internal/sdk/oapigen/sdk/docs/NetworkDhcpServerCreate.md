# NetworkDhcpServerCreate

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ServerIpAddress** | **string** | Server Address for the DHCP Server | 
**LeaseTime** | **int64** | Optional lease time for the DHCP Server | [default to 86400]
**Name** | **string** | Name | 
**Config** | [**NetworkDhcpServerCreateConfig**](NetworkDhcpServerCreateConfig.md) |  | 

## Usage

Instantiate with a Go composite literal:

```go
obj := &NetworkDhcpServerCreate{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


