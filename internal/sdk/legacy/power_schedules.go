package morpheus

import (
	"fmt"
)

var (
	// PowerSchedulesPath is the API endpoint for power schedules
	PowerSchedulesPath = "/api/power-schedules"
)

// PowerSchedule structures for use in request and response payloads
type PowerSchedule struct {
	ID               int64   `json:"id"`
	Name             string  `json:"name"`
	Description      string  `json:"description"`
	ScheduleType     string  `json:"scheduleType"`
	ScheduleTimeZone string  `json:"scheduleTimezone"`
	Enabled          bool    `json:"enabled"`
	SundayOn         float64 `json:"sundayOn"`
	SundayOff        float64 `json:"sundayOff"`
	MondayOn         float64 `json:"mondayOn"`
	MondayOff        float64 `json:"mondayOff"`
	TuesdayOn        float64 `json:"tuesdayOn"`
	TuesdayOff       float64 `json:"tuesdayOff"`
	WednesdayOn      float64 `json:"wednesdayOn"`
	WednesdayOff     float64 `json:"wednesdayOff"`
	ThursdayOn       float64 `json:"thursdayOn"`
	ThursdayOff      float64 `json:"thursdayOff"`
	FridayOn         float64 `json:"fridayOn"`
	FridayOff        float64 `json:"fridayOff"`
	SaturdayOn       float64 `json:"saturdayOn"`
	SaturdayOff      float64 `json:"saturdayOff"`
}

// ListPowerSchedulesResult structure parses the list power schedules response payload
type ListPowerSchedulesResult struct {
	PowerSchedules *[]PowerSchedule `json:"schedules"`
	Meta           *MetaResult      `json:"meta"`
}

// GetPowerScheduleResult structure parses the get power schedule response payload
type GetPowerScheduleResult struct {
	PowerSchedule *PowerSchedule `json:"schedule"`
}

// CreatePowerScheduleResult structure parses the create power schedule response payload
type CreatePowerScheduleResult struct {
	Success       bool              `json:"success"`
	Message       string            `json:"msg"`
	Errors        map[string]string `json:"errors"`
	PowerSchedule *PowerSchedule    `json:"schedule"`
}

// UpdatePowerScheduleResult structure parses the update power schedule response payload
type UpdatePowerScheduleResult struct {
	CreatePowerScheduleResult
}

// DeletePowerScheduleResult structure parses the delete power schedule response payload
type DeletePowerScheduleResult struct {
	DeleteResult
}

// ListPowerSchedules lists all power schedules
// https://apidocs.morpheusdata.com/#power-schedules
func (client *Client) ListPowerSchedules(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        PowerSchedulesPath,
		QueryParams: req.QueryParams,
		Result:      &ListPowerSchedulesResult{},
	})
}

// GetPowerSchedule gets an existing power schedule
// https://apidocs.morpheusdata.com/#get-a-specific-power-schedule
func (client *Client) GetPowerSchedule(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", PowerSchedulesPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetPowerScheduleResult{},
	})
}

// CreatePowerSchedule creates a new power schedule
// https://apidocs.morpheusdata.com/#create-a-power-schedule
func (client *Client) CreatePowerSchedule(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        PowerSchedulesPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreatePowerScheduleResult{},
	})
}

// UpdatePowerSchedule updates an existing power schedule
// https://apidocs.morpheusdata.com/#update-a-power-schedule
func (client *Client) UpdatePowerSchedule(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", PowerSchedulesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdatePowerScheduleResult{},
	})
}

// DeletePowerSchedule deletes an existing power schedule
// https://apidocs.morpheusdata.com/#delete-a-power-schedule
func (client *Client) DeletePowerSchedule(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", PowerSchedulesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeletePowerScheduleResult{},
	})
}

// FindPowerScheduleByName gets an existing power schedule by name
func (client *Client) FindPowerScheduleByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListPowerSchedules(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListPowerSchedulesResult)
	powerScheduleCount := len(*listResult.PowerSchedules)
	if powerScheduleCount != 1 {
		return resp, fmt.Errorf("found %d Power Schedules for %v", powerScheduleCount, name)
	}
	firstRecord := (*listResult.PowerSchedules)[0]
	powerScheduleID := firstRecord.ID
	return client.GetPowerSchedule(powerScheduleID, &Request{})
}
