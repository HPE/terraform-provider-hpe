package client

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	consts "github.com/HPE/terraform-provider-hpe/greenlake/cloud/sdk/vmaascmp/common"
	"github.com/HPE/terraform-provider-hpe/greenlake/cloud/sdk/vmaascmp/models"
	gomock "github.com/golang/mock/gomock"
)

func Test_loadBalancerAPIService_CreateLoadBalancer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name    string
		args    models.CreateLoadBalancerRequest
		given   func(m *MockAPIClientHandler)
		want    models.CreateNetworkLoadBalancerResp
		wantErr bool
	}{
		{
			name: "Normal test case 1: Create LB",
			args: models.CreateLoadBalancerRequest{
				NetworkLoadBalancer: models.CreateNetworkLoadBalancerRequest{
					Name: "tf_LB",
				},
			},

			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/load-balancers"
				method := "POST"
				headers := getDefaultHeaders()
				req, _ := http.NewRequest(method, path, nil)

				m.EXPECT().getVersion().Return(999999)
				m.EXPECT().prepareRequest(gomock.Any(), path, method,
					models.CreateLoadBalancerRequest{
						NetworkLoadBalancer: models.CreateNetworkLoadBalancerRequest{
							Name: "tf_LB",
						},
					},
					headers, url.Values{}, url.Values{}, "", nil).Return(req, nil)

				m.EXPECT().callAPI(req).Return(&http.Response{
					StatusCode: 200,
					Body: io.NopCloser(bytes.NewReader([]byte(`
					{
						"success": true
					}
				`))),
				}, nil)
			},
			want: models.CreateNetworkLoadBalancerResp{
				Success:                 true,
				NetworkLoadBalancerResp: models.NetworkLoadBalancerResp{},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPIClient := NewMockAPIClientHandler(ctrl)
			tt.given(mockAPIClient)
			lb := LoadBalancerAPIService{
				Client: mockAPIClient,
				Cfg: Configuration{
					Host: mockHost,
				},
			}
			got, err := lb.CreateLoadBalancer(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadBalancerAPIService.CreateLoadBalancer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadBalancerAPIService.CreateLoadBalancer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_loadBalancerAPIService_UpdateLoadBalancer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name    string
		args    models.CreateLoadBalancerRequest
		given   func(m *MockAPIClientHandler)
		want    models.CreateNetworkLoadBalancerResp
		wantErr bool
	}{
		{
			name: "Normal test case 1: Create LB",
			args: models.CreateLoadBalancerRequest{
				NetworkLoadBalancer: models.CreateNetworkLoadBalancerRequest{
					Name: "tf_LB",
				},
			},

			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/load-balancers/1"
				method := "PUT"
				headers := getDefaultHeaders()
				req, _ := http.NewRequest(method, path, nil)

				m.EXPECT().getVersion().Return(999999)
				m.EXPECT().prepareRequest(gomock.Any(), path, method,
					models.CreateLoadBalancerRequest{
						NetworkLoadBalancer: models.CreateNetworkLoadBalancerRequest{
							Name: "tf_LB",
						},
					},
					headers, url.Values{}, url.Values{}, "", nil).Return(req, nil)

				m.EXPECT().callAPI(req).Return(&http.Response{
					StatusCode: 200,
					Body: io.NopCloser(bytes.NewReader([]byte(`
					{
						"success": true
					}
				`))),
				}, nil)
			},
			want: models.CreateNetworkLoadBalancerResp{
				Success:                 true,
				NetworkLoadBalancerResp: models.NetworkLoadBalancerResp{},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPIClient := NewMockAPIClientHandler(ctrl)
			tt.given(mockAPIClient)
			lb := LoadBalancerAPIService{
				Client: mockAPIClient,
				Cfg: Configuration{
					Host: mockHost,
				},
			}
			got, err := lb.UpdateLoadBalancer(context.Background(), 1, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadBalancerAPIService.UpdateLoadBalancer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadBalancerAPIService.UpdateLoadBalancer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_loadBalancerAPIService_DeleteLoadBalancer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name    string
		lbID    int
		given   func(m *MockAPIClientHandler)
		want    models.SuccessOrErrorMessage
		wantErr bool
	}{
		{
			name: "Normal Test case 1: Delete a LB",
			lbID: 1,
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/load-balancers/1"
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
					getURLValues(nil), url.Values{}, "", nil).Return(req, nil)

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
			tt.given(mockAPIClient)
			lb := LoadBalancerAPIService{
				Client: mockAPIClient,
				Cfg: Configuration{
					Host: mockHost,
				},
			}

			got, err := lb.DeleteLoadBalancer(context.Background(), tt.lbID)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadBalancerAPIService.DeleteLoadBalancer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadBalancerAPIService.DeleteLoadBalancer() = %v, want %v", got, tt.want)
			}

		})
	}
}

func Test_loadBalancerAPIService_GetSpecificLoadBalancers(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	templateName := "test_template_get_a_specific_LB"
	tests := []struct {
		name    string
		lbID    int
		given   func(m *MockAPIClientHandler)
		want    models.GetSpecificNetworkLoadBalancer
		wantErr bool
	}{
		{
			name: "Normal Test case 1: Get a specific LB",
			lbID: 1,
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/load-balancers/1"
				method := "GET"
				headers := getDefaultHeaders()
				req, _ := http.NewRequest(method, path, nil)
				respBody := io.NopCloser(bytes.NewReader([]byte(`
					{
						"loadBalancer":{
							"id": 1,
							"name":"test_template_get_a_specific_LB"
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
			want: models.GetSpecificNetworkLoadBalancer{
				GetSpecificNetworkLoadBalancerResp: models.GetSpecificNetworkLoadBalancerResp{
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
			lb := LoadBalancerAPIService{
				Client: mockAPIClient,
				Cfg: Configuration{
					Host: mockHost,
				},
			}
			tt.given(mockAPIClient)
			got, err := lb.GetSpecificLoadBalancers(ctx, tt.lbID)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadBalancerAPIService.GetSpecificLoadBalancers() = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadBalancerAPIService.GetSpecificLoadBalancers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_loadBalancerAPIService_CreateLBMonitor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name    string
		args    models.CreateLBMonitor
		given   func(m *MockAPIClientHandler)
		want    models.CreateLBMonitorResp
		wantErr bool
	}{
		{
			name: "Normal test case 1: Create LB-Monitor",
			args: models.CreateLBMonitor{
				CreateLBMonitorReq: models.CreateLBMonitorReq{
					Name: "tf_LBMonitor",
				},
			},

			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/load-balancers/1/monitors"
				method := "POST"
				headers := getDefaultHeaders()
				req, _ := http.NewRequest(method, path, nil)

				m.EXPECT().getVersion().Return(999999)
				m.EXPECT().prepareRequest(gomock.Any(), path, method,
					models.CreateLBMonitor{
						CreateLBMonitorReq: models.CreateLBMonitorReq{
							Name: "tf_LBMonitor",
						},
					},
					headers, url.Values{}, url.Values{}, "", nil).Return(req, nil)

				m.EXPECT().callAPI(req).Return(&http.Response{
					StatusCode: 200,
					Body: io.NopCloser(bytes.NewReader([]byte(`
					{
						"success": true
					}
				`))),
				}, nil)
			},
			want: models.CreateLBMonitorResp{
				Success:       true,
				LBMonitorResp: models.LBMonitorResp{},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPIClient := NewMockAPIClientHandler(ctrl)
			tt.given(mockAPIClient)
			lb := LoadBalancerAPIService{
				Client: mockAPIClient,
				Cfg: Configuration{
					Host: mockHost,
				},
			}
			got, err := lb.CreateLBMonitor(context.Background(), tt.args, 1)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadBalancerAPIService.CreateLBMonitor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadBalancerAPIService.CreateLBMonitor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_loadBalancerAPIService_UpdateLBMonitor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name    string
		args    models.CreateLBMonitor
		given   func(m *MockAPIClientHandler)
		want    models.CreateLBMonitorResp
		wantErr bool
	}{
		{
			name: "Normal test case 1: Create LB-Monitor",
			args: models.CreateLBMonitor{
				CreateLBMonitorReq: models.CreateLBMonitorReq{
					Name: "tf_LBMonitor",
				},
			},

			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/load-balancers/1/monitors/1"
				method := "PUT"
				headers := getDefaultHeaders()
				req, _ := http.NewRequest(method, path, nil)

				m.EXPECT().getVersion().Return(999999)
				m.EXPECT().prepareRequest(gomock.Any(), path, method,
					models.CreateLBMonitor{
						CreateLBMonitorReq: models.CreateLBMonitorReq{
							Name: "tf_LBMonitor",
						},
					},
					headers, url.Values{}, url.Values{}, "", nil).Return(req, nil)

				m.EXPECT().callAPI(req).Return(&http.Response{
					StatusCode: 200,
					Body: io.NopCloser(bytes.NewReader([]byte(`
					{
						"success": true
					}
				`))),
				}, nil)
			},
			want: models.CreateLBMonitorResp{
				Success:       true,
				LBMonitorResp: models.LBMonitorResp{},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPIClient := NewMockAPIClientHandler(ctrl)
			tt.given(mockAPIClient)
			lb := LoadBalancerAPIService{
				Client: mockAPIClient,
				Cfg: Configuration{
					Host: mockHost,
				},
			}
			got, err := lb.UpdateLBMonitor(context.Background(), tt.args, 1, 1)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadBalancerAPIService.UpdateLBMonitor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadBalancerAPIService.UpdateLBMonitor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_loadBalancerAPIService_DeleteLBMonitor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name    string
		lbID    int
		given   func(m *MockAPIClientHandler)
		want    models.SuccessOrErrorMessage
		wantErr bool
	}{
		{
			name: "Normal Test case 1: Delete a LB-M",
			lbID: 1,
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/load-balancers/1/monitors/1"
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
					getURLValues(nil), url.Values{}, "", nil).Return(req, nil)

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
			tt.given(mockAPIClient)
			lb := LoadBalancerAPIService{
				Client: mockAPIClient,
				Cfg: Configuration{
					Host: mockHost,
				},
			}

			got, err := lb.DeleteLBMonitor(context.Background(), tt.lbID, 1)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadBalancerAPIService.DeleteLBMonitor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadBalancerAPIService.DeleteLBMonitor() = %v, want %v", got, tt.want)
			}

		})
	}
}

func Test_loadBalancerAPIService_GetSpecificLBMonitor(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name    string
		lbID    int
		given   func(m *MockAPIClientHandler)
		want    models.GetSpecificLBMonitor
		wantErr bool
	}{
		{
			name: "Normal Test case 1: Get a specific LB-Monitor",
			lbID: 1,
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/load-balancers/1/monitors/1"
				method := "GET"
				headers := getDefaultHeaders()
				req, _ := http.NewRequest(method, path, nil)
				respBody := io.NopCloser(bytes.NewReader([]byte(`
					{
						"loadBalancerMonitor":{
							"name":"test_template_get_a_specific_LB"
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
			want: models.GetSpecificLBMonitor{
				GetSpecificLBMonitorResp: models.GetSpecificLBMonitorResp{
					Name: "test_template_get_a_specific_LB",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPIClient := NewMockAPIClientHandler(ctrl)
			lb := LoadBalancerAPIService{
				Client: mockAPIClient,
				Cfg: Configuration{
					Host: mockHost,
				},
			}
			tt.given(mockAPIClient)
			got, err := lb.GetSpecificLBMonitor(ctx, tt.lbID, 1)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadBalancerAPIService.GetSpecificLBMonitor() = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadBalancerAPIService.GetSpecificLBMonitor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_loadBalancerAPIService_CreateLBProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name    string
		args    models.CreateLBProfile
		given   func(m *MockAPIClientHandler)
		want    models.CreateLBProfileResp
		wantErr bool
	}{
		{
			name: "Normal test case 1: Create LB Profile",
			args: models.CreateLBProfile{
				CreateLBProfileReq: models.CreateLBProfileReq{
					Name: "tf_LB-Profile",
				},
			},

			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/load-balancers/1/profiles"
				method := "POST"
				headers := getDefaultHeaders()
				req, _ := http.NewRequest(method, path, nil)

				m.EXPECT().getVersion().Return(999999)
				m.EXPECT().prepareRequest(gomock.Any(), path, method,
					models.CreateLBProfile{
						CreateLBProfileReq: models.CreateLBProfileReq{
							Name: "tf_LB-Profile",
						},
					},
					headers, url.Values{}, url.Values{}, "", nil).Return(req, nil)

				m.EXPECT().callAPI(req).Return(&http.Response{
					StatusCode: 200,
					Body: io.NopCloser(bytes.NewReader([]byte(`
					{
						"success": true
					}
				`))),
				}, nil)
			},
			want: models.CreateLBProfileResp{
				Success:       true,
				LBProfileResp: models.LBProfileResp{},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPIClient := NewMockAPIClientHandler(ctrl)
			tt.given(mockAPIClient)
			lb := LoadBalancerAPIService{
				Client: mockAPIClient,
				Cfg: Configuration{
					Host: mockHost,
				},
			}
			got, err := lb.CreateLBProfile(context.Background(), tt.args, 1)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadBalancerAPIService.CreateLBProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadBalancerAPIService.CreateLBProfile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_loadBalancerAPIService_UpdateLBProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name    string
		args    models.CreateLBProfile
		given   func(m *MockAPIClientHandler)
		want    models.CreateLBProfileResp
		wantErr bool
	}{
		{
			name: "Normal test case 1: Create LB Profile",
			args: models.CreateLBProfile{
				CreateLBProfileReq: models.CreateLBProfileReq{
					Name: "tf_LB-Profile",
				},
			},

			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/load-balancers/1/profiles/1"
				method := "PUT"
				headers := getDefaultHeaders()
				req, _ := http.NewRequest(method, path, nil)

				m.EXPECT().getVersion().Return(999999)
				m.EXPECT().prepareRequest(gomock.Any(), path, method,
					models.CreateLBProfile{
						CreateLBProfileReq: models.CreateLBProfileReq{
							Name: "tf_LB-Profile",
						},
					},
					headers, url.Values{}, url.Values{}, "", nil).Return(req, nil)

				m.EXPECT().callAPI(req).Return(&http.Response{
					StatusCode: 200,
					Body: io.NopCloser(bytes.NewReader([]byte(`
					{
						"success": true
					}
				`))),
				}, nil)
			},
			want: models.CreateLBProfileResp{
				Success:       true,
				LBProfileResp: models.LBProfileResp{},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPIClient := NewMockAPIClientHandler(ctrl)
			tt.given(mockAPIClient)
			lb := LoadBalancerAPIService{
				Client: mockAPIClient,
				Cfg: Configuration{
					Host: mockHost,
				},
			}
			got, err := lb.UpdateLBProfile(context.Background(), tt.args, 1, 1)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadBalancerAPIService.UpdateLBProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadBalancerAPIService.UpdateLBProfile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_loadBalancerAPIService_DeleteLBProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name    string
		lbID    int
		given   func(m *MockAPIClientHandler)
		want    models.SuccessOrErrorMessage
		wantErr bool
	}{
		{
			name: "Normal Test case 1: Delete a LB-P",
			lbID: 1,
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/load-balancers/1/profiles/1"
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
					getURLValues(nil), url.Values{}, "", nil).Return(req, nil)

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
			tt.given(mockAPIClient)
			lb := LoadBalancerAPIService{
				Client: mockAPIClient,
				Cfg: Configuration{
					Host: mockHost,
				},
			}

			got, err := lb.DeleteLBProfile(context.Background(), tt.lbID, 1)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadBalancerAPIService.DeleteLBProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadBalancerAPIService.DeleteLBProfile() = %v, want %v", got, tt.want)
			}

		})
	}
}

func Test_loadBalancerAPIService_GetSpecificLBProfile(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	templateName := "test_template_get_a_specific_LB-P"
	tests := []struct {
		name    string
		lbID    int
		given   func(m *MockAPIClientHandler)
		want    models.GetLBSpecificProfile
		wantErr bool
	}{
		{
			name: "Normal Test case 1: Get a specific LB-P",
			lbID: 1,
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/load-balancers/1/profiles/1"
				method := "GET"
				headers := getDefaultHeaders()
				req, _ := http.NewRequest(method, path, nil)
				respBody := io.NopCloser(bytes.NewReader([]byte(`
					{
						"loadBalancerProfile":{
							"id": 1,
							"name":"test_template_get_a_specific_LB-P"
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
			want: models.GetLBSpecificProfile{
				GetLBSpecificProfilesResp: models.GetLBSpecificProfilesResp{
					ID:   0,
					Name: templateName,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPIClient := NewMockAPIClientHandler(ctrl)
			lb := LoadBalancerAPIService{
				Client: mockAPIClient,
				Cfg: Configuration{
					Host: mockHost,
				},
			}
			tt.given(mockAPIClient)
			got, err := lb.GetSpecificLBProfile(ctx, tt.lbID, 1)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadBalancerAPIService.GetSpecificLBProfile() = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadBalancerAPIService.GetSpecificLBProfile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_loadBalancerAPIService_CreateLBPool(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name    string
		args    models.CreateLBPool
		given   func(m *MockAPIClientHandler)
		want    models.CreateLBPoolResp
		wantErr bool
	}{
		{
			name: "Normal test case 1: Create LB-Pool",
			args: models.CreateLBPool{
				CreateLBPoolReq: models.CreateLBPoolReq{
					Name: "tf_LB-Pool",
				},
			},

			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/load-balancers/1/pools"
				method := "POST"
				headers := getDefaultHeaders()
				req, _ := http.NewRequest(method, path, nil)

				m.EXPECT().getVersion().Return(999999)
				m.EXPECT().prepareRequest(gomock.Any(), path, method,
					models.CreateLBPool{
						CreateLBPoolReq: models.CreateLBPoolReq{
							Name: "tf_LB-Pool",
						},
					},
					headers, url.Values{}, url.Values{}, "", nil).Return(req, nil)

				m.EXPECT().callAPI(req).Return(&http.Response{
					StatusCode: 200,
					Body: io.NopCloser(bytes.NewReader([]byte(`
					{
						"success": true
					}
				`))),
				}, nil)
			},
			want: models.CreateLBPoolResp{
				Success:    true,
				LBPoolResp: models.LBPoolResp{},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPIClient := NewMockAPIClientHandler(ctrl)
			tt.given(mockAPIClient)
			lb := LoadBalancerAPIService{
				Client: mockAPIClient,
				Cfg: Configuration{
					Host: mockHost,
				},
			}
			got, err := lb.CreateLBPool(context.Background(), tt.args, 1)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadBalancerAPIService.CreateLBPool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadBalancerAPIService.CreateLBPool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_loadBalancerAPIService_UpdateLBPool(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name    string
		args    models.CreateLBPool
		given   func(m *MockAPIClientHandler)
		want    models.CreateLBPoolResp
		wantErr bool
	}{
		{
			name: "Normal test case 1: Create LB-Pool",
			args: models.CreateLBPool{
				CreateLBPoolReq: models.CreateLBPoolReq{
					Name: "tf_LB-Pool",
				},
			},

			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/load-balancers/1/pools/1"
				method := "PUT"
				headers := getDefaultHeaders()
				req, _ := http.NewRequest(method, path, nil)

				m.EXPECT().getVersion().Return(999999)
				m.EXPECT().prepareRequest(gomock.Any(), path, method,
					models.CreateLBPool{
						CreateLBPoolReq: models.CreateLBPoolReq{
							Name: "tf_LB-Pool",
						},
					},
					headers, url.Values{}, url.Values{}, "", nil).Return(req, nil)

				m.EXPECT().callAPI(req).Return(&http.Response{
					StatusCode: 200,
					Body: io.NopCloser(bytes.NewReader([]byte(`
					{
						"success": true
					}
				`))),
				}, nil)
			},
			want: models.CreateLBPoolResp{
				Success:    true,
				LBPoolResp: models.LBPoolResp{},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPIClient := NewMockAPIClientHandler(ctrl)
			tt.given(mockAPIClient)
			lb := LoadBalancerAPIService{
				Client: mockAPIClient,
				Cfg: Configuration{
					Host: mockHost,
				},
			}
			got, err := lb.UpdateLBPool(context.Background(), tt.args, 1, 1)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadBalancerAPIService.UpdateLBPool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadBalancerAPIService.UpdateLBPool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_loadBalancerAPIService_DeleteLBPool(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name    string
		lbID    int
		given   func(m *MockAPIClientHandler)
		want    models.SuccessOrErrorMessage
		wantErr bool
	}{
		{
			name: "Normal Test case 1: Delete a LB-P",
			lbID: 1,
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/load-balancers/1/pools/1"
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
					getURLValues(nil), url.Values{}, "", nil).Return(req, nil)

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
			tt.given(mockAPIClient)
			lb := LoadBalancerAPIService{
				Client: mockAPIClient,
				Cfg: Configuration{
					Host: mockHost,
				},
			}

			got, err := lb.DeleteLBPool(context.Background(), tt.lbID, 1)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadBalancerAPIService.DeleteLBPool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadBalancerAPIService.DeleteLBPool() = %v, want %v", got, tt.want)
			}

		})
	}
}

func Test_loadBalancerAPIService_GetSpecificLBPool(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	templateName := "test_template_get_a_specific_LB-P"
	tests := []struct {
		name    string
		lbID    int
		given   func(m *MockAPIClientHandler)
		want    models.GetSpecificLBPool
		wantErr bool
	}{
		{
			name: "Normal Test case 1: Get a specific LB-P",
			lbID: 1,
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/load-balancers/1/pools/1"
				method := "GET"
				headers := getDefaultHeaders()
				req, _ := http.NewRequest(method, path, nil)
				respBody := io.NopCloser(bytes.NewReader([]byte(`
					{
						"loadBalancerPool":{
							"id": 1,
							"name":"test_template_get_a_specific_LB-P"
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
			want: models.GetSpecificLBPool{
				GetSpecificLBPoolResp: models.GetSpecificLBPoolResp{
					ID:   0,
					Name: templateName,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPIClient := NewMockAPIClientHandler(ctrl)
			lb := LoadBalancerAPIService{
				Client: mockAPIClient,
				Cfg: Configuration{
					Host: mockHost,
				},
			}
			tt.given(mockAPIClient)
			got, err := lb.GetSpecificLBPool(ctx, tt.lbID, 1)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadBalancerAPIService.GetSpecificLBPool() = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadBalancerAPIService.GetSpecificLBPool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_loadBalancerAPIService_CreateLBVirtualServers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name    string
		args    models.CreateLBVirtualServers
		given   func(m *MockAPIClientHandler)
		want    models.LBVirtualServersResp
		wantErr bool
	}{
		{
			name: "Normal test case 1: Create LB-VS",
			args: models.CreateLBVirtualServers{
				CreateLBVirtualServersReq: models.CreateLBVirtualServersReq{
					VipName: "tf_LB",
				},
			},

			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/load-balancers/1/virtual-servers"
				method := "POST"
				headers := getDefaultHeaders()
				req, _ := http.NewRequest(method, path, nil)

				m.EXPECT().getVersion().Return(999999)
				m.EXPECT().prepareRequest(gomock.Any(), path, method,
					models.CreateLBVirtualServers{
						CreateLBVirtualServersReq: models.CreateLBVirtualServersReq{
							VipName: "tf_LB",
						},
					},
					headers, url.Values{}, url.Values{}, "", nil).Return(req, nil)

				m.EXPECT().callAPI(req).Return(&http.Response{
					StatusCode: 200,
					Body: io.NopCloser(bytes.NewReader([]byte(`
					{
						"success": true
					}
				`))),
				}, nil)
			},
			want: models.LBVirtualServersResp{
				Success:                    true,
				CreateLBVirtualServersResp: models.CreateLBVirtualServersResp{},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPIClient := NewMockAPIClientHandler(ctrl)
			tt.given(mockAPIClient)
			lb := LoadBalancerAPIService{
				Client: mockAPIClient,
				Cfg: Configuration{
					Host: mockHost,
				},
			}
			got, err := lb.CreateLBVirtualServers(context.Background(), tt.args, 1)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadBalancerAPIService.CreateLBVirtualServers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadBalancerAPIService.CreateLBVirtualServers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_loadBalancerAPIService_UpdateLBVirtualServers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name    string
		args    models.CreateLBVirtualServers
		given   func(m *MockAPIClientHandler)
		want    models.LBVirtualServersResp
		wantErr bool
	}{
		{
			name: "Normal test case 1: Create LB-VS",
			args: models.CreateLBVirtualServers{
				CreateLBVirtualServersReq: models.CreateLBVirtualServersReq{
					VipName: "tf_LB",
				},
			},

			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/load-balancers/1/virtual-servers/1"
				method := "PUT"
				headers := getDefaultHeaders()
				req, _ := http.NewRequest(method, path, nil)

				m.EXPECT().getVersion().Return(999999)
				m.EXPECT().prepareRequest(gomock.Any(), path, method,
					models.CreateLBVirtualServers{
						CreateLBVirtualServersReq: models.CreateLBVirtualServersReq{
							VipName: "tf_LB",
						},
					},
					headers, url.Values{}, url.Values{}, "", nil).Return(req, nil)

				m.EXPECT().callAPI(req).Return(&http.Response{
					StatusCode: 200,
					Body: io.NopCloser(bytes.NewReader([]byte(`
					{
						"success": true
					}
				`))),
				}, nil)
			},
			want: models.LBVirtualServersResp{
				Success:                    true,
				CreateLBVirtualServersResp: models.CreateLBVirtualServersResp{},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPIClient := NewMockAPIClientHandler(ctrl)
			tt.given(mockAPIClient)
			lb := LoadBalancerAPIService{
				Client: mockAPIClient,
				Cfg: Configuration{
					Host: mockHost,
				},
			}
			got, err := lb.UpdateLBVirtualServers(context.Background(), tt.args, 1, 1)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadBalancerAPIService.UpdateLBVirtualServers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadBalancerAPIService.UpdateLBVirtualServers() = %v, want %v", got, tt.want)
			}
		})
	}
}
func Test_loadBalancerAPIService_DeleteLBVirtualServers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name    string
		lbID    int
		given   func(m *MockAPIClientHandler)
		want    models.SuccessOrErrorMessage
		wantErr bool
	}{
		{
			name: "Normal Test case 1: Delete a LB-VS",
			lbID: 1,
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/load-balancers/1/virtual-servers/1"
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
					getURLValues(nil), url.Values{}, "", nil).Return(req, nil)

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
			tt.given(mockAPIClient)
			lb := LoadBalancerAPIService{
				Client: mockAPIClient,
				Cfg: Configuration{
					Host: mockHost,
				},
			}

			got, err := lb.DeleteLBVirtualServers(context.Background(), tt.lbID, 1)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadBalancerAPIService.DeleteLBVirtualServers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadBalancerAPIService.DeleteLBVirtualServers() = %v, want %v", got, tt.want)
			}

		})
	}
}

func Test_loadBalancerAPIService_GetSpecificLBVirtualServers(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	templateName := "test_template_get_a_specific_LB-VS"
	tests := []struct {
		name    string
		lbID    int
		given   func(m *MockAPIClientHandler)
		want    models.GetSpecificLBVirtualServers
		wantErr bool
	}{
		{
			name: "Normal Test case 1: Get a specific LB-VS",
			lbID: 1,
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/load-balancers/1/virtual-servers/1"
				method := "GET"
				headers := getDefaultHeaders()
				req, _ := http.NewRequest(method, path, nil)
				respBody := io.NopCloser(bytes.NewReader([]byte(`
					{
						"loadBalancerInstance":{
							"id": 1,
							"name":"test_template_get_a_specific_LB-VS"
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
			want: models.GetSpecificLBVirtualServers{
				GetSpecificLBVirtualServersResp: models.GetSpecificLBVirtualServersResp{
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
			lb := LoadBalancerAPIService{
				Client: mockAPIClient,
				Cfg: Configuration{
					Host: mockHost,
				},
			}
			tt.given(mockAPIClient)
			got, err := lb.GetSpecificLBVirtualServer(ctx, tt.lbID, 1)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadBalancerAPIService.GetSpecificLBVirtualServer() = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadBalancerAPIService.GetSpecificLBVirtualServer() = %v, want %v", got, tt.want)
			}
		})
	}
}
