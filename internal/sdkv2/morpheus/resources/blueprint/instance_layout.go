// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package blueprint

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"

	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/helpers"
)

//nolint:lll
func ResourceInstanceLayout() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus instance layout resource",
		CreateContext: resourceInstanceLayoutCreate,
		ReadContext:   resourceInstanceLayoutRead,
		UpdateContext: resourceInstanceLayoutUpdate,
		DeleteContext: resourceInstanceLayoutDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the instance layout",
				Computed:    true,
			},
			"instance_type_id": {
				Type:        schema.TypeInt,
				Description: "The ID of the associated instance type",
				Required:    true,
				ForceNew:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the instance layout",
				Required:    true,
			},
			"version": {
				Type:        schema.TypeString,
				Description: "The version of the instance layout",
				Required:    true,
			},
			"labels": {
				Type:        schema.TypeSet,
				Description: "The organization labels associated with the script template (Only supported on Morpheus 5.5.3 or higher)",
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The instance layout category",
				Optional:    true,
				Computed:    true,
			},
			"creatable": {
				Type:        schema.TypeBool,
				Description: "Whether the instance layout can be used to create an instance",
				Optional:    true,
				Computed:    true,
			},
			"technology": {
				Type:         schema.TypeString,
				Description:  "The technology of the instance layout (alibaba, amazon, arm, azure, maas, cloudFormation, docker, esxi, fusion, google, huawei, hyperv, kubernetes, kvm, nutanix, opentelekom, openstack, oraclecloud, oraclevm, scvmm, terraform, upcloud, vcd.vapp, vcd, vmware, workflow, xen)",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"alibaba", "amazon", "arm", "azure", "maas", "cloudFormation", "docker", "esxi", "fusion", "google", "huawei", "hyperv", "kubernetes", "kvm", "nutanix", "opentelekom", "openstack", "oraclecloud", "oraclevm", "scvmm", "terraform", "upcloud", "vcd.vapp", "vcd", "vmware", "workflow", "xen"}, false),
			},
			"minimum_memory": {
				Type:        schema.TypeInt,
				Description: "The memory requirement in megabytes",
				Optional:    true,
				Computed:    true,
			},
			"workflow_id": {
				Type:        schema.TypeInt,
				Description: "The id of the provisioning workflow associated with the instance layout",
				Optional:    true,
				Computed:    true,
			},
			"support_convert_to_managed": {
				Type:        schema.TypeBool,
				Description: "Whether the instance layout supports deployed instances to be converted to managed",
				Optional:    true,
				Computed:    true,
			},
			/* AWAITING API SUPPORT
			"enable_scaling": {
				Type:        schema.TypeBool,
				Description: "The instance layout category",
				Optional:    true,
			},
			*/
			"evar": {
				Type:        schema.TypeList,
				Description: "The environment variables to create",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "The name of the environment variable",
							Optional:    true,
						},
						"value": {
							Type:        schema.TypeString,
							Description: "The environment variable value when the value can be in plaintext",
							Optional:    true,
							Computed:    true,
						},
						"masked_value": {
							Type:        schema.TypeString,
							Description: "The environment variable value when the value needs to be masked",
							Optional:    true,
							Sensitive:   true,
							Computed:    true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								if old == "" {
									return true
								}
								h := sha256.New()
								h.Write([]byte(new))
								sha256Hash := hex.EncodeToString(h.Sum(nil))

								return strings.EqualFold(strings.ToLower(old), strings.ToLower(sha256Hash))
							},
						},
						"export": {
							Type:        schema.TypeBool,
							Description: "Whether the environment variable is exported as an instance tag",
							Optional:    true,
						},
					},
				},
			},
			"option_type_ids": {
				Type:        schema.TypeList,
				Description: "A list of option type ids associated with the instance layout",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return new == old
				},
				Computed: true,
			},
			"node_type_ids": {
				Type:        schema.TypeList,
				Description: "A list of node type ids associated with the instance layout",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return new == old
				},
				Computed: true,
			},
			"spec_template_ids": {
				Type:        schema.TypeList,
				Description: "A list of spec template ids associated with the instance layout",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return new == old
				},
				Computed: true,
			},
			"price_set_ids": {
				Type:        schema.TypeList,
				Description: "A list of price set ids associated with the instance layout",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return new == old
				},
				Computed: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

//nolint:goconst
func resourceInstanceLayoutCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics

	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	instanceLayout := make(map[string]any)

	var name string
	if v, ok := d.Get("name").(string); ok {
		name = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
	}
	instanceLayout["name"] = name

	var version string
	if v, ok := d.Get("version").(string); ok {
		version = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("version", d.Get("version")))
	}
	instanceLayout["instanceVersion"] = version

	var description string
	if v, ok := d.Get("description").(string); ok {
		description = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("description", d.Get("description")))
	}
	instanceLayout["description"] = description

	var creatable bool
	if v, ok := d.Get("creatable").(bool); ok {
		creatable = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("creatable", d.Get("creatable")))
	}
	instanceLayout["creatable"] = creatable

	var technology string
	if v, ok := d.Get("technology").(string); ok {
		technology = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("technology", d.Get("technology")))
	}
	instanceLayout["provisionTypeCode"] = technology

	var minimumMemory int
	if v, ok := d.Get("minimum_memory").(int); ok {
		minimumMemory = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("minimum_memory", d.Get("minimum_memory")))
	}
	memoryRequirement := convert.IntToString(minimumMemory)
	instanceLayout["memoryRequirement"] = memoryRequirement

	var workflowID int
	if v, ok := d.Get("workflow_id").(int); ok {
		workflowID = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("workflow_id", d.Get("workflow_id")))
	}
	instanceLayout["taskSetId"] = workflowID

	var supportConvertToManaged bool
	if v, ok := d.Get("support_convert_to_managed").(bool); ok {
		supportConvertToManaged = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("support_convert_to_managed", d.Get("support_convert_to_managed")))
	}
	instanceLayout["supportsConvertToManaged"] = supportConvertToManaged

	instanceLayout["optionTypes"] = d.Get("option_type_ids")

	var evarRaw []any
	if v, ok := d.Get("evar").([]any); ok {
		evarRaw = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("evar", d.Get("evar")))
	}
	instanceLayout["environmentVariables"] = parseInstanceLayoutEnvironmentVariables(evarRaw, d)

	//nolint:lll
	switch technology {
	case "alibaba", "amazon", "azure", "maas", "docker", "esxi", "fusion", "google", "huawei", "hyperv", "kubernetes", "kvm", "nutanix", "opentelekom", "openstack", "oraclecloud", "oraclevm", "scvmm", "upcloud", "vcd.vapp", "vcd", "vmware", "xen":
		instanceLayout["containerTypes"] = d.Get("node_type_ids")
	case "arm", "cloudFormation", "terraform":
		instanceLayout["specTemplates"] = d.Get("spec_template_ids")
	case "workflow":
		break
	}

	labelsPayload := make([]string, 0)
	if attr, ok := d.GetOk("labels"); ok {
		var labelSet *schema.Set
		if v, ok := attr.(*schema.Set); ok {
			labelSet = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("label", attr))
		}

		labelList := labelSet.List()
		for _, s := range labelList {
			var label string
			if v, ok := s.(string); ok {
				label = v
			} else {
				return diag.FromErr(helpers.TypeAssertFailError("label", s))
			}
			labelsPayload = append(labelsPayload, label)
		}
	}
	instanceLayout["labels"] = labelsPayload

	// priceSets
	var priceSets []map[string]any
	if d.Get("price_set_ids") != nil {
		var priceSetList []any
		if v, ok := d.Get("price_set_ids").([]any); ok {
			priceSetList = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("price_set_ids", d.Get("price_set_ids")))
		}

		// iterate over the array of tasks
		for i := 0; i < len(priceSetList); i++ {
			row := make(map[string]any)
			row["id"] = priceSetList[i]
			priceSets = append(priceSets, row)
		}
	}
	instanceLayout["priceSets"] = priceSets

	req := &morpheus.Request{
		Body: map[string]any{
			"instanceTypeLayout": instanceLayout,
		},
	}

	var instanceTypeID int
	if v, ok := d.Get("instance_type_id").(int); ok {
		instanceTypeID = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("instance_type_id", d.Get("instance_type_id")))
	}

	resp, err := client.CreateInstanceLayout(int64(instanceTypeID), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	var result *morpheus.CreateInstanceLayoutResult
	if v, ok := resp.Result.(*morpheus.CreateInstanceLayoutResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("CreateInstanceLayoutResult", resp.Result))
	}

	if result.InstanceLayout == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("InstanceLayout"))
	}
	instanceLayoutResponse := result.InstanceLayout
	// Successfully created resource, now set id
	d.SetId(convert.Int64ToString(instanceLayoutResponse.ID))

	diags = append(diags, resourceInstanceLayoutRead(ctx, d, meta)...)

	return diags
}

//nolint:staticcheck
func resourceInstanceLayoutRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics

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

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindInstanceLayoutByName(name)
	} else if id != "" {
		resp, err = client.GetInstanceLayout(convert.StringToInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Instance type cannot be read without name or id")
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
	// log.Printf("API RESPONSE: %s", resp)

	// store resource data
	var instanceLayout InstanceLayoutPayload
	json.Unmarshal(resp.Body, &instanceLayout)

	d.SetId(convert.Int64ToString(instanceLayout.InstanceLayout.ID))

	if instanceLayout.InstanceLayout.InstanceType.ID == 0 {
		return diag.FromErr(helpers.NotFoundInResponseError("InstanceType"))
	}
	d.Set("instance_type_id", instanceLayout.InstanceLayout.InstanceType.ID)

	d.Set("name", instanceLayout.InstanceLayout.Name)
	d.Set("version", instanceLayout.InstanceLayout.ContainerVersion)
	d.Set("description", instanceLayout.InstanceLayout.Description)
	d.Set("labels", instanceLayout.Labels)

	if instanceLayout.InstanceLayout.ProvisionType.ID == 0 {
		return diag.FromErr(helpers.NotFoundInResponseError("ProvisionType"))
	}
	d.Set("technology", instanceLayout.InstanceLayout.ProvisionType.Code)

	d.Set("creatable", instanceLayout.InstanceLayout.Creatable)
	memoryRequirement := instanceLayout.InstanceLayout.MemoryRequirement / 1024 / 1024
	d.Set("minimum_memory", memoryRequirement)

	if len(instanceLayout.InstanceLayout.TaskSets) > 0 {
		d.Set("workflow_id", instanceLayout.InstanceLayout.TaskSets[0].ID)
	}
	d.Set("support_convert_to_managed", instanceLayout.InstanceLayout.SupportsConvertToManaged)

	var evars []map[string]any
	if instanceLayout.InstanceLayout.EnvironmentVariables != nil {
		// iterate over the array of environment variables
		for i := 0; i < len(instanceLayout.InstanceLayout.EnvironmentVariables); i++ {
			environmentVariable := instanceLayout.InstanceLayout.EnvironmentVariables[i]
			envPayload := make(map[string]any)
			envPayload["name"] = environmentVariable.Name
			if environmentVariable.Masked {
				envPayload["masked_value"] = environmentVariable.DefaultValueHash
			} else {
				envPayload["value"] = environmentVariable.DefaultValue
			}
			envPayload["export"] = environmentVariable.Export
			evars = append(evars, envPayload)
		}
	}
	d.Set("evar", evars)

	// inputs
	var inputs []int64
	if instanceLayout.InstanceLayout.OptionTypes != nil {
		// iterate over the array of option types
		for i := 0; i < len(instanceLayout.InstanceLayout.OptionTypes); i++ {
			input := instanceLayout.InstanceLayout.OptionTypes[i]
			inputs = append(inputs, input.ID)
		}
	}

	var optionTypeIDsRaw []any
	if v, ok := d.Get("option_type_ids").([]any); ok {
		optionTypeIDsRaw = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("option_type_ids", d.Get("option_type_ids")))
	}
	stateInputs := matchTemplatesWithSchema(inputs, optionTypeIDsRaw)
	d.Set("option_type_ids", stateInputs)

	// spec templates
	if d.Get("spec_template_ids") != nil {
		var specTemplates []int64
		if instanceLayout.InstanceLayout.SpecTemplates != nil {
			// iterate over the array of script templates
			for i := 0; i < len(instanceLayout.InstanceLayout.SpecTemplates); i++ {
				specTemplate := instanceLayout.InstanceLayout.SpecTemplates[i]
				specTemplates = append(specTemplates, specTemplate.ID)
			}
		}

		var specTemplateIDsRaw []any
		if v, ok := d.Get("spec_template_ids").([]any); ok {
			specTemplateIDsRaw = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("spec_template_ids", d.Get("spec_template_ids")))
		}
		stateSpecTemplates := matchTemplatesWithSchema(specTemplates, specTemplateIDsRaw)
		d.Set("spec_template_ids", stateSpecTemplates)
	}

	// node types
	if d.Get("node_type_ids") != nil {
		var nodeTypes []int64
		if instanceLayout.InstanceLayout.ContainerTypes != nil {
			// iterate over the array of node types
			for i := 0; i < len(instanceLayout.InstanceLayout.ContainerTypes); i++ {
				nodeType := instanceLayout.InstanceLayout.ContainerTypes[i]
				nodeTypes = append(nodeTypes, nodeType.ID)
			}
		}

		var nodeTypeIDsRaw []any
		if v, ok := d.Get("node_type_ids").([]any); ok {
			nodeTypeIDsRaw = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("node_type_ids", d.Get("node_type_ids")))
		}
		stateNodeTypes := matchTemplatesWithSchema(nodeTypes, nodeTypeIDsRaw)
		d.Set("node_type_ids", stateNodeTypes)
	}

	// priceSets
	var priceSets []int64
	if instanceLayout.PriceSets != nil {
		// iterate over the array of price sets
		for i := 0; i < len(instanceLayout.PriceSets); i++ {
			priceSet := instanceLayout.PriceSets[i]
			priceSets = append(priceSets, priceSet.ID)
		}
	}

	var priceSetIDsRaw []any
	if v, ok := d.Get("price_set_ids").([]any); ok {
		priceSetIDsRaw = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("price_set_ids", d.Get("price_set_ids")))
	}
	priceSetData := matchTemplatesWithSchema(priceSets, priceSetIDsRaw)
	d.Set("price_set_ids", priceSetData)

	return diags
}

//nolint:goconst
func resourceInstanceLayoutUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	id := d.Id()

	instanceLayout := make(map[string]any)

	var name string
	if v, ok := d.Get("name").(string); ok {
		name = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
	}
	instanceLayout["name"] = name

	var version string
	if v, ok := d.Get("version").(string); ok {
		version = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("version", d.Get("version")))
	}
	instanceLayout["instanceVersion"] = version

	var description string
	if v, ok := d.Get("description").(string); ok {
		description = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("description", d.Get("description")))
	}
	instanceLayout["description"] = description

	var creatable bool
	if v, ok := d.Get("creatable").(bool); ok {
		creatable = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("creatable", d.Get("creatable")))
	}
	instanceLayout["creatable"] = creatable

	var technology string
	if v, ok := d.Get("technology").(string); ok {
		technology = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("technology", d.Get("technology")))
	}
	instanceLayout["provisionTypeCode"] = technology

	var minimumMemory int
	if v, ok := d.Get("minimum_memory").(int); ok {
		minimumMemory = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("minimum_memory", d.Get("minimum_memory")))
	}
	memoryRequirement := convert.IntToString(minimumMemory)
	instanceLayout["memoryRequirement"] = memoryRequirement

	var workflowID int
	if v, ok := d.Get("workflow_id").(int); ok {
		workflowID = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("workflow_id", d.Get("workflow_id")))
	}
	instanceLayout["taskSetId"] = workflowID

	var supportConvertToManaged bool
	if v, ok := d.Get("support_convert_to_managed").(bool); ok {
		supportConvertToManaged = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("support_convert_to_managed", d.Get("support_convert_to_managed")))
	}
	instanceLayout["supportsConvertToManaged"] = supportConvertToManaged

	instanceLayout["optionTypes"] = d.Get("option_type_ids")

	var evarRaw []any
	if v, ok := d.Get("evar").([]any); ok {
		evarRaw = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("evar", d.Get("evar")))
	}
	instanceLayout["environmentVariables"] = parseInstanceLayoutEnvironmentVariables(evarRaw, d)

	//nolint:lll
	switch technology {
	case "alibaba", "amazon", "azure", "maas", "docker", "esxi", "fusion", "google", "huawei", "hyperv", "kubernetes", "kvm", "nutanix", "opentelekom", "openstack", "oraclecloud", "oraclevm", "scvmm", "upcloud", "vcd.vapp", "vcd", "vmware", "xen":
		instanceLayout["containerTypes"] = d.Get("node_type_ids")
	case "arm", "cloudFormation", "terraform":
		instanceLayout["specTemplates"] = d.Get("spec_template_ids")
	case "workflow":
		break
	}

	labelsPayload := make([]string, 0)
	if attr, ok := d.GetOk("labels"); ok {
		var labelSet *schema.Set
		if v, ok := attr.(*schema.Set); ok {
			labelSet = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("label", attr))
		}

		labelList := labelSet.List()
		for _, s := range labelList {
			var label string
			if v, ok := s.(string); ok {
				label = v
			} else {
				return diag.FromErr(helpers.TypeAssertFailError("label", s))
			}
			labelsPayload = append(labelsPayload, label)
		}
	}
	instanceLayout["labels"] = labelsPayload

	// priceSets
	var priceSets []map[string]any
	if d.Get("price_set_ids") != nil {
		var priceSetList []any
		if v, ok := d.Get("price_set_ids").([]any); ok {
			priceSetList = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("price_set_ids", d.Get("price_set_ids")))
		}

		// iterate over the array of tasks
		for i := 0; i < len(priceSetList); i++ {
			row := make(map[string]any)
			row["id"] = priceSetList[i]
			priceSets = append(priceSets, row)
		}
	}
	instanceLayout["priceSets"] = priceSets

	req := &morpheus.Request{
		Body: map[string]any{
			"instanceTypeLayout": instanceLayout,
		},
	}

	resp, err := client.UpdateInstanceLayout(convert.StringToInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	var result *morpheus.UpdateInstanceLayoutResult
	if v, ok := resp.Result.(*morpheus.UpdateInstanceLayoutResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("UpdateInstanceLayoutResult", resp.Result))
	}

	if result.InstanceLayout == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("InstanceLayout"))
	}
	instanceLayoutResponse := result.InstanceLayout
	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(convert.Int64ToString(instanceLayoutResponse.ID))

	return resourceInstanceLayoutRead(ctx, d, meta)
}

func resourceInstanceLayoutDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics

	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteInstanceLayout(convert.StringToInt64(id), req)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %s", resp, err)

			return nil
		}
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	// log.Printf("API RESPONSE: %s", resp)
	d.SetId("")

	return diags
}

func parseInstanceLayoutEnvironmentVariables(variables []any, d *schema.ResourceData) []map[string]any {
	var evars []map[string]any
	if variables == nil {
		return evars
	}

	// iterate over the array of evars
	for i := 0; i < len(variables); i++ {
		row := make(map[string]any)
		evarconfig := variables[i].(map[string]any)
		for k, v := range evarconfig {
			switch k {
			case "name":
				row["name"] = v.(string)
				row["evarName"] = v.(string)
				row["valueType"] = "fixed"
			case "value":
				if v.(string) != "" {
					row["value"] = v.(string)
					row["masked"] = false
				}
			case "masked_value":
				if v.(string) != "" {
					row["value"] = v.(string)
					row["masked"] = true
				}
			case "export":
				row["export"] = v.(bool)
			}
		}
		evars = append(evars, row)
	}

	return evars
}

type InstanceLayoutPayload struct {
	morpheus.InstanceLayout `json:"instanceTypeLayout"`
}
