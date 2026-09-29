package storageserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/wodigest"
	"github.com/HPE/terraform-provider-hpe/utils/cleanup"
	"github.com/HPE/terraform-provider-hpe/utils/schemadefaults"
)

const (
	statusOk    = "ok"
	statusError = "error"
)

var (
	targetStatuses = []string{statusOk}
	errorStatuses  = []string{statusError}
)

func checkStatusDone(status string) error {
	switch {
	case slices.Contains(errorStatuses, status):
		return backoff.Permanent(errors.New("reached error status: " + status))
	case slices.Contains(targetStatuses, status):
		return nil
	default:
		return backoff.RetryAfter(5)
	}
}

var (
	_ resource.Resource                = &storageServerResource{}
	_ resource.ResourceWithConfigure   = &storageServerResource{}
	_ resource.ResourceWithImportState = &storageServerResource{}
	_ resource.ResourceWithModifyPlan  = &storageServerResource{}
)

type storageServerResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &storageServerResource{}
}

func (r *storageServerResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_" + "storage_server"
}

func (r *storageServerResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = StorageServerResourceSchema(ctx)
}

func (r *storageServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan StorageServerModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get timeout from HCL if set, the default is 5 minutes
	createTimeout, diags := plan.Timeouts.Create(ctx, 5*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, createTimeout)
	defer cancel()

	body := sdk.AddStorageServersRequestStorageServer{
		Name:   plan.Name.ValueString(),
		Type:   plan.Type.ValueString(),
		Config: map[string]interface{}{},
	}
	// Note: servicePort defaults to 22 in the API and is not exposed as a
	// configurable attribute because the GET response does not return it.
	if !plan.Description.IsNull() {
		body.Description = plan.Description.ValueStringPointer()
	}
	if !plan.Enabled.IsNull() {
		body.Enabled = plan.Enabled.ValueBoolPointer()
	}
	if !plan.Visibility.IsNull() {
		body.Visibility = plan.Visibility.ValueStringPointer()
	}
	if !plan.ServiceHost.IsNull() {
		body.ServiceHost = plan.ServiceHost.ValueStringPointer()
	}
	if !plan.ServiceUrl.IsNull() {
		body.ServiceUrl = plan.ServiceUrl.ValueStringPointer()
	}
	if !plan.ServiceUsername.IsNull() {
		body.ServiceUsername = plan.ServiceUsername.ValueStringPointer()
	}
	// service_password_wo is write-only: its value is only in the configuration.
	var servicePasswordWo types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("service_password_wo"), &servicePasswordWo)...)
	if !servicePasswordWo.IsNull() && !servicePasswordWo.IsUnknown() {
		body.ServicePassword = servicePasswordWo.ValueStringPointer()
	}

	// Credential: use credential_id for stored credentials
	if !plan.CredentialId.IsNull() {
		idStr := strconv.FormatInt(plan.CredentialId.ValueInt64(), 10)
		body.Credential = &sdk.AddStorageServersRequestStorageServerCredential{
			Type: &idStr,
		}
	}

	// Tenants
	if !plan.Tenants.IsNull() && !plan.Tenants.IsUnknown() {
		var tenantIds []int64
		resp.Diagnostics.Append(plan.Tenants.ElementsAs(ctx, &tenantIds, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		tenants := make([]sdk.AddStorageServersRequestStorageServerTenantsInner, 0, len(tenantIds))
		for _, id := range tenantIds {
			tenants = append(tenants, sdk.AddStorageServersRequestStorageServerTenantsInner{
				Id: &id,
			})
		}
		body.Tenants = tenants
	}

	result, httpResp, err := client.StorageAPI.AddStorageServers(ctx).
		AddStorageServersRequest(sdk.AddStorageServersRequest{
			StorageServer: body,
		}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "storage_server", plan.Name.ValueString(), err, httpResp)

		return
	}

	ss := result.StorageServer
	if ss == nil {
		resp.Diagnostics.AddError("API returned nil", "StorageServer is nil in the response")

		return
	}

	if ss.Id == nil {
		resp.Diagnostics.AddError("Create Error", "Storage server ID not returned")

		return
	}
	id := *ss.Id
	plan.Id = types.Int64Value(id)

	// Helper to taint the resource state on an error after the POST request
	taintResourceState := func(id int64) {
		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: "storage_server",
			ResourceID:   id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})
	}

	// Poll until the storage server initialization completes
	waitForReady := func() (*sdk.GetStorageServers200Response, error) {
		response, hresp, err := client.StorageAPI.GetStorageServers(ctx, id).Execute()
		if err != nil {
			// The SDK may return a non-nil error alongside a 200 response for
			// non-fatal issues (e.g., response body warnings). Only treat as
			// permanent failure if the status code is not 200.
			if hresp == nil || hresp.StatusCode != http.StatusOK {
				return nil, backoff.Permanent(err)
			}
		}

		ss := response.StorageServer
		if ss == nil {
			return response, backoff.RetryAfter(5)
		}
		if status := ss.Status.Get(); status != nil {
			return response, checkStatusDone(*status)
		}

		return response, backoff.RetryAfter(5)
	}

	if r, err := backoff.Retry(
		ctx,
		waitForReady,
		backoff.WithBackOff(backoff.NewConstantBackOff(5*time.Second)),
		backoff.WithMaxElapsedTime(createTimeout),
	); err != nil {
		var status string
		if r != nil && r.StorageServer != nil {
			ss := r.StorageServer
			if s := ss.Status.Get(); s != nil {
				status = *s
			}
		}

		resp.Diagnostics.AddError(
			"storage server initialization failed",
			fmt.Sprintf("Storage server %d failed to reach ok status. Current status: %s", id, status),
		)

		taintResourceState(id)

		return
	} else {
		storageServer := r.StorageServer
		if storageServer == nil {
			resp.Diagnostics.AddError("API returned nil", "StorageServer is nil in the response")

			return
		}

		mapGetResponseToModel(ctx, &plan, storageServer, &resp.Diagnostics)
	}

	plan.WriteOnlyRevision = wodigest.NewRevision()
	resp.Diagnostics.Append(wo.Store(ctx, resp.Private, req.Config)...)

	resp.Diagnostics.Append(clearable.Record(ctx, req.Config, resp.Private, clearableAttributes)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *storageServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state StorageServerModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueInt64()

	result, httpResp, err := client.StorageAPI.GetStorageServers(ctx, id).Execute()
	if errfmt.IsNotFound(httpResp) {
		resp.State.RemoveResource(ctx)

		return
	}
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "storage_server", "", err, httpResp)

		return
	}

	ss := result.StorageServer
	if ss == nil {
		resp.Diagnostics.AddError("API returned nil", "StorageServer is nil in the response")

		return
	}

	mapGetResponseToModel(ctx, &state, ss, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

	// Fill any schema-declared default the API omitted (null in state) so an
	// imported resource does not plan a change nobody made. MORPH-16192.
	resp.Diagnostics.Append(
		schemadefaults.Apply(ctx, StorageServerResourceSchema(ctx), &resp.State)...,
	)
}

func (r *storageServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan StorageServerModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get timeout from HCL if set, the default is 5 minutes
	updateTimeout, diags := plan.Timeouts.Update(ctx, 5*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, updateTimeout)
	defer cancel()

	id := plan.Id.ValueInt64()

	// Attributes removed from the configuration since the last apply, to be
	// cleared. Read before Record, below, overwrites what it reads.
	pending, pdiags := clearable.Compute(ctx, req.Config, req.State, req.Private, clearableAttributes)
	resp.Diagnostics.Append(pdiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := sdk.UpdateStorageServersRequestStorageServer{
		Name: plan.Name.ValueStringPointer(),
		Type: plan.Type.ValueStringPointer(),
	}
	clearable.NullableString(pending, "description", plan.Description, &body.Description)
	if !plan.Enabled.IsNull() {
		body.Enabled = plan.Enabled.ValueBoolPointer()
	}
	if !plan.Visibility.IsNull() {
		body.Visibility = plan.Visibility.ValueStringPointer()
	}
	clearable.NullableString(pending, "service_host", plan.ServiceHost, &body.ServiceHost)
	clearable.NullableString(pending, "service_url", plan.ServiceUrl, &body.ServiceUrl)
	clearable.NullableString(pending, "service_username", plan.ServiceUsername, &body.ServiceUsername)
	// service_password_wo is write-only: its value is only in the configuration.
	// Re-sent when ModifyPlan marked write_only_revision unknown (its digest
	// changed, or the deprecated marker did). See write_only.go.
	var servicePasswordWo types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("service_password_wo"), &servicePasswordWo)...)
	if wodigest.Resend(ctx, req.Plan) && !servicePasswordWo.IsNull() && !servicePasswordWo.IsUnknown() {
		body.ServicePassword = servicePasswordWo.ValueStringPointer()
	}

	// Credential
	if !plan.CredentialId.IsNull() {
		idStr := strconv.FormatInt(plan.CredentialId.ValueInt64(), 10)
		body.Credential = &sdk.UpdateStorageServersRequestStorageServerCredential{
			Type: &idStr,
		}
	}

	// Tenants: the planned list, or an empty list when being cleared
	if tenantIds, ok, ldiags := pending.Int64s(ctx, "tenants", plan.Tenants); ok {
		resp.Diagnostics.Append(ldiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		tenants := make([]sdk.UpdateStorageServersRequestStorageServerTenantsInner, 0, len(tenantIds))
		for _, tid := range tenantIds {
			tenants = append(tenants, sdk.UpdateStorageServersRequestStorageServerTenantsInner{
				Id: &tid,
			})
		}
		body.Tenants = tenants
	}

	_, httpResp, err := client.StorageAPI.UpdateStorageServers(ctx, id).
		UpdateStorageServersRequest(sdk.UpdateStorageServersRequest{
			StorageServer: body,
		}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "storage_server", plan.Name.ValueString(), err, httpResp)

		return
	}

	// Poll until the storage server re-initialization completes
	waitForReady := func() (*sdk.GetStorageServers200Response, error) {
		response, hresp, err := client.StorageAPI.GetStorageServers(ctx, id).Execute()
		if err != nil {
			// The SDK may return a non-nil error alongside a 200 response for
			// non-fatal issues (e.g., response body warnings). Only treat as
			// permanent failure if the status code is not 200.
			if hresp == nil || hresp.StatusCode != http.StatusOK {
				return nil, backoff.Permanent(err)
			}
		}

		ss := response.StorageServer
		if ss == nil {
			return response, backoff.RetryAfter(5)
		}
		if status := ss.Status.Get(); status != nil {
			return response, checkStatusDone(*status)
		}

		return response, backoff.RetryAfter(5)
	}

	if r, err := backoff.Retry(
		ctx,
		waitForReady,
		backoff.WithBackOff(backoff.NewConstantBackOff(5*time.Second)),
		backoff.WithMaxElapsedTime(updateTimeout),
	); err != nil {
		var status string
		if r != nil && r.StorageServer != nil {
			ss := r.StorageServer
			if s := ss.Status.Get(); s != nil {
				status = *s
			}
		}

		resp.Diagnostics.AddError(
			"storage server initialization failed",
			fmt.Sprintf("Storage server %d failed to reach ok status after update. Current status: %s", id, status),
		)

		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: "storage_server",
			ResourceID:   id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})

		return
	} else {
		storageServer := r.StorageServer
		if storageServer == nil {
			resp.Diagnostics.AddError("API returned nil", "StorageServer is nil in the response")

			return
		}

		mapGetResponseToModel(ctx, &plan, storageServer, &resp.Diagnostics)
	}

	// write_only_revision is not in the API. Keep the prior value unless a
	// write-only value was re-sent (ModifyPlan marked it unknown); and always
	// record the digests of the configured values. See write_only.go.
	var priorState StorageServerModel
	resp.Diagnostics.Append(req.State.Get(ctx, &priorState)...)
	plan.WriteOnlyRevision = priorState.WriteOnlyRevision
	if wodigest.Resend(ctx, req.Plan) {
		plan.WriteOnlyRevision = wodigest.NewRevision()
	}
	resp.Diagnostics.Append(wo.Store(ctx, resp.Private, req.Config)...)

	resp.Diagnostics.Append(clearable.Record(ctx, req.Config, resp.Private, clearableAttributes)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *storageServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state StorageServerModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get timeout from HCL if set, the default is 5 minutes
	deleteTimeout, diags := state.Timeouts.Delete(ctx, 5*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, deleteTimeout)
	defer cancel()

	id := state.Id.ValueInt64()

	_, httpResp, err := client.StorageAPI.RemoveStorageServers(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "storage_server", "", err, httpResp)

		return
	}

	// Wait for the storage server to be deleted (404 response)
	waitForDeleted := func() (int, error) {
		_, hresp, err := client.StorageAPI.GetStorageServers(ctx, id).Execute()
		if err != nil {
			if hresp == nil || hresp.StatusCode != http.StatusNotFound {
				return 0, backoff.Permanent(err)
			}
		}

		switch hresp.StatusCode {
		case http.StatusNotFound:
			return hresp.StatusCode, nil
		default:
			return hresp.StatusCode, backoff.RetryAfter(5)
		}
	}

	if statusCode, err := backoff.Retry(
		ctx,
		waitForDeleted,
		backoff.WithBackOff(backoff.NewConstantBackOff(5*time.Second)),
		backoff.WithMaxElapsedTime(deleteTimeout),
	); err != nil {
		resp.Diagnostics.AddError(
			"storage server deletion failed",
			fmt.Sprintf("Storage server %d: DELETE failed, current status code: %d", id, statusCode),
		)
	}
}

func (r *storageServerResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	// "<id>" or "<id>#service_password_wo=<sha256 hex>": see wodigest.
	rawID, ok := wo.Import(ctx, req, resp)
	if !ok {
		return
	}

	id, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse ID %q as integer: %s", rawID, err))

		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func mapGetResponseToModel(
	ctx context.Context, model *StorageServerModel,
	ss *sdk.GetStorageServers200ResponseStorageServer, diags *diag.Diagnostics,
) {
	if ss.Id != nil {
		model.Id = types.Int64Value(*ss.Id)
	}
	if ss.Name != nil {
		model.Name = types.StringValue(*ss.Name)
	}
	if ss.Enabled != nil {
		model.Enabled = types.BoolValue(*ss.Enabled)
	}
	if v := ss.Description.Get(); v != nil {
		model.Description = types.StringValue(*v)
	} else {
		model.Description = types.StringNull()
	}
	if ss.Visibility != nil {
		model.Visibility = types.StringValue(*ss.Visibility)
	}
	if v := ss.ServiceUrl.Get(); v != nil {
		model.ServiceUrl = types.StringValue(*v)
	} else {
		model.ServiceUrl = types.StringNull()
	}
	if v := ss.ServiceHost.Get(); v != nil {
		model.ServiceHost = types.StringValue(*v)
	} else {
		model.ServiceHost = types.StringNull()
	}
	if v := ss.ServiceUsername.Get(); v != nil {
		model.ServiceUsername = types.StringValue(*v)
	}

	// Credential: extract ID from response if it's a stored credential
	if ss.Credential != nil && ss.Credential.Type != nil && *ss.Credential.Type != "local" {
		if ss.Credential.Id != nil {
			model.CredentialId = types.Int64Value(*ss.Credential.Id)
		}
	}

	// Tenants
	if ss.Tenants != nil {
		tenantIds := make([]int64, 0, len(ss.Tenants))
		for _, t := range ss.Tenants {
			if t.Id != nil {
				tenantIds = append(tenantIds, *t.Id)
			}
		}
		tenantList, d := types.ListValueFrom(ctx, types.Int64Type, tenantIds)
		diags.Append(d...)
		model.Tenants = tenantList
	} else {
		model.Tenants = types.ListNull(types.Int64Type)
	}

	// Write-only field (service_password_wo): preserve state value
}
