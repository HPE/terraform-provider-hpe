package morpheus

import "fmt"

var (
	// LogSettingsPath is the API endpoint for log settings
	LogSettingsPath = "/api/log-settings"
	SyslogRulesPath = "/api/log-settings/syslog-rules"
)

// LogSettings structures for use in request and response payloads
type LogSettings struct {
	Enabled       bool          `json:"enabled"`
	Retentiondays string        `json:"retentionDays"`
	SyslogRules   []SyslogRule  `json:"syslogRules"`
	Integrations  []interface{} `json:"integrations"`
}

type SyslogRule struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Rule string `json:"rule"`
}

// GetLogSettingsResult structure parses the get logSettings response payload
type GetLogSettingsResult struct {
	LogSettings *LogSettings `json:"logSettings"`
}

// UpdateLogSettingsResult structure parses the get logSettings response payload
type UpdateLogSettingsResult struct {
	Success bool              `json:"success"`
	Message string            `json:"msg"`
	Errors  map[string]string `json:"errors"`
}

// GetLogSettings gets the appliance logs settings
// https://apidocs.morpheusdata.com/#get-log-settings
func (client *Client) GetLogSettings(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        LogSettingsPath,
		QueryParams: req.QueryParams,
		Result:      &GetLogSettingsResult{},
	})
}

// UpdateLogSettings updates the appliance log settings
// https://apidocs.morpheusdata.com/#update-log-settings
func (client *Client) UpdateLogSettings(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        LogSettingsPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateLogSettingsResult{},
	})
}

// AddSyslogRule adds a syslog rule
// https://apidocs.morpheusdata.com/#add-syslog-rule
func (client *Client) AddSyslogRule(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        SyslogRulesPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateLogSettingsResult{},
	})
}

// DeleteSyslogRule deletes a syslog rule
// https://apidocs.morpheusdata.com/#delete-syslog-rules
func (client *Client) DeleteSyslogRule(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", SyslogRulesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateLogSettingsResult{},
	})
}
