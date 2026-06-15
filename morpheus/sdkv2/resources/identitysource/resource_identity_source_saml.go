package identitysource

import (
	"context"
	"encoding/json"
	"log"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/convert"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/helpers"
)

//nolint:lll
func ResourceIdentitySourceSAML() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a saml identity source resource",
		CreateContext: resourceIdentitySourceSAMLCreate,
		ReadContext:   resourceIdentitySourceSAMLRead,
		UpdateContext: resourceIdentitySourceSAMLUpdate,
		DeleteContext: resourceIdentitySourceSAMLDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the SAML identity source",
				Computed:    true,
			},
			"tenant_id": {
				Type:        schema.TypeInt,
				Description: "The ID of the Morpheus tenant to associate the identity source with",
				Required:    true,
				ForceNew:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the SAML identity source",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the SAML identity source",
				Optional:    true,
				Computed:    true,
			},
			"login_redirect_url": {
				Type:        schema.TypeString,
				Description: "This is the SAML endpoint Morpheus will redirect to when a user signs into Morpheus via SAML",
				Optional:    true,
				Computed:    true,
			},
			"logout_redirect_url": {
				Type:        schema.TypeString,
				Description: "The URL Morpheus will POST to when a SAML user logs out of Morpheus",
				Optional:    true,
				Computed:    true,
			},
			"include_saml_request_parameter": {
				Type:        schema.TypeBool,
				Description: "Whether to include the SAML request as a parameter",
				Optional:    true,
				Computed:    true,
			},
			"saml_request": {
				Type:         schema.TypeString,
				Description:  "The SAML request configuration (NoSignature, SelfSigned, CustomSignature)",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"NoSignature", "SelfSigned", "CustomSignature"}, false),
			},
			"validate_assertion_signature": {
				Type:        schema.TypeBool,
				Description: "Whether to validate the assertion signature (SAML RESPONSE field in the UI)",
				Optional:    true,
				Computed:    true,
			},
			"given_name_attribute": {
				Type:        schema.TypeString,
				Description: "SAML SP field value to map to Morpheus user First Name",
				Optional:    true,
				Computed:    true,
			},
			"surname_attribute": {
				Type:        schema.TypeString,
				Description: "SAML SP field value to map to Morpheus user Last Name",
				Optional:    true,
				Computed:    true,
			},
			"email_attribute": {
				Type:        schema.TypeString,
				Description: "SAML SP field value to map to Morpheus user email address",
				Optional:    true,
				Computed:    true,
			},
			"default_account_role_id": {
				Type:        schema.TypeInt,
				Description: "The id of the default role a user is assigned when they are in the required group or if no specific group mapping applies to the user",
				Required:    true,
			},
			"role_attribute_name": {
				Type:        schema.TypeString,
				Description: "The name of the attribute/assertion field that will map to Morpheus roles, such a MemberOf",
				Optional:    true,
				Computed:    true,
			},
			"required_role_attribute_value": {
				Type:        schema.TypeString,
				Description: "The name of the attribute/assertion field that maps to the required role",
				Optional:    true,
				Computed:    true,
			},
			"role_mapping": {
				Description: "The SAML to Morpheus Role mapping",
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
						"assertion_attribute": {
							Description: "The assertion attribute to map the role to",
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
			"enable_role_mapping_permission": {
				Type:        schema.TypeBool,
				Description: "When enabled, Tenant users with appropriate rights to view and edit Roles will have the ability to set role mapping for the Identity Source integration",
				Optional:    true,
				Computed:    true,
			},
			"entity_id": {
				Type:        schema.TypeString,
				Description: "The SAML Service Provider entity ID (audience) generated by Morpheus, to register with the IdP",
				Computed:    true,
			},
			"acs_url": {
				Type:        schema.TypeString,
				Description: "The SAML Assertion Consumer Service (ACS) URL generated by Morpheus, to register with the IdP",
				Computed:    true,
			},
			"sp_metadata": {
				Type:        schema.TypeString,
				Description: "The SAML Service Provider metadata XML document generated by Morpheus",
				Computed:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceIdentitySourceSAMLCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	var client *morpheus.Client
	if clientAssert, ok := meta.(*morpheus.Client); ok {
		client = clientAssert
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", clientAssert))
	}

	identitySource := make(map[string]any)

	if name, ok := d.Get("name").(string); ok {
		identitySource["name"] = name
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
	}

	if description, ok := d.Get("description").(string); ok {
		identitySource["description"] = description
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("description", d.Get("description")))
	}

	identitySource["type"] = "saml"

	config := make(map[string]any)

	if loginRedirectURL, ok := d.Get("login_redirect_url").(string); ok {
		config["url"] = loginRedirectURL
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("login_redirect_url", d.Get("login_redirect_url")))
	}

	if logoutRedirectURL, ok := d.Get("logout_redirect_url").(string); ok {
		config["logoutUrl"] = logoutRedirectURL
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("logout_redirect_url", d.Get("logout_redirect_url")))
	}

	if includeSAMLRequest, ok := d.Get("include_saml_request_parameter").(bool); ok {
		if includeSAMLRequest {
			config["doNotIncludeSAMLRequest"] = false
		} else {
			config["doNotIncludeSAMLRequest"] = true
		}
	} else {
		return diag.FromErr(
			helpers.TypeAssertFailError(
				"include_saml_request_parameter",
				d.Get("include_saml_request_parameter"),
			),
		)
	}

	if samlRequest, ok := d.Get("saml_request").(string); ok {
		config["SAMLSignatureMode"] = samlRequest
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("saml_request", d.Get("saml_request")))
	}

	if givenNameAttribute, ok := d.Get("given_name_attribute").(string); ok {
		config["givenNameAttribute"] = givenNameAttribute
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("given_name_attribute", d.Get("given_name_attribute")))
	}

	if surnameAttribute, ok := d.Get("surname_attribute").(string); ok {
		config["surnameAttribute"] = surnameAttribute
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("surname_attribute", d.Get("surname_attribute")))
	}

	if emailAttribute, ok := d.Get("email_attribute").(string); ok {
		config["emailAttribute"] = emailAttribute
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("email_attribute", d.Get("email_attribute")))
	}

	if roleAttributeName, ok := d.Get("role_attribute_name").(string); ok {
		config["roleAttributeName"] = roleAttributeName
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("role_attribute_name", d.Get("role_attribute_name")))
	}

	if requiredAttributeValue, ok := d.Get("required_role_attribute_value").(string); ok {
		config["requiredAttributeValue"] = requiredAttributeValue
	} else {
		return diag.FromErr(
			helpers.TypeAssertFailError(
				"required_role_attribute_value",
				d.Get("required_role_attribute_value"),
			),
		)
	}

	if validateAssertionSignature, ok := d.Get("validate_assertion_signature").(bool); ok {
		if validateAssertionSignature {
			config["doNotValidateSignature"] = false
		} else {
			config["doNotValidateSignature"] = true
		}
	} else {
		return diag.FromErr(
			helpers.TypeAssertFailError(
				"validate_assertion_signature",
				d.Get("validate_assertion_signature"),
			),
		)
	}

	identitySource["config"] = config

	defaultAccountRole := make(map[string]any)

	if defaultAccountRoleID, ok := d.Get("default_account_role_id").(int); ok {
		defaultAccountRole["id"] = defaultAccountRoleID
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("default_account_role_id", d.Get("default_account_role_id")))
	}

	identitySource["defaultAccountRole"] = defaultAccountRole

	// Role Mappings
	if roleMappingValue := d.Get("role_mapping"); roleMappingValue != "" {
		if roleMappingSet, ok := roleMappingValue.(*schema.Set); ok {
			identitySource["roleMappings"] = parseSAMLRoleMappings(roleMappingSet)
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("role_mapping", roleMappingValue))
		}
	}

	if enableRoleMappingPermission, ok := d.Get("enable_role_mapping_permission").(bool); ok {
		identitySource["allowCustomMappings"] = enableRoleMappingPermission
	} else {
		return diag.FromErr(
			helpers.TypeAssertFailError(
				"enable_role_mapping_permission",
				d.Get("enable_role_mapping_permission"),
			),
		)
	}

	req := &morpheus.Request{
		Body: map[string]any{
			"userSource": identitySource,
		},
	}

	if tenantID, ok := d.Get("tenant_id").(int); ok {
		resp, err := client.CreateIdentitySource(int64(tenantID), req)
		if err != nil {
			log.Printf("API FAILURE: %s - %s", resp, err)

			return diag.FromErr(err)
		}
		log.Printf("API RESPONSE: %s", resp)

		if result, ok := resp.Result.(*morpheus.CreateIdentitySourceResult); ok {
			// Successfully created resource, now set id
			d.SetId(convert.Int64ToString(result.IdentitySource.ID))
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("resp.Result", resp.Result))
		}
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("tenant_id", d.Get("tenant_id")))
	}

	diags = append(diags, resourceIdentitySourceSAMLRead(ctx, d, meta)...)

	return diags
}

func resourceIdentitySourceSAMLRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	var client *morpheus.Client
	if clientAssert, ok := meta.(*morpheus.Client); ok {
		client = clientAssert
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", clientAssert))
	}

	id := d.Id()

	var name string
	if nameValue, ok := d.Get("name").(string); ok {
		name = nameValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
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
	log.Printf("API RESPONSE: %s", resp)

	// store resource data
	var result *morpheus.GetIdentitySourceResult
	if resultAssert, ok := resp.Result.(*morpheus.GetIdentitySourceResult); ok {
		result = resultAssert
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("resp.Result", resp.Result))
	}

	identitySource := result.IdentitySource
	if identitySource == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("IdentitySource"))
	}

	d.SetId(convert.Int64ToString(identitySource.ID))
	d.Set("name", identitySource.Name)
	d.Set("description", identitySource.Description)
	d.Set("login_redirect_url", identitySource.Config.URL)
	d.Set("logout_redirect_url", identitySource.Config.LogoutURL)
	if identitySource.Config.DoNotIncludeSAMLRequest {
		d.Set("include_saml_request_parameter", false)
	} else {
		d.Set("include_saml_request_parameter", true)
	}
	d.Set("saml_request", identitySource.Config.SAMLSignatureMode)
	if identitySource.Config.DoNotValidateSignature {
		d.Set("validate_assertion_signature", false)
	} else {
		d.Set("validate_assertion_signature", true)
	}
	d.Set("given_name_attribute", identitySource.Config.GivenNameAttribute)
	d.Set("surname_attribute", identitySource.Config.SurnameAttribute)
	d.Set("email_attribute", identitySource.Config.EmailAttribute)
	d.Set("default_account_role_id", identitySource.DefaultAccountRole.ID)
	d.Set("role_attribute_name", identitySource.Config.RoleAttributeName)
	d.Set("required_role_attribute_value", identitySource.Config.RequiredAttributeValue)
	d.Set("enable_role_mapping_permission", identitySource.AllowCustomMappings)

	var roleMappingPayload []map[string]any

	for _, roleMapping := range identitySource.RoleMappings {
		roleOutput := make(map[string]any)
		roleOutput["assertion_attribute"] = roleMapping.SourceRoleName
		roleOutput["role_id"] = roleMapping.MappedRole.ID
		roleOutput["role_name"] = roleMapping.MappedRole.Authority
		roleMappingPayload = append(roleMappingPayload, roleOutput)
	}
	d.Set("role_mapping", roleMappingPayload)

	// entity_id, acs_url and sp_metadata are computed by Morpheus and returned
	// under userSource.providerSettings. The legacy SDK model does not expose
	// them, so read them from the raw response body.
	entityID, acsURL, spMetadata, err := samlProviderSettings(resp.Body)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("entity_id", entityID)
	d.Set("acs_url", acsURL)
	d.Set("sp_metadata", spMetadata)

	return diags
}

func resourceIdentitySourceSAMLUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if clientAssert, ok := meta.(*morpheus.Client); ok {
		client = clientAssert
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", clientAssert))
	}

	id := d.Id()

	identitySource := make(map[string]any)

	if name, ok := d.Get("name").(string); ok {
		identitySource["name"] = name
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
	}

	if description, ok := d.Get("description").(string); ok {
		identitySource["description"] = description
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("description", d.Get("description")))
	}

	identitySource["type"] = "saml"

	config := make(map[string]any)

	if loginRedirectURL, ok := d.Get("login_redirect_url").(string); ok {
		config["url"] = loginRedirectURL
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("login_redirect_url", d.Get("login_redirect_url")))
	}

	if logoutRedirectURL, ok := d.Get("logout_redirect_url").(string); ok {
		config["logoutUrl"] = logoutRedirectURL
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("logout_redirect_url", d.Get("logout_redirect_url")))
	}

	if includeSAMLRequest, ok := d.Get("include_saml_request_parameter").(bool); ok {
		if includeSAMLRequest {
			config["doNotIncludeSAMLRequest"] = false
		} else {
			config["doNotIncludeSAMLRequest"] = true
		}
	} else {
		return diag.FromErr(
			helpers.TypeAssertFailError(
				"include_saml_request_parameter",
				d.Get("include_saml_request_parameter"),
			),
		)
	}

	if samlRequest, ok := d.Get("saml_request").(string); ok {
		config["SAMLSignatureMode"] = samlRequest
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("saml_request", d.Get("saml_request")))
	}

	if givenNameAttribute, ok := d.Get("given_name_attribute").(string); ok {
		config["givenNameAttribute"] = givenNameAttribute
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("given_name_attribute", d.Get("given_name_attribute")))
	}

	if surnameAttribute, ok := d.Get("surname_attribute").(string); ok {
		config["surnameAttribute"] = surnameAttribute
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("surname_attribute", d.Get("surname_attribute")))
	}

	if emailAttribute, ok := d.Get("email_attribute").(string); ok {
		config["emailAttribute"] = emailAttribute
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("email_attribute", d.Get("email_attribute")))
	}

	if roleAttributeName, ok := d.Get("role_attribute_name").(string); ok {
		config["roleAttributeName"] = roleAttributeName
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("role_attribute_name", d.Get("role_attribute_name")))
	}

	if requiredAttributeValue, ok := d.Get("required_role_attribute_value").(string); ok {
		config["requiredAttributeValue"] = requiredAttributeValue
	} else {
		return diag.FromErr(
			helpers.TypeAssertFailError(
				"required_role_attribute_value",
				d.Get("required_role_attribute_value"),
			),
		)
	}

	if validateAssertionSignature, ok := d.Get("validate_assertion_signature").(bool); ok {
		if validateAssertionSignature {
			config["doNotValidateSignature"] = false
		} else {
			config["doNotValidateSignature"] = true
		}
	} else {
		return diag.FromErr(
			helpers.TypeAssertFailError(
				"validate_assertion_signature",
				d.Get("validate_assertion_signature"),
			),
		)
	}

	identitySource["config"] = config

	defaultAccountRole := make(map[string]any)

	if defaultAccountRoleID, ok := d.Get("default_account_role_id").(int); ok {
		defaultAccountRole["id"] = defaultAccountRoleID
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("default_account_role_id", d.Get("default_account_role_id")))
	}

	identitySource["defaultAccountRole"] = defaultAccountRole

	// Role Mappings
	if roleMappingValue := d.Get("role_mapping"); roleMappingValue != "" {
		if roleMappingSet, ok := roleMappingValue.(*schema.Set); ok {
			identitySource["roleMappings"] = parseSAMLRoleMappings(roleMappingSet)
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("role_mapping", roleMappingValue))
		}
	}

	if enableRoleMappingPermission, ok := d.Get("enable_role_mapping_permission").(bool); ok {
		identitySource["allowCustomMappings"] = enableRoleMappingPermission
	} else {
		return diag.FromErr(
			helpers.TypeAssertFailError(
				"enable_role_mapping_permission",
				d.Get("enable_role_mapping_permission"),
			),
		)
	}

	req := &morpheus.Request{
		Body: map[string]any{
			"userSource": identitySource,
		},
	}

	resp, err := client.UpdateIdentitySource(convert.StringToInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	var result *morpheus.UpdateIdentitySourceResult
	if resultAssert, ok := resp.Result.(*morpheus.UpdateIdentitySourceResult); ok {
		result = resultAssert
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("resp.Result", resp.Result))
	}

	identitySourceResult := result.IdentitySource
	if identitySourceResult == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("IdentitySource"))
	}

	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(convert.Int64ToString(identitySourceResult.ID))

	return resourceIdentitySourceSAMLRead(ctx, d, meta)
}

func resourceIdentitySourceSAMLDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	var client *morpheus.Client
	if clientAssert, ok := meta.(*morpheus.Client); ok {
		client = clientAssert
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", clientAssert))
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
	log.Printf("API RESPONSE: %s", resp)
	d.SetId("")

	return diags
}

func parseSAMLRoleMappings(mappings *schema.Set) []map[string]any {
	var roleMappings []map[string]any
	// iterate over the array of roleMappings
	for _, mapping := range mappings.List() {
		row := make(map[string]any)
		mappedRole := make(map[string]any)

		if mappingConfig, ok := mapping.(map[string]any); ok {
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
				case "assertion_attribute":
					if sourceRoleName, ok := v.(string); ok {
						row["sourceRoleName"] = sourceRoleName
					}
				}
			}
		}

		row["mappedRole"] = mappedRole
		roleMappings = append(roleMappings, row)
	}

	return roleMappings
}

// samlProviderSettings extracts the SAML service-provider metadata that Morpheus
// computes and returns under userSource.providerSettings. The legacy SDK model
// does not expose these fields, so they are read from the raw response body.
func samlProviderSettings(body []byte) (entityID, acsURL, spMetadata string, err error) {
	var parsed struct {
		UserSource struct {
			ProviderSettings struct {
				EntityID   string `json:"entityId"`
				AcsURL     string `json:"acsUrl"`
				SpMetadata string `json:"spMetadata"`
			} `json:"providerSettings"`
		} `json:"userSource"`
	}

	if err = json.Unmarshal(body, &parsed); err != nil {
		return "", "", "", err
	}

	ps := parsed.UserSource.ProviderSettings

	return ps.EntityID, ps.AcsURL, ps.SpMetadata, nil
}
