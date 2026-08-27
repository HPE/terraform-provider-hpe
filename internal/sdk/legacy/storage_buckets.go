package morpheus

import (
	"fmt"
)

var (
	// StorageBucketsPath is the API endpoint for storage buckets
	StorageBucketsPath = "/api/storage-buckets"
)

// StorageBucket structures for use in request and response payloads
type StorageBucket struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Active       bool   `json:"active"`
	AccountID    int64  `json:"accountId"`
	ProviderType string `json:"providerType"`
	Config       struct {
		LocationType       string      `json:"locationType"`
		Location           string      `json:"location"`
		StorageClass       string      `json:"storageClass"`
		ClientEmail        string      `json:"clientEmail"`
		PrivateKey         string      `json:"privateKey"`
		ProjectId          string      `json:"projectId"`
		AccessKey          string      `json:"accessKey"`
		SecretKey          string      `json:"secretKey"`
		SecretKeyHash      string      `json:"secretKeyHash"`
		Endpoint           string      `json:"endpoint"`
		BasePath           string      `json:"basePath"`
		Host               string      `json:"host"`
		StsAssumeRole      string      `json:"stsAssumeRole"`
		UseHostCredentials interface{} `json:"useHostCredentials"`
		ExportFolder       string      `json:"exportFolder"`
		Permissions        []string    `json:"permissions"`
		AdminPermissions   []string    `json:"adminPermissions"`
	} `json:"config"`
	BucketName                string      `json:"bucketName"`
	ReadOnly                  bool        `json:"readOnly"`
	DefaultBackupTarget       bool        `json:"defaultBackupTarget"`
	DefaultDeploymentTarget   bool        `json:"defaultDeploymentTarget"`
	DefaultVirtualImageTarget bool        `json:"defaultVirtualImageTarget"`
	CopyToStore               bool        `json:"copyToStore"`
	RetentionPolicyType       interface{} `json:"retentionPolicyType"`
	RetentionPolicyDays       interface{} `json:"retentionPolicyDays"`
	RetentionProvider         interface{} `json:"retentionProvider"`
}

type ListStorageBucketsResult struct {
	StorageBuckets *[]StorageBucket `json:"storageBuckets"`
	Meta           *MetaResult      `json:"meta"`
}

type GetStorageBucketResult struct {
	StorageBucket *StorageBucket `json:"storageBucket"`
}

type CreateStorageBucketResult struct {
	Success       bool              `json:"success"`
	Message       string            `json:"msg"`
	Errors        map[string]string `json:"errors"`
	StorageBucket *StorageBucket    `json:"storageBucket"`
}

type UpdateStorageBucketResult struct {
	CreateStorageBucketResult
}

type DeleteStorageBucketResult struct {
	DeleteResult
}

// Client request methods
func (client *Client) ListStorageBuckets(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        StorageBucketsPath,
		QueryParams: req.QueryParams,
		Result:      &ListStorageBucketsResult{},
	})
}

func (client *Client) GetStorageBucket(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", StorageBucketsPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetStorageBucketResult{},
	})
}

func (client *Client) CreateStorageBucket(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        StorageBucketsPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateStorageBucketResult{},
	})
}

func (client *Client) UpdateStorageBucket(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", StorageBucketsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateStorageBucketResult{},
	})
}

func (client *Client) DeleteStorageBucket(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", StorageBucketsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteStorageBucketResult{},
	})
}

// FindStorageBucketByName gets an existing storageBucket by name
func (client *Client) FindStorageBucketByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListStorageBuckets(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListStorageBucketsResult)
	storageBucketCount := len(*listResult.StorageBuckets)
	if storageBucketCount != 1 {
		return resp, fmt.Errorf("found %d Storage Buckets for %v", storageBucketCount, name)
	}
	firstRecord := (*listResult.StorageBuckets)[0]
	storageBucketID := firstRecord.ID
	return client.GetStorageBucket(storageBucketID, &Request{})
}
