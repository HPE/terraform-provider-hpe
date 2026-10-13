// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

//go:build sweep

package sweep

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	testsweep "github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/sweep"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/getsafe"
)

const sweeperName = "hpe_morpheus_task_powershell_script"

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

			return getsafe.Get(&resp.Tasks), hresp, err
		},
		// Is this a test task of this type? Both the name prefix and the task
		// type code must match, so the sweeper never removes another test's
		// fixture task of a different type.
		func(item sdk.ListTasks200ResponseAllOfTasksInner) bool {
			name, _, typeCode, ok := taskNameAndID(item)
			if !ok {
				return false
			}

			return strings.HasPrefix(name, testsweep.TestResourcePrefix) && typeCode == "winrmTask"
		},
		// Delete the test task.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item sdk.ListTasks200ResponseAllOfTasksInner,
		) (*http.Response, error) {
			_, id, _, ok := taskNameAndID(item)
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

func taskNameAndID(item sdk.ListTasks200ResponseAllOfTasksInner) (name string, id int64, typeCode string, ok bool) {
	payload, err := json.Marshal(item)
	if err != nil {
		return "", 0, "", false
	}

	var parsed struct {
		ID       *int64  `json:"id"`
		Name     *string `json:"name"`
		TaskType *struct {
			Code *string `json:"code"`
		} `json:"taskType"`
	}

	if err := json.Unmarshal(payload, &parsed); err != nil {
		return "", 0, "", false
	}

	if parsed.ID == nil || parsed.Name == nil {
		return "", 0, "", false
	}

	typeCode = ""
	if parsed.TaskType != nil && parsed.TaskType.Code != nil {
		typeCode = *parsed.TaskType.Code
	}

	return *parsed.Name, *parsed.ID, typeCode, true
}
