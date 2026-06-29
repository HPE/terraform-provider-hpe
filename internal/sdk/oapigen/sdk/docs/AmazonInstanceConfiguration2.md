# AmazonInstanceConfiguration2

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**NoAgent** | Pointer to **NullableBool** | Skipping Agent installation will result in a lack of logging and guest operating system statistics. Automation scripts may also be adversely affected. | [optional] [default to false]
**CreateUser** | Pointer to **NullableBool** | Create user | [optional] [default to false]
**IsEC2** | Pointer to **string** | Amazon Cloud Type | [optional] [default to "false"]
**AvailabilityId** | Pointer to **string** | Amazon Zone | [optional] 
**SecurityId** | Pointer to **string** | Security Group | [optional] 
**PublicIpType** | Pointer to **string** | Public IP | [optional] 
**InstanceProfile** | Pointer to **string** | IAM Profile | [optional] 
**KmsKeyId** | Pointer to **string** | KMS Key ID | [optional] 
**ResourcePoolId** | Pointer to **string** | Resource Pool ID | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AmazonInstanceConfiguration2{
    // Set fields directly
}
```

### NoAgent (Nullable)

Use the Nullable wrapper methods:
- `obj.NoAgent.IsSet()` — check if set
- `obj.NoAgent.Get()` — get the inner value (returns pointer)
- `obj.NoAgent.Set(&val)` — set the value
- `obj.NoAgent.Unset()` — clear the value
### CreateUser (Nullable)

Use the Nullable wrapper methods:
- `obj.CreateUser.IsSet()` — check if set
- `obj.CreateUser.Get()` — get the inner value (returns pointer)
- `obj.CreateUser.Set(&val)` — set the value
- `obj.CreateUser.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


