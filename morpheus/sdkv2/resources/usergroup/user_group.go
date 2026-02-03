// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package usergroup

import (
	"context"
	"log"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/convert"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/helpers"
)

func ResourceUserGroup() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus user group resource",
		CreateContext: resourceUserGroupCreate,
		ReadContext:   resourceUserGroupRead,
		UpdateContext: resourceUserGroupUpdate,
		DeleteContext: resourceUserGroupDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the user group",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the user group",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the user group",
				Optional:    true,
				Computed:    true,
			},
			"server_group": {
				Type:        schema.TypeString,
				Description: "The name of the Linux group to add the users to",
				Optional:    true,
				Computed:    true,
			},
			"sudo_access": {
				Type:        schema.TypeBool,
				Description: "Whether the users in the group are granted sudo permissions",
				Optional:    true,
				Computed:    true,
			},
			"user_ids": {
				Type:        schema.TypeList,
				Description: "A list of Morpheus user IDs to add to the user group",
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceUserGroupCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	userGroup := make(map[string]any)

	var name string
	if v, ok := d.Get("name").(string); ok {
		name = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
	}
	userGroup["name"] = name

	var description string
	if v, ok := d.Get("description").(string); ok {
		description = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("description", d.Get("description")))
	}
	userGroup["description"] = description

	var sudoAccess bool
	if v, ok := d.Get("sudo_access").(bool); ok {
		sudoAccess = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("sudo_access", d.Get("sudo_access")))
	}
	userGroup["sudoUser"] = sudoAccess

	var serverGroup string
	if v, ok := d.Get("server_group").(string); ok {
		serverGroup = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("server_group", d.Get("server_group")))
	}
	userGroup["serverGroup"] = serverGroup

	userGroup["users"] = d.Get("user_ids")

	req := &morpheus.Request{
		Body: map[string]any{
			"userGroup": userGroup,
		},
	}
	resp, err := client.CreateUserGroup(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	var result *morpheus.CreateUserGroupResult
	if v, ok := resp.Result.(*morpheus.CreateUserGroupResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("resp.Result", resp.Result))
	}

	if result.UserGroup == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("UserGroup"))
	}
	userGroupResult := result.UserGroup

	d.SetId(convert.Int64ToString(userGroupResult.ID))

	return resourceUserGroupRead(ctx, d, meta)
}

func resourceUserGroupRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
		resp, err = client.FindUserGroupByName(name)
	} else if id != "" {
		resp, err = client.GetUserGroup(convert.StringToInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("User Group cannot be read without name or id")
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

	var result *morpheus.GetUserGroupResult
	if v, ok := resp.Result.(*morpheus.GetUserGroupResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("resp.Result", resp.Result))
	}

	if result.UserGroup == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("UserGroup"))
	}
	userGroup := result.UserGroup

	d.SetId(convert.IntToString(int(userGroup.ID)))
	d.Set("name", userGroup.Name)
	d.Set("description", userGroup.Description)
	d.Set("server_group", userGroup.ServerGroup)
	d.Set("sudo_access", userGroup.SudoUser)

	var users []int64
	if userGroup.Users != nil {
		for i := 0; i < len(userGroup.Users); i++ {
			users = append(users, userGroup.Users[i].ID)
		}
	}

	var declaredUserIDs []any
	if v, ok := d.Get("user_ids").([]any); ok {
		declaredUserIDs = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("user_ids", d.Get("user_ids")))
	}

	userIDs := matchUserIDsWithSchema(users, declaredUserIDs)
	d.Set("user_ids", userIDs)

	return diags
}

func resourceUserGroupUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	id := d.Id()

	userGroup := make(map[string]any)

	var name string
	if v, ok := d.Get("name").(string); ok {
		name = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
	}
	userGroup["name"] = name

	var description string
	if v, ok := d.Get("description").(string); ok {
		description = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("description", d.Get("description")))
	}
	userGroup["description"] = description

	var sudoAccess bool
	if v, ok := d.Get("sudo_access").(bool); ok {
		sudoAccess = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("sudo_access", d.Get("sudo_access")))
	}
	userGroup["sudoUser"] = sudoAccess

	var serverGroup string
	if v, ok := d.Get("server_group").(string); ok {
		serverGroup = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("server_group", d.Get("server_group")))
	}
	userGroup["serverGroup"] = serverGroup

	userGroup["users"] = d.Get("user_ids")

	req := &morpheus.Request{
		Body: map[string]any{
			"userGroup": userGroup,
		},
	}
	resp, err := client.UpdateUserGroup(convert.StringToInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	var result *morpheus.UpdateUserGroupResult
	if v, ok := resp.Result.(*morpheus.UpdateUserGroupResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("resp.Result", resp.Result))
	}

	if result.UserGroup == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("UserGroup"))
	}
	userGroupResult := result.UserGroup

	d.SetId(convert.Int64ToString(userGroupResult.ID))

	return resourceUserGroupRead(ctx, d, meta)
}

func resourceUserGroupDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteUserGroup(convert.StringToInt64(id), req)
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

// This cannot currently be handled efficiently by a DiffSuppressFunc.
// See: https://github.com/hashicorp/terraform-plugin-sdk/issues/477
func matchUserIDsWithSchema(userIDs []int64, declaredUserIDs []any) []int64 {
	if declaredUserIDs == nil {
		return userIDs
	}

	result := make([]int64, len(declaredUserIDs))

	rMap := make(map[int64]int64, len(userIDs))
	for _, userID := range userIDs {
		rMap[userID] = userID
	}

	for i, definedUserID := range declaredUserIDs {
		definedUserID := int64(definedUserID.(int))

		if v, ok := rMap[definedUserID]; ok {
			result[i] = v
			delete(rMap, v)
		}
	}

	for _, rcpt := range rMap {
		result = append(result, rcpt)
	}

	return result
}
