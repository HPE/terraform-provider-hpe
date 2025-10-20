package testhelpers

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func CreateInstanceTypeLayout(t *testing.T, count int64) (
	[]sdk.GetInstanceType200ResponseInstanceTypeInstanceTypeLayoutsInner,
	error,
) {
	t.Helper()

	ctx := context.TODO()

	client := newClient(ctx, t)

	name := fmt.Sprintf("testacc-%s-%s", t.Name(), rand.Text())

	its, resp, err := client.LibraryAPI.ListInstanceTypes(ctx).Execute()
	if its == nil || err != nil || resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET failed for instance types %w", err)
	}

	if len(its.InstanceTypes) == 0 {
		return nil, errors.New("no instance type returned")
	}

	itID := its.InstanceTypes[len(its.InstanceTypes)-1].Id

	var layouts []sdk.GetInstanceType200ResponseInstanceTypeInstanceTypeLayoutsInner

	for i := range count {
		addLayout := sdk.NewAddLayoutRequestInstanceTypeLayoutWithDefaults()
		addLayout.SetName(name)
		addLayout.SetInstanceVersion("1")
		addLayout.SetProvisionTypeCode("kvm")
		addLayout.SetSortOrder(i)

		addLayoutReq := sdk.NewAddLayoutRequest()
		addLayoutReq.InstanceTypeLayout = addLayout

		req := client.LibraryAPI.AddLayout(ctx, *itID).AddLayoutRequest(*addLayoutReq)

		l, resp, err := req.Execute()
		if err != nil || resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("POST failed for instance layout %w", err)
		}

		layouts = append(layouts, l.GetInstanceTypeLayout())
	}

	return layouts, nil
}

func DeleteInstanceTypeLayout(t *testing.T, id int64) error {
	t.Helper()

	ctx := context.TODO()

	client := newClient(ctx, t)

	_, resp, err := client.LibraryAPI.DeleteLayout(ctx, id).Execute()
	if err != nil || resp.StatusCode != http.StatusOK {
		return fmt.Errorf("DELETE failed for instance layout %d: %w", id, err)
	}

	for range 6 {
		_, resp, _ := client.LibraryAPI.GetLayout(ctx, id).Execute()
		if resp.StatusCode == http.StatusNotFound {
			return nil
		}

		t.Log("Waiting for instance layout to be deleted")
		time.Sleep(time.Second * 10)
	}

	return fmt.Errorf("DELETE failed for instance layout %d: %w", id, err)
}
