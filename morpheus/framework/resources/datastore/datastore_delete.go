// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package datastore

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var data DatastoreModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	id := data.Id.ValueInt64()

	client, _ := r.NewClient(ctx)

	_, hresp, err := client.DatastoresAPI.DeleteDatastores(ctx, id).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"delete datastore resource",
			fmt.Sprintf("datastore %d: DELETE failed ", id)+errfmt.ErrMsg(err, hresp),
		)

		return
	}

	// Wait for the datastore to be "not found" which means it has been deleted
	waitForReady := func() (int, error) {
		_, hresp, err := client.DatastoresAPI.GetDatastores(ctx, id).Execute()
		if err != nil {
			if hresp == nil || hresp.StatusCode != http.StatusNotFound {
				return 0, backoff.Permanent(err)
			}
		}

		switch hresp.StatusCode {
		case http.StatusNotFound:
			return hresp.StatusCode, nil
		default:
			return hresp.StatusCode, backoff.RetryAfter(5)
		}
	}

	if statusCode, err := backoff.Retry(
		ctx,
		waitForReady,
		backoff.WithBackOff(backoff.NewConstantBackOff(5*time.Second)),
		backoff.WithMaxElapsedTime(5*time.Minute),
	); err != nil {
		resp.Diagnostics.AddError(
			"delete datastore resource",
			fmt.Sprintf(
				"datastore %d: DELETE failed current statusCode is: %v",
				id,
				statusCode,
			),
		)
	}
}
