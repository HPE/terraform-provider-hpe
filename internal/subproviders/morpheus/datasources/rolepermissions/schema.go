// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package rolepermissions

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

//nolint:revive,lll
func RolePermissionsDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"json": schema.StringAttribute{
				CustomType:  jsontypes.NormalizedType{},
				Description: "Normalized permissions JSON data",
				Computed:    true,
			},
			"feature_permissions": schema.SetNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"access": schema.StringAttribute{
							Required:            true,
							Description:         "The new access level.",
							MarkdownDescription: "The new access level.",
							Validators: []validator.String{
								stringvalidator.OneOf(
									"full",
									"full_decrypted",
									"group",
									"listfiles",
									"managerules",
									"no",
									"none",
									"provision",
									"read",
									"rolemappings",
									"user",
									"view",
									"yes",
								),
							},
						},
						"code": schema.StringAttribute{
							Required:            true,
							Description:         "`code` of the feature permission",
							MarkdownDescription: "`code` of the feature permission",
							Validators: []validator.String{
								stringvalidator.OneOf(
									"account-usage",
									"activity",
									"admin-accounts",
									"admin-accounts-users",
									"admin-appliance",
									"admin-backupSettings",
									"admin-certificates",
									"admin-clients",
									"admin-cm",
									"admin-containers",
									"admin-distributed-workers",
									"admin-environments",
									"admin-global-policies",
									"admin-groups",
									"admin-guidanceSettings",
									"admin-health",
									"admin-identity-sources",
									"admin-keypairs",
									"admin-licenses",
									"admin-logSettings",
									"admin-monitorSettings",
									"admin-motd",
									"admin-packages",
									"admin-plugins",
									"admin-policies",
									"admin-profiles",
									"admin-provisioningSettings",
									"admin-roles",
									"admin-servers",
									"admin-servicePlans",
									"admin-users",
									"admin-whitelabel",
									"admin-zones",
									"app-templates",
									"apps",
									"arm-template",
									"automation-services",
									"backup-services",
									"backups",
									"billing",
									"catalog",
									"cloudFormation-template",
									"code-repositories",
									"credentials",
									"dashboard",
									"deployment-services",
									"deployments",
									"execution-request",
									"executions",
									"guidance",
									"helm-template",
									"infrastructure-boot",
									"infrastructure-cluster",
									"infrastructure-dhcp-pool",
									"infrastructure-domains",
									"infrastructure-ippools",
									"infrastructure-kube-cntl",
									"infrastructure-loadbalancer",
									"infrastructure-move-server",
									"infrastructure-nat",
									"infrastructure-network-dhcp-relay",
									"infrastructure-network-dhcp-routes",
									"infrastructure-network-dhcp-server",
									"infrastructure-network-firewalls",
									"infrastructure-network-integrations",
									"infrastructure-network-router-firewalls",
									"infrastructure-network-router-interfaces",
									"infrastructure-network-router-redistribution",
									"infrastructure-network-router-routes",
									"infrastructure-network-server-groups",
									"infrastructure-networks",
									"infrastructure-proxies",
									"infrastructure-router-dhcp-binding",
									"infrastructure-router-dhcp-relay",
									"infrastructure-routers",
									"infrastructure-securityGroups",
									"infrastructure-state",
									"infrastructure-storage",
									"infrastructure-storage-browser",
									"integrations-ansible",
									"job-executions",
									"job-templates",
									"kubernetes-template",
									"library-advanced-node-type-options",
									"library-options",
									"library-templates",
									"logs",
									"monitoring",
									"operations-alarms",
									"operations-approvals",
									"operations-budgets",
									"operations-invoices",
									"operations-wiki",
									"projects",
									"provisioning",
									"provisioning-add",
									"provisioning-admin",
									"provisioning-clone",
									"provisioning-delete",
									"provisioning-edit",
									"provisioning-environment",
									"provisioning-execute-script",
									"provisioning-execute-task",
									"provisioning-execute-workflow",
									"provisioning-force-delete",
									"provisioning-import-image",
									"provisioning-lock",
									"provisioning-power",
									"provisioning-reconfigure",
									"provisioning-reconfigure-add-disk",
									"provisioning-reconfigure-add-network",
									"provisioning-reconfigure-change-plan",
									"provisioning-reconfigure-disk-type",
									"provisioning-reconfigure-modify-disk",
									"provisioning-reconfigure-modify-network",
									"provisioning-reconfigure-remove-disk",
									"provisioning-reconfigure-remove-network",
									"provisioning-remove-control",
									"provisioning-scale",
									"provisioning-settings",
									"provisioning-state",
									"reports",
									"reports-analytics",
									"scheduling-execute",
									"scheduling-power",
									"security-scan",
									"service-catalog",
									"service-catalog-dashboard",
									"service-catalog-inventory",
									"services-archives",
									"services-cypher",
									"services-image-builder",
									"services-kubernetes",
									"services-network-registry",
									"services-vdi-copy",
									"services-vdi-pools",
									"services-vdi-printer",
									"snapshots",
									"task-scripts",
									"tasks",
									"terminal",
									"terminal-access",
									"terraform-template",
									"thresholds",
									"trust-services",
									"virtual-images",
								),
							},
						},
					},
				},
				Optional:            true,
				Description:         "Set the access level for the specified permissions.",
				MarkdownDescription: "Set the access level for the specified permissions.",
			},
			"default_group_access": schema.StringAttribute{
				Optional:            true,
				Description:         "Set the default access level for for groups (sites). Only applies to user roles.",
				MarkdownDescription: "Set the default access level for for groups (sites). Only applies to user roles.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"default",
						"full",
						"read",
						"none",
					),
				},
			},
			"default_cloud_access": schema.StringAttribute{
				Optional:            true,
				Description:         "Set the default access level for for clouds (zones). Only applies to base account (tenant) roles.",
				MarkdownDescription: "Set the default access level for for clouds (zones). Only applies to base account (tenant) roles.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"default",
						"full",
						"read",
						"none",
					),
				},
			},
			"default_blueprint_access": schema.StringAttribute{
				Optional:            true,
				Description:         "Set the default access level for blueprints",
				MarkdownDescription: "Set the default access level for blueprints",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"full",
						"none",
					),
				},
			},
			"default_catalog_item_type_access": schema.StringAttribute{
				Optional:            true,
				Description:         "Set the default access level for catalog item types",
				MarkdownDescription: "Set the default access level for catalog item types",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"full",
						"none",
					),
				},
			},
			"default_instance_type_access": schema.StringAttribute{
				Optional:            true,
				Description:         "Set the default access level for for instance types",
				MarkdownDescription: "Set the default access level for for instance types",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"full",
						"none",
					),
				},
			},
			"default_persona_access": schema.StringAttribute{
				Optional:            true,
				Description:         "Set the default access level for personas",
				MarkdownDescription: "Set the default access level for personas",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"full",
						"none",
					),
				},
			},
			"default_report_type_access": schema.StringAttribute{
				Optional:            true,
				Description:         "Set the default access level for report types",
				MarkdownDescription: "Set the default access level for report types",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"full",
						"none",
					),
				},
			},
			"default_task_access": schema.StringAttribute{
				Optional:            true,
				Description:         "Set the default access level for tasks",
				MarkdownDescription: "Set the default access level for tasks",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"full",
						"none",
					),
				},
			},
			"default_workflow_access": schema.StringAttribute{
				Optional:            true,
				Description:         "Set the default access level for workflows (taskSets)",
				MarkdownDescription: "Set the default access level for workflows (taskSets)",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"full",
						"none",
					),
				},
			},
			"default_vdi_pool_access": schema.StringAttribute{
				Optional:            true,
				Description:         "Set the default access level for VDI pools",
				MarkdownDescription: "Set the default access level for VDI pools",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"full",
						"none",
					),
				},
			},
		},
	}
}

//nolint:revive
type RolePermissionsModel struct {
	Json                         jsontypes.Normalized `tfsdk:"json"`
	FeaturePermissions           types.Set            `tfsdk:"feature_permissions"`
	DefaultGroupAccess           types.String         `tfsdk:"default_group_access"`
	DefaultCloudAccess           types.String         `tfsdk:"default_cloud_access"`
	DefaultBlueprintAccess       types.String         `tfsdk:"default_blueprint_access"`
	DefaultCatalogItemTypeAccess types.String         `tfsdk:"default_catalog_item_type_access"`
	DefaultInstanceTypeAccess    types.String         `tfsdk:"default_instance_type_access"`
	DefaultPersonaAccess         types.String         `tfsdk:"default_persona_access"`
	DefaultReportTypeAccess      types.String         `tfsdk:"default_report_type_access"`
	DefaultTaskAccess            types.String         `tfsdk:"default_task_access"`
	DefaultWorkflowAccess        types.String         `tfsdk:"default_workflow_access"`
	DefaultVdiPoolAccess         types.String         `tfsdk:"default_vdi_pool_access"`
}
