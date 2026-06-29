# CreateNetworkDomainRequestNetworkDomain

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Name | [optional] 
**Description** | Pointer to **string** | Description | [optional] 
**DisplayName** | Pointer to **string** | Overrides displayed name in domain selection components. Useful if using many OU Paths. | [optional] 
**PublicZone** | Pointer to **bool** | Public Zone | [optional] [default to false]
**TaskSetId** | Pointer to **int64** | Workflow ID. Workflows can be applied to an instance when associated with a domain. Useful for custom domain related scripting. (Important if wanting to ensure the computer is removed from the domain using teardown phased workflows.)  | [optional] 
**Active** | Pointer to **bool** | Active | [optional] 
**DomainController** | Pointer to **bool** | Join Domain Controller | [optional] [default to true]
**DomainUsername** | Pointer to **string** | Domain Username | [optional] 
**DomainPassword** | Pointer to **string** | Domain Password | [optional] 
**DcServer** | Pointer to **string** | DC Server. If specified, must be the server name and not an IP Address | [optional] 
**OuPath** | Pointer to **string** | OU Path. (i.e. &#39;OU&#x3D;staging,DC&#x3D;ad,DC&#x3D;yourdomain,DC&#x3D;com&#39;) | [optional] 
**GuestUsername** | Pointer to **string** | Guest Username. If set, will change the instances RPC Service User after joining a Domain. | [optional] 
**GuestPassword** | Pointer to **string** | Guest Password | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &CreateNetworkDomainRequestNetworkDomain{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


