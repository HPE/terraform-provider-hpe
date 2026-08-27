// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package configure

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/clientfactory"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/constants"
)

type ResourceWithMorpheusConfigure struct {
	cf clientfactory.ClientFactory
}

func (r *ResourceWithMorpheusConfigure) BlockName() string {
	return constants.ProviderName
}

func (r *ResourceWithMorpheusConfigure) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	// provider.Configure is not guaranteed to have run yet
	if req.ProviderData == nil {
		return
	}

	cf, ok := req.ProviderData.(*clientfactory.ClientFactory)
	if !ok {
		tflog.Debug(ctx, "Nil ProviderData sub block")
		msg := `
Morpheus resource present, but possible missing morpheus provider block.

provider "hpe" {
  morpheus { <- missing or duplicate?
    url = "https://example.com"
  }
}`
		resp.Diagnostics.AddError(
			constants.ProviderName+" client creation failed",
			msg,
		)

		return
	}

	r.cf = *cf
}

func (r *ResourceWithMorpheusConfigure) NewClient(
	ctx context.Context,
) (*sdk.APIClient, error) {
	return r.cf.NewClient(ctx)
}
