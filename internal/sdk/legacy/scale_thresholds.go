package morpheus

import (
	"fmt"
	"time"
)

var (
	// ScaleThresholdsPath is the API endpoint for scale thresholds
	ScaleThresholdsPath = "/api/scale-thresholds"
)

// ScaleThreshold structures for use in request and response payloads
type ScaleThreshold struct {
	ID             int64     `json:"id"`
	Name           string    `json:"name"`
	Type           string    `json:"type"`
	AutoUp         bool      `json:"autoUp"`
	AutoDown       bool      `json:"autoDown"`
	MinCount       int64     `json:"minCount"`
	MaxCount       int64     `json:"maxCount"`
	ScaleIncrement int64     `json:"scaleIncrement"`
	CpuEnabled     bool      `json:"cpuEnabled"`
	MinCpu         float64   `json:"minCpu"`
	MaxCpu         float64   `json:"maxCpu"`
	MemoryEnabled  bool      `json:"memoryEnabled"`
	MinMemory      float64   `json:"minMemory"`
	MaxMemory      float64   `json:"maxMemory"`
	DiskEnabled    bool      `json:"diskEnabled"`
	MinDisk        float64   `json:"minDisk"`
	MaxDisk        float64   `json:"maxDisk"`
	DateCreated    time.Time `json:"dateCreated"`
	LastUpdated    time.Time `json:"lastUpdated"`
}

type ListScaleThresholdsResult struct {
	ScaleThresholds *[]ScaleThreshold `json:"scaleThresholds"`
	Meta            *MetaResult       `json:"meta"`
}

type GetScaleThresholdResult struct {
	ScaleThreshold *ScaleThreshold `json:"scaleThreshold"`
}

type CreateScaleThresholdResult struct {
	Success        bool              `json:"success"`
	Message        string            `json:"msg"`
	Errors         map[string]string `json:"errors"`
	ScaleThreshold *ScaleThreshold   `json:"scaleThreshold"`
}

type UpdateScaleThresholdResult struct {
	CreateScaleThresholdResult
}

type DeleteScaleThresholdResult struct {
	DeleteResult
}

// Client request methods
func (client *Client) ListScaleThresholds(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        ScaleThresholdsPath,
		QueryParams: req.QueryParams,
		Result:      &ListScaleThresholdsResult{},
	})
}

func (client *Client) GetScaleThreshold(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", ScaleThresholdsPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetScaleThresholdResult{},
	})
}

func (client *Client) CreateScaleThreshold(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        ScaleThresholdsPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateScaleThresholdResult{},
	})
}

func (client *Client) UpdateScaleThreshold(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", ScaleThresholdsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateScaleThresholdResult{},
	})
}

func (client *Client) DeleteScaleThreshold(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", ScaleThresholdsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteScaleThresholdResult{},
	})
}

// FindScaleThresholdByName gets an existing scaleThreshold by name
func (client *Client) FindScaleThresholdByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListScaleThresholds(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListScaleThresholdsResult)
	scaleThresholdCount := len(*listResult.ScaleThresholds)
	if scaleThresholdCount != 1 {
		return resp, fmt.Errorf("found %d Scale Thresholds for %v", scaleThresholdCount, name)
	}
	firstRecord := (*listResult.ScaleThresholds)[0]
	scaleThresholdID := firstRecord.ID
	return client.GetScaleThreshold(scaleThresholdID, &Request{})
}
