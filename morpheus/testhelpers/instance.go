// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package testhelpers

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

const (
	instancePollInterval = 10 * time.Second
	instancePollTimeout  = 10 * time.Minute
)

// CreateInstance provisions an instance by cloning the shape of an existing running
// instance that supports the add-node action. This guarantees a mutually valid
// combination of group, cloud, instance type, layout, plan, network, and resource pool.
// It polls until the instance reaches "running" status and returns the instance ID.
func CreateInstance(t *testing.T) (int64, error) {
	t.Helper()

	ctx := context.TODO()
	client := newClient(ctx, t)

	name := fmt.Sprintf("TestAccMorpheus-%s-%s", t.Name(), rand.Text())

	// Step 1: List running instances.
	listResp, httpResp, err := client.InstancesAPI.ListInstances(ctx).
		Status("running").Max(50).Execute()
	if err != nil || httpResp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to list running instances: %w", err)
	}
	if len(listResp.Instances) == 0 {
		return 0, errors.New(
			"no running instances found on appliance — cannot construct a fixture " +
				"(CreateInstance requires a running template instance with an add-node action)")
	}

	// Step 2: Find a running instance whose actions include an add-node action.
	type templateShape struct {
		groupID          int64
		cloudID          int64
		instanceTypeCode string
		layoutID         int64
		planID           int64
		networkID        int64
		resourcePoolID   string
		poolProviderType string
	}

	var tmpl *templateShape

	for i := range listResp.Instances {
		inst := &listResp.Instances[i]
		if inst.Id == nil {
			continue
		}

		// Check actions for this instance.
		actionsResp, actHTTP, actErr := client.InstancesAPI.
			GetInstanceActions(ctx, *inst.Id).Execute()
		if actErr != nil || actHTTP.StatusCode != http.StatusOK {
			continue
		}

		hasAddNode := false
		for j := range actionsResp.Actions {
			action := &actionsResp.Actions[j]
			if action.Code != nil && strings.HasSuffix(*action.Code, "-add-node") {
				hasAddNode = true

				break
			}
			if action.Name != nil && *action.Name == "Add Node" {
				hasAddNode = true

				break
			}
		}
		if !hasAddNode {
			continue
		}

		// Extract shape from this template instance.
		shape := templateShape{}

		// Group
		if g := inst.Group.Get(); g != nil && g.Id != nil {
			shape.groupID = *g.Id
		} else {
			continue
		}

		// Cloud
		if inst.Cloud == nil || inst.Cloud.Id == nil {
			continue
		}
		shape.cloudID = *inst.Cloud.Id

		// Instance type
		if inst.InstanceType == nil || inst.InstanceType.Code == nil {
			continue
		}
		shape.instanceTypeCode = *inst.InstanceType.Code

		// Layout
		if inst.Layout == nil || inst.Layout.Id == nil {
			continue
		}
		shape.layoutID = *inst.Layout.Id

		// Plan
		if inst.Plan == nil || inst.Plan.Id == nil {
			continue
		}
		shape.planID = *inst.Plan.Id

		// Network — from the first interface
		if len(inst.Interfaces) > 0 {
			iface := &inst.Interfaces[0]
			if iface.Network != nil {
				if nid := iface.Network.Id.Get(); nid != nil {
					shape.networkID = *nid
				}
			}
		}
		if shape.networkID == 0 {
			continue
		}

		// Resource pool and pool provider type from config
		if inst.Config != nil {
			if inst.Config.ResourcePoolId != nil {
				if inst.Config.ResourcePoolId.String != nil {
					shape.resourcePoolID = *inst.Config.ResourcePoolId.String
				} else if inst.Config.ResourcePoolId.Int64 != nil {
					shape.resourcePoolID = fmt.Sprintf("pool-%d",
						*inst.Config.ResourcePoolId.Int64)
				}
			}
			if ppt := inst.Config.PoolProviderType.Get(); ppt != nil {
				shape.poolProviderType = *ppt
			}
		}

		tmpl = &shape

		break
	}

	if tmpl == nil {
		return 0, errors.New(
			"no running instance with an add-node action found on appliance — " +
				"the test cannot construct a fixture without a suitable template instance")
	}

	t.Logf("Using template shape: group=%d cloud=%d type=%s layout=%d plan=%d "+
		"network=%d pool=%s",
		tmpl.groupID, tmpl.cloudID, tmpl.instanceTypeCode, tmpl.layoutID,
		tmpl.planID, tmpl.networkID, tmpl.resourcePoolID)

	// Step 3: Build the create request from the template shape.
	noAgent := true
	createUser := false
	hvmConfig := &sdk.HVMInstanceConfiguration{
		AdditionalProperties: map[string]interface{}{
			"createBackup": false,
			"layoutSize":   1,
		},
	}
	hvmConfig.NoAgent.Set(&noAgent)
	hvmConfig.CreateUser.Set(&createUser)
	if tmpl.poolProviderType != "" {
		hvmConfig.PoolProviderType = &tmpl.poolProviderType
	}
	if tmpl.resourcePoolID != "" {
		hvmConfig.ResourcePoolId = &tmpl.resourcePoolID
	}

	rootVolume := true
	volName := "root"
	var volSize int64 = 10
	var volID int64 = -1
	datastoreAuto := "auto"

	networkIDStr := fmt.Sprintf("%d", tmpl.networkID)

	vol := sdk.AddInstanceRequestVolumesInner{
		Id:         &volID,
		RootVolume: &rootVolume,
		Name:       &volName,
		Size:       &volSize,
		DatastoreId: &sdk.InstanceConfigObject1VolumesInnerDatastoreId{
			String: &datastoreAuto,
		},
	}
	var storageType int64 = 1
	vol.StorageType.Set(&storageType)

	req := sdk.AddInstanceRequest{
		Instance: sdk.AddInstanceRequestInstance{
			Name: name,
			Site: sdk.AddInstanceRequestInstanceSite{Id: tmpl.groupID},
			InstanceType: sdk.AddInstanceRequestInstanceInstanceType{
				Code: tmpl.instanceTypeCode,
			},
			Layout: sdk.AddInstanceRequestInstanceLayout{Id: tmpl.layoutID},
			Plan:   sdk.AddInstanceRequestInstancePlan{Id: tmpl.planID},
		},
		ZoneId: &tmpl.cloudID,
		Config: sdk.AddInstanceRequestConfig{
			HVMInstanceConfiguration: hvmConfig,
		},
		Volumes: []sdk.AddInstanceRequestVolumesInner{vol},
		NetworkInterfaces: []sdk.InstancesNetworkInterfaces2{
			{
				Network: sdk.InstancesNetworkInterfaces1Network{
					Id: networkIDStr,
				},
			},
		},
	}

	resp, httpResp, err := client.InstancesAPI.AddInstance(ctx).
		AddInstanceRequest(req).Execute()
	if err != nil || httpResp.StatusCode != http.StatusOK {
		// Surface the appliance's own message for diagnosability.
		body := ""

		if httpResp != nil && httpResp.Body != nil {
			b, readErr := io.ReadAll(httpResp.Body)
			if readErr == nil {
				body = strings.TrimSpace(string(b))
			}
		}

		reqJSON, _ := json.Marshal(req)

		return 0, fmt.Errorf(
			"failed to create instance (HTTP %d): %w\n  response: %s\n  request: %s",
			httpResp.StatusCode, err, body, string(reqJSON),
		)
	}
	if resp == nil || resp.Instance.Id == nil {
		return 0, errors.New("create instance returned no id")
	}

	instanceID := *resp.Instance.Id
	t.Logf("Created instance %d (%s), polling for running status", instanceID, name)

	// Poll until running.
	deadline := time.Now().Add(instancePollTimeout)
	for time.Now().Before(deadline) {
		getResp, _, getErr := client.InstancesAPI.GetInstance(ctx, instanceID).Execute()
		if getErr == nil && getResp != nil && getResp.Instance != nil &&
			getResp.Instance.Status != nil {
			status := *getResp.Instance.Status
			switch status {
			case "running":
				t.Logf("Instance %d is running", instanceID)

				return instanceID, nil
			case "failed":
				return instanceID, fmt.Errorf(
					"instance %d reached 'failed' status — check appliance logs",
					instanceID)
			}
			t.Logf("Instance %d status: %s, waiting...", instanceID, status)
		}
		time.Sleep(instancePollInterval)
	}

	return instanceID, fmt.Errorf(
		"instance %d did not reach running within %v", instanceID, instancePollTimeout)
}

// DeleteInstance force-deletes an instance with volume removal and polls until gone.
func DeleteInstance(t *testing.T, id int64) error {
	t.Helper()

	ctx := context.TODO()
	client := newClient(ctx, t)

	_, httpResp, err := client.InstancesAPI.DeleteInstance(ctx, id).
		Force("on").RemoveVolumes("on").Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			return nil // already gone
		}

		return fmt.Errorf("DELETE failed for instance %d: %w", id, err)
	}

	// Poll until gone.
	for range 30 {
		_, resp, _ := client.InstancesAPI.GetInstance(ctx, id).Execute()
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			t.Logf("Instance %d deleted", id)

			return nil
		}
		t.Logf("Waiting for instance %d to be deleted", id)
		time.Sleep(instancePollInterval)
	}

	return fmt.Errorf("instance %d was not deleted within timeout", id)
}
