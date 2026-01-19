// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package template

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"

	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/helpers"
)

func ResourceClusterPackage() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus cluster package resource.",
		CreateContext: resourceClusterPackageCreate,
		ReadContext:   resourceClusterPackageRead,
		UpdateContext: resourceClusterPackageUpdate,
		DeleteContext: resourceClusterPackageDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the cluster package",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the cluster package",
				Required:    true,
			},
			"code": {
				Type:        schema.TypeString,
				Description: "The code for the cluster package",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the cluster package",
				Optional:    true,
			},
			"package_version": {
				Type:        schema.TypeString,
				Description: "The version of the cluster package",
				Required:    true,
			},
			"type": {
				Type: schema.TypeString,
				Description: "The package category type " +
					"(apps, custom, ingress, logging, monitoring, morpheus, network, serviceMesh, storage)",
				Required: true,
			},
			"package_type": {
				Type:        schema.TypeString,
				Description: "A one word descriptor of package, such as calico, rook, prometheus, etc.",
				Required:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Whether the cluster package is enabled",
				Optional:    true,
				Computed:    true,
			},
			"repeat_install": {
				Type:        schema.TypeBool,
				Description: "Whether to support the reinstallation of the package",
				Optional:    true,
				Computed:    true,
			},
			"spec_template_ids": {
				Type:        schema.TypeList,
				Description: "A list of spec template ids associated with the cluster package",
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Required:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceClusterPackageCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	var name string
	if v, ok := d.Get("name").(string); ok {
		name = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
	}

	var code string
	if v, ok := d.Get("code").(string); ok {
		code = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("code", d.Get("code")))
	}

	var description string
	if v, ok := d.Get("description").(string); ok {
		description = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("description", d.Get("description")))
	}

	var enabled bool
	if v, ok := d.Get("enabled").(bool); ok {
		enabled = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("enabled", d.Get("enabled")))
	}

	var repeatInstall bool
	if v, ok := d.Get("repeat_install").(bool); ok {
		repeatInstall = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("repeat_install", d.Get("repeat_install")))
	}

	var packageType string
	if v, ok := d.Get("type").(string); ok {
		packageType = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("type", d.Get("type")))
	}

	var packageTypeDescriptor string
	if v, ok := d.Get("package_type").(string); ok {
		packageTypeDescriptor = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("package_type", d.Get("package_type")))
	}

	var packageVersion string
	if v, ok := d.Get("package_version").(string); ok {
		packageVersion = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("package_version", d.Get("package_version")))
	}

	specTemplateIDs := d.Get("spec_template_ids")

	req := &morpheus.Request{
		Method: "POST",
		Path:   morpheus.ClusterPackagesPath,
		Body: map[string]any{
			"clusterPackage": map[string]any{
				"name":           name,
				"code":           code,
				"description":    description,
				"enabled":        enabled,
				"repeatInstall":  repeatInstall,
				"type":           packageType,
				"packageType":    packageTypeDescriptor,
				"packageVersion": packageVersion,
				"specTemplates":  specTemplateIDs,
			},
		},
		Result: &ClusterPackageCreateResult{},
	}

	resp, err := client.Execute(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	var result *ClusterPackageCreateResult
	if v, ok := resp.Result.(*ClusterPackageCreateResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("ClusterPackageCreateResult", resp.Result))
	}

	if result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("ClusterPackageCreateResult"))
	}

	// Successfully created resource, now set id
	d.SetId(convert.Int64ToString(result.ID))

	diags = append(diags, resourceClusterPackageRead(ctx, d, meta)...)

	return diags
}

func resourceClusterPackageRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
		resp, err = client.FindClusterPackageByName(name)
	} else if id != "" {
		resp, err = client.GetClusterPackage(convert.StringToInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Cluster Package cannot be read without name or id")
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
	var result *morpheus.GetClusterPackageResult
	if v, ok := resp.Result.(*morpheus.GetClusterPackageResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("GetClusterPackageResult", resp.Result))
	}

	if result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("GetClusterPackageResult"))
	}

	clusterPackage := result.ClusterPackage
	if clusterPackage == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("ClusterPackage"))
	}

	d.SetId(convert.Int64ToString(clusterPackage.ID))
	d.Set("name", clusterPackage.Name)
	d.Set("description", clusterPackage.Description)
	d.Set("code", clusterPackage.Code)
	d.Set("package_version", clusterPackage.PackageVersion)
	d.Set("type", clusterPackage.Type)
	d.Set("package_type", clusterPackage.PackageType)
	d.Set("enabled", clusterPackage.Enabled)
	d.Set("repeat_install", clusterPackage.RepeatInstall)
	// spec templates
	var specTemplates []int64
	if clusterPackage.SpecTemplates != nil {
		// iterate over the array of spec templates
		for i := 0; i < len(clusterPackage.SpecTemplates); i++ {
			specTemplate := clusterPackage.SpecTemplates[i]
			specTemplates = append(specTemplates, specTemplate.ID)
		}
	}
	d.Set("spec_template_ids", specTemplates)

	return diags
}

func resourceClusterPackageUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	var code string
	if v, ok := d.Get("code").(string); ok {
		code = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("code", d.Get("code")))
	}

	var description string
	if v, ok := d.Get("description").(string); ok {
		description = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("description", d.Get("description")))
	}

	var enabled bool
	if v, ok := d.Get("enabled").(bool); ok {
		enabled = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("enabled", d.Get("enabled")))
	}

	var repeatInstall bool
	if v, ok := d.Get("repeat_install").(bool); ok {
		repeatInstall = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("repeat_install", d.Get("repeat_install")))
	}

	var packageType string
	if v, ok := d.Get("type").(string); ok {
		packageType = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("type", d.Get("type")))
	}

	var packageTypeDescriptor string
	if v, ok := d.Get("package_type").(string); ok {
		packageTypeDescriptor = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("package_type", d.Get("package_type")))
	}

	var packageVersion string
	if v, ok := d.Get("package_version").(string); ok {
		packageVersion = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("package_version", d.Get("package_version")))
	}

	specTemplateIDs := d.Get("spec_template_ids")

	req := &morpheus.Request{
		Body: map[string]any{
			"clusterPackage": map[string]any{
				"name":           name,
				"code":           code,
				"description":    description,
				"enabled":        enabled,
				"repeatInstall":  repeatInstall,
				"type":           packageType,
				"packageType":    packageTypeDescriptor,
				"packageVersion": packageVersion,
				"specTemplates":  specTemplateIDs,
			},
		},
	}

	resp, err := client.UpdateClusterPackage(convert.StringToInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	return resourceClusterPackageRead(ctx, d, meta)
}

func resourceClusterPackageDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
	resp, err := client.DeleteClusterPackage(convert.StringToInt64(id), req)
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

type ClusterPackageCreateResult struct {
	ID      int64             `json:"id"`
	Message string            `json:"msg"`
	Errors  map[string]string `json:"errors"`
	Success bool              `json:"success"`
}
