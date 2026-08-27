// (C) Copyright 2026 Hewlett Packard Enterprise Development LP
package provider

// Version is the provider version. It is set at build time via -ldflags:
//
//	go build -ldflags="-X 'github.com/HPE/terraform-provider-opsramp.Version=1.0.0'"
//
// The default "dev" value is used when building without ldflags.
var Version = "dev"
