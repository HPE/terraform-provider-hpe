// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancermonitor

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// The HTTP health-check fields of the config block go into the monitor's
// config under the SDN plugin's option type keys; an unset one is not sent.
func TestUnitHTTPHealthCheckConfigIsSentUnderTheConfigKeys(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	cfg, d := NewConfigValue(ConfigValue{}.AttributeTypes(ctx), map[string]attr.Value{
		"monitor":          types.ObjectNull(map[string]attr.Type{"id": types.Int64Type}),
		"monitor_config":   types.DynamicNull(),
		"http_method":      types.StringValue("GET"),
		"http_uri":         types.StringValue("/health"),
		"http_port":        types.Int64Value(8080),
		"http_return_code": types.StringNull(),
	})
	if d.HasError() {
		t.Fatal(d)
	}

	req := sdk.CreateLoadBalancerMonitorRequestLoadBalancerMonitorConfig{}
	req.HttpMethod, req.HttpUri, req.HttpReturnCode, req.HttpPort = httpHealthCheckConfig(cfg)

	b, _ := json.Marshal(req)
	if want := `{"httpMethod":"GET","httpPort":8080,"httpUri":"/health"}`; string(b) != want {
		t.Errorf("config body = %s, want %s", b, want)
	}
}

// On update, a field removed from the configuration is sent as an explicit
// null so its recorded value is cleared; one that was never set stays unsent
// (MORPH-20773).
func TestUnitHTTPHealthCheckUpdateClearsRemovedFields(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	newCfg := func(method, uri, code attr.Value, port attr.Value) ConfigValue {
		cfg, d := NewConfigValue(ConfigValue{}.AttributeTypes(ctx), map[string]attr.Value{
			"monitor":          types.ObjectNull(map[string]attr.Type{"id": types.Int64Type}),
			"monitor_config":   types.DynamicNull(),
			"http_method":      method,
			"http_uri":         uri,
			"http_port":        port,
			"http_return_code": code,
		})
		if d.HasError() {
			t.Fatal(d)
		}

		return cfg
	}

	// Prior state: method, uri and port set; return code never set.
	prior := newCfg(types.StringValue("GET"), types.StringValue("/health"), types.StringNull(), types.Int64Value(8080))
	// Plan: uri and port removed, method kept, return code still unset.
	plan := newCfg(types.StringValue("GET"), types.StringNull(), types.StringNull(), types.Int64Null())

	req := sdk.UpdateLoadBalancerMonitorRequestLoadBalancerMonitorConfig{}
	req.HttpMethod, req.HttpUri, req.HttpReturnCode, req.HttpPort = httpHealthCheckUpdateConfig(plan, prior)

	b, _ := json.Marshal(req)
	if want := `{"httpMethod":"GET","httpPort":null,"httpUri":null}`; string(b) != want {
		t.Errorf("update config body = %s, want %s", b, want)
	}

	// No prior config block at all: nothing to clear, behaves as create.
	req = sdk.UpdateLoadBalancerMonitorRequestLoadBalancerMonitorConfig{}
	req.HttpMethod, req.HttpUri, req.HttpReturnCode, req.HttpPort = httpHealthCheckUpdateConfig(plan, NewConfigValueNull())

	b, _ = json.Marshal(req)
	if want := `{"httpMethod":"GET"}`; string(b) != want {
		t.Errorf("update config body without prior = %s, want %s", b, want)
	}
}

// Removing the whole config block clears every HTTP health-check value the
// prior state carried, as removing them one by one would; a prior block with
// none of them (or no prior block) sends no config at all.
func TestUnitHTTPHealthCheckWholeBlockRemovalClearsPriorFields(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	newCfg := func(method, uri, code attr.Value, port attr.Value) ConfigValue {
		cfg, d := NewConfigValue(ConfigValue{}.AttributeTypes(ctx), map[string]attr.Value{
			"monitor":          types.ObjectNull(map[string]attr.Type{"id": types.Int64Type}),
			"monitor_config":   types.DynamicNull(),
			"http_method":      method,
			"http_uri":         uri,
			"http_port":        port,
			"http_return_code": code,
		})
		if d.HasError() {
			t.Fatal(d)
		}

		return cfg
	}

	prior := newCfg(types.StringValue("GET"), types.StringValue("/health"), types.StringNull(), types.Int64Value(8080))

	cfg := httpHealthCheckClearConfig(prior)
	if cfg == nil {
		t.Fatal("prior HTTP values set, want a clearing config")
	}

	b, _ := json.Marshal(cfg)
	if want := `{"httpMethod":null,"httpPort":null,"httpUri":null}`; string(b) != want {
		t.Errorf("clearing config body = %s, want %s", b, want)
	}

	if top := topLevelConfig(cfg.ToMap()); string(mustJSON(t, top["config"])) != `{"httpMethod":null,"httpPort":null,"httpUri":null}` {
		t.Errorf("clearing config not sent at the top level: %v", top)
	}

	empty := newCfg(types.StringNull(), types.StringNull(), types.StringNull(), types.Int64Null())
	if cfg := httpHealthCheckClearConfig(empty); cfg != nil {
		t.Errorf("prior block without HTTP values produced a config: %v", cfg)
	}

	if cfg := httpHealthCheckClearConfig(NewConfigValueNull()); cfg != nil {
		t.Errorf("no prior block produced a config: %v", cfg)
	}
}

func mustJSON(t *testing.T, v any) []byte {
	t.Helper()

	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}

	return b
}

// The appliance stores a monitor's config from the request's top-level
// `config`, not from `loadBalancerMonitor.config`, so the block is sent at
// both places with the same content (MORPH-21017). An explicit null (a field
// cleared on update, MORPH-20773) reaches the top level too, and a request
// without a config block has no top-level key at all.
func TestUnitConfigBlockIsAlsoSentAtTheTopLevel(t *testing.T) {
	t.Parallel()

	createCfg := &sdk.CreateLoadBalancerMonitorRequestLoadBalancerMonitorConfig{}
	createCfg.HttpMethod.Set(sdk.PtrString("GET"))
	createCfg.HttpUri.Set(sdk.PtrString("/health"))
	createCfg.HttpPort.Set(sdk.PtrInt64(8080))
	createCfg.HttpReturnCode.Set(sdk.PtrString("200"))

	createReq := sdk.CreateLoadBalancerMonitorRequest{
		LoadBalancerMonitor: &sdk.CreateLoadBalancerMonitorRequestLoadBalancerMonitor{
			Name:   sdk.PtrString("m"),
			Config: createCfg,
		},
	}
	createReq.AdditionalProperties = topLevelConfig(createCfg.ToMap())

	var got map[string]json.RawMessage
	b, _ := json.Marshal(createReq)
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatal(err)
	}

	want := `{"httpMethod":"GET","httpPort":8080,"httpReturnCode":"200","httpUri":"/health"}`
	if string(got["config"]) != want {
		t.Errorf("create top-level config = %s, want %s", got["config"], want)
	}

	var nested struct {
		Config json.RawMessage `json:"config"`
	}
	if err := json.Unmarshal(got["loadBalancerMonitor"], &nested); err != nil {
		t.Fatal(err)
	}
	if string(nested.Config) != want {
		t.Errorf("create nested config = %s, want %s", nested.Config, want)
	}

	updateCfg := &sdk.UpdateLoadBalancerMonitorRequestLoadBalancerMonitorConfig{}
	updateCfg.HttpMethod.Set(sdk.PtrString("GET"))
	updateCfg.HttpUri.Set(nil)

	updateReq := sdk.UpdateLoadBalancerMonitorRequest{
		LoadBalancerMonitor: &sdk.UpdateLoadBalancerMonitorRequestLoadBalancerMonitor{Config: updateCfg},
	}
	updateReq.AdditionalProperties = topLevelConfig(updateCfg.ToMap())

	b, _ = json.Marshal(updateReq)
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatal(err)
	}
	if want := `{"httpMethod":"GET","httpUri":null}`; string(got["config"]) != want {
		t.Errorf("update top-level config = %s, want %s", got["config"], want)
	}

	if top := topLevelConfig((&sdk.CreateLoadBalancerMonitorRequestLoadBalancerMonitorConfig{}).ToMap()); top != nil {
		t.Errorf("empty config block produced a top-level key: %v", top)
	}
}
