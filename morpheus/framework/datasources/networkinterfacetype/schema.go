// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Hand-written schema for the hpe_morpheus_network_interface_type data source.
//
// This lookup has no singular REST endpoint to generate from: a network
// interface (NIC) type is only discoverable, for a given cloud and provision
// type, via the options endpoint /api/options/zoneNetworkOptions, which returns
// the available NIC types under the generic name "networkTypes". The schema is
// therefore maintained here directly rather than produced by the code-spec
// generator.
//
// It is the hpe equivalent of hpegl_vmaas_network_interface and supplies the
// value consumed by network_type_id (hpe_morpheus_instance) and
// network_interface_type_id (hpe_morpheus_instance_clone / hpe_morpheus_cluster).

package networkinterfacetype

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// defaultProvisionTypeCode mirrors hpegl_vmaas_network_interface, which resolves
// NIC types against the VMware provision type. It is exposed as an optional
// argument so non-VMware clouds remain reachable.
const defaultProvisionTypeCode = "vmware"

func NetworkInterfaceTypeDataSourceSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Retrieves the id of a Morpheus network interface (NIC) type available " +
			"for a given cloud and provision type, for use with network_type_id on " +
			"hpe_morpheus_instance or network_interface_type_id on hpe_morpheus_instance_clone.",
		MarkdownDescription: "Retrieves the id of a Morpheus network interface (NIC) type available " +
			"for a given cloud and provision type, for use with `network_type_id` on " +
			"`hpe_morpheus_instance` or `network_interface_type_id` on `hpe_morpheus_instance_clone`.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
				Description: "The name of the network interface type to look up " +
					"(for example \"E1000\" or \"VMXNET 3\").",
				MarkdownDescription: "The name of the network interface type to look up " +
					"(for example `E1000` or `VMXNET 3`).",
			},
			"cloud_id": schema.Int64Attribute{
				Required:            true,
				Description:         "The id of the cloud (zone) the network interface type is available in.",
				MarkdownDescription: "The id of the cloud (zone) the network interface type is available in.",
			},
			"provision_type_code": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Description: "The provision type code the network interface type belongs to. " +
					"Defaults to \"vmware\".",
				MarkdownDescription: "The provision type code the network interface type belongs to. " +
					"Defaults to `vmware`.",
			},
			"id": schema.Int64Attribute{
				Computed:            true,
				Description:         "The id of the network interface type.",
				MarkdownDescription: "The id of the network interface type.",
			},
			"code": schema.StringAttribute{
				Computed:            true,
				Description:         "The code of the network interface type.",
				MarkdownDescription: "The code of the network interface type.",
			},
		},
	}
}

type NetworkInterfaceTypeModel struct {
	Name              types.String `tfsdk:"name"`
	CloudId           types.Int64  `tfsdk:"cloud_id"`
	ProvisionTypeCode types.String `tfsdk:"provision_type_code"`
	Id                types.Int64  `tfsdk:"id"`
	Code              types.String `tfsdk:"code"`
}
