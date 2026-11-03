// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package plugin

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
)

const (
	resourceName = "plugin"

	// Plugin lifecycle statuses reported by the API.
	statusLoaded     = "loaded"
	statusError      = "error"
	statusInstalling = "installing"
	statusUpdating   = "updating"

	// settleTimeout bounds the wait for a plugin to report loaded again after
	// its enabled flag or configuration changed.
	settleTimeout = 5 * time.Minute
)

type Resource struct {
	configure.ResourceWithMorpheusConfigure
}

var (
	_ resource.Resource                     = &Resource{}
	_ resource.ResourceWithConfigure        = &Resource{}
	_ resource.ResourceWithImportState      = &Resource{}
	_ resource.ResourceWithModifyPlan       = &Resource{}
	_ resource.ResourceWithConfigValidators = &Resource{}
)

func NewResource() resource.Resource {
	return &Resource{}
}

func (r *Resource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_" + resourceName
}

func (r *Resource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = PluginResourceSchema(ctx)
}
