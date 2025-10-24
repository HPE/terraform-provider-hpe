// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package morpheus

import (
	"os"

	sdklegacy "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
)

// Config is the configuration structure used to instantiate the Morpheus
// provider.  Only Url and AccessToken are required.
type Config struct {
	Url          string
	AccessToken  string
	RefreshToken string // optional and unused
	Username     string
	Password     string
	ClientId     string
	// TODO need to decide if we want to keep tenant subdomain here
	// TenantSubdomain string
	Insecure bool

	client *sdklegacy.Client
}

func (c *Config) Client() (*sdklegacy.Client, diag.Diagnostics) {
	debug := logging.IsDebugOrHigher() && os.Getenv("MORPHEUS_API_HTTPTRACE") == "true"
	diags := diag.Diagnostics{}

	if c.client == nil {

		var client *sdklegacy.Client
		if !c.Insecure {
			// default: secure
			client = sdklegacy.NewClient(
				c.Url,
				sdklegacy.SkipLogin(), // required for custom auth transport
				sdklegacy.WithDebug(debug),
			)
		} else {
			// insecure
			client = sdklegacy.NewClient(
				c.Url,
				sdklegacy.SkipLogin(), // required for custom auth transport
				sdklegacy.WithDebug(debug),
				sdklegacy.Insecure(),
			)
		}

		c.client = client
	}

	return c.client, diags
}
