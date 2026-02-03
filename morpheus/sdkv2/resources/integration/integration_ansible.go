// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package integration

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"

	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/convert"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/helpers"
)

func ResourceIntegrationAnsible() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides an ansible integration resource",
		CreateContext: resourceIntegrationAnsibleCreate,
		ReadContext:   resourceIntegrationAnsibleRead,
		UpdateContext: resourceIntegrationAnsibleUpdate,
		DeleteContext: resourceIntegrationAnsibleDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the ansible integration",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the ansible integration",
				Required:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Whether the ansible integration is enabled",
				Optional:    true,
				Computed:    true,
			},
			"url": {
				Type:        schema.TypeString,
				Description: "The url of the ansible repository",
				Required:    true,
			},
			"default_branch": {
				Type:        schema.TypeString,
				Description: "The default branch of the ansible repository",
				Optional:    true,
				Computed:    true,
			},
			"playbooks_path": {
				Type:        schema.TypeString,
				Description: "The path in the repository of the Ansible playbooks relative to the Git url",
				Optional:    true,
				Computed:    true,
			},
			"roles_path": {
				Type:        schema.TypeString,
				Description: "The path in the repository of the Ansible roles relative to the Git url",
				Optional:    true,
				Computed:    true,
			},
			"group_variables_path": {
				Type:        schema.TypeString,
				Description: "The path in the repository of the Ansible group variables relative to the Git url",
				Optional:    true,
				Computed:    true,
			},
			"host_variables_path": {
				Type:        schema.TypeString,
				Description: "The path in the repository of the Ansible host variables relative to the Git url",
				Optional:    true,
				Computed:    true,
			},
			"enable_ansible_galaxy_install": {
				Type:        schema.TypeBool,
				Description: "Whether to install the Ansible roles defined in the requirements.yml",
				Optional:    true,
				Computed:    true,
			},
			"enable_verbose_logging": {
				Type:        schema.TypeBool,
				Description: "Whether verbose logging is used during the execution of the ansible playbook",
				Optional:    true,
				Computed:    true,
			},
			"enable_agent_command_bus": {
				Type:        schema.TypeBool,
				Description: "Whether the agent command bus is used to execute the ansible playbook",
				Optional:    true,
				Computed:    true,
			},
			"username": {
				Type:        schema.TypeString,
				Description: "The username of the account used to authenticate to the ansible repository",
				Optional:    true,
				Computed:    true,
			},
			"password": {
				Type:        schema.TypeString,
				Description: "The password of the account used to authenticate to the ansible repository",
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					h := sha256.New()
					h.Write([]byte(new))
					sha256Hash := hex.EncodeToString(h.Sum(nil))

					return strings.EqualFold(old, sha256Hash)
				},
				DiffSuppressOnRefresh: true,
			},
			"access_token": {
				Type:        schema.TypeString,
				Description: "The access token of the account used to authenticate to the ansible repository",
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					h := sha256.New()
					h.Write([]byte(new))
					sha256Hash := hex.EncodeToString(h.Sum(nil))

					return strings.EqualFold(old, sha256Hash)
				},
			},
			"key_pair_id": {
				Type:        schema.TypeInt,
				Description: "The ID of the key pair used to authenticate to the ansible repository",
				Optional:    true,
				Computed:    true,
			},
			"enable_git_caching": {
				Type:        schema.TypeBool,
				Description: "Whether the git repository is cached",
				Optional:    true,
				Computed:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceIntegrationAnsibleCreate(
	ctx context.Context,
	d *schema.ResourceData,
	meta any,
) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	var diags diag.Diagnostics

	integration := make(map[string]any)

	var name string
	if v, ok := d.Get("name").(string); ok {
		name = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
	}
	integration["name"] = name

	var enabled bool
	if v, ok := d.Get("enabled").(bool); ok {
		enabled = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("enabled", d.Get("enabled")))
	}
	integration["enabled"] = enabled

	integration["type"] = "ansible"

	var url string
	if v, ok := d.Get("url").(string); ok {
		url = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("url", d.Get("url")))
	}
	integration["serviceUrl"] = url

	var username string
	if v, ok := d.Get("username").(string); ok {
		username = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("username", d.Get("username")))
	}
	integration["serviceUsername"] = username

	var password string
	if v, ok := d.Get("password").(string); ok {
		password = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("password", d.Get("password")))
	}
	integration["servicePassword"] = password

	var accessToken string
	if v, ok := d.Get("access_token").(string); ok {
		accessToken = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("access_token", d.Get("access_token")))
	}
	integration["serviceToken"] = accessToken

	var keyPairID int
	if v, ok := d.Get("key_pair_id").(int); ok {
		keyPairID = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("key_pair_id", d.Get("key_pair_id")))
	}
	integration["serviceKey"] = keyPairID

	config := make(map[string]any)

	var defaultBranch string
	if v, ok := d.Get("default_branch").(string); ok {
		defaultBranch = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("default_branch", d.Get("default_branch")))
	}
	config["defaultBranch"] = defaultBranch
	config["ansibleDefaultBranch"] = defaultBranch

	var cacheEnabled bool
	if v, ok := d.Get("enable_git_caching").(bool); ok {
		cacheEnabled = v
	} else {
		return diag.FromErr(
			helpers.TypeAssertFailError("enable_git_caching", d.Get("enable_git_caching")),
		)
	}
	config["cacheEnabled"] = cacheEnabled

	var playbooksPath string
	if v, ok := d.Get("playbooks_path").(string); ok {
		playbooksPath = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("playbooks_path", d.Get("playbooks_path")))
	}
	config["ansiblePlaybooks"] = playbooksPath

	var rolesPath string
	if v, ok := d.Get("roles_path").(string); ok {
		rolesPath = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("roles_path", d.Get("roles_path")))
	}
	config["ansibleRoles"] = rolesPath

	var groupVariablesPath string
	if v, ok := d.Get("group_variables_path").(string); ok {
		groupVariablesPath = v
	} else {
		return diag.FromErr(
			helpers.TypeAssertFailError("group_variables_path", d.Get("group_variables_path")),
		)
	}
	config["ansibleGroupVars"] = groupVariablesPath

	var hostVariablesPath string
	if v, ok := d.Get("host_variables_path").(string); ok {
		hostVariablesPath = v
	} else {
		return diag.FromErr(
			helpers.TypeAssertFailError("host_variables_path", d.Get("host_variables_path")),
		)
	}
	config["ansibleHostVars"] = hostVariablesPath

	var galaxyEnabled bool
	if v, ok := d.Get("enable_ansible_galaxy_install").(bool); ok {
		galaxyEnabled = v
	} else {
		return diag.FromErr(
			helpers.TypeAssertFailError(
				"enable_ansible_galaxy_install",
				d.Get("enable_ansible_galaxy_install"),
			),
		)
	}
	config["ansibleGalaxyEnabled"] = galaxyEnabled

	var verboseLogging bool
	if v, ok := d.Get("enable_verbose_logging").(bool); ok {
		verboseLogging = v
	} else {
		return diag.FromErr(
			helpers.TypeAssertFailError("enable_verbose_logging", d.Get("enable_verbose_logging")),
		)
	}
	config["ansibleVerbose"] = verboseLogging

	var commandBus bool
	if v, ok := d.Get("enable_agent_command_bus").(bool); ok {
		commandBus = v
	} else {
		return diag.FromErr(
			helpers.TypeAssertFailError(
				"enable_agent_command_bus",
				d.Get("enable_agent_command_bus"),
			),
		)
	}
	config["ansibleCommandBus"] = commandBus

	integration["config"] = config

	req := &morpheus.Request{
		Body: map[string]any{
			"integration": integration,
		},
	}

	resp, err := client.CreateIntegration(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	var result *morpheus.CreateIntegrationResult
	if v, ok := resp.Result.(*morpheus.CreateIntegrationResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.Integration == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Integration"))
	}

	integrationResult := result.Integration
	d.SetId(convert.Int64ToString(integrationResult.ID))

	diags = append(diags, resourceIntegrationAnsibleRead(ctx, d, meta)...)

	return diags
}

func resourceIntegrationAnsibleRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	var diags diag.Diagnostics

	id := d.Id()

	var name string
	if v, ok := d.Get("name").(string); ok {
		name = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
	}

	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindIntegrationByName(name)
	} else if id != "" {
		resp, err = client.GetIntegration(convert.StringToInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Integration cannot be read without name or id")
	}

	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %s", resp, err)
			log.Printf("Forcing recreation of resource")
			d.SetId("")

			return diags
		}

		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	var result *morpheus.GetIntegrationResult
	if v, ok := resp.Result.(*morpheus.GetIntegrationResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.Integration == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Integration"))
	}

	integration := result.Integration
	d.SetId(convert.Int64ToString(integration.ID))
	d.Set("name", integration.Name)
	d.Set("enabled", integration.Enabled)
	d.Set("url", integration.URL)
	d.Set("username", integration.Username)
	d.Set("password", integration.PasswordHash)
	d.Set("access_token", integration.TokenHash)

	if integration.ServiceKey.ID != 0 {
		d.Set("key_pair_id", integration.ServiceKey.ID)
	}

	d.Set("default_branch", integration.Config.DefaultBranch)
	d.Set("enable_git_caching", integration.Config.CacheEnabled)
	d.Set("playbooks_path", integration.Config.AnsiblePlaybooks)
	d.Set("roles_path", integration.Config.AnsibleRoles)
	d.Set("group_variables_path", integration.Config.AnsibleGroupVars)
	d.Set("host_variables_path", integration.Config.AnsibleHostVars)
	d.Set("enable_ansible_galaxy_install", integration.Config.AnsibleGalaxyEnabled)
	d.Set("enable_verbose_logging", integration.Config.AnsibleVerbose)
	d.Set("enable_agent_command_bus", integration.Config.AnsibleCommandBus)

	return diags
}

func resourceIntegrationAnsibleUpdate(
	ctx context.Context,
	d *schema.ResourceData,
	meta any,
) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	id := d.Id()

	integration := make(map[string]any)

	var name string
	if v, ok := d.Get("name").(string); ok {
		name = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
	}
	integration["name"] = name

	var enabled bool
	if v, ok := d.Get("enabled").(bool); ok {
		enabled = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("enabled", d.Get("enabled")))
	}
	integration["enabled"] = enabled

	integration["type"] = "ansible"

	var url string
	if v, ok := d.Get("url").(string); ok {
		url = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("url", d.Get("url")))
	}
	integration["serviceUrl"] = url

	var username string
	if v, ok := d.Get("username").(string); ok {
		username = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("username", d.Get("username")))
	}
	integration["serviceUsername"] = username

	var password string
	if v, ok := d.Get("password").(string); ok {
		password = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("password", d.Get("password")))
	}
	integration["servicePassword"] = password

	var accessToken string
	if v, ok := d.Get("access_token").(string); ok {
		accessToken = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("access_token", d.Get("access_token")))
	}
	integration["serviceToken"] = accessToken

	var keyPairID int
	if v, ok := d.Get("key_pair_id").(int); ok {
		keyPairID = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("key_pair_id", d.Get("key_pair_id")))
	}
	integration["serviceKey"] = keyPairID

	config := make(map[string]any)

	var defaultBranch string
	if v, ok := d.Get("default_branch").(string); ok {
		defaultBranch = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("default_branch", d.Get("default_branch")))
	}
	config["defaultBranch"] = defaultBranch
	config["ansibleDefaultBranch"] = defaultBranch

	var cacheEnabled bool
	if v, ok := d.Get("enable_git_caching").(bool); ok {
		cacheEnabled = v
	} else {
		return diag.FromErr(
			helpers.TypeAssertFailError("enable_git_caching", d.Get("enable_git_caching")),
		)
	}
	config["cacheEnabled"] = cacheEnabled

	var playbooksPath string
	if v, ok := d.Get("playbooks_path").(string); ok {
		playbooksPath = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("playbooks_path", d.Get("playbooks_path")))
	}
	config["ansiblePlaybooks"] = playbooksPath

	var rolesPath string
	if v, ok := d.Get("roles_path").(string); ok {
		rolesPath = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("roles_path", d.Get("roles_path")))
	}
	config["ansibleRoles"] = rolesPath

	var groupVariablesPath string
	if v, ok := d.Get("group_variables_path").(string); ok {
		groupVariablesPath = v
	} else {
		return diag.FromErr(
			helpers.TypeAssertFailError("group_variables_path", d.Get("group_variables_path")),
		)
	}
	config["ansibleGroupVars"] = groupVariablesPath

	var hostVariablesPath string
	if v, ok := d.Get("host_variables_path").(string); ok {
		hostVariablesPath = v
	} else {
		return diag.FromErr(
			helpers.TypeAssertFailError("host_variables_path", d.Get("host_variables_path")),
		)
	}
	config["ansibleHostVars"] = hostVariablesPath

	var galaxyEnabled bool
	if v, ok := d.Get("enable_ansible_galaxy_install").(bool); ok {
		galaxyEnabled = v
	} else {
		return diag.FromErr(
			helpers.TypeAssertFailError(
				"enable_ansible_galaxy_install",
				d.Get("enable_ansible_galaxy_install"),
			),
		)
	}
	config["ansibleGalaxyEnabled"] = galaxyEnabled

	var verboseLogging bool
	if v, ok := d.Get("enable_verbose_logging").(bool); ok {
		verboseLogging = v
	} else {
		return diag.FromErr(
			helpers.TypeAssertFailError("enable_verbose_logging", d.Get("enable_verbose_logging")),
		)
	}
	config["ansibleVerbose"] = verboseLogging

	var commandBus bool
	if v, ok := d.Get("enable_agent_command_bus").(bool); ok {
		commandBus = v
	} else {
		return diag.FromErr(
			helpers.TypeAssertFailError(
				"enable_agent_command_bus",
				d.Get("enable_agent_command_bus"),
			),
		)
	}
	config["ansibleCommandBus"] = commandBus

	integration["config"] = config

	req := &morpheus.Request{
		Body: map[string]any{
			"integration": integration,
		},
	}

	resp, err := client.UpdateIntegration(convert.StringToInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	var result *morpheus.UpdateIntegrationResult
	if v, ok := resp.Result.(*morpheus.UpdateIntegrationResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.Integration == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Integration"))
	}

	integrationResult := result.Integration
	d.SetId(convert.Int64ToString(integrationResult.ID))

	return resourceIntegrationAnsibleRead(ctx, d, meta)
}

func resourceIntegrationAnsibleDelete(
	ctx context.Context,
	d *schema.ResourceData,
	meta any,
) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteIntegration(convert.StringToInt64(id), req)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %s", resp, err)

			return diag.FromErr(err)
		}

		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	d.SetId("")

	return diags
}
