// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package pce implements the PCE Identity token exchange. It trades GreenLake
// API client credentials for the URL and access token of the Morpheus instance
// sitting behind the VMaaS broker.
//
// The same exchange serves both PCE deployment types. They differ only in the
// IAM dialect used to mint the token (GLCS for Connected, GLP for Disconnected)
// and in how the broker request is scoped; see Config.Version.
package pce

import (
	"context"
	"fmt"
	"sync"

	brokerclient "github.com/HPE/terraform-provider-hpe/morpheus/pce/sdk/broker/client"
	"github.com/HPE/terraform-provider-hpe/morpheus/pce/sdk/token/iamversion"
	"github.com/HPE/terraform-provider-hpe/morpheus/pce/sdk/token/retrieve"
	"github.com/HPE/terraform-provider-hpe/morpheus/pce/sdk/token/serviceclient"
)

// DefaultBrokerURL is the VMaaS broker used when a Config does not name one.
const DefaultBrokerURL = "https://vmaas-broker.us1.greenlake-hpe.com"

// Config is the PCE Identity configuration needed to resolve Morpheus
// connection details. It deliberately holds plain strings rather than Terraform
// types so that both the framework and the SDKv2 Morpheus provider can use it.
//
// Config is also used as the key of the result cache, so any field that changes
// which Morpheus instance is resolved must be represented here. Adding a
// targeting field without adding it to this struct would let two different
// configurations share a cached result.
type Config struct {
	ClientID     string
	ClientSecret string
	Location     string
	Space        string
	IssuerURL    string
	IAMToken     string
	BrokerURL    string

	// Version selects the IAM dialect and how the broker request is scoped.
	// Connected PCE uses iamversion.GLCS and scopes by Space; Disconnected PCE
	// uses iamversion.GLP and scopes by WorkspaceID. There is deliberately no
	// default: an unset Version would silently pick a dialect for the caller.
	Version iamversion.Version

	// WorkspaceID is the GreenLake Platform workspace ID. It scopes the broker
	// request for Disconnected PCE and is unused for Connected PCE, which
	// scopes by Space instead.
	WorkspaceID string
}

type result struct {
	url   string
	token string
}

// The cache holds credentials for the lifetime of the process. Terraform runs a
// provider process per provider configuration and destroys it at the end of
// each graph walk, so entries are short lived.
var (
	cacheMu sync.Mutex
	cache   = map[Config]result{}
)

// TokenExchange returns the Morpheus URL and access token for cfg.
//
// Results are memoised per configuration. The "hpe" provider muxes a framework
// provider and an SDKv2 provider, and Terraform configures both from the same
// configuration, so without memoisation the exchange would run twice per graph
// walk. Keying on the configuration also keeps provider blocks that share a
// process, which happens under TF_REATTACH_PROVIDERS, independent of each
// other.
//
// The returned access token is used as-is and is not refreshed, so callers are
// expected to finish their work within its validity period.
func TokenExchange(ctx context.Context, cfg Config) (string, string, error) {
	// Held across the exchange so that concurrent callers for the same
	// configuration do not each perform it. Terraform normally runs a single
	// provider configuration per process, so contention is not expected.
	cacheMu.Lock()
	defer cacheMu.Unlock()

	if r, ok := cache[cfg]; ok {
		return r.url, r.token, nil
	}

	url, token, err := cfg.exchange(ctx)
	if err != nil {
		// Deliberately not cached, so that a later attempt can retry.
		return "", "", err
	}

	cache[cfg] = result{url: url, token: token}

	return url, token, nil
}

// exchange obtains a GreenLake IAM token and trades it with the VMaaS broker
// for Morpheus connection details.
func (c Config) exchange(ctx context.Context) (string, string, error) {
	iamToken, err := c.iamToken(ctx)
	if err != nil {
		return "", "", err
	}

	brokerURL := c.BrokerURL
	if brokerURL == "" {
		brokerURL = DefaultBrokerURL
	}

	brokerCfg := brokerclient.NewConfiguration()
	brokerCfg.Host = brokerURL

	if c.Location != "" {
		brokerCfg.DefaultQueryParams["location"] = c.Location
	}

	// The two deployment types scope the broker request differently: Connected
	// PCE by GLCS space, Disconnected PCE by GLP workspace. The workspace is
	// sent both as a query parameter and as a header, which is what the broker
	// expects for GLP.
	switch c.Version {
	case iamversion.GLP:
		if c.WorkspaceID != "" {
			brokerCfg.DefaultQueryParams["tenantID"] = c.WorkspaceID
			brokerCfg.AddDefaultHeader("X-Tenant-ID", c.WorkspaceID)
		}
	default:
		if c.Space != "" {
			brokerCfg.DefaultQueryParams["space"] = c.Space
		}
	}

	// Note: this client is deliberately not wrapped with utils/httptrace. That
	// tracer dumps full request and response bodies, which here would write the
	// GreenLake client secret and the Morpheus access token to the Terraform
	// log. Do not add tracing to this path without redacting them.
	brokerClient := brokerclient.NewAPIClient(brokerCfg)

	// Inject the GreenLake IAM token on every request's context, which
	// prepareRequest reads for Bearer auth.
	brokerClient.SetMetaFnAndVersion(nil, 0, func(ctx *context.Context, _ interface{}) {
		*ctx = context.WithValue(*ctx, brokerclient.ContextAccessToken, iamToken)
	})

	details, err := brokerClient.GetCMPDetails(ctx)
	if err != nil {
		return "", "", fmt.Errorf("VMaaS broker exchange failed: %w", err)
	}

	return details.URL, details.AccessToken, nil
}

// iamToken returns the pre-generated GreenLake IAM token when one is
// configured, and otherwise exchanges the API client credentials for one.
func (c Config) iamToken(ctx context.Context) (string, error) {
	opts := []serviceclient.CreateOpt{
		serviceclient.WithIAMVersion(c.Version),
	}

	if c.IAMToken != "" {
		opts = append(opts, serviceclient.WithIAMToken(c.IAMToken))
	} else {
		opts = append(opts,
			serviceclient.WithIAMServiceURL(c.IssuerURL),
			serviceclient.WithClientID(c.ClientID),
			serviceclient.WithClientSecret(c.ClientSecret),
		)
	}

	handler, err := serviceclient.NewHandler(opts...)
	if err != nil {
		return "", fmt.Errorf("could not create GreenLake IAM token handler: %w", err)
	}

	token, err := retrieve.NewTokenRetrieveFunc(handler)(ctx)
	if err != nil {
		return "", fmt.Errorf("could not obtain GreenLake IAM token: %w", err)
	}

	return token, nil
}
