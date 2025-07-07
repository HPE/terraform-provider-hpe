// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build generate && experimental

package tools

import (
	_ "github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs"
)

// Format Terraform code for use in documentation.
// If you do not have Terraform installed, you can remove the formatting command, but it is suggested
// to ensure the documentation is formatted properly.
//go:generate terraform fmt -recursive ../examples/

//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --website-source-dir templates-combined-temp --rendered-website-dir docs-experimental --examples-dir examples --provider-dir ..
