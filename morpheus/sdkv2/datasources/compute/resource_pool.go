// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package compute

import (
	"context"
	"errors"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	morpheus "github.com/HPE/terraform-provider-hpe/internal/sdk/legacy"

	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/convert"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/helpers"
)

func DataSourceResourcePool() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus resource pool data source.\n\n" +
			"Looks a pool up by `id` or by `name` within a cloud. A lookup by name uses the cloud's " +
			"resource-pool listing, which returns cloud-level pools only: the pool Morpheus creates for " +
			"an HVM cluster is attached to the cluster rather than to the cloud and is not listed, so it " +
			"must be looked up here by `id`, or read from the `hpe_morpheus_cluster` data source as " +
			"`permissions.resource_pool.id`.",
		ReadContext: dataSourceResourcePoolRead,
		Schema: map[string]*schema.Schema{
			"cloud_id": {
				Type:        schema.TypeInt,
				Description: "The id of the Morpheus cloud to search for the resource pool.",
				Required:    true,
			},
			"name": {
				Type: schema.TypeString,
				Description: "The name of the Morpheus resource pool. Matches cloud-level pools only; " +
					"an HVM cluster's own pool is not found by name (see the data source description).",
				Optional: true,
			},
			"type": {
				Type: schema.TypeString,
				Description: "The kind of pool, as reported by Morpheus. `default` for a pool created in the " +
					"cloud; `namespace` for the pool Morpheus creates for an HVM cluster; `cluster` for a " +
					"vSphere cluster; `vpc` for an AWS VPC; `resourceGroup` for an Azure resource group.",
				Computed: true,
			},
			"active": {
				Type:        schema.TypeBool,
				Description: "Whether the resource pool is enabled or not",
				Computed:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the resource pool",
				Computed:    true,
			},
			"id": {
				Type:        schema.TypeInt,
				Description: "The id of the resource pool",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func dataSourceResourcePoolRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	var id int
	if v, ok := d.Get("id").(int); ok {
		id = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("id", d.Get("id")))
	}

	var name string
	if v, ok := d.Get("name").(string); ok {
		name = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
	}

	var cloudID int
	if v, ok := d.Get("cloud_id").(int); ok {
		cloudID = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("cloud_id", d.Get("cloud_id")))
	}

	// Ensure that either the id or name is provided
	if id == 0 && name == "" {
		return diag.Errorf(
			"Either 'id' or 'name' must be provided to search for the resource pool",
		)
	}

	var resp *morpheus.Response
	var err error

	if id != 0 {
		resp, err = client.GetResourcePool(
			int64(cloudID),
			int64(id),
			&morpheus.Request{},
		)
		if err != nil && resp != nil && resp.StatusCode == 404 {
			return diag.Errorf("%s", resourcePoolNotFoundByIDMessage(id, cloudID))
		}
	} else {
		resp, err = findResourcePoolByName(client, name, cloudID)
	}

	if err != nil {
		log.Printf("API FAILURE: %s - %v", resp, err)

		return diag.FromErr(err)
	}

	log.Printf("API RESPONSE: %s", resp)

	var result *morpheus.GetResourcePoolResult
	if v, ok := resp.Result.(*morpheus.GetResourcePoolResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.ResourcePool == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("ResourcePool"))
	}

	resourcePool := result.ResourcePool

	d.SetId(convert.Int64ToString(resourcePool.ID))
	d.Set("name", resourcePool.Name)
	d.Set("active", resourcePool.Active)
	d.Set("type", resourcePool.Type)
	d.Set("description", resourcePool.Description)

	return diags
}

// findResourcePoolByName lists the cloud's pools filtered by name, requires
// exactly one exact match, and fetches it by id. It owns the three outcomes so
// that a miss explains the listing's scope instead of just reporting a count —
// and, when a cluster of that name exists in the cloud, gives the id of the
// cluster's pool, which is what the user was almost certainly looking for.
func findResourcePoolByName(client *morpheus.Client, name string, cloudID int) (*morpheus.Response, error) {
	resp, err := client.ListResourcePools(int64(cloudID), &morpheus.Request{
		QueryParams: map[string]string{"name": name},
	})
	if err != nil {
		return resp, err
	}

	list, ok := resp.Result.(*morpheus.ListResourcePoolsResult)
	if !ok {
		return resp, helpers.TypeAssertFailError("Result", resp.Result)
	}

	var matchIDs []int64

	if list.ResourcePools != nil {
		for _, pool := range *list.ResourcePools {
			if pool.Name == name {
				matchIDs = append(matchIDs, pool.ID)
			}
		}
	}

	switch len(matchIDs) {
	case 1:
		return client.GetResourcePool(int64(cloudID), matchIDs[0], &morpheus.Request{})
	case 0:
		var hint string
		if poolID, found := findClusterResourcePoolID(client, name, cloudID); found {
			hint = clusterResourcePoolHint(name, cloudID, poolID)
		}

		return resp, errors.New(resourcePoolNotFoundByNameMessage(name, cloudID, hint))
	default:
		return resp, errors.New(resourcePoolMultipleMatchesMessage(name, cloudID, matchIDs))
	}
}
