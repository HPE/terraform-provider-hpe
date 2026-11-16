// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package hostprofile

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

// hostProfileAsState maps the SDK host profile into the Terraform model.
func hostProfileAsState(
	ctx context.Context,
	hp *sdk.GetHostProfile200ResponseHostProfile,
) (HostProfileModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	state := HostProfileModel{
		Id:               convert.Int64ToType(hp.Id),
		Code:             convert.StrToType(hp.Code.Get()),
		Name:             convert.StrToType(hp.Name),
		Description:      convert.StrToType(hp.Description.Get()),
		Enabled:          convert.BoolToType(hp.Enabled),
		ServerCount:      convert.Int32ToType(hp.ServerCount),
		ComplianceStatus: convert.StrToType(hp.ComplianceStatus),
	}

	cfg, cfgDiags := configAsValue(ctx, hp.Config)
	diags.Append(cfgDiags...)
	state.Config = cfg

	versions, verDiags := versionsAsSet(ctx, hp.Versions)
	diags.Append(verDiags...)
	state.Versions = versions

	return state, diags
}

func configAsValue(
	ctx context.Context,
	cfg *sdk.GetHostProfile200ResponseHostProfileAllOfConfig,
) (ConfigValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if cfg == nil {
		return NewConfigValueNull(), diags
	}

	baseline, d := baselineAsObject(ctx, cfg.Baseline)
	diags.Append(d...)

	storage, d := storageAsObject(ctx, cfg.Storage)
	diags.Append(d...)

	system, d := systemAsObject(ctx, cfg.System)
	diags.Append(d...)

	return ConfigValue{
		Baseline: baseline,
		Storage:  storage,
		System:   system,
		state:    attr.ValueStateKnown,
	}, diags
}

func baselineAsObject(
	ctx context.Context,
	b *sdk.GetHostProfile200ResponseHostProfileAllOfConfigBaseline,
) (basetypes.ObjectValue, diag.Diagnostics) {
	if b == nil {
		return types.ObjectNull(BaselineValue{}.AttributeTypes(ctx)), nil
	}

	pv := NewPackageVersionsValueNull()
	if b.PackageVersions != nil {
		pv = PackageVersionsValue{
			Corosync: convert.StrToType(b.PackageVersions.Corosync),
			Libvirt:  convert.StrToType(b.PackageVersions.Libvirt),
			Qemu:     convert.StrToType(b.PackageVersions.Qemu),
			state:    attr.ValueStateKnown,
		}
	}

	pvObj, diags := pv.ToObjectValue(ctx)
	if diags.HasError() {
		// Don't fold a half-formed object into the parent: that would compound
		// into a second, misleading conversion error. Surface the real one.
		return types.ObjectNull(BaselineValue{}.AttributeTypes(ctx)), diags
	}

	v := BaselineValue{
		AdditionalPackages: convert.StrToType(b.AdditionalPackages),
		KernelVersion:      convert.StrToType(b.KernelVersion),
		PackageVersions:    pvObj,
		state:              attr.ValueStateKnown,
	}

	obj, d := v.ToObjectValue(ctx)
	diags.Append(d...)

	return obj, diags
}

func storageAsObject(
	ctx context.Context,
	s *sdk.GetHostProfile200ResponseHostProfileAllOfConfigStorage,
) (basetypes.ObjectValue, diag.Diagnostics) {
	if s == nil {
		return types.ObjectNull(StorageValue{}.AttributeTypes(ctx)), nil
	}

	var diags diag.Diagnostics

	iscsi := NewIscsiValueNull()
	if s.Iscsi != nil {
		iscsi = IscsiValue{
			ChapEnabled: convert.BoolToType(s.Iscsi.ChapEnabled),
			Transport:   convert.StrToType(s.Iscsi.Transport),
			state:       attr.ValueStateKnown,
		}
	}

	iscsiObj, d := iscsi.ToObjectValue(ctx)
	diags.Append(d...)

	multipath := NewMultipathValueNull()
	if s.Multipath != nil {
		multipath = MultipathValue{
			Enabled:           convert.BoolToType(s.Multipath.Enabled),
			FindMultipaths:    convert.BoolToType(s.Multipath.FindMultipaths),
			UserFriendlyNames: convert.BoolToType(s.Multipath.UserFriendlyNames),
			state:             attr.ValueStateKnown,
		}
	}

	multipathObj, d := multipath.ToObjectValue(ctx)
	diags.Append(d...)

	v := StorageValue{
		Iscsi:     iscsiObj,
		Multipath: multipathObj,
		state:     attr.ValueStateKnown,
	}

	obj, d := v.ToObjectValue(ctx)
	diags.Append(d...)

	return obj, diags
}

func systemAsObject(
	ctx context.Context,
	s *sdk.GetHostProfile200ResponseHostProfileAllOfConfigSystem,
) (basetypes.ObjectValue, diag.Diagnostics) {
	if s == nil {
		return types.ObjectNull(SystemValue{}.AttributeTypes(ctx)), nil
	}

	var diags diag.Diagnostics

	cpu := NewCpuValueNull()
	if s.Cpu != nil {
		cpu = CpuValue{
			HyperThreading:  convert.BoolToType(s.Cpu.HyperThreading),
			PowerManagement: convert.StrToType(s.Cpu.PowerManagement),
			state:           attr.ValueStateKnown,
		}
	}

	cpuObj, d := cpu.ToObjectValue(ctx)
	diags.Append(d...)

	logging := NewLoggingValueNull()
	if s.Logging != nil {
		logging = LoggingValue{
			Level: convert.StrToType(s.Logging.Level),
			state: attr.ValueStateKnown,
		}
	}

	loggingObj, d := logging.ToObjectValue(ctx)
	diags.Append(d...)

	memory := NewMemoryValueNull()
	if s.Memory != nil {
		hugePages := NewHugePagesValueNull()
		if s.Memory.HugePages != nil {
			hugePages = HugePagesValue{
				Enabled: convert.BoolToType(s.Memory.HugePages.Enabled),
				Size:    convert.StrToType(s.Memory.HugePages.Size),
				state:   attr.ValueStateKnown,
			}
		}

		hugePagesObj, hd := hugePages.ToObjectValue(ctx)
		diags.Append(hd...)

		memory = MemoryValue{
			HugePages: hugePagesObj,
			state:     attr.ValueStateKnown,
		}
	}

	memoryObj, d := memory.ToObjectValue(ctx)
	diags.Append(d...)

	ntp := NewNtpValueNull()
	if s.Ntp != nil {
		ntp = NtpValue{
			Enabled: convert.BoolToType(s.Ntp.Enabled),
			Servers: convert.StrSliceToSet(s.Ntp.Servers),
			state:   attr.ValueStateKnown,
		}
	}

	ntpObj, d := ntp.ToObjectValue(ctx)
	diags.Append(d...)

	security := NewSecurityValueNull()
	if s.Security != nil {
		security = SecurityValue{
			Fips:    convert.BoolToType(s.Security.Fips),
			Selinux: convert.StrToType(s.Security.Selinux),
			state:   attr.ValueStateKnown,
		}
	}

	securityObj, d := security.ToObjectValue(ctx)
	diags.Append(d...)

	syslog := NewSyslogValueNull()
	if s.Syslog != nil {
		syslog = SyslogValue{
			Enabled:  convert.BoolToType(s.Syslog.Enabled),
			Port:     convert.Int32ToType(s.Syslog.Port),
			Protocol: convert.StrToType(s.Syslog.Protocol),
			Server:   convert.StrToType(s.Syslog.Server),
			state:    attr.ValueStateKnown,
		}
	}

	syslogObj, d := syslog.ToObjectValue(ctx)
	diags.Append(d...)

	v := SystemValue{
		Cpu:      cpuObj,
		Locale:   convert.StrToType(s.Locale),
		Logging:  loggingObj,
		Memory:   memoryObj,
		Ntp:      ntpObj,
		Security: securityObj,
		Syslog:   syslogObj,
		Timezone: convert.StrToType(s.Timezone),
		state:    attr.ValueStateKnown,
	}

	obj, d := v.ToObjectValue(ctx)
	diags.Append(d...)

	return obj, diags
}

func versionsAsSet(
	ctx context.Context,
	versions []sdk.GetHostProfile200ResponseHostProfileAllOfVersionsInner,
) (types.Set, diag.Diagnostics) {
	elemType := VersionsValue{}.Type(ctx)

	if len(versions) == 0 {
		return types.SetNull(elemType), nil
	}

	elems := make([]VersionsValue, 0, len(versions))
	for i := range versions {
		v := versions[i]
		elems = append(elems, VersionsValue{
			DateCreated:   convert.TimeToType(v.DateCreated),
			Id:            convert.Int64ToType(v.Id),
			VersionNumber: convert.Int32ToType(v.VersionNumber),
			state:         attr.ValueStateKnown,
		})
	}

	return types.SetValueFrom(ctx, elemType, elems)
}
