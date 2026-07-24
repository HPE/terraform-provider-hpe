// (C) Copyright 2021-2024 Hewlett Packard Enterprise Development LP

package serviceclient_test

import (
	"context"
	"errors"
	"log"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/HPE/terraform-provider-hpe/greenlake/sdk/token/mocks"
	"github.com/HPE/terraform-provider-hpe/greenlake/sdk/token/retrieve"
	"github.com/HPE/terraform-provider-hpe/greenlake/sdk/token/serviceclient"

	tokenutil "github.com/HPE/terraform-provider-hpe/greenlake/sdk/token/token-util"

	jose "github.com/go-jose/go-jose/v3"
	"github.com/go-jose/go-jose/v3/jwt"
	"github.com/stretchr/testify/assert"
)

func generateTestToken(timeToExpiry int64) string {
	timeNow := int64(0)
	pars := tokenutil.Token{
		Issuer:  "https://hpe-greenlake-tenant.okta.com/oauth2/default",
		Subject: "clients/subject",
		Expiry:  timeNow + timeToExpiry, IssuedAt: timeNow,
		ClientID: "clientID",
		TenantID: "tenantID",
	}

	sign, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.HS256, Key: []byte("secret")}, nil)
	if err != nil {
		log.Fatal(err)
	}

	retSign, err := jwt.Signed(sign).Claims(pars).CompactSerialize()
	if err != nil {
		log.Fatal()
	}

	return retSign
}

func TestHandler(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	ctx, cancel := context.WithCancel(context.Background())
	testcases := []struct {
		name              string
		token             string
		useAPIVendedToken bool
		err               error
		ctx               context.Context
		cancelFunc        context.CancelFunc
	}{
		{
			name:              "success api vended",
			token:             generateTestToken(600),
			useAPIVendedToken: true,
		},
		{
			name:              "success service client",
			token:             generateTestToken(600),
			useAPIVendedToken: false,
		},
		{
			name:  "no token",
			token: "",
			err:   errors.New("oidc: malformed jwt: square/go-jose: compact JWS format must have three parts"),
		},
		{
			name:  "renew token",
			token: generateTestToken(10),
		},
		{
			name: "network timeout",
			err:  testNetError{},
		},
		{
			name: "non-retryable error",
			err:  errors.New(""),
		},
		{
			name:       "cancelled context",
			ctx:        ctx,
			cancelFunc: cancel,
		},
	}
	for _, testcase := range testcases {
		tc := testcase
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mock := mocks.NewMockIdentityAPI(ctrl)

			testToken := generateTestToken(600)
			mock.EXPECT().GenerateToken(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(testToken, tc.err).MaxTimes(8)

			handler, err := serviceclient.NewHandler(
				serviceclient.WithIdentityAPI(mock),
				serviceclient.WithAPIVendedServiceClient(tc.useAPIVendedToken),
			)
			assert.NoError(t, err)
			if handler != nil {
				getToken := retrieve.NewTokenRetrieveFunc(handler)
				var token string
				var err error
				if !isNil(tc.ctx) {
					tc.cancelFunc()

					token, err = getToken(tc.ctx)
				} else {
					token, err = getToken(context.Background())
				}

				if tc.err != nil {
					assert.EqualError(t, err, tc.err.Error())
				}

				if tc.name != "renew token" {
					assert.Equal(t, tc.token, token)
				}
			}
		})
	}
}

// isNil we've have to add this function to avoid a Github action error
func isNil(i interface{}) bool {
	return i == nil || reflect.ValueOf(i).IsNil()
}

type testNetError struct {
}

func (e testNetError) Timeout() bool {
	return true
}

func (e testNetError) Temporary() bool {
	return true
}

func (e testNetError) Error() string {
	return ""
}
