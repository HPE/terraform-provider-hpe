// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package model

type SubModel struct {
	URL         string `tfsdk:"url"`
	Username    string `tfsdk:"username"`
	Password    string `tfsdk:"password"`
	AccessToken string `tfsdk:"access_token"`
}
