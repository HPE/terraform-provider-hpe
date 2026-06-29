# AddUserRequestUser

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**FirstName** | Pointer to **string** | The user&#39;s first name (optional) | [optional] 
**LastName** | Pointer to **string** | The user&#39;s last name (optional) | [optional] 
**Username** | **string** | Username (unique per tenant). | 
**Email** | **string** | Email address | 
**Password** | **string** | Password to apply to the user | 
**Roles** | [**[]AddUserRequestUserRolesInner**](AddUserRequestUserRolesInner.md) | Array of objects with id of the role(s) to assign to the user. | 
**ReceiveNotifications** | Pointer to **bool** | Receive Notifications? | [optional] [default to true]
**LinuxUsername** | Pointer to **string** | Linux Username, user settings for provisioning | [optional] 
**LinuxPassword** | Pointer to **string** | Linux Password, user settings for provisioning | [optional] 
**LinuxKeyPairId** | Pointer to **int64** | Linux SSH Key, user settings for provisioning | [optional] 
**WindowsUsername** | Pointer to **string** | Windows Username, user settings for provisioning | [optional] 
**WindowsPassword** | Pointer to **string** | Windows Password, user settings for provisioning | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddUserRequestUser{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


