// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package ostype

import (
	"context"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"

	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

func mapImage(img sdk.GetOsType200ResponseOsTypeImagesInner) ImagesValue {
	ctx := context.Background()

	v, _ := NewImagesValue(
		ImagesValue{}.AttributeTypes(ctx),
		map[string]attr.Value{
			"account":            convert.Int64ToType(img.Account.Get()),
			"compute_zone_type":  convert.Int64ToType(img.ComputeZoneType.Get()),
			"id":                 convert.Int64ToType(img.Id),
			"provision_type":     convert.Int64ToType(img.ProvisionType.Get()),
			"virtual_image_id":   convert.Int64ToType(img.VirtualImageId),
			"virtual_image_name": convert.StrToType(img.VirtualImageName),
			"zone":               convert.Int64ToType(img.Zone.Get()),
		},
	)

	return v
}
