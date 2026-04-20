// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package instance

import (
	"context"
	"fmt"
	"net/http"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

// Instances whose name begins with this string will be eligible for deletion
const testInstancePrefix = "TestAccMorpheusInstance"

// All of these tags must be present with value "true" for the instance to be deleted
var requiredSweepTags = []string{
	"terraform",
	"acctest",
	"hpe_morpheus_instance",
	"sweepable",
}

func hasRequiredTags(tags []sdk.AddInstance200ResponseAllOfOneOfInstanceTagsInner) bool {
	if tags == nil {
		return false
	}

	tagMap := make(map[string]string)
	for _, tag := range tags {
		if name, ok := tag.GetNameOk(); ok && name != nil {
			if value, ok := tag.GetValueOk(); ok && value != nil {
				tagMap[*name] = *value
			}
		}
	}

	for _, requiredTag := range requiredSweepTags {
		if value, exists := tagMap[requiredTag]; !exists || value != "true" {
			return false
		}
	}

	return true
}

func init() {
	testhelpers.RegisterAPISweeper(
		"hpe_morpheus_instance",
		testInstancePrefix,
		func(ctx context.Context, client *sdk.APIClient, prefix string) (any, *http.Response, error) {
			return client.InstancesAPI.ListInstances(ctx).Phrase(prefix).Execute()
		},
		"GetInstances",
		func(ctx context.Context, client *sdk.APIClient, id int64, item any) (*http.Response, error) {
			instance := item.(sdk.Instance)
			serverIDs, err := getServerIDs(instance)
			if err != nil {
				return nil, err
			}

			for _, serverID := range serverIDs {
				stopID := sdk.UpdateHostIdParameter{Int64: &serverID}
				_, hresp, err := client.HostsAPI.StopHost(ctx, stopID).Execute()
				if err != nil || (hresp.StatusCode != http.StatusOK && hresp.StatusCode != http.StatusConflict) {
					return hresp, err
				}
			}

			for _, serverID := range serverIDs {
				deleteID := sdk.UpdateHostIdParameter{Int64: &serverID}
				_, hresp, err := client.HostsAPI.RemoveHost(ctx, deleteID).Force("on").
					RemoveResources("on").RemoveInstances("on").Execute()
				if err != nil || hresp.StatusCode != http.StatusOK {
					return hresp, err
				}
			}

			return &http.Response{StatusCode: http.StatusOK}, nil
		},
		testhelpers.WithFilter(func(ctx context.Context, client *sdk.APIClient, item any) (bool, string, error) {
			instance := item.(sdk.Instance)
			id, ok := instance.GetIdOk()
			if !ok || id == nil {
				return false, "id", nil
			}

			instanceDetail, hresp, err := client.InstancesAPI.GetInstance(ctx, *id).Execute()
			if err != nil || hresp == nil || hresp.StatusCode != http.StatusOK {
				return false, "failed to get details", nil
			}

			detail := instanceDetail.GetInstance()
			tags, ok := detail.GetTagsOk()
			if !ok || !hasRequiredTags(tags) {
				return false, "tags", nil
			}

			return true, "", nil
		}),
	)
}

func getServerIDs(instance sdk.Instance) ([]int64, error) {
	containers, ok := instance.GetContainerDetailsOk()
	if !ok || containers == nil {
		return nil, fmt.Errorf("failed to get container details")
	}

	serverIDs := make([]int64, 0, len(containers))
	for _, container := range containers {
		server, ok := container.GetServerOk()
		if !ok || server == nil {
			return nil, fmt.Errorf("failed to get server details from container")
		}

		serverID, ok := server.GetIdOk()
		if !ok || serverID == nil {
			return nil, fmt.Errorf("failed to get server ID from container")
		}

		serverIDs = append(serverIDs, *serverID)
	}

	return serverIDs, nil
}
