// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package compute

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	morpheus "github.com/HPE/terraform-provider-hpe/internal/sdk/legacy"

	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/convert"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/helpers"
)

func ResourceResourcePoolGroup() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus resource pool group resource",
		CreateContext: resourceResourcePoolGroupCreate,
		ReadContext:   resourceResourcePoolGroupRead,
		UpdateContext: resourceResourcePoolGroupUpdate,
		DeleteContext: resourceResourcePoolGroupDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the resource pool group",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the resource pool group",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the resource pool group",
				Optional:    true,
				Computed:    true,
			},
			"mode": {
				Type:         schema.TypeString,
				Description:  "The load balancing mode of the resource pool group (roundrobin, availablecapacity)",
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"roundrobin", "availablecapacity"}, false),
			},
			"resource_pool_ids": {
				Type:        schema.TypeSet,
				Description: "A list of resource pool ids associated with the resource pool group",
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"all_group_access": {
				Type:        schema.TypeBool,
				Description: "Whether all groups will be granted access to the resource pool group",
				Optional:    true,
			},
			"group_access": {
				Type:        schema.TypeList,
				Description: "A list of Morpheus group configuration to enable group access to the resource pool group",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_id": {
							Type:        schema.TypeInt,
							Description: "The ID of the Morpheus group to grant access to the resource pool group",
							Required:    true,
						},
						"default": {
							Type:        schema.TypeBool,
							Description: "Whether the resource pool group will be a default for the associated group",
							Required:    true,
						},
					},
				},
			},
			"visibility": {
				Type:         schema.TypeString,
				Description:  "Whether the resource pool group is visible in sub-tenants or not",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"private", "public"}, false),
				Default:      "private",
			},
			"tenant_ids": {
				Type: schema.TypeSet,
				Description: "A list of tenant ids granted access to the resource pool group. " +
					"The owning tenant always has access and does not need to be listed.",
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceResourcePoolGroupCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	var diags diag.Diagnostics

	resourcePermissions := make(map[string]any)

	// MORPH-8207: tenants must be sent INSIDE the resourcePoolGroup envelope as
	// [{"id": N}, ...] so the /api/resource-pools/groups controller persists them.
	// A top-level tenantPermissions key is never consumed by the controller.
	tenantsPayload := make([]map[string]any, 0)
	if attr, ok := d.GetOk("tenant_ids"); ok {
		var tenantSet *schema.Set
		if v, ok := attr.(*schema.Set); ok {
			tenantSet = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("tenant_ids", attr))
		}
		for _, s := range tenantSet.List() {
			var tenantID int
			if v, ok := s.(int); ok {
				tenantID = v
			} else {
				return diag.FromErr(helpers.TypeAssertFailError("tenant_id", s))
			}
			tenantsPayload = append(tenantsPayload, map[string]any{"id": tenantID})
		}
	}

	var allGroupAccess bool
	if v, ok := d.Get("all_group_access").(bool); ok {
		allGroupAccess = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("all_group_access", d.Get("all_group_access")))
	}
	resourcePermissions["all"] = allGroupAccess

	// Group Access
	if d.Get("group_access") != "" {
		var groupAccess []any
		if v, ok := d.Get("group_access").([]any); ok {
			groupAccess = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("group_access", d.Get("group_access")))
		}
		resourcePermissions["sites"] = parseGroupAccess(groupAccess)
	}

	poolsPayload := make([]int, 0)
	if attr, ok := d.GetOk("resource_pool_ids"); ok {
		var poolSet *schema.Set
		if v, ok := attr.(*schema.Set); ok {
			poolSet = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("resource_pool_ids", attr))
		}
		for _, s := range poolSet.List() {
			var poolID int
			if v, ok := s.(int); ok {
				poolID = v
			} else {
				return diag.FromErr(helpers.TypeAssertFailError("resource_pool_ids element", s))
			}
			poolsPayload = append(poolsPayload, poolID)
		}
	}

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

	var mode string
	if v, ok := d.Get("mode").(string); ok {
		mode = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("mode", d.Get("mode")))
	}

	var visibility string
	if v, ok := d.Get("visibility").(string); ok {
		visibility = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("visibility", d.Get("visibility")))
	}

	req := &morpheus.Request{
		Body: map[string]any{
			"resourcePoolGroup": map[string]any{
				"name":        name,
				"description": description,
				"mode":        mode,
				"visibility":  visibility,
				"pools":       poolsPayload,
				// MORPH-8207: tenants live inside the envelope, not top-level.
				"tenants": tenantsPayload,
			},
			"resourcePermissions": resourcePermissions,
		},
	}
	resp, err := client.CreateResourcePoolGroup(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	var result *morpheus.CreateResourcePoolGroupResult
	if v, ok := resp.Result.(*morpheus.CreateResourcePoolGroupResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.ResourcePoolGroup == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("ResourcePoolGroup"))
	}

	resourceResourcePoolGroup := result.ResourcePoolGroup
	d.SetId(convert.Int64ToString(resourceResourcePoolGroup.ID))

	diags = append(diags, resourceResourcePoolGroupRead(ctx, d, meta)...)

	return diags
}

func resourceResourcePoolGroupRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
		resp, err = client.FindResourcePoolGroupByName(name)
	} else if id != "" {
		resp, err = client.GetResourcePoolGroup(convert.StringToInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Resource pool group cannot be read without name or id")
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

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	// store resource data
	var result *morpheus.GetResourcePoolGroupResult
	if v, ok := resp.Result.(*morpheus.GetResourcePoolGroupResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.ResourcePoolGroup == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("ResourcePoolGroup"))
	}

	resourceResourcePoolGroup := result.ResourcePoolGroup
	d.SetId(convert.Int64ToString(resourceResourcePoolGroup.ID))
	d.Set("name", resourceResourcePoolGroup.Name)
	d.Set("description", resourceResourcePoolGroup.Description)
	d.Set("mode", resourceResourcePoolGroup.Mode)

	var resourcePools []int64
	if len(resourceResourcePoolGroup.Pools) > 0 {
		resourcePools = append(resourcePools, resourceResourcePoolGroup.Pools...)
	}
	d.Set("resource_pool_ids", resourcePools)

	d.Set("all_group_access", resourceResourcePoolGroup.ResourcePermission.All)

	// Group Access
	var groupAccess []map[string]any
	if len(resourceResourcePoolGroup.ResourcePermission.Sites) != 0 {
		for _, group := range resourceResourcePoolGroup.ResourcePermission.Sites {
			groupData := make(map[string]any)
			groupData["group_id"] = group.ID
			groupData["default"] = group.Default
			groupAccess = append(groupAccess, groupData)
		}
	}
	d.Set("group_access", groupAccess)
	d.Set("visibility", resourceResourcePoolGroup.Visibility)

	// tenant ids
	//
	// MORPH-8207: the API always lists the owning tenant in `tenants`. Saving
	// group access creates a ResourcePermission row for the owner, and the read
	// serializer returns every row for the group without flagging which one is
	// the owner. Writing that list straight into state produced a perpetual
	// plan whenever the config did not list the owner (including when
	// tenant_ids was not set at all), so the owner is reconciled out here.
	apiTenantIDs := make([]int64, 0, len(resourceResourcePoolGroup.Tenants))
	for _, tenant := range resourceResourcePoolGroup.Tenants {
		apiTenantIDs = append(apiTenantIDs, tenant.ID)
	}

	ownerID, err := resourcePoolGroupOwnerID(client)
	if err != nil {
		// Read managed without whoami before this reconciliation existed, so a
		// failed lookup must not take refresh, plan or import down with it.
		// Keep the API's list, as Read always did, and say why a plan may show
		// the owner being removed from tenant_ids.
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Could not determine the owning tenant of the resource pool group",
			Detail: fmt.Sprintf("Resource pool group %s: %s. tenant_ids was read as the full list "+
				"Morpheus reports, which includes the owning tenant, so a plan may show it being "+
				"removed until the lookup succeeds.", id, err),
		})
		d.Set("tenant_ids", apiTenantIDs)

		return diags
	}

	var priorTenantIDs []int
	if attr, ok := d.GetOk("tenant_ids"); ok {
		var priorSet *schema.Set
		if v, ok := attr.(*schema.Set); ok {
			priorSet = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("tenant_ids", attr))
		}
		for _, s := range priorSet.List() {
			var tenantID int
			if v, ok := s.(int); ok {
				tenantID = v
			} else {
				return diag.FromErr(helpers.TypeAssertFailError("tenant_id", s))
			}
			priorTenantIDs = append(priorTenantIDs, tenantID)
		}
	}
	d.Set("tenant_ids", reconcileTenantIDs(priorTenantIDs, apiTenantIDs, ownerID))

	return diags
}

// resourcePoolGroupOwnerID resolves the tenant that owns the caller's resource
// pool groups: the caller's own tenant, via whoami. A variable so tests can
// drive Read's owner-lookup fallback without a whoami endpoint.
var resourcePoolGroupOwnerID = helpers.CallerAccountID

// reconcileTenantIDs returns the tenant_ids value to store in state given the
// prior (configured or previously stored) IDs, the IDs the API reported, and
// the ID of the tenant that owns the resource pool group.
//
// The API injects the owner into every response, so the owner is dropped
// unless the user listed it explicitly, in which case it is kept so the state
// matches the config. Every other reported tenant is kept, and every tenant
// the API no longer reports is dropped, so changes made outside Terraform
// still surface as drift. With no prior value (import) the result is the
// reported tenants other than the owner.
func reconcileTenantIDs(prior []int, api []int64, ownerID int64) []int {
	keepOwner := false
	for _, id := range prior {
		if int64(id) == ownerID {
			keepOwner = true

			break
		}
	}

	out := make([]int, 0, len(api))
	for _, id := range api {
		if id == ownerID && !keepOwner {
			continue
		}
		out = append(out, int(id))
	}

	return out
}

func resourceResourcePoolGroupUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	id := d.Id()

	resourcePermissions := make(map[string]any)

	// MORPH-8207: tenants must be sent INSIDE the resourcePoolGroup envelope as
	// [{"id": N}, ...] on update as well; top-level tenantPermissions is ignored.
	tenantsPayload := make([]map[string]any, 0)
	if attr, ok := d.GetOk("tenant_ids"); ok {
		var tenantSet *schema.Set
		if v, ok := attr.(*schema.Set); ok {
			tenantSet = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("tenant_ids", attr))
		}
		for _, s := range tenantSet.List() {
			var tenantID int
			if v, ok := s.(int); ok {
				tenantID = v
			} else {
				return diag.FromErr(helpers.TypeAssertFailError("tenant_id", s))
			}
			tenantsPayload = append(tenantsPayload, map[string]any{"id": tenantID})
		}
	}

	var allGroupAccess bool
	if v, ok := d.Get("all_group_access").(bool); ok {
		allGroupAccess = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("all_group_access", d.Get("all_group_access")))
	}
	resourcePermissions["all"] = allGroupAccess

	// Group Access
	if d.Get("group_access") != "" {
		var groupAccess []any
		if v, ok := d.Get("group_access").([]any); ok {
			groupAccess = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("group_access", d.Get("group_access")))
		}
		resourcePermissions["sites"] = parseGroupAccess(groupAccess)
	}

	poolsPayload := make([]int, 0)
	if attr, ok := d.GetOk("resource_pool_ids"); ok {
		var poolSet *schema.Set
		if v, ok := attr.(*schema.Set); ok {
			poolSet = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("resource_pool_ids", attr))
		}
		for _, s := range poolSet.List() {
			var poolID int
			if v, ok := s.(int); ok {
				poolID = v
			} else {
				return diag.FromErr(helpers.TypeAssertFailError("resource_pool_ids element", s))
			}
			poolsPayload = append(poolsPayload, poolID)
		}
	}

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

	var mode string
	if v, ok := d.Get("mode").(string); ok {
		mode = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("mode", d.Get("mode")))
	}

	var visibility string
	if v, ok := d.Get("visibility").(string); ok {
		visibility = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("visibility", d.Get("visibility")))
	}

	req := &morpheus.Request{
		Body: map[string]any{
			"resourcePoolGroup": map[string]any{
				"name":        name,
				"description": description,
				"mode":        mode,
				"visibility":  visibility,
				"pools":       poolsPayload,
				// MORPH-8207: tenants live inside the envelope, not top-level.
				"tenants": tenantsPayload,
			},
			"resourcePermissions": resourcePermissions,
		},
	}
	resp, err := client.UpdateResourcePoolGroup(convert.StringToInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	var result *morpheus.UpdateResourcePoolGroupResult
	if v, ok := resp.Result.(*morpheus.UpdateResourcePoolGroupResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.ResourcePoolGroup == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("ResourcePoolGroup"))
	}

	resourcePoolGroup := result.ResourcePoolGroup
	d.SetId(convert.Int64ToString(resourcePoolGroup.ID))

	return resourceResourcePoolGroupRead(ctx, d, meta)
}

func resourceResourcePoolGroupDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteResourcePoolGroup(convert.StringToInt64(id), req)
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

func parseGroupAccess(variables []any) []map[string]any {
	var accessData []map[string]any
	// iterate over the array of group access
	for i := 0; i < len(variables); i++ {
		row := make(map[string]any)
		groupconfig := variables[i].(map[string]any)
		for k, v := range groupconfig {
			switch k {
			case "group_id":
				row["id"] = v.(int)
			case "default":
				row["default"] = v.(bool)
			}
		}
		accessData = append(accessData, row)
	}

	return accessData
}
