// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package testhelpers

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// DeleteOutOfBand returns a TestStep that deletes the resource at
// resourceAddress directly through the SDK — bypassing Terraform — and then
// refreshes, asserting the refresh leaves a non-empty plan.
//
// It reproduces a resource deleted outside Terraform (MORPH-18584): after the
// out-of-band delete the next Read must answer the API's 404 by removing the
// resource from state, so the refreshed plan proposes recreating it rather than
// failing. A Read that reported the 404 as an error would fail this step; a Read
// that left stale state would produce an empty plan and also fail it.
//
// deleteFn performs the SDK delete for the resource's own API; it receives the
// id parsed from the resource's "id" attribute in state. The normal destroy at
// the end of the test still runs and tolerates the already-deleted object.
func DeleteOutOfBand(
	t *testing.T,
	state **terraform.State,
	resourceAddress string,
	deleteFn func(ctx context.Context, client *sdk.APIClient, id int64) error,
) resource.TestStep {
	t.Helper()

	return resource.TestStep{
		RefreshState:       true,
		ExpectNonEmptyPlan: true,
		PreConfig: func() {
			ctx := context.Background()
			client := newClient(ctx, t)

			id, err := idFromState(*state, resourceAddress)
			if err != nil {
				t.Fatalf("DeleteOutOfBand: %v", err)
			}

			if err := deleteFn(ctx, client, id); err != nil {
				t.Fatalf("DeleteOutOfBand: deleting %s (id %d) failed: %v", resourceAddress, id, err)
			}
		},
	}
}

// StashState returns a Check that records the current state into dst so a
// following DeleteOutOfBand step can look up the resource's id; the testing
// framework does not pass state into PreConfig. Each test owns its own dst
// variable, so parallel tests do not share it.
func StashState(dst **terraform.State) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		*dst = s

		return nil
	}
}

func idFromState(state *terraform.State, resourceAddress string) (int64, error) {
	if state == nil {
		return 0, fmt.Errorf("no state stashed; add testhelpers.StashState(&state) to the preceding step's Check")
	}

	rs, ok := state.RootModule().Resources[resourceAddress]
	if !ok {
		return 0, fmt.Errorf("resource not found in state: %s", resourceAddress)
	}

	raw, ok := rs.Primary.Attributes["id"]
	if !ok {
		return 0, fmt.Errorf("resource %s has no id attribute", resourceAddress)
	}

	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("resource %s id %q is not an integer: %w", resourceAddress, raw, err)
	}

	return id, nil
}
