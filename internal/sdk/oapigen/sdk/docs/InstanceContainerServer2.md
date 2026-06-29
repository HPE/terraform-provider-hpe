# InstanceContainerServer2

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Uuid** | Pointer to **string** |  | [optional] 
**ExternalId** | Pointer to **NullableString** |  | [optional] 
**InternalId** | Pointer to **NullableString** |  | [optional] 
**HostName** | Pointer to **string** |  | [optional] 
**AccountId** | Pointer to **int64** |  | [optional] 
**Account** | Pointer to [**InstanceContainerServer2Account**](InstanceContainerServer2Account.md) |  | [optional] 
**Owner** | Pointer to [**InstanceContainerServer2Owner**](InstanceContainerServer2Owner.md) |  | [optional] 
**Zone** | Pointer to [**InstanceContainerServer2Zone**](InstanceContainerServer2Zone.md) |  | [optional] 
**Plan** | Pointer to [**InstanceContainerServer2Plan**](InstanceContainerServer2Plan.md) |  | [optional] 
**ComputeServerType** | Pointer to [**InstanceContainerServer2ComputeServerType**](InstanceContainerServer2ComputeServerType.md) |  | [optional] 
**Visibility** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**ZoneId** | Pointer to **int64** |  | [optional] 
**SiteId** | Pointer to **int64** |  | [optional] 
**ResourcePoolId** | Pointer to **int64** |  | [optional] 
**SshHost** | Pointer to **NullableString** |  | [optional] 
**SshPort** | Pointer to **int64** |  | [optional] 
**ExternalIp** | Pointer to **NullableString** |  | [optional] 
**InternalIp** | Pointer to **NullableString** |  | [optional] 
**Platform** | Pointer to **string** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**StatusMessage** | Pointer to **NullableString** |  | [optional] 
**ErrorMessage** | Pointer to **NullableString** |  | [optional] 
**StatusDate** | Pointer to **time.Time** |  | [optional] 
**StatusPercent** | Pointer to **NullableInt64** |  | [optional] 
**StatusEta** | Pointer to **NullableInt64** |  | [optional] 
**PowerState** | Pointer to [**InstanceContainerServer2PowerState**](InstanceContainerServer2PowerState.md) |  | [optional] 
**AgentInstalled** | Pointer to **bool** |  | [optional] 
**LastAgentUpdate** | Pointer to **NullableTime** |  | [optional] 
**MaxCores** | Pointer to **int64** |  | [optional] 
**MaxMemory** | Pointer to **int64** |  | [optional] 
**MaxStorage** | Pointer to **int64** |  | [optional] 
**SourceImage** | Pointer to [**InstanceContainerServer2SourceImage**](InstanceContainerServer2SourceImage.md) |  | [optional] 
**ServerOs** | Pointer to [**InstanceContainerServer2ServerOs**](InstanceContainerServer2ServerOs.md) |  | [optional] 
**Volumes** | Pointer to [**[]InstanceContainerServerVolume1**](InstanceContainerServerVolume1.md) |  | [optional] 
**Interfaces** | Pointer to [**[]InstanceContainerServerInterfacesInner1**](InstanceContainerServerInterfacesInner1.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &InstanceContainerServer2{
    // Set fields directly
}
```

### ExternalId (Nullable)

Use the Nullable wrapper methods:
- `obj.ExternalId.IsSet()` — check if set
- `obj.ExternalId.Get()` — get the inner value (returns pointer)
- `obj.ExternalId.Set(&val)` — set the value
- `obj.ExternalId.Unset()` — clear the value
### InternalId (Nullable)

Use the Nullable wrapper methods:
- `obj.InternalId.IsSet()` — check if set
- `obj.InternalId.Get()` — get the inner value (returns pointer)
- `obj.InternalId.Set(&val)` — set the value
- `obj.InternalId.Unset()` — clear the value
### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### SshHost (Nullable)

Use the Nullable wrapper methods:
- `obj.SshHost.IsSet()` — check if set
- `obj.SshHost.Get()` — get the inner value (returns pointer)
- `obj.SshHost.Set(&val)` — set the value
- `obj.SshHost.Unset()` — clear the value
### ExternalIp (Nullable)

Use the Nullable wrapper methods:
- `obj.ExternalIp.IsSet()` — check if set
- `obj.ExternalIp.Get()` — get the inner value (returns pointer)
- `obj.ExternalIp.Set(&val)` — set the value
- `obj.ExternalIp.Unset()` — clear the value
### InternalIp (Nullable)

Use the Nullable wrapper methods:
- `obj.InternalIp.IsSet()` — check if set
- `obj.InternalIp.Get()` — get the inner value (returns pointer)
- `obj.InternalIp.Set(&val)` — set the value
- `obj.InternalIp.Unset()` — clear the value
### StatusMessage (Nullable)

Use the Nullable wrapper methods:
- `obj.StatusMessage.IsSet()` — check if set
- `obj.StatusMessage.Get()` — get the inner value (returns pointer)
- `obj.StatusMessage.Set(&val)` — set the value
- `obj.StatusMessage.Unset()` — clear the value
### ErrorMessage (Nullable)

Use the Nullable wrapper methods:
- `obj.ErrorMessage.IsSet()` — check if set
- `obj.ErrorMessage.Get()` — get the inner value (returns pointer)
- `obj.ErrorMessage.Set(&val)` — set the value
- `obj.ErrorMessage.Unset()` — clear the value
### StatusPercent (Nullable)

Use the Nullable wrapper methods:
- `obj.StatusPercent.IsSet()` — check if set
- `obj.StatusPercent.Get()` — get the inner value (returns pointer)
- `obj.StatusPercent.Set(&val)` — set the value
- `obj.StatusPercent.Unset()` — clear the value
### StatusEta (Nullable)

Use the Nullable wrapper methods:
- `obj.StatusEta.IsSet()` — check if set
- `obj.StatusEta.Get()` — get the inner value (returns pointer)
- `obj.StatusEta.Set(&val)` — set the value
- `obj.StatusEta.Unset()` — clear the value
### LastAgentUpdate (Nullable)

Use the Nullable wrapper methods:
- `obj.LastAgentUpdate.IsSet()` — check if set
- `obj.LastAgentUpdate.Get()` — get the inner value (returns pointer)
- `obj.LastAgentUpdate.Set(&val)` — set the value
- `obj.LastAgentUpdate.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


