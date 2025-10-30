// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package morpheus

import (
	"context"
	"os"

	sdklegacy "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"

	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/client"
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
		c.client = client.NewLegacyClient(
			context.Background(),
			c.Url,
			c.Username,
			c.Password,
			c.AccessToken,
			sdklegacy.SkipLogin(),
			sdklegacy.WithDebug(debug),
			sdklegacy.WithInsecure(c.Insecure),
		)
	}

	return c.client, diags
}
