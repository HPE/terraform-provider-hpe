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

func CreateGroup(t *testing.T) (*sdk.ListGroups200ResponseAllOfGroupsInner, error) {
	t.Helper()

	name := fmt.Sprintf("testacc-%s-%s", t.Name(), rand.Text())

	addGroup := sdk.NewAddGroupsRequestGroupWithDefaults()
	addGroup.SetName(name)
	addGroup.SetCode(strings.ToLower(name))
	addGroup.SetLocation("here")

	addGroupReq := sdk.NewAddGroupsRequest(*addGroup)

	ctx := context.TODO()

	client := newClient(ctx, t)

	g, hresp, err := client.GroupsAPI.AddGroups(ctx).AddGroupsRequest(*addGroupReq).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("POST failed for group %w", err)
	}

	group := g.GetGroup()

	return &group, nil
}

func DeleteGroup(t *testing.T, id int64) error {
	t.Helper()

	ctx := context.TODO()

	client := newClient(ctx, t)

	_, resp, err := client.GroupsAPI.RemoveGroups(ctx, id).Execute()
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("DELETE failed for group %d: %v", id, err)
	}

	for range 6 {
		_, resp, _ := client.GroupsAPI.GetGroups(ctx, id).Execute()
		if resp.StatusCode == http.StatusNotFound {
			return nil
		}

		t.Log("Waiting for group to be deleted")
		time.Sleep(time.Second * 10)
	}

	return fmt.Errorf("DELETE failed for group %d: %w", id, err)
}
