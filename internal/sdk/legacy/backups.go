package morpheus

import (
	"fmt"
	"time"
)

var (
	// BackupsPath is the API endpoint for backups
	BackupsPath = "/api/backups"
)

// Backup structures for use in request and response payloads
type Backup struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	LocationType string `json:"locationType"`
	Instance     struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"instance"`
	ContainerId int64 `json:"containerId"`
	Job         struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"job"`
	Schedule struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
		Cron string `json:"cron"`
	} `json:"schedule"`
	RetentionCount int64 `json:"retentionCount"`
	BackupType     struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
		Code string `json:"code"`
	} `json:"backupType"`
	BackupProvider struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
		Code string `json:"code"`
	} `json:"backupProvider"`
	StorageProvider struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"StorageProvider"`
	BackupRepository string `json:"backupRepository"`
	CronExpression   string `json:"cronExpression"`
	NextFire         string `json:"nextFire"`
	LastStatus       string `json:"lastStatus"`
	LastResult       struct {
		ID          int64  `json:"id"`
		Status      string `json:"status"`
		DateCreated string `json:"dateCreated"`
	} `json:"lastResult"`
	Stats struct {
		TotalSize       int64    `json:"totalSize"`
		Success         int64    `json:"success"`
		TotalCompleted  int64    `json:"totalCompleted"`
		SuccessRate     float64  `json:"successRate"`
		FailRate        float64  `json:"failRate"`
		AvgSize         int64    `json:"avgSize"`
		LastFiveResults []string `json:"lastFiveResults"`
		FailedRate      float64  `json:"failedRate"`
	} `json:"stats"`
	Enabled     bool      `json:"enabled"`
	DateCreated time.Time `json:"dateCreated"`
	LastUpdated time.Time `json:"lastUpdated"`
}

// ListBackupsResult structure parses the list backups response payload
type ListBackupsResult struct {
	Backups *[]Backup   `json:"backups"`
	Meta    *MetaResult `json:"meta"`
}

type GetBackupResult struct {
	Backup *Backup `json:"backup"`
}

type CreateBackupResult struct {
	Success bool              `json:"success"`
	Message string            `json:"msg"`
	Errors  map[string]string `json:"errors"`
	Backup  *Backup           `json:"backup"`
}

type UpdateBackupResult struct {
	CreateBackupResult
}

type DeleteBackupResult struct {
	DeleteResult
}

// ListBackups lists all backups
func (client *Client) ListBackups(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        BackupsPath,
		QueryParams: req.QueryParams,
		Result:      &ListBackupsResult{},
	})
}

// GetBackup gets an backup
func (client *Client) GetBackup(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", BackupsPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetBackupResult{},
	})
}

// CreateBackup creates a new backup
func (client *Client) CreateBackup(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        BackupsPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateBackupResult{},
	})
}

// UpdateBackup updates an existing backup
func (client *Client) UpdateBackup(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", BackupsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateBackupResult{},
	})
}

// DeleteBackup deletes an existing backup
func (client *Client) DeleteBackup(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", BackupsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteBackupResult{},
	})
}

// ExecuteBackup executes a backup
func (client *Client) ExecuteBackup(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        fmt.Sprintf("%s/%d/execute", BackupsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateBackupResult{},
	})
}

// FindBackupByName gets an existing backup by name
func (client *Client) FindBackupByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListBackups(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListBackupsResult)
	backupsCount := len(*listResult.Backups)
	if backupsCount != 1 {
		return resp, fmt.Errorf("found %d backups for %v", backupsCount, name)
	}
	firstRecord := (*listResult.Backups)[0]
	backupID := firstRecord.ID
	return client.GetBackup(backupID, &Request{})
}
