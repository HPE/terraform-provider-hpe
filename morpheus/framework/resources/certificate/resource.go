package certificate

import (
	"context"
	"fmt"
	"strconv"

	sdk "github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/cleanup"
)

var (
	_ resource.Resource                = &certificateResource{}
	_ resource.ResourceWithConfigure   = &certificateResource{}
	_ resource.ResourceWithImportState = &certificateResource{}
)

type certificateResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &certificateResource{}
}

func (r *certificateResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_certificate"
}

func (r *certificateResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = CertificateSchema(ctx)
}

func (r *certificateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan certificateModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := sdk.AddCertificateRequestCertificate{
		Name:     plan.Name.ValueStringPointer(),
		CertFile: plan.CertFile.ValueStringPointer(),
		KeyFile:  plan.KeyFile.ValueStringPointer(),
	}
	if !plan.Description.IsNull() {
		body.Description = plan.Description.ValueStringPointer()
	}
	if !plan.DomainName.IsNull() {
		body.DomainName = plan.DomainName.ValueStringPointer()
	}

	result, httpResp, err := client.SSLCertificatesAPI.AddCertificate(ctx).AddCertificateRequest(sdk.AddCertificateRequest{
		Certificate: &body,
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "certificate", plan.Name.ValueString(), err, httpResp)

		return
	}

	createCert := result.GetCertificate()
	id := (&createCert).GetId()

	readResult, httpResp, err := client.SSLCertificatesAPI.GetCertificate(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "certificate", plan.Name.ValueString(), err, httpResp)
		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: "certificate",
			ResourceID:   id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})

		return
	}

	readCert := readResult.GetCertificate()
	mapGetResponseToModel(&plan, &readCert)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *certificateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state certificateModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	result, httpResp, err := client.SSLCertificatesAPI.GetCertificate(ctx, id).Execute()
	if errfmt.IsNotFound(httpResp) {
		resp.State.RemoveResource(ctx)

		return
	}
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "certificate", "", err, httpResp)

		return
	}

	cert := result.Certificate
	if cert == nil {
		resp.Diagnostics.AddError("API returned nil", "Certificate is nil in the response")

		return
	}
	mapGetResponseToModel(&state, cert)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *certificateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan certificateModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueInt64()

	body := sdk.UpdateCertificateRequestCertificate{
		Name:     plan.Name.ValueStringPointer(),
		CertFile: plan.CertFile.ValueStringPointer(),
		KeyFile:  plan.KeyFile.ValueStringPointer(),
	}
	if !plan.Description.IsNull() {
		body.Description = plan.Description.ValueStringPointer()
	}
	if !plan.DomainName.IsNull() {
		body.DomainName = plan.DomainName.ValueStringPointer()
	}

	_, httpResp, err := client.SSLCertificatesAPI.UpdateCertificate(ctx, id).
		UpdateCertificateRequest(sdk.UpdateCertificateRequest{
			Certificate: &body,
		}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "certificate", plan.Name.ValueString(), err, httpResp)

		return
	}

	readResult, httpResp, err := client.SSLCertificatesAPI.GetCertificate(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "certificate", plan.Name.ValueString(), err, httpResp)

		return
	}

	readCert := readResult.GetCertificate()
	mapGetResponseToModel(&plan, &readCert)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *certificateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state certificateModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	_, httpResp, err := client.SSLCertificatesAPI.DeleteCertificate(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "certificate", "", err, httpResp)

		return
	}
}

func (r *certificateResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse ID %q as integer: %s", req.ID, err))

		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func mapAddResponseToModel(model *certificateModel, cert *sdk.AddCertificate200ResponseAllOfCertificate) {
	if cert.Id != nil {
		model.ID = types.Int64Value(*cert.Id)
	}
	if cert.Name != nil {
		model.Name = types.StringValue(*cert.Name)
	}
	// Only update description if the API returns a value; otherwise keep the
	// plan value (no-op by not modifying model.Description).
	if cert.Description.IsSet() && cert.Description.Get() != nil {
		model.Description = types.StringValue(*cert.Description.Get())
	}
	if cert.DomainName.IsSet() && cert.DomainName.Get() != nil && *cert.DomainName.Get() != "" {
		model.DomainName = types.StringValue(*cert.DomainName.Get())
	} else {
		model.DomainName = types.StringNull()
	}
	if cert.Enabled != nil {
		model.Enabled = types.BoolValue(*cert.Enabled)
	}
	// cert_file and key_file are not returned by the API; keep plan values
}

func mapGetResponseToModel(model *certificateModel, cert *sdk.GetCertificate200ResponseCertificate) {
	if cert.Id != nil {
		model.ID = types.Int64Value(*cert.Id)
	}
	if cert.Name != nil {
		model.Name = types.StringValue(*cert.Name)
	}
	if cert.Description.IsSet() && cert.Description.Get() != nil {
		model.Description = types.StringValue(*cert.Description.Get())
	} else {
		model.Description = types.StringNull()
	}
	if cert.DomainName.IsSet() && cert.DomainName.Get() != nil && *cert.DomainName.Get() != "" {
		model.DomainName = types.StringValue(*cert.DomainName.Get())
	} else {
		model.DomainName = types.StringNull()
	}
	if cert.Enabled != nil {
		model.Enabled = types.BoolValue(*cert.Enabled)
	}
	// cert_file and key_file are not returned by the API; keep existing state values
}

func mapUpdateResponseToModel(model *certificateModel, cert *sdk.GetCertificate200ResponseCertificate) {
	if cert.Id != nil {
		model.ID = types.Int64Value(*cert.Id)
	}
	if cert.Name != nil {
		model.Name = types.StringValue(*cert.Name)
	}
	// Only update description if the API returns a value; otherwise keep the
	// plan value (no-op by not modifying model.Description).
	if cert.Description.IsSet() && cert.Description.Get() != nil {
		model.Description = types.StringValue(*cert.Description.Get())
	}
	if cert.DomainName.IsSet() && cert.DomainName.Get() != nil && *cert.DomainName.Get() != "" {
		model.DomainName = types.StringValue(*cert.DomainName.Get())
	} else {
		model.DomainName = types.StringNull()
	}
	if cert.Enabled != nil {
		model.Enabled = types.BoolValue(*cert.Enabled)
	}
	// cert_file and key_file are not returned by the API; keep plan values
}
