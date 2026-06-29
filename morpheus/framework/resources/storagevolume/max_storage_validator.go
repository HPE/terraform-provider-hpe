// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package storagevolume

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	// alletraVolumeTypeCodePrefix identifies HPE Alletra MP / 9000 block volume
	// type codes (hpealletraMPLUN, hpealletraMPLUN-active-pp,
	// hpealletraMPLUN-classic-pp).
	alletraVolumeTypeCodePrefix = "hpealletraMPLUN"
	// alletraMaxVolumeSizeGiB is the maximum HPE Alletra MP / 9000 volume size,
	// in GiB.
	alletraMaxVolumeSizeGiB int64 = 65536
	// minVolumeSizeGiB is the minimum volume size, in GiB. (HPE Alletra 9000
	// models 9060/9080 enforce a 16 GiB minimum on the array itself.)
	minVolumeSizeGiB int64 = 1
)

// maxStorageSizeValidator validates max_storage (in GiB) using limits that
// depend on the storage volume type. HPE Alletra MP / 9000 volumes must be
// between 1 and 65536 GiB; other types only require a positive size, leaving any
// type-specific maximum to be enforced by the backend.
type maxStorageSizeValidator struct{}

// maxStorageSize returns a type-aware validator for the max_storage attribute.
func maxStorageSize() validator.Int64 {
	return maxStorageSizeValidator{}
}

func (v maxStorageSizeValidator) Description(_ context.Context) string {
	return "max_storage must be a positive number of GiB; HPE Alletra MP and " +
		"Alletra 9000 volumes must be between 1 and 65536 GiB"
}

func (v maxStorageSizeValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v maxStorageSizeValidator) ValidateInt64(
	ctx context.Context,
	req validator.Int64Request,
	resp *validator.Int64Response,
) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	// The type is identified by type_code; when only type_id is set (or the code
	// is unknown), the Alletra upper bound is left to the backend to enforce.
	var typeCode types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("type_code"), &typeCode)...)
	if resp.Diagnostics.HasError() {
		return
	}

	typeCodeKnown := !typeCode.IsNull() && !typeCode.IsUnknown()

	if title, detail := maxStorageSizeError(
		req.ConfigValue.ValueInt64(), typeCode.ValueString(), typeCodeKnown,
	); title != "" {
		resp.Diagnostics.AddAttributeError(req.Path, title, detail)
	}
}

// maxStorageSizeError returns a non-empty (title, detail) when size (GiB) is
// invalid for the given storage volume type, or empty strings when it is valid.
// A positive size is required for every type; HPE Alletra MP / 9000 block
// volumes additionally must not exceed 65536 GiB.
func maxStorageSizeError(size int64, typeCode string, typeCodeKnown bool) (string, string) {
	if size < minVolumeSizeGiB {
		return "Invalid max_storage",
			fmt.Sprintf("max_storage must be at least %d GiB, got %d.", minVolumeSizeGiB, size)
	}

	if typeCodeKnown &&
		strings.HasPrefix(typeCode, alletraVolumeTypeCodePrefix) &&
		size > alletraMaxVolumeSizeGiB {
		return "Invalid max_storage for HPE Alletra volume",
			fmt.Sprintf(
				"HPE Alletra MP / 9000 volumes must be between %d and %d GiB, got %d.",
				minVolumeSizeGiB, alletraMaxVolumeSizeGiB, size,
			)
	}

	return "", ""
}
