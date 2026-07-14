// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package validators

import (
	"bytes"
	"context"
	"fmt"
	"image"

	// Register the standard-library image decoders so DecodeConfig can validate
	// the corresponding formats.
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = ImageFileValidator{}

// decodableExtensions are the image extensions the Go standard library can
// decode, and whose contents are therefore validated. Other allowed extensions
// (for example .svg or .ico) are only checked for a valid extension and a
// non-empty file.
var decodableExtensions = map[string]bool{
	".png":  true,
	".jpg":  true,
	".jpeg": true,
	".gif":  true,
}

// ImageFileValidator validates that a configured string is a path to a readable,
// non-empty image file with one of the allowed extensions. For formats the Go
// standard library can decode it also verifies the file is a valid image, so an
// empty or corrupt file is rejected at plan time instead of silently resetting
// the image on the Morpheus appliance.
type ImageFileValidator struct {
	AllowedExtensions []string
}

func (v ImageFileValidator) Description(_ context.Context) string {
	return fmt.Sprintf(
		"must be a path to a non-empty image file with one of these extensions: %s",
		strings.Join(v.AllowedExtensions, ", "),
	)
}

func (v ImageFileValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v ImageFileValidator) ValidateString(
	_ context.Context,
	request validator.StringRequest,
	response *validator.StringResponse,
) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()
	ext := strings.ToLower(filepath.Ext(value))

	if !v.extensionAllowed(ext) {
		response.Diagnostics.AddAttributeError(
			request.Path,
			"Invalid Image File Extension",
			fmt.Sprintf("Attribute %s must reference an image file with one of these "+
				"extensions: %s, got: %s", request.Path, strings.Join(v.AllowedExtensions, ", "), value),
		)

		return
	}

	data, err := os.ReadFile(value)
	if err != nil {
		response.Diagnostics.AddAttributeError(
			request.Path,
			"Cannot Read Image File",
			fmt.Sprintf("Attribute %s could not read image file %q: %s", request.Path, value, err),
		)

		return
	}

	if len(data) == 0 {
		response.Diagnostics.AddAttributeError(
			request.Path,
			"Empty Image File",
			fmt.Sprintf("Attribute %s image file %q is empty; uploading an empty file "+
				"resets the image on the appliance.", request.Path, value),
		)

		return
	}

	// Only formats the standard library supports can have their contents checked.
	if !decodableExtensions[ext] {
		return
	}

	if _, _, err := image.DecodeConfig(bytes.NewReader(data)); err != nil {
		response.Diagnostics.AddAttributeError(
			request.Path,
			"Invalid Image File",
			fmt.Sprintf("Attribute %s file %q is not a valid image: %s", request.Path, value, err),
		)
	}
}

func (v ImageFileValidator) extensionAllowed(ext string) bool {
	for _, allowed := range v.AllowedExtensions {
		if strings.EqualFold(ext, allowed) {
			return true
		}
	}

	return false
}

// ImageFile returns a validator that ensures a string is a path to a readable,
// non-empty image file with one of the allowed extensions (for example ".png"
// or ".jpg"). For standard-library-decodable formats it also validates that the
// file contains a valid image.
func ImageFile(allowedExtensions ...string) ImageFileValidator {
	return ImageFileValidator{AllowedExtensions: allowedExtensions}
}
