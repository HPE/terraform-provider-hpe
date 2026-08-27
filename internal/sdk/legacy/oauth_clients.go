package morpheus

import (
	"fmt"
)

var (
	// OauthClientsPath is the API endpoint for oauth client
	OauthClientsPath = "/api/clients"
)

// OauthClient structures for use in request and response payloads
type OauthClient struct {
	ID                          int64    `json:"id"`
	ClientID                    string   `json:"clientId"`
	AccessTokenValiditySeconds  int64    `json:"accessTokenValiditySeconds"`
	RefreshTokenValiditySeconds int64    `json:"refreshTokenValiditySeconds"`
	Authorities                 []string `json:"authorities"`
	AuthorizedGrantTypes        []string `json:"authorizedGrantTypes"`
	Scopes                      []string `json:"scopes"`
}

// ListOauthClientsResult structure parses the list oauth client response payload
type ListOauthClientsResult struct {
	OauthClients *[]OauthClient `json:"clients"`
	Meta         *MetaResult    `json:"meta"`
}

type GetOauthClientResult struct {
	OauthClient *OauthClient `json:"client"`
}

type CreateOauthClientResult struct {
	Success     bool              `json:"success"`
	Message     string            `json:"msg"`
	Errors      map[string]string `json:"errors"`
	OauthClient *OauthClient      `json:"client"`
}

type UpdateOauthClientResult struct {
	CreateOauthClientResult
}

type DeleteOauthClientResult struct {
	DeleteResult
}

// ListOauthClients lists all oauth client
func (client *Client) ListOauthClients(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        OauthClientsPath,
		QueryParams: req.QueryParams,
		Result:      &ListOauthClientsResult{},
	})
}

// GetOauthClient gets an existing oauth client
func (client *Client) GetOauthClient(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", OauthClientsPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetOauthClientResult{},
	})
}

// CreateOauthClient creates a new oauth client
func (client *Client) CreateOauthClient(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        OauthClientsPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateOauthClientResult{},
	})
}

// UpdateOauthClient updates an existing oauth client
func (client *Client) UpdateOauthClient(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", OauthClientsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateOauthClientResult{},
	})
}

// DeleteOauthClient deletes an existing oauth client
func (client *Client) DeleteOauthClient(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", OauthClientsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteOauthClientResult{},
	})
}

// FindOauthClientByName gets an existing oauth client by name
func (client *Client) FindOauthClientByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListOauthClients(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListOauthClientsResult)
	clientsCount := len(*listResult.OauthClients)
	if clientsCount != 1 {
		return resp, fmt.Errorf("found %d oAuth Clients for %v", clientsCount, name)
	}
	firstRecord := (*listResult.OauthClients)[0]
	clientID := firstRecord.ID
	return client.GetOauthClient(clientID, &Request{})
}
