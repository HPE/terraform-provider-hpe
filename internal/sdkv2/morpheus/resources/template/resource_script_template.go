// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package template

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"

	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/helpers"
)

func ResourceScriptTemplate() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus script template resource",
		CreateContext: resourceScriptTemplateCreate,
		ReadContext:   resourceScriptTemplateRead,
		UpdateContext: resourceScriptTemplateUpdate,
		DeleteContext: resourceScriptTemplateDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the script template",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the script template",
				Required:    true,
			},
			"labels": {
				Type: schema.TypeSet,
				Description: "The organization labels associated with the script template " +
					"(Only supported on Morpheus 5.5.3 or higher)",
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"script_type": {
				Type:         schema.TypeString,
				Description:  "The type of the script template (powershell, bash)",
				ValidateFunc: validation.StringInSlice([]string{"powershell", "bash"}, false),
				Required:     true,
			},
			"script_phase": {
				Type: schema.TypeString,
				Description: "The phase that the script should be run during " +
					"(start, stop, preProvision, provision, postProvision, preDeploy, deploy, reconfigure, teardown)",
				ValidateFunc: validation.StringInSlice([]string{
					"start", "stop", "preProvision", "provision", "postProvision",
					"preDeploy", "deploy", "reconfigure", "teardown",
				}, false),
				Required: true,
			},
			"script_content": {
				Type:        schema.TypeString,
				Description: "The content of the script template",
				Optional:    true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					oldPayload := strings.TrimSpace(old)
					newPayload := strings.TrimSpace(new)
					return oldPayload == newPayload
				},
				StateFunc: func(v any) string {
					if str, ok := v.(string); ok {
						return strings.TrimSpace(str)
					}

					return ""
				},
			},
			"run_as_user": {
				Type:        schema.TypeString,
				Description: "The name of the user account the script should run as",
				Optional:    true,
			},
			"sudo": {
				Type:        schema.TypeBool,
				Description: "Whether the script should run with sudo privileges",
				Optional:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceScriptTemplateCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics

	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	var name string
	if v, ok := d.Get("name").(string); ok {
		name = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
	}

	labelsPayload := make([]string, 0)
	if attr, ok := d.GetOk("labels"); ok {
		var labelSet *schema.Set
		if v, ok := attr.(*schema.Set); ok {
			labelSet = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("labels", attr))
		}
		labelList := labelSet.List()
		for _, s := range labelList {
			var labelStr string
			if v, ok := s.(string); ok {
				labelStr = v
			} else {
				return diag.FromErr(helpers.TypeAssertFailError("label", s))
			}
			labelsPayload = append(labelsPayload, labelStr)
		}
	}

	var scriptType string
	if v, ok := d.Get("script_type").(string); ok {
		scriptType = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("script_type", d.Get("script_type")))
	}

	var scriptPhase string
	if v, ok := d.Get("script_phase").(string); ok {
		scriptPhase = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("script_phase", d.Get("script_phase")))
	}

	var scriptContent string
	if v, ok := d.Get("script_content").(string); ok {
		scriptContent = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("script_content", d.Get("script_content")))
	}

	var runAsUser string
	if v, ok := d.Get("run_as_user").(string); ok {
		runAsUser = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("run_as_user", d.Get("run_as_user")))
	}

	var sudo bool
	if v, ok := d.Get("sudo").(bool); ok {
		sudo = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("sudo", d.Get("sudo")))
	}

	req := &morpheus.Request{
		Body: map[string]any{
			"containerScript": map[string]any{
				"name":        name,
				"labels":      labelsPayload,
				"scriptType":  scriptType,
				"scriptPhase": scriptPhase,
				"script":      scriptContent,
				"runAsUser":   runAsUser,
				"sudoUser":    sudo,
			},
		},
	}
	resp, err := client.CreateScriptTemplate(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	var result *morpheus.CreateScriptTemplateResult
	if v, ok := resp.Result.(*morpheus.CreateScriptTemplateResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("CreateScriptTemplateResult", resp.Result))
	}

	if result.ScriptTemplate == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("ScriptTemplate"))
	}
	scriptTemplate := result.ScriptTemplate
	d.SetId(convert.Int64ToString(scriptTemplate.ID))

	diags = append(diags, resourceScriptTemplateRead(ctx, d, meta)...)

	return diags
}

func resourceScriptTemplateRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics

	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

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
		resp, err = client.FindScriptTemplateByName(name)
	} else if id != "" {
		resp, err = client.GetScriptTemplate(convert.StringToInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Script template cannot be read without name or id")
	}

	if err != nil {
		// 404 is ok?
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

	// store resource data
	var result *morpheus.GetScriptTemplateResult
	if v, ok := resp.Result.(*morpheus.GetScriptTemplateResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("GetScriptTemplateResult", resp.Result))
	}

	if result.ScriptTemplate == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("ScriptTemplate"))
	}
	scriptTemplate := result.ScriptTemplate
	d.SetId(convert.Int64ToString(scriptTemplate.ID))
	d.Set("name", scriptTemplate.Name)
	d.Set("labels", scriptTemplate.Labels)
	d.Set("script_phase", scriptTemplate.ScriptPhase)
	d.Set("script_type", scriptTemplate.ScriptType)
	d.Set("script_content", scriptTemplate.Script)
	d.Set("run_as_user", scriptTemplate.RunAsUser)
	d.Set("sudo", scriptTemplate.SudoUser)

	return diags
}

func resourceScriptTemplateUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}
	id := d.Id()

	var name string
	if v, ok := d.Get("name").(string); ok {
		name = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
	}

	labelsPayload := make([]string, 0)
	if attr, ok := d.GetOk("labels"); ok {
		var labelSet *schema.Set
		if v, ok := attr.(*schema.Set); ok {
			labelSet = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("labels", attr))
		}
		labelList := labelSet.List()
		for _, s := range labelList {
			var labelStr string
			if v, ok := s.(string); ok {
				labelStr = v
			} else {
				return diag.FromErr(helpers.TypeAssertFailError("label", s))
			}
			labelsPayload = append(labelsPayload, labelStr)
		}
	}

	var scriptType string
	if v, ok := d.Get("script_type").(string); ok {
		scriptType = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("script_type", d.Get("script_type")))
	}

	var scriptPhase string
	if v, ok := d.Get("script_phase").(string); ok {
		scriptPhase = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("script_phase", d.Get("script_phase")))
	}

	var scriptContent string
	if v, ok := d.Get("script_content").(string); ok {
		scriptContent = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("script_content", d.Get("script_content")))
	}

	var runAsUser string
	if v, ok := d.Get("run_as_user").(string); ok {
		runAsUser = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("run_as_user", d.Get("run_as_user")))
	}

	var sudo bool
	if v, ok := d.Get("sudo").(bool); ok {
		sudo = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("sudo", d.Get("sudo")))
	}

	req := &morpheus.Request{
		Body: map[string]any{
			"containerScript": map[string]any{
				"name":        name,
				"labels":      labelsPayload,
				"scriptType":  scriptType,
				"scriptPhase": scriptPhase,
				"script":      scriptContent,
				"runAsUser":   runAsUser,
				"sudoUser":    sudo,
			},
		},
	}

	resp, err := client.UpdateScriptTemplate(convert.StringToInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	var result *morpheus.UpdateScriptTemplateResult
	if v, ok := resp.Result.(*morpheus.UpdateScriptTemplateResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("UpdateScriptTemplateResult", resp.Result))
	}

	if result.ScriptTemplate == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("ScriptTemplate"))
	}
	scriptTemplate := result.ScriptTemplate
	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(convert.Int64ToString(scriptTemplate.ID))

	return resourceScriptTemplateRead(ctx, d, meta)
}

func resourceScriptTemplateDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics

	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteScriptTemplate(convert.StringToInt64(id), req)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %s", resp, err)

			return nil
		}
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	d.SetId("")

	return diags
}
