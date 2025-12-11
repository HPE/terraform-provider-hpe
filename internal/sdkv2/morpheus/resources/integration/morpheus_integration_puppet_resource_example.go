// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package integration

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_integration_puppet/resource.tf morpheus_integration_puppet_resource.tf.tmpl Name "tfexample puppet integration" Enabled true PuppetMasterHostname peserver01.morpheusdata.com AllowImmediateExecution true PuppetMasterSshUsername root PuppetMasterSshPassword password123
