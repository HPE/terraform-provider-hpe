package backup_job

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type backupJobModel struct {
	ID             types.Int64  `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	ScheduleID     types.Int64  `tfsdk:"schedule_id"`
	RetentionCount types.Int64  `tfsdk:"retention_count"`
	Enabled        types.Bool   `tfsdk:"enabled"`
	Code           types.String `tfsdk:"code"`
}

func BackupJobSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Manages a Morpheus Backup Job resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "The ID of the backup job.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the backup job.",
			},
			"schedule_id": schema.Int64Attribute{
				Optional:    true,
				Description: "The execute schedule ID.",
			},
			"retention_count": schema.Int64Attribute{
				Optional:    true,
				Description: "The retention count.",
			},
			"enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "Whether the backup job is enabled.",
			},
			"code": schema.StringAttribute{
				Optional:    true,
				Description: "The code of the backup job.",
			},
		},
	}
}
