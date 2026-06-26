package morpheus

var (
	// MonitoringSettingsPath is the API endpoint for monitoring settings
	MonitoringSettingsPath = "/api/monitoring-settings"
)

// MonitoringSettings structures for use in request and response payloads
type MonitoringSettings struct {
	AutoManageChecks      bool  `json:"autoManageChecks"`
	AvailabilityTimeFrame int64 `json:"availabilityTimeFrame"`
	AvailabilityPrecision int64 `json:"availabilityPrecision"`
	DefaultCheckInterval  int64 `json:"defaultCheckInterval"`
	ServiceNow            struct {
		Enabled     bool `json:"enabled"`
		Integration struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
		} `json:"integration"`
		NewIncidentAction   string `json:"newIncidentAction"`
		CloseIncidentAction string `json:"closeIncidentAction"`
		InfoMapping         string `json:"infoMapping"`
		WarningMapping      string `json:"warningMapping"`
		CriticalMapping     string `json:"criticalMapping"`
	} `json:"serviceNow"`
	NewRelic struct {
		Enabled    bool   `json:"enabled"`
		LicenseKey string `json:"licenseKey"`
	} `json:"newRelic"`
}

type GetMonitoringSettingsResult struct {
	MonitoringSettings *MonitoringSettings `json:"monitoringSettings"`
}

type UpdateMonitoringSettingsResult struct {
	Success            bool                `json:"success"`
	Message            string              `json:"msg"`
	Errors             map[string]string   `json:"errors"`
	MonitoringSettings *MonitoringSettings `json:"monitoringSettings"`
}

// Client request methods
func (client *Client) GetMonitoringSettings(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        MonitoringSettingsPath,
		QueryParams: req.QueryParams,
		Result:      &GetMonitoringSettingsResult{},
	})
}

func (client *Client) UpdateMonitoringSettings(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        MonitoringSettingsPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateMonitoringSettingsResult{},
	})
}
