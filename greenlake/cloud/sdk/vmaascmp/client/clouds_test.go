//  (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package client

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"testing"

	consts "github.com/HPE/terraform-provider-hpe/greenlake/cloud/sdk/vmaascmp/common"
	"github.com/HPE/terraform-provider-hpe/greenlake/cloud/sdk/vmaascmp/models"
	"github.com/golang/mock/gomock"
)

func TestCloudsAPIService_GetAllCloudDataStores(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	templateName := "test_template_cloud_datastore"
	tests := []struct {
		name    string
		cloudID int
		param   map[string]string
		// All expectaion captures here
		given   func(m *MockAPIClientHandler)
		want    models.DataStoresResp
		wantErr bool
	}{
		{
			name: "Normal Test case 1: Get all cloud datastore",
			param: map[string]string{
				"name": templateName,
			},
			cloudID: 1,
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/zones/1/data-stores"
				method := "GET"
				headers := getDefaultHeaders()
				req, _ := http.NewRequest(method, path, nil)
				respBody := io.NopCloser(bytes.NewReader([]byte(`
					{
						"datastores": [{
							"id": 1,
							"name": "test_template_cloud_datastore"
						}]
					}
				`)))
				// mock the context only since it is not validated in this function
				m.EXPECT().getVersion().Return(999999)
				m.EXPECT().prepareRequest(gomock.Any(), path, method, nil, headers,
					getURLValues(map[string]string{
						"name": templateName,
					}), url.Values{}, "", nil).Return(req, nil)

				m.EXPECT().callAPI(req).Return(&http.Response{
					StatusCode: 200,
					Body:       respBody,
				}, nil)
			},
			want: models.DataStoresResp{
				Datastores: []models.DataStoresRespBody{
					{
						ID:   1,
						Name: templateName,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Failed Test case 2: error in prepare request",
			param: map[string]string{
				"name": templateName,
			},
			cloudID: 1,
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/zones/1/data-stores"
				method := "GET"
				headers := getDefaultHeaders()
				m.EXPECT().getVersion().Return(999999)
				m.EXPECT().prepareRequest(gomock.Any(), path, method, nil, headers,
					getURLValues(map[string]string{
						"name": templateName,
					}), url.Values{}, "", nil).Return(nil, errors.New("prepare request error"))
			},
			want:    models.DataStoresResp{},
			wantErr: true,
		},
		{
			name: "Failed Test case 3: Error in callAPI",
			param: map[string]string{
				"name": templateName,
			},
			cloudID: 1,
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/zones/1/data-stores"
				method := "GET"
				headers := getDefaultHeaders()
				req, _ := http.NewRequest(method, path, nil)
				respBody := io.NopCloser(bytes.NewReader([]byte(`
					{
						"message": "Internal Server Error",
						"recommendedActions": [
							"Unknown error occurred. Please contact the administrator"
						]
					}
				`)))
				// mock the context only since it is not validated in this function
				m.EXPECT().getVersion().Return(999999)
				m.EXPECT().prepareRequest(gomock.Any(), path, method, nil, headers,
					getURLValues(map[string]string{
						"name": templateName,
					}), url.Values{}, "", nil).Return(req, nil)

				m.EXPECT().callAPI(req).Return(&http.Response{
					StatusCode: 500,
					Body:       respBody,
				}, nil)
			},
			want:    models.DataStoresResp{},
			wantErr: true,
		},
		{
			name: "Failed Test case 4: Cloud ID can not be less than 1",
			param: map[string]string{
				"name": templateName,
			},
			cloudID: 0,
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				m.EXPECT().getVersion().Return(999999)
			},
			want:    models.DataStoresResp{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPIClient := NewMockAPIClientHandler(ctrl)
			a := CloudsAPIService{
				Client: mockAPIClient,
				Cfg: Configuration{
					Host: mockHost,
				},
			}
			tt.given(mockAPIClient)
			got, err := a.GetAllCloudDataStores(ctx, tt.cloudID, tt.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("CloudsAPIService.GetAllCloudDataStores() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CloudsAPIService.GetAllCloudDataStores() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCloudsAPIService_GetAllCloudResourcePools(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	templateName := "test_template_cloud_resource_pools"
	tests := []struct {
		name    string
		cloudID int
		param   map[string]string
		given   func(m *MockAPIClientHandler)
		want    models.ResourcePoolsResp
		wantErr bool
	}{
		{
			name: "Normal Test case 1: Get all cloud resource pool",
			param: map[string]string{
				"name": templateName,
			},
			cloudID: 1,
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/zones/1/resource-pools"
				method := "GET"
				headers := getDefaultHeaders()
				req, _ := http.NewRequest(method, path, nil)
				respBody := io.NopCloser(bytes.NewReader([]byte(`
					{
						"resourcePools": [{
							"id": 1,
							"name": "test_template_cloud_resource_pools"
						}]
					}
				`)))
				m.EXPECT().getVersion().Return(999999)
				m.EXPECT().prepareRequest(gomock.Any(), path, method, nil, headers,
					getURLValues(map[string]string{
						"name": templateName,
					}), url.Values{}, "", nil).Return(req, nil)

				m.EXPECT().callAPI(req).Return(&http.Response{
					StatusCode: 200,
					Body:       respBody,
				}, nil)
			},
			want: models.ResourcePoolsResp{
				ResourcePools: []models.ResourcePoolRespBody{
					{
						ID:   1,
						Name: templateName,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Failed Test case 2: Error in prepare request",
			param: map[string]string{
				"name": templateName,
			},
			cloudID: 1,
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/zones/1/resource-pools"
				method := "GET"
				headers := getDefaultHeaders()
				m.EXPECT().getVersion().Return(999999)
				m.EXPECT().prepareRequest(gomock.Any(), path, method, nil, headers,
					getURLValues(map[string]string{
						"name": templateName,
					}), url.Values{}, "", nil).Return(nil, errors.New("prepare request error"))
			},
			want:    models.ResourcePoolsResp{},
			wantErr: true,
		},
		{
			name: "Failed Test case 3: Error in callAPI",
			param: map[string]string{
				"name": templateName,
			},
			cloudID: 1,
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/zones/1/resource-pools"
				method := "GET"
				headers := getDefaultHeaders()
				req, _ := http.NewRequest(method, path, nil)
				respBody := io.NopCloser(bytes.NewReader([]byte(`
					{
						"message": "Internal Server Error",
						"recommendedActions": [
							"Unknown error occurred. Please contact the administrator"
						]
					}
				`)))
				m.EXPECT().getVersion().Return(999999)
				m.EXPECT().prepareRequest(gomock.Any(), path, method, nil, headers,
					getURLValues(map[string]string{
						"name": templateName,
					}), url.Values{}, "", nil).Return(req, nil)

				m.EXPECT().callAPI(req).Return(&http.Response{
					StatusCode: 500,
					Body:       respBody,
				}, nil)
			},
			want:    models.ResourcePoolsResp{},
			wantErr: true,
		},
		{
			name: "Failed Test case 4: Cloud ID can not be less than 1",
			param: map[string]string{
				"name": templateName,
			},
			cloudID: 0,
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				m.EXPECT().getVersion().Return(999999)
			},
			want:    models.ResourcePoolsResp{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPIClient := NewMockAPIClientHandler(ctrl)
			a := CloudsAPIService{
				Client: mockAPIClient,
				Cfg: Configuration{
					Host: mockHost,
				},
			}
			tt.given(mockAPIClient)
			got, err := a.GetAllCloudResourcePools(ctx, tt.cloudID, tt.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("CloudsAPIService.GetAllCloudResourcePools() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CloudsAPIService.GetAllCloudResourcePools() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCloudsAPIService_GetAllClouds(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	templateName := "test_template_clouds"
	tests := []struct {
		name    string
		param   map[string]string
		given   func(m *MockAPIClientHandler)
		want    models.CloudsResp
		wantErr bool
	}{
		{
			name: "Normal Test case 1: Get all clouds",
			param: map[string]string{
				"name": templateName,
			},
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/zones"
				method := "GET"
				headers := getDefaultHeaders()
				req, _ := http.NewRequest(method, path, nil)
				respBody := io.NopCloser(bytes.NewReader([]byte(`
					{
						"zones": [{
							"id": 1,
							"name": "test_template_clouds"
						}]
					}
				`)))
				m.EXPECT().getVersion().Return(999999)
				m.EXPECT().prepareRequest(gomock.Any(), path, method, nil, headers,
					getURLValues(map[string]string{
						"name": templateName,
					}), url.Values{}, "", nil).Return(req, nil)

				m.EXPECT().callAPI(req).Return(&http.Response{
					StatusCode: 200,
					Body:       respBody,
				}, nil)
			},
			want: models.CloudsResp{
				Clouds: []models.CloudRespBody{
					{
						ID:   1,
						Name: templateName,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Failed Test case 2: error in prepare request",
			param: map[string]string{
				"name": templateName,
			},
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/zones"
				method := "GET"
				headers := getDefaultHeaders()
				m.EXPECT().getVersion().Return(999999)
				m.EXPECT().prepareRequest(gomock.Any(), path, method, nil, headers,
					getURLValues(map[string]string{
						"name": templateName,
					}), url.Values{}, "", nil).Return(nil, errors.New("prepare request error"))
			},
			want:    models.CloudsResp{},
			wantErr: true,
		},
		{
			name: "Failed Test case 3: error in callAPI",
			param: map[string]string{
				"name": templateName,
			},
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/zones"
				method := "GET"
				headers := getDefaultHeaders()
				req, _ := http.NewRequest(method, path, nil)
				respBody := io.NopCloser(bytes.NewReader([]byte(`
					{
						"message": "Internal Server Error",
						"recommendedActions": [
							"Unknown error occurred. Please contact the administrator"
						]
					}
				`)))
				m.EXPECT().getVersion().Return(999999)
				m.EXPECT().prepareRequest(gomock.Any(), path, method, nil, headers,
					getURLValues(map[string]string{
						"name": templateName,
					}), url.Values{}, "", nil).Return(req, nil)

				m.EXPECT().callAPI(req).Return(&http.Response{
					StatusCode: 500,
					Body:       respBody,
				}, nil)
			},
			want:    models.CloudsResp{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPIClient := NewMockAPIClientHandler(ctrl)
			a := CloudsAPIService{
				Client: mockAPIClient,
				Cfg: Configuration{
					Host: mockHost,
				},
			}
			tt.given(mockAPIClient)
			got, err := a.GetAllClouds(ctx, tt.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("CloudsAPIService.GetAllClouds() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CloudsAPIService.GetAllClouds() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TODO add UTs for get all cloud folder(s)
func TestCloudsAPIService_GetAllCloudFolders(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	templateName := "test_template_get_all_folders"
	tests := []struct {
		name    string
		cloudID int
		param   map[string]string
		given   func(m *MockAPIClientHandler)
		want    models.GetAllCloudFolders
		wantErr bool
	}{
		{
			name: "Normal Test case 1: Get all folders",
			param: map[string]string{
				"name": templateName,
			},
			cloudID: 1,
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/zones/1/folders"
				method := "GET"
				headers := getDefaultHeaders()
				req, _ := http.NewRequest(method, path, nil)
				respBody := io.NopCloser(bytes.NewReader([]byte(`
					{
						"folders": [{
							"id": 1,
							"name": "test_template_get_all_folders"
						}]
					}
				`)))
				m.EXPECT().getVersion().Return(999999)
				m.EXPECT().prepareRequest(gomock.Any(), path, method, nil, headers,
					getURLValues(map[string]string{
						"name": templateName,
					}), url.Values{}, "", nil).Return(req, nil)

				m.EXPECT().callAPI(req).Return(&http.Response{
					StatusCode: 200,
					Body:       respBody,
				}, nil)
			},
			want: models.GetAllCloudFolders{
				Folders: []models.GetCloudFolder{
					{
						ID:   1,
						Name: templateName,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Failed Test case 2: Error in prepare request",
			param: map[string]string{
				"name": templateName,
			},
			cloudID: 1,
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/zones/1/folders"
				method := "GET"
				headers := getDefaultHeaders()
				m.EXPECT().getVersion().Return(999999)
				m.EXPECT().prepareRequest(gomock.Any(), path, method, nil, headers,
					getURLValues(map[string]string{
						"name": templateName,
					}), url.Values{}, "", nil).Return(nil, errors.New("prepare error request"))
			},
			want:    models.GetAllCloudFolders{},
			wantErr: true,
		},
		{
			name: "Failed Test case 3: Error in callAPI",
			param: map[string]string{
				"name": templateName,
			},
			cloudID: 1,
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/zones/1/folders"
				method := "GET"
				headers := getDefaultHeaders()
				req, _ := http.NewRequest(method, path, nil)
				respBody := io.NopCloser(bytes.NewReader([]byte(`
					{
						"message": "Internal Server Error",
						"recommendedActions": [
							"Unknown error occurred. Please contact the administrator"
						]
					}
				`)))
				m.EXPECT().getVersion().Return(999999)
				m.EXPECT().prepareRequest(gomock.Any(), path, method, nil, headers,
					getURLValues(map[string]string{
						"name": templateName,
					}), url.Values{}, "", nil).Return(req, nil)

				m.EXPECT().callAPI(req).Return(&http.Response{
					StatusCode: 500,
					Body:       respBody,
				}, nil)
			},
			want:    models.GetAllCloudFolders{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPIClient := NewMockAPIClientHandler(ctrl)
			a := CloudsAPIService{
				Client: mockAPIClient,
				Cfg: Configuration{
					Host: mockHost,
				},
			}
			tt.given(mockAPIClient)
			got, err := a.GetAllCloudFolders(ctx, tt.cloudID, tt.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("CloudsAPIService.GetAllCloudFolders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CloudsAPIService.GetAllCloudFolders() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCloudsAPIService_GetAllCloudNetworks(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	templateName := "test_template_all_cloud_networks"
	tests := []struct {
		name            string
		cloudID         int
		provisionTypeID int
		given           func(m *MockAPIClientHandler)
		want            models.GetAllCloudNetworks
		wantErr         bool
	}{
		{
			name:            "Normal Test case 1: Get all cloud networks",
			cloudID:         1,
			provisionTypeID: 1,
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/options/zoneNetworkOptions"
				method := "GET"
				headers := getDefaultHeaders()
				req, _ := http.NewRequest(method, path, nil)
				respBody := io.NopCloser(bytes.NewReader([]byte(`
					{
						"data": {
							"networkTypes": [{
								"id": 1,
								"name": "test_template_all_cloud_networks"
							}]
						}
					}
				`)))
				m.EXPECT().getVersion().Return(999999)
				m.EXPECT().prepareRequest(gomock.Any(), path, method, nil, headers,
					getURLValues(map[string]string{
						"zoneId":          strconv.Itoa(1),
						"provisionTypeId": strconv.Itoa(1),
					}), url.Values{}, "", nil).Return(req, nil)

				m.EXPECT().callAPI(req).Return(&http.Response{
					StatusCode: 200,
					Body:       respBody,
				}, nil)
			},
			want: models.GetAllCloudNetworks{
				Data: models.DataGetNetworkInterface{
					NetworkTypes: []models.GetNetworkInterfaceNetworkTypes{
						{
							ID:   1,
							Name: templateName,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:            "Failed Test case 2: Error in prepare request",
			cloudID:         1,
			provisionTypeID: 1,
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/options/zoneNetworkOptions"
				method := "GET"
				headers := getDefaultHeaders()
				m.EXPECT().getVersion().Return(999999)
				m.EXPECT().prepareRequest(gomock.Any(), path, method, nil, headers,
					getURLValues(map[string]string{
						"zoneId":          strconv.Itoa(1),
						"provisionTypeId": strconv.Itoa(1),
					}), url.Values{}, "", nil).Return(nil, errors.New("prepare request error"))
			},
			want:    models.GetAllCloudNetworks{},
			wantErr: true,
		},
		{
			name:            "Failed Test case 3: error in callAPI",
			cloudID:         1,
			provisionTypeID: 1,
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/options/zoneNetworkOptions"
				method := "GET"
				headers := getDefaultHeaders()
				req, _ := http.NewRequest(method, path, nil)
				respBody := io.NopCloser(bytes.NewReader([]byte(`
					{
						"message": "Internal Server Error",
						"recommendedActions": [
							"Unknown error occurred. Please contact the administrator"
						]
					}
				`)))
				m.EXPECT().getVersion().Return(999999)
				m.EXPECT().prepareRequest(gomock.Any(), path, method, nil, headers,
					getURLValues(map[string]string{
						"zoneId":          strconv.Itoa(1),
						"provisionTypeId": strconv.Itoa(1),
					}), url.Values{}, "", nil).Return(req, nil)

				m.EXPECT().callAPI(req).Return(&http.Response{
					StatusCode: 500,
					Body:       respBody,
				}, nil)
			},
			want:    models.GetAllCloudNetworks{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPIClient := NewMockAPIClientHandler(ctrl)
			a := CloudsAPIService{
				Client: mockAPIClient,
				Cfg: Configuration{
					Host: mockHost,
				},
			}
			tt.given(mockAPIClient)
			got, err := a.GetAllCloudNetworks(ctx, tt.cloudID, tt.provisionTypeID)
			if (err != nil) != tt.wantErr {
				t.Errorf("CloudsAPIService.GetAllCloudNetworks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CloudsAPIService.GetAllCloudNetworks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCloudsAPIService_GetSpecificCloudFolder(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	templateName := "test_template_all_specific_cloud_folder"
	tests := []struct {
		name     string
		cloudID  int
		folderID int
		given    func(m *MockAPIClientHandler)
		want     models.GetSpecificCloudFolder
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			name:     "Normal Test case 1: Get Specific Cloud Folder",
			cloudID:  1,
			folderID: 1,
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/zones/1/folders/1"
				method := "GET"
				headers := getDefaultHeaders()
				req, _ := http.NewRequest(method, path, nil)
				respBody := io.NopCloser(bytes.NewReader([]byte(`
					{
						"folder": {
								"id": 1,
								"name": "test_template_all_specific_cloud_folder"
							}
					}
				`)))
				m.EXPECT().getVersion().Return(999999)
				m.EXPECT().prepareRequest(gomock.Any(), path, method, nil, headers,
					getURLValues(nil), url.Values{}, "", nil).Return(req, nil)

				m.EXPECT().callAPI(req).Return(&http.Response{
					StatusCode: 200,
					Body:       respBody,
				}, nil)
			},
			want: models.GetSpecificCloudFolder{
				Folder: models.GetCloudFolder{
					ID:   1,
					Name: templateName,
				},
			},
			wantErr: false,
		},
		{
			name:     "Failed Test case 2: Error in prepare request",
			cloudID:  1,
			folderID: 1,
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/zones/1/folders/1"
				method := "GET"
				headers := getDefaultHeaders()
				m.EXPECT().getVersion().Return(999999)
				m.EXPECT().prepareRequest(gomock.Any(), path, method, nil, headers,
					getURLValues(nil), url.Values{}, "", nil).Return(nil, errors.New("prepare request error"))
			},
			want:    models.GetSpecificCloudFolder{},
			wantErr: true,
		},
		{
			name:     "Failed Test case 3: error in callAPI",
			cloudID:  1,
			folderID: 1,
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/zones/1/folders/1"
				method := "GET"
				headers := getDefaultHeaders()
				req, _ := http.NewRequest(method, path, nil)
				respBody := io.NopCloser(bytes.NewReader([]byte(`
					{
						"message": "Internal Server Error",
						"recommendedActions": [
							"Unknown error occurred. Please contact the administrator"
						]
					}
				`)))
				m.EXPECT().getVersion().Return(999999)
				m.EXPECT().prepareRequest(gomock.Any(), path, method, nil, headers,
					getURLValues(nil), url.Values{}, "", nil).Return(req, nil)

				m.EXPECT().callAPI(req).Return(&http.Response{
					StatusCode: 500,
					Body:       respBody,
				}, nil)
			},
			want:    models.GetSpecificCloudFolder{},
			wantErr: true,
		},
		{
			name:     "Failed Test case 4: Cloud ID can not be less than 1",
			cloudID:  0,
			folderID: 1,
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				m.EXPECT().getVersion().Return(999999)
			},
			want:    models.GetSpecificCloudFolder{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPIClient := NewMockAPIClientHandler(ctrl)
			a := CloudsAPIService{
				Client: mockAPIClient,
				Cfg: Configuration{
					Host: mockHost,
				},
			}
			tt.given(mockAPIClient)
			got, err := a.GetSpecificCloudFolder(ctx, tt.cloudID, tt.folderID)
			if (err != nil) != tt.wantErr {
				t.Errorf("CloudsAPIService.GetSpecificCloudFolder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CloudsAPIService.GetSpecificCloudFolder() = %v, want %v", got, tt.want)
			}
		})
	}
}
