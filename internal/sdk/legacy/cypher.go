package morpheus

import (
	"fmt"
	"time"
)

var (
	// CypherPath is the API endpoint for cypher
	CypherPath = "/api/cypher"
)

// Cypher structures for use in request and response payloads
type Cypher struct {
	ID           int64     `json:"id"`
	ItemKey      string    `json:"itemKey"`
	LeaseTimeout int64     `json:"leaseTimeout"`
	ExpireDate   time.Time `json:"expireDate"`
	DateCreated  time.Time `json:"dateCreated"`
	LastUpdated  time.Time `json:"lastUpdated"`
	LastAccessed time.Time `json:"lastAccessed"`
	CreatedBy    string    `json:"createdBy"`
}

type CypherData struct {
	Keys []string `json:"keys"`
}

type ListCypherResult struct {
	Success bool        `json:"success"`
	Data    CypherData  `json:"data"`
	Cyphers *[]Cypher   `json:"cyphers"`
	Meta    *MetaResult `json:"meta"`
}

type GetCypherResult struct {
	Success       bool              `json:"success"`
	Data          interface{}       `json:"data"`
	Type          string            `json:"type"`
	LeaseDuration int64             `json:"lease_duration"`
	Cypher        *Cypher           `json:"cypher"`
	Message       string            `json:"msg"`
	Errors        map[string]string `json:"errors"`
}

type CreateCypherResult struct {
	Success       bool              `json:"success"`
	Data          string            `json:"data"`
	Type          string            `json:"type"`
	LeaseDuration int64             `json:"lease_duration"`
	Message       string            `json:"msg"`
	Errors        map[string]string `json:"errors"`
	Cypher        *Cypher           `json:"cypher"`
}

type DeleteCypherResult struct {
	DeleteResult
}

// Client request methods
func (client *Client) ListCyphers(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        CypherPath,
		QueryParams: req.QueryParams,
		Result:      &ListCypherResult{},
	})
}

func (client *Client) GetCypher(path string, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%s", CypherPath, path),
		QueryParams: req.QueryParams,
		Result:      &GetCypherResult{},
	})
}

func (client *Client) CreateCypher(path string, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        fmt.Sprintf("%s/%s", CypherPath, path),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateCypherResult{},
	})
}

func (client *Client) DeleteCypher(path string, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%s", CypherPath, path),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteCypherResult{},
	})
}
