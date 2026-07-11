// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0
package provider

// Version is the provider version. It is set at build time via -ldflags:
//
//	go build -ldflags="-X 'github.com/HPE/terraform-provider-opsramp.Version=1.0.0'"
//
// The default "dev" value is used when building without ldflags.
var Version = "dev"
