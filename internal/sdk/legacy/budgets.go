package morpheus

import (
	"fmt"
)

var (
	// BudgetsPath is the API endpoint for budgets
	BudgetsPath = "/api/budgets"
)

// Budget structures for use in request and response payloads
type Budget struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Account     struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"account"`
	RefScope      string      `json:"refScope"`
	RefType       interface{} `json:"refType"`
	RefId         interface{} `json:"refId"`
	RefName       string      `json:"refName"`
	Interval      string      `json:"interval"`
	Period        string      `json:"period"`
	Year          string      `json:"year"`
	ResourceType  string      `json:"resourceType"`
	TimeZone      string      `json:"timezone"`
	StartDate     string      `json:"startDate"`
	EndDate       string      `json:"endDate"`
	Active        bool        `json:"active"`
	Enabled       bool        `json:"enabled"`
	Rollover      bool        `json:"rollover"`
	Costs         []float64   `json:"costs"`
	AverageCost   float64     `json:"averageCost"`
	TotalCost     float64     `json:"totalCost"`
	Currency      string      `json:"currency"`
	WarningLimit  interface{} `json:"warningLimit"`
	OverLimit     interface{} `json:"overLimit"`
	ExternalId    interface{} `json:"externalId"`
	InternalId    interface{} `json:"internalId"`
	CreatedById   int64       `json:"createdById"`
	CreatedByName string      `json:"createdByName"`
	UpdatedById   interface{} `json:"updatedById"`
	UpdatedByName interface{} `json:"updatedByName"`
	DateCreated   string      `json:"dateCreated"`
	LastUpdated   string      `json:"lastUpdated"`
	Stats         Stats       `json:"stats"`
}

type Intervals struct {
	Index     int64   `json:"index"`
	Year      string  `json:"year"`
	ShortYear string  `json:"shortYear"`
	Budget    float64 `json:"budget"`
	Cost      float64 `json:"cost"`
}

type Current struct {
	EstimatedCost float64 `json:"estimatedCost"`
	LastCost      float64 `json:"lastCost"`
}
type Stats struct {
	AverageCost    float64     `json:"averageCost"`
	TotalCost      float64     `json:"totalCost"`
	Currency       string      `json:"currency"`
	ConversionRate int64       `json:"conversionRate"`
	Intervals      []Intervals `json:"intervals"`
	Current        Current     `json:"current"`
}

type ListBudgetsResult struct {
	Budgets *[]Budget   `json:"budgets"`
	Meta    *MetaResult `json:"meta"`
}

type GetBudgetResult struct {
	Budget *Budget `json:"budget"`
}

type CreateBudgetResult struct {
	Success bool              `json:"success"`
	Message string            `json:"msg"`
	Errors  map[string]string `json:"errors"`
	Budget  *Budget           `json:"budget"`
}

type UpdateBudgetResult struct {
	CreateBudgetResult
}

type DeleteBudgetResult struct {
	DeleteResult
}

// Client request methods
func (client *Client) ListBudgets(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        BudgetsPath,
		QueryParams: req.QueryParams,
		Result:      &ListBudgetsResult{},
	})
}

func (client *Client) GetBudget(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", BudgetsPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetBudgetResult{},
	})
}

func (client *Client) CreateBudget(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        BudgetsPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateBudgetResult{},
	})
}

func (client *Client) UpdateBudget(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", BudgetsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateBudgetResult{},
	})
}

func (client *Client) DeleteBudget(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", BudgetsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteBudgetResult{},
	})
}

// FindBudgetByName gets an existing budget by name
func (client *Client) FindBudgetByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListBudgets(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListBudgetsResult)
	budgetCount := len(*listResult.Budgets)
	if budgetCount != 1 {
		return resp, fmt.Errorf("found %d Budgets for %v", budgetCount, name)
	}
	firstRecord := (*listResult.Budgets)[0]
	budgetID := firstRecord.ID
	return client.GetBudget(budgetID, &Request{})
}
