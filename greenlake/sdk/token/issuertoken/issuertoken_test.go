// (C) Copyright 2024 Hewlett Packard Enterprise Development LP

package issuertoken

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/HPE/terraform-provider-hpe/greenlake/sdk/token/iamversion"
)

func generateExpParams(iamVersion iamversion.Version) url.Values {
	expParams := url.Values{}
	expParams.Add("client_id", "clientID")
	expParams.Add("client_secret", "clientSecret")
	expParams.Add("grant_type", "client_credentials")
	if iamVersion == iamversion.GLCS {
		expParams.Add("scope", "hpe-tenant")
	}

	return expParams
}

func TestGenerateParamsAndURL(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name       string
		iamVersion iamversion.Version
		expParams  url.Values
		hasError   bool
	}{
		{
			name:       "valid IAM version GLCS",
			iamVersion: iamversion.GLCS,
			expParams:  generateExpParams(iamversion.GLCS),
			hasError:   false,
		},
		{
			name:       "valid IAM version GLP",
			iamVersion: iamversion.GLP,
			expParams:  generateExpParams(iamversion.GLP),
			hasError:   false,
		},
		{
			name:       "invalid IAM version",
			iamVersion: "invalid",
			expParams:  nil,
			hasError:   true,
		},
	}

	for _, testcase := range testcases {
		tc := testcase
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			params, _, err := generateParamsAndURL("clientID", "clientSecret", "identityServiceURL", string(tc.iamVersion))
			assert.Equal(t, tc.expParams, params)
			if tc.hasError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
