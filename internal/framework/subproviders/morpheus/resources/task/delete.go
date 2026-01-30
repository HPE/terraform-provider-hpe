package task

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/errors"
)

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TaskModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.Id
	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("creating client failed", err.Error())

		return
	}

	deleteReq := client.AutomationAPI.RemoveTasks(ctx, id.ValueInt64())
	_, hresp, err := deleteReq.Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"delete task resource",
			fmt.Sprintf("task %d: DELETE failed ", id)+errors.ErrMsg(err, hresp),
		)

		return
	}

	waitForDeleted := func() (*sdk.GetInstance200Response, error) {
		_, httpResp, err := client.AutomationAPI.GetTasks(ctx, id.ValueInt64()).Execute()
		// 404 status code counts as a successful delete
		if err != nil && httpResp.StatusCode != http.StatusNotFound {
			return nil, backoff.Permanent(err)
		}

		return nil, nil
	}

	if _, err := backoff.Retry(
		ctx,
		waitForDeleted,
		backoff.WithBackOff(backoff.NewConstantBackOff(5*time.Second)),
		backoff.WithMaxElapsedTime(45*time.Minute),
	); err != nil {
		resp.Diagnostics.AddError(
			"delete task resource",
			fmt.Sprintf("task %d: DELETE failed ", id)+err.Error(),
		)
	}
}
