// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package images

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

func strp(s string) *string { return &s }
func i64p(i int64) *int64   { return &i }
func boolp(b bool) *bool    { return &b }

// img builds a minimal image with the fields the filters read.
func img(name, imageType, status, visibility string) *sdk.ListVirtualImages200ResponseAllOfVirtualImagesInner {
	return &sdk.ListVirtualImages200ResponseAllOfVirtualImagesInner{
		Name:       strp(name),
		ImageType:  strp(imageType),
		Status:     strp(status),
		Visibility: strp(visibility),
	}
}

func TestFieldValue(t *testing.T) {
	t.Parallel()

	image := img("ubuntu-22", "qcow2", "active", "public")
	image.Description = *sdk.NewNullableString(strp("an image"))

	tests := []struct {
		field   string
		want    string
		wantOK  bool
		comment string
	}{
		{field: "name", want: "ubuntu-22", wantOK: true},
		{field: "description", want: "an image", wantOK: true},
		{field: "image_type", want: "qcow2", wantOK: true},
		{field: "status", want: "active", wantOK: true},
		{field: "visibility", want: "public", wantOK: true},
		{
			field:   "unknown_field",
			wantOK:  false,
			comment: "the schema restricts the name, so this is unreachable in practice",
		},
	}

	for _, tt := range tests {
		t.Run(tt.field, func(t *testing.T) {
			t.Parallel()

			got, ok := fieldValue(image, tt.field)
			if ok != tt.wantOK {
				t.Fatalf("ok = %v, want %v", ok, tt.wantOK)
			}

			if ok && got != tt.want {
				t.Errorf("value = %q, want %q", got, tt.want)
			}
		})
	}
}

// An image with no value for a field cannot match any expression, including one
// that would match the empty string.
func TestFieldValueAbsent(t *testing.T) {
	t.Parallel()

	empty := &sdk.ListVirtualImages200ResponseAllOfVirtualImagesInner{}

	for _, field := range []string{
		"name", "description", "image_type", "status", "visibility",
	} {
		if _, ok := fieldValue(empty, field); ok {
			t.Errorf("%s: reported a value on an image that has none", field)
		}
	}
}

func TestCompileFilters(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("null set compiles to nothing", func(t *testing.T) {
		t.Parallel()

		var diags diag.Diagnostics

		if got := compileFilters(ctx, types.SetNull(FilterValue{}.Type(ctx)), &diags); got != nil {
			t.Errorf("compiled %v, want nil", got)
		}

		if diags.HasError() {
			t.Errorf("unexpected diagnostics: %v", diags)
		}
	})

	t.Run("invalid expression is reported, not ignored", func(t *testing.T) {
		t.Parallel()

		var diags diag.Diagnostics

		set := filterSet(t, ctx, &diags, "name", []string{"[unterminated"})
		if diags.HasError() {
			t.Fatalf("building the fixture failed: %v", diags)
		}

		compileFilters(ctx, set, &diags)

		if !diags.HasError() {
			t.Fatal("an invalid regular expression compiled without error")
		}
	})
}

func TestMatchesFilters(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	image := img("ubuntu-22.04", "qcow2", "active", "public")

	tests := []struct {
		name    string
		filters []compiledFilter
		want    bool
	}{
		{
			name: "no filters matches everything",
			want: true,
		},
		{
			name:    "single matching expression",
			filters: compiled(t, "name", "^ubuntu"),
			want:    true,
		},
		{
			// Documented behaviour: expressions are unanchored, so a bare
			// substring matches anywhere in the value.
			name:    "expressions are unanchored",
			filters: compiled(t, "name", "22"),
			want:    true,
		},
		{
			name:    "non-matching expression excludes",
			filters: compiled(t, "name", "^debian"),
			want:    false,
		},
		{
			name:    "values within a block are ORed",
			filters: compiled(t, "name", "^debian", "^ubuntu"),
			want:    true,
		},
		{
			name: "blocks are ANDed",
			filters: append(
				compiled(t, "name", "^ubuntu"),
				compiled(t, "status", "active")...),
			want: true,
		},
		{
			name: "one failing block excludes",
			filters: append(
				compiled(t, "name", "^ubuntu"),
				compiled(t, "status", "^inactive$")...),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := matchesFilters(image, tt.filters); got != tt.want {
				t.Errorf("matched = %v, want %v", got, tt.want)
			}
		})
	}

	_ = ctx
}

// min_disk and min_ram are reported in gigabytes while the API returns bytes.
// The singular data source does the same, and the two must agree.
func TestBytesToGB(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *int64
		want any
	}{
		{name: "absent stays null", in: nil, want: nil},
		{name: "exact gigabyte", in: i64p(10 * 1024 * 1024 * 1024), want: int64(10)},
		{
			// Integer division truncates, which is what the singular does.
			name: "part gigabyte truncates",
			in:   i64p(10*1024*1024*1024 + 1),
			want: int64(10),
		},
		{name: "under a gigabyte becomes zero", in: i64p(1024), want: int64(0)},
		{name: "zero", in: i64p(0), want: int64(0)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := bytesToGB(tt.in)

			if tt.want == nil {
				if !got.IsNull() {
					t.Errorf("value = %v, want null", got)
				}

				return
			}

			if got.ValueInt64() != tt.want.(int64) {
				t.Errorf("value = %d, want %d", got.ValueInt64(), tt.want)
			}
		})
	}
}

// config_azure is populated only for azure-reference images. Every other image
// reports it null rather than as a set of empty strings.
func TestConfigAzureValue(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("non-azure image is null", func(t *testing.T) {
		t.Parallel()

		if got := configAzureValue(ctx, img("x", "qcow2", "active", "public")); !got.IsNull() {
			t.Errorf("value = %v, want null", got)
		}
	})

	t.Run("azure image with no config is null", func(t *testing.T) {
		t.Parallel()

		image := img("x", "azure-reference", "active", "public")

		if got := configAzureValue(ctx, image); !got.IsNull() {
			t.Errorf("value = %v, want null", got)
		}
	})

	t.Run("azure image with config is populated", func(t *testing.T) {
		t.Parallel()

		image := img("x", "azure-reference", "active", "public")
		image.Config = &sdk.ListVirtualImages200ResponseAllOfVirtualImagesInnerConfig{
			AzureReferenceVirtualImageConfiguration: &sdk.AzureReferenceVirtualImageConfiguration{
				Publisher: "canonical",
				Offer:     "ubuntu",
				Version:   "latest",
				Sku:       "22_04-lts",
			},
		}

		got := configAzureValue(ctx, image)
		if got.IsNull() {
			t.Fatal("value is null, want populated")
		}

		attrs := got.Attributes()
		for key, want := range map[string]string{
			"publisher": "canonical",
			"offer":     "ubuntu",
			"version":   "latest",
			"sku":       "22_04-lts",
		} {
			v, ok := attrs[key].(types.String)
			if !ok {
				t.Errorf("%s: not a string", key)

				continue
			}

			if v.ValueString() != want {
				t.Errorf("%s = %q, want %q", key, v.ValueString(), want)
			}
		}
	})
}

// The API orders by name, names are not unique across clouds, and offset paging
// can therefore serve the same record twice. The read de-duplicates by id.
func TestImageToValueOnMinimalImage(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// Nothing set beyond an id: every other field must come back null rather
	// than panicking on a nil pointer.
	minimal := &sdk.ListVirtualImages200ResponseAllOfVirtualImagesInner{Id: i64p(7)}

	v, diags := imageToValue(ctx, minimal)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	if v.Id.ValueInt64() != 7 {
		t.Errorf("id = %d, want 7", v.Id.ValueInt64())
	}

	if !v.Name.IsNull() {
		t.Errorf("name = %v, want null", v.Name)
	}

	if !v.MinDisk.IsNull() {
		t.Errorf("min_disk = %v, want null", v.MinDisk)
	}

	if !v.ConfigAzure.IsNull() {
		t.Errorf("config_azure = %v, want null", v.ConfigAzure)
	}
}

func TestImageToValueMapsFields(t *testing.T) {
	t.Parallel()

	image := img("ubuntu-22", "qcow2", "active", "public")
	image.Id = i64p(42)
	image.MinDisk = *sdk.NewNullableInt64(i64p(10 * 1024 * 1024 * 1024))
	image.UserUploaded = boolp(true)
	image.VirtioSupported = boolp(true)

	v, diags := imageToValue(context.Background(), image)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	if v.Name.ValueString() != "ubuntu-22" {
		t.Errorf("name = %q, want ubuntu-22", v.Name.ValueString())
	}

	if v.MinDisk.ValueInt64() != 10 {
		t.Errorf("min_disk = %d, want 10 (GB)", v.MinDisk.ValueInt64())
	}

	if !v.UserUploaded.ValueBool() {
		t.Error("user_uploaded = false, want true")
	}

	// Guards the bug this replaced on the singular data source, where the
	// virtio_supported block assigned user_data a second time.
	if !v.VirtioSupported.ValueBool() {
		t.Error("virtio_supported = false, want true")
	}
}

// Offset paging over a non-unique sort key can serve the same record twice: the
// API orders by name, names are not unique across clouds, and rows move between
// pages as ties are broken differently. selectImages is what stops a duplicate
// reaching state.
func TestSelectImagesDeduplicatesById(t *testing.T) {
	t.Parallel()

	dup := func(id int64, name string) sdk.ListVirtualImages200ResponseAllOfVirtualImagesInner {
		return sdk.ListVirtualImages200ResponseAllOfVirtualImagesInner{
			Id: i64p(id), Name: strp(name),
		}
	}

	tests := []struct {
		name    string
		images  []sdk.ListVirtualImages200ResponseAllOfVirtualImagesInner
		wantIDs []int64
	}{
		{
			name:    "no duplicates passes everything through",
			images:  []sdk.ListVirtualImages200ResponseAllOfVirtualImagesInner{dup(1, "a"), dup(2, "b")},
			wantIDs: []int64{1, 2},
		},
		{
			// The same row served on two pages.
			name:    "repeated id is kept once",
			images:  []sdk.ListVirtualImages200ResponseAllOfVirtualImagesInner{dup(1, "a"), dup(1, "a"), dup(2, "b")},
			wantIDs: []int64{1, 2},
		},
		{
			// Distinct images sharing a name are not duplicates.
			name:    "same name, different ids are both kept",
			images:  []sdk.ListVirtualImages200ResponseAllOfVirtualImagesInner{dup(1, "ubuntu"), dup(2, "ubuntu")},
			wantIDs: []int64{1, 2},
		},
		{
			// The first occurrence wins, so order is the order served.
			name:    "order is preserved",
			images:  []sdk.ListVirtualImages200ResponseAllOfVirtualImagesInner{dup(3, "c"), dup(1, "a"), dup(3, "c")},
			wantIDs: []int64{3, 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := selectImages(tt.images, nil)

			if len(got) != len(tt.wantIDs) {
				t.Fatalf("kept %d images, want %d", len(got), len(tt.wantIDs))
			}

			for i, want := range tt.wantIDs {
				if got[i].Id == nil || *got[i].Id != want {
					t.Errorf("position %d has id %v, want %d", i, got[i].Id, want)
				}
			}
		})
	}
}

// An image without an id cannot be de-duplicated. Keeping it is the lesser
// failure: dropping records on a missing field would lose data outright.
func TestSelectImagesKeepsImagesWithoutAnId(t *testing.T) {
	t.Parallel()

	got := selectImages([]sdk.ListVirtualImages200ResponseAllOfVirtualImagesInner{
		{Name: strp("no-id-1")},
		{Name: strp("no-id-2")},
	}, nil)

	if len(got) != 2 {
		t.Errorf("kept %d images, want 2 — images without an id were dropped", len(got))
	}
}

// De-duplication and filtering happen together, so a duplicate that the filters
// would reject must not consume the id that a later valid row needs.
func TestSelectImagesAppliesFiltersAfterDeduplication(t *testing.T) {
	t.Parallel()

	images := []sdk.ListVirtualImages200ResponseAllOfVirtualImagesInner{
		{Id: i64p(1), Name: strp("ubuntu-22")},
		{Id: i64p(1), Name: strp("ubuntu-22")},
		{Id: i64p(2), Name: strp("debian-12")},
	}

	got := selectImages(images, compiled(t, "name", "^ubuntu"))

	if len(got) != 1 {
		t.Fatalf("kept %d images, want 1", len(got))
	}

	if got[0].Id == nil || *got[0].Id != 1 {
		t.Errorf("kept id %v, want 1", got[0].Id)
	}
}
