// (C) Copyright 2026 Hewlett Packard Enterprise Development LP


//go:build sweep

package sweep

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	testsweep "github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/sweep"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/getsafe"
)

const sweeperName = "hpe_morpheus_task"

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List task resources.
		func(ctx context.Context, client *sdk.APIClient) (
			[]sdk.ListTasks200ResponseAllOfTasksInner,
			*http.Response,
			error,
		) {
			resp, hresp, err := client.AutomationAPI.ListTasks(ctx).Execute()
			if resp == nil {
				return nil, hresp, err
			}

			return getsafe.GetSafe(&resp.Tasks), hresp, err
		},
		// Is this a test task?
		func(item sdk.ListTasks200ResponseAllOfTasksInner) bool {
			name, _, ok := getTaskNameAndID(item)
			if !ok {
				return false
			}

			return strings.HasPrefix(name, testsweep.TestResourcePrefix)
		},
		// Delete the test task.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item sdk.ListTasks200ResponseAllOfTasksInner,
		) (*http.Response, error) {
			_, id, ok := getTaskNameAndID(item)
			if !ok {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.AutomationAPI.RemoveTasks(ctx, id).Execute()

			return hresp, err
		},
		testsweep.WithIgnoreListStatuses[sdk.ListTasks200ResponseAllOfTasksInner](
			http.StatusNotFound,
			http.StatusForbidden,
		),
	)
}

func getTaskNameAndID(item sdk.ListTasks200ResponseAllOfTasksInner) (string, int64, bool) {
	// Task list responses are polymorphic; parse common top-level fields.
	payload, err := json.Marshal(item)
	if err != nil {
		return "", 0, false
	}

	var parsed struct {
		ID   *int64  `json:"id"`
		Name *string `json:"name"`
	}

	if err := json.Unmarshal(payload, &parsed); err != nil {
		return "", 0, false
	}

	if parsed.ID == nil || parsed.Name == nil {
		return "", 0, false
	}

	return *parsed.Name, *parsed.ID, true
}
