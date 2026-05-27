// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

//go:build sweep

package cluster

import (
	_ "github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/cluster/sweep"
	_ "github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/cluster_affinity_group/sweep"
	_ "github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/cluster_datastore/sweep"
	_ "github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/cluster_namespace/sweep"
)
