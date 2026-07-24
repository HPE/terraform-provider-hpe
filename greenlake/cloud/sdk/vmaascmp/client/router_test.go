//  (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package client

import (
	"bytes"
	"context"
	"encoding/json"
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

func TestRouterAPIService_GetAllRouters(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	templateName := "test_template_get_all_routers"

	tests := []struct {
		name    string
		param   map[string]string
		given   func(m *MockAPIClientHandler)
		want    models.GetAllNetworkRouter
		wantErr bool
	}{
		{
			name: "Normal Test case 1: Get all Routers",
			param: map[string]string{
				"name": templateName,
			},
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/networks/routers"
				method := "GET"
				headers := getDefaultHeaders()
				req, _ := http.NewRequest(method, path, nil)
				respBody := io.NopCloser(bytes.NewReader([]byte(`
					{
						"networkRouters": [{
							"id": 1,
							"name": "test_template_get_all_routers"
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
			want: models.GetAllNetworkRouter{
				NetworkRouters: []models.GetNetworkRouter{
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
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/networks/routers"
				method := "GET"
				headers := getDefaultHeaders()
				m.EXPECT().getVersion().Return(999999)
				m.EXPECT().prepareRequest(gomock.Any(), path, method, nil, headers,
					getURLValues(map[string]string{
						"name": templateName,
					}), url.Values{}, "", nil).Return(nil, errors.New("prepare error request"))
			},
			want:    models.GetAllNetworkRouter{},
			wantErr: true,
		},
		{
			name: "Failed Test case 3: Error in callAPI",
			param: map[string]string{
				"name": templateName,
			},
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/networks/routers"
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
			want:    models.GetAllNetworkRouter{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPIClient := NewMockAPIClientHandler(ctrl)
			n := RouterAPIService{
				Client: mockAPIClient,
				Cfg: Configuration{
					Host: mockHost,
				},
			}
			tt.given(mockAPIClient)
			got, err := n.GetAllRouter(ctx, tt.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("RouterAPIService.GetAllRouter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RouterAPIService.GetAllRouter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouterAPIService_GetNetworkRouterTypes(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	templateName := "test_template_get_all_router_types"

	tests := []struct {
		name    string
		param   map[string]string
		given   func(m *MockAPIClientHandler)
		want    models.GetNetworlRouterTypes
		wantErr bool
	}{
		{
			name: "Normal Test case 1: Get all Router Types",
			param: map[string]string{
				"name": templateName,
			},
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/network-router-types"
				method := "GET"
				headers := getDefaultHeaders()
				req, _ := http.NewRequest(method, path, nil)
				respBody := io.NopCloser(bytes.NewReader([]byte(`
					{
						"networkRouterTypes": [{
							"id": 1,
							"name": "test_template_get_all_router_types"
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
			want: models.GetNetworlRouterTypes{
				NetworkRouterTypes: []models.NetworkRouterTypes{
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
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/network-router-types"
				method := "GET"
				headers := getDefaultHeaders()
				m.EXPECT().getVersion().Return(999999)
				m.EXPECT().prepareRequest(gomock.Any(), path, method, nil, headers,
					getURLValues(map[string]string{
						"name": templateName,
					}), url.Values{}, "", nil).Return(nil, errors.New("prepare error request"))
			},
			want:    models.GetNetworlRouterTypes{},
			wantErr: true,
		},
		{
			name: "Failed Test case 3: Error in callAPI",
			param: map[string]string{
				"name": templateName,
			},
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/network-router-types"
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
			want:    models.GetNetworlRouterTypes{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPIClient := NewMockAPIClientHandler(ctrl)
			n := RouterAPIService{
				Client: mockAPIClient,
				Cfg: Configuration{
					Host: mockHost,
				},
			}
			tt.given(mockAPIClient)
			got, err := n.GetRouterTypes(ctx, tt.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("RouterAPIService.GetRouterTypes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RouterAPIService.GetRouterTypes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouterAPIService_CreateRouter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name    string
		args    models.CreateRouterRequest
		given   func(m *MockAPIClientHandler)
		want    models.CreateRouterResp
		wantErr bool
	}{
		{
			name: "Normal test case 1: Create Router",
			args: models.CreateRouterRequest{
				NetworkRouter: models.CreateRouterRequestRouter{
					Name: "tf_router",
				},
			},
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/networks/routers"
				method := "POST"
				headers := getDefaultHeaders()
				req, _ := http.NewRequest(method, path, nil)

				m.EXPECT().getVersion().Return(999999)
				m.EXPECT().prepareRequest(gomock.Any(), path, method,
					models.CreateRouterRequest{
						NetworkRouter: models.CreateRouterRequestRouter{
							Name: "tf_router",
						},
					},
					headers, url.Values{}, url.Values{}, "", nil).Return(req, nil)

				m.EXPECT().callAPI(req).Return(&http.Response{
					StatusCode: 200,
					Body: io.NopCloser(bytes.NewReader([]byte(`
					{
						"success": true,
						"id": 16
					}
				`))),
				}, nil)
			},
			want: models.CreateRouterResp{
				Success: true,
				ID:      16,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPIClient := NewMockAPIClientHandler(ctrl)
			tt.given(mockAPIClient)
			n := RouterAPIService{
				Client: mockAPIClient,
				Cfg: Configuration{
					Host: mockHost,
				},
			}
			got, err := n.CreateRouter(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("RouterAPIService.CreateRouter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RouterAPIService.CreateRouter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouterAPIService_UpdateRouter(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	headers := getDefaultHeaders()
	defer ctrl.Finish()
	tests := []struct {
		name     string
		routerID int
		param    models.CreateRouterRequest
		given    func(m *MockAPIClientHandler)
		want     models.SuccessOrErrorMessage
		wantErr  bool
	}{
		{
			name:     "Normal Test case 1: Update a Router",
			routerID: 1,
			param: models.CreateRouterRequest{
				NetworkRouter: models.CreateRouterRequestRouter{
					Name: "test_update_router_name",
				},
			},
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/networks/routers/1"
				method := "PUT"
				postBody := io.NopCloser(bytes.NewReader([]byte(`
				{
					"networkRouter": {
						"name": "test_update_router_name"
					}
				}`)))
				req, _ := http.NewRequest(method, path, postBody)
				respBody := io.NopCloser(bytes.NewReader([]byte(`{
						"success": true
				}`)))
				pBody := models.CreateRouterRequest{
					NetworkRouter: models.CreateRouterRequestRouter{
						Name: "test_update_router_name",
					},
				}
				m.EXPECT().getVersion().Return(999999)
				m.EXPECT().prepareRequest(gomock.Any(), path, method, pBody, headers, url.Values{},
					url.Values{}, "", nil).Return(req, nil)
				m.EXPECT().callAPI(req).Return(&http.Response{
					StatusCode: 200,
					Body:       respBody,
				}, nil)
			},
			want: models.SuccessOrErrorMessage{
				Success: true,
			},
			wantErr: false,
		},
		{
			name:     "Failed test case 2: error in prepare request",
			routerID: 1,
			param: models.CreateRouterRequest{
				NetworkRouter: models.CreateRouterRequestRouter{
					Name: "test_update_router_name",
				},
			},
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/networks/routers/1"
				method := "PUT"

				pBody := models.CreateRouterRequest{
					NetworkRouter: models.CreateRouterRequestRouter{
						Name: "test_update_router_name",
					},
				}
				m.EXPECT().getVersion().Return(999999)
				m.EXPECT().prepareRequest(gomock.Any(), path, method, pBody, headers, url.Values{},
					url.Values{}, "", nil).Return(nil, errors.New("prepare request error"))
			},
			want:    models.SuccessOrErrorMessage{},
			wantErr: true,
		},
		{
			name:     "Failed test case 3: error in callAPI",
			routerID: 1,
			param: models.CreateRouterRequest{
				NetworkRouter: models.CreateRouterRequestRouter{
					Name: "test_update_router_name",
				},
			},
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/networks/routers/1"
				method := "PUT"
				postBody := io.NopCloser(bytes.NewReader([]byte(`
				{
					"networkRouter": {
						"name": "test_update_router_name"
					}
				}`)))
				req, _ := http.NewRequest(method, path, postBody)
				respBody := io.NopCloser(bytes.NewReader([]byte(`{
					{
						"message": "Internal Server Error",
						"recommendedActions": [
							"Unknown error occurred. Please contact the administrator"
						]
					}
				}`)))
				pBody := models.CreateRouterRequest{
					NetworkRouter: models.CreateRouterRequestRouter{
						Name: "test_update_router_name",
					},
				}
				m.EXPECT().getVersion().Return(999999)
				m.EXPECT().prepareRequest(gomock.Any(), path, method, pBody, headers, url.Values{},
					url.Values{}, "", nil).Return(req, nil)
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
			a := RouterAPIService{
				Client: mockAPIClient,
				Cfg: Configuration{
					Host: mockHost,
				},
			}
			tt.given(mockAPIClient)
			got, err := a.UpdateRouter(ctx, tt.routerID, tt.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("RouterAPIService.UpdateRouter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RouterAPIService.UpdateRouter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouterAPIService_GetSpecificRouter(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	templateName := "test_template_get_a_specific_router"
	tests := []struct {
		name     string
		routerID int
		given    func(m *MockAPIClientHandler)
		want     models.GetSpecificRouterResp
		wantErr  bool
	}{
		{
			name:     "Normal Test case 1: Get a specific router",
			routerID: 1,
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/networks/routers/1"
				method := "GET"
				headers := getDefaultHeaders()
				req, _ := http.NewRequest(method, path, nil)
				respBody := io.NopCloser(bytes.NewReader([]byte(`
					{
						"networkRouter":{
							"id": 1,
							"name": "test_template_get_a_specific_router"
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
			want: models.GetSpecificRouterResp{
				NetworkRouter: models.GetNetworkRouter{
					ID:   1,
					Name: templateName,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPIClient := NewMockAPIClientHandler(ctrl)
			r := RouterAPIService{
				Client: mockAPIClient,
				Cfg: Configuration{
					Host: mockHost,
				},
			}
			tt.given(mockAPIClient)
			got, err := r.GetSpecificRouter(ctx, tt.routerID)
			if (err != nil) != tt.wantErr {
				t.Errorf("RouterAPIService.GetSpecificRouter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RouterAPIService.GetSpecificRouter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouterAPIService_DeleteRouter(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	templateName := "test_template_delete_a_router"
	tests := []struct {
		name     string
		routerID int
		given    func(m *MockAPIClientHandler)
		want     models.SuccessOrErrorMessage
		wantErr  bool
	}{
		{
			name:     "Normal Test case 1: Delete a  router",
			routerID: 1,
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/networks/routers/1"
				method := "DELETE"
				headers := getDefaultHeaders()
				req, _ := http.NewRequest(method, path, nil)
				respBody := io.NopCloser(bytes.NewReader([]byte(`
					{
						"success": true,
						"message": "test_template_delete_a_router"
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
			name:     "Failed Test case 2: Error in prepare request",
			routerID: 1,
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/networks/routers/1"
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
			name:     "Failed Test case 3: Error in callAPI",
			routerID: 1,
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/networks/routers/1"
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
			r := RouterAPIService{
				Client: mockAPIClient,
				Cfg: Configuration{
					Host: mockHost,
				},
			}
			tt.given(mockAPIClient)
			got, err := r.DeleteRouter(ctx, tt.routerID)
			if (err != nil) != tt.wantErr {
				t.Errorf("RouterAPIService.DeleteRouter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RouterAPIService.DeleteRouter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouterAPIService_GetNetworkServices(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	templateName := "test_template_get_network_services"
	tests := []struct {
		name    string
		param   map[string]string
		given   func(m *MockAPIClientHandler)
		want    models.GetNetworkServicesResp
		wantErr bool
	}{
		{
			name: "Normal Test case 1: Get Network Services",
			param: map[string]string{
				"name": templateName,
			},
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/networks/services"
				method := "GET"
				headers := getDefaultHeaders()
				req, _ := http.NewRequest(method, path, nil)
				respBody := io.NopCloser(bytes.NewReader([]byte(`
					{
						"networkServices": [{
							"id": 1,
							"name": "test_template_get_network_services"
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
			want: models.GetNetworkServicesResp{
				NetworkServices: []models.GetNetworkServices{
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
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/networks/services"
				method := "GET"
				headers := getDefaultHeaders()
				m.EXPECT().getVersion().Return(999999)
				m.EXPECT().prepareRequest(gomock.Any(), path, method, nil, headers,
					getURLValues(map[string]string{
						"name": templateName,
					}), url.Values{}, "", nil).Return(nil, errors.New("prepare error request"))
			},
			want:    models.GetNetworkServicesResp{},
			wantErr: true,
		},
		{
			name: "Failed Test case 3: Error in callAPI",
			param: map[string]string{
				"name": templateName,
			},
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/networks/services"
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
			want:    models.GetNetworkServicesResp{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPIClient := NewMockAPIClientHandler(ctrl)
			r := RouterAPIService{
				Client: mockAPIClient,
				Cfg: Configuration{
					Host: mockHost,
				},
			}
			tt.given(mockAPIClient)
			got, err := r.GetNetworkServices(ctx, tt.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("RouterAPIService.GetNetworkServices() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RouterAPIService.GetNetworkServices() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouterAPIService_DeleteRouterNat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		routerID int
		natID    int
	}
	tests := []struct {
		name    string
		args    args
		given   func(m *MockAPIClientHandler)
		want    models.SuccessOrErrorMessage
		wantErr bool
	}{
		{
			name: "Normal test case 1",
			args: args{
				routerID: 1,
				natID:    2,
			},
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/networks/routers/1/nats/2"
				method := "DELETE"
				headers := getDefaultHeaders()
				req, _ := http.NewRequest(method, path, nil)
				respBody := io.NopCloser(bytes.NewReader([]byte(`
					{
						"success": true
					}
				`)))
				m.EXPECT().getVersion().Return(999999)
				m.EXPECT().prepareRequest(gomock.Any(), path, method, nil, headers,
					url.Values{}, url.Values{}, "", nil).Return(req, nil)

				m.EXPECT().callAPI(req).Return(&http.Response{
					StatusCode: 200,
					Body:       respBody,
				}, nil)
			},
			want: models.SuccessOrErrorMessage{
				Success: true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPIClient := NewMockAPIClientHandler(ctrl)
			r := RouterAPIService{
				Client: mockAPIClient,
				Cfg: Configuration{
					Host: mockHost,
				},
			}
			tt.given(mockAPIClient)
			got, err := r.DeleteRouterNat(context.Background(), tt.args.routerID, tt.args.natID)
			if (err != nil) != tt.wantErr {
				t.Errorf("RouterAPIService.DeleteRouterNat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RouterAPIService.DeleteRouterNat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouterAPIService_GetSpecificRouterNat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		routerID int
		natID    int
	}
	tests := []struct {
		name    string
		args    args
		given   func(m *MockAPIClientHandler)
		want    models.GetSpecificRouterNatResponse
		wantErr bool
	}{
		{
			name: "Normal test case 1",
			args: args{
				routerID: 1,
				natID:    2,
			},
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/networks/routers/1/nats/2"
				method := "GET"
				headers := getDefaultHeaders()
				req, _ := http.NewRequest(method, path, nil)
				respBody := io.NopCloser(bytes.NewReader([]byte(`
					{
						"networkRouterNAT": {
							"id": 1
						}
					}
				`)))
				m.EXPECT().getVersion().Return(999999)
				m.EXPECT().prepareRequest(gomock.Any(), path, method, nil, headers,
					url.Values{}, url.Values{}, "", nil).Return(req, nil)

				m.EXPECT().callAPI(req).Return(&http.Response{
					StatusCode: 200,
					Body:       respBody,
				}, nil)
			},
			want: models.GetSpecificRouterNatResponse{
				GetSpecificRouterNat: models.GetSpecificRouterNat{
					ID: 1,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPIClient := NewMockAPIClientHandler(ctrl)
			r := RouterAPIService{
				Client: mockAPIClient,
				Cfg: Configuration{
					Host: mockHost,
				},
			}
			tt.given(mockAPIClient)
			got, err := r.GetSpecificRouterNat(context.Background(), tt.args.routerID, tt.args.natID)
			if (err != nil) != tt.wantErr {
				t.Errorf("RouterAPIService.GetSpecificRouterNat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RouterAPIService.GetSpecificRouterNat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouterAPIService_CreateRouterNat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		ctx      context.Context
		routerID int
		request  models.CreateRouterNatRequest
	}
	tests := []struct {
		name    string
		args    args
		given   func(m *MockAPIClientHandler)
		want    models.SuccessOrErrorMessage
		wantErr bool
	}{
		{
			name: "Normal test case 1",
			args: args{
				routerID: 1,
				ctx:      context.Background(),
				request: models.CreateRouterNatRequest{
					CreateRouterNat: models.CreateRouterNat{
						Name: "test_nat",
					},
				},
			},
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/networks/routers/1/nats"
				method := "POST"
				headers := getDefaultHeaders()
				reqModel := models.CreateRouterNatRequest{
					CreateRouterNat: models.CreateRouterNat{
						Name: "test_nat",
					},
				}
				jsonByte, _ := json.Marshal(reqModel)
				postBody := io.NopCloser(bytes.NewReader(jsonByte))
				req, _ := http.NewRequest(method, path, postBody)
				respBody := io.NopCloser(bytes.NewReader([]byte(`
					{
						"success": true,
						"id": 2
					}
				`)))
				m.EXPECT().getVersion().Return(999999)
				m.EXPECT().prepareRequest(gomock.Any(), path, method, reqModel, headers,
					url.Values{}, url.Values{}, "", nil).Return(req, nil)

				m.EXPECT().callAPI(req).Return(&http.Response{
					StatusCode: 200,
					Body:       respBody,
				}, nil)
			},
			want: models.SuccessOrErrorMessage{
				IDModel: models.IDModel{ID: 2},
				Success: true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPIClient := NewMockAPIClientHandler(ctrl)
			r := RouterAPIService{
				Client: mockAPIClient,
				Cfg: Configuration{
					Host: mockHost,
				},
			}
			tt.given(mockAPIClient)
			got, err := r.CreateRouterNat(tt.args.ctx, tt.args.routerID, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("RouterAPIService.CreateRouterNat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RouterAPIService.CreateRouterNat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouterAPIService_UpdateRouterNat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		ctx      context.Context
		routerID int
		natID    int
		req      models.CreateRouterNatRequest
	}
	tests := []struct {
		name    string
		args    args
		given   func(m *MockAPIClientHandler)
		want    models.SuccessOrErrorMessage
		wantErr bool
	}{
		{
			name: "Normal test case 1",
			args: args{
				ctx:      context.Background(),
				routerID: 1,
				natID:    2,
				req: models.CreateRouterNatRequest{
					CreateRouterNat: models.CreateRouterNat{
						Name: "test-router-nat",
					},
				},
			},
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/networks/routers/1/nats/2"
				method := "PUT"
				headers := getDefaultHeaders()
				reqModel := models.CreateRouterNatRequest{
					CreateRouterNat: models.CreateRouterNat{
						Name: "test-router-nat",
					},
				}
				jsonByte, _ := json.Marshal(reqModel)
				postBody := io.NopCloser(bytes.NewReader(jsonByte))
				req, _ := http.NewRequest(method, path, postBody)
				respBody := io.NopCloser(bytes.NewReader([]byte(`
					{
						"success": true,
						"id": 2
					}
				`)))
				m.EXPECT().getVersion().Return(999999)
				m.EXPECT().prepareRequest(gomock.Any(), path, method, reqModel, headers,
					url.Values{}, url.Values{}, "", nil).Return(req, nil)

				m.EXPECT().callAPI(req).Return(&http.Response{
					StatusCode: 200,
					Body:       respBody,
				}, nil)
			},
			want: models.SuccessOrErrorMessage{
				Success: true,
				IDModel: models.IDModel{
					ID: 2,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPIClient := NewMockAPIClientHandler(ctrl)
			r := RouterAPIService{
				Client: mockAPIClient,
				Cfg: Configuration{
					Host: mockHost,
				},
			}
			tt.given(mockAPIClient)

			got, err := r.UpdateRouterNat(tt.args.ctx, tt.args.routerID, tt.args.natID, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("RouterAPIService.UpdateRouterNat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RouterAPIService.UpdateRouterNat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouterAPIService_DeleteRouterFirewallRuleGroup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		routerID        int
		firewallGroupID int
	}
	tests := []struct {
		name    string
		args    args
		given   func(m *MockAPIClientHandler)
		want    models.SuccessOrErrorMessage
		wantErr bool
	}{
		{
			name: "Normal test case 1",
			args: args{
				routerID:        1,
				firewallGroupID: 2,
			},
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/networks/routers/1/firewall-rule-groups/2"
				method := "DELETE"
				headers := getDefaultHeaders()
				req, _ := http.NewRequest(method, path, nil)
				respBody := io.NopCloser(bytes.NewReader([]byte(`
					{
						"success": true
					}
				`)))
				m.EXPECT().getVersion().Return(999999)
				m.EXPECT().prepareRequest(gomock.Any(), path, method, nil, headers,
					url.Values{}, url.Values{}, "", nil).Return(req, nil)

				m.EXPECT().callAPI(req).Return(&http.Response{
					StatusCode: 200,
					Body:       respBody,
				}, nil)
			},
			want: models.SuccessOrErrorMessage{
				Success: true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPIClient := NewMockAPIClientHandler(ctrl)
			r := RouterAPIService{
				Client: mockAPIClient,
				Cfg: Configuration{
					Host: mockHost,
				},
			}
			tt.given(mockAPIClient)
			got, err := r.DeleteRouterFirewallRuleGroup(context.Background(), tt.args.routerID, tt.args.firewallGroupID)
			if (err != nil) != tt.wantErr {
				t.Errorf("RouterAPIService.DeleteRouterFirewallRuleGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RouterAPIService.DeleteRouterFirewallRuleGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouterAPIService_GetSpecificRouterFirewallRuleGroup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		routerID        int
		firewallGroupID int
	}
	tests := []struct {
		name    string
		args    args
		given   func(m *MockAPIClientHandler)
		want    models.GetSpecificRouterFirewallRuleGroupResponse
		wantErr bool
	}{
		{
			name: "Normal test case 1",
			args: args{
				routerID:        1,
				firewallGroupID: 2,
			},
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/networks/routers/1/firewall-rule-groups/2"
				method := "GET"
				headers := getDefaultHeaders()
				req, _ := http.NewRequest(method, path, nil)
				respBody := io.NopCloser(bytes.NewReader([]byte(`
					{
						"ruleGroup": {
							"id": 1
						}
					}
				`)))
				m.EXPECT().getVersion().Return(999999)
				m.EXPECT().prepareRequest(gomock.Any(), path, method, nil, headers,
					url.Values{}, url.Values{}, "", nil).Return(req, nil)

				m.EXPECT().callAPI(req).Return(&http.Response{
					StatusCode: 200,
					Body:       respBody,
				}, nil)
			},
			want: models.GetSpecificRouterFirewallRuleGroupResponse{
				GetSpecificRouterFirewallRuleGroup: models.GetSpecificRouterFirewallRuleGroup{
					ID: 1,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPIClient := NewMockAPIClientHandler(ctrl)
			r := RouterAPIService{
				Client: mockAPIClient,
				Cfg: Configuration{
					Host: mockHost,
				},
			}
			tt.given(mockAPIClient)
			got, err := r.GetSpecificRouterFirewallRuleGroup(context.Background(), tt.args.routerID, tt.args.firewallGroupID)
			if (err != nil) != tt.wantErr {
				t.Errorf("RouterAPIService.GetSpecificRouterFirewallRuleGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RouterAPIService.GetSpecificRouterFirewallRuleGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouterAPIService_CreateRouterFirewallRuleGroup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		ctx      context.Context
		routerID int
		request  models.CreateRouterFirewallRuleGroupRequest
	}
	tests := []struct {
		name    string
		args    args
		given   func(m *MockAPIClientHandler)
		want    models.SuccessOrErrorMessage
		wantErr bool
	}{
		{
			name: "Normal test case 1",
			args: args{
				routerID: 1,
				ctx:      context.Background(),
				request: models.CreateRouterFirewallRuleGroupRequest{
					CreateRouterFirewallRuleGroup: models.CreateRouterFirewallRuleGroup{
						Name: "test_firewall_rule_group",
					},
				},
			},
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/networks/routers/1/firewall-rule-groups"
				method := "POST"
				headers := getDefaultHeaders()
				reqModel := models.CreateRouterFirewallRuleGroupRequest{
					CreateRouterFirewallRuleGroup: models.CreateRouterFirewallRuleGroup{
						Name: "test_firewall_rule_group",
					},
				}
				jsonByte, _ := json.Marshal(reqModel)
				postBody := io.NopCloser(bytes.NewReader(jsonByte))
				req, _ := http.NewRequest(method, path, postBody)
				respBody := io.NopCloser(bytes.NewReader([]byte(`
					{
						"success": true,
						"id": 2
					}
				`)))
				m.EXPECT().getVersion().Return(999999)
				m.EXPECT().prepareRequest(gomock.Any(), path, method, reqModel, headers,
					url.Values{}, url.Values{}, "", nil).Return(req, nil)

				m.EXPECT().callAPI(req).Return(&http.Response{
					StatusCode: 200,
					Body:       respBody,
				}, nil)
			},
			want: models.SuccessOrErrorMessage{
				IDModel: models.IDModel{ID: 2},
				Success: true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPIClient := NewMockAPIClientHandler(ctrl)
			r := RouterAPIService{
				Client: mockAPIClient,
				Cfg: Configuration{
					Host: mockHost,
				},
			}
			tt.given(mockAPIClient)
			got, err := r.CreateRouterFirewallRuleGroup(tt.args.ctx, tt.args.routerID, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("RouterAPIService.CreateRouterFirewallRuleGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RouterAPIService.CreateRouterFirewallRuleGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouterAPIService_DeleteRouterBgpNeighbor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		routerID      int
		bgpNeighborID int
	}
	tests := []struct {
		name    string
		args    args
		given   func(m *MockAPIClientHandler)
		want    models.SuccessOrErrorMessage
		wantErr bool
	}{
		{
			name: "Normal test case 1",
			args: args{
				routerID:      1,
				bgpNeighborID: 2,
			},
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/networks/routers/1/bgp-neighbors/2"
				method := "DELETE"
				headers := getDefaultHeaders()
				req, _ := http.NewRequest(method, path, nil)
				respBody := io.NopCloser(bytes.NewReader([]byte(`
					{
						"success": true
					}
				`)))
				m.EXPECT().getVersion().Return(999999)
				m.EXPECT().prepareRequest(gomock.Any(), path, method, nil, headers,
					url.Values{}, url.Values{}, "", nil).Return(req, nil)

				m.EXPECT().callAPI(req).Return(&http.Response{
					StatusCode: 200,
					Body:       respBody,
				}, nil)
			},
			want: models.SuccessOrErrorMessage{
				Success: true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPIClient := NewMockAPIClientHandler(ctrl)
			r := RouterAPIService{
				Client: mockAPIClient,
				Cfg: Configuration{
					Host: mockHost,
				},
			}
			tt.given(mockAPIClient)
			got, err := r.DeleteRouterBgpNeighbor(context.Background(), tt.args.routerID, tt.args.bgpNeighborID)
			if (err != nil) != tt.wantErr {
				t.Errorf("RouterAPIService.DeleteRouterBgpNeighbor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RouterAPIService.DeleteRouterBgpNeighbor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouterAPIService_GetSpecificRouterBgpNeighbor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		routerID      int
		bgpNeighborID int
	}
	tests := []struct {
		name    string
		args    args
		given   func(m *MockAPIClientHandler)
		want    models.GetSpecificNetworkRouterBgpNeighbor
		wantErr bool
	}{
		{
			name: "Normal test case 1",
			args: args{
				routerID:      1,
				bgpNeighborID: 2,
			},
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/networks/routers/1/bgp-neighbors/2"
				method := "GET"
				headers := getDefaultHeaders()
				req, _ := http.NewRequest(method, path, nil)
				respBody := io.NopCloser(bytes.NewReader([]byte(`
					{
						"networkRouterBgpNeighbor": {
							"id": 1
						}
					}
				`)))
				m.EXPECT().getVersion().Return(999999)
				m.EXPECT().prepareRequest(gomock.Any(), path, method, nil, headers,
					url.Values{}, url.Values{}, "", nil).Return(req, nil)

				m.EXPECT().callAPI(req).Return(&http.Response{
					StatusCode: 200,
					Body:       respBody,
				}, nil)
			},
			want: models.GetSpecificNetworkRouterBgpNeighbor{
				NetworkRouterBgpNeighbor: models.NetworkRouterBgpNeighborBody{
					ID: 1,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPIClient := NewMockAPIClientHandler(ctrl)
			r := RouterAPIService{
				Client: mockAPIClient,
				Cfg: Configuration{
					Host: mockHost,
				},
			}
			tt.given(mockAPIClient)
			got, err := r.GetSpecificRouterBgpNeighbor(context.Background(), tt.args.routerID, tt.args.bgpNeighborID)
			if (err != nil) != tt.wantErr {
				t.Errorf("RouterAPIService.GetSpecificRouterBgpNeighbor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RouterAPIService.GetSpecificRouterBgpNeighbor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouterAPIService_CreateRouterBgpNeighbor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		ctx      context.Context
		routerID int
		request  models.CreateNetworkRouterBgpNeighborRequest
	}
	tests := []struct {
		name    string
		args    args
		given   func(m *MockAPIClientHandler)
		want    models.SuccessOrErrorMessage
		wantErr bool
	}{
		{
			name: "Normal test case 1",
			args: args{
				routerID: 1,
				ctx:      context.Background(),
				request: models.CreateNetworkRouterBgpNeighborRequest{
					NetworkRouterBgpNeighbor: models.CreateRouterRequestBgpNeighborBody{
						IPAddress: "10.201.227.84",
					},
				},
			},
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/networks/routers/1/bgp-neighbors"
				method := "POST"
				headers := getDefaultHeaders()
				reqModel := models.CreateNetworkRouterBgpNeighborRequest{
					NetworkRouterBgpNeighbor: models.CreateRouterRequestBgpNeighborBody{
						IPAddress: "10.201.227.84",
					},
				}
				jsonByte, _ := json.Marshal(reqModel)
				postBody := io.NopCloser(bytes.NewReader(jsonByte))
				req, _ := http.NewRequest(method, path, postBody)
				respBody := io.NopCloser(bytes.NewReader([]byte(`
					{
						"success": true,
						"id": 2
					}
				`)))
				m.EXPECT().getVersion().Return(999999)
				m.EXPECT().prepareRequest(gomock.Any(), path, method, reqModel, headers,
					url.Values{}, url.Values{}, "", nil).Return(req, nil)

				m.EXPECT().callAPI(req).Return(&http.Response{
					StatusCode: 200,
					Body:       respBody,
				}, nil)
			},
			want: models.SuccessOrErrorMessage{
				IDModel: models.IDModel{ID: 2},
				Success: true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPIClient := NewMockAPIClientHandler(ctrl)
			r := RouterAPIService{
				Client: mockAPIClient,
				Cfg: Configuration{
					Host: mockHost,
				},
			}
			tt.given(mockAPIClient)
			got, err := r.CreateRouterBgpNeighbor(tt.args.ctx, tt.args.routerID, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("RouterAPIService.CreateRouterBgpNeighbor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RouterAPIService.CreateRouterBgpNeighbor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouterAPIService_UpdateRouterBgpNeighbor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		ctx           context.Context
		routerID      int
		bgpNeighborID int
		req           models.CreateNetworkRouterBgpNeighborRequest
	}
	tests := []struct {
		name    string
		args    args
		given   func(m *MockAPIClientHandler)
		want    models.SuccessOrErrorMessage
		wantErr bool
	}{
		{
			name: "Normal test case 1",
			args: args{
				ctx:           context.Background(),
				routerID:      1,
				bgpNeighborID: 2,
				req: models.CreateNetworkRouterBgpNeighborRequest{
					NetworkRouterBgpNeighbor: models.CreateRouterRequestBgpNeighborBody{
						IPAddress: "10.201.227.85",
					},
				},
			},
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/networks/routers/1/bgp-neighbors/2"
				method := "PUT"
				headers := getDefaultHeaders()
				reqModel := models.CreateNetworkRouterBgpNeighborRequest{
					NetworkRouterBgpNeighbor: models.CreateRouterRequestBgpNeighborBody{
						IPAddress: "10.201.227.85",
					},
				}
				jsonByte, _ := json.Marshal(reqModel)
				postBody := io.NopCloser(bytes.NewReader(jsonByte))
				req, _ := http.NewRequest(method, path, postBody)
				respBody := io.NopCloser(bytes.NewReader([]byte(`
					{
						"success": true
					}
				`)))
				m.EXPECT().getVersion().Return(999999)
				m.EXPECT().prepareRequest(gomock.Any(), path, method, reqModel, headers,
					url.Values{}, url.Values{}, "", nil).Return(req, nil)

				m.EXPECT().callAPI(req).Return(&http.Response{
					StatusCode: 200,
					Body:       respBody,
				}, nil)
			},
			want: models.SuccessOrErrorMessage{
				Success: true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPIClient := NewMockAPIClientHandler(ctrl)
			r := RouterAPIService{
				Client: mockAPIClient,
				Cfg: Configuration{
					Host: mockHost,
				},
			}
			tt.given(mockAPIClient)

			got, err := r.UpdateRouterBgpNeighbor(tt.args.ctx, tt.args.routerID, tt.args.bgpNeighborID, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("RouterAPIService.UpdateRouterBgpNeighbor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RouterAPIService.UpdateRouterBgpNeighbor() = %v, want %v", got, tt.want)
			}
		})
	}
}
