# ListEmailTemplates200ResponseAllOfEmailTemplatesInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** | The name of the email template. This is set by morpheus.  | [optional] 
**Code** | Pointer to **string** | A unique code for the email template. This code is used to reference the email template and as a reference of the templates type.  | [optional] 
**Owner** | Pointer to [**ListEmailTemplates200ResponseAllOfEmailTemplatesInnerOwner**](ListEmailTemplates200ResponseAllOfEmailTemplatesInnerOwner.md) |  | [optional] 
**Accounts** | Pointer to [**[]ListEmailTemplates200ResponseAllOfEmailTemplatesInnerAccountsInner**](ListEmailTemplates200ResponseAllOfEmailTemplatesInnerAccountsInner.md) |  | [optional] 
**Template** | Pointer to **string** | The email template. This is the actual email template that is sent to the user. This uses handlebars notation (not javascript)  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListEmailTemplates200ResponseAllOfEmailTemplatesInner{
    // Set fields directly
}
```

### Accounts (Nullable)

Use the Nullable wrapper methods:
- `obj.Accounts.IsSet()` — check if set
- `obj.Accounts.Get()` — get the inner value (returns pointer)
- `obj.Accounts.Set(&val)` — set the value
- `obj.Accounts.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


