# InstallLicense200ResponseLicense

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** | ID | [optional] 
**KeyId** | Pointer to **string** | Key ID (only the first 8 characters are required to identify license to uninstall) | [optional] 
**Hash** | Pointer to **string** | Hash of the license content which can be used as a fingerprint identifier | [optional] 
**LicenseVersion** | Pointer to **int64** | License Version which determines the required appliance version to install this license. | [optional] 
**ProductTier** | Pointer to **string** | Product Tier | [optional] 
**StartDate** | Pointer to **time.Time** | The start date of the applied license. | [optional] 
**EndDate** | Pointer to **time.Time** | The expiration date of the applied license. | [optional] 
**MaxInstances** | Pointer to **int64** | Workload Limit. 0 is used for unlimited. | [optional] 
**MaxMemory** | Pointer to **int64** | Memory Limit. 0 is used for unlimited. | [optional] 
**MaxStorage** | Pointer to **int64** | Storage Limit. 0 is used for unlimited. | [optional] 
**LimitType** | Pointer to **string** | The limit type determines which limits apply to the license, the new &#39;standard&#39; or legacy &#39;workload&#39;. | [optional] 
**MaxManagedServers** | Pointer to **NullableInt64** | Managed Servers Limit. 0 is enforced and null is used for unlimited. | [optional] 
**MaxDiscoveredServers** | Pointer to **NullableInt64** | Discovered Servers Limit. 0 is enforced and null is used for unlimited. | [optional] 
**MaxHosts** | Pointer to **NullableInt64** | Host Limit. 0 is enforced and null is used for unlimited. | [optional] 
**MaxMvm** | Pointer to **NullableInt64** | HPE VM Host Limit. 0 is enforced and null is used for unlimited. | [optional] 
**MaxMvmSockets** | Pointer to **NullableInt64** | HPE VM Host Limit. 0 is enforced and null is used for unlimited. | [optional] 
**MaxSockets** | Pointer to **NullableInt64** | Global Socket Limit. 0 is enforced and null is used for unlimited. | [optional] 
**MaxIac** | Pointer to **NullableInt64** | IAC Deployments Limit. 0 is enforced and null is used for unlimited. | [optional] 
**MaxXaas** | Pointer to **NullableInt64** | Xaas Instances Limit. 0 is enforced and null is used for unlimited. | [optional] 
**MaxExecutions** | Pointer to **NullableInt64** | Execution Limit. 0 is enforced and null is used for unlimited. | [optional] 
**MaxDistributedWorkers** | Pointer to **NullableInt64** | Distributed Workers Limit. 0 is enforced and null is used for unlimited. | [optional] 
**MaxDiscoveredObjects** | Pointer to **NullableInt64** | Discovered Objects Limit. Not yet enforced. | [optional] 
**HardLimit** | Pointer to **bool** | Hard Limit | [optional] 
**FreeTrial** | Pointer to **bool** | Free Trial (Community License) | [optional] 
**MultiTenant** | Pointer to **bool** | Multi-Tenant Enabled | [optional] 
**Whitelabel** | Pointer to **bool** | White Label Enabled | [optional] 
**ReportStatus** | Pointer to **bool** | Stats Reporting. This is true when the appliance registers and sends usage stats to the hub. | [optional] 
**SupportLevel** | Pointer to **string** | Support Level | [optional] 
**AccountName** | Pointer to **string** | Account Name | [optional] 
**Config** | Pointer to **map[string]interface{}** | License Configuration Object | [optional] 
**AmazonProductCodes** | Pointer to **NullableString** |  | [optional] 
**Features** | Pointer to [**InstallLicense200ResponseLicenseFeatures**](InstallLicense200ResponseLicenseFeatures.md) |  | [optional] 
**ZoneTypes** | Pointer to **NullableString** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**RecalculationDate** | Pointer to **NullableTime** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &InstallLicense200ResponseLicense{
    // Set fields directly
}
```

### MaxManagedServers (Nullable)

Use the Nullable wrapper methods:
- `obj.MaxManagedServers.IsSet()` — check if set
- `obj.MaxManagedServers.Get()` — get the inner value (returns pointer)
- `obj.MaxManagedServers.Set(&val)` — set the value
- `obj.MaxManagedServers.Unset()` — clear the value
### MaxDiscoveredServers (Nullable)

Use the Nullable wrapper methods:
- `obj.MaxDiscoveredServers.IsSet()` — check if set
- `obj.MaxDiscoveredServers.Get()` — get the inner value (returns pointer)
- `obj.MaxDiscoveredServers.Set(&val)` — set the value
- `obj.MaxDiscoveredServers.Unset()` — clear the value
### MaxHosts (Nullable)

Use the Nullable wrapper methods:
- `obj.MaxHosts.IsSet()` — check if set
- `obj.MaxHosts.Get()` — get the inner value (returns pointer)
- `obj.MaxHosts.Set(&val)` — set the value
- `obj.MaxHosts.Unset()` — clear the value
### MaxMvm (Nullable)

Use the Nullable wrapper methods:
- `obj.MaxMvm.IsSet()` — check if set
- `obj.MaxMvm.Get()` — get the inner value (returns pointer)
- `obj.MaxMvm.Set(&val)` — set the value
- `obj.MaxMvm.Unset()` — clear the value
### MaxMvmSockets (Nullable)

Use the Nullable wrapper methods:
- `obj.MaxMvmSockets.IsSet()` — check if set
- `obj.MaxMvmSockets.Get()` — get the inner value (returns pointer)
- `obj.MaxMvmSockets.Set(&val)` — set the value
- `obj.MaxMvmSockets.Unset()` — clear the value
### MaxSockets (Nullable)

Use the Nullable wrapper methods:
- `obj.MaxSockets.IsSet()` — check if set
- `obj.MaxSockets.Get()` — get the inner value (returns pointer)
- `obj.MaxSockets.Set(&val)` — set the value
- `obj.MaxSockets.Unset()` — clear the value
### MaxIac (Nullable)

Use the Nullable wrapper methods:
- `obj.MaxIac.IsSet()` — check if set
- `obj.MaxIac.Get()` — get the inner value (returns pointer)
- `obj.MaxIac.Set(&val)` — set the value
- `obj.MaxIac.Unset()` — clear the value
### MaxXaas (Nullable)

Use the Nullable wrapper methods:
- `obj.MaxXaas.IsSet()` — check if set
- `obj.MaxXaas.Get()` — get the inner value (returns pointer)
- `obj.MaxXaas.Set(&val)` — set the value
- `obj.MaxXaas.Unset()` — clear the value
### MaxExecutions (Nullable)

Use the Nullable wrapper methods:
- `obj.MaxExecutions.IsSet()` — check if set
- `obj.MaxExecutions.Get()` — get the inner value (returns pointer)
- `obj.MaxExecutions.Set(&val)` — set the value
- `obj.MaxExecutions.Unset()` — clear the value
### MaxDistributedWorkers (Nullable)

Use the Nullable wrapper methods:
- `obj.MaxDistributedWorkers.IsSet()` — check if set
- `obj.MaxDistributedWorkers.Get()` — get the inner value (returns pointer)
- `obj.MaxDistributedWorkers.Set(&val)` — set the value
- `obj.MaxDistributedWorkers.Unset()` — clear the value
### MaxDiscoveredObjects (Nullable)

Use the Nullable wrapper methods:
- `obj.MaxDiscoveredObjects.IsSet()` — check if set
- `obj.MaxDiscoveredObjects.Get()` — get the inner value (returns pointer)
- `obj.MaxDiscoveredObjects.Set(&val)` — set the value
- `obj.MaxDiscoveredObjects.Unset()` — clear the value
### Config (Nullable)

Use the Nullable wrapper methods:
- `obj.Config.IsSet()` — check if set
- `obj.Config.Get()` — get the inner value (returns pointer)
- `obj.Config.Set(&val)` — set the value
- `obj.Config.Unset()` — clear the value
### AmazonProductCodes (Nullable)

Use the Nullable wrapper methods:
- `obj.AmazonProductCodes.IsSet()` — check if set
- `obj.AmazonProductCodes.Get()` — get the inner value (returns pointer)
- `obj.AmazonProductCodes.Set(&val)` — set the value
- `obj.AmazonProductCodes.Unset()` — clear the value
### ZoneTypes (Nullable)

Use the Nullable wrapper methods:
- `obj.ZoneTypes.IsSet()` — check if set
- `obj.ZoneTypes.Get()` — get the inner value (returns pointer)
- `obj.ZoneTypes.Set(&val)` — set the value
- `obj.ZoneTypes.Unset()` — clear the value
### RecalculationDate (Nullable)

Use the Nullable wrapper methods:
- `obj.RecalculationDate.IsSet()` — check if set
- `obj.RecalculationDate.Get()` — get the inner value (returns pointer)
- `obj.RecalculationDate.Set(&val)` — set the value
- `obj.RecalculationDate.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


