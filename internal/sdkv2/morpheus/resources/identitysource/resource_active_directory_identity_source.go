package identitysource

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"log"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/helpers"
)

func ResourceActiveDirectoryIdentitySource() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides an active directory identity source resource",
		CreateContext: resourceActiveDirectoryIdentitySourceCreate,
		ReadContext:   resourceActiveDirectoryIdentitySourceRead,
		UpdateContext: resourceActiveDirectoryIdentitySourceUpdate,
		DeleteContext: resourceActiveDirectoryIdentitySourceDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the active directory identity source",
				Computed:    true,
			},
			"tenant_id": {
				Type:        schema.TypeInt,
				Description: "The ID of the Morpheus tenant to associate the identity source with",
				Required:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the active directory identity source",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the active directory identity source",
				Optional:    true,
				Computed:    true,
			},
			"ad_server": {
				Type:        schema.TypeString,
				Description: "The IP address or hostname of the active directory domain controller",
				Required:    true,
			},
			"domain": {
				Type:        schema.TypeString,
				Description: "The name of the active directory domain",
				Required:    true,
			},
			"use_ssl": {
				Type:        schema.TypeBool,
				Description: "Whether to use SSL when connecting to the domain controller",
				Optional:    true,
				Computed:    true,
			},
			"binding_username": {
				Type:        schema.TypeString,
				Description: "The username of the account used to authenticate to the domain",
				Required:    true,
			},
			"binding_password": {
				Type:        schema.TypeString,
				Description: "The password of the account used to authenticate to the domain",
				Required:    true,
				Sensitive:   true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					h := sha256.New()
					h.Write([]byte(new))
					sha256_hash := hex.EncodeToString(h.Sum(nil))
					return strings.ToLower(old) == strings.ToLower(sha256_hash)
				},
			},
			"required_group": {
				Type:        schema.TypeString,
				Description: "The active directory group users must be in to access Morpheus",
				Optional:    true,
				Computed:    true,
			},
			"search_member_groups": {
				Type:        schema.TypeBool,
				Description: "Whether groups nested inside the required group will also be included",
				Optional:    true,
				Computed:    true,
			},
			"default_account_role_id": {
				Type:        schema.TypeInt,
				Description: "The id of the default role a user is assigned when they are in the required group or if no specific group mapping applies to the user",
				Required:    true,
			},
			"enable_role_mapping_permission": {
				Type:        schema.TypeBool,
				Description: "When enabled, Tenant users with appropriate rights to view and edit Roles will have the ability to set role mapping for the Identity Source integration",
				Optional:    true,
				Computed:    true,
			},
			"role_mapping": {
				Description: "The Active Directory to Morpheus Role mapping",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role_id": {
							Description: "The id of the Morpheus role to map to",
							Type:        schema.TypeInt,
							Optional:    true,
						},
						"role_name": {
							Description: "The name or authority of the Morpheus role to map to",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"active_directory_group_name": {
							Description: "The name of the active directory role to map to",
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
						},
						"active_directory_group_fqn": {
							Description: "The fully qualified name of the active directory role to map to (i.e. - CN=Administrators,CN=Builtin,DC=contoso,DC=com)",
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceActiveDirectoryIdentitySourceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	var client *morpheus.Client
	if clientAssert, ok := meta.(*morpheus.Client); ok {
		client = clientAssert
	} else {
		return diag.FromErr(helpers.TypeAssertFail("client", clientAssert))
	}

	identitySource := make(map[string]interface{})

	if name, ok := d.Get("name").(string); ok {
		identitySource["name"] = name
	} else {
		return diag.FromErr(helpers.TypeAssertFail("name", d.Get("name")))
	}

	if description, ok := d.Get("description").(string); ok {
		identitySource["description"] = description
	} else {
		return diag.FromErr(helpers.TypeAssertFail("description", d.Get("description")))
	}

	identitySource["type"] = "activeDirectory"

	config := make(map[string]interface{})

	if adServer, ok := d.Get("ad_server").(string); ok {
		config["url"] = adServer
	} else {
		return diag.FromErr(helpers.TypeAssertFail("ad_server", d.Get("ad_server")))
	}

	if domain, ok := d.Get("domain").(string); ok {
		config["domain"] = domain
	} else {
		return diag.FromErr(helpers.TypeAssertFail("domain", d.Get("domain")))
	}

	if useSSL, ok := d.Get("use_ssl").(bool); ok {
		if useSSL {
			config["useSSL"] = "on"
		} else {
			config["useSSL"] = "off"
		}
	} else {
		return diag.FromErr(helpers.TypeAssertFail("use_ssl", d.Get("use_ssl")))
	}

	if bindingUsername, ok := d.Get("binding_username").(string); ok {
		config["bindingUsername"] = bindingUsername
	} else {
		return diag.FromErr(helpers.TypeAssertFail("binding_username", d.Get("binding_username")))
	}

	if bindingPassword, ok := d.Get("binding_password").(string); ok {
		config["bindingPassword"] = bindingPassword
	} else {
		return diag.FromErr(helpers.TypeAssertFail("binding_password", d.Get("binding_password")))
	}

	if requiredGroup, ok := d.Get("required_group").(string); ok {
		config["requiredGroup"] = requiredGroup
	} else {
		return diag.FromErr(helpers.TypeAssertFail("required_group", d.Get("required_group")))
	}

	if searchMemberGroups, ok := d.Get("search_member_groups").(bool); ok {
		config["searchMemberGroups"] = searchMemberGroups
	} else {
		return diag.FromErr(helpers.TypeAssertFail("search_member_groups", d.Get("search_member_groups")))
	}

	if allowCustomMappings, ok := d.Get("enable_role_mapping_permission").(bool); ok {
		config["allowCustomMappings"] = allowCustomMappings
	} else {
		return diag.FromErr(helpers.TypeAssertFail("enable_role_mapping_permission", d.Get("enable_role_mapping_permission")))
	}

	identitySource["config"] = config

	defaultAccountRole := make(map[string]interface{})

	if defaultAccountRoleID, ok := d.Get("default_account_role_id").(int); ok {
		defaultAccountRole["id"] = defaultAccountRoleID
	} else {
		return diag.FromErr(helpers.TypeAssertFail("default_account_role_id", d.Get("default_account_role_id")))
	}

	identitySource["defaultAccountRole"] = defaultAccountRole

	// Role Mappings
	if roleMappingValue := d.Get("role_mapping"); roleMappingValue != "" {
		if roleMappingSet, ok := roleMappingValue.(*schema.Set); ok {
			identitySource["roleMappings"] = parseRoleMappings(roleMappingSet)
		} else {
			return diag.FromErr(helpers.TypeAssertFail("role_mapping", roleMappingValue))
		}
	}

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"userSource": identitySource,
		},
	}

	if tenantID, ok := d.Get("tenant_id").(int); ok {
		resp, err := client.CreateIdentitySource(int64(tenantID), req)
		if err != nil {
			log.Printf("API FAILURE: %s - %s", resp, err)
			return diag.FromErr(err)
		}

		if result, ok := resp.Result.(*morpheus.CreateIdentitySourceResult); ok {
			// Successfully created resource, now set id
			d.SetId(convert.Int64ToString(result.IdentitySource.ID))
		} else {
			return diag.FromErr(helpers.TypeAssertFail("resp.Result", resp.Result))
		}

	} else {
		return diag.FromErr(helpers.TypeAssertFail("tenant_id", d.Get("tenant_id")))
	}

	resourceActiveDirectoryIdentitySourceRead(ctx, d, meta)
	return diags
}

func resourceActiveDirectoryIdentitySourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	var client *morpheus.Client
	if clientAssert, ok := meta.(*morpheus.Client); ok {
		client = clientAssert
	} else {
		return diag.FromErr(helpers.TypeAssertFail("client", clientAssert))
	}

	id := d.Id()

	var name string
	if nameValue, ok := d.Get("name").(string); ok {
		name = nameValue
	} else {
		return diag.FromErr(helpers.TypeAssertFail("name", d.Get("name")))
	}

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindIdentitySourceByName(name)
	} else if id != "" {
		resp, err = client.GetIdentitySource(convert.StringToInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Identity source cannot be read without name or id")
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

	// store resource data
	result := resp.Result.(*morpheus.GetIdentitySourceResult)
	identitySource := result.IdentitySource
	d.SetId(convert.Int64ToString(identitySource.ID))
	d.Set("name", identitySource.Name)
	d.Set("description", identitySource.Description)
	d.Set("ad_server", identitySource.Config.URL)
	d.Set("domain", identitySource.Config.Domain)
	if identitySource.Config.UseSSL == "off" {
		d.Set("use_ssl", false)
	} else {
		d.Set("use_ssl", true)
	}
	d.Set("binding_username", identitySource.Config.BindingUsername)
	d.Set("binding_password", identitySource.Config.BindingPasswordHash)
	d.Set("required_group", identitySource.Config.RequiredGroup)
	d.Set("search_member_groups", identitySource.Config.SearchMemberGroups)
	d.Set("enable_role_mapping_permission", identitySource.AllowCustomMappings)
	d.Set("default_account_role_id", identitySource.DefaultAccountRole.ID)

	var roleMappingPayload []map[string]interface{}

	for _, roleMapping := range identitySource.RoleMappings {
		roleOutput := make(map[string]interface{})
		roleOutput["active_directory_group_fqn"] = roleMapping.SourceRoleFqn
		roleOutput["active_directory_group_name"] = roleMapping.SourceRoleName
		roleOutput["role_id"] = roleMapping.MappedRole.ID
		roleOutput["role_name"] = roleMapping.MappedRole.Authority
		roleMappingPayload = append(roleMappingPayload, roleOutput)
	}
	d.Set("role_mapping", roleMappingPayload)
	return diags
}

func resourceActiveDirectoryIdentitySourceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var client *morpheus.Client
	if clientAssert, ok := meta.(*morpheus.Client); ok {
		client = clientAssert
	} else {
		return diag.FromErr(helpers.TypeAssertFail("client", clientAssert))
	}

	id := d.Id()

	identitySource := make(map[string]interface{})

	if name, ok := d.Get("name").(string); ok {
		identitySource["name"] = name
	} else {
		return diag.FromErr(helpers.TypeAssertFail("name", d.Get("name")))
	}

	if description, ok := d.Get("description").(string); ok {
		identitySource["description"] = description
	} else {
		return diag.FromErr(helpers.TypeAssertFail("description", d.Get("description")))
	}

	identitySource["type"] = "activeDirectory"

	config := make(map[string]interface{})

	if adServer, ok := d.Get("ad_server").(string); ok {
		config["url"] = adServer
	} else {
		return diag.FromErr(helpers.TypeAssertFail("ad_server", d.Get("ad_server")))
	}

	if domain, ok := d.Get("domain").(string); ok {
		config["domain"] = domain
	} else {
		return diag.FromErr(helpers.TypeAssertFail("domain", d.Get("domain")))
	}

	if useSSL, ok := d.Get("use_ssl").(bool); ok {
		if useSSL {
			config["useSSL"] = "on"
		} else {
			config["useSSL"] = "off"
		}
	} else {
		return diag.FromErr(helpers.TypeAssertFail("use_ssl", d.Get("use_ssl")))
	}

	if bindingUsername, ok := d.Get("binding_username").(string); ok {
		config["bindingUsername"] = bindingUsername
	} else {
		return diag.FromErr(helpers.TypeAssertFail("binding_username", d.Get("binding_username")))
	}

	if d.HasChange("binding_password") {
		if bindingPassword, ok := d.Get("binding_password").(string); ok {
			config["bindingPassword"] = bindingPassword
		} else {
			return diag.FromErr(helpers.TypeAssertFail("binding_password", d.Get("binding_password")))
		}
	}

	if requiredGroup, ok := d.Get("required_group").(string); ok {
		config["requiredGroup"] = requiredGroup
	} else {
		return diag.FromErr(helpers.TypeAssertFail("required_group", d.Get("required_group")))
	}

	if searchMemberGroups, ok := d.Get("search_member_groups").(bool); ok {
		config["searchMemberGroups"] = searchMemberGroups
	} else {
		return diag.FromErr(helpers.TypeAssertFail("search_member_groups", d.Get("search_member_groups")))
	}

	if allowCustomMappings, ok := d.Get("enable_role_mapping_permission").(bool); ok {
		config["allowCustomMappings"] = allowCustomMappings
	} else {
		return diag.FromErr(helpers.TypeAssertFail("enable_role_mapping_permission", d.Get("enable_role_mapping_permission")))
	}

	identitySource["config"] = config

	defaultAccountRole := make(map[string]interface{})

	if defaultAccountRoleID, ok := d.Get("default_account_role_id").(int); ok {
		defaultAccountRole["id"] = defaultAccountRoleID
	} else {
		return diag.FromErr(helpers.TypeAssertFail("default_account_role_id", d.Get("default_account_role_id")))
	}

	identitySource["defaultAccountRole"] = defaultAccountRole

	// Role Mappings
	if roleMappingValue := d.Get("role_mapping"); roleMappingValue != "" {
		if roleMappingSet, ok := roleMappingValue.(*schema.Set); ok {
			identitySource["roleMappings"] = parseRoleMappings(roleMappingSet)
		} else {
			return diag.FromErr(helpers.TypeAssertFail("role_mapping", roleMappingValue))
		}
	}

	req := &morpheus.Request{
		Body: map[string]interface{}{
			"userSource": identitySource,
		},
	}

	resp, err := client.UpdateIdentitySource(convert.StringToInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)
		return diag.FromErr(err)
	}
	result := resp.Result.(*morpheus.UpdateIdentitySourceResult)
	identitySourceResult := result.IdentitySource

	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(convert.Int64ToString(identitySourceResult.ID))
	return resourceActiveDirectoryIdentitySourceRead(ctx, d, meta)
}

func resourceActiveDirectoryIdentitySourceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	var client *morpheus.Client
	if clientAssert, ok := meta.(*morpheus.Client); ok {
		client = clientAssert
	} else {
		return diag.FromErr(helpers.TypeAssertFail("client", clientAssert))
	}

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteIdentitySource(convert.StringToInt64(id), req)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %s", resp, err)
			return diag.FromErr(err)
		} else {
			log.Printf("API FAILURE: %s - %s", resp, err)
			return diag.FromErr(err)
		}
	}
	d.SetId("")
	return diags
}

func parseRoleMappings(mappings *schema.Set) []map[string]interface{} {
	var roleMappings []map[string]interface{}
	// iterate over the array of roleMappings
	for _, mapping := range mappings.List() {
		row := make(map[string]interface{})
		mappedRole := make(map[string]interface{})
		mappingConfig := mapping.(map[string]interface{})
		for k, v := range mappingConfig {
			switch k {
			case "role_id":
				if id, ok := v.(int); ok {
					mappedRole["id"] = id
				}
			case "role_name":
				if authority, ok := v.(string); ok {
					mappedRole["authority"] = authority
				}
			case "active_directory_group_name":
				if sourceRoleName, ok := v.(string); ok {
					row["sourceRoleName"] = sourceRoleName
				}
			case "active_directory_group_fqn":
				if sourceRoleFqn, ok := v.(string); ok {
					row["sourceRoleFqn"] = sourceRoleFqn
				}
			}
		}
		row["mappedRole"] = mappedRole
		roleMappings = append(roleMappings, row)
	}
	return roleMappings
}
