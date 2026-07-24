// (C) Copyright 2021-2025 Hewlett Packard Enterprise Development LP

package client

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/antihax/optional"

	consts "github.com/HPE/terraform-provider-hpe/greenlake/cloud/sdk/vmaascmp/common"
	models "github.com/HPE/terraform-provider-hpe/greenlake/cloud/sdk/vmaascmp/models"
)

type InstancesAPIService struct {
	Client APIClientHandler
	Cfg    Configuration
}

/*
InstancesAPIService
Clone an instance and all VM within that instance.
 * @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc.
 	Passed from http.Request or context.Background().
 * @param serviceInstanceID
 * @param instanceID
 * @param optional nil or *InstancesAPICloneAnInstanceOpts - Optional Parameters:
     * @param "Body" (optional.Interface of CloneInstanceBody) -

*/

type InstancesAPICloneAnInstanceOpts struct {
	Body optional.Interface
}

/*
InstancesAPIService
Create an Instance
 * @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc.
 	Passed from http.Request or context.Background().
 * @param serviceInstanceID
 * @param optional nil or *InstancesAPICreateAnInstanceOpts - Optional Parameters:
     * @param "Body" (optional.Interface of CreateInstanceBody) -
@return models.GetInstanceResponse
*/

func (a *InstancesAPIService) CreateAnInstance(ctx context.Context,
	localVarOptionals *models.CreateInstanceBody) (models.GetInstanceResponse, error) {
	CreateInstanceResp := models.GetInstanceResponse{}

	// Pre-pending 'pool-' to ResourcePoolId in 6.0.3 and above
	if v, _ := parseVersion("6.0.3"); v <= a.Client.getVersion() {
		// To check if the key exists
		localVarOptionalsValue := reflect.ValueOf(localVarOptionals).Elem()
		if field := localVarOptionalsValue.FieldByName("Config"); field.IsValid() && !field.IsZero() {
			configValue := reflect.ValueOf(localVarOptionals.Config).Elem()
			if configField := configValue.FieldByName("ResourcePoolID"); configField.IsValid() && !configField.IsZero() {
				localVarOptionals.Config.ResourcePoolID = fmt.Sprintf("pool-%v",
					localVarOptionals.Config.ResourcePoolID)
			}
		}
	}
	createInstanceAPI := &api{
		method: "POST",
		path:   consts.InstancesPath,
		client: a.Client,

		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &CreateInstanceResp)
		},
	}

	err := createInstanceAPI.do(ctx, localVarOptionals, nil)

	return CreateInstanceResp, err
}

/*
InstancesAPIService
Will delete an instance and all associated monitors and backups.
 * @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc.
 	Passed from http.Request or context.Background().
 * @param serviceInstanceID
 * @param instanceID
 * @param optional nil or *InstancesAPIDeleteAnIstanceOpts - Optional Parameters:
     * @param "Force" (optional.String) -

*/

type InstancesAPIDeleteAnIstanceOpts struct {
	Force optional.String
}

func (a *InstancesAPIService) DeleteAnInstance(ctx context.Context,
	instanceID int) (models.SuccessOrErrorMessage, error) {
	DelInstanceResp := models.SuccessOrErrorMessage{}

	delInstanceAPI := &api{
		method: "DELETE",
		path: fmt.Sprintf("%s/%d",
			consts.InstancesPath, instanceID),
		client: a.Client,

		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &DelInstanceResp)
		},
	}
	err := delInstanceAPI.do(ctx, nil, nil)

	return DelInstanceResp, err
}

/*
InstancesAPIService
Get a Specific Instance
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc.
    Passed from http.Request or context.Background().
  - @param serviceInstanceID
  - @param instanceID

@return models.GetInstanceResponse
*/
func (a *InstancesAPIService) GetASpecificInstance(ctx context.Context,
	instanceID int) (models.GetInstanceResponse, error) {
	SpecificInstResp := models.GetInstanceResponse{}

	specificInstanceAPI := &api{
		method: "GET",
		path: fmt.Sprintf("%s/%d",
			consts.InstancesPath, instanceID),
		client: a.Client,

		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &SpecificInstResp)
		},
	}
	err := specificInstanceAPI.do(ctx, nil, nil)

	return SpecificInstResp, err
}

/*
InstancesAPIService
Get All Instances
 * @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc.
 	Passed from http.Request or context.Background().
 * @param serviceInstanceID
*/

func (a *InstancesAPIService) GetAllInstances(ctx context.Context,
	queryParams map[string]string) (models.Instances, error) {
	GetAllInstance := models.Instances{}

	instanceAPI := &api{
		method: "GET",
		path:   consts.InstancesPath,
		client: a.Client,

		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &GetAllInstance)
		},
	}
	err := instanceAPI.do(ctx, nil, queryParams)

	return GetAllInstance, err
}

/*
InstancesAPIService
Lists VMware Snapshot of the instance
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc.
    Passed from http.Request or context.Background().
  - @param serviceInstanceID
  - @param instanceID
*/
func (a *InstancesAPIService) GetListOfSnapshotsForAnInstance(ctx context.Context,
	instanceID int) (models.ListSnapshotResponse, error) {
	ListSnapshotResp := models.ListSnapshotResponse{}

	listSnapshotAPI := &api{
		method: "GET",
		path: fmt.Sprintf("%s/%d/snapshots",
			consts.InstancesPath, instanceID),
		client: a.Client,

		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &ListSnapshotResp)
		},
	}
	err := listSnapshotAPI.do(ctx, nil, nil)

	return ListSnapshotResp, err
}

/*
InstancesAPIService
It is possible to import a snapshot of an instance. This creates a Virtual Image of the instance as
it currently exists.
 * @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc.
 	Passed from http.Request or context.Background().
 * @param serviceInstanceID
 * @param instanceID
 * @param optional nil or *InstancesAPIImportSnapshotOfAnInstanceOpts - Optional Parameters:
     * @param "Body" (optional.Interface of ImportSnapshotBody) -

*/

func (a *InstancesAPIService) ImportSnapshotOfAnInstance(ctx context.Context, instanceID int,
	localVarOptionals *models.ImportSnapshotBody) (models.SuccessOrErrorMessage, error) {
	ImportSnapshotResp := models.SuccessOrErrorMessage{}

	importSnapshotAPI := &api{
		method: "PUT",
		path: fmt.Sprintf("%s/%d/import-snapshot",
			consts.InstancesPath, instanceID),
		client: a.Client,

		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &ImportSnapshotResp)
		},
	}
	err := importSnapshotAPI.do(ctx, localVarOptionals, nil)

	return ImportSnapshotResp, err
}

/*
InstancesAPIService
Restarts all VM running within an instance
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc.
    Passed from http.Request or context.Background().
  - @param serviceInstanceID
  - @param instanceID
*/
func (a *InstancesAPIService) RestartAnInstance(ctx context.Context,
	instanceID int) (models.InstancePowerResponse, error) {
	RestartInstResp := models.InstancePowerResponse{}

	restartInstAPI := &api{
		method: "PUT",
		path: fmt.Sprintf("%s/%d/restart",
			consts.InstancesPath, instanceID),
		client: a.Client,

		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &RestartInstResp)
		},
	}
	err := restartInstAPI.do(ctx, nil, nil)

	return RestartInstResp, err
}

/*
InstancesAPIService
Creates VMware Snapshot of the instance
 * @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc.
 	Passed from http.Request or context.Background().
 * @param serviceInstanceID
 * @param instanceID
 * @param optional nil or *InstancesAPISnapshotAnInstanceOpts - Optional Parameters:
     * @param "Body" (optional.Interface of SnapshotBody) -

*/

func (a *InstancesAPIService) SnapshotAnInstance(ctx context.Context, instanceID int,
	localVarOptionals *models.SnapshotBody) (models.Instances, error) {
	SnapshotInstResp := models.Instances{}

	instanceAPI := &api{
		method: "PUT",
		path: fmt.Sprintf("%s/%d/snapshot",
			consts.InstancesPath, instanceID),
		client: a.Client,

		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &SnapshotInstResp)
		},
	}
	err := instanceAPI.do(ctx, localVarOptionals, nil)

	return SnapshotInstResp, err
}

/*
InstancesAPIService
Starts all VM running within an instance
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc.
    Passed from http.Request or context.Background().
  - @param serviceInstanceID
  - @param instanceID
*/
func (a *InstancesAPIService) StartAnInstance(ctx context.Context,
	instanceID int) (models.InstancePowerResponse, error) {
	StartInstanceResp := models.InstancePowerResponse{}

	startInstanceAPI := &api{
		method: "PUT",
		path: fmt.Sprintf("%s/%d/start",
			consts.InstancesPath, instanceID),
		client: a.Client,

		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &StartInstanceResp)
		},
	}
	err := startInstanceAPI.do(ctx, nil, nil)

	return StartInstanceResp, err
}

/*
InstancesAPIService
Stops all VM running within an instance
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc.
    Passed from http.Request or context.Background().
  - @param serviceInstanceID
  - @param instanceID
*/
func (a *InstancesAPIService) StopAnInstance(ctx context.Context,
	instanceID int) (models.InstancePowerResponse, error) {
	StopInstanceResp := models.InstancePowerResponse{}

	stopInstanceAPI := &api{
		method: "PUT",
		path: fmt.Sprintf("%s/%d/stop",
			consts.InstancesPath, instanceID),
		client: a.Client,

		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &StopInstanceResp)
		},
	}
	err := stopInstanceAPI.do(ctx, nil, nil)

	return StopInstanceResp, err
}

/*
InstancesAPIService
Suspends all VM running within an instance
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc.
    Passed from http.Request or context.Background().
  - @param serviceInstanceID
  - @param instanceID
*/
func (a *InstancesAPIService) SuspendAnInstance(ctx context.Context,
	instanceID int) (models.InstancePowerResponse, error) {
	SuspendResp := models.InstancePowerResponse{}
	suspendInstanceAPI := &api{
		method: "PUT",
		path: fmt.Sprintf("%s/%d/suspend",
			consts.InstancesPath, instanceID),
		client: a.Client,

		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &SuspendResp)
		},
	}
	err := suspendInstanceAPI.do(ctx, nil, nil)

	return SuspendResp, err
}

func (a *InstancesAPIService) ResizeAnInstance(ctx context.Context, instanceID int,
	request *models.ResizeInstanceBody) (models.ResizeInstanceResponse, error) {
	ResizeResp := models.ResizeInstanceResponse{}

	instanceAPI := &api{
		method: "PUT",
		path: fmt.Sprintf("%s/%d/resize",
			consts.InstancesPath, instanceID),
		client: a.Client,

		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &ResizeResp)
		},
	}
	err := instanceAPI.do(ctx, request, nil)

	return ResizeResp, err
}

func (a *InstancesAPIService) UpdatingAnInstance(
	ctx context.Context,
	instanceID int,
	request *models.UpdateInstanceBody,
) (models.UpdateInstanceResponse, error) {
	Instance := models.UpdateInstanceResponse{}

	instanceAPI := &api{
		method: "PUT",
		path: fmt.Sprintf("%s/%d",
			consts.InstancesPath, instanceID),
		client: a.Client,

		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &Instance)
		},
	}
	err := instanceAPI.do(ctx, request, nil)

	return Instance, err
}

func (a *InstancesAPIService) GetInstanceHistory(
	ctx context.Context, instanceID int) (models.GetInstanceHistory, error) {
	History := models.GetInstanceHistory{}

	historyAPI := &api{
		method: "GET",
		path: fmt.Sprintf("%s/%d/history",
			consts.InstancesPath, instanceID),
		client: a.Client,

		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &History)
		},
	}
	err := historyAPI.do(ctx, nil, nil)

	return History, err
}

func (a *InstancesAPIService) CloneAnInstance(ctx context.Context, instanceID int,
	cloneRequest models.CreateInstanceCloneBody) (models.SuccessOrErrorMessage, error) {
	var CloneResp models.SuccessOrErrorMessage
	if v, _ := parseVersion("5.2.12"); v <= a.Client.getVersion() {
		cloneRequest.Tags = cloneRequest.Metadata
		cloneRequest.Metadata = nil
		cloneRequest.Instance.Labels = cloneRequest.Instance.Tags
		cloneRequest.Instance.Tags = nil
	}
	instanceClone := &api{
		method: "PUT",
		path: fmt.Sprintf("%s/%d/clone",
			consts.InstancesPath, instanceID),
		client: a.Client,

		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &CloneResp)
		},
	}
	err := instanceClone.do(ctx, cloneRequest, nil)

	return CloneResp, err
}

func (a *InstancesAPIService) GetStorageVolTypeID(ctx context.Context, cloudID, layoutID, siteID string,
) (models.InstancePlanResponse, error) {
	StorageVol := models.InstancePlanResponse{}

	apiCaller := &api{
		method: "GET",
		path:   consts.InstancePlanPath,
		client: a.Client,

		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &StorageVol)
		},
	}
	queryParams := map[string]string{
		"zoneId":   cloudID,
		"layoutId": layoutID,
		"siteId":   siteID,
	}
	err := apiCaller.do(ctx, nil, queryParams)

	return StorageVol, err
}

func (a *InstancesAPIService) GetStorageControllerTypes(ctx context.Context, layoutID string,
) (models.StorageControllerTypesResp, error) {
	StorageControllers := models.StorageControllerTypesResp{}

	apiCaller := &api{
		method: "GET",
		path:   consts.StorageControllerTypesPath,
		client: a.Client,

		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &StorageControllers)
		},
	}
	queryParams := map[string]string{
		"layoutId": layoutID,
	}
	err := apiCaller.do(ctx, nil, queryParams)

	return StorageControllers, err
}
func (a *InstancesAPIService) GetStorageControllerMount(ctx context.Context, layoutID string, controllerName string,
	busNumber, unitNumber int) (ControllerMount string, err error) {
	controllerNameInput := strings.TrimSpace(strings.ToLower(controllerName))
	controllerRest, err := a.GetStorageControllerTypes(ctx, layoutID)
	for _, controller := range controllerRest.Data {
		controllerName := strings.TrimSpace(strings.ToLower(controller.Name))
		if controllerName == controllerNameInput {

			ControllerMount = fmt.Sprintf("%d:%d:%d:%d", -1, busNumber, controller.Value, unitNumber)
			break
		}
	}
	if ControllerMount == "" {
		err = fmt.Errorf("invalid controller type '%s'", controllerName)
	}
	return
}
