package backup

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type backupModel struct {
	ID             types.Int64  `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	InstanceID     types.Int64  `tfsdk:"instance_id"`
	ContainerID    types.Int64  `tfsdk:"container_id"`
	BackupType     types.String `tfsdk:"backup_type"`
	ScheduleID     types.Int64  `tfsdk:"schedule_id"`
	RetentionCount types.Int64  `tfsdk:"retention_count"`
	Enabled        types.Bool   `tfsdk:"enabled"`
}

func BackupSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Manages a Morpheus Backup resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "The ID of the backup.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the backup.",
			},
			"instance_id": schema.Int64Attribute{
				Optional:    true,
				Description: "The ID of the instance to backup.",
			},
			"container_id": schema.Int64Attribute{
				Optional:    true,
				Description: "The ID of the container to backup.",
			},
			"backup_type": schema.StringAttribute{
				Optional:    true,
				Description: "The backup type code.",
			},
			"schedule_id": schema.Int64Attribute{
				Optional:    true,
				Description: "The ID of the execute schedule.",
			},
			"retention_count": schema.Int64Attribute{
				Optional:    true,
				Description: "The retention count.",
			},
			"enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "Whether the backup is enabled.",
			},
		},
	}
}
