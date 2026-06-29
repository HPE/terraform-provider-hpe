# ListInvoices200ResponseAllOfInvoicesInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**OwnerId** | Pointer to **int64** |  | [optional] 
**Account** | Pointer to [**NullableListInvoices200ResponseAllOfInvoicesInnerAccount**](ListInvoices200ResponseAllOfInvoicesInnerAccount.md) |  | [optional] 
**Group** | Pointer to **map[string]interface{}** |  | [optional] 
**Cloud** | Pointer to [**ListInvoices200ResponseAllOfInvoicesInnerCloud**](ListInvoices200ResponseAllOfInvoicesInnerCloud.md) |  | [optional] 
**Instance** | Pointer to **map[string]interface{}** |  | [optional] 
**Server** | Pointer to **NullableString** |  | [optional] 
**Cluster** | Pointer to **NullableString** |  | [optional] 
**User** | Pointer to **map[string]interface{}** |  | [optional] 
**Plan** | Pointer to **map[string]interface{}** |  | [optional] 
**Tags** | Pointer to **[]map[string]interface{}** |  | [optional] 
**Project** | Pointer to **NullableString** |  | [optional] 
**RefType** | Pointer to **string** |  | [optional] 
**RefId** | Pointer to **int64** |  | [optional] 
**RefUuid** | Pointer to **NullableString** |  | [optional] 
**RefName** | Pointer to **string** |  | [optional] 
**RefCategory** | Pointer to **string** |  | [optional] 
**ResourceId** | Pointer to **NullableString** |  | [optional] 
**ResourceUuid** | Pointer to **NullableString** |  | [optional] 
**ResourceType** | Pointer to **NullableString** |  | [optional] 
**ResourceName** | Pointer to **NullableString** |  | [optional] 
**ResourceExternalId** | Pointer to **NullableString** |  | [optional] 
**ResourceInternalId** | Pointer to **NullableString** |  | [optional] 
**Interval** | Pointer to **string** |  | [optional] 
**Period** | Pointer to **string** |  | [optional] 
**Estimate** | Pointer to **bool** |  | [optional] 
**SummaryInvoice** | Pointer to **bool** |  | [optional] 
**Active** | Pointer to **bool** |  | [optional] 
**StartDate** | Pointer to **time.Time** |  | [optional] 
**EndDate** | Pointer to **time.Time** |  | [optional] 
**RefStart** | Pointer to **time.Time** |  | [optional] 
**RefEnd** | Pointer to **time.Time** |  | [optional] 
**EstimatedComputePrice** | Pointer to **float32** |  | [optional] 
**EstimatedComputeCost** | Pointer to **float32** |  | [optional] 
**EstimatedMemoryPrice** | Pointer to **float32** |  | [optional] 
**EstimatedMemoryCost** | Pointer to **float32** |  | [optional] 
**EstimatedStoragePrice** | Pointer to **float32** |  | [optional] 
**EstimatedStorageCost** | Pointer to **float32** |  | [optional] 
**EstimatedNetworkPrice** | Pointer to **float32** |  | [optional] 
**EstimatedNetworkCost** | Pointer to **float32** |  | [optional] 
**EstimatedLicensePrice** | Pointer to **float32** |  | [optional] 
**EstimatedLicenseCost** | Pointer to **float32** |  | [optional] 
**EstimatedExtraPrice** | Pointer to **float32** |  | [optional] 
**EstimatedExtraCost** | Pointer to **float32** |  | [optional] 
**EstimatedTotalPrice** | Pointer to **float32** |  | [optional] 
**EstimatedTotalCost** | Pointer to **float32** |  | [optional] 
**EstimatedRunningPrice** | Pointer to **float32** |  | [optional] 
**EstimatedRunningCost** | Pointer to **float32** |  | [optional] 
**EstimatedCurrency** | Pointer to **string** |  | [optional] 
**EstimatedConversionRate** | Pointer to **float32** |  | [optional] 
**ActualComputePrice** | Pointer to **float32** |  | [optional] 
**ActualComputeCost** | Pointer to **float32** |  | [optional] 
**ActualMemoryPrice** | Pointer to **float32** |  | [optional] 
**ActualMemoryCost** | Pointer to **float32** |  | [optional] 
**ActualStoragePrice** | Pointer to **float32** |  | [optional] 
**ActualStorageCost** | Pointer to **float32** |  | [optional] 
**ActualNetworkPrice** | Pointer to **float32** |  | [optional] 
**ActualNetworkCost** | Pointer to **float32** |  | [optional] 
**ActualLicensePrice** | Pointer to **float32** |  | [optional] 
**ActualLicenseCost** | Pointer to **float32** |  | [optional] 
**ActualExtraPrice** | Pointer to **float32** |  | [optional] 
**ActualExtraCost** | Pointer to **float32** |  | [optional] 
**ActualTotalPrice** | Pointer to **float32** |  | [optional] 
**ActualTotalCost** | Pointer to **float32** |  | [optional] 
**ActualRunningPrice** | Pointer to **float32** |  | [optional] 
**ActualRunningCost** | Pointer to **float32** |  | [optional] 
**ActualCurrency** | Pointer to **string** |  | [optional] 
**ActualConversionRate** | Pointer to **float32** |  | [optional] 
**ComputePrice** | Pointer to **float32** |  | [optional] 
**ComputeCost** | Pointer to **float32** |  | [optional] 
**MemoryPrice** | Pointer to **float32** |  | [optional] 
**MemoryCost** | Pointer to **float32** |  | [optional] 
**StoragePrice** | Pointer to **float32** |  | [optional] 
**StorageCost** | Pointer to **float32** |  | [optional] 
**NetworkPrice** | Pointer to **float32** |  | [optional] 
**NetworkCost** | Pointer to **float32** |  | [optional] 
**LicensePrice** | Pointer to **float32** |  | [optional] 
**LicenseCost** | Pointer to **float32** |  | [optional] 
**ExtraPrice** | Pointer to **float32** |  | [optional] 
**ExtraCost** | Pointer to **float32** |  | [optional] 
**TotalPrice** | Pointer to **float32** |  | [optional] 
**TotalCost** | Pointer to **float32** |  | [optional] 
**RunningPrice** | Pointer to **float32** |  | [optional] 
**RunningCost** | Pointer to **float32** |  | [optional] 
**Currency** | Pointer to **string** |  | [optional] 
**ConversionRate** | Pointer to **float32** |  | [optional] 
**CostType** | Pointer to **string** |  | [optional] 
**OffTime** | Pointer to **int64** |  | [optional] 
**PowerState** | Pointer to **NullableString** |  | [optional] 
**PowerDate** | Pointer to **time.Time** |  | [optional] 
**RunningMultiplier** | Pointer to **float32** |  | [optional] 
**UsageType** | Pointer to **NullableString** |  | [optional] 
**UsageCategory** | Pointer to **NullableString** |  | [optional] 
**LastCostDate** | Pointer to **NullableTime** |  | [optional] 
**LastActualDate** | Pointer to **NullableTime** |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 
**LineItemCount** | Pointer to **int64** |  | [optional] 
**LineItems** | Pointer to [**[]ListInvoices200ResponseAllOfInvoicesInnerLineItemsInner**](ListInvoices200ResponseAllOfInvoicesInnerLineItemsInner.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListInvoices200ResponseAllOfInvoicesInner{
    // Set fields directly
}
```

### Account (Nullable)

Use the Nullable wrapper methods:
- `obj.Account.IsSet()` — check if set
- `obj.Account.Get()` — get the inner value (returns pointer)
- `obj.Account.Set(&val)` — set the value
- `obj.Account.Unset()` — clear the value
### Group (Nullable)

Use the Nullable wrapper methods:
- `obj.Group.IsSet()` — check if set
- `obj.Group.Get()` — get the inner value (returns pointer)
- `obj.Group.Set(&val)` — set the value
- `obj.Group.Unset()` — clear the value
### Instance (Nullable)

Use the Nullable wrapper methods:
- `obj.Instance.IsSet()` — check if set
- `obj.Instance.Get()` — get the inner value (returns pointer)
- `obj.Instance.Set(&val)` — set the value
- `obj.Instance.Unset()` — clear the value
### Server (Nullable)

Use the Nullable wrapper methods:
- `obj.Server.IsSet()` — check if set
- `obj.Server.Get()` — get the inner value (returns pointer)
- `obj.Server.Set(&val)` — set the value
- `obj.Server.Unset()` — clear the value
### Cluster (Nullable)

Use the Nullable wrapper methods:
- `obj.Cluster.IsSet()` — check if set
- `obj.Cluster.Get()` — get the inner value (returns pointer)
- `obj.Cluster.Set(&val)` — set the value
- `obj.Cluster.Unset()` — clear the value
### User (Nullable)

Use the Nullable wrapper methods:
- `obj.User.IsSet()` — check if set
- `obj.User.Get()` — get the inner value (returns pointer)
- `obj.User.Set(&val)` — set the value
- `obj.User.Unset()` — clear the value
### Plan (Nullable)

Use the Nullable wrapper methods:
- `obj.Plan.IsSet()` — check if set
- `obj.Plan.Get()` — get the inner value (returns pointer)
- `obj.Plan.Set(&val)` — set the value
- `obj.Plan.Unset()` — clear the value
### Tags (Nullable)

Use the Nullable wrapper methods:
- `obj.Tags.IsSet()` — check if set
- `obj.Tags.Get()` — get the inner value (returns pointer)
- `obj.Tags.Set(&val)` — set the value
- `obj.Tags.Unset()` — clear the value
### Project (Nullable)

Use the Nullable wrapper methods:
- `obj.Project.IsSet()` — check if set
- `obj.Project.Get()` — get the inner value (returns pointer)
- `obj.Project.Set(&val)` — set the value
- `obj.Project.Unset()` — clear the value
### RefUuid (Nullable)

Use the Nullable wrapper methods:
- `obj.RefUuid.IsSet()` — check if set
- `obj.RefUuid.Get()` — get the inner value (returns pointer)
- `obj.RefUuid.Set(&val)` — set the value
- `obj.RefUuid.Unset()` — clear the value
### ResourceId (Nullable)

Use the Nullable wrapper methods:
- `obj.ResourceId.IsSet()` — check if set
- `obj.ResourceId.Get()` — get the inner value (returns pointer)
- `obj.ResourceId.Set(&val)` — set the value
- `obj.ResourceId.Unset()` — clear the value
### ResourceUuid (Nullable)

Use the Nullable wrapper methods:
- `obj.ResourceUuid.IsSet()` — check if set
- `obj.ResourceUuid.Get()` — get the inner value (returns pointer)
- `obj.ResourceUuid.Set(&val)` — set the value
- `obj.ResourceUuid.Unset()` — clear the value
### ResourceType (Nullable)

Use the Nullable wrapper methods:
- `obj.ResourceType.IsSet()` — check if set
- `obj.ResourceType.Get()` — get the inner value (returns pointer)
- `obj.ResourceType.Set(&val)` — set the value
- `obj.ResourceType.Unset()` — clear the value
### ResourceName (Nullable)

Use the Nullable wrapper methods:
- `obj.ResourceName.IsSet()` — check if set
- `obj.ResourceName.Get()` — get the inner value (returns pointer)
- `obj.ResourceName.Set(&val)` — set the value
- `obj.ResourceName.Unset()` — clear the value
### ResourceExternalId (Nullable)

Use the Nullable wrapper methods:
- `obj.ResourceExternalId.IsSet()` — check if set
- `obj.ResourceExternalId.Get()` — get the inner value (returns pointer)
- `obj.ResourceExternalId.Set(&val)` — set the value
- `obj.ResourceExternalId.Unset()` — clear the value
### ResourceInternalId (Nullable)

Use the Nullable wrapper methods:
- `obj.ResourceInternalId.IsSet()` — check if set
- `obj.ResourceInternalId.Get()` — get the inner value (returns pointer)
- `obj.ResourceInternalId.Set(&val)` — set the value
- `obj.ResourceInternalId.Unset()` — clear the value
### PowerState (Nullable)

Use the Nullable wrapper methods:
- `obj.PowerState.IsSet()` — check if set
- `obj.PowerState.Get()` — get the inner value (returns pointer)
- `obj.PowerState.Set(&val)` — set the value
- `obj.PowerState.Unset()` — clear the value
### UsageType (Nullable)

Use the Nullable wrapper methods:
- `obj.UsageType.IsSet()` — check if set
- `obj.UsageType.Get()` — get the inner value (returns pointer)
- `obj.UsageType.Set(&val)` — set the value
- `obj.UsageType.Unset()` — clear the value
### UsageCategory (Nullable)

Use the Nullable wrapper methods:
- `obj.UsageCategory.IsSet()` — check if set
- `obj.UsageCategory.Get()` — get the inner value (returns pointer)
- `obj.UsageCategory.Set(&val)` — set the value
- `obj.UsageCategory.Unset()` — clear the value
### LastCostDate (Nullable)

Use the Nullable wrapper methods:
- `obj.LastCostDate.IsSet()` — check if set
- `obj.LastCostDate.Get()` — get the inner value (returns pointer)
- `obj.LastCostDate.Set(&val)` — set the value
- `obj.LastCostDate.Unset()` — clear the value
### LastActualDate (Nullable)

Use the Nullable wrapper methods:
- `obj.LastActualDate.IsSet()` — check if set
- `obj.LastActualDate.Get()` — get the inner value (returns pointer)
- `obj.LastActualDate.Set(&val)` — set the value
- `obj.LastActualDate.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


