package image

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestUnitLocalUploadLifecycle(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	store := newFakeImageStore(nil)
	client := newFakeClient(t, store)
	first := []byte("first image")
	plan := ImageModel{
		Id:         types.Int64Value(42),
		FilePath:   types.StringValue(writeTempFile(t, "disk.qcow2", first)),
		FileSha256: types.StringValue(strings.ToUpper(sumOf(first))),
	}
	if err := applyLocalFile(ctx, client, &plan, ImageModel{}); err != nil {
		t.Fatal(err)
	}
	if plan.Id.ValueInt64() != 42 || plan.FileSha256.ValueString() != strings.ToUpper(sumOf(first)) ||
		plan.FileSize.ValueInt64() != int64(len(first)) {
		t.Fatal("upload changed identity or checksum casing")
	}
	state := plan
	store.calls = nil
	if err := applyLocalFile(ctx, client, &plan, state); err != nil {
		t.Fatal(err)
	}
	if len(store.calls) != 1 || !strings.HasPrefix(store.calls[0], "GET ") {
		t.Fatalf("unchanged file was uploaded: %v", store.calls)
	}
	second := []byte("other image")
	plan.FilePath = types.StringValue(writeTempFile(t, "disk.qcow2", second))
	plan.FileSha256 = types.StringValue(sumOf(second))
	if err := applyLocalFile(ctx, client, &plan, state); err != nil {
		t.Fatal(err)
	}
	if string(store.files["disk.qcow2"]) != string(second) || plan.Id.ValueInt64() != 42 {
		t.Fatal("same-size content change did not preserve image identity and replace bytes")
	}
	state = plan
	delete(store.files, "disk.qcow2")
	if err := applyLocalFile(ctx, client, &plan, state); err != nil {
		t.Fatal(err)
	}
	if string(store.files["disk.qcow2"]) != string(second) {
		t.Fatal("missing file was not repaired")
	}
	plan.FilePath = types.StringNull()
	if err := applyLocalFile(ctx, client, &plan, state); err != nil {
		t.Fatal(err)
	}
	if !plan.FileSha256.IsNull() || len(store.files) != 1 {
		t.Fatal("removing file_path must relinquish management without deleting content")
	}
}

func TestUnitLocalUploadAdoptionAndFailure(t *testing.T) {
	t.Parallel()
	for _, adopt := range []bool{false, true} {
		store := newFakeImageStore(map[string][]byte{"older.iso": []byte("keep until verified")})
		data := []byte("image bytes")
		if adopt {
			store.files["disk.iso"] = data
		} else {
			store.uploadStatus = 500
		}
		client := newFakeClient(t, store)
		plan := ImageModel{
			Id:         types.Int64Value(42),
			FilePath:   types.StringValue(writeTempFile(t, "disk.iso", data)),
			FileSha256: types.StringUnknown(),
		}
		err := applyLocalFile(context.Background(), client, &plan, ImageModel{})
		if adopt {
			if err != nil || len(store.files) != 1 || plan.FileSha256.ValueString() != sumOf(data) {
				t.Fatalf("adoption failed: %v", err)
			}
			for _, call := range store.calls {
				if strings.HasPrefix(call, "POST ") {
					t.Fatal("adoption uploaded existing bytes")
				}
			}
		} else if err == nil || store.files["older.iso"] == nil {
			t.Fatal("failed upload removed existing differently named content")
		}
	}
}
