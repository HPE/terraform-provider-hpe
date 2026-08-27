package testhelpers

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

func GetRole(t *testing.T, roleID int64) (*sdk.GetRole200Response, error) {
	t.Helper()

	ctx := context.TODO()

	client := newClient(ctx, t)

	r, hresp, err := client.RolesAPI.GetRole(ctx, roleID).Execute()
	if r == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET failed for role %w", err)
	}

	return r, nil
}
