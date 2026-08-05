// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package morpheus

import (
	"context"
	"fmt"

	vmaascmpclient "github.com/HPE/terraform-provider-hpe/morpheus/greenlake/connected/sdk/vmaascmp/client"
	"github.com/HPE/terraform-provider-hpe/morpheus/greenlake/sdk/token/iamversion"
	"github.com/HPE/terraform-provider-hpe/morpheus/greenlake/sdk/token/retrieve"
	"github.com/HPE/terraform-provider-hpe/morpheus/greenlake/sdk/token/serviceclient"
	"github.com/HPE/terraform-provider-hpe/morpheus/model"
)

// defaultBrokerURL is the VMaaS broker used when a greenlake_connected block
// does not specify one.
const defaultBrokerURL = "https://vmaas-broker.us1.greenlake-hpe.com"

// greenlakeConnectedTokenExchange trades GreenLake credentials for the Morpheus
// connection details of a GreenLake Connected deployment. It first obtains a
// GreenLake IAM token, then calls the VMaaS broker to exchange it for a
// Morpheus url and access token.
//
// The returned access token is used as-is for the lifetime of the provider;
// it is not refreshed, as the token round tripper has no way to renew a bare
// access token. Callers are expected to complete their work within the token's
// validity period.
func greenlakeConnectedTokenExchange(
	ctx context.Context,
	m *model.GreenLakeConnectedModel,
) (string, string, error) {
	iamToken, err := greenlakeIAMToken(ctx, m)
	if err != nil {
		return "", "", err
	}

	brokerURL := defaultBrokerURL
	if !m.BrokerURL.IsNull() && m.BrokerURL.ValueString() != "" {
		brokerURL = m.BrokerURL.ValueString()
	}

	brokerCfg := vmaascmpclient.NewConfiguration()
	brokerCfg.Host = brokerURL

	if !m.Location.IsNull() {
		brokerCfg.DefaultQueryParams["location"] = m.Location.ValueString()
	}

	if !m.Space.IsNull() {
		brokerCfg.DefaultQueryParams["space"] = m.Space.ValueString()
	}

	brokerClient := vmaascmpclient.NewAPIClient(brokerCfg)

	// Inject the GreenLake IAM token on every request's context, which
	// prepareRequest reads for Bearer auth.
	brokerClient.SetMetaFnAndVersion(nil, 0, func(ctx *context.Context, _ interface{}) {
		*ctx = context.WithValue(*ctx, vmaascmpclient.ContextAccessToken, iamToken)
	})

	details, err := brokerClient.GetCMPDetails(ctx)
	if err != nil {
		return "", "", fmt.Errorf("VMaaS broker exchange failed: %w", err)
	}

	return details.URL, details.AccessToken, nil
}

// greenlakeIAMToken obtains a GreenLake IAM token, either by returning a
// pre-generated one or by exchanging the configured API client credentials.
func greenlakeIAMToken(
	ctx context.Context,
	m *model.GreenLakeConnectedModel,
) (string, error) {
	opts := []serviceclient.CreateOpt{
		serviceclient.WithIAMVersion(iamversion.GLCS),
	}

	if !m.IAMToken.IsNull() && m.IAMToken.ValueString() != "" {
		opts = append(opts, serviceclient.WithIAMToken(m.IAMToken.ValueString()))
	} else {
		opts = append(opts,
			serviceclient.WithIAMServiceURL(m.IssuerURL.ValueString()),
			serviceclient.WithClientID(m.ClientID.ValueString()),
			serviceclient.WithClientSecret(m.ClientSecret.ValueString()),
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
