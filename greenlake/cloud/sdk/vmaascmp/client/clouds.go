// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package client

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	consts "github.com/HPE/terraform-provider-hpe/greenlake/cloud/sdk/vmaascmp/common"
	"github.com/HPE/terraform-provider-hpe/greenlake/cloud/sdk/vmaascmp/models"
)

type CloudsAPIService struct {
	Client APIClientHandler
	Cfg    Configuration
}

/*
CloudsAPIService
Get All Cloud Data Stores
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc.
    Passed from http.Request or context.Background().
  - @param serviceInstanceID
  - @param cloudID The cloud ID
*/
func (a *CloudsAPIService) GetAllCloudDataStores(ctx context.Context, cloudID int,
	queryParams map[string]string) (models.DataStoresResp, error) {
	AllCloudDSResp := models.DataStoresResp{}

	allCloudDSAPI := &api{
		method: "GET",
		path: fmt.Sprintf("%s/%d/%s",
			consts.ZonePath, cloudID, consts.DatstorePath),
		client: a.Client,

		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &AllCloudDSResp)
		},
		validations: []validationFunc{
			func() error {
				if cloudID < 1 {
					return fmt.Errorf("%s", "cloud id should be greater than or equal to 1")
				}

				return nil
			},
		},
	}
	err := allCloudDSAPI.do(ctx, nil, queryParams)

	return AllCloudDSResp, err
}

/*
CloudsAPIService
Get All Cloud Resource Pools
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc.
    Passed from http.Request or context.Background().
  - @param serviceInstanceID
  - @param cloudID The cloud ID
*/
func (a *CloudsAPIService) GetAllCloudResourcePools(ctx context.Context, cloudID int,
	queryParams map[string]string) (models.ResourcePoolsResp, error) {
	AllCloudRPoolResp := models.ResourcePoolsResp{}

	allCloudRPoolAPI := &api{
		method: "GET",
		path: fmt.Sprintf("%s/%d/%s",
			consts.ZonePath, cloudID, consts.ResourcePoolPath),
		client: a.Client,

		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &AllCloudRPoolResp)
		},
		validations: []validationFunc{
			func() error {
				if cloudID < 1 {
					return fmt.Errorf("%s", "cloud id should be greater than or equal to 1")
				}

				return nil
			},
		},
	}
	err := allCloudRPoolAPI.do(ctx, nil, queryParams)

	return AllCloudRPoolResp, err
}

/*
CloudsAPIService
Get All Clouds
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc.
    Passed from http.Request or context.Background().
  - @param serviceInstanceID
*/
func (a *CloudsAPIService) GetAllClouds(ctx context.Context,
	queryParams map[string]string) (models.CloudsResp, error) {
	AllCloudResp := models.CloudsResp{}

	allCloudAPI := &api{
		method: "GET",
		path:   consts.ZonePath,
		client: a.Client,

		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &AllCloudResp)
		},
	}
	err := allCloudAPI.do(ctx, nil, queryParams)

	return AllCloudResp, err
}

func (a *CloudsAPIService) GetAllCloudNetworks(ctx context.Context, cloudID,
	provisionTypeID int) (models.GetAllCloudNetworks, error) {
	AllCloudNetworkResp := models.GetAllCloudNetworks{}

	allCloudNetworkAPI := &api{
		method: "GET",
		path: fmt.Sprintf("%s/%s",
			consts.OptionsPath, consts.ZoneNetworkOptionsPath),
		client: a.Client,

		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &AllCloudNetworkResp)
		},
	}
	queryParams := map[string]string{
		"zoneId":          strconv.Itoa(cloudID),
		"provisionTypeId": strconv.Itoa(provisionTypeID),
	}
	err := allCloudNetworkAPI.do(ctx, nil, queryParams)

	return AllCloudNetworkResp, err
}

func (a *CloudsAPIService) GetAllCloudFolders(
	ctx context.Context,
	cloudID int,
	queryParams map[string]string,
) (models.GetAllCloudFolders, error) {
	Folders := models.GetAllCloudFolders{}

	folderAPI := &api{
		method: "GET",
		path: fmt.Sprintf("%s/%d/%s",
			consts.ZonePath, cloudID, consts.FolderPath),
		client: a.Client,
		// jsonParser also stores folder response, since folders is not a local variable
		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &Folders)
		},
		// validate cloud id > 1
		validations: []validationFunc{
			func() error {
				if cloudID < 1 {
					return fmt.Errorf("%s", "cloud id should be greater than or equal to 1")
				}

				return nil
			},
		},
	}

	err := folderAPI.do(ctx, nil, queryParams)

	return Folders, err
}

func (a *CloudsAPIService) GetSpecificCloudFolder(
	ctx context.Context,
	cloudID int,
	folderID int,
) (models.GetSpecificCloudFolder, error) {
	Folder := models.GetSpecificCloudFolder{}

	folderAPI := &api{
		method: "GET",
		path: fmt.Sprintf("%s/%d/%s/%d",
			consts.ZonePath, cloudID, consts.FolderPath, folderID),
		client: a.Client,

		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &Folder)
		},
		// validate cloud id > 1
		validations: []validationFunc{
			func() error {
				if cloudID < 1 {
					return fmt.Errorf("%s", "cloud id should be greater than or equal to 1")
				}

				return nil
			},
		},
	}

	err := folderAPI.do(ctx, nil, nil)

	return Folder, err
}

func (a *CloudsAPIService) GetSpecificCloud(
	ctx context.Context,
	cloudID int,
) (models.GetSpecificCloud, error) {
	Folder := models.GetSpecificCloud{}

	folderAPI := &api{
		method: "GET",
		path: fmt.Sprintf("%s/%d",
			consts.ZonePath, cloudID),
		client: a.Client,

		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &Folder)
		},
	}
	err := folderAPI.do(ctx, nil, nil)

	return Folder, err
}

func (a *CloudsAPIService) GetSpecificCloudDataStores(
	ctx context.Context,
	cloudID, datastoreID int,
) (models.GetSpecificCloudDataStores, error) {
	Folder := models.GetSpecificCloudDataStores{}

	folderAPI := &api{
		method: "GET",
		path: fmt.Sprintf("%s/%d/%s/%d",
			consts.ZonePath, cloudID, consts.DatstorePath, datastoreID),
		client: a.Client,

		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &Folder)
		},
	}
	err := folderAPI.do(ctx, nil, nil)

	return Folder, err
}

func (a *CloudsAPIService) GetSpecificCloudResourcePool(
	ctx context.Context,
	cloudID int,
	resourcePoolID int,
) (models.GetSpecificCloudResourcePool, error) {
	Resp := models.GetSpecificCloudResourcePool{}

	api := &api{
		method: "GET",
		path: fmt.Sprintf("%s/%d/%s/%d",
			consts.ZonePath, cloudID, consts.ResourcePoolPath, resourcePoolID),
		client: a.Client,
		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &Resp)
		},
	}
	err := api.do(ctx, nil, nil)

	return Resp, err
}
