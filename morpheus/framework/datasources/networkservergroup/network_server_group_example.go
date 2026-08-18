// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkservergroup

//go:generate ../../../../bin/render -out examples/data-sources/morpheus_network_server_group/example-name.tf example-name.tf.tmpl Name "Example group"
//go:generate ../../../../bin/render -out examples/data-sources/morpheus_network_server_group/example-name-with-server.tf example-name-with-server.tf.tmpl Name "Example group" NetworkServerId 99
