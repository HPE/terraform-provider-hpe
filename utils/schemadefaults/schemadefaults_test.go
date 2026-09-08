package schemadefaults

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// syntheticSchema exercises every shape: flat string/bool/int64 defaults, an
// attribute with no default, single-nested, list-nested, set-nested, and two
// nested blocks sharing a leaf name (config_a/mode, config_b/mode) to prove
// distinct parents do not collide.
func syntheticSchema() schema.Schema {
	str := func(d string) schema.StringAttribute {
		return schema.StringAttribute{
			Optional: true,
			Computed: true,
			Default:  stringdefault.StaticString(d),
		}
	}

	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"visibility": str("private"),
			"enabled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(true),
			},
			"count": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(1),
			},
			"name": schema.StringAttribute{Optional: true}, // no default
			"config": schema.SingleNestedAttribute{
				Optional:   true,
				Computed:   true,
				Attributes: map[string]schema.Attribute{"cert": str("internal")},
			},
			"ifaces": schema.ListNestedAttribute{
				Optional: true,
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{"ip_mode": str("")},
				},
			},
			"rules": schema.SetNestedAttribute{
				Optional: true,
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{"action": str("allow")},
				},
			},
			"config_a": schema.SingleNestedAttribute{
				Optional:   true,
				Computed:   true,
				Attributes: map[string]schema.Attribute{"mode": str("a")},
			},
			"config_b": schema.SingleNestedAttribute{
				Optional:   true,
				Computed:   true,
				Attributes: map[string]schema.Attribute{"mode": str("b")},
			},
		},
	}
}

func TestCollectDefaults(t *testing.T) {
	ctx := context.Background()
	got := collectDefaults(ctx, syntheticSchema())

	want := map[string]string{
		"visibility":     `"private"`,
		"enabled":        "true",
		"count":          "1",
		"config/cert":    `"internal"`,
		"ifaces/ip_mode": `""`,
		"rules/action":   `"allow"`,
		"config_a/mode":  `"a"`, // distinct parent
		"config_b/mode":  `"b"`, // distinct parent, same leaf -> no collision
	}

	if len(got) != len(want) {
		t.Fatalf("got %d defaults, want %d: keys=%v", len(got), len(want), keys(got))
	}

	for k, w := range want {
		v, ok := got[k]
		if !ok {
			t.Errorf("missing default for %q", k)

			continue
		}

		if v.String() != w {
			t.Errorf("default[%q] = %s, want %s", k, v.String(), w)
		}
	}

	if _, ok := got["name"]; ok {
		t.Errorf("attribute without a default must not appear")
	}
}

func TestApplyFillsNullsAllShapes(t *testing.T) {
	ctx := context.Background()
	s := syntheticSchema()
	ty := s.Type().TerraformType(ctx)
	at := ty.(tftypes.Object).AttributeTypes

	nul := func(t tftypes.Type) tftypes.Value { return tftypes.NewValue(t, nil) }
	ifaceElem := at["ifaces"].(tftypes.List).ElementType
	ruleElem := at["rules"].(tftypes.Set).ElementType
	static := tftypes.NewValue(tftypes.String, "static")

	raw := tftypes.NewValue(ty, map[string]tftypes.Value{
		"visibility": nul(tftypes.String), // -> "private"
		"enabled":    nul(tftypes.Bool),   // -> true
		"count":      nul(tftypes.Number), // -> 1
		"name":       nul(tftypes.String), // no default -> stays null
		"config": tftypes.NewValue(at["config"], map[string]tftypes.Value{
			"cert": nul(tftypes.String), // -> "internal"
		}),
		"ifaces": tftypes.NewValue(at["ifaces"], []tftypes.Value{
			tftypes.NewValue(ifaceElem, map[string]tftypes.Value{"ip_mode": nul(tftypes.String)}), // -> ""
			tftypes.NewValue(ifaceElem, map[string]tftypes.Value{"ip_mode": static}),              // preserved
		}),
		"rules": tftypes.NewValue(at["rules"], []tftypes.Value{
			tftypes.NewValue(ruleElem, map[string]tftypes.Value{"action": nul(tftypes.String)}), // set: -> "allow"
		}),
		// config_a/mode null -> "a"
		"config_a": tftypes.NewValue(at["config_a"], map[string]tftypes.Value{
			"mode": nul(tftypes.String),
		}),
		// config_b is a null block -> its child must stay null.
		"config_b": nul(at["config_b"]),
	})

	state := &tfsdk.State{Schema: s, Raw: raw}
	if d := Apply(ctx, s, state); d.HasError() {
		t.Fatalf("Apply: %v", d)
	}

	leaves := leavesOf(state.Raw)

	check := func(match, want string) {
		t.Helper()

		for k, v := range leaves {
			if strings.Contains(k, match) {
				if v != want {
					t.Errorf("leaf %s = %s, want %s", k, v, want)
				}

				return
			}
		}

		t.Errorf("no leaf matching %q", match)
	}

	check(`AttributeName("visibility")`, `tftypes.String<"private">`)
	check(`AttributeName("enabled")`, `tftypes.Bool<"true">`)
	check(`AttributeName("count")`, `tftypes.Number<"1">`)
	check(`AttributeName("config").AttributeName("cert")`, `tftypes.String<"internal">`)
	check(`AttributeName("rules")`, `tftypes.String<"allow">`) // set element filled
	check(`AttributeName("config_a").AttributeName("mode")`, `tftypes.String<"a">`)
	check(`AttributeName("name")`, "tftypes.String<null>") // no default -> null

	// config_b was a null block -> its child must not have been visited/filled.
	for k := range leaves {
		if strings.Contains(k, `AttributeName("config_b")`) {
			t.Errorf("null-block child %s was filled; must be skipped", k)
		}
	}

	// ifaces[0] filled "", ifaces[1] preserved "static".
	var g0, g1 string

	for k, v := range leaves {
		if strings.Contains(k, "ElementKeyInt(0)") {
			g0 = v
		}

		if strings.Contains(k, "ElementKeyInt(1)") {
			g1 = v
		}
	}

	if g0 != `tftypes.String<"">` {
		t.Errorf("ifaces[0].ip_mode = %s, want empty string", g0)
	}

	if g1 != `tftypes.String<"static">` {
		t.Errorf("ifaces[1].ip_mode = %s, want preserved static", g1)
	}
}

func TestApplyNoOpOnNullState(t *testing.T) {
	ctx := context.Background()
	s := syntheticSchema()
	ty := s.Type().TerraformType(ctx)
	state := &tfsdk.State{Schema: s, Raw: tftypes.NewValue(ty, nil)}

	if d := Apply(ctx, s, state); d.HasError() {
		t.Fatalf("Apply: %v", d)
	}

	if !state.Raw.IsNull() {
		t.Errorf("Apply must leave a null state null")
	}
}

func keys(m map[string]attr.Value) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}

	return out
}

func leavesOf(v tftypes.Value) map[string]string {
	out := map[string]string{}

	_, _ = tftypes.Transform(v, func(p *tftypes.AttributePath, val tftypes.Value) (tftypes.Value, error) {
		ty := val.Type()
		if ty.Is(tftypes.String) || ty.Is(tftypes.Bool) || ty.Is(tftypes.Number) {
			out[p.String()] = val.String()
		}

		return val, nil
	})

	return out
}
