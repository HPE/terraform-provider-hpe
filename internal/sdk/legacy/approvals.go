package morpheus

import (
	"fmt"
)

var (
	// ApprovalsPath is the API endpoint for approvals
	ApprovalsPath     = "/api/approvals"
	ApprovalItemsPath = "/api/approval-items"
)

// Approval structures for use in request and response payloads
type Approval struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	InternalID   string `json:"internalId"`
	ExternalID   string `json:"externalId"`
	ExternalName string `json:"externalName"`
	RequestType  string `json:"requestType"`
	Account      struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"account"`
	Approver struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"approver"`
	AccountIntegration interface{} `json:"accountIntegration"`
	Status             string      `json:"status"`
	ErrorMessage       string      `json:"errorMessage"`
	DateCreated        string      `json:"dateCreated"`
	LastUpdated        string      `json:"lastUpdated"`
	RequestBy          string      `json:"requestBy"`
}

type ApprovalItem struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	InternalID   string `json:"internalId"`
	ExternalID   string `json:"externalId"`
	ExternalName string `json:"externalName"`
	ApprovedBy   string `json:"approvedBy"`
	DeniedBy     string `json:"deniedBy"`
	Status       string `json:"status"`
	ErrorMessage string `json:"errorMessage"`
	DateCreated  string `json:"dateCreated"`
	LastUpdated  string `json:"lastUpdated"`
	DateApproved string `json:"dateApproved"`
	DateDenied   string `json:"dateDenied"`
	Approval     struct {
		ID int64 `json:"id"`
	} `json:"approal"`
	Reference struct {
		ID          int64  `json:"id"`
		Type        string `json:"type"`
		Name        string `json:"name"`
		DisplayName string `json:"displayName"`
	} `json:"reference"`
}

// ListApprovalsResult structure parses the list approvals response payload
type ListApprovalsResult struct {
	Approvals *[]Approval `json:"approvals"`
	Meta      *MetaResult `json:"meta"`
}

// GetApprovalResult structure parses the get approval response payload
type GetApprovalResult struct {
	Approval *Approval `json:"approval"`
}

// GetApprovalItemResult structure parses the get approval response payload
type GetApprovalItemResult struct {
	ApprovalItem *ApprovalItem `json:"approvalItem"`
}

// UpdateApprovalItemResult structure parses the create approval response payload
type UpdateApprovalItemResult struct {
	Success bool              `json:"success"`
	Message string            `json:"msg"`
	Errors  map[string]string `json:"errors"`
}

// ListApprovals lists all approvals
func (client *Client) ListApprovals(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        ApprovalsPath,
		QueryParams: req.QueryParams,
		Result:      &ListApprovalsResult{},
	})
}

// GetApproval gets an existing approval
func (client *Client) GetApproval(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", ApprovalsPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetApprovalResult{},
	})
}

// UpdateApproval updates an existing approval
func (client *Client) UpdateApprovalItem(id int64, action string, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d/%s", ApprovalItemsPath, id, action),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateApprovalItemResult{},
	})
}

// GetApprovalItem gets an existing approval item
func (client *Client) GetApprovalItem(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", ApprovalItemsPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetApprovalItemResult{},
	})
}

// FindApprovalByName gets an existing approval by name
func (client *Client) FindApprovalByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListApprovals(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListApprovalsResult)
	approvalCount := len(*listResult.Approvals)
	if approvalCount != 1 {
		return resp, fmt.Errorf("found %d Approvals for %v", approvalCount, name)
	}
	firstRecord := (*listResult.Approvals)[0]
	approvalID := firstRecord.ID
	return client.GetApproval(approvalID, &Request{})
}
