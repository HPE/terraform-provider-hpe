// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package host

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

// notCreatableHints explains, for server types a practitioner is likely to try,
// which resource actually creates them. The host resource can only add the
// types Morpheus marks creatable (GET /api/server-types?creatable=true); the
// others come into being through the object that owns them.
var notCreatableHints = map[string]string{
	"mvm":                 "an HVM hypervisor is created by hpe_morpheus_cluster as a member of its cluster; to add one to an existing cluster use the cluster's Add Worker (no resource yet)",
	"morpheusVm":          "an HVM virtual machine is created by hpe_morpheus_instance",
	"morpheusWindowsVm":   "an HVM virtual machine is created by hpe_morpheus_instance",
	"vmwareVm":            "a VMware virtual machine is created by hpe_morpheus_instance",
	"vmwareWindowsVm":     "a VMware virtual machine is created by hpe_morpheus_instance",
	"morpheusDockerLinux": "a Docker host is created through its cluster (hpe_morpheus_cluster)",
	"morpheusKubeMaster":  "a Kubernetes master is created through its cluster (hpe_morpheus_cluster)",
	"morpheusKubeWorker":  "a Kubernetes worker is created through its cluster (hpe_morpheus_cluster)",
}

// checkCreatable refuses a compute_server_type that Morpheus does not let a
// user add as a standalone host, before the create is attempted. The API's own
// answer for such a type is an unhelpful generic error; this names the type,
// says which resource owns it where known, and lists what is creatable.
// Failing to read the type list is not fatal: the create proceeds and the API
// decides.
func checkCreatable(ctx context.Context, client *sdk.APIClient, typeCode string) diag.Diagnostics {
	var diags diag.Diagnostics

	// Ask for the type by code, so the answer does not depend on a page size.
	byCode, httpResp, err := client.HostsAPI.ListHostTypes(ctx).Creatable(true).Code(typeCode).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil || byCode == nil {
		return diags // fail open
	}

	for _, t := range byCode.ServerTypes {
		if t.Code != nil && *t.Code == typeCode {
			return diags
		}
	}

	// Not creatable. The list of what is, for the message; best effort.
	var codes []string

	all, httpResp, err := client.HostsAPI.ListHostTypes(ctx).Creatable(true).Max(500).Execute()
	if errfmt.CheckResponse(err, httpResp) == nil && all != nil {
		for _, t := range all.ServerTypes {
			if t.Code != nil {
				codes = append(codes, *t.Code)
			}
		}
	}
	sort.Strings(codes)

	detail := fmt.Sprintf("Morpheus does not allow a host of type %q to be added on its own", typeCode)
	if hint, ok := notCreatableHints[typeCode]; ok {
		detail += ": " + hint
	}
	detail += fmt.Sprintf(". The hpe_morpheus_host resource can create the types the appliance marks creatable "+
		"(GET /api/server-types?creatable=true); on this appliance: %s. Existing hosts of any type can be read "+
		"with the hpe_morpheus_host and hpe_morpheus_hosts data sources.", strings.Join(codes, ", "))
	diags.AddAttributeError(path.Root("compute_server_type"), "Server type cannot be added as a host", detail)

	return diags
}
