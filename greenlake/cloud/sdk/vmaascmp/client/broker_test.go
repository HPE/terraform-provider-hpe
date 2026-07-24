// (C) Copyright 2024-2025 Hewlett Packard Enterprise Development LP

package client

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/golang/mock/gomock"

	consts "github.com/HPE/terraform-provider-hpe/greenlake/cloud/sdk/vmaascmp/common"
	"github.com/HPE/terraform-provider-hpe/greenlake/cloud/sdk/vmaascmp/models"
)

const (
	testServiceInstanceID    = "18ba6409-ac59-4eac-9414-0147e72d615e"
	testAccessToken          = "2b9fba7f-7c14-4773-a970-a9ad393811ac"
	testRefreshToken         = "7806acfb-f847-48b1-a6d5-6119dccb3ffe"
	testMorpheusURL          = "https://1234-mp.private.greenlake.hpe-gl-intg.com"
	testAccessTokenExpires   = 1758034360176
	testAccessTokenExpiresIn = 3600
)

func TestBrokerAPIService_GetMorpheusDetails(t *testing.T) {
	ctx := context.Background()
	testCtrl := gomock.NewController(t)
	defer testCtrl.Finish()

	headers := getDefaultHeaders()

	queryParams := map[string]string{
		"location":   "BLR",
		"space_name": "default",
	}

	clientCfg := Configuration{
		DefaultHeader:      headers,
		DefaultQueryParams: queryParams,
	}

	tests := []struct {
		name    string
		given   func(m *MockAPIClientHandler)
		want    models.TFMorpheusDetails
		wantErr bool
	}{
		{
			name: "Test GetMorpheusDetails success",
			want: models.TFMorpheusDetails{
				ID:          testServiceInstanceID,
				AccessToken: testAccessToken,
				ValidTill:   testAccessTokenExpiresIn,
				URL:         testMorpheusURL,
			},
			wantErr: false,
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				pathSubscription := mockHost + "/" + consts.CMPDetails
				method := "GET"
				reqSubscription, _ := http.NewRequest(method, pathSubscription, nil)
				respBodySubscription := io.NopCloser(bytes.NewReader([]byte(`
					{
						"ServiceInstanceID": "` + testServiceInstanceID + `",
						"TenantID": "1234",
						"TenantName": "tenant",
						"LocationName": "BLR",	
						"URL": "` + testMorpheusURL + `",
						"TokenDetails": {
							"access_token": "` + testAccessToken + `",
							"expires": ` + fmt.Sprintf("%d", testAccessTokenExpires) + `,
							"refresh_token": "` + testRefreshToken + `",
							"expires_in": ` + fmt.Sprintf("%d", testAccessTokenExpiresIn) + `	
						}
					}
				`)))
				// mock the context only since it is not validated in this function
				m.EXPECT().getVersion().Return(999999)
				m.EXPECT().prepareRequest(gomock.Any(), pathSubscription, method, nil, headers,
					getURLValues(queryParams), url.Values{}, "", nil).Return(reqSubscription, nil)

				m.EXPECT().callAPI(reqSubscription).Return(&http.Response{
					StatusCode: 200,
					Body:       respBodySubscription,
				}, nil)
			},
		},

		{
			name:    "Test GetMorpheusDetails error in get subscription details prepare request",
			want:    models.TFMorpheusDetails{},
			wantErr: true,
			given: func(m *MockAPIClientHandler) {
				m.EXPECT().getHost().Return(mockHost)
				pathSubscription := mockHost + "/" + consts.CMPDetails
				method := "GET"
				// mock the context only since it is not validated in this function
				m.EXPECT().getVersion().Return(999999)
				m.EXPECT().prepareRequest(gomock.Any(), pathSubscription, method, nil, headers,
					getURLValues(queryParams), url.Values{}, "", nil).
					Return(nil, errors.New("error in prepare request"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := NewMockAPIClientHandler(testCtrl)
			tt.given(mockClient)
			a := &BrokerAPIService{
				Cfg:    clientCfg,
				Client: mockClient,
			}
			got, err := a.GetMorpheusDetails(ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("BrokerAPIService.GetMorpheusDetails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("BrokerAPIService.GetMorpheusDetails() = %v, want %v", got, tt.want)
			}
		})
	}
}
