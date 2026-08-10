// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package pce implements the PCE Identity token exchange. It trades GreenLake
// API client credentials for the URL and access token of the Morpheus instance
// sitting behind the PCE broker.
//
// The same exchange serves both PCE deployment types. They differ only in the
// IAM dialect used to mint the token (GLCS for Connected, GLP for Disconnected)
// and in how the broker request is scoped; see Config.Version.
package pce

import (
	"context"
	"fmt"

	brokerclient "github.com/HPE/terraform-provider-hpe/morpheus/pce/sdk/broker/client"
	"github.com/HPE/terraform-provider-hpe/morpheus/pce/sdk/token/iamversion"
	"github.com/HPE/terraform-provider-hpe/morpheus/pce/sdk/token/serviceclient"
)

// DefaultBrokerURL is the HPE-hosted PCE broker, used when a Config does not
// name one.
const DefaultBrokerURL = "https://vmaas-broker.us1.greenlake-hpe.com"

// Config is the PCE Identity configuration needed to resolve Morpheus
// connection details. It deliberately holds plain strings rather than Terraform
// types so that both the framework and the SDKv2 Morpheus provider can use it.
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

// TokenExchange obtains a GreenLake IAM token and trades it with the PCE broker
// for the Morpheus URL and access token.
//
// Every call performs the exchange. The "hpe" provider muxes a framework
// provider and an SDKv2 provider and Terraform configures both, so a
// configuration using an identity block exchanges once per provider rather than
// once per graph walk.
//
// Whether a repeated exchange returns the same Morpheus token or a new one is
// the broker's decision, not something this package arranges. Either is safe:
// issuing a token does not invalidate one already issued, so the two providers
// cannot disturb each other's credentials.
//
// The returned access token is used as-is and is not refreshed, so callers are
// expected to finish their work within its validity period.
func TokenExchange(ctx context.Context, c Config) (string, string, error) {
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
		return "", "", fmt.Errorf("PCE broker exchange failed: %w", err)
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

	token, err := handler.Token(ctx)
	if err != nil {
		return "", fmt.Errorf("could not obtain GreenLake IAM token: %w", err)
	}

	return token, nil
}
