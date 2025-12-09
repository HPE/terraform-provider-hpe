// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package network

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_ip_pool_ipv4/resource.tf morpheus_ip_pool_ipv4_resource.tf.tmpl Name "\"Terraform Example IPv4 IP pool\"" StartingAddress1 "\"192.168.1.1\"" EndingAddress1 "\"192.168.1.10\"" StartingAddress2 "\"10.0.0.1\"" EndingAddress2 "\"10.0.0.10\""
