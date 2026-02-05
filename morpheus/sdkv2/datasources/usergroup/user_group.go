// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package usergroup

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"

	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/convert"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/helpers"
)

func DataSourceUserGroup() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus user group data source.",
		ReadContext: dataSourceUserGroupRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeInt,
				Description:   "The ID of the user group",
				Optional:      true,
				ConflictsWith: []string{"name"},
				Computed:      true,
			},
			"name": {
				Type:          schema.TypeString,
				Description:   "The name of the user group",
				Optional:      true,
				ConflictsWith: []string{"id"},
			},
		},
	}
}

func dataSourceUserGroupRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	var id int
	if v, ok := d.Get("id").(int); ok {
		id = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("id", d.Get("id")))
	}

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == 0 && name != "" {
		resp, err = client.FindUserGroupByName(name)
	} else if id != 0 {
		resp, err = client.GetUserGroup(int64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("User Group cannot be read without name or id")
	}
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %v", resp, err)

			return nil
		}

		log.Printf("API FAILURE: %s - %v", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	// store resource data
	var result *morpheus.GetUserGroupResult
	if v, ok := resp.Result.(*morpheus.GetUserGroupResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("GetUserGroupResult", resp.Result))
	}

	if result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("GetUserGroupResult"))
	}

	userGroup := result.UserGroup
	if userGroup == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("UserGroup"))
	}

	d.SetId(convert.Int64ToString(userGroup.ID))
	d.Set("name", userGroup.Name)

	return diags
}
