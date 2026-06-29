# AddCloudResourcePoolRequestResourcePoolConfig

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CidrBlock** | Pointer to **string** | Provide the base CIDR Block to use for this VPC (must be between a /16 and /28 Block) | [optional] 
**Tenancy** | Pointer to **string** | default or dedicated | [optional] [default to "default"]
**Managers** | Pointer to **[]string** | Array of manager usernames | [optional] 
**Developers** | Pointer to **[]string** | Array of developer usernames | [optional] 
**Auditors** | Pointer to **[]string** | Array of auditor usernames | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddCloudResourcePoolRequestResourcePoolConfig{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


