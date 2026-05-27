// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

//go:build sweep

package networkrouter

import (
	_ "github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/network_router_firewall_rule/sweep"
	_ "github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/network_router_nat/sweep"
	_ "github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/network_router_route/sweep"
	_ "github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/networkrouter/sweep"
	_ "github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/networkrouterbgpneighbor/sweep"
)
