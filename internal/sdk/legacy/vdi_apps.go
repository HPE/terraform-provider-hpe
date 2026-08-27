package morpheus

import (
	"fmt"
	"time"
)

var (
	// VDIAppsPath is the API endpoint for vdi apps
	VDIAppsPath = "/api/vdi-apps"
)

// VDIApp structures for use in request and response payloads
type VDIApp struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	LaunchPrefix string    `json:"launchPrefix"`
	Logo         string    `json:"logo"`
	DateCreated  time.Time `json:"dateCreated"`
	LastUpdated  time.Time `json:"lastUpdated"`
}

type ListVDIAppsResult struct {
	VDIApps *[]VDIApp   `json:"vdiApps"`
	Meta    *MetaResult `json:"meta"`
}

type GetVDIAppResult struct {
	VDIApp *VDIApp `json:"vdiApp"`
}

type CreateVDIAppResult struct {
	Success bool              `json:"success"`
	Message string            `json:"msg"`
	Errors  map[string]string `json:"errors"`
	VDIApp  *VDIApp           `json:"vdiApp"`
}

type UpdateVDIAppResult struct {
	CreateVDIAppResult
}

type DeleteVDIAppResult struct {
	DeleteResult
}

// Client request methods
func (client *Client) ListVDIApps(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        VDIAppsPath,
		QueryParams: req.QueryParams,
		Result:      &ListVDIAppsResult{},
	})
}

func (client *Client) GetVDIApp(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", VDIAppsPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetVDIAppResult{},
	})
}

func (client *Client) CreateVDIApp(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        VDIAppsPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateVDIAppResult{},
	})
}

func (client *Client) UpdateVDIApp(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", VDIAppsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateVDIAppResult{},
	})
}

func (client *Client) DeleteVDIApp(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", VDIAppsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteVDIAppResult{},
	})
}

func (client *Client) UploadVDIAppLogo(id int64, filePayload []*FilePayload, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:         "PUT",
		Path:           fmt.Sprintf("/api/vdi-apps/%d", id),
		IsMultiPart:    true,
		MultiPartFiles: filePayload,
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
		Result: &UpdateVDIAppResult{},
	})
}

// FindVDIAppByName gets an existing vdi app by name
func (client *Client) FindVDIAppByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListVDIApps(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListVDIAppsResult)
	vdiAppCount := len(*listResult.VDIApps)
	if vdiAppCount != 1 {
		return resp, fmt.Errorf("found %d VDI Apps for %v", vdiAppCount, name)
	}
	firstRecord := (*listResult.VDIApps)[0]
	vdiAppID := firstRecord.ID
	return client.GetVDIApp(vdiAppID, &Request{})
}
