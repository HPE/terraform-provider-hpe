package image

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

func planLocalFile(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var config ImageModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if config.FilePath.IsUnknown() {
		return
	}
	if config.FilePath.IsNull() {
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("file_sha256"), types.StringNull())...)
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("file_name"), types.StringNull())...)
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("file_size"), types.Int64Null())...)

		return
	}
	name := filepath.Base(config.FilePath.ValueString())
	if name == "." || name == ".." || !filenameRE.MatchString(name) {
		resp.Diagnostics.AddAttributeError(path.Root("file_path"), "Invalid image filename",
			"The local file must have a bare filename without path separators.")

		return
	}
	resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("file_name"), name)...)
	if config.FileSha256.IsUnknown() {
		return
	}
	local, err := inspectLocalFile(config.FilePath.ValueString(), config.FileSha256.ValueString())
	if errors.Is(err, os.ErrNotExist) {
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("file_sha256"), config.FileSha256)...)
		if config.FileSha256.IsNull() {
			resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("file_sha256"), types.StringUnknown())...)
		}
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("file_size"), types.Int64Unknown())...)

		return
	}
	if err != nil {
		resp.Diagnostics.AddAttributeError(path.Root("file_path"), "Invalid local image file", err.Error())

		return
	}
	checksum := config.FileSha256.ValueString()
	if checksum == "" {
		checksum = local.Sha256
	}
	resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("file_sha256"), checksum)...)
	resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("file_size"), local.Size)...)
}

func applyLocalFile(ctx context.Context, client *sdk.APIClient, plan *ImageModel, state ImageModel) error {
	if plan.FilePath.IsNull() {
		plan.FileName = types.StringNull()
		plan.FileSha256 = types.StringNull()
		plan.FileSize = types.Int64Null()

		return nil
	}
	local, err := inspectLocalFile(plan.FilePath.ValueString(), plan.FileSha256.ValueString())
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(ctx, 2*time.Hour)
	defer cancel()
	name := filepath.Base(local.Path)
	files, _, err := listFiles(ctx, client, plan.Id.ValueInt64())
	if err != nil {
		return err
	}
	existing, present := findFile(files, name)
	changed := !state.FileSha256.IsNull() && !state.FileSha256.IsUnknown() &&
		!strings.EqualFold(state.FileSha256.ValueString(), local.Sha256)
	if changed || !present || existing.Size != local.Size {
		if present {
			if err := removeFile(ctx, client, plan.Id.ValueInt64(), name); err != nil {
				return err
			}
		}
		if err := uploadFile(ctx, client, plan.Id.ValueInt64(), name, local); err != nil {
			return err
		}
		files, err = awaitListed(ctx, client, plan.Id.ValueInt64(), name, local.Size, time.Minute)
		if err != nil {
			return err
		}
	}
	if err := removeOthers(ctx, client, plan.Id.ValueInt64(), files, name); err != nil {
		return err
	}
	plan.FileName = types.StringValue(name)
	plan.FileSize = types.Int64Value(local.Size)
	if plan.FileSha256.IsNull() || plan.FileSha256.IsUnknown() {
		plan.FileSha256 = types.StringValue(local.Sha256)
	}

	return nil
}
