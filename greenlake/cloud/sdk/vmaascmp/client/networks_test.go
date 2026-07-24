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
	"testing"

	consts "github.com/HPE/terraform-provider-hpe/greenlake/cloud/sdk/vmaascmp/common"
	"github.com/HPE/terraform-provider-hpe/greenlake/cloud/sdk/vmaascmp/models"
	gomock "github.com/golang/mock/gomock"
)

func TestNetworksAPIService_GetAllNetworks(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	templateName := "test_template_get_all_networks"

	tests := []struct {
		name    string
		param   map[string]string
		given   func(m *MockAPIClientHandler)
		want    models.ListNetworksBody
		wantErr bool
	}{
		{
			name: "Normal Test case 1: Get all networks",
			param: map[string]string{
				"name": templateName,
			},
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/networks"
				method := "GET"
				headers := getDefaultHeaders()
				req, _ := http.NewRequest(method, path, nil)
				respBody := io.NopCloser(bytes.NewReader([]byte(`
					{
						"networks": [{
							"id": 1,
							"name": "test_template_get_all_networks"
						}],
						"networkCount": 1
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
			want: models.ListNetworksBody{
				Networks: []models.GetSpecificNetwork{
					{
						ID:   1,
						Name: templateName,
					},
				},
				NetworkCount: 1,
			},
			wantErr: false,
		},
		{
			name: "Failed Test case 2: Error in prepare request",
			param: map[string]string{
				"name": templateName,
			},
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/networks"
				method := "GET"
				headers := getDefaultHeaders()
				m.EXPECT().getVersion().Return(999999)
				m.EXPECT().prepareRequest(gomock.Any(), path, method, nil, headers,
					getURLValues(map[string]string{
						"name": templateName,
					}), url.Values{}, "", nil).Return(nil, errors.New("prepare error request"))
			},
			want:    models.ListNetworksBody{},
			wantErr: true,
		},
		{
			name: "Failed Test case 3: Error in callAPI",
			param: map[string]string{
				"name": templateName,
			},
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/networks"
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
			want:    models.ListNetworksBody{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPIClient := NewMockAPIClientHandler(ctrl)
			n := NetworksAPIService{
				Client: mockAPIClient,
				Cfg: Configuration{
					Host: mockHost,
				},
			}
			tt.given(mockAPIClient)
			got, err := n.GetAllNetworks(ctx, tt.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetworksAPIService.GetAllNetworks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NetworksAPIService.GetAllNetworks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetworksAPIService_CreateNetwork(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name    string
		args    models.CreateNetworkRequest
		given   func(m *MockAPIClientHandler)
		want    models.CreateNetworkResponse
		wantErr bool
	}{
		{
			name: "Normal test case 1: Create network",
			args: models.CreateNetworkRequest{
				Network: models.CreateNetwork{
					Name: "tf_net",
				},
			},
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/networks"
				method := "POST"
				headers := getDefaultHeaders()
				req, _ := http.NewRequest(method, path, nil)

				m.EXPECT().getVersion().Return(999999).MaxTimes(3)
				m.EXPECT().prepareRequest(gomock.Any(), path, method,
					models.CreateNetworkRequest{
						Network: models.CreateNetwork{
							Name: "tf_net",
						},
					},
					headers, url.Values{}, url.Values{}, "", nil).Return(req, nil)

				m.EXPECT().callAPI(req).Return(&http.Response{
					StatusCode: 200,
					Body: io.NopCloser(bytes.NewReader([]byte(`
					{
						"Success": true,
						"network": {
							"id": 16,
							"name": "tf_net"
						}
					}
				`))),
				}, nil)
			},
			want: models.CreateNetworkResponse{
				Success: true,
				Network: models.GetSpecificNetwork{
					ID:   16,
					Name: "tf_net",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPIClient := NewMockAPIClientHandler(ctrl)
			tt.given(mockAPIClient)
			n := NetworksAPIService{
				Client: mockAPIClient,
				Cfg: Configuration{
					Host: mockHost,
				},
			}
			got, err := n.CreateNetwork(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetworksAPIService.CreateNetwork() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NetworksAPIService.CreateNetwork() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetworksAPIService_UpdateNetwork(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		networkID int
		request   models.CreateNetworkRequest
	}
	tests := []struct {
		name    string
		args    args
		given   func(m *MockAPIClientHandler)
		want    models.SuccessOrErrorMessage
		wantErr bool
	}{
		{
			name: "Normal test case 1: update network",
			args: args{
				networkID: 1,
				request: models.CreateNetworkRequest{
					Network: models.CreateNetwork{
						Name: "tf_net",
					},
				},
			},
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/networks/1"
				method := "PUT"
				headers := getDefaultHeaders()
				req, _ := http.NewRequest(method, path, nil)

				m.EXPECT().getVersion().Return(999999)
				m.EXPECT().prepareRequest(gomock.Any(), path, method,
					models.CreateNetworkRequest{
						Network: models.CreateNetwork{
							Name: "tf_net",
						},
					},
					headers, url.Values{}, url.Values{}, "", nil).Return(req, nil)

				m.EXPECT().callAPI(req).Return(&http.Response{
					StatusCode: 200,
					Body: io.NopCloser(bytes.NewReader([]byte(`
					{
						"Success": true
					}
				`))),
				}, nil)
			},
			want: models.SuccessOrErrorMessage{
				Success: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPIClient := NewMockAPIClientHandler(ctrl)
			tt.given(mockAPIClient)
			n := NetworksAPIService{
				Client: mockAPIClient,
				Cfg: Configuration{
					Host: mockHost,
				},
			}
			got, err := n.UpdateNetwork(context.Background(), tt.args.networkID, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetworksAPIService.UpdateNetwork() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NetworksAPIService.UpdateNetwork() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetworksAPIService_GetNetworkProxy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name    string
		args    map[string]string
		given   func(m *MockAPIClientHandler)
		want    models.GetAllNetworkProxies
		wantErr bool
	}{
		{
			name: "Normal test case 1: Get all proxies",
			args: map[string]string{
				"name": "NSXT",
			},
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/networks/proxies"
				method := "GET"
				headers := getDefaultHeaders()
				req, _ := http.NewRequest(method, path, nil)
				respBody := io.NopCloser(bytes.NewReader([]byte(`
					{
						"networkProxies": [
							{
								"id": 1,
								"name": "test_proxy"
							}
						]
					}
				`)))
				m.EXPECT().getVersion().Return(999999)
				m.EXPECT().prepareRequest(gomock.Any(), path, method, nil, headers,
					getURLValues(map[string]string{
						"name": "NSXT",
					}), url.Values{}, "", nil).Return(req, nil)

				m.EXPECT().callAPI(req).Return(&http.Response{
					StatusCode: 200,
					Body:       respBody,
				}, nil)
			},
			want: models.GetAllNetworkProxies{
				GetNetworkProxies: []models.GetNetworkProxy{
					{
						ID:   1,
						Name: "test_proxy",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPIClient := NewMockAPIClientHandler(ctrl)
			tt.given(mockAPIClient)
			n := NetworksAPIService{
				Client: mockAPIClient,
				Cfg: Configuration{
					Host: mockHost,
				},
			}
			got, err := n.GetNetworkProxy(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetworksAPIService.GetNetworkProxy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NetworksAPIService.GetNetworkProxy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetworksAPIService_GetSpecificNetwork(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	templateName := "test_template_get_specific_network"
	tests := []struct {
		name      string
		networkID int
		given     func(m *MockAPIClientHandler)
		want      models.GetSpecificNetworkBody
		wantErr   bool
	}{
		// TODO: Add test cases.
		{
			name:      "Normal Test case 1: Get specific network",
			networkID: 1,
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/networks/1"
				method := "GET"
				headers := getDefaultHeaders()
				req, _ := http.NewRequest(method, path, nil)
				respBody := io.NopCloser(bytes.NewReader([]byte(`
					{
						"network": {
							"id": 1,
							"name": "test_template_get_specific_network"
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
			want: models.GetSpecificNetworkBody{
				Network: models.GetSpecificNetwork{
					ID:   1,
					Name: templateName,
				},
			},
			wantErr: false,
		},
		{
			name:      "Failed Test case 2: Error in prepare request",
			networkID: 1,
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/networks/1"
				method := "GET"
				headers := getDefaultHeaders()
				m.EXPECT().getVersion().Return(999999)
				m.EXPECT().prepareRequest(gomock.Any(), path, method, nil, headers,
					getURLValues(nil), url.Values{}, "", nil).Return(nil, errors.New("prepare error request"))
			},
			want:    models.GetSpecificNetworkBody{},
			wantErr: true,
		},
		{
			name:      "Failed Test case 3: Error in callAPI",
			networkID: 1,
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/networks/1"
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
			want:    models.GetSpecificNetworkBody{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPIClient := NewMockAPIClientHandler(ctrl)
			n := NetworksAPIService{
				Client: mockAPIClient,
				Cfg: Configuration{
					Host: mockHost,
				},
			}
			tt.given(mockAPIClient)
			got, err := n.GetSpecificNetwork(ctx, tt.networkID)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetworksAPIService.GetSpecificNetwork() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NetworksAPIService.GetSpecificNetwork() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetworksAPIService_DeleteNetwork(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	templateName := "test_template_delete_network"
	tests := []struct {
		name      string
		networkID int
		given     func(m *MockAPIClientHandler)
		want      models.SuccessOrErrorMessage
		wantErr   bool
	}{
		// TODO: Add test cases.
		{
			name:      "Normal Test case 1: Delete a network",
			networkID: 1,
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/networks/1"
				method := "DELETE"
				headers := getDefaultHeaders()
				req, _ := http.NewRequest(method, path, nil)
				respBody := io.NopCloser(bytes.NewReader([]byte(`
					{
						"success": true,
						"message": "test_template_delete_network"
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
			want: models.SuccessOrErrorMessage{
				Success: true,
				Message: templateName,
			},
			wantErr: false,
		},
		{
			name:      "Failed Test case 2: Error in prepare request",
			networkID: 1,
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/networks/1"
				method := "DELETE"
				headers := getDefaultHeaders()
				m.EXPECT().getVersion().Return(999999)
				m.EXPECT().prepareRequest(gomock.Any(), path, method, nil, headers,
					getURLValues(nil), url.Values{}, "", nil).Return(nil, errors.New("prepare error request"))
			},
			want:    models.SuccessOrErrorMessage{},
			wantErr: true,
		},
		{
			name:      "Failed Test case 3: Error in callAPI",
			networkID: 1,
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/networks/1"
				method := "DELETE"
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
			want:    models.SuccessOrErrorMessage{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPIClient := NewMockAPIClientHandler(ctrl)
			n := NetworksAPIService{
				Client: mockAPIClient,
				Cfg: Configuration{
					Host: mockHost,
				},
			}
			tt.given(mockAPIClient)
			got, err := n.DeleteNetwork(ctx, tt.networkID)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetworksAPIService.DeleteNetwork() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NetworksAPIService.DeleteNetwork() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetworksAPIService_GetNetworkType(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	templateName := "test_template_get_network_type"

	tests := []struct {
		name    string
		param   map[string]string
		given   func(m *MockAPIClientHandler)
		want    models.GetNetworkTypesResponse
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Normal Test case 1: Get network types",
			param: map[string]string{
				"name": templateName,
			},
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/network-types"
				method := "GET"
				headers := getDefaultHeaders()
				req, _ := http.NewRequest(method, path, nil)
				respBody := io.NopCloser(bytes.NewReader([]byte(`
					{
						"networkTypes": [{
							"id": 1,
							"name": "test_template_get_network_type"
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
			want: models.GetNetworkTypesResponse{
				NetworkTypes: []models.GetSpecificNetworkType{
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
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/network-types"
				method := "GET"
				headers := getDefaultHeaders()
				m.EXPECT().getVersion().Return(999999)
				m.EXPECT().prepareRequest(gomock.Any(), path, method, nil, headers,
					getURLValues(map[string]string{
						"name": templateName,
					}), url.Values{}, "", nil).Return(nil, errors.New("prepare error request"))
			},
			want:    models.GetNetworkTypesResponse{},
			wantErr: true,
		},
		{
			name: "Failed Test case 3: Error in callAPI",
			param: map[string]string{
				"name": templateName,
			},
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/network-types"
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
			want:    models.GetNetworkTypesResponse{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPIClient := NewMockAPIClientHandler(ctrl)
			n := NetworksAPIService{
				Client: mockAPIClient,
				Cfg: Configuration{
					Host: mockHost,
				},
			}
			tt.given(mockAPIClient)
			got, err := n.GetNetworkType(ctx, tt.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetworksAPIService.GetNetworkType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NetworksAPIService.GetNetworkType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetworksAPIService_GetNetworkPool(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	templateName := "test_template_get_network_pool"

	tests := []struct {
		name    string
		param   map[string]string
		given   func(m *MockAPIClientHandler)
		want    models.GetNetworkPoolsResp
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Normal Test case 1: Get network pool",
			param: map[string]string{
				"name": templateName,
			},
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/networks/pools"
				method := "GET"
				headers := getDefaultHeaders()
				req, _ := http.NewRequest(method, path, nil)
				respBody := io.NopCloser(bytes.NewReader([]byte(`
					{
						"networkPools": [{
							"id": 1,
							"name": "test_template_get_network_pool"
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
			want: models.GetNetworkPoolsResp{
				NetworkPools: []models.GetNetworkPool{
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
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/networks/pools"
				method := "GET"
				headers := getDefaultHeaders()
				m.EXPECT().getVersion().Return(999999)
				m.EXPECT().prepareRequest(gomock.Any(), path, method, nil, headers,
					getURLValues(map[string]string{
						"name": templateName,
					}), url.Values{}, "", nil).Return(nil, errors.New("prepare error request"))
			},
			want:    models.GetNetworkPoolsResp{},
			wantErr: true,
		},
		{
			name: "Failed Test case 3: Error in callAPI",
			param: map[string]string{
				"name": templateName,
			},
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/networks/pools"
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
			want:    models.GetNetworkPoolsResp{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPIClient := NewMockAPIClientHandler(ctrl)
			n := NetworksAPIService{
				Client: mockAPIClient,
				Cfg: Configuration{
					Host: mockHost,
				},
			}
			tt.given(mockAPIClient)
			got, err := n.GetNetworkPool(ctx, tt.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetworksAPIService.GetNetworkPool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NetworksAPIService.GetNetworkPool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetworksAPIService_GetSpecificNetworkPool(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	templateName := "test_template_get_specific_network"
	tests := []struct {
		name          string
		networkPoolID int
		given         func(m *MockAPIClientHandler)
		want          models.GetSpecificNetworkPool
		wantErr       bool
	}{
		// TODO: Add test cases.
		{
			name:          "Normal Test case 1: Get specific network",
			networkPoolID: 1,
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/networks/pools/1"
				method := "GET"
				headers := getDefaultHeaders()
				req, _ := http.NewRequest(method, path, nil)
				respBody := io.NopCloser(bytes.NewReader([]byte(`
					{
						"networkPool": {
							"id": 1,
							"name": "test_template_get_specific_network"
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
			want: models.GetSpecificNetworkPool{
				NetworkPool: models.GetNetworkPool{
					ID:   1,
					Name: templateName,
				},
			},
			wantErr: false,
		},
		{
			name:          "Failed Test case 2: Error in prepare request",
			networkPoolID: 1,
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/networks/pools/1"
				method := "GET"
				headers := getDefaultHeaders()
				m.EXPECT().getVersion().Return(999999)
				m.EXPECT().prepareRequest(gomock.Any(), path, method, nil, headers,
					getURLValues(nil), url.Values{}, "", nil).Return(nil, errors.New("prepare error request"))
			},
			want:    models.GetSpecificNetworkPool{},
			wantErr: true,
		},
		{
			name:          "Failed Test case 3: Error in callAPI",
			networkPoolID: 1,
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/networks/pools/1"
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
			want:    models.GetSpecificNetworkPool{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPIClient := NewMockAPIClientHandler(ctrl)
			n := NetworksAPIService{
				Client: mockAPIClient,
				Cfg: Configuration{
					Host: mockHost,
				},
			}
			tt.given(mockAPIClient)
			got, err := n.GetSpecificNetworkPool(ctx, tt.networkPoolID)
			if (err != nil) != tt.wantErr {
				t.Errorf("NetworksAPIService.GetSpecificNetworkPool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NetworksAPIService.GetSpecificNetworkPool() = %v, want %v", got, tt.want)
			}
		})
	}
}
