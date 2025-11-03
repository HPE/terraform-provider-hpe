// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

//go:build experimental

package policy

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/errors"
)

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var data PolicyModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.Id.ValueInt64()

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"delete policy resource",
			fmt.Sprintf("policy %d: failed to create client: %s", id, err.Error()),
		)

		return
	}

	_, hresp, err := client.PoliciesAPI.RemovePolicies(ctx, id).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"delete policy resource",
			fmt.Sprintf("policy %d: DELETE failed ", id)+errors.ErrMsg(err, hresp),
		)

		return
	}
}
