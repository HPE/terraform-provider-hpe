// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package cloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/convert"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/helpers"

	morpheus "github.com/HPE/terraform-provider-hpe/internal/sdk/legacy"
)

func DataSourceCloudType() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus cloud type data source. It reads " +
			"`/api/zone-types`, which requires no appliance-level permission and " +
			"so is reachable by sub-tenant callers. Only enabled cloud types are " +
			"returned (the endpoint defaults to enabled cloud types).",
		ReadContext: dataSourceCloudTypeRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the Morpheus cloud type",
				Required:    true,
			},
		},
	}
}

func dataSourceCloudTypeRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	// Use the tenant-reachable /api/zone-types endpoint (via ListCloudTypes)
	// instead of the appliance-scoped /api/appliance-settings/zone-types, which
	// is gated on the admin-appliance permission and 403s for sub-tenants
	// (MORPH-16399). The server "name" filter matches on code OR name, so we
	// still confirm an exact Name match locally.
	resp, err := client.ListCloudTypes(&morpheus.Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %v", resp, err)

			return diag.Errorf("cloud type %q not found", name)
		}

		log.Printf("API FAILURE: %s - %v", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	var result *morpheus.ListCloudTypesResult
	if v, ok := resp.Result.(*morpheus.ListCloudTypesResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.CloudTypes != nil {
		for _, cType := range *result.CloudTypes {
			if cType.Name == name {
				d.SetId(convert.Int64ToString(cType.ID))
				d.Set("name", cType.Name)

				return diags
			}
		}
	}

	return diag.Errorf("cloud type %q not found", name)
}
