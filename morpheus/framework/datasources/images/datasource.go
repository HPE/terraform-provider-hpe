// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package images implements a data source for images
package images

import (
	"context"
	"fmt"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	providererrors "github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
	"github.com/HPE/terraform-provider-hpe/utils/paging"
)

const summary = "read images data source"

// defaultFilterType diverges from the platform's own default of "user".
//
// The API applies eq('userUploaded', true) when no filterType is given, so a
// faithful default would return only images a user uploaded — hiding synced and
// system images, which are exactly the cloud-provided images most instances are
// provisioned from. A data source whose purpose is finding an image to
// provision from is not much use if it cannot see them.
const defaultFilterType = "All"

// Ensure the implementation satisfies the expected interfaces.
var _ datasource.DataSource = &DataSource{}

// NewDataSource is a helper function to simplify the provider implementation.
func NewDataSource() datasource.DataSource {
	return &DataSource{}
}

// DataSource is the data source implementation.
type DataSource struct {
	configure.DataSourceWithMorpheusConfigure
	datasource.DataSource
}

// Metadata returns the data source type name.
func (d *DataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_" + "images"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = ImagesDataSourceSchema(ctx)
}

// compiledFilter is a filter block with its values pre-compiled as regular
// expressions.
type compiledFilter struct {
	field string
	res   []*regexp.Regexp
}

// Read refreshes the Terraform state with the latest data.
//
// Filtering happens in two stages, and the split matters. The structured
// arguments become query parameters and decide how much is fetched; the filter
// blocks are regular expressions evaluated here, on whatever arrived. Sending a
// regular expression to the server would fail silently — matching there is a
// SQL `like`, so `^Test.*$` would be searched for as a literal and match
// nothing.
//
// That makes a configuration filtering only by block expensive: it fetches
// everything the arguments admit, which by default is the whole library.
// phraseFromFilters recovers some of that where a block offers a literal to
// narrow on.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config ImagesModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Compile before fetching, so an unusable expression fails immediately
	// rather than after walking the whole estate.
	filters := compileFilters(ctx, config.Filter, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	apiClient, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(summary, "could not create sdk client")

		return
	}

	params := buildServerSideParams(ctx, &config, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	// A configuration that filters only by block would otherwise fetch the
	// whole library. Where a name block starts with a literal, ask the server
	// for names containing it and let the block narrow what comes back.
	if params.phrase == nil {
		if phrase, ok := phraseFromFilters(filters); ok {
			params.phrase = &phrase
		}
	}

	// Every page is fetched before the filter blocks are applied. A regular
	// expression evaluated over a truncated fetch would report "no such image"
	// for one that exists, so completeness is a correctness requirement here.
	images, err := paging.Collect(
		ctx,
		func(ctx context.Context, offset, max int64) (
			[]sdk.ListVirtualImages200ResponseAllOfVirtualImagesInner, int64, error,
		) {
			apiReq := params.apply(
				apiClient.LibraryAPI.ListVirtualImages(ctx).
					Max(max).
					Offset(offset),
			)

			rs, hresp, err := apiReq.Execute()
			if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
				return nil, 0, fmt.Errorf("LIST failed for images: %s",
					providererrors.ErrMsg(err, hresp))
			}

			var total int64
			if rs.Meta != nil && rs.Meta.Total != nil {
				total = *rs.Meta.Total
			}

			return rs.VirtualImages, total, nil
		},
	)
	if err != nil {
		resp.Diagnostics.AddError(summary, err.Error())

		return
	}

	selected := selectImages(images, filters)

	objs := make([]attr.Value, 0, len(selected))

	for _, img := range selected {
		v, diags := imageToValue(ctx, img)
		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		objs = append(objs, v)
	}

	// The generated schema declares the element type of images as the custom
	// ImagesType, so both the element type used here and every element must be
	// that type. A bare types.Object element fails the set's element type check.
	setVal, diags := types.SetValue(ImagesValue{}.Type(ctx), objs)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	config.Images = setVal

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// phraseFromFilters derives a server-side phrase from a name filter block,
// where one can be derived soundly.
//
// A filter block is evaluated here, after everything matching the server-side
// arguments has been fetched, so a configuration filtering only by block reads
// the whole library. Where the block's expression begins with a literal, that
// literal must appear in every name it can match, and the API can be asked to
// return only names containing it.
//
// This is sound in one direction only, which is the direction that matters.
// phrase is ilike('name', '%phrase%'), a case-insensitive substring; a regular
// expression's literal prefix is case-sensitive. The server therefore returns a
// superset of what the block will accept, and the block still runs afterwards.
// Narrowing further server-side would risk discarding a row the block would
// have matched.
//
// Only a single-valued block qualifies. Values within a block are ORed, and one
// phrase cannot express that union unless the values share a prefix — not worth
// the subtlety for the gain.
func phraseFromFilters(filters []compiledFilter) (string, bool) {
	for _, f := range filters {
		if f.field != "name" || len(f.res) != 1 {
			continue
		}

		// LiteralPrefix reports a string every match must begin with. Whether
		// it is the complete match does not matter: a prefix is contained in
		// the match either way, which is all a substring search needs.
		prefix, _ := f.res[0].LiteralPrefix()
		if prefix != "" {
			return prefix, true
		}
	}

	return "", false
}

// serverSideParams is the set of query parameters a configuration maps to.
//
// It is computed separately from being applied so the mapping can be asserted
// directly: the generated request builder keeps its fields unexported, so a
// request that has been built cannot be inspected.
type serverSideParams struct {
	filterType         string
	imageTypes         []string
	descriptions       []string
	imageIDs           []int64
	phrase             *string
	systemImage        *bool
	includeSystemImage *bool
	labels             *string
	allLabels          *string
	sort               *string
	direction          *string
}

// buildServerSideParams maps the configured arguments onto query parameters.
//
// Only structured values are included. A filter block's regular expression must
// never reach the server: matching there is a SQL `like`, so forwarding
// `^Test.*$` would search for that literal and match nothing.
func buildServerSideParams(
	ctx context.Context,
	config *ImagesModel,
	diags *diag.Diagnostics,
) serverSideParams {
	p := serverSideParams{
		filterType:   defaultFilterType,
		imageTypes:   stringsFromSet(ctx, config.ImageType, diags),
		descriptions: stringsFromSet(ctx, config.Description, diags),
		imageIDs:     int64sFromSet(ctx, config.ImageId, diags),
	}

	if !config.FilterType.IsNull() && !config.FilterType.IsUnknown() {
		p.filterType = config.FilterType.ValueString()
	}

	if !config.Phrase.IsNull() && !config.Phrase.IsUnknown() {
		p.phrase = config.Phrase.ValueStringPointer()
	}

	if !config.SystemImage.IsNull() && !config.SystemImage.IsUnknown() {
		p.systemImage = config.SystemImage.ValueBoolPointer()
	}

	if !config.IncludeSystemImage.IsNull() && !config.IncludeSystemImage.IsUnknown() {
		p.includeSystemImage = config.IncludeSystemImage.ValueBoolPointer()
	}

	if !config.Labels.IsNull() && !config.Labels.IsUnknown() {
		p.labels = config.Labels.ValueStringPointer()
	}

	if !config.AllLabels.IsNull() && !config.AllLabels.IsUnknown() {
		p.allLabels = config.AllLabels.ValueStringPointer()
	}

	if !config.Sort.IsNull() && !config.Sort.IsUnknown() {
		p.sort = config.Sort.ValueStringPointer()
	}

	if !config.Direction.IsNull() && !config.Direction.IsUnknown() {
		p.direction = config.Direction.ValueStringPointer()
	}

	return p
}

// apply sets the parameters on a request.
func (p serverSideParams) apply(
	apiReq sdk.ApiListVirtualImagesRequest,
) sdk.ApiListVirtualImagesRequest {
	apiReq = apiReq.FilterType(p.filterType)

	if len(p.imageTypes) > 0 {
		apiReq = apiReq.ImageType(p.imageTypes)
	}

	if len(p.descriptions) > 0 {
		apiReq = apiReq.Description(p.descriptions)
	}

	if len(p.imageIDs) > 0 {
		apiReq = apiReq.Id(p.imageIDs)
	}

	if p.phrase != nil {
		apiReq = apiReq.Phrase(*p.phrase)
	}

	if p.systemImage != nil {
		apiReq = apiReq.SystemImage(*p.systemImage)
	}

	if p.includeSystemImage != nil {
		apiReq = apiReq.IncludeSystemImage(*p.includeSystemImage)
	}

	if p.labels != nil {
		apiReq = apiReq.Labels(*p.labels)
	}

	if p.allLabels != nil {
		apiReq = apiReq.AllLabels(*p.allLabels)
	}

	if p.sort != nil {
		apiReq = apiReq.Sort(*p.sort)
	}

	if p.direction != nil {
		apiReq = apiReq.Direction(*p.direction)
	}

	return apiReq
}

func stringsFromSet(
	ctx context.Context,
	set types.Set,
	diags *diag.Diagnostics,
) []string {
	if set.IsNull() || set.IsUnknown() {
		return nil
	}

	var out []string
	diags.Append(set.ElementsAs(ctx, &out, false)...)

	return out
}

func int64sFromSet(
	ctx context.Context,
	set types.Set,
	diags *diag.Diagnostics,
) []int64 {
	if set.IsNull() || set.IsUnknown() {
		return nil
	}

	var out []int64
	diags.Append(set.ElementsAs(ctx, &out, false)...)

	return out
}

func compileFilters(
	ctx context.Context,
	filterSet types.Set,
	diags *diag.Diagnostics,
) []compiledFilter {
	if filterSet.IsNull() || filterSet.IsUnknown() {
		return nil
	}

	var filterBlocks []FilterValue

	diags.Append(filterSet.ElementsAs(ctx, &filterBlocks, false)...)

	if diags.HasError() {
		return nil
	}

	compiled := make([]compiledFilter, 0, len(filterBlocks))

	for _, b := range filterBlocks {
		field := b.Name.ValueString()

		var values []string

		diags.Append(b.Values.ElementsAs(ctx, &values, false)...)

		if diags.HasError() {
			return nil
		}

		res := make([]*regexp.Regexp, 0, len(values))

		for _, v := range values {
			re, err := regexp.Compile(v)
			if err != nil {
				diags.AddError(summary, fmt.Sprintf(
					"invalid regular expression %q for filter %q: %s", v, field, err))

				return nil
			}

			res = append(res, re)
		}

		compiled = append(compiled, compiledFilter{field: field, res: res})
	}

	return compiled
}

// selectImages de-duplicates the fetched images by id and keeps those the
// filter blocks accept, preserving the order the endpoint served them.
//
// De-duplication is needed because offset paging over a non-unique sort key can
// serve the same record twice: the API orders by name, names are not unique
// across clouds, and rows can move between pages as ties are broken differently
// from one query to the next.
//
// An image with no id cannot be de-duplicated and is kept as it comes. That
// should not happen — every image the API returns has one — but dropping
// records on a missing field would be a worse failure than passing a duplicate
// through.
func selectImages(
	images []sdk.ListVirtualImages200ResponseAllOfVirtualImagesInner,
	filters []compiledFilter,
) []*sdk.ListVirtualImages200ResponseAllOfVirtualImagesInner {
	out := make([]*sdk.ListVirtualImages200ResponseAllOfVirtualImagesInner, 0, len(images))
	seen := make(map[int64]struct{}, len(images))

	for i := range images {
		img := &images[i]

		if img.Id != nil {
			if _, dup := seen[*img.Id]; dup {
				continue
			}

			seen[*img.Id] = struct{}{}
		}

		if !matchesFilters(img, filters) {
			continue
		}

		out = append(out, img)
	}

	return out
}

// matchesFilters reports whether an image satisfies every filter block. Within
// a block any value may match; across blocks all must.
func matchesFilters(
	img *sdk.ListVirtualImages200ResponseAllOfVirtualImagesInner,
	filters []compiledFilter,
) bool {
	for _, f := range filters {
		val, ok := fieldValue(img, f.field)
		if !ok {
			return false
		}

		matched := false

		for _, re := range f.res {
			if re.MatchString(val) {
				matched = true

				break
			}
		}

		if !matched {
			return false
		}
	}

	return true
}

// fieldValue resolves a filter field name to the image's value for it. The
// second result is false when the image has no value, which cannot match any
// expression.
func fieldValue(
	img *sdk.ListVirtualImages200ResponseAllOfVirtualImagesInner,
	field string,
) (string, bool) {
	switch field {
	case "name":
		if img.Name == nil {
			return "", false
		}

		return *img.Name, true
	case "description":
		if v := img.Description.Get(); v != nil {
			return *v, true
		}

		return "", false
	case "image_type":
		if img.ImageType == nil {
			return "", false
		}

		return *img.ImageType, true
	case "status":
		if img.Status == nil {
			return "", false
		}

		return *img.Status, true
	case "visibility":
		if img.Visibility == nil {
			return "", false
		}

		return *img.Visibility, true
	default:
		// Unreachable: the schema restricts the field name to the cases above.
		return "", false
	}
}

// imageToValue converts one API image into the generated element type.
//
// The fields are mapped directly rather than re-decoded through the
// single-image type. Round-tripping via JSON is tempting, since the two
// endpoints describe the same image, but the generated config field is a oneOf
// whose MarshalJSON fails outright when no variant is set — which is the normal
// case for every image that is not an Azure reference. The element shape is
// instead kept in step with hpe_morpheus_image by the spec: both are generated
// from the same element definition, and their field sets are asserted equal.
func imageToValue(
	ctx context.Context,
	img *sdk.ListVirtualImages200ResponseAllOfVirtualImagesInner,
) (ImagesValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	labels := convert.StrSliceToSet(img.Labels)

	tags, d := convert.ToSetType(
		ctx,
		img.Tags,
		func(in sdk.ListVirtualImages200ResponseAllOfVirtualImagesInnerTagsInner) TagsValue {
			return TagsValue{
				Name:  convert.StrToType(in.Name),
				Value: convert.StrToType(in.Value),
				state: attr.ValueStateKnown,
			}
		},
	)
	diags.Append(d...)

	tenants, d := convert.ToSetType(
		ctx,
		img.Accounts,
		func(in sdk.ListVirtualImages200ResponseAllOfVirtualImagesInnerAccountsInner) TenantsValue {
			return TenantsValue{
				Name:  convert.StrToType(in.Name),
				Id:    convert.Int64ToType(in.Id),
				state: attr.ValueStateKnown,
			}
		},
	)
	diags.Append(d...)

	v := ImagesValue{
		AutoJoinDomain:     convert.BoolToType(img.IsAutoJoinDomain),
		CloudInit:          convert.BoolToType(img.IsCloudInit),
		ConfigAzure:        configAzureValue(ctx, img),
		ConsoleKeymap:      convert.StrToType(img.ConsoleKeymap.Get()),
		Description:        convert.StrToType(img.Description.Get()),
		ExternalId:         convert.StrToType(img.ExternalId.Get()),
		FipsEnabled:        convert.BoolToType(img.FipsEnabled),
		ForceCustomization: convert.BoolToType(img.IsForceCustomization),
		Id:                 convert.Int64ToType(img.Id),
		ImageType:          convert.StrToType(img.ImageType),
		InstallAgent:       convert.BoolToType(img.InstallAgent),
		Labels:             labels,
		MinDisk:            bytesToGB(img.MinDisk.Get()),
		MinRam:             bytesToGB(img.MinRam.Get()),
		Name:               convert.StrToType(img.Name),
		OsTypeId:           convert.Int64ToType(osTypeID(img)),
		OwnerId:            convert.Int64ToType(img.OwnerId),
		RawSize:            convert.Int64ToType(img.RawSize.Get()),
		SshKey:             convert.StrToType(img.SshKey.Get()),
		SshUsername:        convert.StrToType(img.SshUsername.Get()),
		Status:             convert.StrToType(img.Status),
		StorageProviderId:  convert.Int64ToType(storageProviderID(img)),
		Sysprep:            convert.BoolToType(img.IsSysprep),
		SystemImage:        convert.BoolToType(img.SystemImage),
		Tags:               tags,
		Tenants:            tenants,
		TrialVersion:       convert.BoolToType(img.TrialVersion),
		Uefi:               convert.BoolToType(img.Uefi.Get()),
		UserData:           convert.StrToType(img.UserData.Get()),
		UserDefined:        convert.BoolToType(img.UserDefined),
		UserUploaded:       convert.BoolToType(img.UserUploaded),
		VirtioSupported:    convert.BoolToType(img.VirtioSupported),
		Visibility:         convert.StrToType(img.Visibility),
		VmToolsInstalled:   convert.BoolToType(img.VmToolsInstalled),
		state:              attr.ValueStateKnown,
	}

	return v, diags
}

// bytesToGB matches the singular data source, which reports min_disk and
// min_ram in gigabytes while the API returns bytes.
func bytesToGB(v *int64) basetypes.Int64Value {
	if v == nil {
		return basetypes.NewInt64Null()
	}

	gb := *v / 1024 / 1024 / 1024

	return convert.Int64ToType(&gb)
}

// configAzureValue is populated only for azure-reference images. Every other
// image reports it null rather than as a set of empty strings.
func configAzureValue(
	ctx context.Context,
	img *sdk.ListVirtualImages200ResponseAllOfVirtualImagesInner,
) basetypes.ObjectValue {
	nullValue := types.ObjectNull(ConfigAzureValue{}.AttributeTypes(ctx))

	if img.ImageType == nil || *img.ImageType != "azure-reference" {
		return nullValue
	}

	if img.Config == nil || img.Config.AzureReferenceVirtualImageConfiguration == nil {
		return nullValue
	}

	azure := img.Config.AzureReferenceVirtualImageConfiguration

	cfg := ConfigAzureValue{
		Publisher: convert.StrToType(&azure.Publisher),
		Offer:     convert.StrToType(&azure.Offer),
		Version:   convert.StrToType(&azure.Version),
		Sku:       convert.StrToType(&azure.Sku),
		state:     attr.ValueStateKnown,
	}

	obj, diags := cfg.ToObjectValue(ctx)
	if diags.HasError() {
		return nullValue
	}

	return obj
}

func osTypeID(img *sdk.ListVirtualImages200ResponseAllOfVirtualImagesInner) *int64 {
	if img.OsType == nil {
		return nil
	}

	return img.OsType.Id
}

func storageProviderID(
	img *sdk.ListVirtualImages200ResponseAllOfVirtualImagesInner,
) *int64 {
	if img.StorageProvider == nil {
		return nil
	}

	return img.StorageProvider.Id
}
