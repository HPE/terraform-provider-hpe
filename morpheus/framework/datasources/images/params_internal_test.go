// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package images

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// emptyModel is a configuration with nothing set, which is what an empty
// `data "hpe_morpheus_images" "x" {}` block produces.
func emptyModel() ImagesModel {
	return ImagesModel{
		AllLabels:          types.StringNull(),
		Description:        types.SetNull(types.StringType),
		Direction:          types.StringNull(),
		FilterType:         types.StringNull(),
		ImageId:            types.SetNull(types.Int64Type),
		ImageType:          types.SetNull(types.StringType),
		IncludeSystemImage: types.BoolNull(),
		Labels:             types.StringNull(),
		Phrase:             types.StringNull(),
		Sort:               types.StringNull(),
		SystemImage:        types.BoolNull(),
	}
}

func strSet(t *testing.T, values ...string) types.Set {
	t.Helper()

	vals := make([]attr.Value, 0, len(values))
	for _, v := range values {
		vals = append(vals, types.StringValue(v))
	}

	set, d := types.SetValue(types.StringType, vals)
	if d.HasError() {
		t.Fatalf("building the fixture failed: %v", d)
	}

	return set
}

func int64Set(t *testing.T, values ...int64) types.Set {
	t.Helper()

	vals := make([]attr.Value, 0, len(values))
	for _, v := range values {
		vals = append(vals, types.Int64Value(v))
	}

	set, d := types.SetValue(types.Int64Type, vals)
	if d.HasError() {
		t.Fatalf("building the fixture failed: %v", d)
	}

	return set
}

// filter_type defaults to All, deliberately diverging from the platform's own
// default of User. Left to itself the API returns only user-uploaded images,
// hiding the synced and system images most instances are provisioned from.
func TestBuildServerSideParamsFilterTypeDefault(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("unset defaults to All", func(t *testing.T) {
		t.Parallel()

		var diags diag.Diagnostics

		config := emptyModel()

		if got := buildServerSideParams(ctx, &config, &diags).filterType; got != "All" {
			t.Errorf("filter_type = %q, want All", got)
		}

		if diags.HasError() {
			t.Errorf("unexpected diagnostics: %v", diags)
		}
	})

	t.Run("an explicit value is honoured", func(t *testing.T) {
		t.Parallel()

		var diags diag.Diagnostics

		config := emptyModel()
		config.FilterType = types.StringValue("User")

		if got := buildServerSideParams(ctx, &config, &diags).filterType; got != "User" {
			t.Errorf("filter_type = %q, want User", got)
		}
	})
}

// An unset argument must not be sent at all. Sending a zero value would filter
// on it — system_image=false is a real filter, not an absence.
func TestBuildServerSideParamsOmitsUnsetArguments(t *testing.T) {
	t.Parallel()

	var diags diag.Diagnostics

	config := emptyModel()
	p := buildServerSideParams(context.Background(), &config, &diags)

	if p.phrase != nil {
		t.Errorf("phrase = %v, want unset", *p.phrase)
	}

	if p.systemImage != nil {
		t.Errorf("system_image = %v, want unset", *p.systemImage)
	}

	if p.includeSystemImage != nil {
		t.Errorf("include_system_image = %v, want unset", *p.includeSystemImage)
	}

	if p.sort != nil {
		t.Errorf("sort = %v, want unset", *p.sort)
	}

	if p.direction != nil {
		t.Errorf("direction = %v, want unset", *p.direction)
	}

	if len(p.imageTypes) != 0 {
		t.Errorf("image_type = %v, want empty", p.imageTypes)
	}

	if len(p.imageIDs) != 0 {
		t.Errorf("image_id = %v, want empty", p.imageIDs)
	}
}

// system_image = false must be forwarded, not mistaken for unset. It selects
// non-system images, which is a different query from not filtering at all.
func TestBuildServerSideParamsForwardsFalse(t *testing.T) {
	t.Parallel()

	var diags diag.Diagnostics

	config := emptyModel()
	config.SystemImage = types.BoolValue(false)

	p := buildServerSideParams(context.Background(), &config, &diags)

	if p.systemImage == nil {
		t.Fatal("system_image was dropped, want false to be sent")
	}

	if *p.systemImage {
		t.Error("system_image = true, want false")
	}
}

// The repeatable arguments must survive as sets, since the API ORs them.
func TestBuildServerSideParamsRepeatableArguments(t *testing.T) {
	t.Parallel()

	var diags diag.Diagnostics

	config := emptyModel()
	config.ImageType = strSet(t, "qcow2", "raw")
	config.Description = strSet(t, "ubuntu%")
	config.ImageId = int64Set(t, 7, 42)

	p := buildServerSideParams(context.Background(), &config, &diags)

	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	if len(p.imageTypes) != 2 {
		t.Errorf("image_type = %v, want two values", p.imageTypes)
	}

	if len(p.descriptions) != 1 {
		t.Errorf("description = %v, want one value", p.descriptions)
	}

	if len(p.imageIDs) != 2 {
		t.Errorf("image_id = %v, want two values", p.imageIDs)
	}
}

// Every argument set at once, so no branch of the mapping goes unexercised.
func TestBuildServerSideParamsAllArguments(t *testing.T) {
	t.Parallel()

	var diags diag.Diagnostics

	config := emptyModel()
	config.FilterType = types.StringValue("Synced")
	config.ImageType = strSet(t, "vmware")
	config.Description = strSet(t, "prod%")
	config.ImageId = int64Set(t, 1)
	config.Phrase = types.StringValue("ubuntu")
	config.SystemImage = types.BoolValue(true)
	config.IncludeSystemImage = types.BoolValue(true)
	config.Labels = types.StringValue("gold")
	config.AllLabels = types.StringValue("gold,tested")
	config.Sort = types.StringValue("name")
	config.Direction = types.StringValue("desc")

	p := buildServerSideParams(context.Background(), &config, &diags)

	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	checks := map[string]struct{ got, want string }{
		"filter_type": {p.filterType, "Synced"},
		"phrase":      {derefStr(p.phrase), "ubuntu"},
		"labels":      {derefStr(p.labels), "gold"},
		"all_labels":  {derefStr(p.allLabels), "gold,tested"},
		"sort":        {derefStr(p.sort), "name"},
		"direction":   {derefStr(p.direction), "desc"},
	}

	for name, c := range checks {
		if c.got != c.want {
			t.Errorf("%s = %q, want %q", name, c.got, c.want)
		}
	}

	if p.systemImage == nil || !*p.systemImage {
		t.Error("system_image was not forwarded")
	}

	if p.includeSystemImage == nil || !*p.includeSystemImage {
		t.Error("include_system_image was not forwarded")
	}
}

func derefStr(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}

func TestStringsAndInt64sFromSet(t *testing.T) {
	ctx := context.Background()

	t.Run("null sets yield nothing", func(t *testing.T) {
		t.Parallel()

		var diags diag.Diagnostics

		if got := stringsFromSet(ctx, types.SetNull(types.StringType), &diags); got != nil {
			t.Errorf("strings = %v, want nil", got)
		}

		if got := int64sFromSet(ctx, types.SetNull(types.Int64Type), &diags); got != nil {
			t.Errorf("int64s = %v, want nil", got)
		}

		if diags.HasError() {
			t.Errorf("unexpected diagnostics: %v", diags)
		}
	})

	t.Run("unknown sets yield nothing", func(t *testing.T) {
		t.Parallel()

		var diags diag.Diagnostics

		if got := stringsFromSet(ctx, types.SetUnknown(types.StringType), &diags); got != nil {
			t.Errorf("strings = %v, want nil", got)
		}
	})
}
