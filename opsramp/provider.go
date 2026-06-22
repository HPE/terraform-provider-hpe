package opsramp

import (
// "context"
//
// "github.com/HPE/terraform-provider-hpe/provider/subprovider"
// opsrampprovider "github.com/HPE/terraform-provider-opsramp/src/provider"
// "github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
// "github.com/hashicorp/terraform-plugin-framework/list/schema"
// "github.com/hashicorp/terraform-plugin-framework/path"
// "github.com/hashicorp/terraform-plugin-framework/provider"
// "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// var _ provider.Provider = &OpsRampSubProvider{}
// var _ subprovider.SubProvider = &OpsRampSubProvider{}
//
// const subProviderName = "opsramp"
//
// type OpsRampSubProvider struct {
// 	p provider.Provider
// }
//
// func New() *OpsRampSubProvider {
// 	return &OpsRampSubProvider{
// 		p: opsrampprovider.New("foo"),
// 	}
// }
//
// func (p *OpsRampSubProvider) Configure(_ context.Context, f func(any)) (any, error) {}
//
// func (OpsRampSubProvider) GetName(_ context.Context) string {
// 	return subProviderName
// }
//
// func (OpsRampSubProvider) GetSchema(_ context.Context) map[string]schema.Attribute {
// 	parentBlock := path.MatchRelative().AtParent()
//
// 	return map[string]schema.Attribute{
// 		"url": schema.StringAttribute{
// 			Description: "Morpheus instance URL",
// 			Required:    true,
// 			Validators: []validator.String{
// 				stringvalidator.Any(
// 					stringvalidator.AlsoRequires(parentBlock.AtName("username")),
// 					stringvalidator.AlsoRequires(parentBlock.AtName("access_token")),
// 				),
// 			},
// 		},
// 		"username": schema.StringAttribute{
// 			Description: "Morpheus username for authentication, required if password is set",
// 			Optional:    true,
// 			Validators: []validator.String{
// 				stringvalidator.AlsoRequires(parentBlock.AtName("password")),
// 			},
// 		},
// 		"password": schema.StringAttribute{
// 			Description: "Morpheus password for authentication, required if username is set",
// 			Optional:    true,
// 			Sensitive:   true,
// 			Validators: []validator.String{
// 				stringvalidator.AlsoRequires(parentBlock.AtName("username")),
// 			},
// 		},
// 		"access_token": schema.StringAttribute{
// 			Description: "Morpheus access token for authentication",
// 			Optional:    true,
// 			Sensitive:   true,
// 			Validators: []validator.String{
// 				stringvalidator.ConflictsWith(parentBlock.AtName("username")),
// 				stringvalidator.ConflictsWith(parentBlock.AtName("password")),
// 				stringvalidator.ConflictsWith(parentBlock.AtName("tenant_subdomain")),
// 			},
// 		},
// 		"tenant_subdomain": schema.StringAttribute{
// 			Description: "Morpheus tenant subdomain used for authentication",
// 			Optional:    true,
// 		},
// 		"insecure": schema.BoolAttribute{
// 			Description: "Explicitly allow the provider to perform " +
// 				"\"insecure\" SSL requests. If omitted, " +
// 				"default value is `false`",
// 			Optional: true,
// 		},
// 	}
// }
