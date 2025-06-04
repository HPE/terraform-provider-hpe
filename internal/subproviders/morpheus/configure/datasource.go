// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package configure

import (
	"context"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/clientfactory"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/constants"
)

type DataSourceWithMorpheusConfigure struct {
	cf clientfactory.ClientFactory
}

func (r *DataSourceWithMorpheusConfigure) BlockName() string {
	return constants.SubProviderName
}

func (r *DataSourceWithMorpheusConfigure) Configure(
	ctx context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) {
	// provider.Configure is not guaranteed to have run yet
	if req.ProviderData == nil {
		return
	}

	m, _ := req.ProviderData.(map[string]any)
	cf, ok := m[constants.SubProviderName].(*clientfactory.ClientFactory)
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
			constants.SubProviderName+" client creation failed",
			msg,
		)

		return
	}

	r.cf = *cf
}

func (r *DataSourceWithMorpheusConfigure) NewClient(
	ctx context.Context,
) (*sdk.APIClient, error) {
	return r.cf.NewClient(ctx)
}
