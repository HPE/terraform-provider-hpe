// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package blueprint

import (
	"context"
	"encoding/json"
	"log"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"

	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/helpers"
)

const (
	shortNameCharactersWarning = "Short names may not contain spaces or underscores."
)

var shortNameCharacters, _ = regexp.Compile("^[^ _]*$")

func ResourceBlueprintNodeType() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus node type resource",
		CreateContext: resourceBlueprintNodeTypeCreate,
		ReadContext:   resourceBlueprintNodeTypeRead,
		UpdateContext: resourceBlueprintNodeTypeUpdate,
		DeleteContext: resourceBlueprintNodeTypeDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the node type",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the node type",
				Required:    true,
			},
			"short_name": {
				Type:         schema.TypeString,
				Description:  "The short name of the node type",
				Required:     true,
				ValidateFunc: validation.StringMatch(shortNameCharacters, shortNameCharactersWarning),
			},
			"labels": {
				Type: schema.TypeSet,
				Description: "The organization labels associated with the script template " +
					"(Only supported on Morpheus 5.5.3 or higher)",
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"technology": {
				Type: schema.TypeString,
				Description: "The technology of the node type (alibaba, amazon, azure, maas, esxi, fusion, " +
					"google, huawei, hyperv, kvm, nutanix, opentelekom, openstack, oraclecloud, oraclevm, " +
					"scvmm, upcloud, vcd.vapp, vcd, vmware, xen)",
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"alibaba", "amazon", "azure", "maas", "esxi", "fusion", "google", "huawei",
					"hyperv", "kvm", "nutanix", "opentelekom", "openstack", "oraclecloud",
					"oraclevm", "scvmm", "upcloud", "vcd.vapp", "vcd", "vmware", "xen",
				}, false),
			},
			/* AWAITING API SUPPORT TO AVOID DUPLICATE ENTRIES
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
							Description: "The value of the environment variable",
							Optional:    true,
						},
						"export": {
							Type:        schema.TypeBool,
							Description: "Whether the environment variable is exported as an instance tag",
							Optional:    true,
						},
						"masked": {
							Type:        schema.TypeBool,
							Description: "Whether the environment variable is masked for security purposes",
							Optional:    true,
						},
					},
				},
			},*/
			"version": {
				Type:        schema.TypeString,
				Description: "The version of the node type",
				Required:    true,
			},
			"virtual_image_id": {
				Type:        schema.TypeInt,
				Description: "The ID of the virtual image associated with the node type",
				Optional:    true,
				Computed:    true,
			},
			"stat_type_code": {
				Type: schema.TypeString,
				Description: "Supported technology of the node type (server, container, amazon, azure, esxi, " +
					"google, hyperv, nutanix, openstack, scvmm, vmware, xen, docker, virtualbox, vm)",
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"server", "container", "amazon", "azure", "esxi", "google", "hyperv",
					"nutanix", "openstack", "scvmm", "vmware", "xen", "docker", "virtualbox", "vm",
				}, false),
			},
			"log_type_code": {
				Type: schema.TypeString,
				Description: "Supported technology of the node type (server, amazon, azure, esxi, google, " +
					"hyperv, nutanix, openstack, scvmm, vmware, xen, docker, virtualbox, vm)",
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"server", "amazon", "azure", "esxi", "google", "hyperv", "nutanix",
					"openstack", "scvmm", "vmware", "xen", "docker", "virtualbox", "vm",
				}, false),
			},
			/* AWAITING API SUPPORT
			"logs_folder": {
				Type:        schema.TypeString,
				Description: "The log folder associated with the node type",
				Optional:    true,
			},
			"config_folder": {
				Type:        schema.TypeString,
				Description: "The config folder associated with the node type",
				Optional:    true,
			},
			"deploy_folder": {
				Type:        schema.TypeString,
				Description: "The deploy folder associated with the node type",
				Optional:    true,
			},
			*/
			/* Waiting to add support for kubernetes
			"kubernetes_manifest": {
				Type:        schema.TypeString,
				Description: "The kubernetes manifest associated with the node type",
				Optional:    true,
				Computed:    true,
			},
			*/
			"service_port": {
				Type:        schema.TypeList,
				Description: "Service ports associated with the node type",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"port": {
							Type:        schema.TypeString,
							Description: "The port number of the service",
							Optional:    true,
						},
						"name": {
							Type:        schema.TypeString,
							Description: "The name of the service port",
							Optional:    true,
						},
						"protocol": {
							Type:        schema.TypeString,
							Description: "The load balancer protocol (HTTP, HTTPS, TCP)",
							Optional:    true,
						},
					},
				},
			},
			"extra_options": {
				Type:        schema.TypeMap,
				Description: "VMware custom options associated with the node type",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"script_template_ids": {
				Type:        schema.TypeList,
				Description: "A list of script template ids associated with the node type",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return new == old
				},
				Computed: true,
			},
			"file_template_ids": {
				Type:        schema.TypeList,
				Description: "A list of file template ids associated with the node type",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return new == old
				},
				Computed: true,
			},
			"category": {
				Type:        schema.TypeString,
				Description: "The node type category",
				Optional:    true,
				Computed:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

//nolint:goconst
func resourceBlueprintNodeTypeCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	var statTypeCode string
	if v, ok := d.Get("stat_type_code").(string); ok {
		statTypeCode = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("stat_type_code", d.Get("stat_type_code")))
	}

	var logTypeCode string
	if v, ok := d.Get("log_type_code").(string); ok {
		logTypeCode = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("log_type_code", d.Get("log_type_code")))
	}
	config := make(map[string]any)
	if d.Get("extra_options") != nil {
		config["extraOptions"] = d.Get("extra_options")
	}

	containerType := make(map[string]any)
	containerType["name"] = name

	var shortName string
	if v, ok := d.Get("short_name").(string); ok {
		shortName = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("short_name", d.Get("short_name")))
	}
	containerType["shortName"] = shortName

	var version string
	if v, ok := d.Get("version").(string); ok {
		version = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("version", d.Get("version")))
	}
	containerType["containerVersion"] = version

	var technology string
	if v, ok := d.Get("technology").(string); ok {
		technology = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("technology", d.Get("technology")))
	}
	containerType["provisionTypeCode"] = technology

	var virtualImageID int
	if v, ok := d.Get("virtual_image_id").(int); ok {
		virtualImageID = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("virtual_image_id", d.Get("virtual_image_id")))
	}
	if virtualImageID != 0 {
		containerType["virtualImageId"] = virtualImageID
	}

	containerType["config"] = config

	var servicePorts []any
	if v, ok := d.Get("service_port").([]any); ok {
		servicePorts = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("service_port", d.Get("service_port")))
	}
	containerType["containerPorts"] = parseNodeTypeServicePorts(servicePorts)

	containerType["scripts"] = d.Get("script_template_ids")
	containerType["containerTemplates"] = d.Get("file_template_ids")

	var category string
	if v, ok := d.Get("category").(string); ok {
		category = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("category", d.Get("category")))
	}
	containerType["category"] = category
	containerType["serverType"] = "vm"
	containerType["statTypeCode"] = "vm"
	containerType["logTypeCode"] = "vm"

	if statTypeCode != "" {
		containerType["statTypeCode"] = statTypeCode
	}
	if logTypeCode != "" {
		containerType["logTypeCode"] = logTypeCode
	}

	labelsPayload := make([]string, 0)
	if attr, ok := d.GetOk("labels"); ok {
		if labelSet, ok := attr.(*schema.Set); ok {
			for _, s := range labelSet.List() {
				if labelStr, ok := s.(string); ok {
					labelsPayload = append(labelsPayload, labelStr)
				}
			}
		}
	}
	containerType["labels"] = labelsPayload

	req := &morpheus.Request{
		Body: map[string]any{
			"containerType": containerType,
		},
	}

	resp, err := client.CreateNodeType(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	var result *morpheus.CreateNodeTypeResult
	if v, ok := resp.Result.(*morpheus.CreateNodeTypeResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.NodeType == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("NodeType"))
	}

	nodeType := result.NodeType
	d.SetId(convert.Int64ToString(nodeType.ID))

	diags = append(diags, resourceBlueprintNodeTypeRead(ctx, d, meta)...)

	return diags
}

//nolint:goconst
func resourceBlueprintNodeTypeRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
		resp, err = client.FindNodeTypeByName(name)
	} else if id != "" {
		resp, err = client.GetNodeType(convert.StringToInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Node type cannot be read without name or id")
	}

	if err != nil {
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

	var nodeType NodeTypePayload
	json.Unmarshal(resp.Body, &nodeType)

	log.Printf("RESPONSE_PAYLOAD: %v", nodeType)

	if nodeType.ID == 0 {
		return diag.FromErr(helpers.NotFoundInResponseError("ID"))
	}
	d.SetId(convert.Int64ToString(nodeType.ID))

	if nodeType.Name == "" {
		return diag.FromErr(helpers.NotFoundInResponseError("Name"))
	}
	d.Set("name", nodeType.Name)

	if nodeType.ShortName == "" {
		return diag.FromErr(helpers.NotFoundInResponseError("ShortName"))
	}
	d.Set("short_name", nodeType.ShortName)
	d.Set("labels", nodeType.Labels)

	if nodeType.ContainerVersion == "" {
		return diag.FromErr(helpers.NotFoundInResponseError("ContainerVersion"))
	}
	d.Set("version", nodeType.ContainerVersion)
	d.Set("technology", nodeType.ProvisionType.Code)
	d.Set("virtual_image_id", nodeType.VirtualImage.ID)

	if nodeType.ContainerPorts != nil {
		d.Set("service_port", parseServicePortPayload(nodeType.ContainerPorts))
	}

	var scriptTemplates []int64
	if nodeType.ContainerScripts != nil {
		for i := 0; i < len(nodeType.ContainerScripts); i++ {
			scriptTemplate := nodeType.ContainerScripts[i]
			scriptTemplates = append(scriptTemplates, scriptTemplate.ID)
		}
	}

	var fileTemplates []int64
	if nodeType.ContainerTemplates != nil {
		for i := 0; i < len(nodeType.ContainerTemplates); i++ {
			fileTemplate := nodeType.ContainerTemplates[i]
			fileTemplates = append(fileTemplates, fileTemplate.ID)
		}
	}

	var scriptTemplateIDs []any
	if v, ok := d.Get("script_template_ids").([]any); ok {
		scriptTemplateIDs = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("script_template_ids", d.Get("script_template_ids")))
	}
	stateScriptTemplates := matchTemplatesWithSchema(scriptTemplates, scriptTemplateIDs)

	var fileTemplateIDs []any
	if v, ok := d.Get("file_template_ids").([]any); ok {
		fileTemplateIDs = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("file_template_ids", d.Get("file_template_ids")))
	}
	stateFileTemplates := matchTemplatesWithSchema(fileTemplates, fileTemplateIDs)

	d.Set("script_template_ids", stateScriptTemplates)
	d.Set("file_template_ids", stateFileTemplates)
	if nodeType.ProvisionType.Code == "vmware" {
		extraOptions := make(map[string]any)
		if nodeType.Config.ExtraOptions != nil {
			log.Printf("FoundExtraOptions: %s", nodeType.Config.ExtraOptions)
			for k, v := range nodeType.Config.ExtraOptions {
				extraOptions[k] = v
			}
			d.Set("extra_options", extraOptions)
		}
	}
	d.Set("category", nodeType.Category)

	return diags
}

func resourceBlueprintNodeTypeUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	statTypeCode := "vm"
	logTypeCode := "vm"

	var statTypeCodeInput string
	if v, ok := d.Get("stat_type_code").(string); ok {
		statTypeCodeInput = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("stat_type_code", d.Get("stat_type_code")))
	}
	if statTypeCodeInput != "" {
		statTypeCode = statTypeCodeInput
	}

	var logTypeCodeInput string
	if v, ok := d.Get("log_type_code").(string); ok {
		logTypeCodeInput = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("log_type_code", d.Get("log_type_code")))
	}
	if logTypeCodeInput != "" {
		logTypeCode = logTypeCodeInput
	}

	config := make(map[string]any)
	if d.Get("extra_options") != nil {
		config["extraOptions"] = d.Get("extra_options")
	}

	var shortName string
	if v, ok := d.Get("short_name").(string); ok {
		shortName = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("short_name", d.Get("short_name")))
	}

	var version string
	if v, ok := d.Get("version").(string); ok {
		version = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("version", d.Get("version")))
	}

	var technology string
	if v, ok := d.Get("technology").(string); ok {
		technology = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("technology", d.Get("technology")))
	}

	var virtualImageID int
	if v, ok := d.Get("virtual_image_id").(int); ok {
		virtualImageID = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("virtual_image_id", d.Get("virtual_image_id")))
	}

	var servicePorts []any
	if v, ok := d.Get("service_port").([]any); ok {
		servicePorts = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("service_port", d.Get("service_port")))
	}

	var category string
	if v, ok := d.Get("category").(string); ok {
		category = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("category", d.Get("category")))
	}

	req := &morpheus.Request{
		Body: map[string]any{
			"containerType": map[string]any{
				"name":               name,
				"shortName":          shortName,
				"containerVersion":   version,
				"provisionTypeCode":  technology,
				"virtualImageId":     virtualImageID,
				"config":             config,
				"containerPorts":     parseNodeTypeServicePorts(servicePorts),
				"containerScripts":   d.Get("script_template_ids"),
				"containerTemplates": d.Get("file_template_ids"),
				"category":           category,
				"serverType":         "vm",
				"statTypeCode":       statTypeCode,
				"logTypeCode":        logTypeCode,
			},
		},
	}

	resp, err := client.UpdateNodeType(convert.StringToInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	var result *morpheus.UpdateNodeTypeResult
	if v, ok := resp.Result.(*morpheus.UpdateNodeTypeResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.NodeType == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("NodeType"))
	}

	nodeType := result.NodeType
	d.SetId(convert.Int64ToString(nodeType.ID))

	return resourceBlueprintNodeTypeRead(ctx, d, meta)
}

func resourceBlueprintNodeTypeDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteNodeType(convert.StringToInt64(id), req)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %s", resp, err)

			return diag.FromErr(err)
		}

		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	d.SetId("")

	return diags
}

func parseNodeTypeServicePorts(variables []any) []map[string]any {
	if variables == nil {
		return nil
	}

	var svcports []map[string]any
	for i := 0; i < len(variables); i++ {
		row := make(map[string]any)
		svcportconfig := variables[i].(map[string]any)
		for k, v := range svcportconfig {
			switch k {
			case "name":
				row["name"] = v.(string)
			case "port":
				row["port"] = v.(string)
			case "protocol":
				row["loadBalanceProtocol"] = v.(string)
			}
		}
		svcports = append(svcports, row)
	}

	return svcports
}

/* AWAITING API SUPPORT TO AVOID DUPLICATE ENTRIES
func parseNodeTypeEnvironmentVariables(variables []any) []map[string]any {
	var evars []map[string]any
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
			case "port":
				row["value"] = v.(string)
			case "export":
				row["export"] = v.(bool)
			case "masked":
				row["masked"] = v
			}
		}
		evars = append(evars, row)
	}
	return evars
}
*/

func parseServicePortPayload(variables []morpheus.ContainerPort) []map[string]any {
	if variables == nil {
		return nil
	}

	var svcports []map[string]any
	for i := 0; i < len(variables); i++ {
		row := make(map[string]any)
		row["name"] = variables[i].Name
		row["port"] = strconv.Itoa(int(variables[i].Port))
		row["protocol"] = variables[i].LoadBalanceProtocol
		svcports = append(svcports, row)
	}

	return svcports
}

type NodeTypePayload struct {
	morpheus.NodeType `json:"containerType"`
}
