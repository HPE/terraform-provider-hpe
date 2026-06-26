package morpheus

import (
	"fmt"
)

var (
	// SoftwareLicensesPath is the API endpoint for software licenses
	SoftwareLicensesPath = "/api/provisioning-licenses"
)

// SoftwareLicense structures for use in request and response payloads
type SoftwareLicense struct {
	ID               int64           `json:"id"`
	Name             string          `json:"name"`
	Description      string          `json:"description"`
	LicenseType      LicenseType     `json:"licenseType"`
	LicenseKey       string          `json:"licenseKey"`
	OrgName          string          `json:"orgName"`
	FullName         string          `json:"fullName"`
	LicenseVersion   string          `json:"licenseVersion"`
	Account          Account         `json:"account"`
	Copies           int64           `json:"copies"`
	Reservationcount int64           `json:"reservationCount"`
	Tenants          []interface{}   `json:"tenants"`
	Virtualimages    []Virtualimages `json:"virtualImages"`
}

type LicenseType struct {
	ID   int64  `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type Virtualimages struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Account struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Reservation struct {
	ResourceID   int64  `json:"resourceId"`
	ResourceType string `json:"resourceType"`
}

type ListSoftwareLicensesResult struct {
	SoftwareLicenses *[]SoftwareLicense `json:"licenses"`
	Meta             *MetaResult        `json:"meta"`
}

type GetSoftwareLicenseResult struct {
	SoftwareLicense *SoftwareLicense `json:"license"`
}

type GetSoftwareLicenseReservationsResult struct {
	SoftwareLicenses *[]Reservation `json:"licenses"`
}

type CreateSoftwareLicenseResult struct {
	Success         bool              `json:"success"`
	Message         string            `json:"msg"`
	Errors          map[string]string `json:"errors"`
	SoftwareLicense *SoftwareLicense  `json:"license"`
}

type UpdateSoftwareLicenseResult struct {
	CreateSoftwareLicenseResult
}

type DeleteSoftwareLicenseResult struct {
	DeleteResult
}

// Client request methods
func (client *Client) ListSoftwareLicenses(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        SoftwareLicensesPath,
		QueryParams: req.QueryParams,
		Result:      &ListSoftwareLicensesResult{},
	})
}

func (client *Client) GetSoftwareLicense(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", SoftwareLicensesPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetSoftwareLicenseResult{},
	})
}

func (client *Client) GetSoftwareLicenseReservations(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d/reservations", SoftwareLicensesPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetSoftwareLicenseReservationsResult{},
	})
}

func (client *Client) CreateSoftwareLicense(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        SoftwareLicensesPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateSoftwareLicenseResult{},
	})
}

func (client *Client) UpdateSoftwareLicense(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", SoftwareLicensesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateSoftwareLicenseResult{},
	})
}

func (client *Client) DeleteSoftwareLicense(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", SoftwareLicensesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteSoftwareLicenseResult{},
	})
}

// FindSoftwareLicenseByName gets an existing software license by name
func (client *Client) FindSoftwareLicenseByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListSoftwareLicenses(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListSoftwareLicensesResult)
	softwareLicenseCount := len(*listResult.SoftwareLicenses)
	if softwareLicenseCount != 1 {
		return resp, fmt.Errorf("found %d Software Licenses for %v", softwareLicenseCount, name)
	}
	firstRecord := (*listResult.SoftwareLicenses)[0]
	softwareLicenseID := firstRecord.ID
	return client.GetSoftwareLicense(softwareLicenseID, &Request{})
}
