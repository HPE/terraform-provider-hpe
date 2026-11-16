// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package hostprofile

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

// snakeToCamel converts a snake_case key to camelCase to match the SDK's JSON
// field tags (e.g. find_multipaths -> findMultipaths).
func snakeToCamel(s string) string {
	parts := strings.Split(s, "_")
	for i := 1; i < len(parts); i++ {
		if parts[i] == "" {
			continue
		}
		parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
	}

	return strings.Join(parts, "")
}

// renameKeysToCamel recursively rewrites the keys of an arbitrary decoded
// value from snake_case (Terraform) to camelCase (SDK JSON).
func renameKeysToCamel(v any) any {
	switch t := v.(type) {
	case map[string]any:
		out := make(map[string]any, len(t))
		for k, val := range t {
			out[snakeToCamel(k)] = renameKeysToCamel(val)
		}

		return out
	case []any:
		for i := range t {
			t[i] = renameKeysToCamel(t[i])
		}

		return t
	default:
		return v
	}
}

// expandConfigMap converts the planned config object into a camelCase map
// ready to be marshalled into any of the SDK request config structs.
func expandConfigMap(
	ctx context.Context,
	cfg ConfigValue,
) (map[string]any, diag.Diagnostics) {
	var diags diag.Diagnostics

	if cfg.IsNull() || cfg.IsUnknown() {
		return nil, diags
	}

	obj, objDiags := cfg.ToObjectValue(ctx)
	diags.Append(objDiags...)
	if diags.HasError() {
		return nil, diags
	}

	raw, err := convert.ObjectToAny(ctx, obj)
	if err != nil {
		diags.AddError(
			"expand host profile config",
			fmt.Sprintf("failed to convert config object: %s", err),
		)

		return nil, diags
	}

	camel, ok := renameKeysToCamel(raw).(map[string]any)
	if !ok {
		diags.AddError(
			"expand host profile config",
			"config did not convert to an object",
		)

		return nil, diags
	}

	return camel, diags
}

// expandCreateConfig builds the SDK create request config from plan.
func expandCreateConfig(
	ctx context.Context,
	cfg ConfigValue,
) (*sdk.CreateHostProfileRequestConfig, diag.Diagnostics) {
	m, diags := expandConfigMap(ctx, cfg)
	if diags.HasError() || m == nil {
		return nil, diags
	}

	encoded, err := json.Marshal(m)
	if err != nil {
		diags.AddError("expand host profile config", err.Error())

		return nil, diags
	}

	var out sdk.CreateHostProfileRequestConfig
	if err := json.Unmarshal(encoded, &out); err != nil {
		diags.AddError("expand host profile config", err.Error())

		return nil, diags
	}

	return &out, diags
}

// expandUpdateConfig builds the SDK update request config from plan.
func expandUpdateConfig(
	ctx context.Context,
	cfg ConfigValue,
) (*sdk.UpdateHostProfileRequestConfig, diag.Diagnostics) {
	m, diags := expandConfigMap(ctx, cfg)
	if diags.HasError() || m == nil {
		return nil, diags
	}

	encoded, err := json.Marshal(m)
	if err != nil {
		diags.AddError("expand host profile config", err.Error())

		return nil, diags
	}

	var out sdk.UpdateHostProfileRequestConfig
	if err := json.Unmarshal(encoded, &out); err != nil {
		diags.AddError("expand host profile config", err.Error())

		return nil, diags
	}

	return &out, diags
}

// flattenConfig builds the Terraform ConfigValue from the SDK response config.
// When plan.Config is null the user did not manage config, so we keep it null
// to avoid provider-produced inconsistency. The one exception is import, which
// has no prior plan and goes through flattenConfigAlways directly.
func flattenConfig(
	ctx context.Context,
	apiCfg *sdk.GetHostProfile200ResponseHostProfileAllOfConfig,
	plan ConfigValue,
) (ConfigValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if plan.IsNull() || plan.IsUnknown() || apiCfg == nil {
		return NewConfigValueNull(), diags
	}

	return flattenConfigAlways(ctx, apiCfg)
}

// flattenConfigAlways reconstructs the full ConfigValue from the API response
// unconditionally, without the plan-null guard flattenConfig applies. It backs
// import, where there is no prior plan but the stored config should still be
// populated into state.
func flattenConfigAlways(
	ctx context.Context,
	apiCfg *sdk.GetHostProfile200ResponseHostProfileAllOfConfig,
) (ConfigValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if apiCfg == nil {
		return NewConfigValueNull(), diags
	}

	baseline, d := flattenBaseline(ctx, apiCfg.Baseline)
	diags.Append(d...)
	storage, d := flattenStorage(ctx, apiCfg.Storage)
	diags.Append(d...)
	system, d := flattenSystem(ctx, apiCfg.System)
	diags.Append(d...)
	if diags.HasError() {
		return NewConfigValueNull(), diags
	}

	return ConfigValue{
		Baseline: baseline,
		Storage:  storage,
		System:   system,
		state:    attr.ValueStateKnown,
	}, diags
}

func objVal(
	attrTypes map[string]attr.Type,
	attrs map[string]attr.Value,
	diags *diag.Diagnostics,
) types.Object {
	o, d := types.ObjectValue(attrTypes, attrs)
	diags.Append(d...)

	return o
}

func flattenBaseline(
	ctx context.Context,
	b *sdk.GetHostProfile200ResponseHostProfileAllOfConfigBaseline,
) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	at := BaselineValue{}.AttributeTypes(ctx)
	if b == nil {
		return types.ObjectNull(at), diags
	}

	pvTypes := PackageVersionsValue{}.AttributeTypes(ctx)
	var pv types.Object
	if b.PackageVersions == nil {
		pv = types.ObjectNull(pvTypes)
	} else {
		pv = objVal(pvTypes, map[string]attr.Value{
			"corosync": convert.StrToType(b.PackageVersions.Corosync),
			"libvirt":  convert.StrToType(b.PackageVersions.Libvirt),
			"qemu":     convert.StrToType(b.PackageVersions.Qemu),
		}, &diags)
	}

	return objVal(at, map[string]attr.Value{
		"additional_packages": convert.StrToType(b.AdditionalPackages),
		"kernel_version":      convert.StrToType(b.KernelVersion),
		"package_versions":    pv,
	}, &diags), diags
}

func flattenStorage(
	ctx context.Context,
	s *sdk.GetHostProfile200ResponseHostProfileAllOfConfigStorage,
) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	at := StorageValue{}.AttributeTypes(ctx)
	if s == nil {
		return types.ObjectNull(at), diags
	}

	iscsiTypes := IscsiValue{}.AttributeTypes(ctx)
	var iscsi types.Object
	if s.Iscsi == nil {
		iscsi = types.ObjectNull(iscsiTypes)
	} else {
		iscsi = objVal(iscsiTypes, map[string]attr.Value{
			"chap_enabled": convert.BoolToType(s.Iscsi.ChapEnabled),
			"transport":    convert.StrToType(s.Iscsi.Transport),
		}, &diags)
	}

	mpTypes := MultipathValue{}.AttributeTypes(ctx)
	var mp types.Object
	if s.Multipath == nil {
		mp = types.ObjectNull(mpTypes)
	} else {
		mp = objVal(mpTypes, map[string]attr.Value{
			"enabled":             convert.BoolToType(s.Multipath.Enabled),
			"find_multipaths":     convert.BoolToType(s.Multipath.FindMultipaths),
			"user_friendly_names": convert.BoolToType(s.Multipath.UserFriendlyNames),
		}, &diags)
	}

	return objVal(at, map[string]attr.Value{
		"iscsi":     iscsi,
		"multipath": mp,
	}, &diags), diags
}

func flattenSystem(
	ctx context.Context,
	s *sdk.GetHostProfile200ResponseHostProfileAllOfConfigSystem,
) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	at := SystemValue{}.AttributeTypes(ctx)
	if s == nil {
		return types.ObjectNull(at), diags
	}

	cpuTypes := CpuValue{}.AttributeTypes(ctx)
	var cpu types.Object
	if s.Cpu == nil {
		cpu = types.ObjectNull(cpuTypes)
	} else {
		cpu = objVal(cpuTypes, map[string]attr.Value{
			"hyper_threading":  convert.BoolToType(s.Cpu.HyperThreading),
			"power_management": convert.StrToType(s.Cpu.PowerManagement),
		}, &diags)
	}

	loggingTypes := LoggingValue{}.AttributeTypes(ctx)
	var logging types.Object
	if s.Logging == nil {
		logging = types.ObjectNull(loggingTypes)
	} else {
		logging = objVal(loggingTypes, map[string]attr.Value{
			"level": convert.StrToType(s.Logging.Level),
		}, &diags)
	}

	memTypes := MemoryValue{}.AttributeTypes(ctx)
	var memory types.Object
	if s.Memory == nil {
		memory = types.ObjectNull(memTypes)
	} else {
		hpTypes := HugePagesValue{}.AttributeTypes(ctx)
		var hp types.Object
		if s.Memory.HugePages == nil {
			hp = types.ObjectNull(hpTypes)
		} else {
			hp = objVal(hpTypes, map[string]attr.Value{
				"enabled": convert.BoolToType(s.Memory.HugePages.Enabled),
				"size":    convert.StrToType(s.Memory.HugePages.Size),
			}, &diags)
		}
		memory = objVal(memTypes, map[string]attr.Value{
			"huge_pages": hp,
		}, &diags)
	}

	ntpTypes := NtpValue{}.AttributeTypes(ctx)
	var ntp types.Object
	if s.Ntp == nil {
		ntp = types.ObjectNull(ntpTypes)
	} else {
		ntp = objVal(ntpTypes, map[string]attr.Value{
			"enabled": convert.BoolToType(s.Ntp.Enabled),
			"servers": serversToSet(s.Ntp.Servers, &diags),
		}, &diags)
	}

	secTypes := SecurityValue{}.AttributeTypes(ctx)
	var security types.Object
	if s.Security == nil {
		security = types.ObjectNull(secTypes)
	} else {
		security = objVal(secTypes, map[string]attr.Value{
			"fips":    convert.BoolToType(s.Security.Fips),
			"selinux": convert.StrToType(s.Security.Selinux),
		}, &diags)
	}

	syslogTypes := SyslogValue{}.AttributeTypes(ctx)
	var syslog types.Object
	if s.Syslog == nil {
		syslog = types.ObjectNull(syslogTypes)
	} else {
		syslog = objVal(syslogTypes, map[string]attr.Value{
			"enabled":  convert.BoolToType(s.Syslog.Enabled),
			"port":     convert.Int32ToType(s.Syslog.Port),
			"protocol": convert.StrToType(s.Syslog.Protocol),
			"server":   convert.StrToType(s.Syslog.Server),
		}, &diags)
	}

	return objVal(at, map[string]attr.Value{
		"cpu":      cpu,
		"locale":   convert.StrToType(s.Locale),
		"logging":  logging,
		"memory":   memory,
		"ntp":      ntp,
		"security": security,
		"syslog":   syslog,
		"timezone": convert.StrToType(s.Timezone),
	}, &diags), diags
}

// serversToSet maps the API's ntp.servers list onto the plain-Optional
// `servers` set, keeping nil and empty apart: an absent key (nil slice) is
// null, which is what an omitted attribute plans as, while a present empty list
// is a known empty set, which is what `servers = []` plans as. The shared
// convert.StrSliceToSet folds both into null, which made an accepted
// `servers = []` read back as null and fail the apply as an inconsistent
// result. The SDK decoder preserves the distinction ("servers":[] decodes to
// an empty, non-nil slice), so it only has to be carried through here.
func serversToSet(servers []string, diags *diag.Diagnostics) types.Set {
	if servers == nil {
		return types.SetNull(types.StringType)
	}

	elems := make([]attr.Value, 0, len(servers))
	for _, s := range servers {
		elems = append(elems, types.StringValue(s))
	}

	set, d := types.SetValue(types.StringType, elems)
	diags.Append(d...)

	return set
}

// flattenVersions builds the computed versions set from the SDK response.
func flattenVersions(
	ctx context.Context,
	versions []sdk.GetHostProfile200ResponseHostProfileAllOfVersionsInner,
) (types.Set, diag.Diagnostics) {
	var diags diag.Diagnostics
	elemType := basetypes.ObjectType{AttrTypes: VersionsValue{}.AttributeTypes(ctx)}

	// Treat an absent and an empty versions list identically (null), matching
	// the data source's versionsAsSet so the same API shape maps the same way.
	if len(versions) == 0 {
		return types.SetNull(elemType), diags
	}

	vt := VersionsValue{}.AttributeTypes(ctx)
	elems := make([]attr.Value, 0, len(versions))
	for _, v := range versions {
		o := objVal(vt, map[string]attr.Value{
			"date_created":   convert.TimeToType(v.DateCreated),
			"id":             convert.Int64ToType(v.Id),
			"version_number": convert.Int32ToType(v.VersionNumber),
		}, &diags)
		elems = append(elems, o)
	}

	set, d := types.SetValue(elemType, elems)
	diags.Append(d...)

	return set, diags
}
