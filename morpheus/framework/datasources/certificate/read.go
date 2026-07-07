// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package certificate

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	providererrors "github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

func certificateAsState(
	cert *sdk.GetCertificate200ResponseCertificate,
) CertificateModel {
	state := CertificateModel{
		Id:          convert.Int64ToType(cert.Id),
		Name:        convert.StrToType(cert.Name),
		Description: convert.StrToType(cert.Description.Get()),
		DomainName:  convert.StrToType(cert.DomainName.Get()),
		Enabled:     convert.BoolToType(cert.Enabled),
		AccountId:   convert.Int64ToType(cert.AccountId),
		Generated:   convert.BoolToType(cert.Generated),
		SelfSigned:  convert.BoolToType(cert.SelfSigned),
		Wildcard:    convert.BoolToType(cert.Wildcard),
	}

	// category (nullable)
	if cert.Category.IsSet() {
		state.Category = convert.StrToType(cert.Category.Get())
	} else {
		state.Category = types.StringNull()
	}

	// cert_type (nullable)
	if cert.CertType.IsSet() {
		state.CertType = convert.StrToType(cert.CertType.Get())
	} else {
		state.CertType = types.StringNull()
	}

	// common_name (nullable)
	if cert.CommonName.IsSet() {
		state.CommonName = convert.StrToType(cert.CommonName.Get())
	} else {
		state.CommonName = types.StringNull()
	}

	// integration_id (nullable)
	if cert.IntegrationId.IsSet() {
		state.IntegrationId = convert.Int64ToType(cert.IntegrationId.Get())
	} else {
		state.IntegrationId = types.Int64Null()
	}

	// type nested object {id, code}
	state.Type = NewTypeValueNull()
	if cert.Type != nil {
		state.Type = TypeValue{
			Id:    convert.Int64ToType(cert.Type.Id),
			Code:  convert.StrToType(cert.Type.Code),
			state: attr.ValueStateKnown,
		}
	}

	return state
}

func getByID(
	ctx context.Context,
	id int64,
	apiClient *sdk.APIClient,
) (*sdk.GetCertificate200ResponseCertificate, error) {
	r, hresp, err := apiClient.SSLCertificatesAPI.GetCertificate(ctx, id).Execute()
	if r == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET failed for certificate %d: %s", id, providererrors.ErrMsg(err, hresp))
	}

	if r.Certificate == nil {
		return nil, fmt.Errorf("GET failed for certificate %d: response missing certificate", id)
	}

	cert := *r.Certificate

	return &cert, nil
}

func getByName(
	ctx context.Context,
	name string,
	apiClient *sdk.APIClient,
) (*sdk.GetCertificate200ResponseCertificate, error) {
	rs, hresp, err := apiClient.SSLCertificatesAPI.ListCertificates(ctx).Name(name).Execute()
	if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET failed for certificate %s: %s", name, providererrors.ErrMsg(err, hresp))
	}

	var matched []sdk.ListCertificates200ResponseCertificatesInner

	for _, o := range rs.Certificates {
		if o.Name != nil && *o.Name == name {
			matched = append(matched, o)
		}
	}

	if len(matched) == 0 {
		return nil, errors.New(ErrorNoCertificateFound)
	} else if len(matched) > 1 {
		return nil, errors.New(ErrorMultipleCerts)
	}

	if matched[0].Id == nil {
		return nil, fmt.Errorf("GET failed for certificate %s: response missing id", name)
	}

	return getByID(ctx, *matched[0].Id, apiClient)
}

func getCertificate(
	ctx context.Context,
	config *CertificateModel,
	apiClient *sdk.APIClient,
) (*sdk.GetCertificate200ResponseCertificate, error) {
	if !config.Id.IsNull() {
		return getByID(ctx, config.Id.ValueInt64(), apiClient)
	} else if !config.Name.IsNull() {
		return getByName(ctx, config.Name.ValueString(), apiClient)
	}

	return nil, errors.New(ErrorNoValidSearchTerms)
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config CertificateModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiClient, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			summary,
			"could not create sdk client",
		)

		return
	}

	cert, err := getCertificate(ctx, &config, apiClient)
	if err != nil {
		resp.Diagnostics.AddError(
			summary,
			err.Error(),
		)

		return
	}

	state := certificateAsState(cert)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
