package morpheus

var (
	// ProvisioningSettingsPath is the API endpoint for provisioning settings
	ProvisioningSettingsPath = "/api/provisioning-settings"
)

// ProvisioningSettings structures for use in request and response payloads
type ProvisioningSettings struct {
	AllowZoneSelection          bool   `json:"allowZoneSelection"`
	AllowServerSelection        bool   `json:"allowServerSelection"`
	RequireEnvironments         bool   `json:"requireEnvironments"`
	ShowPricing                 bool   `json:"showPricing"`
	HideDatastoreStats          bool   `json:"hideDatastoreStats"`
	CrossTenantNamingPolicies   bool   `json:"crossTenantNamingPolicies"`
	ReuseSequence               bool   `json:"reuseSequence"`
	ShowConsoleKeyboardSettings bool   `json:"showConsoleKeyboardSettings"`
	CloudInitUsername           string `json:"cloudInitUsername"`
	CloudInitPassword           string `json:"cloudInitPassword"`
	CloudInitPasswordHash       string `json:"cloudInitPasswordHash"`
	Cloudinitkeypair            struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"cloudInitKeyPair"`
	WindowsPassword     string `json:"windowsPassword"`
	WindowsPasswordHash string `json:"windowsPasswordHash"`
	PXERootPassword     string `json:"pxeRootPassword"`
	PXERootPasswordHash string `json:"pxeRootPasswordHash"`
	DefaultTemplateType struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
		Code string `json:"code"`
	} `json:"defaultTemplateType"`
	DeployStorageProvider struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"deployStorageProvider"`
}

type GetProvisioningSettingsResult struct {
	ProvisioningSettings *ProvisioningSettings `json:"provisioningSettings"`
}

type UpdateProvisioningSettingsResult struct {
	Success              bool                  `json:"success"`
	Message              string                `json:"msg"`
	Errors               map[string]string     `json:"errors"`
	ProvisioningSettings *ProvisioningSettings `json:"provisioningSettings"`
}

// Client request methods
func (client *Client) GetProvisioningSettings(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        ProvisioningSettingsPath,
		QueryParams: req.QueryParams,
		Result:      &GetProvisioningSettingsResult{},
	})
}

func (client *Client) UpdateProvisioningSettings(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        ProvisioningSettingsPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateProvisioningSettingsResult{},
	})
}
