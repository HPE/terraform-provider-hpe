// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkpoolip

//go:generate ../../../../bin/render -out examples/data-sources/morpheus_network_pool_ip/example-id.tf example-id.tf.tmpl NetworkPoolId 99 Id 1
//go:generate ../../../../bin/render -out examples/data-sources/morpheus_network_pool_ip/example-ip-address.tf example-ip-address.tf.tmpl NetworkPoolId 99 IpAddress "10.0.0.5"
