// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package hosttypes answers questions about compute server types that the
// servers API does not: /api/servers renders a server's computeServerType
// with only id, code, name, managed and externalDelete, so whether a host is
// bare metal has to come from /api/server-types.
package hosttypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/paging"
)

// BareMetalCodes returns the codes of the server types that are bare-metal
// hosts (computeServerType.bareMetalHost), from every page of
// /api/server-types.
func BareMetalCodes(ctx context.Context, client *sdk.APIClient) (map[string]bool, error) {
	types, err := paging.Collect(ctx, func(
		ctx context.Context,
		offset, max int64,
	) ([]sdk.ListHostTypes200ResponseAllOfServerTypesInner, int64, error) {
		rs, hresp, err := client.HostsAPI.ListHostTypes(ctx).Max(max).Offset(offset).Execute()
		if err != nil || rs == nil {
			return nil, 0, fmt.Errorf("ListHostTypes failed: %s", errfmt.ErrMsg(err, hresp))
		}

		total := int64(0)
		if rs.Meta != nil && rs.Meta.Total != nil {
			total = *rs.Meta.Total
		}

		return rs.ServerTypes, total, nil
	})
	if err != nil {
		return nil, err
	}

	codes := map[string]bool{}
	for _, t := range types {
		if t.Code != nil && t.BareMetalHost != nil && *t.BareMetalHost {
			codes[*t.Code] = true
		}
	}

	return codes, nil
}

// IsBareMetal reports whether the server type code is a bare-metal host type.
func IsBareMetal(codes map[string]bool, code *string) bool {
	return code != nil && codes[*code]
}

// BareMetalValue is the bare_metal attribute value for a server type code:
// null when the codes are not known (the server types list could not be read),
// so "not known" is never reported as false.
func BareMetalValue(codes map[string]bool, code *string) types.Bool {
	if codes == nil {
		return types.BoolNull()
	}

	return types.BoolValue(IsBareMetal(codes, code))
}
