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

	name := fmt.Sprintf("testacc-%s-%s", t.Name(), rand.Text())

	addCloud := sdk.NewAddCloudsRequestZoneWithDefaults()
	addCloud.SetName(name)
	addCloud.SetCode(strings.ToLower(name))
	addCloud.SetLocation("here")
	addCloud.SetGroupId(groupID)

	// This is the ID of a Morpheus zone type
	ztID := int64(1)
	zt := sdk.AddCloudsRequestZoneZoneType{
		AddCloudsRequestZoneZoneTypeAnyOf: &sdk.AddCloudsRequestZoneZoneTypeAnyOf{
			Id: &ztID,
		},
	}

	addCloud.SetZoneType(zt)

	addCloudReq := sdk.NewAddCloudsRequest(*addCloud)

	ctx := context.TODO()

	client := newClient(ctx, t)

	c, hresp, err := client.CloudsAPI.AddClouds(ctx).AddCloudsRequest(*addCloudReq).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		t.Fatalf("POST failed for group %v", err)
	}

	cloud := c.GetZone()

	return cloud
}

func DeleteCloud(t *testing.T, id int64) {
	t.Helper()

	ctx := context.TODO()

	client := newClient(ctx, t)

	counter := 0

	// Clouds go into initializing state and can't be deleted immediately
	_, _, err := client.CloudsAPI.RemoveClouds(ctx, id).Execute()
	for counter < 10 && err != nil {
		time.Sleep(5 * time.Second)
		counter++

		_, _, err = client.CloudsAPI.RemoveClouds(ctx, id).Execute()
	}
}
