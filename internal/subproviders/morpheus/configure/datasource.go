package configure

import (
	"context"

	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/clientfactory"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/constants"
	"github.com/HewlettPackard/hpe-morpheus-client/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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
) (*client.APIClient, error) {
	return r.cf.NewClient(ctx)
}
