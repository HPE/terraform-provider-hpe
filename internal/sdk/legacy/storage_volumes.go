package morpheus

import (
	"fmt"
)

// StorageVolumesPath is the API endpoint for Storage Volumes
var StorageVolumesPath = "/api/storage-volumes"

// StorageVolume structures for use in request and response payloads
type StorageVolume struct {
	ID                   int64       `json:"id"`
	Name                 string      `json:"name"`
	ShortName            string      `json:"shortName"`
	Description          string      `json:"description"`
	ControllerId         int64       `json:"controllerId"`
	ControllerMountPoint string      `json:"controllerMountPoint"`
	Resizeable           interface{} `json:"resizeable"`
	PlanResizable        interface{} `json:"planResizable"`
	Size                 interface{} `json:"size"`
	StorageType          interface{} `json:"storageType"`
	RootVolume           interface{} `json:"rootVolume"`
	UnitNumber           string      `json:"unitNumber"`
	DeviceName           string      `json:"deviceName"`
	DeviceDisplayName    string      `json:"deviceDisplayName"`
	Type                 struct {
		ID   int64  `json:"id"`
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"type"`
	TypeId           int64  `json:"typeId"`
	Category         string `json:"category"`
	Status           string `json:"status"`
	StatusMessage    string `json:"statusMessage"`
	ConfigurableIOPS bool   `json:"configurableIOPS"`
	MaxStorage       int64  `json:"maxStorage"`
	DisplayOrder     int64  `json:"displayOrder"`
	MaxIOPS          string `json:"maxIOPS"`
	Uuid             string `json:"uuid"`
	Active           bool   `json:"active"`
	Zone             struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"zone"`
	ZoneId    int64 `json:"zoneId"`
	Datastore struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"datastore"`
	DatastoreId   int64  `json:"datastoreId"`
	StorageGroup  string `json:"storageGroup"`
	Namespace     string `json:"namespace"`
	StorageServer string `json:"storageServer"`
	Source        string `json:"source"`
	Owner         struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"owner"`
}

type ListStorageVolumesResult struct {
	StorageVolumes *[]StorageVolume `json:"storageVolumes"`
	Meta           *MetaResult      `json:"meta"`
}

type GetStorageVolumeResult struct {
	StorageVolume *StorageVolume `json:"storageVolume"`
}

type CreateStorageVolumeResult struct {
	Success       bool              `json:"success"`
	Message       string            `json:"msg"`
	Errors        map[string]string `json:"errors"`
	StorageVolume *StorageVolume    `json:"storageVolume"`
}

type UpdateStorageVolumeResult struct {
	CreateStorageVolumeResult
}

type DeleteStorageVolumeResult struct {
	DeleteResult
}

// Client request methods
func (client *Client) ListStorageVolumes(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        StorageVolumesPath,
		QueryParams: req.QueryParams,
		Result:      &ListStorageVolumesResult{},
	})
}

func (client *Client) GetStorageVolume(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", StorageVolumesPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetStorageVolumeResult{},
	})
}

func (client *Client) CreateStorageVolume(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        StorageVolumesPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateStorageVolumeResult{},
	})
}

func (client *Client) UpdateStorageVolume(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", StorageVolumesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateStorageVolumeResult{},
	})
}

func (client *Client) DeleteStorageVolume(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", StorageVolumesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteStorageVolumeResult{},
	})
}

func (client *Client) FindStorageVolumeByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListStorageVolumes(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListStorageVolumesResult)
	StorageVolumeCount := len(*listResult.StorageVolumes)
	if StorageVolumeCount != 1 {
		return resp, fmt.Errorf("found %d storage volumes named %v", StorageVolumeCount, name)
	}
	firstRecord := (*listResult.StorageVolumes)[0]
	StorageVolumeID := firstRecord.ID
	return client.GetStorageVolume(StorageVolumeID, &Request{})
}
