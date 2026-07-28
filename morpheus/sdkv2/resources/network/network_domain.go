// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package network

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	morpheus "github.com/HPE/terraform-provider-hpe/internal/sdk/legacy"

	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/convert"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/helpers"
)

func ResourceNetworkDomain() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus network domain resource.",
		CreateContext: resourceNetworkDomainCreate,
		ReadContext:   resourceNetworkDomainRead,
		UpdateContext: resourceNetworkDomainUpdate,
		DeleteContext: resourceNetworkDomainDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of the network domain",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "The name of the network domain",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "The user friendly description of the network domain",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"public_zone": {
				Description: "Whether the domain will be public or private",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"auto_join_domain": {
				Description: "Whether to automatically join machines to the domain",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Deprecated:  "Has no effect for network domains. Use domain_controller instead.",
			},
			"domain_controller": {
				Description: "The domain controller used to facilitate an automated domain join operation",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"domain_username": {
				Description: "The username of the account used to facilitate an automated domain join operation",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"domain_password": {
				Description:   "The password of the account used to facilitate an automated domain join operation",
				Type:          schema.TypeString,
				Sensitive:     true,
				Optional:      true,
				Deprecated:    "Use domain_password_wo instead; write-only values are not stored in state.",
				ConflictsWith: []string{"domain_password_wo"},
			},
			"domain_password_wo": {
				Description:   "The account password for an automated domain join (write-only, not stored in state).",
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				WriteOnly:     true,
				ConflictsWith: []string{"domain_password"},
			},
			"domain_password_wo_version": {
				Description: "Increment to re-send domain_password_wo; write-only values are not stored in state.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"active": {
				Description: "The state of the network domain",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"visibility": {
				Description:  "Determines whether the resource is visible in sub-tenants or not",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"private", "public", ""}, false),
				Default:      "private",
			},
			"tenant_id": {
				Description: "The tenant to assign the network domain",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceNetworkDomainCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	var diags diag.Diagnostics

	var name string
	if v, ok := d.Get("name").(string); ok {
		name = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
	}

	var description string
	if v, ok := d.Get("description").(string); ok {
		description = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("description", d.Get("description")))
	}

	var publicZone bool
	if v, ok := d.Get("public_zone").(bool); ok {
		publicZone = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("public_zone", d.Get("public_zone")))
	}

	var visibility string
	if v, ok := d.Get("visibility").(string); ok {
		visibility = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("visibility", d.Get("visibility")))
	}

	var active bool
	if v, ok := d.Get("active").(bool); ok {
		active = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("active", d.Get("active")))
	}

	var domainController bool
	if v, ok := d.Get("domain_controller").(bool); ok {
		domainController = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("domain_controller", d.Get("domain_controller")))
	}

	var domainUsername string
	if v, ok := d.Get("domain_username").(string); ok {
		domainUsername = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("domain_username", d.Get("domain_username")))
	}

	domainBody := map[string]any{
		"name":             name,
		"description":      description,
		"publicZone":       publicZone,
		"visibility":       visibility,
		"active":           active,
		"domainController": domainController,
		"domainUsername":   domainUsername,
	}
	// tenant_id maps to the API "account" association (master tenant only).
	if v, ok := d.GetOk("tenant_id"); ok {
		if tenantID, isInt := v.(int); isInt {
			domainBody["account"] = map[string]any{"id": tenantID}
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("tenant_id", v))
		}
	}
	// domain_password is deprecated in favour of the write-only
	// domain_password_wo; send whichever is set.
	if v, ok := d.GetOk("domain_password"); ok {
		if domainPassword, isString := v.(string); isString {
			domainBody["domainPassword"] = domainPassword
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("domain_password", v))
		}
	}
	if rawConfig := d.GetRawConfig(); !rawConfig.IsNull() && rawConfig.IsKnown() {
		wo := rawConfig.GetAttr("domain_password_wo")
		if !wo.IsNull() && wo.IsKnown() && wo.AsString() != "" {
			domainBody["domainPassword"] = wo.AsString()
		}
	}
	req := &morpheus.Request{
		Body: map[string]any{
			"networkDomain": domainBody,
		},
	}
	resp, err := client.CreateNetworkDomain(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	var result *morpheus.CreateNetworkDomainResult
	if v, ok := resp.Result.(*morpheus.CreateNetworkDomainResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.NetworkDomain == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("NetworkDomain"))
	}

	networkDomain := result.NetworkDomain
	// Successfully created resource, now set id
	d.SetId(convert.Int64ToString(networkDomain.ID))
	diags = append(diags, resourceNetworkDomainRead(ctx, d, meta)...)

	return diags
}

func resourceNetworkDomainRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindNetworkDomainByName(name)
	} else if id != "" {
		resp, err = client.GetNetworkDomain(convert.StringToInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("NetworkDomain cannot be read without name or id")
	}
	if err != nil {
		// 404 is ok?
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

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	// store resource data
	var result *morpheus.GetNetworkDomainResult
	if v, ok := resp.Result.(*morpheus.GetNetworkDomainResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.NetworkDomain == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("NetworkDomain"))
	}

	networkDomain := result.NetworkDomain
	d.SetId(convert.Int64ToString(networkDomain.ID))
	if err := d.Set("name", networkDomain.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", networkDomain.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("active", networkDomain.Active); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("public_zone", networkDomain.PublicZone); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("domain_controller", networkDomain.DomainController); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("visibility", networkDomain.Visibility); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("domain_username", networkDomain.DomainUsername); err != nil {
		return diag.FromErr(err)
	}
	// auto_join_domain has no Morpheus API counterpart; preserve the prior
	// state value so imported state matches applied state.
	if err := d.Set("auto_join_domain", d.Get("auto_join_domain")); err != nil {
		return diag.FromErr(err)
	}
	// tenant_id maps to the API "account" association.
	if networkDomain.Account != nil {
		if err := d.Set("tenant_id", int(networkDomain.Account.ID)); err != nil {
			return diag.FromErr(err)
		}
	}
	// d.Set("fqdn", networkDomain.Fqdn)

	return diags
}

func resourceNetworkDomainUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	var description string
	if v, ok := d.Get("description").(string); ok {
		description = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("description", d.Get("description")))
	}

	var publicZone bool
	if v, ok := d.Get("public_zone").(bool); ok {
		publicZone = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("public_zone", d.Get("public_zone")))
	}

	var visibility string
	if v, ok := d.Get("visibility").(string); ok {
		visibility = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("visibility", d.Get("visibility")))
	}

	var active bool
	if v, ok := d.Get("active").(bool); ok {
		active = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("active", d.Get("active")))
	}

	var domainController bool
	if v, ok := d.Get("domain_controller").(bool); ok {
		domainController = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("domain_controller", d.Get("domain_controller")))
	}

	var domainUsername string
	if v, ok := d.Get("domain_username").(string); ok {
		domainUsername = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("domain_username", d.Get("domain_username")))
	}

	domainBody := map[string]any{
		"name":             name,
		"description":      description,
		"publicZone":       publicZone,
		"visibility":       visibility,
		"active":           active,
		"domainController": domainController,
		"domainUsername":   domainUsername,
	}
	// tenant_id maps to the API "account" association (master tenant only).
	if v, ok := d.GetOk("tenant_id"); ok {
		if tenantID, isInt := v.(int); isInt {
			domainBody["account"] = map[string]any{"id": tenantID}
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("tenant_id", v))
		}
	}
	// domain_password is deprecated in favour of the write-only
	// domain_password_wo; send whichever is configured. Omitting the field
	// clears the stored credential, so a configured value is always sent.
	if domainPassword, ok := d.Get("domain_password").(string); ok && domainPassword != "" {
		domainBody["domainPassword"] = domainPassword
	}
	if rawConfig := d.GetRawConfig(); !rawConfig.IsNull() && rawConfig.IsKnown() {
		wo := rawConfig.GetAttr("domain_password_wo")
		switch {
		case !wo.IsNull() && wo.IsKnown() && wo.AsString() != "":
			domainBody["domainPassword"] = wo.AsString()
		case d.HasChange("domain_password_wo_version"):
			return diag.Errorf("domain_password_wo_version changed but domain_password_wo is not set")
		}
	}
	req := &morpheus.Request{
		Body: map[string]any{
			"networkDomain": domainBody,
		},
	}
	resp, err := client.UpdateNetworkDomain(convert.StringToInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	var result *morpheus.UpdateNetworkDomainResult
	if v, ok := resp.Result.(*morpheus.UpdateNetworkDomainResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.NetworkDomain == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("NetworkDomain"))
	}

	networkDomain := result.NetworkDomain
	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(convert.Int64ToString(networkDomain.ID))

	return resourceNetworkDomainRead(ctx, d, meta)
}

func resourceNetworkDomainDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteNetworkDomain(convert.StringToInt64(id), req)
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
