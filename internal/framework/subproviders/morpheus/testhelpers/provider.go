package testhelpers

import (
	"context"
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/HPE/terraform-provider-hpe/internal/framework/provider"
	"github.com/HPE/terraform-provider-hpe/subprovider"
)

const NoProvidersErr = "at least one provider must be non-nil"

// GetAccTestFactories returns a map of provider factories for acceptance testing.
// The first parameter is a testing.T test state struct.
// The second parameter is a SubProvider implemented using the Terraform Plugin Framework.
// The third parameter is a Terraform SDK v2 Provider.
// At least one of the parameters must be non-nil.
func GetAccTestFactories(
	t *testing.T,
	fw subprovider.SubProvider,
	sdkv2 *schema.Provider,
) map[string]func() (tfprotov6.ProviderServer, error) {
	t.Helper()

	factories, err := getAccTestFactories(fw, sdkv2)
	if err != nil {
		t.Fatal(err)
	}

	return factories
}

func getAccTestFactories(fw subprovider.SubProvider, sdkv2 *schema.Provider) (
	map[string]func() (tfprotov6.ProviderServer, error), error,
) {
	if fw == nil && sdkv2 == nil {
		return nil, errors.New(NoProvidersErr)
	}

	var providers []func() tfprotov6.ProviderServer

	if fw != nil {
		frameworkProvider := provider.New("test", []subprovider.SubProvider{fw}...)()

		frameworkServer, err := providerserver.NewProtocol6WithError(frameworkProvider)()
		if err != nil {
			return nil, err
		}

		providers = append(providers, func() tfprotov6.ProviderServer { return frameworkServer })
	}

	if sdkv2 != nil {
		sdkv2Server, err := tf5to6server.UpgradeServer(context.Background(), sdkv2.GRPCProvider)
		if err != nil {
			return nil, err
		}

		providers = append(providers, func() tfprotov6.ProviderServer { return sdkv2Server })
	}

	factory := func() (tfprotov6.ProviderServer, error) {
		return tf6muxserver.NewMuxServer(context.Background(), providers...)
	}

	testAccProtoV6ProviderFactories := map[string]func() (
		tfprotov6.ProviderServer, error,
	){
		"hpe": factory,
	}

	return testAccProtoV6ProviderFactories, nil
}
