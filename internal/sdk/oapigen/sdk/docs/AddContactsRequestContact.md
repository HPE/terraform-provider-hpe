# AddContactsRequestContact

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | Unique name scoped to your account for the contact | 
**EmailAddress** | Pointer to **string** | Email notification address | [optional] 
**SmsAddress** | Pointer to **string** | SMS notification address | [optional] 
**SlackHook** | Pointer to **string** | Slack Hook | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddContactsRequestContact{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


