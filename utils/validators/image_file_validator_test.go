// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package validators_test

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/utils/validators"
)

// TestUnitImageFile exercises the ImageFile validator directly: a valid image
// passes, null/unknown are skipped, a disallowed extension, an empty file, a
// corrupt image, and a missing file are rejected, and formats the standard
// library cannot decode (.svg, .ico) are only checked for a non-empty file.
func TestUnitImageFile(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	dir := t.TempDir()

	writeFile := func(name string, data []byte) string {
		p := filepath.Join(dir, name)
		if err := os.WriteFile(p, data, 0o600); err != nil {
			t.Fatalf("writing %s: %v", name, err)
		}

		return p
	}

	var pngBuf bytes.Buffer
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.RGBA{R: 1, G: 2, B: 3, A: 255})
	if err := png.Encode(&pngBuf, img); err != nil {
		t.Fatal(err)
	}

	validPNG := writeFile("logo.png", pngBuf.Bytes())
	emptyPNG := writeFile("empty.png", nil)
	garbagePNG := writeFile("garbage.png", []byte("not a real png"))
	svgFile := writeFile("logo.svg", []byte("<svg xmlns='http://www.w3.org/2000/svg'></svg>"))
	emptySVG := writeFile("empty.svg", nil)
	icoFile := writeFile("favicon.ico", []byte{0x00, 0x00, 0x01, 0x00})

	cases := []struct {
		name      string
		exts      []string
		value     types.String
		wantError bool
	}{
		{"valid png", []string{".png", ".jpg"}, types.StringValue(validPNG), false},
		{"null skipped", []string{".png"}, types.StringNull(), false},
		{"unknown skipped", []string{".png"}, types.StringUnknown(), false},
		{"disallowed extension", []string{".png"}, types.StringValue(svgFile), true},
		{"empty png", []string{".png"}, types.StringValue(emptyPNG), true},
		{"corrupt png", []string{".png"}, types.StringValue(garbagePNG), true},
		{"missing file", []string{".png"}, types.StringValue(filepath.Join(dir, "nope.png")), true},
		{"svg not decoded", []string{".svg"}, types.StringValue(svgFile), false},
		{"empty svg rejected", []string{".svg"}, types.StringValue(emptySVG), true},
		{"ico not decoded", []string{".ico", ".png"}, types.StringValue(icoFile), false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			v := validators.ImageFile(tc.exts...)
			req := validator.StringRequest{
				Path:        path.Root("header_logo"),
				ConfigValue: tc.value,
			}
			resp := &validator.StringResponse{}
			v.ValidateString(ctx, req, resp)

			if got := resp.Diagnostics.HasError(); got != tc.wantError {
				t.Fatalf("value %v (exts %v): got error=%v, want error=%v (diagnostics: %v)",
					tc.value, tc.exts, got, tc.wantError, resp.Diagnostics)
			}
		})
	}
}
