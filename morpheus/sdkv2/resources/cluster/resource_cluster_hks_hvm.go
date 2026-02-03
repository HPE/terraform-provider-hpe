// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package cluster

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/convert"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/helpers"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func ResourceClusterHKSHVM() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a HPE Kubernetes Service (HKS) cluster on a HVM cloud",
		CreateContext: resourceClusterHKSHVMCreate,
		ReadContext:   resourceClusterHKSHVMRead,
		UpdateContext: resourceClusterHKSHVMUpdate,
		DeleteContext: resourceClusterHKSHVMDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(45 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(45 * time.Minute),
			Delete: schema.DefaultTimeout(45 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of the cluster",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"api_endpoint": {
				Description: "The API URL of the cluster",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"kubernetes_version": {
				Description: "The Kubernetes version of the cluster",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "The name of the cluster",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"resource_prefix": {
				Description: "The prefix used for the virtual machine name of the master and worker nodes",
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
			},
			"hostname_prefix": {
				Description: "The prefix used for the guest operating system hostname of the master and worker nodes",
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
			},
			"description": {
				Description: "The user friendly description of the cluster",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"cloud_id": {
				Description: "The ID of the cloud associated with the cluster",
				Type:        schema.TypeInt,
				ForceNew:    true,
				Required:    true,
			},
			"group_id": {
				Description: "The ID of the group associated with the cluster",
				Type:        schema.TypeInt,
				ForceNew:    true,
				Required:    true,
			},
			"cluster_layout_id": {
				Description: "The ID of the cluster layout to provision the cluster from",
				Type:        schema.TypeInt,
				ForceNew:    true,
				Required:    true,
			},
			"api_proxy_id": {
				Description: "The ID of the api proxy associated with the cluster",
				Type:        schema.TypeInt,
				ForceNew:    true,
				Optional:    true,
			},
			// AWAITING API Support
			// "visibility": {
			//	Type:         schema.TypeString,
			//	Description:  "The visibility of the cluster (public or private)",
			//	Required:     true,
			//	ValidateFunc: validation.StringInSlice([]string{"public", "private"}, false),
			//},
			"pod_cidr": {
				Description: "The cluster pod cidr (default - 172.20.0.0/16)",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "172.20.0.0/16",
			},
			"service_cidr": {
				Description: "The cluster service cidr (default - 172.30.0.0/16)",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "172.30.0.0/16",
			},
			// AWAITING API Support
			//"labels": {
			//	Type:        schema.TypeList,
			//	Description: "The list of labels to add to the cluster",
			//	Optional:    true,
			//	Elem: &schema.Schema{
			//		Type: schema.TypeString,
			//	},
			//	Computed: true,
			//},
			"workflow_id": {
				Description: "The ID of the provisioning workflow to execute",
				Type:        schema.TypeInt,
				ForceNew:    true,
				Optional:    true,
			},
			"workers": {
				Description:  "The number of worker nodes",
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(3),
			},
			"server": {
				Type:        schema.TypeList,
				Description: "Master node pool configuration",
				ForceNew:    true,
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"plan_id": {
							Description: "The ID of the service plan associated with the master nodes in the cluster",
							Type:        schema.TypeInt,
							ForceNew:    true,
							Required:    true,
						},
						"resource_pool_id": {
							Description: "The ID of the resource pool to provision the cluster master nodes to",
							Type:        schema.TypeInt,
							ForceNew:    true,
							Optional:    true,
							Computed:    true,
						},
						"storage_volume": {
							Description: "The storage volumes to create for the cluster master nodes",
							Type:        schema.TypeList,
							ForceNew:    true,
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"uuid": {
										Description: "The storage volume uuid",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"root": {
										Description: "Whether the volume is the root volume of the instance",
										Type:        schema.TypeBool,
										ForceNew:    true,
										Required:    true,
									},
									"name": {
										Description: "The name of the volume",
										Type:        schema.TypeString,
										ForceNew:    true,
										Required:    true,
									},
									"size": {
										Description: "The size of the volume in GB",
										Type:        schema.TypeInt,
										ForceNew:    true,
										Required:    true,
									},
									"storage_type": {
										Description: "The storage volume type ID",
										Type:        schema.TypeInt,
										ForceNew:    true,
										Required:    true,
									},
									"datastore_id": {
										Description: "The ID of the datastore",
										Type:        schema.TypeInt,
										ForceNew:    true,
										Required:    true,
									},
								},
							},
						},
						"network_interface": {
							Description: "The network interfaces to create for the cluster master nodes",
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"network_id": {
										Description: "The ID of the network to assign the network interface to",
										Type:        schema.TypeInt,
										ForceNew:    true,
										Required:    true,
									},
									/* AWAITING API Support for the master node pool for consistency
									"network_interface_type_id": {
										Description: "The id of the network interface type",
										Type:        schema.TypeInt,
										Optional:    true,
									},
									*/
								},
							},
						},
						"tags": {
							Description: "Tags to assign to the cluster master nodes",
							Type:        schema.TypeMap,
							ForceNew:    false,
							Optional:    true,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
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

func resourceClusterHKSHVMCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var client *morpheus.Client
	if clientAssert, ok := meta.(*morpheus.Client); ok {
		client = clientAssert
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	clusterPayload := map[string]interface{}{}

	var name string
	if nameValue, ok := d.Get("name").(string); ok {
		name = nameValue
		clusterPayload["name"] = name
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
	}
	clusterPayload["type"] = "kubernetes-cluster"
	clusterPayload["autoRecoverPowerState"] = false

	// Cloud
	var cloudID int
	if cloudIDValue, ok := d.Get("cloud_id").(int); ok {
		cloudID = cloudIDValue
		clusterPayload["cloud"] = map[string]interface{}{
			"id": cloudID,
		}
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("cloud_id", d.Get("cloud_id")))
	}

	// Group
	var groupID int
	if groupIDValue, ok := d.Get("group_id").(int); ok {
		groupID = groupIDValue
		clusterPayload["group"] = map[string]interface{}{
			"id": groupID,
		}
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("group_id", d.Get("group_id")))
	}

	// Labels - AWAITING API support
	//if d.Get("labels") != nil {
	//	clusterPayload["labels"] = d.Get("labels")
	//}

	// Description
	if d.Get("description") != nil {
		if description, ok := d.Get("description").(string); ok {
			clusterPayload["description"] = description
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("description", d.Get("description")))
		}
	}

	// Cluster Layout
	var clusterLayoutID int
	if clusterLayoutIDValue, ok := d.Get("cluster_layout_id").(int); ok {
		clusterLayoutID = clusterLayoutIDValue
		clusterPayload["layout"] = map[string]interface{}{
			"id": clusterLayoutID,
		}
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("cluster_layout_id", d.Get("cluster_layout_id")))
	}

	// Workflow
	var workflowID int
	if workflowIDValue, ok := d.Get("workflow_id").(int); ok {
		workflowID = workflowIDValue
		clusterPayload["taskSetId"] = workflowID
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("workflow_id", d.Get("workflow_id")))
	}

	var masterpool map[string]interface{}

	if masterPoolValue := d.Get("server"); masterPoolValue != nil {
		if masterPoolSlice, ok := masterPoolValue.([]interface{}); ok {
			if len(masterPoolSlice) == 0 {
				return diag.FromErr(helpers.EmptySliceError("server"))
			}
			if masterPoolMap, ok := masterPoolSlice[0].(map[string]interface{}); ok {
				masterpool = masterPoolMap
			} else {
				return diag.FromErr(helpers.TypeAssertFailError("server[0]", masterPoolSlice[0]))
			}
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("server", masterPoolValue))
		}
	}

	nodeCount := 0
	if nc, ok := d.Get("workers").(int); ok {
		nodeCount = nc
	}

	serverPayload := map[string]interface{}{}

	var podCIDR, serviceCIDR string

	if podCIDRValue, ok := d.Get("pod_cidr").(string); ok {
		podCIDR = podCIDRValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("pod_cidr", d.Get("pod_cidr")))
	}

	if serviceCIDRValue, ok := d.Get("service_cidr").(string); ok {
		serviceCIDR = serviceCIDRValue
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("service_cidr", d.Get("service_cidr")))
	}

	serverPayloadConfig := make(map[string]interface{})
	serverPayloadConfig["podCidr"] = podCIDR
	serverPayloadConfig["serviceCidr"] = serviceCIDR
	serverPayloadConfig["resourcePoolId"] = masterpool["resource_pool_id"]
	serverPayloadConfig["nodeCount"] = nodeCount

	serverPayload["config"] = serverPayloadConfig
	serverPayload["nodeCount"] = nodeCount
	// serverPayload["visibility"] = d.Get("visibility").(string)

	if storageVolumeValue, ok := masterpool["storage_volume"].([]interface{}); ok {
		serverPayload["volumes"] = parseStorageVolumes(storageVolumeValue)
	} else {
		return diag.FromErr(
			helpers.TypeAssertFailError(
				"server.storage_volume",
				masterpool["storage_volume"],
			),
		)
	}

	if networkInterfaceValue, ok := masterpool["network_interface"].([]interface{}); ok {
		serverPayload["networkInterfaces"] = parseMasterNetworkInterfaces(networkInterfaceValue)
	} else {
		return diag.FromErr(
			helpers.TypeAssertFailError(
				"server.network_interface",
				masterpool["network_interface"],
			),
		)
	}

	if masterpool["tags"] != nil {
		if tagsValue, ok := masterpool["tags"].(map[string]interface{}); ok {
			serverPayload["tags"] = parseTags(tagsValue)
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("server.tags", masterpool["tags"]))
		}
	}

	serverPayload["plan"] = map[string]interface{}{
		"id": masterpool["plan_id"],
	}

	var apiProxyID int
	if apiProxyIDValue, ok := d.Get("api_proxy_id").(int); ok {
		apiProxyID = apiProxyIDValue
		if apiProxyID > 0 {
			serverPayload["apiProxy"] = map[string]interface{}{
				"id": apiProxyID,
			}
		}
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("api_proxy_id", d.Get("api_proxy_id")))
	}

	var hostnamePrefix, resourcePrefix string

	if hostnamePrefixValue, ok := d.Get("hostname_prefix").(string); ok {
		hostnamePrefix = hostnamePrefixValue
		if hostnamePrefix != "" {
			serverPayload["hostname"] = hostnamePrefix
		}
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("hostname_prefix", d.Get("hostname_prefix")))
	}

	if resourcePrefixValue, ok := d.Get("resource_prefix").(string); ok {
		resourcePrefix = resourcePrefixValue
		if resourcePrefix != "" {
			serverPayload["name"] = resourcePrefix
		}
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("resource_prefix", d.Get("resource_prefix")))
	}

	serverPayload["softwareRaid"] = "off"
	serverPayload["visibility"] = "public"
	serverPayload["dataDevice"] = "/dev/sdb"
	serverPayload["lvmEnabled"] = "off"
	serverPayload["network"] = map[string]interface{}{"name": "eth0"}

	clusterPayload["server"] = serverPayload

	req := &morpheus.Request{Body: map[string]interface{}{
		"cluster": clusterPayload,
	}}

	// bs, _ := json.MarshalIndent(clusterPayload, "", "  ")
	// log.Printf("API PAYLOAD: %s\n", string(bs))
	// // return diag.FromErr(fmt.Errorf("debug stop"))

	resp, err := client.CreateCluster(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	var result *morpheus.CreateClusterResult
	if resultAssert, ok := resp.Result.(*morpheus.CreateClusterResult); ok {
		result = resultAssert
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("resp.Result", resp.Result))
	}

	if result.Cluster == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Cluster"))
	}
	cluster := result.Cluster
	clusterStatus := statusProvisioning

	stateConf := &retry.StateChangeConf{
		Pending: []string{statusProvisioning, statusStarting, statusStopping, statusPending, statusSyncing},
		//nolint:lll
		Target: []string{statusRunning, statusFailed, statusWarning, statusDenied, statusCancelled, statusSuspended, statusOk},
		Refresh: func() (interface{}, string, error) {
			clusterDetails, err := client.GetCluster(cluster.ID, &morpheus.Request{})
			if err != nil {
				return "", "", err
			}
			log.Printf("API RESPONSE: %s", clusterDetails)

			var clusterResult *morpheus.GetClusterResult
			if clusterResultAssert, ok := clusterDetails.Result.(*morpheus.GetClusterResult); ok {
				clusterResult = clusterResultAssert
			} else {
				return "", "", helpers.TypeAssertFailError("clusterDetails.Result", clusterDetails.Result)
			}

			cluster := clusterResult.Cluster
			clusterStatus = cluster.Status
			if clusterStatus == statusFailed {
				hostsDetails, err := client.ListHosts(&morpheus.Request{
					QueryParams: map[string]string{
						"clusterId": strconv.Itoa(int(cluster.ID)),
					},
				})
				if err != nil {
					log.Printf("API FAILURE: %s - %s", resp, err)
				}

				var hostsResult *morpheus.ListHostsResult
				if hostsResultAssert, ok := hostsDetails.Result.(*morpheus.ListHostsResult); ok {
					hostsResult = hostsResultAssert
				} else {
					return clusterResult, clusterStatus, helpers.TypeAssertFailError("hostsDetails.Result", hostsDetails.Result)
				}

				if hostsResult.Hosts == nil {
					return clusterResult, clusterStatus, helpers.NotFoundInResponseError("Hosts")
				}

				for _, host := range *hostsResult.Hosts {
					// Override the cluster status if the worker nodes are still provisioning
					// to avoid a false failure while the cluster is still being deployed. This is
					// a workaround that has been fixed in 8.0.4 but has been added for legacy support.
					if host.Status == statusProvisioning {
						clusterStatus = statusProvisioning
					}
				}
			}
			// Added an arbitrary wait period for cluster refresh.
			// This should probably trigger a cluster refresh and then poll
			// the cluster to reach a definitive state.
			if clusterStatus == statusFailed {
				time.Sleep(3 * time.Minute)
				clusterStatus = statusOk
			}

			return clusterResult, clusterStatus, nil
		},
		Timeout:      3 * time.Hour,
		MinTimeout:   1 * time.Minute,
		Delay:        3 * time.Minute,
		PollInterval: 1 * time.Minute,
	}

	// Wait, catching any errors
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error creating cluster: %s", err)
	}

	// Successfully created resource, now set id
	d.SetId(convert.Int64ToString(cluster.ID))
	diags = append(diags, resourceClusterHKSHVMRead(ctx, d, meta)...)

	// Fail the cluster deployment if the cluster status is in a failed state
	if clusterStatus == statusFailed {
		return diag.Errorf("error creating cluster: failed to create cluster")
	}

	return diags
}

func resourceClusterHKSHVMRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var client *morpheus.Client
	if clientAssert, ok := meta.(*morpheus.Client); ok {
		client = clientAssert
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

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
		resp, err = client.FindClusterByName(name)
	} else if id != "" {
		resp, err = client.GetCluster(convert.StringToInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Cluster cannot be read without name or id")
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
	var result *morpheus.GetClusterResult
	if resultAssert, ok := resp.Result.(*morpheus.GetClusterResult); ok {
		result = resultAssert
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("resp.Result", resp.Result))
	}

	cluster := result.Cluster
	if cluster == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Cluster")) // should not happen
	}

	d.SetId(convert.Int64ToString(cluster.ID))
	d.Set("name", cluster.Name)
	d.Set("description", cluster.Description)
	d.Set("cloud_id", cluster.Zone.Id)
	d.Set("group_id", cluster.Site.Id)
	d.Set("cluster_layout_id", cluster.Layout.Id)
	// d.Set("visibility", cluster.Visibility)
	d.Set("kubernetes_version", cluster.ServiceVersion)
	d.Set("api_endpoint", cluster.ServiceUrl)

	workers, err := getClusterWorkers(client, cluster.ID)
	if err != nil {
		return diag.FromErr(err)
	}
	workers = filterOutClusterWorkersByStatus(workers, statusDeprovisioning)

	if len(workers) == 0 {
		return diag.FromErr(helpers.EmptySliceError("workers"))
	}

	d.Set("workers", len(workers))

	return diags
}

func resourceClusterHKSHVMUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var client *morpheus.Client
	if clientAssert, ok := meta.(*morpheus.Client); ok {
		client = clientAssert
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	clusterId := convert.StringToInt64(d.Id())

	o, n := d.GetChange("workers")

	var oldCount, newCount int

	if oc, ok := o.(int); ok {
		oldCount = oc
	}

	if nc, ok := n.(int); ok {
		newCount = nc
	}

	if newCount != oldCount {
		countDelta := newCount - oldCount

		if countDelta > 0 {
			err := doHVMClusterWorkerAdd(ctx, client, clusterId, countDelta, d)
			if err != nil {
				return diag.Errorf("error adding cluster worker node(s): %s", err)
			}
		} else {
			err := doClusterWorkerDelete(ctx, client, clusterId, countDelta)
			if err != nil {
				return diag.Errorf("error deleting cluster worker node(s): %s", err)
			}
		}
	}

	clusterPayload := map[string]interface{}{}

	if d.HasChange("name") {
		if nameValue, ok := d.Get("name").(string); ok {
			clusterPayload["name"] = nameValue
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
		}
	}

	if d.HasChange("description") {
		if descriptionValue, ok := d.Get("description").(string); ok {
			clusterPayload["description"] = descriptionValue
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("description", d.Get("description")))
		}
	}

	if len(clusterPayload) > 0 {
		req := &morpheus.Request{Body: map[string]interface{}{
			"cluster": clusterPayload,
		}}

		resp, err := client.UpdateCluster(clusterId, req)
		if err != nil {
			log.Printf("API FAILURE: %s - %s", resp, err)

			return diag.FromErr(err)
		}
		log.Printf("API RESPONSE: %s", resp)
	}

	return resourceClusterHKSHVMRead(ctx, d, meta)
}

func resourceClusterHKSHVMDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var client *morpheus.Client
	if clientAssert, ok := meta.(*morpheus.Client); ok {
		client = clientAssert
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{
		QueryParams: map[string]string{
			"removeInstances": "on",
			"removeResources": "on",
		},
	}
	req.QueryParams["force"] = "true"
	resp, err := client.DeleteCluster(convert.StringToInt64(id), req)
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

	stateConf := &retry.StateChangeConf{
		//nolint:lll
		Pending: []string{statusRemoving, statusPendingRemoval, statusStopping, statusPending, statusWarning, statusDeprovisioning},
		Target:  []string{statusRemoved},
		Refresh: func() (interface{}, string, error) {
			clusterDetails, err := client.GetCluster(convert.StringToInt64(id), &morpheus.Request{})
			if clusterDetails.StatusCode == 404 {
				return "", "removed", nil
			}
			if err != nil {
				return "", "", err
			}

			var result *morpheus.GetClusterResult
			if resultAssert, ok := clusterDetails.Result.(*morpheus.GetClusterResult); ok {
				result = resultAssert
			} else {
				return "", "", helpers.TypeAssertFailError("clusterDetails.Result", clusterDetails.Result)
			}

			if result.Cluster == nil {
				return result, "error", helpers.NotFoundInResponseError("Cluster")
			}

			cluster := result.Cluster

			return result, cluster.Status, nil
		},
		Timeout:      30 * time.Minute,
		MinTimeout:   1 * time.Minute,
		Delay:        1 * time.Minute,
		PollInterval: 30 * time.Second,
	}

	// Wait, catching any errors
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error deleting cluster: %s", err)
	}

	d.SetId("")

	return diags
}

func doHVMClusterWorkerAdd(
	ctx context.Context,
	client *morpheus.Client,
	clusterId int64,
	nodeCount int,
	d *schema.ResourceData,
) error {
	workers, err := getClusterWorkers(client, clusterId)
	if err != nil {
		return err
	}

	if len(workers) == 0 {
		return helpers.EmptySliceError("workers")
	}

	worker := workers[0]
	desiredWorkerCount := len(workers) + nodeCount

	var server map[string]interface{}
	if serverValue := d.Get("server"); serverValue != nil {
		if serverSlice, ok := serverValue.([]interface{}); ok {
			if len(serverSlice) == 0 {
				return helpers.EmptySliceError("server")
			}
			if serverMap, ok := serverSlice[0].(map[string]interface{}); ok {
				server = serverMap
			} else {
				return helpers.TypeAssertFailError("server[0]", serverSlice[0])
			}
		} else {
			return helpers.TypeAssertFailError("server", serverValue)
		}
	}

	serverPayload := map[string]interface{}{}

	var podCIDR, serviceCIDR string
	var clusterRepoAccountID, cloudID int

	if podCIDRValue, ok := d.Get("pod_cidr").(string); ok {
		podCIDR = podCIDRValue
	} else {
		return helpers.TypeAssertFailError("pod_cidr", d.Get("pod_cidr"))
	}

	if serviceCIDRValue, ok := d.Get("service_cidr").(string); ok {
		serviceCIDR = serviceCIDRValue
	} else {
		return helpers.TypeAssertFailError("service_cidr", d.Get("service_cidr"))
	}

	if cloudIDValue, ok := d.Get("cloud_id").(int); ok {
		cloudID = cloudIDValue
	} else {
		return helpers.TypeAssertFailError("cloud_id", d.Get("cloud_id"))
	}

	serverPayload["config"] = map[string]interface{}{
		"podCidr":            podCIDR,
		"serviceCidr":        serviceCIDR,
		"nodeCount":          nodeCount, // Might need to go in serverPayload.server
		"resourcePoolId":     worker.ResourcePoolId,
		"defaultRepoAccount": clusterRepoAccountID,
	}

	// We will let Morpheus set the name for us.

	serverPayload["serverType"] = map[string]interface{}{
		"id": worker.ComputeServerType.ID,
	}
	serverPayload["cloud"] = map[string]interface{}{
		"id": cloudID,
	}
	serverPayload["plan"] = map[string]interface{}{
		"id": worker.Plan.ID,
	}

	if storageVolumeValue, ok := server["storage_volume"].([]interface{}); ok {
		serverPayload["volumes"] = parseStorageVolumes(storageVolumeValue)
	} else {
		return helpers.TypeAssertFailError("worker_node_pool.storage_volume", server["storage_volume"])
	}

	if networkInterfaceValue, ok := server["network_interface"].([]interface{}); ok {
		serverPayload["networkInterfaces"] = parseWorkerNetworkInterfacesForWorkerPayload(networkInterfaceValue)
	} else {
		return helpers.TypeAssertFailError("worker_node_pool.network_interface", server["network_interface"])
	}

	serverPayload["nodeCount"] = nodeCount

	if tagsValue, ok := server["tags"].(map[string]interface{}); ok {
		serverPayload["tags"] = parseTags(tagsValue)
	} else {
		return helpers.TypeAssertFailError("worker_node_pool.tags", server["tags"])
	}

	// NOTE: Not needed from Morpheus 8.05 onward
	serverPayload["server"] = map[string]interface{}{
		"network": map[string]interface{}{},
	}

	req := &morpheus.Request{Body: map[string]interface{}{
		"server": serverPayload,
	}}

	resp, err := client.AddClusterWorker(clusterId, req)
	if err != nil {
		log.Printf("API FAILURE - Error in creating cluster worker node(s): %s - %s", resp, err)

		return err
	}

	stateConf := &retry.StateChangeConf{
		Pending: []string{statusProvisioning},
		Target:  []string{statusProvisioned},
		Refresh: func() (interface{}, string, error) {
			log.Printf("Waiting for all cluster worker nodes to be provisioned...")

			workers, err := getClusterWorkers(client, clusterId)
			if err != nil {
				return "", "", err
			}

			failedWorkers := filterClusterWorkersByStatus(workers, statusFailed)
			if len(failedWorkers) > 0 {
				return "", "", fmt.Errorf("failed to provision all cluster worker nodes")
			}

			provisionedWorkers := filterClusterWorkersByStatus(workers, statusProvisioned)
			if len(provisionedWorkers) == desiredWorkerCount {
				return "", statusProvisioned, nil
			}

			return "", statusProvisioning, nil
		},
		Timeout:      30 * time.Minute,
		MinTimeout:   1 * time.Minute,
		Delay:        1 * time.Minute,
		PollInterval: pollIntervalSeconds * time.Second,
	}

	// Wait, catching any errors
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return err
	}

	return nil
}
