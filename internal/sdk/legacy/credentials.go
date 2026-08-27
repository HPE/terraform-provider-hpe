package morpheus

import (
	"fmt"
	"time"
)

var (
	// CredentialsPath is the API endpoint for credentials
	CredentialsPath = "/api/credentials"
)

// Credential structures for use in request and response payloads
type Credential struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Type struct {
		ID   int64  `json:"id"`
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"type"`
	Integration struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"integration"`
	Description  string `json:"description"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	PasswordHash string `json:"passwordHash"`
	AuthKey      struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"authKey"`
	AuthPath      string      `json:"authPath"`
	ExternalID    interface{} `json:"externalId"`
	RefType       interface{} `json:"refType"`
	RefID         interface{} `json:"refId"`
	Category      interface{} `json:"category"`
	Scope         string      `json:"scope"`
	Status        string      `json:"status"`
	StatusMessage interface{} `json:"statusMessage"`
	StatusDate    interface{} `json:"statusDate"`
	Enabled       bool        `json:"enabled"`
	Account       struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"account"`
	User struct {
		ID          int64  `json:"id"`
		Username    string `json:"username"`
		DisplayName string `json:"displayName"`
	} `json:"user"`
	DateCreated time.Time `json:"dateCreated"`
	LastUpdated time.Time `json:"lastUpdated"`
	Config      struct {
		GrantType        string `json:"grantType"`
		AccessTokenUrl   string `json:"accessTokenUrl"`
		ClientAuth       string `json:"clientAuth"`
		ClientSecret     string `json:"clientSecret"`
		Scope            string `json:"scope"`
		ClientId         string `json:"clientId"`
		ClientSecretHash string `json:"clientSecretHash"`
	} `json:"config"`
}

type ListCredentialsResult struct {
	Credentials *[]Credential `json:"credentials"`
	Meta        *MetaResult   `json:"meta"`
}

type GetCredentialResult struct {
	Credential *Credential `json:"credential"`
}

type CreateCredentialResult struct {
	Success    bool              `json:"success"`
	Message    string            `json:"msg"`
	Errors     map[string]string `json:"errors"`
	Credential *Credential       `json:"credential"`
}

type UpdateCredentialResult struct {
	CreateCredentialResult
}

type DeleteCredentialResult struct {
	DeleteResult
}

// Client request methods
func (client *Client) ListCredentials(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        CredentialsPath,
		QueryParams: req.QueryParams,
		Result:      &ListCredentialsResult{},
	})
}

func (client *Client) GetCredential(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", CredentialsPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetCredentialResult{},
	})
}

func (client *Client) CreateCredential(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        CredentialsPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateCredentialResult{},
	})
}

func (client *Client) UpdateCredential(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", CredentialsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateCredentialResult{},
	})
}

func (client *Client) DeleteCredential(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", CredentialsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteCredentialResult{},
	})
}

// FindCredentialByName gets an existing credential by name
func (client *Client) FindCredentialByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListCredentials(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListCredentialsResult)
	credentialCount := len(*listResult.Credentials)
	if credentialCount != 1 {
		return resp, fmt.Errorf("found %d Credentials for %v", credentialCount, name)
	}
	firstRecord := (*listResult.Credentials)[0]
	credentialID := firstRecord.ID
	return client.GetCredential(credentialID, &Request{})
}
