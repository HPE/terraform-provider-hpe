package morpheus

var (
	// BackupSettingsPath is the API endpoint for backup settings
	BackupSettingsPath = "/api/backup-settings"
)

// BackupSettings structures for use in request and response payloads
type BackupSettings struct {
	BackupsEnabled       bool `json:"backupsEnabled"`
	CreateBackups        bool `json:"createBackups"`
	BackupAppliance      bool `json:"backupAppliance"`
	DefaultStorageBucket struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"defaultStorageBucket"`
	DefaultSchedule struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"defaultSchedule"`
	RetentionCount int64 `json:"retentionCount"`
}

type GetBackupSettingsResult struct {
	BackupSettings *BackupSettings `json:"backupSettings"`
}

type UpdateBackupSettingsResult struct {
	Success        bool              `json:"success"`
	Message        string            `json:"msg"`
	Errors         map[string]string `json:"errors"`
	BackupSettings *BackupSettings   `json:"backupSettings"`
}

// Client request methods
func (client *Client) GetBackupSettings(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        BackupSettingsPath,
		QueryParams: req.QueryParams,
		Result:      &GetBackupSettingsResult{},
	})
}

func (client *Client) UpdateBackupSettings(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        BackupSettingsPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateBackupSettingsResult{},
	})
}
