// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SubModel struct {
	URL         types.String `tfsdk:"url"`
	Username    types.String `tfsdk:"username"`
	Password    types.String `tfsdk:"password"`
	AccessToken types.String `tfsdk:"access_token"`
}
