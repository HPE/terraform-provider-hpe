package morpheus

import (
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/constants"
	"github.com/HPE/terraform-provider-hpe/subprovider"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ subprovider.SubProvider = (*SubProvider)(nil)

type SubProvider struct{}

func New() subprovider.SubProvider {
	return SubProvider{}
}

func (SubProvider) Configure(_ func(x any)) error {
	return nil
}

func (SubProvider) GetName() string {
	return constants.SubProviderName
}

func (SubProvider) GetSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{}
}

func (SubProvider) GetDataSources() []func() datasource.DataSource {
	return nil
}

func (SubProvider) GetResources() []func() resource.Resource {
	return nil
}
