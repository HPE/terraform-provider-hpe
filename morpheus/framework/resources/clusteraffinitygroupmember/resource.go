// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package clusteraffinitygroupmember manages a single server's membership of a
// cluster affinity group.
//
// It exists because the affinity group resource's own membership handling is
// authoritative: it owns the whole set, so anything added by other means -- an
// instance provisioned with config_hvm.affinity_group_id, a node added by
// hpe_morpheus_instance_node -- is treated as drift and removed on the next
// apply. This resource manages one membership at a time and ignores the rest,
// so it can coexist with membership created elsewhere.
package clusteraffinitygroupmember

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
)

// Interface compliance assertions.
var (
	_ resource.Resource                = &Resource{}
	_ resource.ResourceWithImportState = &Resource{}
)

// memberModel is the Terraform state model.
//
// The group's full server list is deliberately absent. If every membership held
// a copy of the whole list it would go stale as soon as any other member joined
// or left, producing a permanent diff on every membership at once -- which is
// the authoritative-set problem this resource exists to avoid, reintroduced at
// N times the scale. State records only which server belongs to which group;
// Read answers that by checking for presence.
type memberModel struct {
	ID              types.String `tfsdk:"id"`
	ClusterID       types.Int64  `tfsdk:"cluster_id"`
	AffinityGroupID types.Int64  `tfsdk:"affinity_group_id"`
	ServerID        types.Int64  `tfsdk:"server_id"`
}

// Resource implements hpe_morpheus_cluster_affinity_group_member.
type Resource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &Resource{}
}

// Metadata implements resource.Resource.
func (r *Resource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_cluster_affinity_group_member"
}

// Schema implements resource.Resource.
//
// Every attribute forces replacement. A membership is a fact rather than a
// thing with settings: it either exists or it does not, so there is nothing to
// update in place. Pointing it at a different group or a different server is a
// different membership.
func (r *Resource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Places a single compute server in a cluster affinity group, " +
			"without taking ownership of the group's other members.",
		MarkdownDescription: "Places a single compute server in a cluster affinity group, " +
			"without taking ownership of the group's other members.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				Description: "Identifier for this membership, in the form " +
					"cluster_id/affinity_group_id/server_id.",
				MarkdownDescription: "Identifier for this membership, in the form " +
					"`cluster_id/affinity_group_id/server_id`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"cluster_id": schema.Int64Attribute{
				Required:            true,
				Description:         "The cluster the affinity group belongs to.",
				MarkdownDescription: "The cluster the affinity group belongs to.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"affinity_group_id": schema.Int64Attribute{
				Required:            true,
				Description:         "The affinity group to place the server in.",
				MarkdownDescription: "The affinity group to place the server in.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"server_id": schema.Int64Attribute{
				Required: true,
				Description: "The compute server to place in the group. Affinity group " +
					"members are guests, not the hypervisor hosts they run on.",
				MarkdownDescription: "The compute server to place in the group. Affinity group " +
					"members are guests, not the hypervisor hosts they run on.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
		},
	}
}

// membershipID builds the composite identifier used for state and import.
func membershipID(clusterID, groupID, serverID int64) string {
	return fmt.Sprintf("%d/%d/%d", clusterID, groupID, serverID)
}

// containsServer reports whether the server is in the membership list.
func containsServer(servers []int64, serverID int64) bool {
	for _, s := range servers {
		if s == serverID {
			return true
		}
	}

	return false
}

// withoutServer returns the list with the server removed, preserving order.
func withoutServer(servers []int64, serverID int64) []int64 {
	out := make([]int64, 0, len(servers))

	for _, s := range servers {
		if s != serverID {
			out = append(out, s)
		}
	}

	return out
}

// importIDPath is the attribute reported when an import identifier is malformed.
var importIDPath = path.Root("id")
