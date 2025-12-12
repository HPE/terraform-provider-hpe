// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package integration

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_integration_ansible_tower/resource.tf morpheus_integration_ansible_tower_resource.tf.tmpl Name "tfexample ansible tower integration" Enabled true Url "https://ansibletower01.morpheusdata.com" Username admin Password password123
