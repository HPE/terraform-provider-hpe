package storage_bucket

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type storageBucketModel struct {
	ID                  types.Int64  `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	ProviderType        types.String `tfsdk:"provider_type"`
	BucketName          types.String `tfsdk:"bucket_name"`
	AccessKey           types.String `tfsdk:"access_key"`
	SecretKey           types.String `tfsdk:"secret_key"`
	Endpoint            types.String `tfsdk:"endpoint"`
	DefaultBackupTarget types.Bool   `tfsdk:"default_backup_target"`
	RetentionDays       types.Int64  `tfsdk:"retention_days"`
}

func StorageBucketSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Manages a Morpheus Storage Bucket resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "The ID of the storage bucket.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the storage bucket.",
			},
			"provider_type": schema.StringAttribute{
				Required:    true,
				Description: "The type of storage bucket (e.g. s3, azure, etc.).",
			},
			"bucket_name": schema.StringAttribute{
				Optional:    true,
				Description: "The name of the bucket.",
			},
			"access_key": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "The access key for the storage bucket.",
			},
			"secret_key": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "The secret key for the storage bucket.",
			},
			"endpoint": schema.StringAttribute{
				Optional:    true,
				Description: "The endpoint URL for the storage bucket.",
			},
			"default_backup_target": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether this is the default backup target.",
			},
			"retention_days": schema.Int64Attribute{
				Optional:    true,
				Description: "The number of days to retain files before deletion.",
			},
		},
	}
}
