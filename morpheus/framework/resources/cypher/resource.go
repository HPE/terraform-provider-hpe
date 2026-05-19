package cypher

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &cypherResource{}
	_ resource.ResourceWithConfigure   = &cypherResource{}
	_ resource.ResourceWithImportState = &cypherResource{}
)

type cypherResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &cypherResource{}
}

func (r *cypherResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_cypher"
}

func (r *cypherResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = CypherSchema(ctx)
}

// cypherCreateResponse represents the JSON response from cypher write
type cypherCreateResponse struct {
	Success       bool   `json:"success"`
	LeaseDuration *int64 `json:"lease_duration"`
}

// cypherReadResponse represents the JSON response from cypher read
type cypherReadResponse struct {
	Success       bool   `json:"success"`
	LeaseDuration *int64 `json:"lease_duration"`
	Data          string `json:"data"`
}

func (r *cypherResource) doCypherRequest(ctx context.Context, method, cypherPath string, value string, ttl int64) (*http.Response, []byte, error) {
	client, err := r.NewClient(ctx)
	if err != nil {
		return nil, nil, err
	}

	cfg := client.GetConfig()
	baseURL, err := cfg.ServerURLWithContext(ctx, "")
	if err != nil {
		return nil, nil, err
	}

	// Build URL without path-escaping the cypherPath (it contains slashes)
	reqURL := fmt.Sprintf("%s/api/cypher/%s?type=string", baseURL, cypherPath)
	if ttl > 0 {
		reqURL += fmt.Sprintf("&ttl=%d", ttl)
	}

	var body io.Reader
	if method == http.MethodPost || method == http.MethodPut {
		jsonBody, _ := json.Marshal(map[string]string{"value": value})
		body = strings.NewReader(string(jsonBody))
	}

	req, err := http.NewRequestWithContext(ctx, method, reqURL, body)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Add auth headers from the SDK config
	if cfg.DefaultHeader != nil {
		for k, v := range cfg.DefaultHeader {
			req.Header.Set(k, v)
		}
	}

	httpResp, err := cfg.HTTPClient.Do(req)
	if err != nil {
		return httpResp, nil, err
	}
	defer httpResp.Body.Close()

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return httpResp, nil, err
	}

	return httpResp, respBody, nil
}

func (r *cypherResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan cypherModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cypherPath := plan.ID.ValueString()
	value := plan.Value.ValueString()
	var ttl int64
	if !plan.TTL.IsNull() {
		ttl = plan.TTL.ValueInt64()
	}

	httpResp, body, err := r.doCypherRequest(ctx, http.MethodPost, cypherPath, value, ttl)
	if err != nil {
		resp.Diagnostics.AddError("Error creating cypher", err.Error())
		return
	}
	if httpResp.StatusCode >= 400 {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error creating cypher: HTTP %d", httpResp.StatusCode),
			string(body),
		)
		return
	}

	var result cypherCreateResponse
	if err := json.Unmarshal(body, &result); err != nil {
		resp.Diagnostics.AddError("Error parsing cypher response", err.Error())
		return
	}

	plan.ID = types.StringValue(cypherPath)
	if result.LeaseDuration != nil {
		plan.LeaseDuration = types.Int64Value(*result.LeaseDuration)
	} else {
		plan.LeaseDuration = types.Int64Value(0)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *cypherResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state cypherModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cypherPath := state.ID.ValueString()

	httpResp, body, err := r.doCypherRequest(ctx, http.MethodGet, cypherPath, "", 0)
	if err != nil {
		resp.Diagnostics.AddError("Error reading cypher", err.Error())
		return
	}
	if errfmt.IsNotFound(httpResp) {
		resp.State.RemoveResource(ctx)
		return
	}
	if httpResp.StatusCode >= 400 {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading cypher: HTTP %d", httpResp.StatusCode),
			string(body),
		)
		return
	}

	var result cypherReadResponse
	if err := json.Unmarshal(body, &result); err != nil {
		resp.Diagnostics.AddError("Error parsing cypher response", err.Error())
		return
	}

	if result.LeaseDuration != nil {
		state.LeaseDuration = types.Int64Value(*result.LeaseDuration)
	} else {
		state.LeaseDuration = types.Int64Value(0)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *cypherResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan cypherModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cypherPath := plan.ID.ValueString()

	// Delete first
	httpResp, body, err := r.doCypherRequest(ctx, http.MethodDelete, cypherPath, "", 0)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting cypher for update", err.Error())
		return
	}
	if httpResp.StatusCode >= 400 && !errfmt.IsNotFound(httpResp) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error deleting cypher for update: HTTP %d", httpResp.StatusCode),
			string(body),
		)
		return
	}

	// Recreate
	value := plan.Value.ValueString()
	var ttl int64
	if !plan.TTL.IsNull() {
		ttl = plan.TTL.ValueInt64()
	}

	httpResp, body, err = r.doCypherRequest(ctx, http.MethodPost, cypherPath, value, ttl)
	if err != nil {
		resp.Diagnostics.AddError("Error creating cypher", err.Error())
		return
	}
	if httpResp.StatusCode >= 400 {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error creating cypher: HTTP %d", httpResp.StatusCode),
			string(body),
		)
		return
	}

	var result cypherCreateResponse
	if err := json.Unmarshal(body, &result); err != nil {
		resp.Diagnostics.AddError("Error parsing cypher response", err.Error())
		return
	}

	if result.LeaseDuration != nil {
		plan.LeaseDuration = types.Int64Value(*result.LeaseDuration)
	} else {
		plan.LeaseDuration = types.Int64Value(0)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *cypherResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state cypherModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cypherPath := state.ID.ValueString()

	httpResp, body, err := r.doCypherRequest(ctx, http.MethodDelete, cypherPath, "", 0)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting cypher", err.Error())
		return
	}
	if httpResp.StatusCode >= 400 && !errfmt.IsNotFound(httpResp) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error deleting cypher: HTTP %d", httpResp.StatusCode),
			string(body),
		)
		return
	}
}

func (r *cypherResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import by cypher key path
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}
