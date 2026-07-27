// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

//go:build sweep

package networkrouter

import (
	_ "github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/networkrouter/sweep"
	_ "github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/networkrouterbgpneighbor/sweep"
	_ "github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/networkrouterfirewallrule/sweep"
	_ "github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/networkrouterfirewallrulegroup/sweep"
	_ "github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/networkrouternat/sweep"
	_ "github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/networkrouterroute/sweep"
)
