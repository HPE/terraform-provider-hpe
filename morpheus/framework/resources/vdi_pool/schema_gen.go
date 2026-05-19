package vdi_pool

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type vdiPoolModel struct {
	ID                               types.Int64 `tfsdk:"id"`
	Name                             types.String `tfsdk:"name"`
	Description                      types.String `tfsdk:"description"`
	Enabled                          types.Bool   `tfsdk:"enabled"`
	AutoCreateLocalUserOnReservation types.Bool   `tfsdk:"auto_create_local_user_on_reservation"`
	PersistentUser                   types.Bool   `tfsdk:"persistent_user"`
	Recyclable                       types.Bool   `tfsdk:"recyclable"`
	AllowHypervisorConsole           types.Bool   `tfsdk:"allow_hypervisor_console"`
	AllowCopy                        types.Bool   `tfsdk:"allow_copy"`
	AllowPrinter                     types.Bool   `tfsdk:"allow_printer"`
	AllowFileshare                   types.Bool   `tfsdk:"allow_fileshare"`
	IdleTimeout                      types.Int64  `tfsdk:"idle_timeout"`
	MaxSessionTimeout                types.Int64  `tfsdk:"max_session_timeout"`
	MaxPoolSize                      types.Int64  `tfsdk:"max_pool_size"`
	MinIdle                          types.Int64  `tfsdk:"min_idle"`
	InitialPoolSize                  types.Int64  `tfsdk:"initial_pool_size"`
	MaxIdle                          types.Int64  `tfsdk:"max_idle"`
}

func VdiPoolSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Manages a Morpheus VDI Pool resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "The ID of the VDI pool.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the VDI pool.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "The description of the VDI pool.",
			},
			"enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "Whether the VDI pool is enabled.",
			},
			"auto_create_local_user_on_reservation": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether to auto-create a local user on reservation.",
			},
			"persistent_user": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether users are persistent across sessions.",
			},
			"recyclable": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether VDI pool instances are recyclable.",
			},
			"allow_hypervisor_console": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether to allow hypervisor console access.",
			},
			"allow_copy": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether to allow copy operations.",
			},
			"allow_printer": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether to allow printer access.",
			},
			"allow_fileshare": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether to allow file sharing.",
			},
			"idle_timeout": schema.Int64Attribute{
				Optional:    true,
				Description: "The idle timeout in minutes.",
			},
			"max_session_timeout": schema.Int64Attribute{
				Optional:    true,
				Description: "The maximum session timeout in minutes.",
			},
			"max_pool_size": schema.Int64Attribute{
				Required:    true,
				Description: "Max limit on number of allocations and instances within the pool.",
			},
			"min_idle": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Sets the minimum number of idle instances on standby in the pool.",
			},
			"initial_pool_size": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The initial size of the pool to be allocated on creation.",
			},
			"max_idle": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Sets the maximum number of idle instances on standby in the pool.",
			},
		},
	}
}
