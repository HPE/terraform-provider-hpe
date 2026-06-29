# UpdateUserSettingsRequestUser

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Username** | Pointer to **string** | Username | [optional] 
**Email** | Pointer to **string** | Email | [optional] 
**FirstName** | Pointer to **string** | First Name | [optional] 
**LastName** | Pointer to **string** | Last Name | [optional] 
**Password** | Pointer to **string** | Change your password | [optional] 
**LinuxUsername** | Pointer to **string** | Linux Username | [optional] 
**LinuxPassword** | Pointer to **string** | Linux Password | [optional] 
**LinuxKeyPairId** | Pointer to **int64** | Linux Key Pair ID | [optional] 
**WindowsUsername** | Pointer to **string** | Windows Username | [optional] 
**WindowsPassword** | Pointer to **string** | Windows Password | [optional] 
**ReceiveNotifications** | Pointer to **bool** | Receive Notifications (true or false) | [optional] 
**DefaultGroup** | Pointer to [**UpdateUserSettingsRequestUserDefaultGroup**](UpdateUserSettingsRequestUserDefaultGroup.md) |  | [optional] 
**DefaultCloud** | Pointer to [**UpdateUserSettingsRequestUserDefaultCloud**](UpdateUserSettingsRequestUserDefaultCloud.md) |  | [optional] 
**DefaultPersona** | Pointer to [**UpdateUserSettingsRequestUserDefaultPersona**](UpdateUserSettingsRequestUserDefaultPersona.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateUserSettingsRequestUser{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


