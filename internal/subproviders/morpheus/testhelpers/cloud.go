// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package testhelpers

import (
	"context"
	"crypto/rand"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func CreateCloud(t *testing.T, groupID int64) sdk.ListClouds200ResponseAllOfZonesInner {
	t.Helper()

	ctx := context.TODO()

	client := newClient(ctx, t)

	name := fmt.Sprintf("testacc-%s-%s", t.Name(), rand.Text())

	zts, hresp, err := client.CloudsAPI.ListCloudTypes(ctx).Execute()
	if zts == nil || err != nil || hresp.StatusCode != http.StatusOK {
		t.Fatalf("GET failed for cloud types %v", err)
	}

	if len(zts.ZoneTypes) < 1 {
		t.Fatal("no cloud type returned")
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
		t.Fatalf("POST failed for cloud %v", err)
	}

	cloud := c.GetZone()

	return cloud
}

func DeleteCloud(t *testing.T, id int64) {
	t.Helper()

	ctx := context.TODO()

	client := newClient(ctx, t)

	_, resp, err := client.CloudsAPI.RemoveClouds(ctx, id).Force(true).Execute()
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("DELETE failed for cloud %d: %v", id, err)
	}

	for range 6 {
		_, resp, _ := client.CloudsAPI.GetClouds(ctx, id).Execute()
		if resp.StatusCode == http.StatusNotFound {
			return
		}

		t.Log("Waiting for cloud to be deleted")
		time.Sleep(time.Second * 10)
	}

	t.Fatalf("DELETE failed for cloud %d: %v", id, err)
}
