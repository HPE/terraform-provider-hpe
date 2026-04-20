// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package ostype

import (
	"context"
	"net/http"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

const testOsTypePrefix = "TestAccMorpheusOsType"

func init() {
	testhelpers.RegisterAPISweeper(
		"hpe_morpheus_os_type",
		testOsTypePrefix,
		func(ctx context.Context, client *sdk.APIClient, prefix string) (any, *http.Response, error) {
			return client.LibraryAPI.ListOsTypes(ctx).Phrase(prefix).Execute()
		},
		"GetOsTypes",
		func(ctx context.Context, client *sdk.APIClient, id int64, _ any) (*http.Response, error) {
			_, hresp, err := client.LibraryAPI.DeleteOsType(ctx, id).Execute()
			return hresp, err
		},
	)
}
