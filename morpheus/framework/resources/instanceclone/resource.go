// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instanceclone

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/cenkalti/backoff/v5"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &Resource{}
	_ resource.ResourceWithImportState = &Resource{}
)

var (
	cloneTargetStatuses = []string{"running", "stopped"}
	cloneErrorStatuses  = []string{"failed", "errored"}
)

func NewResource() resource.Resource {
	return &Resource{}
}

// Resource defines the resource implementation.
type Resource struct {
	configure.ResourceWithMorpheusConfigure
}

func (r *Resource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_instance_clone"
}

func (r *Resource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = InstanceCloneResourceSchema(ctx)
}

func checkStatusDone(status string, targetStatuses, errorStatuses []string) error {
	for _, s := range errorStatuses {
		if strings.EqualFold(status, s) {
			return backoff.Permanent(errors.New("reached error status: " + status))
		}
	}
	for _, s := range targetStatuses {
		if strings.EqualFold(status, s) {
			return nil
		}
	}

	return fmt.Errorf("instance status %q not yet in target set", status)
}

// pollAPIError classifies an error hit while polling. Client errors (4xx) are
// permanent - retrying will not help - while 5xx responses and transport errors
// are returned as-is so backoff retries them.
func pollAPIError(msg string, err error, hresp *http.Response) error {
	e := fmt.Errorf("%s: %s", msg, errfmt.ErrMsg(err, hresp))
	if hresp != nil && hresp.StatusCode >= 400 && hresp.StatusCode < 500 {
		return backoff.Permanent(e)
	}

	return e
}
