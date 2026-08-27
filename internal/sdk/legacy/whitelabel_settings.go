package morpheus

import "fmt"

var (
	// WhitelabelSettingsPath is the API endpoint for whitelabel settings
	WhitelabelSettingsPath = "/api/whitelabel-settings"
)

// WhitelabelSettings structures for use in request and response payloads
type WhitelabelSettings struct {
	Account struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"account"`
	Enabled                   bool   `json:"enabled"`
	ApplianceName             string `json:"applianceName"`
	DisableSupportMenu        bool   `json:"disableSupportMenu"`
	HeaderLogo                string `json:"headerLogo"`
	FooterLogo                string `json:"footerLogo"`
	LoginLogo                 string `json:"loginLogo"`
	Favicon                   string `json:"favicon"`
	HeaderBgColor             string `json:"headerBgColor"`
	HeaderFgColor             string `json:"headerFgColor"`
	NavBgColor                string `json:"navBgColor"`
	NavFgColor                string `json:"navFgColor"`
	NavHoverColor             string `json:"navHoverColor"`
	PrimaryButtonBgColor      string `json:"primaryButtonBgColor"`
	PrimaryButtonFgColor      string `json:"primaryButtonFgColor"`
	PrimaryButtonHoverBgColor string `json:"primaryButtonHoverBgColor"`
	PrimaryButtonHoverFgColor string `json:"primaryButtonHoverFgColor"`
	FooterBgColor             string `json:"footerBgColor"`
	FooterFgColor             string `json:"footerFgColor"`
	LoginBgColor              string `json:"loginBgColor"`
	OverrideCSS               string `json:"overrideCss"`
	CopyrightString           string `json:"copyrightString"`
	TermsOfUse                string `json:"termsOfUse"`
	PrivacyPolicy             string `json:"privacyPolicy"`
	SupportMenuLinks          []struct {
		URL       string `json:"url"`
		Label     string `json:"label"`
		LabelCode string `json:"labelCode"`
	} `json:"supportMenuLinks"`
}

type GetWhitelabelSettingsResult struct {
	WhitelabelSettings *WhitelabelSettings `json:"whitelabelSettings"`
}

type UpdateWhitelabelSettingsResult struct {
	Success            bool                `json:"success"`
	Message            string              `json:"msg"`
	Errors             map[string]string   `json:"errors"`
	WhitelabelSettings *WhitelabelSettings `json:"whitelabelSettings"`
}

type UpdateWhiteLabelImageResult struct {
	Success bool              `json:"success"`
	Message string            `json:"msg"`
	Errors  map[string]string `json:"errors"`
}

// Client request methods
func (client *Client) GetWhitelabelSettings(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        WhitelabelSettingsPath,
		QueryParams: req.QueryParams,
		Result:      &GetWhitelabelSettingsResult{},
	})
}

func (client *Client) UpdateWhitelabelSettings(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        WhitelabelSettingsPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateWhitelabelSettingsResult{},
	})
}

func (client *Client) UpdateWhitelabelImages(id int64, filePayload []*FilePayload, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:         "POST",
		Path:           fmt.Sprintf("%s/images", WhitelabelSettingsPath),
		IsMultiPart:    true,
		MultiPartFiles: filePayload,
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
		Result: &UpdateWhiteLabelImageResult{},
	})
}

func (client *Client) ResetWhitelabelImage(req *Request, imageType string) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/images/%s", WhitelabelSettingsPath, imageType),
		QueryParams: req.QueryParams,
		Result:      &UpdateWhiteLabelImageResult{},
	})
}
