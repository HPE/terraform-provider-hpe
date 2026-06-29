# UpdateHostRequestServer

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Unique name scoped to your account for the server. | [optional] 
**Description** | Pointer to **string** | Optional description field. | [optional] 
**Enabled** | Pointer to **bool** | Flag to determine if a host can be selected for provisioning | [optional] [default to true]
**ManageInternalFirewall** | Pointer to **bool** | Flag to enable/disable managment of internal firewall | [optional] [default to true]
**EnableLogs** | Pointer to **bool** | Flag to enable/disable logs | [optional] [default to true]
**SshUsername** | Pointer to **string** | SSH Username | [optional] 
**SshPassword** | Pointer to **NullableString** | SSH Password | [optional] 
**SshKeyPair** | Pointer to [**UpdateHostRequestServerSshKeyPair**](UpdateHostRequestServerSshKeyPair.md) |  | [optional] 
**PowerScheduleType** | Pointer to **int64** | Power schedule ID. | [optional] 
**Labels** | Pointer to **[]string** |  | [optional] 
**ServerOs** | Pointer to [**UpdateHostRequestServerServerOs**](UpdateHostRequestServerServerOs.md) |  | [optional] 
**Tags** | Pointer to [**[]UpdateHostRequestServerTagsInner**](UpdateHostRequestServerTagsInner.md) | Metadata tags, Array of objects having a name and value. | [optional] 
**AddTags** | Pointer to [**[]UpdateHostRequestServerAddTagsInner**](UpdateHostRequestServerAddTagsInner.md) | Add or update value of Metadata tags, Array of objects having a name and value. | [optional] 
**RemoveTags** | Pointer to [**[]UpdateHostRequestServerRemoveTagsInner**](UpdateHostRequestServerRemoveTagsInner.md) | Remove Metadata tags, Array of objects having a name and an optional value. If value is passed, it must match to be removed. | [optional] 
**GuestConsoleType** | Pointer to **string** | The Type of guest console this server provides such as disabled, vnc, rdp, ssh | [optional] 
**GuestConsoleUsername** | Pointer to **string** | The optional guest console username if you don&#39;t want to use the user defaults | [optional] 
**GuestConsolePassword** | Pointer to **string** | The optional guest console password if not using the accessing users creds | [optional] 
**GuestConsolePort** | Pointer to **string** | The port the guest console is being accessed from | [optional] 
**GuestConsolePreferred** | Pointer to **bool** | Can turn off guest console preferences on server in favor of hypervisor console | [optional] [default to true]
**Config** | Pointer to [**UpdateHostRequestServerConfig**](UpdateHostRequestServerConfig.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateHostRequestServer{
    // Set fields directly
}
```

### SshPassword (Nullable)

Use the Nullable wrapper methods:
- `obj.SshPassword.IsSet()` — check if set
- `obj.SshPassword.Get()` — get the inner value (returns pointer)
- `obj.SshPassword.Set(&val)` — set the value
- `obj.SshPassword.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


