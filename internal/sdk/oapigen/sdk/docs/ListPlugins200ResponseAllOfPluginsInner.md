# ListPlugins200ResponseAllOfPluginsInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Code** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **string** |  | [optional] 
**Version** | Pointer to **string** |  | [optional] 
**Enabled** | Pointer to **bool** |  | [optional] 
**Author** | Pointer to **NullableString** |  | [optional] 
**WebsiteUrl** | Pointer to **NullableString** |  | [optional] 
**SourceCodeLocationUrl** | Pointer to **NullableString** |  | [optional] 
**IssueTrackerUrl** | Pointer to **NullableString** |  | [optional] 
**Valid** | Pointer to **bool** |  | [optional] 
**HasValidUpdate** | Pointer to **bool** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**StatusMessage** | Pointer to **NullableString** |  | [optional] 
**Providers** | Pointer to [**[]ListPlugins200ResponseAllOfPluginsInnerProvidersInner**](ListPlugins200ResponseAllOfPluginsInnerProvidersInner.md) |  | [optional] 
**Config** | Pointer to **map[string]interface{}** |  | [optional] 
**OptionTypes** | Pointer to [**[]ListPlugins200ResponseAllOfPluginsInnerOptionTypesInner**](ListPlugins200ResponseAllOfPluginsInnerOptionTypesInner.md) |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListPlugins200ResponseAllOfPluginsInner{
    // Set fields directly
}
```

### Author (Nullable)

Use the Nullable wrapper methods:
- `obj.Author.IsSet()` — check if set
- `obj.Author.Get()` — get the inner value (returns pointer)
- `obj.Author.Set(&val)` — set the value
- `obj.Author.Unset()` — clear the value
### WebsiteUrl (Nullable)

Use the Nullable wrapper methods:
- `obj.WebsiteUrl.IsSet()` — check if set
- `obj.WebsiteUrl.Get()` — get the inner value (returns pointer)
- `obj.WebsiteUrl.Set(&val)` — set the value
- `obj.WebsiteUrl.Unset()` — clear the value
### SourceCodeLocationUrl (Nullable)

Use the Nullable wrapper methods:
- `obj.SourceCodeLocationUrl.IsSet()` — check if set
- `obj.SourceCodeLocationUrl.Get()` — get the inner value (returns pointer)
- `obj.SourceCodeLocationUrl.Set(&val)` — set the value
- `obj.SourceCodeLocationUrl.Unset()` — clear the value
### IssueTrackerUrl (Nullable)

Use the Nullable wrapper methods:
- `obj.IssueTrackerUrl.IsSet()` — check if set
- `obj.IssueTrackerUrl.Get()` — get the inner value (returns pointer)
- `obj.IssueTrackerUrl.Set(&val)` — set the value
- `obj.IssueTrackerUrl.Unset()` — clear the value
### StatusMessage (Nullable)

Use the Nullable wrapper methods:
- `obj.StatusMessage.IsSet()` — check if set
- `obj.StatusMessage.Get()` — get the inner value (returns pointer)
- `obj.StatusMessage.Set(&val)` — set the value
- `obj.StatusMessage.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


