# UpdateVDIPoolsRequestVdiPool

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Virtual Desktop name | [optional] 
**Description** | Pointer to **string** | Virtual Desktop description | [optional] 
**Owner** | Pointer to **int64** | Owner (User) ID | [optional] 
**MinIdle** | Pointer to **float32** | Min Idle - Sets the minimum number of idle instances on standby in the pool. The pool will always try to maintain this number of available instances on standby.  | [optional] 
**InitialPoolSize** | Pointer to **float32** | The initial size of the pool to be allocated on creation | [optional] 
**MaxIdle** | Pointer to **float32** | Sets the maximum number of idle instances on standby in the pool. If the number of idle instances supersedes this, the pool will start removing instances  | [optional] 
**MaxPoolSize** | Pointer to **float32** | Max limit on number of allocations and instances within the pool.  | [optional] 
**AllocationTimeoutMinutes** | Pointer to **float32** | Time (in minutes) after a user disconnects before an allocation is recycled or shutdown depending on persistence.  | [optional] 
**PersistentUser** | Pointer to **bool** | Persistent Desktop Pool | [optional] [default to false]
**Recyclable** | Pointer to **bool** | Recyclable VDI Pools only work with cloud types that support snapshot management (i.e. Vmware, Nutanix, VCD) | [optional] [default to false]
**AllowCopy** | Pointer to **bool** | Allow copy from desktop | [optional] [default to false]
**AllowPrinter** | Pointer to **bool** | Allow local printers from Desktop | [optional] [default to false]
**AllowFileshare** | Pointer to **bool** | Allow File Share | [optional] [default to false]
**AllowHypervisorConsole** | Pointer to **bool** | Allow hypervisor console | [optional] [default to false]
**AutoCreateLocalUserOnReservation** | Pointer to **bool** | Auto create local user upon reservation | [optional] [default to false]
**Enabled** | Pointer to **bool** | Can be used to enable or disable the VDI pool | [optional] [default to true]
**IconPath** | Pointer to **string** | The relative location of an icon image | [optional] 
**Apps** | Pointer to **[]int64** | Array of VDI App IDs | [optional] 
**Gateway** | Pointer to **int64** | VDI Gateway ID | [optional] 
**InstanceConfig** | Pointer to **string** | Instance Config JSON. Passing as a string will preserve property order.  See &#x60;config&#x60; object for required values. | [optional] 
**Config** | Pointer to [**UpdateVDIPoolsRequestVdiPoolConfig**](UpdateVDIPoolsRequestVdiPoolConfig.md) |  | [optional] 
**GuestConsoleJumpHost** | Pointer to **string** | Guest Console Jump Host | [optional] 
**GuestConsoleJumpPort** | Pointer to **int64** | Guest Console Jump Port | [optional] 
**GuestConsoleJumpUsername** | Pointer to **string** | Guest Console Jump Username | [optional] 
**GuestConsoleJumpPassword** | Pointer to **string** | Guest Console Jump Password | [optional] 
**GuestConsoleJumpKeypair** | Pointer to **int64** | Guest Console Jump Key Pair. see &#x60;Key Pair&#x60; | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateVDIPoolsRequestVdiPool{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


