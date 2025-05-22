package testhelpers

import (
	"context"
	"crypto/rand"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func CreateGroup(t *testing.T) sdk.ListGroups200ResponseAllOfGroupsInner {
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
		t.Fatalf("POST failed for group %v", err)
	}

	group := g.GetGroup()

	return group
}

func DeleteGroup(t *testing.T, id int64) {
	t.Helper()

	ctx := context.TODO()

	client := newClient(ctx, t)

	_, hresp, err := client.GroupsAPI.RemoveGroups(ctx, id).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		t.Fatalf("DELETE failed for group %v", err)
	}
}
