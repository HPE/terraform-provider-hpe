package morpheus

var (
	// GuidanceSettingsPath is the API endpoint for guidance settings
	GuidanceSettingsPath = "/api/guidance-settings"
)

// GuidanceSettings structures for use in request and response payloads
type GuidanceSettings struct {
	CpuAvgCutoffPower                    int64 `json:"cpuAvgCutoffPower"`
	CpuMaxCutoffPower                    int64 `json:"cpuMaxCutoffPower"`
	NetworkCutoffPower                   int64 `json:"networkCutoffPower"`
	CpuUpAvgStandardCutoffRightSize      int64 `json:"cpuUpAvgStandardCutoffRightSize"`
	CpuUpMaxStandardCutoffRightSize      int64 `json:"cpuUpMaxStandardCutoffRightSize"`
	MemoryUpAvgStandardCutoffRightSize   int64 `json:"memoryUpAvgStandardCutoffRightSize"`
	MemoryDownAvgStandardCutoffRightSize int64 `json:"memoryDownAvgStandardCutoffRightSize"`
	MemoryDownMaxStandardCutoffRightSize int64 `json:"memoryDownMaxStandardCutoffRightSize"`
}

type GetGuidanceSettingsResult struct {
	GuidanceSettings *GuidanceSettings `json:"guidanceSettings"`
}

type UpdateGuidanceSettingsResult struct {
	Success          bool              `json:"success"`
	Message          string            `json:"msg"`
	Errors           map[string]string `json:"errors"`
	GuidanceSettings *GuidanceSettings `json:"guidanceSettings"`
}

// Client request methods
func (client *Client) GetGuidanceSettings(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        GuidanceSettingsPath,
		QueryParams: req.QueryParams,
		Result:      &GetGuidanceSettingsResult{},
	})
}

func (client *Client) UpdateGuidanceSettings(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        GuidanceSettingsPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateGuidanceSettingsResult{},
	})
}
