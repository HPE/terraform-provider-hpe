package morpheus

import (
	"fmt"
	"time"
)

var (
	// VDIPoolsPath is the API endpoint for vdi pools
	VDIPoolsPath = "/api/vdi-pools"
)

// VDIPool structures for use in request and response payloads
type VDIPool struct {
	ID                               int64       `json:"id"`
	Name                             string      `json:"name"`
	Description                      string      `json:"description"`
	MinIdle                          int64       `json:"minIdle"`
	MaxIdle                          int64       `json:"maxIdle"`
	InitialPoolSize                  int64       `json:"initialPoolSize"`
	MaxPoolSize                      int64       `json:"maxPoolSize"`
	AllocationTimeoutMinutes         int64       `json:"allocationTimeoutMinutes"`
	PersistentUser                   bool        `json:"persistentUser"`
	Recyclable                       bool        `json:"recyclable"`
	Enabled                          bool        `json:"enabled"`
	AutoCreateLocalUserOnReservation bool        `json:"autoCreateLocalUserOnReservation"`
	AllowHypervisorConsole           bool        `json:"allowHypervisorConsole"`
	AllowCopy                        bool        `json:"allowCopy"`
	AllowPrinter                     bool        `json:"allowPrinter"`
	AllowFileShare                   bool        `json:"allowFileshare"`
	GuestConsoleJumpHost             interface{} `json:"guestConsoleJumpHost"`
	GuestConsoleJumpPort             interface{} `json:"guestConsoleJumpPort"`
	GuestConsoleJumpUsername         interface{} `json:"guestConsoleJumpUsername"`
	GuestConsoleJumpPassword         interface{} `json:"guestConsoleJumpPassword"`
	GuestConsoleJumpKeyPair          interface{} `json:"guestConsoleJumpKeypair"`
	Gateway                          interface{} `json:"gateway"`
	IconPath                         string      `json:"iconPath"`
	Logo                             string      `json:"logo"`
	Apps                             []struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"apps"`
	Owner struct {
		ID       int64  `json:"id"`
		Name     string `json:"name"`
		Username string `json:"username"`
	} `json:"owner"`
	Config struct {
		Group struct {
			ID   interface{} `json:"id"`
			Name string      `json:"name"`
		} `json:"group"`
		Cloud struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
		} `json:"cloud"`
		Type   string `json:"type"`
		Config struct {
			IsEc2           bool  `json:"isEC2"`
			IsVpcSelectable bool  `json:"isVpcSelectable"`
			NoAgent         bool  `json:"noAgent"`
			ResourcePoolID  int64 `json:"resourcePoolId"`
		} `json:"config"`
		Name    string `json:"name"`
		Volumes []struct {
			Name        string `json:"name"`
			RootVolume  bool   `json:"rootVolume"`
			TypeID      int64  `json:"typeId"`
			Size        int64  `json:"size"`
			StorageType int64  `json:"storageType"`
			DatastoreID string `json:"datastoreId"`
		} `json:"volumes"`
		HostName string `json:"hostName"`
		Layout   struct {
			ID   int64  `json:"id"`
			Code string `json:"code"`
		} `json:"layout"`
		NetworkInterfaces []struct {
			PrimaryInterface bool `json:"primaryInterface"`
			Network          struct {
				ID string `json:"id"`
			} `json:"network"`
			NetworkInterfaceTypeID     int64  `json:"networkInterfaceTypeId"`
			NetworkInterfaceTypeIDName string `json:"networkInterfaceTypeIdName"`
		} `json:"networkInterfaces"`
		Plan struct {
			ID   int64  `json:"id"`
			Code string `json:"code"`
		} `json:"plan"`
		Version string `json:"version"`
	} `json:"config"`
	Group struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"group"`
	Cloud struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"cloud"`
	UsedCount      int64     `json:"usedCount"`
	ReservedCount  int64     `json:"reservedCount"`
	PreparingCount int64     `json:"preparingCount"`
	IdleCount      int64     `json:"idleCount"`
	Status         string    `json:"status"`
	DateCreated    time.Time `json:"dateCreated"`
	LastUpdated    time.Time `json:"lastUpdated"`
}

type ListVDIPoolsResult struct {
	VDIPools *[]VDIPool  `json:"vdiPools"`
	Meta     *MetaResult `json:"meta"`
}

type GetVDIPoolResult struct {
	VDIPool *VDIPool `json:"vdiPool"`
}

type CreateVDIPoolResult struct {
	Success bool              `json:"success"`
	Message string            `json:"msg"`
	Errors  map[string]string `json:"errors"`
	VDIPool *VDIPool          `json:"vdiPool"`
}

type UpdateVDIPoolResult struct {
	CreateVDIPoolResult
}

type DeleteVDIPoolResult struct {
	DeleteResult
}

// Client request methods
func (client *Client) ListVDIPools(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        VDIPoolsPath,
		QueryParams: req.QueryParams,
		Result:      &ListVDIPoolsResult{},
	})
}

func (client *Client) GetVDIPool(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", VDIPoolsPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetVDIPoolResult{},
	})
}

func (client *Client) CreateVDIPool(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        VDIPoolsPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateVDIPoolResult{},
	})
}

func (client *Client) UpdateVDIPool(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", VDIPoolsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateVDIPoolResult{},
	})
}

func (client *Client) DeleteVDIPool(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", VDIPoolsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteVDIPoolResult{},
	})
}

func (client *Client) UploadVDIPoolLogo(id int64, filePayload []*FilePayload, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:         "PUT",
		Path:           fmt.Sprintf("/api/vdi-pools/%d", id),
		IsMultiPart:    true,
		MultiPartFiles: filePayload,
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
		Result: &UpdateVDIPoolResult{},
	})
}

// FindVDIPoolByName gets an existing vdi pool by name
func (client *Client) FindVDIPoolByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListVDIPools(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListVDIPoolsResult)
	vdiPoolCount := len(*listResult.VDIPools)
	if vdiPoolCount != 1 {
		return resp, fmt.Errorf("found %d VDI Pools for %v", vdiPoolCount, name)
	}
	firstRecord := (*listResult.VDIPools)[0]
	vdiPoolID := firstRecord.ID
	return client.GetVDIPool(vdiPoolID, &Request{})
}
