# GetGuidanceSettings200ResponseGuidanceSettings

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CpuAvgCutoffPower** | Pointer to **NullableInt32** | Power Shutdown Average CPU (%). Lower limit for average CPU usage | [optional] 
**CpuMaxCutoffPower** | Pointer to **NullableInt32** | Power Shutdown Maximum CPU (%). Lower limit for peak CPU usage | [optional] 
**NetworkCutoffPower** | Pointer to **NullableInt32** | Power Shutdown Network threshold (bytes). Lower limit for average network bandwidth | [optional] 
**CpuUpAvgStandardCutoffRightSize** | Pointer to **NullableInt32** | CPU Up-size Average CPU (%). Upper limit for CPU usage | [optional] 
**CpuUpMaxStandardCutoffRightSize** | Pointer to **NullableInt32** | CPU Up-size Maximum CPU (%). Upper limit for peak CPU usage | [optional] 
**MemoryUpAvgStandardCutoffRightSize** | Pointer to **NullableInt32** | Memory Up-size Minimum Free Memory (%). Lower limit for average free memory usage | [optional] 
**MemoryDownAvgStandardCutoffRightSize** | Pointer to **NullableInt32** | Memory Down-size Maximum Free Memory (%). Upper limit for average free memory | [optional] 
**MemoryDownMaxStandardCutoffRightSize** | Pointer to **NullableInt32** | Memory Down-size Maximum Free Memory (%). Upper limit for peak memory usage | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetGuidanceSettings200ResponseGuidanceSettings{
    // Set fields directly
}
```

### CpuAvgCutoffPower (Nullable)

Use the Nullable wrapper methods:
- `obj.CpuAvgCutoffPower.IsSet()` — check if set
- `obj.CpuAvgCutoffPower.Get()` — get the inner value (returns pointer)
- `obj.CpuAvgCutoffPower.Set(&val)` — set the value
- `obj.CpuAvgCutoffPower.Unset()` — clear the value
### CpuMaxCutoffPower (Nullable)

Use the Nullable wrapper methods:
- `obj.CpuMaxCutoffPower.IsSet()` — check if set
- `obj.CpuMaxCutoffPower.Get()` — get the inner value (returns pointer)
- `obj.CpuMaxCutoffPower.Set(&val)` — set the value
- `obj.CpuMaxCutoffPower.Unset()` — clear the value
### NetworkCutoffPower (Nullable)

Use the Nullable wrapper methods:
- `obj.NetworkCutoffPower.IsSet()` — check if set
- `obj.NetworkCutoffPower.Get()` — get the inner value (returns pointer)
- `obj.NetworkCutoffPower.Set(&val)` — set the value
- `obj.NetworkCutoffPower.Unset()` — clear the value
### CpuUpAvgStandardCutoffRightSize (Nullable)

Use the Nullable wrapper methods:
- `obj.CpuUpAvgStandardCutoffRightSize.IsSet()` — check if set
- `obj.CpuUpAvgStandardCutoffRightSize.Get()` — get the inner value (returns pointer)
- `obj.CpuUpAvgStandardCutoffRightSize.Set(&val)` — set the value
- `obj.CpuUpAvgStandardCutoffRightSize.Unset()` — clear the value
### CpuUpMaxStandardCutoffRightSize (Nullable)

Use the Nullable wrapper methods:
- `obj.CpuUpMaxStandardCutoffRightSize.IsSet()` — check if set
- `obj.CpuUpMaxStandardCutoffRightSize.Get()` — get the inner value (returns pointer)
- `obj.CpuUpMaxStandardCutoffRightSize.Set(&val)` — set the value
- `obj.CpuUpMaxStandardCutoffRightSize.Unset()` — clear the value
### MemoryUpAvgStandardCutoffRightSize (Nullable)

Use the Nullable wrapper methods:
- `obj.MemoryUpAvgStandardCutoffRightSize.IsSet()` — check if set
- `obj.MemoryUpAvgStandardCutoffRightSize.Get()` — get the inner value (returns pointer)
- `obj.MemoryUpAvgStandardCutoffRightSize.Set(&val)` — set the value
- `obj.MemoryUpAvgStandardCutoffRightSize.Unset()` — clear the value
### MemoryDownAvgStandardCutoffRightSize (Nullable)

Use the Nullable wrapper methods:
- `obj.MemoryDownAvgStandardCutoffRightSize.IsSet()` — check if set
- `obj.MemoryDownAvgStandardCutoffRightSize.Get()` — get the inner value (returns pointer)
- `obj.MemoryDownAvgStandardCutoffRightSize.Set(&val)` — set the value
- `obj.MemoryDownAvgStandardCutoffRightSize.Unset()` — clear the value
### MemoryDownMaxStandardCutoffRightSize (Nullable)

Use the Nullable wrapper methods:
- `obj.MemoryDownMaxStandardCutoffRightSize.IsSet()` — check if set
- `obj.MemoryDownMaxStandardCutoffRightSize.Get()` — get the inner value (returns pointer)
- `obj.MemoryDownMaxStandardCutoffRightSize.Set(&val)` — set the value
- `obj.MemoryDownMaxStandardCutoffRightSize.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


