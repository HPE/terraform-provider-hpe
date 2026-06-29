# UpdateIntegrationsRequestOneOf1IntegrationConfig

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**DefaultBranch** | Pointer to **string** | default branch | [optional] 
**AnsiblePlaybooks** | Pointer to **string** | Playbooks path | [optional] 
**AnsibleRoles** | Pointer to **string** | Roles path | [optional] 
**AnsibleGroupVars** | Pointer to **string** | Group variables path | [optional] 
**AnsibleHostVars** | Pointer to **string** | Host variables path | [optional] 
**AnsibleGalaxyEnabled** | Pointer to **bool** | Use Ansible Galaxy | [optional] 
**AnsibleVerbose** | Pointer to **bool** | Use verbose logging | [optional] 
**AnsibleCommandBus** | Pointer to **bool** | Use Morpheus Agent Command Bus | [optional] 
**CacheEnabled** | Pointer to **bool** | Enable Git repository caching | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateIntegrationsRequestOneOf1IntegrationConfig{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


