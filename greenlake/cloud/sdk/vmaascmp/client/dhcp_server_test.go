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

func Test_dhcpServerAPIService_CreateDhcpServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name    string
		args    models.CreateNetworkDhcpServerRequest
		given   func(m *MockAPIClientHandler)
		want    models.CreateNetworkDhcpServerResp
		wantErr bool
	}{
		{
			name: "Normal test case 1: Create DHCP",
			args: models.CreateNetworkDhcpServerRequest{
				NetworkDhcpServer: models.CreateNetworkDhcpServer{
					Name: "tf_DHCP_SERVER",
				},
			},

			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/networks/servers/1/dhcp-servers"
				method := "POST"
				headers := getDefaultHeaders()
				req, _ := http.NewRequest(method, path, nil)

				m.EXPECT().getVersion().Return(999999)
				m.EXPECT().prepareRequest(gomock.Any(), path, method,
					models.CreateNetworkDhcpServerRequest{
						NetworkDhcpServer: models.CreateNetworkDhcpServer{
							Name: "tf_DHCP_SERVER",
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
			want: models.CreateNetworkDhcpServerResp{
				Success: true,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPIClient := NewMockAPIClientHandler(ctrl)
			tt.given(mockAPIClient)
			lb := DhcpServerAPIService{
				Client: mockAPIClient,
				Cfg: Configuration{
					Host: mockHost,
				},
			}
			got, err := lb.CreateDhcpServer(context.Background(), 1, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("DhcpServerAPIService.CreateDhcpServer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DhcpServerAPIService.CreateDhcpServer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dhcpServerAPIService_UpdateDhcpServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name    string
		args    models.CreateNetworkDhcpServerRequest
		given   func(m *MockAPIClientHandler)
		want    models.CreateNetworkDhcpServerResp
		wantErr bool
	}{
		{
			name: "Normal test case 1: Update DHCP Server",
			args: models.CreateNetworkDhcpServerRequest{
				NetworkDhcpServer: models.CreateNetworkDhcpServer{
					Name: "tf_DHCP_SERVER",
				},
			},

			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/networks/servers/1/dhcp-servers/1"
				method := "PUT"
				headers := getDefaultHeaders()
				req, _ := http.NewRequest(method, path, nil)

				m.EXPECT().getVersion().Return(999999)
				m.EXPECT().prepareRequest(gomock.Any(), path, method,
					models.CreateNetworkDhcpServerRequest{
						NetworkDhcpServer: models.CreateNetworkDhcpServer{
							Name: "tf_DHCP_SERVER",
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
			want: models.CreateNetworkDhcpServerResp{
				Success: true,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPIClient := NewMockAPIClientHandler(ctrl)
			tt.given(mockAPIClient)
			lb := DhcpServerAPIService{
				Client: mockAPIClient,
				Cfg: Configuration{
					Host: mockHost,
				},
			}
			got, err := lb.UpdateDhcpServer(context.Background(), 1, 1, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("DhcpServerAPIService.UpdateDhcpServer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DhcpServerAPIService.UpdateDhcpServer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dhcpServerAPIService_DeleteDhcpServer(t *testing.T) {
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
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/networks/servers/1/dhcp-servers/1"
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
			lb := DhcpServerAPIService{
				Client: mockAPIClient,
				Cfg: Configuration{
					Host: mockHost,
				},
			}

			got, err := lb.DeleteDhcpServer(context.Background(), 1, tt.lbID)
			if (err != nil) != tt.wantErr {
				t.Errorf("DhcpServerAPIService.DeleteDhcpServer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DhcpServerAPIService.DeleteDhcpServer() = %v, want %v", got, tt.want)
			}

		})
	}
}

func Test_dhcpServerAPIService_GetSpecificDhcpServer(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	templateName := "test_template_get_a_specific_DHCP"
	tests := []struct {
		name    string
		lbID    int
		given   func(m *MockAPIClientHandler)
		want    models.GetSpecificNetworkDhcpServer
		wantErr bool
	}{
		{
			name: "Normal Test case 1: Get a specific DHCP",
			lbID: 1,
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				path := mockHost + "/" + consts.VmaasCmpAPIBasePath + "/networks/servers/1/dhcp-servers/1"
				method := "GET"
				headers := getDefaultHeaders()
				req, _ := http.NewRequest(method, path, nil)
				respBody := io.NopCloser(bytes.NewReader([]byte(`
					{
						"networkDhcpServer":{
							"id": 1,
							"name":"test_template_get_a_specific_DHCP"
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
			want: models.GetSpecificNetworkDhcpServer{
				GetSpecificNetworkDhcpServerResp: models.GetSpecificNetworkDhcpServerResp{
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
			lb := DhcpServerAPIService{
				Client: mockAPIClient,
				Cfg: Configuration{
					Host: mockHost,
				},
			}
			tt.given(mockAPIClient)
			got, err := lb.GetSpecificDhcpServer(ctx, 1, tt.lbID)
			if (err != nil) != tt.wantErr {
				t.Errorf("DhcpServerAPIService.GetSpecificDhcpServer() = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DhcpServerAPIService.GetSpecificDhcpServer() = %v, want %v", got, tt.want)
			}
		})
	}
}
