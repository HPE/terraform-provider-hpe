// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package integration

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"strconv"
	"strings"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/helpers"
)

func ResourceIntegrationServiceNow() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a ServiceNow integration resource",
		CreateContext: resourceIntegrationServiceNowCreate,
		ReadContext:   resourceIntegrationServiceNowRead,
		UpdateContext: resourceIntegrationServiceNowUpdate,
		DeleteContext: resourceIntegrationServiceNowDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The id of the ServiceNow integration",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the ServiceNow integration",
				Required:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Whether the SerivceNow integration is enabled",
				Optional:    true,
				Computed:    true,
			},
			"url": {
				Type:        schema.TypeString,
				Description: "The url of the ServiceNow instance",
				Required:    true,
			},
			"credential_id": {
				Description:   "The id of the credential store entry used for authentication",
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"username", "password"},
			},
			"username": {
				Type:        schema.TypeString,
				Description: "The username of the account used to connect to ServiceNow",
				Required:    true,
			},
			"password": {
				Type:        schema.TypeString,
				Description: "The password of the account used to connect to ServiceNow",
				Required:    true,
				Sensitive:   true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					h := sha256.New()
					h.Write([]byte(new))
					sha256Hash := hex.EncodeToString(h.Sum(nil))

					return strings.EqualFold(old, sha256Hash)
				},
				DiffSuppressOnRefresh: true,
			},
			"cmdb_custom_mapping": {
				Type: schema.TypeString,
				Description: "A JSON encoded payload to populate a specific field in the ServiceNow table and " +
					"with a specific mapping",
				Optional: true,
			},
			"cmdb_class_mapping": {
				Type:        schema.TypeMap,
				Description: "The mapping between Morpheus server types and ServiceNow CI classes",
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"default_cmdb_business_class": {
				Type:        schema.TypeString,
				Description: "The default ServiceNow table that records are written to if they aren't explicitly defined",
				Optional:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceIntegrationServiceNowCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	integration := make(map[string]any)

	integration["type"] = "serviceNow"
	var name string
	if v, ok := d.Get("name").(string); ok {
		name = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
	}
	integration["name"] = name

	var enabled bool
	if v, ok := d.Get("enabled").(bool); ok {
		enabled = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("enabled", d.Get("enabled")))
	}
	integration["enabled"] = enabled

	var url string
	if v, ok := d.Get("url").(string); ok {
		url = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("url", d.Get("url")))
	}
	integration["serviceUrl"] = url

	config := make(map[string]any)

	var credentialID int
	if v, ok := d.Get("credential_id").(int); ok {
		credentialID = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("credential_id", d.Get("credential_id")))
	}

	if credentialID != 0 {
		credential := make(map[string]any)
		credential["type"] = "username-password"
		credential["id"] = credentialID
		integration["credential"] = credential
	} else {
		var username string
		if v, ok := d.Get("username").(string); ok {
			username = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("username", d.Get("username")))
		}
		integration["serviceUsername"] = username

		var password string
		if v, ok := d.Get("password").(string); ok {
			password = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("password", d.Get("password")))
		}
		integration["servicePassword"] = password
	}

	var cmdbClassMapping map[string]any
	if v, ok := d.Get("cmdb_class_mapping").(map[string]any); ok {
		cmdbClassMapping = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("cmdb_class_mapping", d.Get("cmdb_class_mapping")))
	}

	if cmdbClassMapping != nil {
		classMappingResponse, err := client.GetOptionSource("serviceNowServerMappings", &morpheus.Request{})
		if err != nil {
			return diag.FromErr(err)
		}

		var classMappingResult *morpheus.GetOptionSourceResult
		if v, ok := classMappingResponse.Result.(*morpheus.GetOptionSourceResult); ok {
			classMappingResult = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("classMappingResult", classMappingResponse.Result))
		}

		if classMappingResult.Data == nil {
			return diag.FromErr(helpers.NotFoundInResponseError("Data"))
		}

		classMappingsInput := cmdbClassMapping
		var classMappings []Mapping
		for key, value := range classMappingsInput {
			matchStatus := false
			for _, mapping := range *classMappingResult.Data {
				if key == mapping.Name {
					var classMapping Mapping
					classMapping.Name = mapping.Name

					var valueFloat float64
					if v, ok := mapping.Value.(float64); ok {
						valueFloat = v
					} else {
						return diag.FromErr(helpers.TypeAssertFailError("mapping.Value", mapping.Value))
					}
					classMapping.ID = strconv.Itoa(int(valueFloat))

					var nowClass string
					if v, ok := value.(string); ok {
						nowClass = v
					} else {
						return diag.FromErr(helpers.TypeAssertFailError("nowClass", value))
					}
					classMapping.NowClass = nowClass

					classMappings = append(classMappings, classMapping)
					matchStatus = true
				}
			}
			if !matchStatus {
				return diag.Errorf("The %s cmdb mapping class is not a supported class", key)
			}
		}
		config["serviceNowCmdbClassMapping"] = classMappings
	}

	var defaultCMDBBusinessClass string
	if v, ok := d.Get("default_cmdb_business_class").(string); ok {
		defaultCMDBBusinessClass = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("default_cmdb_business_class", d.Get("default_cmdb_business_class")))
	}
	config["serviceNowCMDBBusinessObject"] = defaultCMDBBusinessClass
	config["serviceNowCustomCmdbMapping"] = d.Get("cmdb_custom_mapping")

	integration["config"] = config

	req := &morpheus.Request{
		Body: map[string]any{
			"integration": integration,
		},
	}

	resp, err := client.CreateIntegration(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	var result *morpheus.CreateIntegrationResult
	if v, ok := resp.Result.(*morpheus.CreateIntegrationResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("result", resp.Result))
	}

	if result.Integration == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Integration"))
	}

	integrationResult := result.Integration
	// Successfully created resource, now set id
	d.SetId(convert.Int64ToString(integrationResult.ID))

	return resourceIntegrationServiceNowRead(ctx, d, meta)
}

func resourceIntegrationServiceNowRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
		resp, err = client.FindIntegrationByName(name)
	} else if id != "" {
		resp, err = client.GetIntegration(convert.StringToInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Integration cannot be read without name or id")
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
	var result *morpheus.GetIntegrationResult
	if v, ok := resp.Result.(*morpheus.GetIntegrationResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("result", resp.Result))
	}

	if result.Integration == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Integration"))
	}

	integration := result.Integration
	d.SetId(convert.Int64ToString(integration.ID))
	d.Set("name", integration.Name)
	d.Set("enabled", integration.Enabled)
	d.Set("url", integration.URL)

	if integration.Credential.ID == 0 {
		d.Set("username", integration.Username)
		d.Set("password", integration.PasswordHash)
	} else {
		d.Set("credential_id", integration.Credential.ID)
	}

	d.Set("cmdb_custom_mapping", integration.Config.ServiceNowCustomCmdbMapping)
	classMappings := make(map[string]any)

	if integration.Config.ServiceNowCmdbClassMapping != nil {
		// iterate over the array of classMappings
		for i := 0; i < len(integration.Config.ServiceNowCmdbClassMapping); i++ {
			classMap := integration.Config.ServiceNowCmdbClassMapping[i]
			classMapName := classMap.Name
			classMappings[classMapName] = classMap.NowClass
		}
	}

	d.Set("cmdb_class_mapping", classMappings)
	d.Set("default_cmdb_business_class", integration.Config.ServiceNowCMDBBusinessObject)

	return diags
}

func resourceIntegrationServiceNowUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}
	id := d.Id()

	integration := make(map[string]any)

	var name string
	if v, ok := d.Get("name").(string); ok {
		name = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
	}
	integration["name"] = name

	var enabled bool
	if v, ok := d.Get("enabled").(bool); ok {
		enabled = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("enabled", d.Get("enabled")))
	}
	integration["enabled"] = enabled

	integration["type"] = "serviceNow"

	var url string
	if v, ok := d.Get("url").(string); ok {
		url = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("url", d.Get("url")))
	}
	integration["serviceUrl"] = url

	config := make(map[string]any)

	var credentialID int
	if v, ok := d.Get("credential_id").(int); ok {
		credentialID = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("credential_id", d.Get("credential_id")))
	}

	if credentialID != 0 {
		credential := make(map[string]any)
		credential["type"] = "username-password"
		credential["id"] = credentialID
		integration["credential"] = credential
	} else {
		if d.HasChange("username") {
			var username string
			if v, ok := d.Get("username").(string); ok {
				username = v
			} else {
				return diag.FromErr(helpers.TypeAssertFailError("username", d.Get("username")))
			}
			integration["serviceUsername"] = username
		}
		if d.HasChange("password") {
			var password string
			if v, ok := d.Get("password").(string); ok {
				password = v
			} else {
				return diag.FromErr(helpers.TypeAssertFailError("password", d.Get("password")))
			}
			integration["servicePassword"] = password
		}
	}

	var cmdbClassMapping map[string]any
	if v, ok := d.Get("cmdb_class_mapping").(map[string]any); ok {
		cmdbClassMapping = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("cmdb_class_mapping", d.Get("cmdb_class_mapping")))
	}

	if cmdbClassMapping != nil {
		// Query the API to fetch the ID of the class map
		classMappingResponse, err := client.GetOptionSource("serviceNowServerMappings", &morpheus.Request{})
		if err != nil {
			return diag.FromErr(err)
		}

		var classMappingResult *morpheus.GetOptionSourceResult
		if v, ok := classMappingResponse.Result.(*morpheus.GetOptionSourceResult); ok {
			classMappingResult = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("classMappingResult", classMappingResponse.Result))
		}

		if classMappingResult.Data == nil {
			return diag.FromErr(helpers.NotFoundInResponseError("Data"))
		}

		classMappingsInput := cmdbClassMapping
		var classMappings []Mapping
		for key, value := range classMappingsInput {
			matchStatus := false
			for _, mapping := range *classMappingResult.Data {
				if key == mapping.Name {
					var classMapping Mapping
					classMapping.Name = mapping.Name

					var valueFloat float64
					if v, ok := mapping.Value.(float64); ok {
						valueFloat = v
					} else {
						return diag.FromErr(helpers.TypeAssertFailError("mapping.Value", mapping.Value))
					}
					classMapping.ID = strconv.Itoa(int(valueFloat))

					var nowClass string
					if v, ok := value.(string); ok {
						nowClass = v
					} else {
						return diag.FromErr(helpers.TypeAssertFailError("nowClass", value))
					}
					classMapping.NowClass = nowClass

					classMappings = append(classMappings, classMapping)
					matchStatus = true
				}
			}
			if !matchStatus {
				return diag.Errorf("The %s cmdb mapping class is not a supported class", key)
			}
		}
		config["serviceNowCmdbClassMapping"] = classMappings
	}

	var defaultCMDBBusinessClass string
	if v, ok := d.Get("default_cmdb_business_class").(string); ok {
		defaultCMDBBusinessClass = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("default_cmdb_business_class", d.Get("default_cmdb_business_class")))
	}
	config["serviceNowCMDBBusinessObject"] = defaultCMDBBusinessClass

	if d.HasChange("cmdb_custom_mapping") {
		config["serviceNowCustomCmdbMapping"] = d.Get("cmdb_custom_mapping")
	}
	integration["config"] = config

	req := &morpheus.Request{
		Body: map[string]any{
			"integration": integration,
		},
	}

	resp, err := client.UpdateIntegration(convert.StringToInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	var result *morpheus.UpdateIntegrationResult
	if v, ok := resp.Result.(*morpheus.UpdateIntegrationResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("result", resp.Result))
	}

	if result.Integration == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Integration"))
	}

	integrationResult := result.Integration

	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(convert.Int64ToString(integrationResult.ID))

	return resourceIntegrationServiceNowRead(ctx, d, meta)
}

func resourceIntegrationServiceNowDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
	resp, err := client.DeleteIntegration(convert.StringToInt64(id), req)
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

type Mapping struct {
	ID       string `json:"id"`
	Code     string `json:"code"`
	Name     string `json:"name"`
	NowClass string `json:"nowClass"`
}
