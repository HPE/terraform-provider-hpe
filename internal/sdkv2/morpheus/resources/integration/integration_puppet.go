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

	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/helpers"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

const (
	puppetFireNowTrue = "true"
)

func ResourceIntegrationPuppet() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a puppet integration resource",
		CreateContext: resourceIntegrationPuppetCreate,
		ReadContext:   resourceIntegrationPuppetRead,
		UpdateContext: resourceIntegrationPuppetUpdate,
		DeleteContext: resourceIntegrationPuppetDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the puppet integration",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the puppet integration",
				Required:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Whether the puppet integration is enabled",
				Optional:    true,
				Computed:    true,
			},
			"puppet_master_hostname": {
				Type:        schema.TypeString,
				Description: "The hostname of the puppet server",
				Required:    true,
			},
			"allow_immediate_execution": {
				Type:        schema.TypeBool,
				Description: "Whether to trigger the immediate execution of a puppet agent run",
				Optional:    true,
				Computed:    true,
			},
			"puppet_master_ssh_username": {
				Type: schema.TypeString,
				Description: "The username of the account on the puppet server used to " +
					"trigger the immediate execution of a puppet agent run",
				Optional: true,
				Computed: true,
			},
			"puppet_master_ssh_password": {
				Type: schema.TypeString,
				Description: "The password of the account on the puppet server used to " +
					"trigger the immediate execution of a puppet agent run",
				Optional:  true,
				Computed:  true,
				Sensitive: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					h := sha256.New()
					h.Write([]byte(new))
					sha256Hash := hex.EncodeToString(h.Sum(nil))

					return strings.EqualFold(old, sha256Hash)
				},
				DiffSuppressOnRefresh: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceIntegrationPuppetCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	// Warning or errors can be collected in a slice type
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
	integration["type"] = "puppet"

	config := make(map[string]any)

	var puppetMaster string
	if v, ok := d.Get("puppet_master_hostname").(string); ok {
		puppetMaster = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("puppet_master_hostname", d.Get("puppet_master_hostname")))
	}
	config["puppetMaster"] = puppetMaster

	var allowImmediateExecution bool
	if v, ok := d.Get("allow_immediate_execution").(bool); ok {
		allowImmediateExecution = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("allow_immediate_execution", d.Get("allow_immediate_execution")))
	}
	if allowImmediateExecution {
		config["puppetFireNow"] = puppetFireNowTrue
	} else {
		config["puppetFireNow"] = "false"
	}

	var puppetSSHUser string
	if v, ok := d.Get("puppet_master_ssh_username").(string); ok {
		puppetSSHUser = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("puppet_master_ssh_username", d.Get("puppet_master_ssh_username")))
	}
	config["puppetSshUser"] = puppetSSHUser

	var puppetSSHPassword string
	if v, ok := d.Get("puppet_master_ssh_password").(string); ok {
		puppetSSHPassword = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("puppet_master_ssh_password", d.Get("puppet_master_ssh_password")))
	}
	config["puppetSshPassword"] = puppetSSHPassword

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

	result := resp.Result.(*morpheus.CreateIntegrationResult)
	if result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("result"))
	}
	integrationResult := result.Integration
	if integrationResult == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Integration"))
	}
	// Successfully created resource, now set id
	d.SetId(convert.Int64ToString(integrationResult.ID))

	diags = append(diags, resourceIntegrationPuppetRead(ctx, d, meta)...)

	return diags
}

func resourceIntegrationPuppetRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()

	var name string
	if v, ok := d.Get("name").(string); ok {
		name = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
	}

	// lookup by name if we do not have an id yet
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
		} else {
			log.Printf("API FAILURE: %s - %s", resp, err)

			return diag.FromErr(err)
		}
	}
	log.Printf("API RESPONSE: %s", resp)

	// store resource data
	result := resp.Result.(*morpheus.GetIntegrationResult)
	if result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("result"))
	}
	integration := result.Integration
	if integration == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Integration"))
	}
	d.SetId(convert.Int64ToString(integration.ID))
	d.Set("name", integration.Name)
	d.Set("enabled", integration.Enabled)
	d.Set("puppet_master_hostname", integration.Config.PuppetMaster)
	if integration.Config.PuppetFireNow == puppetFireNowTrue {
		d.Set("allow_immediate_execution", true)
	} else {
		d.Set("allow_immediate_execution", false)
	}
	d.Set("puppet_master_ssh_username", integration.Config.PuppetSshUser)
	d.Set("puppet_master_ssh_password", integration.Config.PuppetSshPasswordHash)

	return diags
}

func resourceIntegrationPuppetUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
	integration["type"] = "puppet"

	config := make(map[string]any)

	var puppetMaster string
	if v, ok := d.Get("puppet_master_hostname").(string); ok {
		puppetMaster = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("puppet_master_hostname", d.Get("puppet_master_hostname")))
	}
	config["puppetMaster"] = puppetMaster

	var allowImmediateExecution bool
	if v, ok := d.Get("allow_immediate_execution").(bool); ok {
		allowImmediateExecution = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("allow_immediate_execution", d.Get("allow_immediate_execution")))
	}
	if allowImmediateExecution {
		config["puppetFireNow"] = puppetFireNowTrue
	} else {
		config["puppetFireNow"] = "false"
	}

	var puppetSSHUser string
	if v, ok := d.Get("puppet_master_ssh_username").(string); ok {
		puppetSSHUser = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("puppet_master_ssh_username", d.Get("puppet_master_ssh_username")))
	}
	config["puppetSshUser"] = puppetSSHUser

	var puppetSSHPassword string
	if v, ok := d.Get("puppet_master_ssh_password").(string); ok {
		puppetSSHPassword = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("puppet_master_ssh_password", d.Get("puppet_master_ssh_password")))
	}
	config["puppetSshPassword"] = puppetSSHPassword

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
	result := resp.Result.(*morpheus.UpdateIntegrationResult)
	if result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("result"))
	}
	integrationResult := result.Integration
	if integrationResult == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Integration"))
	}

	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(convert.Int64ToString(integrationResult.ID))

	return resourceIntegrationPuppetRead(ctx, d, meta)
}

func resourceIntegrationPuppetDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteIntegration(convert.StringToInt64(id), req)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %s", resp, err)

			return diag.FromErr(err)
		} else {
			log.Printf("API FAILURE: %s - %s", resp, err)

			return diag.FromErr(err)
		}
	}
	log.Printf("API RESPONSE: %s", resp)
	d.SetId("")

	return diags
}
