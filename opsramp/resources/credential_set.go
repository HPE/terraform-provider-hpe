// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"
	"strings"

	"github.com/HPE/terraform-provider-hpe/opsramp/client"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource               = &CredentialSetResource{}
	_ resource.ResourceWithModifyPlan = &CredentialSetResource{}
)

type CredentialSetResource struct {
	BaseResource
}

type CredentialSetModel struct {
	Client                          types.String `tfsdk:"client"`
	Id                              types.String `tfsdk:"id"`
	Name                            types.String `tfsdk:"name"`
	Description                     types.String `tfsdk:"description"`
	CredentialType                  types.String `tfsdk:"credential_type"`
	Secure                          types.Bool   `tfsdk:"secure"`
	Port                            types.Int64  `tfsdk:"port"`
	SnmpVersion                     types.String `tfsdk:"snmp_version"`
	AutoEnableMode                  types.Bool   `tfsdk:"auto_enable_mode"`
	Universal                       types.Bool   `tfsdk:"universal"`
	SpSecure                        types.Bool   `tfsdk:"sp_secure"`
	SpPort                          types.Int64  `tfsdk:"sp_port"`
	TimeoutMs                       types.Int64  `tfsdk:"timeout_ms"`
	DomainName                      types.String `tfsdk:"domain_name"`
	UserName                        types.String `tfsdk:"user_name"`
	Password                        types.String `tfsdk:"password"`
	TransportType                   types.String `tfsdk:"transport_type"`
	Community                       types.String `tfsdk:"community"`
	SpUserName                      types.String `tfsdk:"sp_user_name"`
	SpPassword                      types.String `tfsdk:"sp_password"`
	SpAuthScope                     types.String `tfsdk:"sp_auth_scope"`
	FileAuthScope                   types.String `tfsdk:"file_auth_scope"`
	EsxUserName                     types.String `tfsdk:"esx_user_name"`
	EsxPassword                     types.String `tfsdk:"esx_password"`
	SpNameSpace                     types.String `tfsdk:"sp_name_space"`
	AuthProtocol                    types.String `tfsdk:"auth_protocol"`
	EncryptPassword                 types.String `tfsdk:"encrypt_password"`
	SnmpContext                     types.String `tfsdk:"snmp_context"`
	SecurityLevel                   types.String `tfsdk:"security_level"`
	SecurityName                    types.String `tfsdk:"security_name"`
	ApiEndPoint                     types.String `tfsdk:"api_end_point"`
	AccountId                       types.String `tfsdk:"account_id"`
	AccountName                     types.String `tfsdk:"account_name"`
	AccountKey                      types.String `tfsdk:"account_key"`
	ManagementCertificate           types.String `tfsdk:"management_certificate"`
	ManagementCertificatePassphrase types.String `tfsdk:"management_certificate_passphrase"`
	SshCredentialType               types.String `tfsdk:"ssh_credential_type"`
	CollectorType                   types.String `tfsdk:"collector_type"`
	EnablePassword                  types.String `tfsdk:"enable_password"`
}

func NewCredentialSet() resource.Resource {
	return &CredentialSetResource{}
}

func (r *CredentialSetResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_credential_set"
}

func (r *CredentialSetResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an OpsRamp Credential Set.",
		Attributes: map[string]schema.Attribute{
			"client": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The client UUID. If not specified, uses the provider's tenant.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the credential set.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Credential set name.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "The description of the credential set.",
			},
			"credential_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Credential type.",
				Validators: []validator.String{
					stringvalidator.OneOf("AWS", "SNMP", "WINDOWS", "VNC", "VMWARE"),
				},
			},
			"secure": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Whether the connection is secure.",
			},
			"port": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "Port number.",
			},
			"snmp_version": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "SNMP version.",
			},
			"auto_enable_mode": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Whether auto enable mode is on.",
			},
			"universal": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Whether this credential set is universal.",
			},
			"sp_secure": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Whether the service processor connection is secure.",
			},
			"sp_port": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "Service processor port.",
			},
			"timeout_ms": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "Timeout in milliseconds.",
			},
			"domain_name": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Domain name.",
			},
			"user_name": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "User name credentials.",
			},
			"password": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "Password credentials.",
			},
			"transport_type": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Transport layer type.",
				Validators: []validator.String{
					stringvalidator.OneOf("HTTP", "SNMP", "SSH", "TELNET"),
				},
			},
			"community": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "SNMP community string.",
			},
			"sp_user_name": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Service processor user name.",
			},
			"sp_password": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "Service processor password.",
			},
			"sp_auth_scope": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Service processor auth scope.",
			},
			"file_auth_scope": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "File auth scope.",
			},
			"esx_user_name": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "ESX user name.",
			},
			"esx_password": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "ESX password.",
			},
			"sp_name_space": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Service processor namespace.",
			},
			"auth_protocol": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Authentication protocol.",
			},
			"encrypt_password": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "Encryption password.",
			},
			"snmp_context": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "SNMP context.",
			},
			"security_level": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Security level.",
				Validators: []validator.String{
					stringvalidator.OneOf("NOAUTHNOPRIV", "AUTHPRIV", "AUTHNOPRIV", "BASIC", "OAUTH2"),
				},
			},
			"security_name": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Security name.",
			},
			"api_end_point": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "API endpoint.",
			},
			"account_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Account ID.",
			},
			"account_name": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Account name.",
			},
			"account_key": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "Account key.",
			},
			"management_certificate": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "Management certificate.",
			},
			"management_certificate_passphrase": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "Management certificate passphrase.",
			},
			"ssh_credential_type": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "SSH credential type.",
				Validators: []validator.String{
					stringvalidator.OneOf("PASSWORD", "KEYPAIR"),
				},
			},
			"collector_type": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Collector type.",
				Validators: []validator.String{
					stringvalidator.OneOf("API", "CLI", "APIANDCLI", "SMIS", "UNKNOWN"),
				},
			},
			"enable_password": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "Enable password.",
			},
		},
	}
}

func buildCredentialSetRequest(plan CredentialSetModel) client.CredentialSet {
	cs := client.CredentialSet{
		Name:           plan.Name.ValueString(),
		CredentialType: plan.CredentialType.ValueString(),
		Secure:         plan.Secure.ValueBool(),
		AutoEnableMode: plan.AutoEnableMode.ValueBool(),
		Universal:      plan.Universal.ValueBool(),
		SpSecure:       plan.SpSecure.ValueBool(),
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		cs.Description = plan.Description.ValueString()
	}
	if !plan.Port.IsNull() && !plan.Port.IsUnknown() {
		cs.Port = int(plan.Port.ValueInt64())
	}
	if !plan.SnmpVersion.IsNull() && !plan.SnmpVersion.IsUnknown() {
		cs.SnmpVersion = plan.SnmpVersion.ValueString()
	}
	if !plan.SpPort.IsNull() && !plan.SpPort.IsUnknown() {
		cs.SpPort = int(plan.SpPort.ValueInt64())
	}
	if !plan.TimeoutMs.IsNull() && !plan.TimeoutMs.IsUnknown() {
		cs.TimeoutMs = int(plan.TimeoutMs.ValueInt64())
	}
	if !plan.DomainName.IsNull() && !plan.DomainName.IsUnknown() {
		cs.DomainName = plan.DomainName.ValueString()
	}
	if !plan.UserName.IsNull() && !plan.UserName.IsUnknown() {
		cs.UserName = plan.UserName.ValueString()
	}
	if !plan.Password.IsNull() && !plan.Password.IsUnknown() {
		cs.Password = plan.Password.ValueString()
	}
	if !plan.TransportType.IsNull() && !plan.TransportType.IsUnknown() {
		cs.TransportType = plan.TransportType.ValueString()
	}
	if !plan.Community.IsNull() && !plan.Community.IsUnknown() {
		cs.Community = plan.Community.ValueString()
	}
	if !plan.SpUserName.IsNull() && !plan.SpUserName.IsUnknown() {
		cs.SpUserName = plan.SpUserName.ValueString()
	}
	if !plan.SpPassword.IsNull() && !plan.SpPassword.IsUnknown() {
		cs.SpPassword = plan.SpPassword.ValueString()
	}
	if !plan.SpAuthScope.IsNull() && !plan.SpAuthScope.IsUnknown() {
		cs.SpAuthScope = plan.SpAuthScope.ValueString()
	}
	if !plan.FileAuthScope.IsNull() && !plan.FileAuthScope.IsUnknown() {
		cs.FileAuthScope = plan.FileAuthScope.ValueString()
	}
	if !plan.EsxUserName.IsNull() && !plan.EsxUserName.IsUnknown() {
		cs.EsxUserName = plan.EsxUserName.ValueString()
	}
	if !plan.EsxPassword.IsNull() && !plan.EsxPassword.IsUnknown() {
		cs.EsxPassword = plan.EsxPassword.ValueString()
	}
	if !plan.SpNameSpace.IsNull() && !plan.SpNameSpace.IsUnknown() {
		cs.SpNameSpace = plan.SpNameSpace.ValueString()
	}
	if !plan.AuthProtocol.IsNull() && !plan.AuthProtocol.IsUnknown() {
		cs.AuthProtocol = plan.AuthProtocol.ValueString()
	}
	if !plan.EncryptPassword.IsNull() && !plan.EncryptPassword.IsUnknown() {
		cs.EncryptPassword = plan.EncryptPassword.ValueString()
	}
	if !plan.SnmpContext.IsNull() && !plan.SnmpContext.IsUnknown() {
		cs.SnmpContext = plan.SnmpContext.ValueString()
	}
	if !plan.SecurityLevel.IsNull() && !plan.SecurityLevel.IsUnknown() {
		cs.SecurityLevel = plan.SecurityLevel.ValueString()
	}
	if !plan.SecurityName.IsNull() && !plan.SecurityName.IsUnknown() {
		cs.SecurityName = plan.SecurityName.ValueString()
	}
	if !plan.ApiEndPoint.IsNull() && !plan.ApiEndPoint.IsUnknown() {
		cs.ApiEndPoint = plan.ApiEndPoint.ValueString()
	}
	if !plan.AccountId.IsNull() && !plan.AccountId.IsUnknown() {
		cs.AccountId = plan.AccountId.ValueString()
	}
	if !plan.AccountName.IsNull() && !plan.AccountName.IsUnknown() {
		cs.AccountName = plan.AccountName.ValueString()
	}
	if !plan.AccountKey.IsNull() && !plan.AccountKey.IsUnknown() {
		cs.AccountKey = plan.AccountKey.ValueString()
	}
	if !plan.ManagementCertificate.IsNull() && !plan.ManagementCertificate.IsUnknown() {
		cs.ManagementCertificate = plan.ManagementCertificate.ValueString()
	}
	if !plan.ManagementCertificatePassphrase.IsNull() && !plan.ManagementCertificatePassphrase.IsUnknown() {
		cs.ManagementCertificatePassphrase = plan.ManagementCertificatePassphrase.ValueString()
	}
	if !plan.SshCredentialType.IsNull() && !plan.SshCredentialType.IsUnknown() {
		cs.SshCredentialType = plan.SshCredentialType.ValueString()
	}
	if !plan.CollectorType.IsNull() && !plan.CollectorType.IsUnknown() {
		cs.CollectorType = plan.CollectorType.ValueString()
	}
	if !plan.EnablePassword.IsNull() && !plan.EnablePassword.IsUnknown() {
		cs.EnablePassword = plan.EnablePassword.ValueString()
	}

	return cs
}

func mapCredentialSetToState(resp *client.CredentialSet, state *CredentialSetModel) {
	state.Id = types.StringValue(resp.UniqueId)
	state.Name = types.StringValue(resp.Name)

	if resp.Description != "" {
		state.Description = types.StringValue(resp.Description)
	}

	state.CredentialType = types.StringValue(resp.CredentialType)
	state.Secure = types.BoolValue(resp.Secure)
	state.AutoEnableMode = types.BoolValue(resp.AutoEnableMode)
	state.Universal = types.BoolValue(resp.Universal)
	state.SpSecure = types.BoolValue(resp.SpSecure)

	if resp.Port != 0 {
		state.Port = types.Int64Value(int64(resp.Port))
	}
	if resp.SnmpVersion != "" {
		state.SnmpVersion = types.StringValue(resp.SnmpVersion)
	}
	if resp.SpPort != 0 {
		state.SpPort = types.Int64Value(int64(resp.SpPort))
	}
	if resp.TimeoutMs != 0 {
		state.TimeoutMs = types.Int64Value(int64(resp.TimeoutMs))
	}
	if resp.DomainName != "" {
		state.DomainName = types.StringValue(resp.DomainName)
	}
	if resp.UserName != "" {
		state.UserName = types.StringValue(resp.UserName)
	}
	if resp.TransportType != "" {
		state.TransportType = types.StringValue(resp.TransportType)
	}
	if resp.SpUserName != "" {
		state.SpUserName = types.StringValue(resp.SpUserName)
	}
	if resp.SpAuthScope != "" {
		state.SpAuthScope = types.StringValue(resp.SpAuthScope)
	}
	if resp.FileAuthScope != "" {
		state.FileAuthScope = types.StringValue(resp.FileAuthScope)
	}
	if resp.EsxUserName != "" {
		state.EsxUserName = types.StringValue(resp.EsxUserName)
	}
	if resp.SpNameSpace != "" {
		state.SpNameSpace = types.StringValue(resp.SpNameSpace)
	}
	if resp.AuthProtocol != "" {
		state.AuthProtocol = types.StringValue(resp.AuthProtocol)
	}
	if resp.SnmpContext != "" {
		state.SnmpContext = types.StringValue(resp.SnmpContext)
	}
	if resp.SecurityLevel != "" {
		state.SecurityLevel = types.StringValue(resp.SecurityLevel)
	}
	if resp.SecurityName != "" {
		state.SecurityName = types.StringValue(resp.SecurityName)
	}
	if resp.ApiEndPoint != "" {
		state.ApiEndPoint = types.StringValue(resp.ApiEndPoint)
	}
	if resp.AccountId != "" {
		state.AccountId = types.StringValue(resp.AccountId)
	}
	if resp.AccountName != "" {
		state.AccountName = types.StringValue(resp.AccountName)
	}
	if resp.SshCredentialType != "" {
		state.SshCredentialType = types.StringValue(resp.SshCredentialType)
	}
	if resp.CollectorType != "" {
		state.CollectorType = types.StringValue(resp.CollectorType)
	}
}

func (r *CredentialSetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan CredentialSetModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.apiClient.TenantId
	if !plan.Client.IsNull() && plan.Client.ValueString() != "" {
		tenantId = plan.Client.ValueString()
	}

	csReq := buildCredentialSetRequest(plan)
	created, err := r.apiClient.CreateCredentialSet(tenantId, csReq)
	if err != nil {
		resp.Diagnostics.AddError("Creation Error", err.Error())

		return
	}

	// Only set the server-generated ID; preserve plan values for everything else
	plan.Id = types.StringValue(created.UniqueId)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (r *CredentialSetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state CredentialSetModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.apiClient.TenantId
	if !state.Client.IsNull() && state.Client.ValueString() != "" {
		tenantId = state.Client.ValueString()
	}

	existing, err := r.apiClient.GetCredentialSet(tenantId, state.Id.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "not found") {
			resp.State.RemoveResource(ctx)

			return
		}
		resp.Diagnostics.AddError("Read Error", err.Error())

		return
	}

	mapCredentialSetToState(existing, &state)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *CredentialSetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state CredentialSetModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.apiClient.TenantId
	if !state.Client.IsNull() && state.Client.ValueString() != "" {
		tenantId = state.Client.ValueString()
	}

	csReq := buildCredentialSetRequest(plan)
	_, err := r.apiClient.UpdateCredentialSet(tenantId, state.Id.ValueString(), csReq)
	if err != nil {
		resp.Diagnostics.AddError("Update Error", err.Error())

		return
	}

	// Preserve plan values; ID doesn't change on update
	plan.Id = state.Id

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (r *CredentialSetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state CredentialSetModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.apiClient.TenantId
	if !state.Client.IsNull() && state.Client.ValueString() != "" {
		tenantId = state.Client.ValueString()
	}

	err := r.apiClient.DeleteCredentialSet(tenantId, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Delete Error", err.Error())

		return
	}

	resp.State.RemoveResource(ctx)
}
