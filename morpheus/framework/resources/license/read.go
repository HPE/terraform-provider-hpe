// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package license

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var prior LicenseModel
	resp.Diagnostics.Append(req.State.Get(ctx, &prior)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The key is write-only, so match the installed licence by the digest kept
	// in private state (recorded from the API's hash by a prior apply or
	// import); fall back to the key id when there is no digest yet.
	stored, ok := storedDigest(ctx, req.Private, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	hash := ""
	if ok {
		hash = stored
	}

	state, installed, found := getLicenseAsState(ctx, client, hash, prior.Id.ValueString(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// The configured licence was uninstalled or replaced out of band. Remove it
	// from state so the next plan recreates it (MORPH-18584).
	if !found {
		tflog.Info(ctx, fmt.Sprintf("license %s not among the installed licences, removing from state",
			prior.Id.ValueString()))
		resp.State.RemoveResource(ctx)

		return
	}

	// Record the digest from the API's hash. This is how an imported or
	// state-upgraded resource acquires its digest without ever seeing the key.
	if installed.Hash != "" {
		resp.Diagnostics.Append(storeDigest(ctx, resp.Private, installed.Hash)...)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// installedLicense is the subset of an installed licence's fields the resource
// needs. GetLicense200Response reports the same shape under two Go types --
// GetLicense200ResponseInstalledLicensesInner and GetLicense200ResponseLicense
// -- so this small adapter lets findInstalledLicense return either.
type installedLicense struct {
	KeyId        string
	Hash         string
	ProductTier  *string
	AccountName  *string
	StartDate    *time.Time
	EndDate      *time.Time
	MaxInstances *int64
	MaxMemory    *int64
	MaxStorage   *int64
	HardLimit    *bool
	FreeTrial    *bool
	MultiTenant  *bool
	SupportLevel *string
}

func fromInstalledInner(l *sdk.GetLicense200ResponseInstalledLicensesInner) installedLicense {
	return installedLicense{
		KeyId:        strOrEmpty(l.KeyId),
		Hash:         strOrEmpty(l.Hash),
		ProductTier:  l.ProductTier,
		AccountName:  l.AccountName,
		StartDate:    l.StartDate,
		EndDate:      l.EndDate,
		MaxInstances: l.MaxInstances,
		MaxMemory:    l.MaxMemory,
		MaxStorage:   l.MaxStorage,
		HardLimit:    l.HardLimit,
		FreeTrial:    l.FreeTrial,
		MultiTenant:  l.MultiTenant,
		SupportLevel: l.SupportLevel,
	}
}

func fromMergedLicense(l *sdk.GetLicense200ResponseLicense) installedLicense {
	return installedLicense{
		KeyId:        strOrEmpty(l.KeyId),
		Hash:         strOrEmpty(l.Hash),
		ProductTier:  l.ProductTier,
		AccountName:  l.AccountName,
		StartDate:    l.StartDate,
		EndDate:      l.EndDate,
		MaxInstances: l.MaxInstances,
		MaxMemory:    l.MaxMemory,
		MaxStorage:   l.MaxStorage,
		HardLimit:    l.HardLimit,
		FreeTrial:    l.FreeTrial,
		MultiTenant:  l.MultiTenant,
		SupportLevel: l.SupportLevel,
	}
}

func strOrEmpty(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}

// findInstalledLicense returns the installed licence matching hash (when hash
// is set) or else keyID, searching the installed-licences list first and
// falling back to the merged License object only if the list had no match. It
// returns nil when nothing matches or there is nothing to match on.
func findInstalledLicense(result *sdk.GetLicense200Response, hash, keyID string) *installedLicense {
	if result == nil {
		return nil
	}

	for i := range result.InstalledLicenses {
		l := fromInstalledInner(&result.InstalledLicenses[i])
		if matchesLicense(l, hash, keyID) {
			return &l
		}
	}

	if result.License != nil {
		l := fromMergedLicense(result.License)
		if matchesLicense(l, hash, keyID) {
			return &l
		}
	}

	return nil
}

func matchesLicense(l installedLicense, hash, keyID string) bool {
	if hash != "" {
		return strings.EqualFold(l.Hash, hash)
	}

	return keyID != "" && l.KeyId == keyID
}

// licenseToState maps an installed licence into the resource model. The key is
// write-only and is left null (the framework nullifies it anyway).
func licenseToState(l installedLicense) LicenseModel {
	return LicenseModel{
		Id:           convert.StrToType(&l.KeyId),
		Key:          types.StringNull(),
		ProductTier:  convert.StrToType(l.ProductTier),
		AccountName:  convert.StrToType(l.AccountName),
		StartDate:    convert.TimeToType(l.StartDate),
		EndDate:      convert.TimeToType(l.EndDate),
		MaxInstances: convert.Int64ToType(l.MaxInstances),
		MaxMemory:    convert.Int64ToType(l.MaxMemory),
		MaxStorage:   convert.Int64ToType(l.MaxStorage),
		HardLimit:    convert.BoolToType(l.HardLimit),
		FreeTrial:    convert.BoolToType(l.FreeTrial),
		MultiTenant:  convert.BoolToType(l.MultiTenant),
		SupportLevel: convert.StrToType(l.SupportLevel),
	}
}

// provisionalState is the state Create and Update leave behind when the key
// was installed but could not be read back: the key id, which identifies the
// installed licence to Read and Delete, and every detail null. The following
// refresh completes it from the appliance.
func provisionalState(keyID string) *LicenseModel {
	return &LicenseModel{
		Id:           types.StringValue(keyID),
		Key:          types.StringNull(),
		ProductTier:  types.StringNull(),
		AccountName:  types.StringNull(),
		StartDate:    types.StringNull(),
		EndDate:      types.StringNull(),
		MaxInstances: types.Int64Null(),
		MaxMemory:    types.Int64Null(),
		MaxStorage:   types.Int64Null(),
		HardLimit:    types.BoolNull(),
		FreeTrial:    types.BoolNull(),
		MultiTenant:  types.BoolNull(),
		SupportLevel: types.StringNull(),
	}
}

// getLicenseAsState performs a GET /api/license and matches the configured
// licence by hash (preferred) or key id. found is false when the licence is
// not installed; the returned installedLicense is valid only when found.
// Diagnostics are appended to diags (a pointer argument, not a result, so this
// helper is not mistaken for a Read not-found helper by the read-rules test).
func getLicenseAsState(
	ctx context.Context,
	client *sdk.APIClient,
	hash, keyID string,
	diags *diag.Diagnostics,
) (LicenseModel, installedLicense, bool) {
	result, httpResp, err := client.LicenseAPI.GetLicense(ctx).Execute()
	// GET /api/license never 404s, but handle it defensively as not-found.
	if errfmt.IsNotFound(httpResp) {
		return LicenseModel{}, installedLicense{}, false
	}
	if err != nil || httpResp.StatusCode != http.StatusOK {
		diags.AddError(readOperation, errfmt.ErrMsg(err, httpResp))

		return LicenseModel{}, installedLicense{}, false
	}

	installed := findInstalledLicense(result, hash, keyID)
	if installed == nil {
		return LicenseModel{}, installedLicense{}, false
	}

	return licenseToState(*installed), *installed, true
}
