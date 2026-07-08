// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package certificate

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	providererrors "github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const (
	summary                 = "read certificate data source"
	ErrorNoValidSearchTerms = `no valid search terms - an id or name is required`
	ErrorRunningPreApply    = `Error running pre-apply plan: exit status 1`
	ErrorNoCertificateFound = `no certificate found`
	ErrorMultipleCerts      = `multiple certificates were returned`
)

// Ensure the implementation satisfies the expected interfaces.
var _ datasource.DataSource = &DataSource{}

// NewDataSource is a helper function to simplify the provider implementation.
func NewDataSource() datasource.DataSource {
	return &DataSource{}
}

// DataSource is the data source implementation.
type DataSource struct {
	configure.DataSourceWithMorpheusConfigure
	datasource.DataSource
}

// Metadata returns the data source type name.
func (d *DataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_" + "certificate"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = CertificateDataSourceSchema(ctx)
}

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
	state.Category = convert.StrToType(cert.Category.Get())

	// cert_type (nullable)
	state.CertType = convert.StrToType(cert.CertType.Get())

	// common_name (nullable)
	state.CommonName = convert.StrToType(cert.CommonName.Get())

	// integration_id (nullable)
	state.IntegrationId = convert.Int64ToType(cert.IntegrationId.Get())

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
