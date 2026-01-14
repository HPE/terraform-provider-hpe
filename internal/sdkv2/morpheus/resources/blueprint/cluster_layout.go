// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package blueprint

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/helpers"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

const (
	evarNameKey        = "name"
	evarValueKey       = "value"
	evarMaskedValueKey = "masked_value"
	evarExportKey      = "export"
)

func ResourceClusterLayout() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus cluster layout resource",
		CreateContext: resourceClusterLayoutCreate,
		ReadContext:   resourceClusterLayoutRead,
		UpdateContext: resourceClusterLayoutUpdate,
		DeleteContext: resourceClusterLayoutDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the cluster layout",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the cluster layout",
				Required:    true,
			},
			"version": {
				Type:        schema.TypeString,
				Description: "The version of the cluster layout",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the cluster layout",
				Optional:    true,
				Computed:    true,
			},
			"creatable": {
				Type:        schema.TypeBool,
				Description: "Whether the cluster layout can be used to create clusters or not",
				Optional:    true,
				Computed:    true,
			},
			"minimum_memory": {
				Type:        schema.TypeInt,
				Description: "The minimum amount of memory in bytes",
				Optional:    true,
				Computed:    true,
			},
			"cluster_type_id": {
				Type:        schema.TypeInt,
				Description: "The cluster type ID of the cluster layout",
				Required:    true,
			},
			"provision_type_id": {
				Type:        schema.TypeInt,
				Description: "The provision type ID of the cluster layout",
				Required:    true,
			},
			"enable_scaling": {
				Type:        schema.TypeBool,
				Description: "Whether to enable or disable horizontal scaling",
				Optional:    true,
				Computed:    true,
			},
			"install_docker": {
				Type:        schema.TypeBool,
				Description: "Whether to automatically install Docker or not",
				Optional:    true,
				Computed:    true,
			},
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
									return false
								}
								h := sha256.New()
								h.Write([]byte(new))
								sha256Hash := hex.EncodeToString(h.Sum(nil))

								return strings.EqualFold(strings.ToLower(old), strings.ToLower(sha256Hash))
							},
							DiffSuppressOnRefresh: true,
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
				Description: "A list of option type ids associated with the cluster layout",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return new == old
				},
				Computed: true,
			},
			/* WAITING ON API SUPPORT
			"spec_template_ids": {
				Type:        schema.TypeList,
				Description: "A list of spec templates ids associated with the cluster layout",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return new == old
				},
				Computed: true,
			},*/
			"workflow_id": {
				Type:        schema.TypeInt,
				Description: "Workflow ID to associate with the cluster layout",
				Optional:    true,
				Computed:    true,
			},
			"master_node_pool": {
				Type:        schema.TypeList,
				Description: "Master node configuration",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"count": {
							Type:        schema.TypeInt,
							Description: "The number of nodes",
							Required:    true,
						},
						"node_type_id": {
							Type:        schema.TypeInt,
							Description: "The id of the node type",
							Required:    true,
						},
						"priority_order": {
							Type:        schema.TypeInt,
							Description: "The priority order of the node type",
							Required:    true,
						},
					},
				},
			},
			"worker_node_pool": {
				Type:        schema.TypeList,
				Description: "Worker node configuration",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"count": {
							Type:        schema.TypeInt,
							Description: "The number of nodes",
							Required:    true,
						},
						"node_type_id": {
							Type:        schema.TypeInt,
							Description: "The id of the node type",
							Required:    true,
						},
						"priority_order": {
							Type:        schema.TypeInt,
							Description: "The priority order of the node type",
							Required:    true,
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

func resourceClusterLayoutCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	clusterLayout := make(map[string]any)

	var name string
	if v, ok := d.Get("name").(string); ok {
		name = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
	}
	clusterLayout["name"] = name

	var version string
	if v, ok := d.Get("version").(string); ok {
		version = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("version", d.Get("version")))
	}
	clusterLayout["computeVersion"] = version

	var description string
	if v, ok := d.Get("description").(string); ok {
		description = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("description", d.Get("description")))
	}
	clusterLayout["description"] = description

	var creatable bool
	if v, ok := d.Get("creatable").(bool); ok {
		creatable = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("creatable", d.Get("creatable")))
	}
	clusterLayout["creatable"] = creatable

	provisionType := make(map[string]any)
	var provisionTypeID int
	if v, ok := d.Get("provision_type_id").(int); ok {
		provisionTypeID = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("provision_type_id", d.Get("provision_type_id")))
	}
	provisionType["id"] = provisionTypeID
	clusterLayout["provisionType"] = provisionType

	groupType := make(map[string]any)
	var clusterTypeID int
	if v, ok := d.Get("cluster_type_id").(int); ok {
		clusterTypeID = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("cluster_type_id", d.Get("cluster_type_id")))
	}
	groupType["id"] = clusterTypeID
	clusterLayout["groupType"] = groupType

	var minimumMemory int
	if v, ok := d.Get("minimum_memory").(int); ok {
		minimumMemory = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("minimum_memory", d.Get("minimum_memory")))
	}
	clusterLayout["memoryRequirement"] = minimumMemory

	var workflowID int
	if v, ok := d.Get("workflow_id").(int); ok {
		workflowID = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("workflow_id", d.Get("workflow_id")))
	}
	if workflowID > 0 {
		taskSet := make(map[string]any)
		taskSet["id"] = workflowID
		var taskSets [1]map[string]any
		taskSets[0] = taskSet
		clusterLayout["taskSets"] = taskSets
	}

	var enableScaling bool
	if v, ok := d.Get("enable_scaling").(bool); ok {
		enableScaling = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("enable_scaling", d.Get("enable_scaling")))
	}
	clusterLayout["hasAutoScale"] = enableScaling

	var installDocker bool
	if v, ok := d.Get("install_docker").(bool); ok {
		installDocker = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("install_docker", d.Get("install_docker")))
	}
	clusterLayout["installContainerRuntime"] = installDocker

	// input types
	var optionTypes []map[string]any
	if d.Get("option_type_ids") != nil {
		var optionTypeList []any
		if v, ok := d.Get("option_type_ids").([]any); ok {
			optionTypeList = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("option_type_ids", d.Get("option_type_ids")))
		}
		if optionTypeList != nil {
			// iterate over the array of tasks
			for i := 0; i < len(optionTypeList); i++ {
				row := make(map[string]any)
				row["id"] = optionTypeList[i]
				optionTypes = append(optionTypes, row)
			}
		}
	}

	clusterLayout["optionTypes"] = optionTypes

	/* WAITING ON API SUPPORT
	// spec templates
	var specTemplates []map[string]any
	if d.Get("spec_template_ids") != nil {
		var specTemplateList []any
		if v, ok := d.Get("spec_template_ids").([]any); ok {
			specTemplateList = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("spec_template_ids", d.Get("spec_template_ids")))
		}
		if specTemplateList != nil {
			// iterate over the array of spec templates
			for i := 0; i < len(specTemplateList); i++ {
				row := make(map[string]any)
				row["id"] = specTemplateList[i]
				specTemplates = append(specTemplates, row)
			}
		}
	}

	clusterLayout["specTemplates"] = specTemplates
	*/

	var evar []any
	if v, ok := d.Get("evar").([]any); ok {
		evar = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("evar", d.Get("evar")))
	}
	clusterLayout["environmentVariables"] = parseClusterLayoutEnvironmentVariables(evar)

	var masterNodePool []any
	if v, ok := d.Get("master_node_pool").([]any); ok {
		masterNodePool = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("master_node_pool", d.Get("master_node_pool")))
	}
	clusterLayout["masters"] = parseClusterLayoutNodePools(masterNodePool)

	var workerNodePool []any
	if v, ok := d.Get("worker_node_pool").([]any); ok {
		workerNodePool = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("worker_node_pool", d.Get("worker_node_pool")))
	}
	clusterLayout["workers"] = parseClusterLayoutNodePools(workerNodePool)

	req := &morpheus.Request{
		Body: map[string]any{
			"layout": clusterLayout,
		},
	}

	resp, err := client.CreateClusterLayout(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	var result map[string]any
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		log.Println(err)
	}

	clusterLayoutID := fmt.Sprintf("%v", result["id"])

	// Successfully created resource, now set id
	d.SetId(clusterLayoutID)

	diags = append(diags, resourceClusterLayoutRead(ctx, d, meta)...)

	return diags
}

func resourceClusterLayoutRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
		resp, err = client.FindClusterLayoutByName(name)
	} else if id != "" {
		resp, err = client.GetClusterLayout(convert.StringToInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Cluster layout cannot be read without name or id")
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

	// store resource data
	var clusterLayout ClusterLayoutPayload
	json.Unmarshal(resp.Body, &clusterLayout)

	d.SetId(convert.Int64ToString(clusterLayout.ClusterLayout.ID))
	d.Set("name", clusterLayout.ClusterLayout.Name)
	d.Set("version", clusterLayout.ClusterLayout.ComputeVersion)
	d.Set("description", clusterLayout.ClusterLayout.Description)
	d.Set("creatable", clusterLayout.ClusterLayout.Creatable)
	d.Set("minimum_memory", clusterLayout.ClusterLayout.MemoryRequirement)

	if clusterLayout.ClusterLayout.ProvisionType.ID == 0 {
		return diag.FromErr(helpers.NotFoundInResponseError("ProvisionType"))
	}
	d.Set("provision_type_id", clusterLayout.ClusterLayout.ProvisionType.ID)

	if clusterLayout.ClusterLayout.GroupType.ID == 0 {
		return diag.FromErr(helpers.NotFoundInResponseError("GroupType"))
	}
	d.Set("cluster_type_id", clusterLayout.ClusterLayout.GroupType.ID)

	d.Set("install_docker", clusterLayout.ClusterLayout.InstallContainerRuntime)
	d.Set("enable_scaling", clusterLayout.ClusterLayout.HasAutoScale)

	if clusterLayout.ClusterLayout.TaskSets != nil {
		if len(clusterLayout.ClusterLayout.TaskSets) > 0 {
			d.Set("workflow_id", clusterLayout.ClusterLayout.TaskSets[0].ID)
		}
	}

	var evars []map[string]any
	if clusterLayout.ClusterLayout.EnvironmentVariables != nil {
		// iterate over the array of environment variables
		for i := 0; i < len(clusterLayout.ClusterLayout.EnvironmentVariables); i++ {
			environmentVariable := clusterLayout.ClusterLayout.EnvironmentVariables[i]
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
	if clusterLayout.ClusterLayout.OptionTypes != nil {
		// iterate over the array of option types
		for i := 0; i < len(clusterLayout.ClusterLayout.OptionTypes); i++ {
			input := clusterLayout.ClusterLayout.OptionTypes[i]
			inputs = append(inputs, input.ID)
		}
	}

	var optionTypeIDs []any
	if v, ok := d.Get("option_type_ids").([]any); ok {
		optionTypeIDs = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("option_type_ids", d.Get("option_type_ids")))
	}
	stateInputs := matchTemplatesWithSchema(inputs, optionTypeIDs)
	d.Set("option_type_ids", stateInputs)

	/* WAITING ON API SUPPORT
	// spec templates
	var specTemplates []int64
	if clusterLayout.ClusterLayout.SpecTemplates != nil {
		// iterate over the array of spec templates
		for i := 0; i < len(clusterLayout.ClusterLayout.SpecTemplates); i++ {
			specTemplate := clusterLayout.ClusterLayout.SpecTemplates[i]
			specTemplates = append(specTemplates, specTemplate.ID)
		}
	}

	var specTemplateIDs []any
	if v, ok := d.Get("spec_template_ids").([]any); ok {
		specTemplateIDs = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("spec_template_ids", d.Get("spec_template_ids")))
	}
	stateSpecTemplates := matchTemplatesWithSchema(specTemplates, specTemplateIDs)
	d.Set("spec_template_ids", stateSpecTemplates)
	*/

	return diags
}

func resourceClusterLayoutUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}
	id := d.Id()

	clusterLayout := make(map[string]any)

	var name string
	if v, ok := d.Get("name").(string); ok {
		name = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
	}
	clusterLayout["name"] = name

	var version string
	if v, ok := d.Get("version").(string); ok {
		version = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("version", d.Get("version")))
	}
	clusterLayout["computeVersion"] = version

	var description string
	if v, ok := d.Get("description").(string); ok {
		description = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("description", d.Get("description")))
	}
	clusterLayout["description"] = description

	var creatable bool
	if v, ok := d.Get("creatable").(bool); ok {
		creatable = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("creatable", d.Get("creatable")))
	}
	clusterLayout["creatable"] = creatable

	provisionType := make(map[string]any)
	var provisionTypeID int
	if v, ok := d.Get("provision_type_id").(int); ok {
		provisionTypeID = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("provision_type_id", d.Get("provision_type_id")))
	}
	provisionType["id"] = provisionTypeID
	clusterLayout["provisionType"] = provisionType

	groupType := make(map[string]any)
	var clusterTypeID int
	if v, ok := d.Get("cluster_type_id").(int); ok {
		clusterTypeID = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("cluster_type_id", d.Get("cluster_type_id")))
	}
	groupType["id"] = clusterTypeID
	clusterLayout["groupType"] = groupType

	var minimumMemory int
	if v, ok := d.Get("minimum_memory").(int); ok {
		minimumMemory = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("minimum_memory", d.Get("minimum_memory")))
	}
	clusterLayout["memoryRequirement"] = minimumMemory

	var workflowID int
	if v, ok := d.Get("workflow_id").(int); ok {
		workflowID = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("workflow_id", d.Get("workflow_id")))
	}
	if workflowID > 0 {
		taskSet := make(map[string]any)
		taskSet["id"] = workflowID
		var taskSets [1]map[string]any
		taskSets[0] = taskSet
		clusterLayout["taskSets"] = taskSets
	}

	// option types
	var optionTypes []map[string]any
	if d.Get("option_type_ids") != nil {
		var optionTypeList []any
		if v, ok := d.Get("option_type_ids").([]any); ok {
			optionTypeList = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("option_type_ids", d.Get("option_type_ids")))
		}
		if optionTypeList != nil {
			// iterate over the array of option types
			for i := 0; i < len(optionTypeList); i++ {
				row := make(map[string]any)
				row["id"] = optionTypeList[i]
				optionTypes = append(optionTypes, row)
			}
		}
	}

	clusterLayout["optionTypes"] = optionTypes

	/* WAITING ON API SUPPORT
	// spec templates
	var specTemplates []map[string]any
	if d.Get("spec_template_ids") != nil {
		var specTemplateList []any
		if v, ok := d.Get("spec_template_ids").([]any); ok {
			specTemplateList = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("spec_template_ids", d.Get("spec_template_ids")))
		}
		if specTemplateList != nil {
			// iterate over the array of spec templates
			for i := 0; i < len(specTemplateList); i++ {
				row := make(map[string]any)
				row["id"] = specTemplateList[i]
				specTemplates = append(specTemplates, row)
			}
		}
	}

	clusterLayout["specTemplates"] = specTemplates
	*/
	var enableScaling bool
	if v, ok := d.Get("enable_scaling").(bool); ok {
		enableScaling = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("enable_scaling", d.Get("enable_scaling")))
	}
	clusterLayout["hasAutoScale"] = enableScaling

	var evar []any
	if v, ok := d.Get("evar").([]any); ok {
		evar = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("evar", d.Get("evar")))
	}
	clusterLayout["environmentVariables"] = parseClusterLayoutEnvironmentVariables(evar)

	var masterNodePool []any
	if v, ok := d.Get("master_node_pool").([]any); ok {
		masterNodePool = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("master_node_pool", d.Get("master_node_pool")))
	}
	clusterLayout["masters"] = parseClusterLayoutNodePools(masterNodePool)

	var workerNodePool []any
	if v, ok := d.Get("worker_node_pool").([]any); ok {
		workerNodePool = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("worker_node_pool", d.Get("worker_node_pool")))
	}
	clusterLayout["workers"] = parseClusterLayoutNodePools(workerNodePool)

	req := &morpheus.Request{
		Body: map[string]any{
			"layout": clusterLayout,
		},
	}

	resp, err := client.UpdateClusterLayout(convert.StringToInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	// log.Printf("API RESPONSE: %s", resp)

	var result map[string]any
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		log.Println(err)
	}

	clusterLayoutID := fmt.Sprintf("%v", result["id"])

	// Successfully updated resource, now set id
	d.SetId(clusterLayoutID)

	return resourceClusterLayoutRead(ctx, d, meta)
}

func resourceClusterLayoutDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
	resp, err := client.DeleteClusterLayout(convert.StringToInt64(id), req)
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

func parseClusterLayoutNodePools(variables []any) []map[string]any {
	var nodepools []map[string]any
	if variables == nil {
		return nodepools
	}
	// iterate over the array of nodepools
	for i := 0; i < len(variables); i++ {
		row := make(map[string]any)
		nodepoolconfig := variables[i].(map[string]any)
		for k, v := range nodepoolconfig {
			switch k {
			case "count":
				row["nodeCount"] = v.(int)
			case "node_type_id":
				nodeType := make(map[string]any)
				nodeType["id"] = v.(int)
				row["containerType"] = nodeType
			case "priority_order":
				row["priorityOrder"] = v.(int)
			}
		}
		nodepools = append(nodepools, row)
	}

	return nodepools
}

func parseClusterLayoutEnvironmentVariables(variables []any) []map[string]any {
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
			case evarNameKey:
				row[evarNameKey] = v.(string)
				// row["evarName"] = v.(string)
				// row["valueType"] = "fixed"
			case evarValueKey:
				if v.(string) != "" {
					row[evarValueKey] = v.(string)
					row["masked"] = false
				}
			case evarMaskedValueKey:
				if v.(string) != "" {
					row[evarValueKey] = v.(string)
					row["masked"] = true
				}
			case evarExportKey:
				row[evarExportKey] = v.(bool)
			}
		}
		evars = append(evars, row)
	}

	return evars
}

type ClusterLayoutPayload struct {
	ClusterLayout struct {
		ID      int64 `json:"id"`
		Account struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
		} `json:"account"`
		Name                    string `json:"name"`
		Description             string `json:"description"`
		Code                    string `json:"code"`
		ComputeVersion          string `json:"computeVersion"`
		HasAutoScale            bool   `json:"hasAutoScale"`
		Creatable               bool   `json:"creatable"`
		MemoryRequirement       int64  `json:"memoryRequirement"`
		InstallContainerRuntime bool   `json:"installContainerRuntime"`
		ProvisionType           struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
			Code string `json:"code"`
		} `json:"provisionType"`
		GroupType struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
			Code string `json:"code"`
		} `json:"groupType"`
		TaskSets []struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
		} `json:"taskSets"`
		SpecTemplates []struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
		} `json:"specTemplates"`
		OptionTypes []struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
		} `json:"optionTypes"`
		EnvironmentVariables []struct {
			EvarName         string `json:"evarName"`
			Name             string `json:"name"`
			DefaultValue     string `json:"defaultValue"`
			DefaultValueHash string `json:"defaultValueHash"`
			ValueType        string `json:"valueType"`
			Export           bool   `json:"export"`
			Masked           bool   `json:"masked"`
		} `json:"environmentVariables"`
		ComputeServers []struct {
			ID                      int64  `json:"id"`
			PriorityOrder           int64  `json:"priorityOrder"`
			NodeCount               int64  `json:"nodeCount"`
			NodeType                string `json:"nodeType"`
			MinNodeCount            int64  `json:"minNodeCount"`
			MaxNodeCount            any    `json:"maxNodeCount"`
			DynamicCount            bool   `json:"dynamicCount"`
			InstallContainerRuntime bool   `json:"installContainerRuntime"`
			InstallStorageRuntime   bool   `json:"installStorageRuntime"`
			Name                    string `json:"name"`
			Code                    string `json:"code"`
			Category                any    `json:"category"`
			Config                  any    `json:"config"`
			ContainertType          struct {
				ID               int64  `json:"id"`
				Account          any    `json:"account"`
				Name             string `json:"name"`
				ShortName        string `json:"shortName"`
				Code             string `json:"code"`
				ContainerVersion string `json:"containerVersion"`
				ProvisionType    struct {
					ID   int64  `json:"id"`
					Name string `json:"name"`
					Code string `json:"code"`
				} `json:"provisionType"`
				VirtualImage   any      `json:"virtualImage"`
				Category       any      `json:"category"`
				Config         struct{} `json:"config"`
				ContainerPorts []struct {
					ID                  int64  `json:"id"`
					Name                string `json:"name"`
					Port                int64  `json:"port"`
					LoadBalanceProtocol any    `json:"loadBalanceProtocol"`
					ExportName          string `json:"exportName"`
				} `json:"containerPorts"`
				ContainerScripts   []any `json:"containerScripts"`
				ContainerTemplates []struct {
					ID   int64  `json:"id"`
					Name string `json:"name"`
				} `json:"containerTemplates"`
				EnvironmentVariables []any `json:"environmentVariables"`
			} `json:"containerType"`
			ComputeServerType struct {
				ID             any `json:"id"`
				Code           any `json:"code"`
				Name           any `json:"name"`
				Managed        any `json:"managed"`
				ExternalDelete any `json:"externalDelete"`
			} `json:"computeServerType"`
			ProvisionService any    `json:"provisionService"`
			PlanCategory     any    `json:"planCategory"`
			NamePrefix       any    `json:"namePrefix"`
			NameSuffix       string `json:"nameSuffix"`
			ForceNameIndex   bool   `json:"forceNameIndex"`
			LoadBalance      bool   `json:"loadBalance"`
		}
	} `json:"layout"`
}
