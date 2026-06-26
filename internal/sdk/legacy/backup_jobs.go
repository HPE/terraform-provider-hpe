package morpheus

import (
	"fmt"
	"time"
)

var (
	// BackupJobsPath is the API endpoint for backup jobs
	BackupJobsPath = "/api/backups/jobs"
)

// BackupJob structures for use in request and response payloads
type BackupJob struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	Code           string `json:"code"`
	RetentionCount int64  `json:"retentionCount"`
	Schedule       struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
		Cron string `json:"cron"`
	} `json:"schedule"`
	ExternalId     string `json:"externalId"`
	BackupProvider struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
		Code string `json:"code"`
	} `json:"backupProvider"`
	BackupRepository string `json:"backupRepository"`
	CronExpression   string `json:"cronExpression"`
	NextFire         string `json:"nextFire"`
	Source           string `json:"source"`
	Visibility       string `json:"visibility"`
	Account          struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"account"`
	DateCreated time.Time `json:"dateCreated"`
	LastUpdated time.Time `json:"lastUpdated"`
	Backups     []struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"backups"`
}

// ListBackupsResult structure parses the list backups response payload
type ListBackupJobsResult struct {
	BackupJobs *[]BackupJob `json:"jobs"`
	Meta       *MetaResult  `json:"meta"`
}

type GetBackupJobResult struct {
	BackupJob *BackupJob `json:"job"`
}

type CreateBackupJobResult struct {
	Success   bool              `json:"success"`
	Message   string            `json:"msg"`
	Errors    map[string]string `json:"errors"`
	BackupJob *BackupJob        `json:"job"`
}

type UpdateBackupJobResult struct {
	CreateBackupJobResult
}

type DeleteBackupJobResult struct {
	DeleteResult
}

// ListBackupJobs lists all backup jobs
func (client *Client) ListBackupJobs(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        BackupJobsPath,
		QueryParams: req.QueryParams,
		Result:      &ListBackupJobsResult{},
	})
}

// GetBackupJob gets an existing backup job
func (client *Client) GetBackupJob(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", BackupJobsPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetBackupJobResult{},
	})
}

// CreateBackupJob creates a new backup job
func (client *Client) CreateBackupJob(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        BackupJobsPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateBackupJobResult{},
	})
}

// UpdateBackupJob updates an existing backup job
func (client *Client) UpdateBackupJob(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", BackupJobsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateBackupJobResult{},
	})
}

// DeleteBackup deletes an existing backup job
func (client *Client) DeleteBackupJob(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", BackupJobsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteBackupJobResult{},
	})
}

// ExecuteBackupJob executes a backup job
func (client *Client) ExecuteBackupJob(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        fmt.Sprintf("%s/%d/execute", BackupJobsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateBackupJobResult{},
	})
}

// FindBackupJobByName gets an existing backup job by name
func (client *Client) FindBackupJobByName(name string) (*Response, error) {
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
