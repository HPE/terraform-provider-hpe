# \ClustersAPI

All URIs are relative to *https://CHANGEME*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AddCluster**](ClustersAPI.md#AddCluster) | **Post** /api/clusters | Create a Cluster
[**AddClusterNamespace**](ClustersAPI.md#AddClusterNamespace) | **Post** /api/clusters/{clusterId}/namespaces | Add Namespace (Kubernetes)
[**AddClusterWorker**](ClustersAPI.md#AddClusterWorker) | **Post** /api/clusters/{clusterId}/servers | Add Worker
[**ApplyTemplate**](ClustersAPI.md#ApplyTemplate) | **Post** /api/clusters/{clusterId}/apply-template | Apply Template to Cluster (Kubernetes)
[**DeleteCluster**](ClustersAPI.md#DeleteCluster) | **Delete** /api/clusters/{clusterId} | Delete a Cluster
[**DeleteClusterAffinityGroup**](ClustersAPI.md#DeleteClusterAffinityGroup) | **Delete** /api/clusters/{clusterId}/affinity-groups/{id} | Delete a Cluster Affinity Group
[**DeleteClusterContainer**](ClustersAPI.md#DeleteClusterContainer) | **Delete** /api/clusters/{clusterId}/containers/{id} | Delete Container
[**DeleteClusterDatastore**](ClustersAPI.md#DeleteClusterDatastore) | **Delete** /api/clusters/{clusterId}/datastores/{id} | Delete a Cluster Datastore
[**DeleteClusterDeployment**](ClustersAPI.md#DeleteClusterDeployment) | **Delete** /api/clusters/{clusterId}/deployments/{id} | Delete Deployment
[**DeleteClusterJob**](ClustersAPI.md#DeleteClusterJob) | **Delete** /api/clusters/{clusterId}/jobs/{id} | Delete a Job
[**DeleteClusterNamespace**](ClustersAPI.md#DeleteClusterNamespace) | **Delete** /api/clusters/{clusterId}/namespaces/{id} | Delete a Namespace (Kubernetes)
[**DeleteClusterService**](ClustersAPI.md#DeleteClusterService) | **Delete** /api/clusters/{clusterId}/services/{id} | Delete a Service
[**DeleteClusterStatefulSet**](ClustersAPI.md#DeleteClusterStatefulSet) | **Delete** /api/clusters/{clusterId}/statefulsets/{id} | Delete a Stateful Set
[**DeleteClusterVolume**](ClustersAPI.md#DeleteClusterVolume) | **Delete** /api/clusters/{clusterId}/volumes/{id} | Delete a Volume
[**DeleteClusterWorker**](ClustersAPI.md#DeleteClusterWorker) | **Delete** /api/clusters/{clusterId}/servers/{id} | Delete a Worker
[**GetCluster**](ClustersAPI.md#GetCluster) | **Get** /api/clusters/{clusterId} | Get a Specific Cluster
[**GetClusterAffinityGroup**](ClustersAPI.md#GetClusterAffinityGroup) | **Get** /api/clusters/{clusterId}/affinity-groups/{id} | Get a Specific Cluster Affinity Group
[**GetClusterApiConfig**](ClustersAPI.md#GetClusterApiConfig) | **Get** /api/clusters/{clusterId}/api-config | Get API Config
[**GetClusterConfigmap**](ClustersAPI.md#GetClusterConfigmap) | **Get** /api/clusters/{clusterId}/configmaps/{id} | Get a Specific Cluster Configmap
[**GetClusterContainer**](ClustersAPI.md#GetClusterContainer) | **Get** /api/clusters/{clusterId}/containers/{id} | Get a Specific Cluster Container
[**GetClusterDaemonset**](ClustersAPI.md#GetClusterDaemonset) | **Get** /api/clusters/{clusterId}/daemonsets/{id} | Get a Specific Cluster Daemonset
[**GetClusterDatastore**](ClustersAPI.md#GetClusterDatastore) | **Get** /api/clusters/{clusterId}/datastores/{id} | Get a Specific Cluster Datastore
[**GetClusterDeployment**](ClustersAPI.md#GetClusterDeployment) | **Get** /api/clusters/{clusterId}/deployments/{id} | Get a Specific Cluster Deployment
[**GetClusterHistory**](ClustersAPI.md#GetClusterHistory) | **Get** /api/clusters/{clusterId}/history | Get Cluster History
[**GetClusterHistoryDetail**](ClustersAPI.md#GetClusterHistoryDetail) | **Get** /api/clusters/{clusterId}/history/{id} | Get Cluster History Details
[**GetClusterHistoryEventDetail**](ClustersAPI.md#GetClusterHistoryEventDetail) | **Get** /api/clusters/{clusterId}/history/events/{id} | Get Cluster History Event
[**GetClusterIngress**](ClustersAPI.md#GetClusterIngress) | **Get** /api/clusters/{clusterId}/ingresses/{id} | Get a Specific Cluster Ingress
[**GetClusterJob**](ClustersAPI.md#GetClusterJob) | **Get** /api/clusters/{clusterId}/jobs/{id} | Get a Specific Cluster Job
[**GetClusterMasters**](ClustersAPI.md#GetClusterMasters) | **Get** /api/clusters/{clusterId}/masters | Get Masters (Kubernetes)
[**GetClusterNamespace**](ClustersAPI.md#GetClusterNamespace) | **Get** /api/clusters/{clusterId}/namespaces/{id} | Get Namespace (Kubernetes)
[**GetClusterNamespaces**](ClustersAPI.md#GetClusterNamespaces) | **Get** /api/clusters/{clusterId}/namespaces | List Namespaces (Kubernetes)
[**GetClusterNetworkEndpoing**](ClustersAPI.md#GetClusterNetworkEndpoing) | **Get** /api/clusters/{clusterId}/endpoints/{id} | Get a Specific Cluster Endpoint
[**GetClusterPod**](ClustersAPI.md#GetClusterPod) | **Get** /api/clusters/{clusterId}/pods/{id} | Get a Specific Pod for a Cluster
[**GetClusterPolicy**](ClustersAPI.md#GetClusterPolicy) | **Get** /api/clusters/{clusterId}/policies/{id} | Get a Specific Cluster Policy
[**GetClusterReplicaset**](ClustersAPI.md#GetClusterReplicaset) | **Get** /api/clusters/{clusterId}/replicasets/{id} | Get a Specific Cluster Replicaset
[**GetClusterSecret**](ClustersAPI.md#GetClusterSecret) | **Get** /api/clusters/{clusterId}/secrets/{id} | Get a Specific Cluster Secret
[**GetClusterService**](ClustersAPI.md#GetClusterService) | **Get** /api/clusters/{clusterId}/services/{id} | Get a Specific Cluster Service
[**GetClusterStatefulset**](ClustersAPI.md#GetClusterStatefulset) | **Get** /api/clusters/{clusterId}/statefulsets/{id} | Get a Specific Cluster Statefulset
[**GetClusterUpgradeVersions**](ClustersAPI.md#GetClusterUpgradeVersions) | **Get** /api/clusters/{clusterId}/upgrade-cluster | Get Cluster Upgrade Versions (Kubernetes)
[**GetClusterVolume**](ClustersAPI.md#GetClusterVolume) | **Get** /api/clusters/{clusterId}/volumes/{id} | Get a Specific Cluster Volume
[**GetClusterVolumeclaim**](ClustersAPI.md#GetClusterVolumeclaim) | **Get** /api/clusters/{clusterId}/volumeclaims/{id} | Get a Specific Cluster VolumeClaim
[**ListClusterAffinityGroups**](ClustersAPI.md#ListClusterAffinityGroups) | **Get** /api/clusters/{clusterId}/affinity-groups | Get Cluster Affinity Groups
[**ListClusterConfigmaps**](ClustersAPI.md#ListClusterConfigmaps) | **Get** /api/clusters/{clusterId}/configmaps | Get Configmaps
[**ListClusterContainers**](ClustersAPI.md#ListClusterContainers) | **Get** /api/clusters/{clusterId}/containers | Get Containers for a Cluster
[**ListClusterDaemonsets**](ClustersAPI.md#ListClusterDaemonsets) | **Get** /api/clusters/{clusterId}/daemonsets | Get Daemonsets
[**ListClusterDatastores**](ClustersAPI.md#ListClusterDatastores) | **Get** /api/clusters/{clusterId}/datastores | Get Cluster Datastores
[**ListClusterDeployments**](ClustersAPI.md#ListClusterDeployments) | **Get** /api/clusters/{clusterId}/deployments | Get Deployments
[**ListClusterIngresses**](ClustersAPI.md#ListClusterIngresses) | **Get** /api/clusters/{clusterId}/ingresses | Get Ingresses
[**ListClusterJobs**](ClustersAPI.md#ListClusterJobs) | **Get** /api/clusters/{clusterId}/jobs | Get Jobs
[**ListClusterNetworkEndpoints**](ClustersAPI.md#ListClusterNetworkEndpoints) | **Get** /api/clusters/{clusterId}/endpoints | Get Endpoints
[**ListClusterPods**](ClustersAPI.md#ListClusterPods) | **Get** /api/clusters/{clusterId}/pods | Get Pods
[**ListClusterPolicies**](ClustersAPI.md#ListClusterPolicies) | **Get** /api/clusters/{clusterId}/policies | Get Policies
[**ListClusterReplicasets**](ClustersAPI.md#ListClusterReplicasets) | **Get** /api/clusters/{clusterId}/replicasets | Get Replicasets
[**ListClusterSecrets**](ClustersAPI.md#ListClusterSecrets) | **Get** /api/clusters/{clusterId}/secrets | Get Secrets
[**ListClusterServices**](ClustersAPI.md#ListClusterServices) | **Get** /api/clusters/{clusterId}/services | Get Services
[**ListClusterStatefulSets**](ClustersAPI.md#ListClusterStatefulSets) | **Get** /api/clusters/{clusterId}/statefulsets | Get Stateful Sets
[**ListClusterTypes**](ClustersAPI.md#ListClusterTypes) | **Get** /api/cluster-types | Get All Cluster Types
[**ListClusterVolumeclaims**](ClustersAPI.md#ListClusterVolumeclaims) | **Get** /api/clusters/{clusterId}/volumeclaims | Get VolumeClaims
[**ListClusterVolumes**](ClustersAPI.md#ListClusterVolumes) | **Get** /api/clusters/{clusterId}/volumes | Get Volumes
[**ListClusterWorkers**](ClustersAPI.md#ListClusterWorkers) | **Get** /api/clusters/{clusterId}/workers | Get Workers
[**ListClusters**](ClustersAPI.md#ListClusters) | **Get** /api/clusters | Get All Clusters
[**RefreshCluster**](ClustersAPI.md#RefreshCluster) | **Get** /api/clusters/{clusterId}/refresh | Refreshes a Cluster
[**RestartClusterContainer**](ClustersAPI.md#RestartClusterContainer) | **Put** /api/clusters/{clusterId}/containers/{id}/restart | Restart a Container
[**RestartClusterDeployment**](ClustersAPI.md#RestartClusterDeployment) | **Put** /api/clusters/{clusterId}/deployments/{id}/restart | Restart a Deployment
[**RestartClusterPod**](ClustersAPI.md#RestartClusterPod) | **Put** /api/clusters/{clusterId}/pods/{id}/restart | Restart a Pod
[**RestartClusterStatefulSet**](ClustersAPI.md#RestartClusterStatefulSet) | **Put** /api/clusters/{clusterId}/statefulsets/{id}/restart | Restart a Stateful Set
[**SaveClusterAffinityGroup**](ClustersAPI.md#SaveClusterAffinityGroup) | **Post** /api/clusters/{clusterId}/affinity-groups | Create a Cluster Affinity Group
[**SaveClusterDatastore**](ClustersAPI.md#SaveClusterDatastore) | **Post** /api/clusters/{clusterId}/datastores | Create a Cluster Datastore
[**UpdateCluster**](ClustersAPI.md#UpdateCluster) | **Put** /api/clusters/{clusterId} | Update Cluster
[**UpdateClusterAffinityGroup**](ClustersAPI.md#UpdateClusterAffinityGroup) | **Put** /api/clusters/{clusterId}/affinity-groups/{id} | Update Cluster Affinity Group
[**UpdateClusterAssignedLayout**](ClustersAPI.md#UpdateClusterAssignedLayout) | **Put** /api/clusters/{clusterId}/layout | Update Cluster Assigned Layout
[**UpdateClusterDatastore**](ClustersAPI.md#UpdateClusterDatastore) | **Put** /api/clusters/{clusterId}/datastores/{id} | Update Cluster Datastore
[**UpdateClusterNamespace**](ClustersAPI.md#UpdateClusterNamespace) | **Put** /api/clusters/{clusterId}/namespaces/{id} | Update Namespace (Kubernetes)
[**UpdateClusterPermissions**](ClustersAPI.md#UpdateClusterPermissions) | **Put** /api/clusters/{clusterId}/permissions | Update Cluster Permissions
[**UpdateClusterUpgradeVersions**](ClustersAPI.md#UpdateClusterUpgradeVersions) | **Post** /api/clusters/{clusterId}/upgrade-cluster | Upgrade a Cluster (Kubernetes)
[**UpdateClusterWorkerCount**](ClustersAPI.md#UpdateClusterWorkerCount) | **Put** /api/clusters/{clusterId}/worker-count | Update Worker Count



## AddCluster

> AddCluster200Response AddCluster(ctx).AddClusterRequest(addClusterRequest).Execute()

Create a Cluster



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	addClusterRequest := *openapiclient.NewAddClusterRequest() // AddClusterRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.AddCluster(context.Background()).AddClusterRequest(addClusterRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.AddCluster``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddCluster`: AddCluster200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.AddCluster`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAddClusterRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **addClusterRequest** | [**AddClusterRequest**](AddClusterRequest.md) |  | 

### Return type

[**AddCluster200Response**](AddCluster200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AddClusterNamespace

> AddClusterNamespace200Response AddClusterNamespace(ctx, clusterId).AddClusterNamespaceRequest(addClusterNamespaceRequest).Execute()

Add Namespace (Kubernetes)



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster
	addClusterNamespaceRequest := *openapiclient.NewAddClusterNamespaceRequest() // AddClusterNamespaceRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.AddClusterNamespace(context.Background(), clusterId).AddClusterNamespaceRequest(addClusterNamespaceRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.AddClusterNamespace``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddClusterNamespace`: AddClusterNamespace200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.AddClusterNamespace`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 

### Other Parameters

Other parameters are passed through a pointer to a apiAddClusterNamespaceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **addClusterNamespaceRequest** | [**AddClusterNamespaceRequest**](AddClusterNamespaceRequest.md) |  | 

### Return type

[**AddClusterNamespace200Response**](AddClusterNamespace200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AddClusterWorker

> AddClusterWorker200Response AddClusterWorker(ctx, clusterId).AddClusterWorkerRequest(addClusterWorkerRequest).Execute()

Add Worker



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster
	addClusterWorkerRequest := *openapiclient.NewAddClusterWorkerRequest() // AddClusterWorkerRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.AddClusterWorker(context.Background(), clusterId).AddClusterWorkerRequest(addClusterWorkerRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.AddClusterWorker``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddClusterWorker`: AddClusterWorker200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.AddClusterWorker`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 

### Other Parameters

Other parameters are passed through a pointer to a apiAddClusterWorkerRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **addClusterWorkerRequest** | [**AddClusterWorkerRequest**](AddClusterWorkerRequest.md) |  | 

### Return type

[**AddClusterWorker200Response**](AddClusterWorker200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApplyTemplate

> ApplyTemplate200Response ApplyTemplate(ctx, clusterId).ApplyTemplateRequest(applyTemplateRequest).Execute()

Apply Template to Cluster (Kubernetes)



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster
	applyTemplateRequest := *openapiclient.NewApplyTemplateRequest() // ApplyTemplateRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.ApplyTemplate(context.Background(), clusterId).ApplyTemplateRequest(applyTemplateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.ApplyTemplate``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApplyTemplate`: ApplyTemplate200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.ApplyTemplate`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 

### Other Parameters

Other parameters are passed through a pointer to a apiApplyTemplateRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **applyTemplateRequest** | [**ApplyTemplateRequest**](ApplyTemplateRequest.md) |  | 

### Return type

[**ApplyTemplate200Response**](ApplyTemplate200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteCluster

> DeleteCluster200Response DeleteCluster(ctx, clusterId).RemoveInstances(removeInstances).RemoveResources(removeResources).PreserveVolumes(preserveVolumes).ReleaseFloatingIps(releaseFloatingIps).ReleaseEIPs(releaseEIPs).Force(force).Execute()

Delete a Cluster



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster
	removeInstances := "on" // string | Remove Instances (optional) (default to "off")
	removeResources := "off" // string | Remove Resources (optional) (default to "on")
	preserveVolumes := "on" // string | Preserve Volumes (optional) (default to "off")
	releaseFloatingIps := "off" // string | Release Floating IPs (optional) (default to "on")
	releaseEIPs := "off" // string | Alias for releaseFloatingIps (optional) (default to "on")
	force := "on" // string | Force Delete (optional) (default to "off")

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.DeleteCluster(context.Background(), clusterId).RemoveInstances(removeInstances).RemoveResources(removeResources).PreserveVolumes(preserveVolumes).ReleaseFloatingIps(releaseFloatingIps).ReleaseEIPs(releaseEIPs).Force(force).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.DeleteCluster``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteCluster`: DeleteCluster200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.DeleteCluster`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteClusterRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **removeInstances** | **string** | Remove Instances | [default to &quot;off&quot;]
 **removeResources** | **string** | Remove Resources | [default to &quot;on&quot;]
 **preserveVolumes** | **string** | Preserve Volumes | [default to &quot;off&quot;]
 **releaseFloatingIps** | **string** | Release Floating IPs | [default to &quot;on&quot;]
 **releaseEIPs** | **string** | Alias for releaseFloatingIps | [default to &quot;on&quot;]
 **force** | **string** | Force Delete | [default to &quot;off&quot;]

### Return type

[**DeleteCluster200Response**](DeleteCluster200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteClusterAffinityGroup

> DeleteClusterAffinityGroup200Response DeleteClusterAffinityGroup(ctx, clusterId, id).Execute()

Delete a Cluster Affinity Group



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster
	id := int64(1) // int64 | Morpheus ID of the Object being referenced

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.DeleteClusterAffinityGroup(context.Background(), clusterId, id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.DeleteClusterAffinityGroup``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteClusterAffinityGroup`: DeleteClusterAffinityGroup200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.DeleteClusterAffinityGroup`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteClusterAffinityGroupRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**DeleteClusterAffinityGroup200Response**](DeleteClusterAffinityGroup200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteClusterContainer

> DeleteClusterContainer200Response DeleteClusterContainer(ctx, clusterId, id).Force(force).Execute()

Delete Container



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster
	id := int64(1) // int64 | Morpheus ID of the Object being referenced
	force := "on" // string | Force Delete (optional) (default to "off")

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.DeleteClusterContainer(context.Background(), clusterId, id).Force(force).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.DeleteClusterContainer``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteClusterContainer`: DeleteClusterContainer200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.DeleteClusterContainer`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteClusterContainerRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **force** | **string** | Force Delete | [default to &quot;off&quot;]

### Return type

[**DeleteClusterContainer200Response**](DeleteClusterContainer200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteClusterDatastore

> DeleteClusterDatastore200Response DeleteClusterDatastore(ctx, clusterId, id).Execute()

Delete a Cluster Datastore



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster
	id := int64(1) // int64 | Morpheus ID of the Object being referenced

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.DeleteClusterDatastore(context.Background(), clusterId, id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.DeleteClusterDatastore``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteClusterDatastore`: DeleteClusterDatastore200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.DeleteClusterDatastore`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteClusterDatastoreRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**DeleteClusterDatastore200Response**](DeleteClusterDatastore200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteClusterDeployment

> DeleteClusterDeployment200Response DeleteClusterDeployment(ctx, clusterId, id).Force(force).Execute()

Delete Deployment



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster
	id := int64(1) // int64 | Morpheus ID of the Object being referenced
	force := "on" // string | Force Delete (optional) (default to "off")

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.DeleteClusterDeployment(context.Background(), clusterId, id).Force(force).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.DeleteClusterDeployment``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteClusterDeployment`: DeleteClusterDeployment200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.DeleteClusterDeployment`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteClusterDeploymentRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **force** | **string** | Force Delete | [default to &quot;off&quot;]

### Return type

[**DeleteClusterDeployment200Response**](DeleteClusterDeployment200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteClusterJob

> DeleteClusterJob200Response DeleteClusterJob(ctx, clusterId, id).Force(force).Execute()

Delete a Job



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster
	id := int64(1) // int64 | Morpheus ID of the Object being referenced
	force := "on" // string | Force Delete (optional) (default to "off")

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.DeleteClusterJob(context.Background(), clusterId, id).Force(force).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.DeleteClusterJob``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteClusterJob`: DeleteClusterJob200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.DeleteClusterJob`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteClusterJobRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **force** | **string** | Force Delete | [default to &quot;off&quot;]

### Return type

[**DeleteClusterJob200Response**](DeleteClusterJob200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteClusterNamespace

> DeleteClusterNamespace200Response DeleteClusterNamespace(ctx, clusterId, id).Force(force).Execute()

Delete a Namespace (Kubernetes)



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster
	id := int64(1) // int64 | Morpheus ID of the Object being referenced
	force := "on" // string | Force Delete (optional) (default to "off")

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.DeleteClusterNamespace(context.Background(), clusterId, id).Force(force).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.DeleteClusterNamespace``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteClusterNamespace`: DeleteClusterNamespace200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.DeleteClusterNamespace`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteClusterNamespaceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **force** | **string** | Force Delete | [default to &quot;off&quot;]

### Return type

[**DeleteClusterNamespace200Response**](DeleteClusterNamespace200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteClusterService

> DeleteClusterService200Response DeleteClusterService(ctx, clusterId, id).Force(force).Execute()

Delete a Service



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster
	id := int64(1) // int64 | Morpheus ID of the Object being referenced
	force := "on" // string | Force Delete (optional) (default to "off")

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.DeleteClusterService(context.Background(), clusterId, id).Force(force).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.DeleteClusterService``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteClusterService`: DeleteClusterService200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.DeleteClusterService`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteClusterServiceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **force** | **string** | Force Delete | [default to &quot;off&quot;]

### Return type

[**DeleteClusterService200Response**](DeleteClusterService200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteClusterStatefulSet

> RestartClusterDeployment200Response DeleteClusterStatefulSet(ctx, clusterId, id).Force(force).Execute()

Delete a Stateful Set



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster
	id := int64(1) // int64 | Morpheus ID of the Object being referenced
	force := "on" // string | Force Delete (optional) (default to "off")

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.DeleteClusterStatefulSet(context.Background(), clusterId, id).Force(force).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.DeleteClusterStatefulSet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteClusterStatefulSet`: RestartClusterDeployment200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.DeleteClusterStatefulSet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteClusterStatefulSetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **force** | **string** | Force Delete | [default to &quot;off&quot;]

### Return type

[**RestartClusterDeployment200Response**](RestartClusterDeployment200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteClusterVolume

> DeleteClusterVolume200Response DeleteClusterVolume(ctx, clusterId, id).Force(force).Execute()

Delete a Volume



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster
	id := int64(1) // int64 | Morpheus ID of the Object being referenced
	force := "on" // string | Force Delete (optional) (default to "off")

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.DeleteClusterVolume(context.Background(), clusterId, id).Force(force).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.DeleteClusterVolume``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteClusterVolume`: DeleteClusterVolume200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.DeleteClusterVolume`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteClusterVolumeRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **force** | **string** | Force Delete | [default to &quot;off&quot;]

### Return type

[**DeleteClusterVolume200Response**](DeleteClusterVolume200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteClusterWorker

> DeleteClusterWorker200Response DeleteClusterWorker(ctx, clusterId, id).Force(force).Execute()

Delete a Worker



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster
	id := int64(1) // int64 | Morpheus ID of the Object being referenced
	force := "on" // string | Force Delete (optional) (default to "off")

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.DeleteClusterWorker(context.Background(), clusterId, id).Force(force).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.DeleteClusterWorker``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteClusterWorker`: DeleteClusterWorker200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.DeleteClusterWorker`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteClusterWorkerRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **force** | **string** | Force Delete | [default to &quot;off&quot;]

### Return type

[**DeleteClusterWorker200Response**](DeleteClusterWorker200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetCluster

> GetCluster200Response GetCluster(ctx, clusterId).Stats(stats).Execute()

Get a Specific Cluster



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster
	stats := true // bool | This can be used to include the cluster `stats` object in the response by passing `true`. This object includes usage information that varies by cluster type. For example for HVM it may include `cephStats` and `pacemakerStats`. (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.GetCluster(context.Background(), clusterId).Stats(stats).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.GetCluster``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetCluster`: GetCluster200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.GetCluster`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetClusterRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **stats** | **bool** | This can be used to include the cluster &#x60;stats&#x60; object in the response by passing &#x60;true&#x60;. This object includes usage information that varies by cluster type. For example for HVM it may include &#x60;cephStats&#x60; and &#x60;pacemakerStats&#x60;. | 

### Return type

[**GetCluster200Response**](GetCluster200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetClusterAffinityGroup

> GetClusterAffinityGroup200Response GetClusterAffinityGroup(ctx, clusterId, id).Execute()

Get a Specific Cluster Affinity Group



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster
	id := int64(1) // int64 | Morpheus ID of the Object being referenced

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.GetClusterAffinityGroup(context.Background(), clusterId, id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.GetClusterAffinityGroup``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetClusterAffinityGroup`: GetClusterAffinityGroup200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.GetClusterAffinityGroup`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetClusterAffinityGroupRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**GetClusterAffinityGroup200Response**](GetClusterAffinityGroup200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetClusterApiConfig

> GetClusterApiConfig200Response GetClusterApiConfig(ctx, clusterId).Execute()

Get API Config



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.GetClusterApiConfig(context.Background(), clusterId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.GetClusterApiConfig``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetClusterApiConfig`: GetClusterApiConfig200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.GetClusterApiConfig`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetClusterApiConfigRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetClusterApiConfig200Response**](GetClusterApiConfig200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetClusterConfigmap

> GetClusterDeployment200Response GetClusterConfigmap(ctx, clusterId, id).Execute()

Get a Specific Cluster Configmap



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster
	id := int64(1) // int64 | Morpheus ID of the Object being referenced

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.GetClusterConfigmap(context.Background(), clusterId, id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.GetClusterConfigmap``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetClusterConfigmap`: GetClusterDeployment200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.GetClusterConfigmap`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetClusterConfigmapRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**GetClusterDeployment200Response**](GetClusterDeployment200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetClusterContainer

> GetClusterContainer200Response GetClusterContainer(ctx, clusterId, id).Execute()

Get a Specific Cluster Container



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster
	id := int64(1) // int64 | Morpheus ID of the Object being referenced

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.GetClusterContainer(context.Background(), clusterId, id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.GetClusterContainer``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetClusterContainer`: GetClusterContainer200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.GetClusterContainer`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetClusterContainerRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**GetClusterContainer200Response**](GetClusterContainer200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetClusterDaemonset

> GetClusterDeployment200Response GetClusterDaemonset(ctx, clusterId, id).Execute()

Get a Specific Cluster Daemonset



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster
	id := int64(1) // int64 | Morpheus ID of the Object being referenced

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.GetClusterDaemonset(context.Background(), clusterId, id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.GetClusterDaemonset``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetClusterDaemonset`: GetClusterDeployment200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.GetClusterDaemonset`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetClusterDaemonsetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**GetClusterDeployment200Response**](GetClusterDeployment200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetClusterDatastore

> GetClusterDatastore200Response GetClusterDatastore(ctx, clusterId, id).Execute()

Get a Specific Cluster Datastore



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster
	id := int64(1) // int64 | Morpheus ID of the Object being referenced

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.GetClusterDatastore(context.Background(), clusterId, id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.GetClusterDatastore``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetClusterDatastore`: GetClusterDatastore200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.GetClusterDatastore`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetClusterDatastoreRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**GetClusterDatastore200Response**](GetClusterDatastore200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetClusterDeployment

> GetClusterDeployment200Response GetClusterDeployment(ctx, clusterId, id).Execute()

Get a Specific Cluster Deployment



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster
	id := int64(1) // int64 | Morpheus ID of the Object being referenced

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.GetClusterDeployment(context.Background(), clusterId, id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.GetClusterDeployment``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetClusterDeployment`: GetClusterDeployment200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.GetClusterDeployment`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetClusterDeploymentRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**GetClusterDeployment200Response**](GetClusterDeployment200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetClusterHistory

> GetClusterHistory200Response GetClusterHistory(ctx, clusterId).Execute()

Get Cluster History



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.GetClusterHistory(context.Background(), clusterId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.GetClusterHistory``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetClusterHistory`: GetClusterHistory200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.GetClusterHistory`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetClusterHistoryRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetClusterHistory200Response**](GetClusterHistory200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetClusterHistoryDetail

> GetClusterHistoryDetail200Response GetClusterHistoryDetail(ctx, clusterId, id).Execute()

Get Cluster History Details



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster
	id := int64(1) // int64 | Morpheus ID of the Object being referenced

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.GetClusterHistoryDetail(context.Background(), clusterId, id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.GetClusterHistoryDetail``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetClusterHistoryDetail`: GetClusterHistoryDetail200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.GetClusterHistoryDetail`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetClusterHistoryDetailRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**GetClusterHistoryDetail200Response**](GetClusterHistoryDetail200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetClusterHistoryEventDetail

> GetClusterHistoryEventDetail200Response GetClusterHistoryEventDetail(ctx, clusterId, id).Execute()

Get Cluster History Event



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster
	id := int64(1) // int64 | Morpheus ID of the Object being referenced

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.GetClusterHistoryEventDetail(context.Background(), clusterId, id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.GetClusterHistoryEventDetail``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetClusterHistoryEventDetail`: GetClusterHistoryEventDetail200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.GetClusterHistoryEventDetail`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetClusterHistoryEventDetailRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**GetClusterHistoryEventDetail200Response**](GetClusterHistoryEventDetail200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetClusterIngress

> GetClusterDeployment200Response GetClusterIngress(ctx, clusterId, id).Execute()

Get a Specific Cluster Ingress



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster
	id := int64(1) // int64 | Morpheus ID of the Object being referenced

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.GetClusterIngress(context.Background(), clusterId, id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.GetClusterIngress``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetClusterIngress`: GetClusterDeployment200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.GetClusterIngress`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetClusterIngressRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**GetClusterDeployment200Response**](GetClusterDeployment200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetClusterJob

> GetClusterDeployment200Response GetClusterJob(ctx, clusterId, id).Execute()

Get a Specific Cluster Job



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster
	id := int64(1) // int64 | Morpheus ID of the Object being referenced

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.GetClusterJob(context.Background(), clusterId, id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.GetClusterJob``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetClusterJob`: GetClusterDeployment200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.GetClusterJob`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetClusterJobRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**GetClusterDeployment200Response**](GetClusterDeployment200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetClusterMasters

> GetClusterMasters200Response GetClusterMasters(ctx, clusterId).Phrase(phrase).Execute()

Get Masters (Kubernetes)



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster
	phrase := "phrase_example" // string | Search phrase for partial matches on name or description (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.GetClusterMasters(context.Background(), clusterId).Phrase(phrase).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.GetClusterMasters``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetClusterMasters`: GetClusterMasters200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.GetClusterMasters`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetClusterMastersRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **phrase** | **string** | Search phrase for partial matches on name or description | 

### Return type

[**GetClusterMasters200Response**](GetClusterMasters200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetClusterNamespace

> GetClusterNamespace200Response GetClusterNamespace(ctx, clusterId, id).Execute()

Get Namespace (Kubernetes)



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster
	id := int64(1) // int64 | Morpheus ID of the Object being referenced

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.GetClusterNamespace(context.Background(), clusterId, id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.GetClusterNamespace``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetClusterNamespace`: GetClusterNamespace200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.GetClusterNamespace`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetClusterNamespaceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**GetClusterNamespace200Response**](GetClusterNamespace200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetClusterNamespaces

> GetClusterNamespaces200Response GetClusterNamespaces(ctx, clusterId).Execute()

List Namespaces (Kubernetes)



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.GetClusterNamespaces(context.Background(), clusterId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.GetClusterNamespaces``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetClusterNamespaces`: GetClusterNamespaces200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.GetClusterNamespaces`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetClusterNamespacesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetClusterNamespaces200Response**](GetClusterNamespaces200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetClusterNetworkEndpoing

> GetClusterDeployment200Response GetClusterNetworkEndpoing(ctx, clusterId, id).Execute()

Get a Specific Cluster Endpoint



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster
	id := int64(1) // int64 | Morpheus ID of the Object being referenced

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.GetClusterNetworkEndpoing(context.Background(), clusterId, id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.GetClusterNetworkEndpoing``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetClusterNetworkEndpoing`: GetClusterDeployment200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.GetClusterNetworkEndpoing`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetClusterNetworkEndpoingRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**GetClusterDeployment200Response**](GetClusterDeployment200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetClusterPod

> GetClusterDeployment200Response GetClusterPod(ctx, clusterId, id).Execute()

Get a Specific Pod for a Cluster



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster
	id := int64(1) // int64 | Morpheus ID of the Object being referenced

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.GetClusterPod(context.Background(), clusterId, id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.GetClusterPod``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetClusterPod`: GetClusterDeployment200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.GetClusterPod`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetClusterPodRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**GetClusterDeployment200Response**](GetClusterDeployment200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetClusterPolicy

> GetClusterDeployment200Response GetClusterPolicy(ctx, clusterId, id).Execute()

Get a Specific Cluster Policy



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster
	id := int64(1) // int64 | Morpheus ID of the Object being referenced

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.GetClusterPolicy(context.Background(), clusterId, id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.GetClusterPolicy``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetClusterPolicy`: GetClusterDeployment200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.GetClusterPolicy`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetClusterPolicyRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**GetClusterDeployment200Response**](GetClusterDeployment200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetClusterReplicaset

> GetClusterDeployment200Response GetClusterReplicaset(ctx, clusterId, id).Execute()

Get a Specific Cluster Replicaset



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster
	id := int64(1) // int64 | Morpheus ID of the Object being referenced

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.GetClusterReplicaset(context.Background(), clusterId, id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.GetClusterReplicaset``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetClusterReplicaset`: GetClusterDeployment200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.GetClusterReplicaset`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetClusterReplicasetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**GetClusterDeployment200Response**](GetClusterDeployment200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetClusterSecret

> GetClusterDeployment200Response GetClusterSecret(ctx, clusterId, id).Execute()

Get a Specific Cluster Secret



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster
	id := int64(1) // int64 | Morpheus ID of the Object being referenced

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.GetClusterSecret(context.Background(), clusterId, id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.GetClusterSecret``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetClusterSecret`: GetClusterDeployment200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.GetClusterSecret`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetClusterSecretRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**GetClusterDeployment200Response**](GetClusterDeployment200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetClusterService

> GetClusterDeployment200Response GetClusterService(ctx, clusterId, id).Execute()

Get a Specific Cluster Service



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster
	id := int64(1) // int64 | Morpheus ID of the Object being referenced

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.GetClusterService(context.Background(), clusterId, id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.GetClusterService``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetClusterService`: GetClusterDeployment200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.GetClusterService`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetClusterServiceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**GetClusterDeployment200Response**](GetClusterDeployment200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetClusterStatefulset

> GetClusterDeployment200Response GetClusterStatefulset(ctx, clusterId, id).Execute()

Get a Specific Cluster Statefulset



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster
	id := int64(1) // int64 | Morpheus ID of the Object being referenced

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.GetClusterStatefulset(context.Background(), clusterId, id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.GetClusterStatefulset``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetClusterStatefulset`: GetClusterDeployment200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.GetClusterStatefulset`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetClusterStatefulsetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**GetClusterDeployment200Response**](GetClusterDeployment200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetClusterUpgradeVersions

> GetClusterUpgradeVersions200Response GetClusterUpgradeVersions(ctx, clusterId).Execute()

Get Cluster Upgrade Versions (Kubernetes)



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.GetClusterUpgradeVersions(context.Background(), clusterId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.GetClusterUpgradeVersions``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetClusterUpgradeVersions`: GetClusterUpgradeVersions200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.GetClusterUpgradeVersions`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetClusterUpgradeVersionsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetClusterUpgradeVersions200Response**](GetClusterUpgradeVersions200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetClusterVolume

> GetClusterDeployment200Response GetClusterVolume(ctx, clusterId, id).Execute()

Get a Specific Cluster Volume



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster
	id := int64(1) // int64 | Morpheus ID of the Object being referenced

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.GetClusterVolume(context.Background(), clusterId, id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.GetClusterVolume``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetClusterVolume`: GetClusterDeployment200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.GetClusterVolume`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetClusterVolumeRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**GetClusterDeployment200Response**](GetClusterDeployment200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetClusterVolumeclaim

> GetClusterDeployment200Response GetClusterVolumeclaim(ctx, clusterId, id).Execute()

Get a Specific Cluster VolumeClaim



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster
	id := int64(1) // int64 | Morpheus ID of the Object being referenced

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.GetClusterVolumeclaim(context.Background(), clusterId, id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.GetClusterVolumeclaim``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetClusterVolumeclaim`: GetClusterDeployment200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.GetClusterVolumeclaim`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetClusterVolumeclaimRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**GetClusterDeployment200Response**](GetClusterDeployment200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListClusterAffinityGroups

> ListClusterAffinityGroups200Response ListClusterAffinityGroups(ctx, clusterId).Max(max).Offset(offset).Sort(sort).Order(order).Phrase(phrase).Name(name).Code(code).HideInactive(hideInactive).Execute()

Get Cluster Affinity Groups



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster
	max := int64(789) // int64 | Maximum number of records to return (optional) (default to 25)
	offset := int64(789) // int64 | Offset records, the number of records to skip, for paginating requests (optional) (default to 0)
	sort := "sort_example" // string | Sort order, the name of the property to sort by (optional) (default to "name")
	order := "order_example" // string | Sort direction, use 'desc' to reverse sort (optional) (default to "asc")
	phrase := "phrase_example" // string | Search phrase for partial matches on name or description (optional)
	name := "example" // string | Filter by name (optional)
	code := "azr" // string | If specified will return an exact match on code (optional)
	hideInactive := true // bool | If true restricts query to only load active datastores (optional) (default to false)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.ListClusterAffinityGroups(context.Background(), clusterId).Max(max).Offset(offset).Sort(sort).Order(order).Phrase(phrase).Name(name).Code(code).HideInactive(hideInactive).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.ListClusterAffinityGroups``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListClusterAffinityGroups`: ListClusterAffinityGroups200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.ListClusterAffinityGroups`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 

### Other Parameters

Other parameters are passed through a pointer to a apiListClusterAffinityGroupsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **order** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **name** | **string** | Filter by name | 
 **code** | **string** | If specified will return an exact match on code | 
 **hideInactive** | **bool** | If true restricts query to only load active datastores | [default to false]

### Return type

[**ListClusterAffinityGroups200Response**](ListClusterAffinityGroups200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListClusterConfigmaps

> ListClusterConfigmaps200Response ListClusterConfigmaps(ctx, clusterId).Execute()

Get Configmaps



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.ListClusterConfigmaps(context.Background(), clusterId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.ListClusterConfigmaps``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListClusterConfigmaps`: ListClusterConfigmaps200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.ListClusterConfigmaps`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 

### Other Parameters

Other parameters are passed through a pointer to a apiListClusterConfigmapsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**ListClusterConfigmaps200Response**](ListClusterConfigmaps200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListClusterContainers

> ListClusterContainers200Response ListClusterContainers(ctx, clusterId).Max(max).Offset(offset).Sort(sort).Order(order).Phrase(phrase).ResourceLevel(resourceLevel).Execute()

Get Containers for a Cluster



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster
	max := int64(789) // int64 | Maximum number of records to return (optional) (default to 25)
	offset := int64(789) // int64 | Offset records, the number of records to skip, for paginating requests (optional) (default to 0)
	sort := "sort_example" // string | Sort order, the name of the property to sort by (optional)
	order := "order_example" // string | Sort direction, use 'desc' to reverse sort (optional) (default to "asc")
	phrase := "phrase_example" // string | Search phrase for partial matches on name or description (optional)
	resourceLevel := "resourceLevel_example" // string | Resource level filter (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.ListClusterContainers(context.Background(), clusterId).Max(max).Offset(offset).Sort(sort).Order(order).Phrase(phrase).ResourceLevel(resourceLevel).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.ListClusterContainers``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListClusterContainers`: ListClusterContainers200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.ListClusterContainers`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 

### Other Parameters

Other parameters are passed through a pointer to a apiListClusterContainersRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | 
 **order** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **resourceLevel** | **string** | Resource level filter | 

### Return type

[**ListClusterContainers200Response**](ListClusterContainers200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListClusterDaemonsets

> ListClusterDaemonsets200Response ListClusterDaemonsets(ctx, clusterId).Execute()

Get Daemonsets



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.ListClusterDaemonsets(context.Background(), clusterId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.ListClusterDaemonsets``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListClusterDaemonsets`: ListClusterDaemonsets200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.ListClusterDaemonsets`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 

### Other Parameters

Other parameters are passed through a pointer to a apiListClusterDaemonsetsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**ListClusterDaemonsets200Response**](ListClusterDaemonsets200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListClusterDatastores

> ListClusterDatastores200Response ListClusterDatastores(ctx, clusterId).Max(max).Offset(offset).Sort(sort).Order(order).Phrase(phrase).Name(name).Code(code).HideInactive(hideInactive).Execute()

Get Cluster Datastores



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster
	max := int64(789) // int64 | Maximum number of records to return (optional) (default to 25)
	offset := int64(789) // int64 | Offset records, the number of records to skip, for paginating requests (optional) (default to 0)
	sort := "sort_example" // string | Sort order, the name of the property to sort by (optional) (default to "name")
	order := "order_example" // string | Sort direction, use 'desc' to reverse sort (optional) (default to "asc")
	phrase := "phrase_example" // string | Search phrase for partial matches on name or description (optional)
	name := "example" // string | Filter by name (optional)
	code := "azr" // string | If specified will return an exact match on code (optional)
	hideInactive := true // bool | If true restricts query to only load active datastores (optional) (default to false)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.ListClusterDatastores(context.Background(), clusterId).Max(max).Offset(offset).Sort(sort).Order(order).Phrase(phrase).Name(name).Code(code).HideInactive(hideInactive).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.ListClusterDatastores``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListClusterDatastores`: ListClusterDatastores200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.ListClusterDatastores`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 

### Other Parameters

Other parameters are passed through a pointer to a apiListClusterDatastoresRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **order** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **name** | **string** | Filter by name | 
 **code** | **string** | If specified will return an exact match on code | 
 **hideInactive** | **bool** | If true restricts query to only load active datastores | [default to false]

### Return type

[**ListClusterDatastores200Response**](ListClusterDatastores200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListClusterDeployments

> ListClusterDeployments200Response ListClusterDeployments(ctx, clusterId).Max(max).Offset(offset).Sort(sort).Order(order).Phrase(phrase).ResourceLevel(resourceLevel).Execute()

Get Deployments



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster
	max := int64(789) // int64 | Maximum number of records to return (optional) (default to 25)
	offset := int64(789) // int64 | Offset records, the number of records to skip, for paginating requests (optional) (default to 0)
	sort := "sort_example" // string | Sort order, the name of the property to sort by (optional) (default to "name")
	order := "order_example" // string | Sort direction, use 'desc' to reverse sort (optional) (default to "asc")
	phrase := "phrase_example" // string | Search phrase for partial matches on name or description (optional)
	resourceLevel := "resourceLevel_example" // string | Resource level filter (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.ListClusterDeployments(context.Background(), clusterId).Max(max).Offset(offset).Sort(sort).Order(order).Phrase(phrase).ResourceLevel(resourceLevel).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.ListClusterDeployments``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListClusterDeployments`: ListClusterDeployments200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.ListClusterDeployments`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 

### Other Parameters

Other parameters are passed through a pointer to a apiListClusterDeploymentsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **order** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **resourceLevel** | **string** | Resource level filter | 

### Return type

[**ListClusterDeployments200Response**](ListClusterDeployments200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListClusterIngresses

> ListClusterIngresses200Response ListClusterIngresses(ctx, clusterId).Execute()

Get Ingresses



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.ListClusterIngresses(context.Background(), clusterId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.ListClusterIngresses``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListClusterIngresses`: ListClusterIngresses200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.ListClusterIngresses`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 

### Other Parameters

Other parameters are passed through a pointer to a apiListClusterIngressesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**ListClusterIngresses200Response**](ListClusterIngresses200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListClusterJobs

> ListClusterJobs200Response ListClusterJobs(ctx, clusterId).Max(max).Offset(offset).Sort(sort).Order(order).Phrase(phrase).Execute()

Get Jobs



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster
	max := int64(789) // int64 | Maximum number of records to return (optional) (default to 25)
	offset := int64(789) // int64 | Offset records, the number of records to skip, for paginating requests (optional) (default to 0)
	sort := "sort_example" // string | Sort order, the name of the property to sort by (optional) (default to "name")
	order := "order_example" // string | Sort direction, use 'desc' to reverse sort (optional) (default to "asc")
	phrase := "phrase_example" // string | Search phrase for partial matches on name or description (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.ListClusterJobs(context.Background(), clusterId).Max(max).Offset(offset).Sort(sort).Order(order).Phrase(phrase).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.ListClusterJobs``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListClusterJobs`: ListClusterJobs200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.ListClusterJobs`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 

### Other Parameters

Other parameters are passed through a pointer to a apiListClusterJobsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **order** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **phrase** | **string** | Search phrase for partial matches on name or description | 

### Return type

[**ListClusterJobs200Response**](ListClusterJobs200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListClusterNetworkEndpoints

> ListClusterNetworkEndpoints200Response ListClusterNetworkEndpoints(ctx, clusterId).Execute()

Get Endpoints



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.ListClusterNetworkEndpoints(context.Background(), clusterId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.ListClusterNetworkEndpoints``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListClusterNetworkEndpoints`: ListClusterNetworkEndpoints200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.ListClusterNetworkEndpoints`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 

### Other Parameters

Other parameters are passed through a pointer to a apiListClusterNetworkEndpointsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**ListClusterNetworkEndpoints200Response**](ListClusterNetworkEndpoints200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListClusterPods

> ListClusterPods200Response ListClusterPods(ctx, clusterId).Max(max).Offset(offset).Sort(sort).Order(order).Phrase(phrase).ResourceLevel(resourceLevel).Execute()

Get Pods



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster
	max := int64(789) // int64 | Maximum number of records to return (optional) (default to 25)
	offset := int64(789) // int64 | Offset records, the number of records to skip, for paginating requests (optional) (default to 0)
	sort := "sort_example" // string | Sort order, the name of the property to sort by (optional) (default to "name")
	order := "order_example" // string | Sort direction, use 'desc' to reverse sort (optional) (default to "asc")
	phrase := "phrase_example" // string | Search phrase for partial matches on name or description (optional)
	resourceLevel := "resourceLevel_example" // string | Resource level filter (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.ListClusterPods(context.Background(), clusterId).Max(max).Offset(offset).Sort(sort).Order(order).Phrase(phrase).ResourceLevel(resourceLevel).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.ListClusterPods``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListClusterPods`: ListClusterPods200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.ListClusterPods`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 

### Other Parameters

Other parameters are passed through a pointer to a apiListClusterPodsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **order** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **resourceLevel** | **string** | Resource level filter | 

### Return type

[**ListClusterPods200Response**](ListClusterPods200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListClusterPolicies

> ListClusterPolicies200Response ListClusterPolicies(ctx, clusterId).Execute()

Get Policies



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.ListClusterPolicies(context.Background(), clusterId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.ListClusterPolicies``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListClusterPolicies`: ListClusterPolicies200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.ListClusterPolicies`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 

### Other Parameters

Other parameters are passed through a pointer to a apiListClusterPoliciesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**ListClusterPolicies200Response**](ListClusterPolicies200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListClusterReplicasets

> ListClusterReplicasets200Response ListClusterReplicasets(ctx, clusterId).Execute()

Get Replicasets



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.ListClusterReplicasets(context.Background(), clusterId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.ListClusterReplicasets``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListClusterReplicasets`: ListClusterReplicasets200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.ListClusterReplicasets`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 

### Other Parameters

Other parameters are passed through a pointer to a apiListClusterReplicasetsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**ListClusterReplicasets200Response**](ListClusterReplicasets200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListClusterSecrets

> ListClusterSecrets200Response ListClusterSecrets(ctx, clusterId).Execute()

Get Secrets



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.ListClusterSecrets(context.Background(), clusterId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.ListClusterSecrets``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListClusterSecrets`: ListClusterSecrets200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.ListClusterSecrets`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 

### Other Parameters

Other parameters are passed through a pointer to a apiListClusterSecretsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**ListClusterSecrets200Response**](ListClusterSecrets200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListClusterServices

> ListClusterServices200Response ListClusterServices(ctx, clusterId).Max(max).Offset(offset).Sort(sort).Order(order).Phrase(phrase).Execute()

Get Services



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster
	max := int64(789) // int64 | Maximum number of records to return (optional) (default to 25)
	offset := int64(789) // int64 | Offset records, the number of records to skip, for paginating requests (optional) (default to 0)
	sort := "sort_example" // string | Sort order, the name of the property to sort by (optional) (default to "name")
	order := "order_example" // string | Sort direction, use 'desc' to reverse sort (optional) (default to "asc")
	phrase := "phrase_example" // string | Search phrase for partial matches on name or description (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.ListClusterServices(context.Background(), clusterId).Max(max).Offset(offset).Sort(sort).Order(order).Phrase(phrase).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.ListClusterServices``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListClusterServices`: ListClusterServices200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.ListClusterServices`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 

### Other Parameters

Other parameters are passed through a pointer to a apiListClusterServicesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **order** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **phrase** | **string** | Search phrase for partial matches on name or description | 

### Return type

[**ListClusterServices200Response**](ListClusterServices200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListClusterStatefulSets

> ListClusterStatefulSets200Response ListClusterStatefulSets(ctx, clusterId).Max(max).Offset(offset).Sort(sort).Order(order).Phrase(phrase).ResourceLevel(resourceLevel).Execute()

Get Stateful Sets



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster
	max := int64(789) // int64 | Maximum number of records to return (optional) (default to 25)
	offset := int64(789) // int64 | Offset records, the number of records to skip, for paginating requests (optional) (default to 0)
	sort := "sort_example" // string | Sort order, the name of the property to sort by (optional) (default to "name")
	order := "order_example" // string | Sort direction, use 'desc' to reverse sort (optional) (default to "asc")
	phrase := "phrase_example" // string | Search phrase for partial matches on name or description (optional)
	resourceLevel := "resourceLevel_example" // string | Resource level filter (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.ListClusterStatefulSets(context.Background(), clusterId).Max(max).Offset(offset).Sort(sort).Order(order).Phrase(phrase).ResourceLevel(resourceLevel).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.ListClusterStatefulSets``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListClusterStatefulSets`: ListClusterStatefulSets200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.ListClusterStatefulSets`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 

### Other Parameters

Other parameters are passed through a pointer to a apiListClusterStatefulSetsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **order** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **resourceLevel** | **string** | Resource level filter | 

### Return type

[**ListClusterStatefulSets200Response**](ListClusterStatefulSets200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListClusterTypes

> ListClusterTypes200Response ListClusterTypes(ctx).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Code(code).ProviderType(providerType).Execute()

Get All Cluster Types



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	max := int64(789) // int64 | Maximum number of records to return (optional) (default to 25)
	offset := int64(789) // int64 | Offset records, the number of records to skip, for paginating requests (optional) (default to 0)
	sort := "sort_example" // string | Sort order, the name of the property to sort by (optional) (default to "name")
	direction := "asc" // string | Sort direction, use 'desc' to reverse sort (optional) (default to "asc")
	phrase := "phrase_example" // string | Search phrase for partial matches on name or description (optional)
	name := "example" // string | Filter by name (optional)
	code := "azr" // string | If specified will return an exact match on code (optional)
	providerType := "providerType_example" // string | Filter by `Provider Type` code.  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.ListClusterTypes(context.Background()).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Code(code).ProviderType(providerType).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.ListClusterTypes``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListClusterTypes`: ListClusterTypes200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.ListClusterTypes`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListClusterTypesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **direction** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **name** | **string** | Filter by name | 
 **code** | **string** | If specified will return an exact match on code | 
 **providerType** | **string** | Filter by &#x60;Provider Type&#x60; code.  | 

### Return type

[**ListClusterTypes200Response**](ListClusterTypes200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListClusterVolumeclaims

> ListClusterVolumeclaims200Response ListClusterVolumeclaims(ctx, clusterId).Execute()

Get VolumeClaims



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.ListClusterVolumeclaims(context.Background(), clusterId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.ListClusterVolumeclaims``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListClusterVolumeclaims`: ListClusterVolumeclaims200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.ListClusterVolumeclaims`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 

### Other Parameters

Other parameters are passed through a pointer to a apiListClusterVolumeclaimsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**ListClusterVolumeclaims200Response**](ListClusterVolumeclaims200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListClusterVolumes

> ListClusterVolumes200Response ListClusterVolumes(ctx, clusterId).Execute()

Get Volumes



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.ListClusterVolumes(context.Background(), clusterId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.ListClusterVolumes``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListClusterVolumes`: ListClusterVolumes200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.ListClusterVolumes`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 

### Other Parameters

Other parameters are passed through a pointer to a apiListClusterVolumesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**ListClusterVolumes200Response**](ListClusterVolumes200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListClusterWorkers

> ListClusterWorkers200Response ListClusterWorkers(ctx, clusterId).Execute()

Get Workers



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.ListClusterWorkers(context.Background(), clusterId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.ListClusterWorkers``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListClusterWorkers`: ListClusterWorkers200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.ListClusterWorkers`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 

### Other Parameters

Other parameters are passed through a pointer to a apiListClusterWorkersRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**ListClusterWorkers200Response**](ListClusterWorkers200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListClusters

> ListClusters200Response ListClusters(ctx).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).ZoneId(zoneId).TypeId(typeId).Labels(labels).AllLabels(allLabels).Execute()

Get All Clusters



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	max := int64(789) // int64 | Maximum number of records to return (optional) (default to 25)
	offset := int64(789) // int64 | Offset records, the number of records to skip, for paginating requests (optional) (default to 0)
	sort := "sort_example" // string | Sort order, the name of the property to sort by (optional) (default to "name")
	direction := "asc" // string | Sort direction, use 'desc' to reverse sort (optional) (default to "asc")
	phrase := "phrase_example" // string | Search phrase for partial matches on name or description (optional)
	name := "example" // string | Filter by name (optional)
	zoneId := int64(3) // int64 | The Cloud ID (Zone ID) for Filtering (optional)
	typeId := int64(3) // int64 | Type filter, restricts query to only load clusters of a specified cluster type (optional)
	labels := "labels_example" // string | Filter by label(s), matches records that contain any of the specified labels (optional)
	allLabels := "allLabels_example" // string | Filter by label(s), matches records that contain all of the specified labels (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.ListClusters(context.Background()).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).ZoneId(zoneId).TypeId(typeId).Labels(labels).AllLabels(allLabels).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.ListClusters``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListClusters`: ListClusters200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.ListClusters`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListClustersRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **direction** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **name** | **string** | Filter by name | 
 **zoneId** | **int64** | The Cloud ID (Zone ID) for Filtering | 
 **typeId** | **int64** | Type filter, restricts query to only load clusters of a specified cluster type | 
 **labels** | **string** | Filter by label(s), matches records that contain any of the specified labels | 
 **allLabels** | **string** | Filter by label(s), matches records that contain all of the specified labels | 

### Return type

[**ListClusters200Response**](ListClusters200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RefreshCluster

> RefreshCluster200Response RefreshCluster(ctx, clusterId).Execute()

Refreshes a Cluster



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.RefreshCluster(context.Background(), clusterId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.RefreshCluster``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RefreshCluster`: RefreshCluster200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.RefreshCluster`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 

### Other Parameters

Other parameters are passed through a pointer to a apiRefreshClusterRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**RefreshCluster200Response**](RefreshCluster200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RestartClusterContainer

> RestartClusterContainer200Response RestartClusterContainer(ctx, clusterId, id).Execute()

Restart a Container



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster
	id := int64(1) // int64 | Morpheus ID of the Object being referenced

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.RestartClusterContainer(context.Background(), clusterId, id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.RestartClusterContainer``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RestartClusterContainer`: RestartClusterContainer200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.RestartClusterContainer`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiRestartClusterContainerRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**RestartClusterContainer200Response**](RestartClusterContainer200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RestartClusterDeployment

> RestartClusterDeployment200Response RestartClusterDeployment(ctx, clusterId, id).Execute()

Restart a Deployment



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster
	id := int64(1) // int64 | Morpheus ID of the Object being referenced

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.RestartClusterDeployment(context.Background(), clusterId, id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.RestartClusterDeployment``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RestartClusterDeployment`: RestartClusterDeployment200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.RestartClusterDeployment`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiRestartClusterDeploymentRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**RestartClusterDeployment200Response**](RestartClusterDeployment200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RestartClusterPod

> RestartClusterDeployment200Response RestartClusterPod(ctx, clusterId, id).Execute()

Restart a Pod



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster
	id := int64(1) // int64 | Morpheus ID of the Object being referenced

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.RestartClusterPod(context.Background(), clusterId, id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.RestartClusterPod``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RestartClusterPod`: RestartClusterDeployment200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.RestartClusterPod`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiRestartClusterPodRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**RestartClusterDeployment200Response**](RestartClusterDeployment200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RestartClusterStatefulSet

> RestartClusterDeployment200Response RestartClusterStatefulSet(ctx, clusterId, id).Execute()

Restart a Stateful Set



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster
	id := int64(1) // int64 | Morpheus ID of the Object being referenced

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.RestartClusterStatefulSet(context.Background(), clusterId, id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.RestartClusterStatefulSet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RestartClusterStatefulSet`: RestartClusterDeployment200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.RestartClusterStatefulSet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiRestartClusterStatefulSetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**RestartClusterDeployment200Response**](RestartClusterDeployment200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## SaveClusterAffinityGroup

> SaveClusterAffinityGroup200Response SaveClusterAffinityGroup(ctx, clusterId).SaveClusterAffinityGroupRequest(saveClusterAffinityGroupRequest).Execute()

Create a Cluster Affinity Group



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster
	saveClusterAffinityGroupRequest := *openapiclient.NewSaveClusterAffinityGroupRequest() // SaveClusterAffinityGroupRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.SaveClusterAffinityGroup(context.Background(), clusterId).SaveClusterAffinityGroupRequest(saveClusterAffinityGroupRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.SaveClusterAffinityGroup``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `SaveClusterAffinityGroup`: SaveClusterAffinityGroup200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.SaveClusterAffinityGroup`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 

### Other Parameters

Other parameters are passed through a pointer to a apiSaveClusterAffinityGroupRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **saveClusterAffinityGroupRequest** | [**SaveClusterAffinityGroupRequest**](SaveClusterAffinityGroupRequest.md) |  | 

### Return type

[**SaveClusterAffinityGroup200Response**](SaveClusterAffinityGroup200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## SaveClusterDatastore

> SaveClusterDatastore200Response SaveClusterDatastore(ctx, clusterId).SaveClusterDatastoreRequest(saveClusterDatastoreRequest).Execute()

Create a Cluster Datastore



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster
	saveClusterDatastoreRequest := *openapiclient.NewSaveClusterDatastoreRequest() // SaveClusterDatastoreRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.SaveClusterDatastore(context.Background(), clusterId).SaveClusterDatastoreRequest(saveClusterDatastoreRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.SaveClusterDatastore``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `SaveClusterDatastore`: SaveClusterDatastore200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.SaveClusterDatastore`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 

### Other Parameters

Other parameters are passed through a pointer to a apiSaveClusterDatastoreRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **saveClusterDatastoreRequest** | [**SaveClusterDatastoreRequest**](SaveClusterDatastoreRequest.md) |  | 

### Return type

[**SaveClusterDatastore200Response**](SaveClusterDatastore200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateCluster

> UpdateCluster200Response UpdateCluster(ctx, clusterId).UpdateClusterRequest(updateClusterRequest).Execute()

Update Cluster



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster
	updateClusterRequest := *openapiclient.NewUpdateClusterRequest() // UpdateClusterRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.UpdateCluster(context.Background(), clusterId).UpdateClusterRequest(updateClusterRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.UpdateCluster``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateCluster`: UpdateCluster200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.UpdateCluster`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateClusterRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateClusterRequest** | [**UpdateClusterRequest**](UpdateClusterRequest.md) |  | 

### Return type

[**UpdateCluster200Response**](UpdateCluster200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateClusterAffinityGroup

> UpdateClusterAffinityGroup200Response UpdateClusterAffinityGroup(ctx, clusterId, id).UpdateClusterAffinityGroupRequest(updateClusterAffinityGroupRequest).Execute()

Update Cluster Affinity Group



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster
	id := int64(1) // int64 | Morpheus ID of the Object being referenced
	updateClusterAffinityGroupRequest := *openapiclient.NewUpdateClusterAffinityGroupRequest() // UpdateClusterAffinityGroupRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.UpdateClusterAffinityGroup(context.Background(), clusterId, id).UpdateClusterAffinityGroupRequest(updateClusterAffinityGroupRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.UpdateClusterAffinityGroup``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateClusterAffinityGroup`: UpdateClusterAffinityGroup200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.UpdateClusterAffinityGroup`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateClusterAffinityGroupRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **updateClusterAffinityGroupRequest** | [**UpdateClusterAffinityGroupRequest**](UpdateClusterAffinityGroupRequest.md) |  | 

### Return type

[**UpdateClusterAffinityGroup200Response**](UpdateClusterAffinityGroup200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateClusterAssignedLayout

> UpdateClusterAssignedLayout200Response UpdateClusterAssignedLayout(ctx, clusterId).UpdateClusterAssignedLayoutRequest(updateClusterAssignedLayoutRequest).Execute()

Update Cluster Assigned Layout



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster
	updateClusterAssignedLayoutRequest := *openapiclient.NewUpdateClusterAssignedLayoutRequest(int64(123)) // UpdateClusterAssignedLayoutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.UpdateClusterAssignedLayout(context.Background(), clusterId).UpdateClusterAssignedLayoutRequest(updateClusterAssignedLayoutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.UpdateClusterAssignedLayout``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateClusterAssignedLayout`: UpdateClusterAssignedLayout200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.UpdateClusterAssignedLayout`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateClusterAssignedLayoutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateClusterAssignedLayoutRequest** | [**UpdateClusterAssignedLayoutRequest**](UpdateClusterAssignedLayoutRequest.md) |  | 

### Return type

[**UpdateClusterAssignedLayout200Response**](UpdateClusterAssignedLayout200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateClusterDatastore

> UpdateClusterDatastore200Response UpdateClusterDatastore(ctx, clusterId, id).UpdateClusterDatastoreRequest(updateClusterDatastoreRequest).Execute()

Update Cluster Datastore



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster
	id := int64(1) // int64 | Morpheus ID of the Object being referenced
	updateClusterDatastoreRequest := *openapiclient.NewUpdateClusterDatastoreRequest() // UpdateClusterDatastoreRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.UpdateClusterDatastore(context.Background(), clusterId, id).UpdateClusterDatastoreRequest(updateClusterDatastoreRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.UpdateClusterDatastore``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateClusterDatastore`: UpdateClusterDatastore200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.UpdateClusterDatastore`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateClusterDatastoreRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **updateClusterDatastoreRequest** | [**UpdateClusterDatastoreRequest**](UpdateClusterDatastoreRequest.md) |  | 

### Return type

[**UpdateClusterDatastore200Response**](UpdateClusterDatastore200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateClusterNamespace

> UpdateClusterNamespace200Response UpdateClusterNamespace(ctx, clusterId, id).UpdateClusterNamespaceRequest(updateClusterNamespaceRequest).Execute()

Update Namespace (Kubernetes)



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster
	id := int64(1) // int64 | Morpheus ID of the Object being referenced
	updateClusterNamespaceRequest := *openapiclient.NewUpdateClusterNamespaceRequest() // UpdateClusterNamespaceRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.UpdateClusterNamespace(context.Background(), clusterId, id).UpdateClusterNamespaceRequest(updateClusterNamespaceRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.UpdateClusterNamespace``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateClusterNamespace`: UpdateClusterNamespace200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.UpdateClusterNamespace`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateClusterNamespaceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **updateClusterNamespaceRequest** | [**UpdateClusterNamespaceRequest**](UpdateClusterNamespaceRequest.md) |  | 

### Return type

[**UpdateClusterNamespace200Response**](UpdateClusterNamespace200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateClusterPermissions

> UpdateClusterPermissions200Response UpdateClusterPermissions(ctx, clusterId).UpdateClusterPermissionsRequest(updateClusterPermissionsRequest).Execute()

Update Cluster Permissions



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster
	updateClusterPermissionsRequest := *openapiclient.NewUpdateClusterPermissionsRequest(*openapiclient.NewUpdateClusterPermissionsRequestPermissions()) // UpdateClusterPermissionsRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.UpdateClusterPermissions(context.Background(), clusterId).UpdateClusterPermissionsRequest(updateClusterPermissionsRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.UpdateClusterPermissions``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateClusterPermissions`: UpdateClusterPermissions200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.UpdateClusterPermissions`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateClusterPermissionsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateClusterPermissionsRequest** | [**UpdateClusterPermissionsRequest**](UpdateClusterPermissionsRequest.md) |  | 

### Return type

[**UpdateClusterPermissions200Response**](UpdateClusterPermissions200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateClusterUpgradeVersions

> UpdateClusterUpgradeVersions200Response UpdateClusterUpgradeVersions(ctx, clusterId).TargetVersion(targetVersion).Execute()

Upgrade a Cluster (Kubernetes)



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster
	targetVersion := "1.21.14" // string | Target version for cluster after upgrade

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.UpdateClusterUpgradeVersions(context.Background(), clusterId).TargetVersion(targetVersion).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.UpdateClusterUpgradeVersions``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateClusterUpgradeVersions`: UpdateClusterUpgradeVersions200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.UpdateClusterUpgradeVersions`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateClusterUpgradeVersionsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **targetVersion** | **string** | Target version for cluster after upgrade | 

### Return type

[**UpdateClusterUpgradeVersions200Response**](UpdateClusterUpgradeVersions200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateClusterWorkerCount

> UpdateClusterWorkerCount200Response UpdateClusterWorkerCount(ctx, clusterId).WorkerCount(workerCount).Execute()

Update Worker Count



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clusterId := int64(5) // int64 | The ID of the cluster
	workerCount := int64(5) // int64 | The target number of worker nodes

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.UpdateClusterWorkerCount(context.Background(), clusterId).WorkerCount(workerCount).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.UpdateClusterWorkerCount``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateClusterWorkerCount`: UpdateClusterWorkerCount200Response
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.UpdateClusterWorkerCount`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **int64** | The ID of the cluster | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateClusterWorkerCountRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **workerCount** | **int64** | The target number of worker nodes | 

### Return type

[**UpdateClusterWorkerCount200Response**](UpdateClusterWorkerCount200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

