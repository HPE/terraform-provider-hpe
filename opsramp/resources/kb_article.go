// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/HPE/terraform-provider-hpe/opsramp/client"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure implementation satisfies the expected interfaces
var _ resource.Resource = &KBArticleResource{}
var _ resource.ResourceWithImportState = &KBArticleResource{}
var _ resource.ResourceWithModifyPlan = &KBArticleResource{}

// KBArticleResource defines the resource implementation.
type KBArticleResource struct {
	BaseResource
}

// KBArticleModel maps Terraform schema attributes to the provider model.
type KBArticleModel struct {
	Client           types.String `tfsdk:"client"`
	Id               types.String `tfsdk:"id"`
	Subject          types.String `tfsdk:"subject"`
	Content          types.String `tfsdk:"content"`
	CategoryId       types.String `tfsdk:"category_id"`
	State            types.String `tfsdk:"state"`
	LinkedArticleIds types.List   `tfsdk:"linked_article_ids"`
	AttachmentIds    types.List   `tfsdk:"attachment_ids"`
	ExpiryDate       types.String `tfsdk:"expiry_date"`
}

// NewKBArticle creates a new instance of the resource.
func NewKBArticle() resource.Resource {
	return &KBArticleResource{}
}

// Metadata returns the resource type name.
func (r *KBArticleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kb_article"
}

// Schema defines the schema for the resource.
func (r *KBArticleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an OpsRamp Knowledge Base Article.",
		Attributes: map[string]schema.Attribute{
			"client": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The client (tenant) UUID where this article should be created. If not specified, uses the provider's tenant.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The string ID of the KB article, assigned by the API.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"subject": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The subject of the KB article.",
			},
			"content": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The full content/body of the KB article (may contain HTML).",
			},
			"category_id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The string ID of the KB category this article belongs to.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"state": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The publication status of the article. (`PUBLISHED`, `TRASH`).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"linked_article_ids": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Description: "The IDs of linked articles",
			},
			"attachment_ids": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "The IDs of attachments",
			},
			"expiry_date": schema.StringAttribute{
				Optional:    true,
				Description: "The expiry date of the article in ISO 8601 format (e.g. `2025-12-31T23:59:59Z`)",
			},
		},
	}
}

func (r *KBArticleResource) resolveTenantId(clientAttr types.String) string {
	if !clientAttr.IsNull() && clientAttr.ValueString() != "" {
		return clientAttr.ValueString()
	}
	return r.apiClient.TenantId
}

func buildKBArticleRequest(plan KBArticleModel) client.KBArticle {
	article := client.KBArticle{
		Subject: plan.Subject.ValueString(),
		Content: plan.Content.ValueString(),
		State:   "PUBLISHED",
	}

	if !plan.ExpiryDate.IsNull() && !plan.ExpiryDate.IsUnknown() {
		article.ExpiryDate = plan.ExpiryDate.ValueString()
	}

	if !plan.CategoryId.IsNull() && !plan.CategoryId.IsUnknown() && plan.CategoryId.ValueString() != "" {
		article.Category = &client.KBCategoryRef{Id: plan.CategoryId.ValueString()}
	}

	if !plan.AttachmentIds.IsNull() && !plan.AttachmentIds.IsUnknown() {
		var attachmentIds []string
		diags := plan.AttachmentIds.ElementsAs(context.Background(), &attachmentIds, false)
		if !diags.HasError() {
			for _, id := range attachmentIds {
				article.Attachments = append(article.Attachments, client.KBArticleAttachment{Id: id})
			}
		}
	}

	if !plan.LinkedArticleIds.IsNull() && !plan.LinkedArticleIds.IsUnknown() {
		var linkedIds []string
		diags := plan.LinkedArticleIds.ElementsAs(context.Background(), &linkedIds, false)
		if !diags.HasError() {
			for _, id := range linkedIds {
				article.LinkedArticles = append(article.LinkedArticles, client.KBArticleRef{Id: id})
			}
		}
	}

	return article
}

func populateKBArticleModel(model *KBArticleModel, a *client.KBArticle) {
	model.Id = types.StringValue(a.Id)
	model.Subject = types.StringValue(a.Subject)
	model.Content = types.StringValue(a.Content)
	model.State = types.StringValue("PUBLISHED")
	if a.ExpiryDate != "" {
		model.ExpiryDate = types.StringValue(a.ExpiryDate)
	} else {
		model.ExpiryDate = types.StringNull()
	}

	if a.Category != nil {
		model.CategoryId = types.StringValue(a.Category.Id)
	} else {
		model.CategoryId = types.StringNull()
	}

	attachmentIds := make([]string, 0, len(a.Attachments))
	for _, attachment := range a.Attachments {
		attachmentIds = append(attachmentIds, attachment.Id)
	}
	if len(attachmentIds) > 0 {
		model.AttachmentIds = types.ListValueMust(types.StringType, stringSliceToAttrValues(attachmentIds))
	} else {
		model.AttachmentIds = types.ListNull(types.StringType)
	}

	linkedIds := make([]string, 0, len(a.LinkedArticles))
	for _, linkedArticle := range a.LinkedArticles {
		linkedIds = append(linkedIds, linkedArticle.Id)
	}
	if len(linkedIds) > 0 {
		model.LinkedArticleIds = types.ListValueMust(types.StringType, stringSliceToAttrValues(linkedIds))
	} else {
		model.LinkedArticleIds = types.ListNull(types.StringType)
	}
}

func stringSliceToAttrValues(values []string) []attr.Value {
	result := make([]attr.Value, 0, len(values))
	for _, v := range values {
		result = append(result, types.StringValue(v))
	}
	return result
}

// Create handles the creation of the resource.
func (r *KBArticleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan KBArticleModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.resolveTenantId(plan.Client)
	articleReq := buildKBArticleRequest(plan)

	created, err := r.apiClient.CreateKBArticle(tenantId, articleReq)
	if err != nil {
		resp.Diagnostics.AddError("Error Creating KB Article", fmt.Sprintf("Could not create KB article: %s", err))
		return
	}

	populateKBArticleModel(&plan, created)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Read handles reading the resource state.
func (r *KBArticleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state KBArticleModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.resolveTenantId(state.Client)
	articleId := state.Id.ValueString()

	existing, err := r.apiClient.GetKBArticle(tenantId, articleId)
	if err != nil {
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "not found") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error Reading KB Article", fmt.Sprintf("Could not read KB article %s: %s", articleId, err))
		return
	}

	populateKBArticleModel(&state, existing)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update handles updating the resource.
func (r *KBArticleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state KBArticleModel
	req.Plan.Get(ctx, &plan)
	req.State.Get(ctx, &state)

	tenantId := r.resolveTenantId(plan.Client)
	articleId := state.Id.ValueString()

	articleReq := buildKBArticleRequest(plan)

	updated, err := r.apiClient.UpdateKBArticle(tenantId, articleId, articleReq)
	if err != nil {
		resp.Diagnostics.AddError("Error Updating KB Article", fmt.Sprintf("Could not update KB article %s: %s", articleId, err))
		return
	}

	plan.Id = state.Id
	populateKBArticleModel(&plan, updated)

	diags := resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Delete handles deleting the resource.
func (r *KBArticleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state KBArticleModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.resolveTenantId(state.Client)
	articleId := state.Id.ValueString()

	err := r.apiClient.DeleteKBArticle(tenantId, articleId)
	if err != nil {
		resp.Diagnostics.AddError("Error Deleting KB Article", fmt.Sprintf("Could not delete KB article %s: %s", articleId, err))
		return
	}

	resp.State.RemoveResource(ctx)
}

// ImportState handles importing an existing resource.
// Import ID format: <article_id> or <client_id>:<article_id> (MSP only)
func (r *KBArticleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parsed, err := r.ParseImportID(req.ID, 1)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Import ID", err.Error())
		return
	}

	tenantId := r.TenantForImport(parsed)
	articleId := parsed.Parts[0]

	existing, err := r.apiClient.GetKBArticle(tenantId, articleId)
	if err != nil {
		resp.Diagnostics.AddError("Error Importing KB Article", fmt.Sprintf("Could not import KB article: %s", err))
		return
	}

	var state KBArticleModel
	state.Client = parsed.Client
	populateKBArticleModel(&state, existing)

	resp.State.Set(ctx, &state)
}
