package morpheus

import (
	"fmt"
	"time"
)

var (
	// ArchivesPath is the API endpoint for archives
	ArchivesPath = "/api/archives/buckets"
)

// Archive structures for use in request and response payloads
type Archive struct {
	ID              int64  `json:"id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	StorageProvider struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"storageProvider"`
	Owner struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"owner"`
	CreatedBy struct {
		Username string `json:"username"`
	} `json:"createdBy"`
	IsPublic    bool      `json:"isPublic"`
	Code        string    `json:"code"`
	FilePath    string    `json:"filePath"`
	RawSize     int64     `json:"rawSize"`
	FileCount   int64     `json:"fileCount"`
	DateCreated time.Time `json:"dateCreated"`
	LastUpdated time.Time `json:"lastUpdated"`
	IsOwner     bool      `json:"isOwner"`
	Visibility  string    `json:"visibility"`
}

// ListArchivesResult structure parses the list archives response payload
type ListArchivesResult struct {
	Archives *[]Archive  `json:"archiveBuckets"`
	Meta     *MetaResult `json:"meta"`
}

type GetArchiveResult struct {
	Archive *Archive `json:"archiveBucket"`
}

type CreateArchiveResult struct {
	Success bool              `json:"success"`
	Message string            `json:"msg"`
	Errors  map[string]string `json:"errors"`
	Archive *Archive          `json:"archiveBucket"`
}

type UpdateArchiveResult struct {
	CreateArchiveResult
}

type DeleteArchiveResult struct {
	DeleteResult
}

// ListArchives lists all archives
func (client *Client) ListArchives(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        ArchivesPath,
		QueryParams: req.QueryParams,
		Result:      &ListArchivesResult{},
	})
}

// GetArchive gets an archive
func (client *Client) GetArchive(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", ArchivesPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetArchiveResult{},
	})
}

// CreateArchive creates a new archive
func (client *Client) CreateArchive(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        ArchivesPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateArchiveResult{},
	})
}

// UpdateArchive updates an existing archive
func (client *Client) UpdateArchive(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", ArchivesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateArchiveResult{},
	})
}

// DeleteArchive deletes an existing archive
func (client *Client) DeleteArchive(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", ArchivesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteArchiveResult{},
	})
}

// FindArchiveByName gets an existing archive by name
func (client *Client) FindArchiveByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListArchives(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListArchivesResult)
	archivesCount := len(*listResult.Archives)
	if archivesCount != 1 {
		return resp, fmt.Errorf("found %d Archives for %v", archivesCount, name)
	}
	firstRecord := (*listResult.Archives)[0]
	archiveID := firstRecord.ID
	return client.GetArchive(archiveID, &Request{})
}
