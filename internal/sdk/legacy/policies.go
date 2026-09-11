package morpheus

import (
	"fmt"
)

var (
	// PoliciesPath is the API endpoint for policies
	PoliciesPath = "/api/policies"
)

// Policy structures for use in request and response payloads
type Policy struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Enabled     bool   `json:"enabled"`
	Description string `json:"description"`
	EachUser    bool   `json:"eachUser"`
	PolicyType  struct {
		ID   int64  `json:"id"`
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"policyType"`
	// Config is an opaque bag whose field types vary by policy type and are
	// polymorphic in the API (e.g. maxCores/valueListId come back as a number
	// or a string). The provider never reads these fields (the policies data
	// source only uses id/name/policyType), so decode config as a raw map to
	// avoid unmarshal failures when listing real-world policies.
	Config map[string]interface{} `json:"config"`
	Owner struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"owner"`
	RefID   int64  `json:"refId"`
	RefType string `json:"refType"`
	Role    struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"role"`
	Site struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"site"`
	User struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"user"`
	Zone struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"zone"`
	Accounts []struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	}
}

type ListPoliciesResult struct {
	Policies *[]Policy   `json:"policies"`
	Meta     *MetaResult `json:"meta"`
}

type GetPolicyResult struct {
	Policy *Policy `json:"policy"`
}

type CreatePolicyResult struct {
	Success bool              `json:"success"`
	Message string            `json:"msg"`
	Errors  map[string]string `json:"errors"`
	Policy  *Policy           `json:"policy"`
}

type UpdatePolicyResult struct {
	CreatePolicyResult
}

type DeletePolicyResult struct {
	DeleteResult
}

// ListPolicies list all policies
func (client *Client) ListPolicies(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        PoliciesPath,
		QueryParams: req.QueryParams,
		Result:      &ListPoliciesResult{},
	})
}

// GetPolicy gets a policy
func (client *Client) GetPolicy(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", PoliciesPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetPolicyResult{},
	})
}

// CreatePolicy creates a new policy
func (client *Client) CreatePolicy(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        PoliciesPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreatePolicyResult{},
	})
}

// UpdatePolicy updates an existing policy
func (client *Client) UpdatePolicy(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", PoliciesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdatePolicyResult{},
	})
}

// DeletePolicy deletes an existing policy
func (client *Client) DeletePolicy(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", PoliciesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeletePolicyResult{},
	})
}

// FindPolicyByName gets an existing policy by name
func (client *Client) FindPolicyByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListPolicies(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListPoliciesResult)
	policyCount := len(*listResult.Policies)
	if policyCount != 1 {
		return resp, fmt.Errorf("found %d policies named %v", policyCount, name)
	}
	firstRecord := (*listResult.Policies)[0]
	policyID := firstRecord.ID
	return client.GetPolicy(policyID, &Request{})
}
