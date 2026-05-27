package cypher

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/utils/modifiers"
)

type cypherModel struct {
	ID            types.String `tfsdk:"id"`
	Value         types.String `tfsdk:"value"`
	ValueVersion  types.Int64  `tfsdk:"value_version"`
	TTL           types.Int64  `tfsdk:"ttl"`
	LeaseDuration types.Int64  `tfsdk:"lease_duration"`
}

func CypherSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Manages a Morpheus Cypher secret resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "The key path for the cypher secret (e.g. secret/mykey).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"value": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				WriteOnly:   true,
				Description: "The secret value to store.",
				PlanModifiers: []planmodifier.String{
					modifiers.NullableStringUpdateModifier{},
				},
			},
			"value_version": schema.Int64Attribute{
				Optional:    true,
				Description: "Value version. Used to determine if value has been updated.",
			},
			"ttl": schema.Int64Attribute{
				Optional:    true,
				Description: "Time to live in seconds. 0 means no expiry.",
			},
			"lease_duration": schema.Int64Attribute{
				Computed:    true,
				Description: "The lease duration in seconds.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}
