# AwsResourcePoolConfig

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CidrBlock** | Pointer to **string** | Provide the base CIDR Block to use for this VPC (must be between a /16 and /28 Block) | [optional] 
**Tenancy** | Pointer to **string** | default or dedicated | [optional] [default to "default"]

## Usage

Instantiate with a Go composite literal:

```go
obj := &AwsResourcePoolConfig{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


