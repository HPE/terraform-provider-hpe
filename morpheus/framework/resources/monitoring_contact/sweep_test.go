// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package monitoring_contact_test

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	testsweep "github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/sweep"
)

const sweeperName = "hpe_morpheus_monitoring_contact"

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List monitoring contact resources.
		func(ctx context.Context, client *sdk.APIClient) (
			[]sdk.ListContacts200ResponseAllOfContactsInner,
			*http.Response,
			error,
		) {
			resp, hresp, err := client.ContactsAPI.ListContacts(ctx).Execute()
			if resp == nil {
				return nil, hresp, err
			}

			return resp.GetContacts(), hresp, err
		},
		// Is this a test monitoring contact?
		func(item sdk.ListContacts200ResponseAllOfContactsInner) bool {
			name, ok := item.GetNameOk()
			if !ok || name == nil {
				return false
			}

			return strings.HasPrefix(*name, testsweep.TestResourcePrefix)
		},
		// Delete the test monitoring contact.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item sdk.ListContacts200ResponseAllOfContactsInner,
		) (*http.Response, error) {
			id, ok := item.GetIdOk()
			if !ok || id == nil {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.ContactsAPI.DeleteContacts(ctx, *id).Execute()

			return hresp, err
		},
		testsweep.WithIgnoreListStatuses[sdk.ListContacts200ResponseAllOfContactsInner](
			http.StatusNotFound,
			http.StatusForbidden,
		),
	)
}
