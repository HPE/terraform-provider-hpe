// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package testhelpers

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func CreateCloud(t *testing.T, groupID int64) (*sdk.ListClouds200ResponseAllOfZonesInner, error) {
	t.Helper()

	ctx := context.TODO()

	client := newClient(ctx, t)

	name := fmt.Sprintf("testacc-%s-%s", t.Name(), rand.Text())

	zts, hresp, err := client.CloudsAPI.ListCloudTypes(ctx).Execute()
	if zts == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET failed for cloud types %w", err)
	}

	if len(zts.ZoneTypes) < 1 {
		return nil, errors.New("no cloud type returned")
	}

	ztID := zts.ZoneTypes[len(zts.ZoneTypes)-1].Id

	zt := sdk.AddCloudsRequestZoneZoneType{
		AddCloudsRequestZoneZoneTypeAnyOf: &sdk.AddCloudsRequestZoneZoneTypeAnyOf{
			Id: ztID,
		},
	}

	addCloud := sdk.NewAddCloudsRequestZoneWithDefaults()
	addCloud.SetName(name)
	addCloud.SetCode(strings.ToLower(name))
	addCloud.SetLocation("here")
	addCloud.SetGroupId(groupID)
	addCloud.SetZoneType(zt)

	addCloudReq := sdk.NewAddCloudsRequest(*addCloud)

	c, hresp, err := client.CloudsAPI.AddClouds(ctx).AddCloudsRequest(*addCloudReq).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("POST failed for cloud %w", err)
	}

	cloud := c.GetZone()

	return &cloud, nil
}

func DeleteCloud(t *testing.T, id int64) error {
	t.Helper()

	ctx := context.TODO()

	client := newClient(ctx, t)

	_, resp, err := client.CloudsAPI.RemoveClouds(ctx, id).Force(true).Execute()
	if err != nil || resp.StatusCode != http.StatusOK {
		return fmt.Errorf("DELETE failed for cloud %d: %w", id, err)
	}

	for range 6 {
		_, resp, _ := client.CloudsAPI.GetClouds(ctx, id).Execute()
		if resp.StatusCode == http.StatusNotFound {
			return nil
		}

		t.Log("Waiting for cloud to be deleted")
		time.Sleep(time.Second * 10)
	}

	return fmt.Errorf("DELETE failed for cloud %d: %w", id, err)
}
