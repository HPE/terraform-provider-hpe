package task

import (
	"context"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TaskModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.Id.ValueInt64()
	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("creating client failed", err.Error())

		return
	}

	deleteReq := client.AutomationAPI.RemoveTasks(ctx, id)
	result, hresp, err := deleteReq.Execute()
	if err := errfmt.CheckDeleted(result, err, hresp); err != nil {
		resp.Diagnostics.AddError(
			"delete task resource",
			fmt.Sprintf("task %d: DELETE failed ", id)+err.Error(),
		)

		return
	}

	// Wait until the task can no longer be read; a successful GET means it is
	// still there and the loop tries again (see errfmt.PollDeleted).
	what := fmt.Sprintf("task %d", id)
	waitForDeleted := func() (struct{}, error) {
		_, httpResp, err := client.AutomationAPI.GetTasks(ctx, id).Execute()

		return struct{}{}, errfmt.PollDeleted(what, err, httpResp)
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
