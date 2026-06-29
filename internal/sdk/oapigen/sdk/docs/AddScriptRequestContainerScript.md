# AddScriptRequestContainerScript

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | Script name | 
**Labels** | Pointer to **[]string** | Array of label strings, can be used for filtering. | [optional] 
**Category** | Pointer to **string** | Script category | [optional] 
**ScriptVersion** | Pointer to **string** | Version of the script | [optional] [default to "1"]
**ScriptPhase** | Pointer to **string** | Phase for the script, provision, start, etc. | [optional] 
**ScriptType** | Pointer to **string** | Type for the script | [optional] [default to "bash"]
**Script** | Pointer to **string** | Script content, that is, the code itself. | [optional] 
**RunAsUser** | Pointer to **string** | Run as a specific user. | [optional] 
**SudoUser** | Pointer to **bool** | Sudo, whether or not to run with sudo. | [optional] [default to false]

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddScriptRequestContainerScript{
    // Set fields directly
}
```

### Labels (Nullable)

Use the Nullable wrapper methods:
- `obj.Labels.IsSet()` — check if set
- `obj.Labels.Get()` — get the inner value (returns pointer)
- `obj.Labels.Set(&val)` — set the value
- `obj.Labels.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


