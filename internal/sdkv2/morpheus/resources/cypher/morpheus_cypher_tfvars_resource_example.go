// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package cypher

//go:generate sh -c "go run ../../../../../cmd/render -out examples/resources/morpheus_cypher_tfvars/resource.tf hpe_morpheus_cypher_tfvars_resource.tf.tmpl Key securetfvars Value 'account=12345\npassword=supersecure' Ttl 86400"
