// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package morpheus

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/HPE/terraform-provider-hpe/morpheus/pce"
	"github.com/HPE/terraform-provider-hpe/morpheus/pce/sdk/token/iamversion"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/automation"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/blueprint"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/catalogitem"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/cluster"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/compute"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/contact"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/credential"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/cypher"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/environment"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/form"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/identitysource"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/integration"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/job"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/license"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/network"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/optionlist"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/optiontype"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/plan"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/script"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/setting"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/task"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/template"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/tenant"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/trust"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/usergroup"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/wiki"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/workflow"

	automationds "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/datasources/automation"
	blueprintds "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/datasources/blueprint"
	catalogds "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/datasources/catalog"
	cloudds "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/datasources/cloud"
	clusterds "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/datasources/cluster"
	computeds "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/datasources/compute"
	contactds "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/datasources/contact"
	costingds "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/datasources/costing"
	credentialds "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/datasources/credential"
	cypherds "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/datasources/cypher"
	environmentds "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/datasources/environment"
	groupds "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/datasources/group"
	imageds "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/datasources/image"
	integrationds "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/datasources/integration"
	jobds "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/datasources/job"
	networkds "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/datasources/network"
	optionds "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/datasources/option"
	plands "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/datasources/plan"
	policyds "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/datasources/policy"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/datasources/provisiontype"
	storageds "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/datasources/storage"
	taskds "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/datasources/task"
	templateds "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/datasources/template"
	tenantds "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/datasources/tenant"
	trustds "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/datasources/trust"
	usergroupds "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/datasources/usergroup"
	vdids "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/datasources/vdi"
	workflowds "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/datasources/workflow"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: providerSchema(),

		ResourcesMap: map[string]*schema.Resource{
			"hpe_morpheus_app_blueprint_arm":                blueprint.ResourceAppBlueprintARM(),
			"hpe_morpheus_app_blueprint_cloud_formation":    blueprint.ResourceAppBlueprintCloudFormation(),
			"hpe_morpheus_app_blueprint_helm":               blueprint.ResourceAppBlueprintHelm(),
			"hpe_morpheus_app_blueprint_kubernetes":         blueprint.ResourceAppBlueprintKubernetes(),
			"hpe_morpheus_app_blueprint_terraform":          blueprint.ResourceAppBlueprintTerraform(),
			"hpe_morpheus_boot_script":                      script.ResourceBootScript(),
			"hpe_morpheus_catalog_item_app_blueprint":       catalogitem.ResourceCatalogItemAppBlueprint(),
			"hpe_morpheus_catalog_item_instance":            catalogitem.ResourceCatalogItemInstance(),
			"hpe_morpheus_catalog_item_workflow":            catalogitem.ResourceCatalogItemWorkflow(),
			"hpe_morpheus_cluster_layout":                   blueprint.ResourceClusterLayout(),
			"hpe_morpheus_cluster_hks_hvm":                  cluster.ResourceClusterHKSHVM(),
			"hpe_morpheus_cluster_hks_vsphere":              cluster.ResourceClusterHKSVSphere(),
			"hpe_morpheus_cluster_package":                  template.ResourceClusterPackage(),
			"hpe_morpheus_contact":                          contact.ResourceContact(),
			"hpe_morpheus_credential":                       credential.ResourceCredential(),
			"hpe_morpheus_cypher_secret":                    cypher.ResourceCypherSecret(),
			"hpe_morpheus_cypher_tfvars":                    cypher.ResourceCypherTFVars(),
			"hpe_morpheus_environment":                      environment.ResourceEnvironment(),
			"hpe_morpheus_execute_schedule":                 automation.ResourceExecuteSchedule(),
			"hpe_morpheus_file_template":                    template.ResourceFileTemplate(),
			"hpe_morpheus_form":                             form.ResourceForm(),
			"hpe_morpheus_identity_source_active_directory": identitysource.ResourceIdentitySourceActiveDirectory(),
			"hpe_morpheus_identity_source_saml":             identitysource.ResourceIdentitySourceSAML(),
			"hpe_morpheus_instance_type":                    blueprint.ResourceInstanceType(),
			"hpe_morpheus_instance_type_layout":             blueprint.ResourceInstanceTypeLayout(),
			"hpe_morpheus_integration_ansible":              integration.ResourceIntegrationAnsible(),
			"hpe_morpheus_integration_ansible_tower":        integration.ResourceIntegrationAnsibleTower(),
			"hpe_morpheus_integration_chef":                 integration.ResourceIntegrationChef(),
			"hpe_morpheus_integration_docker_registry":      integration.ResourceIntegrationDockerRegistry(),
			"hpe_morpheus_integration_git":                  integration.ResourceIntegrationGit(),
			"hpe_morpheus_integration_puppet":               integration.ResourceIntegrationPuppet(),
			"hpe_morpheus_integration_servicenow":           integration.ResourceIntegrationServiceNow(),
			"hpe_morpheus_integration_vro":                  integration.ResourceIntegrationVRO(),
			"hpe_morpheus_ip_pool_ipv4":                     network.ResourceIPPoolIPv4(),
			"hpe_morpheus_job_task":                         job.ResourceJobTask(),
			"hpe_morpheus_job_workflow":                     job.ResourceJobWorkflow(),
			"hpe_morpheus_key_pair":                         trust.ResourceKeyPair(),
			"hpe_morpheus_license":                          license.ResourceLicense(),
			"hpe_morpheus_network_domain":                   network.ResourceNetworkDomain(),
			"hpe_morpheus_node_type":                        blueprint.ResourceNodeType(),
			"hpe_morpheus_option_list_api":                  optionlist.ResourceOptionListAPI(),
			"hpe_morpheus_option_list_manual":               optionlist.ResourceOptionListManual(),
			"hpe_morpheus_option_list_rest":                 optionlist.ResourceOptionListREST(),
			"hpe_morpheus_option_type_checkbox":             optiontype.ResourceOptionTypeCheckbox(),
			"hpe_morpheus_option_type_hidden":               optiontype.ResourceOptionTypeHidden(),
			"hpe_morpheus_option_type_number":               optiontype.ResourceOptionTypeNumber(),
			"hpe_morpheus_option_type_password":             optiontype.ResourceOptionTypePassword(),
			"hpe_morpheus_option_type_radio_list":           optiontype.ResourceOptionTypeRadioList(),
			"hpe_morpheus_option_type_select_list":          optiontype.ResourceOptionTypeSelectList(),
			"hpe_morpheus_option_type_text":                 optiontype.ResourceOptionTypeText(),
			"hpe_morpheus_option_type_textarea":             optiontype.ResourceOptionTypeTextarea(),
			"hpe_morpheus_option_type_typeahead":            optiontype.ResourceOptionTypeTypeahead(),
			"hpe_morpheus_preseed_script":                   script.ResourcePreseedScript(),
			"hpe_morpheus_price":                            plan.ResourcePrice(),
			"hpe_morpheus_price_set":                        plan.ResourcePriceSet(),
			"hpe_morpheus_resource_pool_group":              compute.ResourceResourcePoolGroup(),
			"hpe_morpheus_scale_threshold":                  automation.ResourceScaleThreshold(),
			"hpe_morpheus_script_template":                  template.ResourceScriptTemplate(),
			"hpe_morpheus_security_package":                 template.ResourceSecurityPackage(),
			"hpe_morpheus_setting_appliance":                setting.ResourceSettingAppliance(),
			"hpe_morpheus_setting_backup":                   setting.ResourceSettingBackup(),
			"hpe_morpheus_setting_guidance":                 setting.ResourceSettingGuidance(),
			"hpe_morpheus_setting_monitoring":               setting.ResourceSettingMonitoring(),
			"hpe_morpheus_setting_provisioning":             setting.ResourceSettingProvisioning(),
			"hpe_morpheus_spec_template_arm":                template.ResourceSpecTemplateARM(),
			"hpe_morpheus_spec_template_cloud_formation":    template.ResourceSpecTemplateCloudFormation(),
			"hpe_morpheus_spec_template_helm":               template.ResourceSpecTemplateHelm(),
			"hpe_morpheus_spec_template_kubernetes":         template.ResourceSpecTemplateKubernetes(),
			"hpe_morpheus_spec_template_terraform":          template.ResourceSpecTemplateTerraform(),
			"hpe_morpheus_task_ansible_playbook":            task.ResourceTaskAnsiblePlaybook(),
			"hpe_morpheus_task_ansible_tower":               task.ResourceTaskAnsibleTower(),
			"hpe_morpheus_task_chef_bootstrap":              task.ResourceTaskChefBootstrap(),
			"hpe_morpheus_task_email":                       task.ResourceTaskEmail(),
			"hpe_morpheus_task_groovy_script":               task.ResourceTaskGroovyScript(),
			"hpe_morpheus_task_javascript":                  task.ResourceTaskJavaScript(),
			"hpe_morpheus_task_library_script":              task.ResourceTaskLibraryScript(),
			"hpe_morpheus_task_library_template":            task.ResourceTaskLibraryTemplate(),
			"hpe_morpheus_task_nested_workflow":             task.ResourceTaskNestedWorkflow(),
			"hpe_morpheus_task_powershell_script":           task.ResourceTaskPowerShellScript(),
			"hpe_morpheus_task_python_script":               task.ResourceTaskPythonScript(),
			"hpe_morpheus_task_restart":                     task.ResourceTaskRestart(),
			"hpe_morpheus_task_ruby_script":                 task.ResourceTaskRubyScript(),
			"hpe_morpheus_task_shell_script":                task.ResourceTaskShellScript(),
			"hpe_morpheus_task_vro":                         task.ResourceTaskVRO(),
			"hpe_morpheus_task_write_attributes":            task.ResourceTaskWriteAttributes(),
			"hpe_morpheus_tenant":                           tenant.ResourceTenant(),
			"hpe_morpheus_user_group":                       usergroup.ResourceUserGroup(),
			"hpe_morpheus_wiki_page":                        wiki.ResourceWikiPage(),
			"hpe_morpheus_workflow_operational":             workflow.ResourceWorkflowOperational(),
			"hpe_morpheus_workflow_provisioning":            workflow.ResourceWorkflowProvisioning(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"hpe_morpheus_ansible_tower_inventory":    integrationds.DataSourceAnsibleTowerInventory(),
			"hpe_morpheus_ansible_tower_job_template": integrationds.DataSourceAnsibleTowerJobTemplate(),
			"hpe_morpheus_blueprint":                  blueprintds.DataSourceBlueprint(),
			"hpe_morpheus_budget":                     costingds.DataSourceBudget(),
			"hpe_morpheus_catalog_item_type":          catalogds.DataSourceCatalogItemType(),
			"hpe_morpheus_cloud_folder":               cloudds.DataSourceCloudFolder(),
			"hpe_morpheus_cloud_type":                 cloudds.DataSourceCloudType(),
			"hpe_morpheus_clouds":                     cloudds.DataSourceClouds(),
			"hpe_morpheus_cluster_type":               clusterds.DataSourceClusterType(),
			"hpe_morpheus_contact":                    contactds.DataSourceContact(),
			"hpe_morpheus_credential":                 credentialds.DataSourceCredential(),
			"hpe_morpheus_cypher_secret":              cypherds.DataSourceCypherSecret(),
			"hpe_morpheus_environments":               environmentds.DataSourceEnvironments(),
			"hpe_morpheus_execute_schedule":           automationds.DataSourceExecuteSchedule(),
			"hpe_morpheus_file_template":              templateds.DataSourceFileTemplate(),
			"hpe_morpheus_groups":                     groupds.DataSourceGroups(),
			"hpe_morpheus_images":                     imageds.DataSourceImages(),
			"hpe_morpheus_instance_type":              blueprintds.DataSourceInstanceType(),
			"hpe_morpheus_integration":                integrationds.DataSourceIntegration(),
			"hpe_morpheus_integration_git":            integrationds.DataSourceIntegrationGit(),
			"hpe_morpheus_job":                        jobds.DataSourceJob(),
			"hpe_morpheus_key_pair":                   trustds.DataSourceKeyPair(),
			"hpe_morpheus_network_group":              networkds.DataSourceNetworkGroup(),
			"hpe_morpheus_network_subnet":             networkds.DataSourceNetworkSubnet(),
			"hpe_morpheus_networks":                   networkds.DataSourceNetworks(),
			"hpe_morpheus_node_type":                  blueprintds.DataSourceNodeType(),
			"hpe_morpheus_option_list":                optionds.DataSourceOptionList(),
			"hpe_morpheus_option_type":                optionds.DataSourceOptionType(),
			"hpe_morpheus_policies":                   policyds.DataSourcePolicies(),
			"hpe_morpheus_power_schedule":             automationds.DataSourcePowerSchedule(),
			"hpe_morpheus_price":                      plands.DataSourcePrice(),
			"hpe_morpheus_price_set":                  plands.DataSourcePriceSet(),
			"hpe_morpheus_provision_type":             provisiontype.DataSourceProvisionType(),
			"hpe_morpheus_resource_pool":              computeds.DataSourceResourcePool(),
			"hpe_morpheus_script_template":            templateds.DataSourceScriptTemplate(),
			"hpe_morpheus_security_package":           templateds.DataSourceSecurityPackage(),
			"hpe_morpheus_servicenow_workflow":        integrationds.DataSourceServiceNowWorkFlow(),
			"hpe_morpheus_spec_template":              templateds.DataSourceSpecTemplate(),
			"hpe_morpheus_storage_bucket":             storageds.DataSourceStorageBucket(),
			"hpe_morpheus_storage_volume_type":        storageds.DataSourceStorageVolumeType(),
			"hpe_morpheus_task":                       taskds.DataSourceTask(),
			"hpe_morpheus_tasks":                      taskds.DataSourceTasks(),
			"hpe_morpheus_tenant":                     tenantds.DataSourceTenant(),
			"hpe_morpheus_tenants":                    tenantds.DataSourceTenants(),
			"hpe_morpheus_user_group":                 usergroupds.DataSourceUserGroup(),
			"hpe_morpheus_user_groups":                usergroupds.DataSourceUserGroups(),
			"hpe_morpheus_vdi_pool":                   vdids.DataSourceVDIPool(),
			"hpe_morpheus_vro_workflow":               integrationds.DataSourceVROWorkflow(),
			"hpe_morpheus_workflow":                   workflowds.DataSourceWorkflow(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

// providerSchema defines the provider schema
func providerSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"morpheus": providerSchemaMorpheus(),
	}
}

// providerSchemaMorpheus defines the Morpheus provider schema
func providerSchemaMorpheus() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"url": {
					Type:     schema.TypeString,
					Optional: true,
					Description: "Morpheus instance URL. May be omitted when it is " +
						"supplied by a pce_identity or pce_disconnected_identity block.",
				},

				"access_token": {
					Type:          schema.TypeString,
					Optional:      true,
					Sensitive:     true,
					Description:   "Morpheus access token for authentication",
					Default:       "",
					ConflictsWith: []string{"morpheus.0.username", "morpheus.0.password", "morpheus.0.tenant_subdomain"},
				},

				"tenant_subdomain": {
					Type:          schema.TypeString,
					Optional:      true,
					Description:   "Morpheus tenant subdomain used for authentication",
					ConflictsWith: []string{"morpheus.0.access_token"},
				},

				"username": {
					Type:          schema.TypeString,
					Optional:      true,
					Description:   "Morpheus username for authentication, required if password is set",
					Default:       "",
					ConflictsWith: []string{"morpheus.0.access_token"},
				},

				"password": {
					Type:          schema.TypeString,
					Optional:      true,
					Sensitive:     true,
					Description:   "Morpheus password for authentication, required if username is set",
					Default:       "",
					ConflictsWith: []string{"morpheus.0.access_token"},
				},

				// defaults to false
				"insecure": {
					Type:     schema.TypeBool,
					Optional: true,
					Description: "Explicitly allow the provider to perform " +
						"\"insecure\" SSL requests. If omitted, " +
						"default value is `false`",
					Default: false,
				},

				// Mirrors the pce_identity block on the framework
				// Morpheus provider. Both providers are muxed together and
				// Terraform requires their schemas to be identical, so this
				// must match what utils/convert produces from the framework
				// schema, including descriptions.
				"pce_identity": {
					Type:        schema.TypeList,
					Optional:    true,
					MaxItems:    1,
					Description: "Configuration block for using Morpheus with PCE (Private Cloud Enterprise) Identity",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"client_id": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "GreenLake API client ID used for authentication.",
							},

							"client_secret": {
								Type:        schema.TypeString,
								Optional:    true,
								Sensitive:   true,
								Description: "GreenLake API client secret used for authentication.",
							},

							"location": {
								Type:        schema.TypeString,
								Required:    true,
								Description: "The PCE instance's Location.",
							},

							"space": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "The name of the GreenLake Space that the PCE instance is in.",
							},

							"issuer_url": {
								Type:     schema.TypeString,
								Optional: true,
								Description: `GreenLake IAM Issuer URL used to generate access tokens. ` +
									`This should be set to the "Issuer" URL of the API client.`,
							},

							"iam_token": {
								Type:      schema.TypeString,
								Optional:  true,
								Sensitive: true,
								Description: "GreenLake IAM access token. If set, token " +
									"generation from credentials is skipped.",
							},

							"broker_url": {
								Type:     schema.TypeString,
								Optional: true,
								Description: "URL of the PCE broker. Defaults to the " +
									"HPE-hosted broker if not set.",
							},
						},
					},
				},

				// Mirrors the pce_disconnected_identity block on the framework
				// Morpheus provider. As above, this must match what
				// utils/convert produces from the framework schema, including
				// descriptions and which attributes are required.
				"pce_disconnected_identity": {
					Type:     schema.TypeList,
					Optional: true,
					MaxItems: 1,
					Description: "Configuration block for using Morpheus with Disconnected PCE " +
						"(Private Cloud Enterprise) Identity",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"client_id": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "GreenLake API client ID used for authentication.",
							},

							"client_secret": {
								Type:        schema.TypeString,
								Optional:    true,
								Sensitive:   true,
								Description: "GreenLake API client secret used for authentication.",
							},

							"issuer_url": {
								Type:     schema.TypeString,
								Optional: true,
								Description: `GreenLake IAM Issuer URL used to generate access tokens. ` +
									`This should be set to the "Issuer" URL of the API client.`,
							},

							"iam_token": {
								Type:      schema.TypeString,
								Optional:  true,
								Sensitive: true,
								Description: "GreenLake IAM access token. If set, token " +
									"generation from credentials is skipped.",
							},

							"location": {
								Type:        schema.TypeString,
								Required:    true,
								Description: "The PCE instance's Location.",
							},

							"workspace_id": {
								Type:        schema.TypeString,
								Required:    true,
								Description: "The GreenLake Workspace ID that the PCE instance is in.",
							},

							"broker_url": {
								Type:     schema.TypeString,
								Required: true,
								Description: "URL of the PCE broker for this deployment. There is " +
									"no default: a Disconnected deployment has no hosted " +
									"broker to fall back to.",
							},
						},
					},
				},
			},
		},
	}
}

// missingMorpheusBlock is returned as the provider meta value when no morpheus
// provider block is configured. Legacy resources surface it as a type
// assertion failure, which reports this text back to the user.
const missingMorpheusBlock = `Morpheus resource or data source present, but possible missing morpheus provider block.
 
 provider "hpe" {
   morpheus { <- missing or duplicate?
     url = "https://example.com"
   }
 }
`

// incompleteMorpheusBlock is returned when a morpheus block is present but no
// connection details could be determined from it, either directly or by way of
// an identity block.
const incompleteMorpheusBlock = `Morpheus resource or data source present, but the morpheus provider block
does not set "url", and no usable identity block was found.

Set the connection details explicitly, or configure a pce_identity or
pce_disconnected_identity block so that they can be obtained from GreenLake:
 
 provider "hpe" {
   morpheus {
     url          = "https://example.com"
     access_token = "..."
   }
 }
`

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	morph, ok := d.GetOk("morpheus")
	if !ok {
		return missingMorpheusBlock, nil
	}

	blocks, ok := morph.([]interface{})
	if !ok || len(blocks) == 0 {
		return missingMorpheusBlock, nil
	}

	// A morpheus block with nothing set at all, as in "morpheus {}", is
	// represented as a nil element rather than an empty map. A block that
	// carries only a nested block, such as pce_identity, is a map, so this
	// guard is for the wholly empty case.
	morpheusConfig, ok := blocks[0].(map[string]interface{})
	if !ok {
		return incompleteMorpheusBlock, nil
	}

	config := Config{
		Url:             morpheusConfig["url"].(string),
		AccessToken:     morpheusConfig["access_token"].(string),
		Username:        morpheusConfig["username"].(string),
		Password:        morpheusConfig["password"].(string),
		TenantSubdomain: morpheusConfig["tenant_subdomain"].(string),
		Insecure:        morpheusConfig["insecure"].(bool), //.(bool),
	}

	// An identity block supplies the url and access token by exchanging
	// GreenLake credentials. The framework Morpheus provider performs the same
	// exchange from the same configuration, so this repeats it rather than
	// sharing a result with it.
	//
	// The framework provider rejects configurations that set both an identity
	// block and the connection details directly, so at most one source of
	// connection details is present here.
	if glc, ok := identityConfig(morpheusConfig); ok {
		url, token, err := pce.TokenExchange(ctx, glc)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		config.Url = url
		config.AccessToken = token
	}

	// Without a url the legacy client cannot reach Morpheus at all, so report
	// it here rather than failing later with a less obvious error.
	if config.Url == "" {
		return incompleteMorpheusBlock, nil
	}

	return config.Client()
}

// identityConfig reads whichever identity block is present from a morpheus
// provider block. It reports false when neither is configured.
//
// The framework provider rejects configurations that set both, so the order
// here only decides which one wins in a configuration that would already have
// been rejected.
func identityConfig(morpheusConfig map[string]interface{}) (pce.Config, bool) {
	if cfg, ok := pceIdentityConfig(morpheusConfig); ok {
		return cfg, true
	}

	return pceDisconnectedIdentityConfig(morpheusConfig)
}

// identityBlock returns the single element of a named identity block. It
// reports false when the block is absent or empty.
func identityBlock(
	morpheusConfig map[string]interface{},
	name string,
) (map[string]interface{}, bool) {
	blocks, ok := morpheusConfig[name].([]interface{})
	if !ok || len(blocks) == 0 {
		return nil, false
	}

	// A block with no attributes set is represented as a nil element.
	block, ok := blocks[0].(map[string]interface{})
	if !ok {
		return nil, false
	}

	return block, true
}

// pceIdentityConfig reads a pce_identity block from a morpheus
// provider block. It reports false when the block is absent or empty.
func pceIdentityConfig(morpheusConfig map[string]interface{}) (pce.Config, bool) {
	block, ok := identityBlock(morpheusConfig, "pce_identity")
	if !ok {
		return pce.Config{}, false
	}

	return pce.Config{
		ClientID:     stringAttr(block, "client_id"),
		ClientSecret: stringAttr(block, "client_secret"),
		Location:     stringAttr(block, "location"),
		Space:        stringAttr(block, "space"),
		IssuerURL:    stringAttr(block, "issuer_url"),
		IAMToken:     stringAttr(block, "iam_token"),
		BrokerURL:    stringAttr(block, "broker_url"),
		Version:      iamversion.GLCS,
	}, true
}

// pceDisconnectedIdentityConfig reads a pce_disconnected_identity block from a
// morpheus provider block. It reports false when the block is absent or empty.
func pceDisconnectedIdentityConfig(morpheusConfig map[string]interface{}) (pce.Config, bool) {
	block, ok := identityBlock(morpheusConfig, "pce_disconnected_identity")
	if !ok {
		return pce.Config{}, false
	}

	return pce.Config{
		ClientID:     stringAttr(block, "client_id"),
		ClientSecret: stringAttr(block, "client_secret"),
		IssuerURL:    stringAttr(block, "issuer_url"),
		IAMToken:     stringAttr(block, "iam_token"),
		Location:     stringAttr(block, "location"),
		WorkspaceID:  stringAttr(block, "workspace_id"),
		BrokerURL:    stringAttr(block, "broker_url"),
		Version:      iamversion.GLP,
	}, true
}

// stringAttr reads a string attribute, returning "" when it is absent or of an
// unexpected type. This matches how the framework provider reads null values,
// which matters because both providers must resolve the same Morpheus instance
// from the same block.
func stringAttr(m map[string]interface{}, key string) string {
	v, _ := m[key].(string)

	return v
}
