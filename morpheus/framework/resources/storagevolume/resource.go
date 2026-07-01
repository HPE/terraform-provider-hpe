// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package storagevolume

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
)

var (
	_ resource.Resource                = &storageVolumeResource{}
	_ resource.ResourceWithConfigure   = &storageVolumeResource{}
	_ resource.ResourceWithImportState = &storageVolumeResource{}
)

type storageVolumeResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &storageVolumeResource{}
}

func (r *storageVolumeResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_" + "storage_volume"
}

func (r *storageVolumeResource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = StorageVolumeResourceSchema(ctx)
}

// oneGibibyte is the number of bytes in 1 GiB. The Morpheus API stores and
// returns storage volume sizes in bytes, while the resource expresses
// max_storage in GiB.
const oneGibibyte int64 = 1024 * 1024 * 1024
