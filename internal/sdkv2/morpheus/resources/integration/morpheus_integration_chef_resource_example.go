// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package integration

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_integration_chef/resource.tf morpheus_integration_chef_resource.tf.tmpl Name "tfexample chef integration" Enabled true Url "https://chef.morpheusdata.com" Version "15.9.38" WindowsVersion "15.9.38" WindowsMsiInstallUrl "https://packages.chef.io" Organization "morpheus" Username "admin" PrivateKey "EXAMPLEPRIVATEKEY" OrganizationValidatorKey "EXAMPLEPRIVATEKEY"
