package morpheus

import (
	"fmt"
)

var (
	// StorageServerTypesPath is the API endpoint for Storage Server types
	StorageServerTypesPath = "/api/storage-server-types"
)

// StorageServerType structures for use in request and response payloads
type StorageServerType struct {
	ID                     int64               `json:"id"`
	Code                   string              `json:"code"`
	Name                   string              `json:"name"`
	Description            string              `json:"description"`
	Enabled                bool                `json:"enabled"`
	Creatable              bool                `json:"creatable"`
	HasNamespaces          bool                `json:"hasNamespaces"`
	HasGroups              bool                `json:"hasGroups"`
	HasBlock               bool                `json:"hasBlock"`
	HasObject              bool                `json:"hasObject"`
	HasFile                bool                `json:"hasFile"`
	HasDatastore           bool                `json:"hasDatastore"`
	HasDisks               bool                `json:"hasDisks"`
	HasHosts               bool                `json:"hasHosts"`
	CreateNamespaces       bool                `json:"createNamespaces"`
	CreateGroup            bool                `json:"createGroup"`
	CreateBlock            bool                `json:"createBlock"`
	CreateObject           bool                `json:"createObject"`
	CreateFile             bool                `json:"createFile"`
	CreateDatastore        bool                `json:"createDatastore"`
	CreateDisk             bool                `json:"createDisk"`
	CreateHost             bool                `json:"createHost"`
	IconCode               string              `json:"iconCode"`
	HasFileBrowser         bool                `json:"hasFileBrowser"`
	OptionTypes            []OptionType        `json:"optionTypes"`
	GroupOptionTypes       []OptionType        `json:"groupOptionTypes"`
	BucketOptionTypes      []OptionType        `json:"bucketOptionTypes"`
	ShareOptionTypes       []OptionType        `json:"shareOptionTypes"`
	ShareAccessOptionTypes []OptionType        `json:"shareAccessOptionTypes"`
	StorageVolumeTypes     []StorageVolumeType `json:"storageVolumeTypes"`
}

type ListStorageServerTypesResult struct {
	StorageServerTypes *[]StorageServerType `json:"storageServerTypes"`
	Meta               *MetaResult          `json:"meta"`
}

type GetStorageServerTypeResult struct {
	StorageServerType *StorageServerType `json:"storageServerType"`
}

// Client request methods
func (client *Client) ListStorageServerTypes(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        StorageServerTypesPath,
		QueryParams: req.QueryParams,
		Result:      &ListStorageServerTypesResult{},
	})
}

func (client *Client) GetStorageServerType(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", StorageServerTypesPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetStorageServerTypeResult{},
	})
}

func (client *Client) FindStorageServerTypeByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListStorageServerTypes(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListStorageServerTypesResult)
	StorageServerTypeCount := len(*listResult.StorageServerTypes)
	if StorageServerTypeCount != 1 {
		return resp, fmt.Errorf("found %d storage server types named %v", StorageServerTypeCount, name)
	}
	firstRecord := (*listResult.StorageServerTypes)[0]
	StorageServerTypeID := firstRecord.ID
	return client.GetStorageServerType(StorageServerTypeID, &Request{})
}
