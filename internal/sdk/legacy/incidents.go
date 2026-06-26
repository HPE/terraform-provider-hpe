package morpheus

import (
	"fmt"
)

var (
	// IncidentsPath is the API endpoint for incidents
	IncidentsPath = "/api/monitoring/incidents"
)

// Incident structures for use in request and response payloads
type Incident struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Comment    string `json:"comment"`
	Resolution string `json:"resolution"`
	Status     string `json:"status"`
	Severity   string `json:"severity"`
	InUptime   bool   `json:"inUptime"`
	StartDate  string `json:"startDate"`
	EndDate    string `json:"endDate"`
}

// ListIncidentsResult structure parses the list incidents response payload
type ListIncidentsResult struct {
	Incidents *[]Incident `json:"incidents"`
	Meta      *MetaResult `json:"meta"`
}

type GetIncidentResult struct {
	Incident *Incident `json:"incident"`
}

type CreateIncidentResult struct {
	Success  bool              `json:"success"`
	Message  string            `json:"msg"`
	Errors   map[string]string `json:"errors"`
	Incident *Incident         `json:"incident"`
}

type UpdateIncidentResult struct {
	CreateIncidentResult
}

type DeleteIncidentResult struct {
	DeleteResult
}

// Client request methods

func (client *Client) ListIncidents(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        IncidentsPath,
		QueryParams: req.QueryParams,
		Result:      &ListIncidentsResult{},
	})
}

func (client *Client) GetIncident(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", IncidentsPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetIncidentResult{},
	})
}

func (client *Client) CreateIncident(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        IncidentsPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateIncidentResult{},
	})
}

func (client *Client) UpdateIncident(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", IncidentsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateIncidentResult{},
	})
}

func (client *Client) DeleteIncident(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", IncidentsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteIncidentResult{},
	})
}

// helper functions
func (client *Client) FindIncidentByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListIncidents(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListIncidentsResult)
	incidentCount := len(*listResult.Incidents)
	if incidentCount != 1 {
		return resp, fmt.Errorf("found %d Incidents for %v", incidentCount, name)
	}
	firstRecord := (*listResult.Incidents)[0]
	incidentID := firstRecord.ID
	return client.GetIncident(incidentID, &Request{})
}
