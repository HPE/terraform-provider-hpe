// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package morpheus

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	taskdatasource "github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/datasources/task"
	tasksdatasource "github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/datasources/tasks"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/resources/cluster"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/resources/identitysource"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/resources/integration"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/resources/task"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/resources/wiki"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/resources/workflow"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: providerSchema(),

		ResourcesMap: map[string]*schema.Resource{
			"hpe_morpheus_cluster_mks_vsphere":              cluster.ResourceClusterMKSVSphere(),
			"hpe_morpheus_identity_source_active_directory": identitysource.ResourceIdentitySourceActiveDirectory(),
			"hpe_morpheus_identity_source_saml":             identitysource.ResourceIdentitySourceSAML(),
			"hpe_morpheus_integration_ansible":              integration.ResourceIntegrationAnsible(),
			"hpe_morpheus_integration_ansible_tower":        integration.ResourceIntegrationAnsibleTower(),
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
			"hpe_morpheus_task_vro":                         task.ResourceTaskVro(),
			"hpe_morpheus_task_write_attributes":            task.ResourceTaskWriteAttributes(),
			"hpe_morpheus_wiki_page":                        wiki.ResourceWikiPage(),
			"hpe_morpheus_workflow_catalog_item":            workflow.ResourceWorkflowCatalogItem(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"hpe_morpheus_task":  taskdatasource.DataSourceMorpheusTask(),
			"hpe_morpheus_tasks": tasksdatasource.DataSourceMorpheusTasks(),
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
					Type:        schema.TypeString,
					Required:    true,
					Description: "Morpheus instance URL",
				},

				"access_token": {
					Type:          schema.TypeString,
					Optional:      true,
					Sensitive:     true,
					Description:   "Morpheus access token for authentication",
					Default:       "",
					ConflictsWith: []string{"morpheus.0.username", "morpheus.0.password"},
					// ConflictsWith: []string{"username", "password", "tenant_subdomain"},
				},

				// TODO check if this is needed in framework and here
				/*
					"tenant_subdomain": {
						Type:          schema.TypeString,
						Optional:      true,
						Description:   "The tenant subdomain used for authentication",
						ConflictsWith: []string{"access_token"},
					},

				*/

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
			},
		},
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	morpheusConfig := d.Get("morpheus").([]interface{})[0].(map[string]interface{})

	config := Config{
		Url:         morpheusConfig["url"].(string),
		AccessToken: morpheusConfig["access_token"].(string),
		Username:    morpheusConfig["username"].(string),
		Password:    morpheusConfig["password"].(string),
		Insecure:    morpheusConfig["insecure"].(bool), //.(bool),
	}

	return config.Client()
}
