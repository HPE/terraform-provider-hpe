package morpheus

import (
	"fmt"
)

var (
	// KeyPairsPath is the API endpoint for key pairs
	KeyPairsPath = "/api/key-pairs"
)

// KeyPair structures for use in request and response payloads
type KeyPair struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	AccountId      int64  `json:"accountId"`
	PublicKey      string `json:"publicKey"`
	HasPrivateKey  bool   `json:"hasPrivateKey"`
	PrivateKeyHash string `json:"privateKeyHash"`
	Fingerprint    string `json:"fingerprint"`
	PrivateKey     string `json:"privateKey"`
	DateCreated    string `json:"dateCreated"`
	LastUpdated    string `json:"lastUpdated"`
}

// ListKeyPairsResult structure parses the list key pairs response payload
type ListKeyPairsResult struct {
	KeyPairs *[]KeyPair  `json:"keyPairs"`
	Meta     *MetaResult `json:"meta"`
}

// GetKeyPairResult structure parses the get key pair response payload
type GetKeyPairResult struct {
	KeyPair *KeyPair `json:"keyPair"`
}

type CreateKeyPairResult struct {
	Success bool              `json:"success"`
	Message string            `json:"msg"`
	Errors  map[string]string `json:"errors"`
	KeyPair *KeyPair          `json:"keyPair"`
}

type UpdateKeyPairResult struct {
	CreateKeyPairResult
}

type DeleteKeyPairResult struct {
	DeleteResult
}

func (client *Client) ListKeyPairs() (*Response, error) {
	return client.Execute(&Request{
		Method: "GET",
		Path:   KeyPairsPath,
		Result: &ListKeyPairsResult{},
	})
}

func (client *Client) GetKeyPair(id int64) (*Response, error) {
	return client.Execute(&Request{
		Method: "GET",
		Path:   fmt.Sprintf("%s/%d", KeyPairsPath, id),
		Result: &GetKeyPairResult{},
	})
}
func (client *Client) CreateKeyPair(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        KeyPairsPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateKeyPairResult{},
	})
}
func (client *Client) DeleteKeyPair(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", KeyPairsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteKeyPairResult{},
	})
}
func (client *Client) GetKeyPairByName(name string) (*Response, error) {
	return client.Execute(&Request{
		Method: "GET",
		QueryParams: map[string]string{
			"name": name,
		},
		Path:   KeyPairsPath,
		Result: &ListKeyPairsResult{},
	})
}
