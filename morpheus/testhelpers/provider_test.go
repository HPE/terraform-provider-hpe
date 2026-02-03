// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package testhelpers_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/provider/subprovider"
)

type providerTestCase struct {
	name              string
	frameworkProvider subprovider.SubProvider
	sdvkv2Provider    *schema.Provider
	expectNoFactories bool
}

func TestGetAccTestFactories(t *testing.T) {
	t.Parallel()

	tests := []providerTestCase{
		{"Framework provider only success", morpheus.New(), nil, false},
		{"SDK v2 provider only success", nil, sdkv2morpheus.Provider(), false},
		{"Both providers success", morpheus.New(), sdkv2morpheus.Provider(), false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			factories := testhelpers.GetAccTestFactories(t, test.frameworkProvider, test.sdvkv2Provider)

			if test.expectNoFactories && factories != nil {
				t.Error("expected factories to be nil")
			}

			if !test.expectNoFactories && factories == nil {
				t.Error("expected factories to be a map")
			}
		})
	}
}
