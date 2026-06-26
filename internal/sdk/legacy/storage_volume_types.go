package morpheus

import (
	"fmt"
)

var (
	// StorageVolumeTypesPath is the API endpoint for Storage Volume types
	StorageVolumeTypesPath = "/api/storage-volume-types"
)

// StorageVolumeType structures for use in request and response payloads
type StorageVolumeType struct {
	ID      int64  `json:"id"`
	Code    string `json:"code"`
	Account struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"account"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	DisplayOrder      int64  `json:"displayOrder"`
	DefaultType       bool   `json:"defaultType"`
	CustomLabel       bool   `json:"customLabel"`
	CustomSize        bool   `json:"customSize"`
	CustomSizeOptions []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
		Size  string `json:"size"`
	} `json:"customSizeOptions"`
	ConfigurableIOPS bool         `json:"configurableIOPS"`
	HasDatastore     bool         `json:"hasDatastore"`
	Category         string       `json:"category"`
	Enabled          bool         `json:"enabled"`
	OptionTypes      []OptionType `json:"optionTypes"`
}

type ListStorageVolumeTypesResult struct {
	StorageVolumeTypes *[]StorageVolumeType `json:"storageVolumeTypes"`
	Meta               *MetaResult          `json:"meta"`
}

type GetStorageVolumeTypeResult struct {
	StorageVolumeType *StorageVolumeType `json:"storageVolumeType"`
}

// Client request methods

func (client *Client) ListStorageVolumeTypes(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        StorageVolumeTypesPath,
		QueryParams: req.QueryParams,
		Result:      &ListStorageVolumeTypesResult{},
	})
}

func (client *Client) GetStorageVolumeType(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", StorageVolumeTypesPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetStorageVolumeTypeResult{},
	})
}

func (client *Client) FindStorageVolumeTypeByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListStorageVolumeTypes(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListStorageVolumeTypesResult)
	StorageVolumeTypeCount := len(*listResult.StorageVolumeTypes)
	if StorageVolumeTypeCount != 1 {
		return resp, fmt.Errorf("found %d storage volume types named %v", StorageVolumeTypeCount, name)
	}
	firstRecord := (*listResult.StorageVolumeTypes)[0]
	StorageVolumeTypeID := firstRecord.ID
	return client.GetStorageVolumeType(StorageVolumeTypeID, &Request{})
}
