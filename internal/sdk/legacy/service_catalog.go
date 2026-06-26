package morpheus

import (
	"fmt"
)

var (
	ServiceCatalogPath      = "/api/catalog/cart"
	ServiceCatalogItemsPath = "/api/catalog/cart/items"
	ServiceCatalogTypesPath = "/api/catalog/types"
)

type CatalogItemType struct {
	Id            int64        `json:"id"`
	Name          string       `json:"name"`
	Description   string       `json:"description"`
	Type          string       `json:"type"`
	Featured      bool         `json:"featured"`
	AllowQuantity bool         `json:"allowQuantity"`
	ImagePath     string       `json:"imagePath"`
	DarkImagePath string       `json:"darkImagePath"`
	OptionTypes   []OptionType `json:"optionTypes"`
}

type Cart struct {
	Id    int64      `json:"id"`
	Name  string     `json:"name"`
	Items []CartItem `json:"items"`
	Stats struct {
		Price    float64 `json:"price"`
		Currency string  `json:"currency"`
		Unit     string  `json:"unit"`
	}
}

type CartItem struct {
	Id   int64 `json:"id"`
	Type struct {
		Id   int64  `json:"id"`
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"type"`
	Quantity    int64   `json:"quantity"`
	Price       float64 `json:"price"`
	Currency    string  `json:"currency"`
	Unit        string  `json:"unit"`
	Valid       bool    `json:"valid"`
	Status      string  `json:"status"`
	DateCreated string  `json:"dateCreated"`
	LastUpdated string  `json:"lastUpdated"`
}

type CatalogOrder struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Items []struct {
		ID   int64 `json:"id"`
		Type struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
			Type string `json:"type"`
		} `json:"type"`
		Quantity    int64   `json:"quantity"`
		Price       float64 `json:"price"`
		Currency    string  `json:"currency"`
		Unit        string  `json:"unit"`
		Valid       bool    `json:"valid"`
		Status      string  `json:"status"`
		DateCreated string  `json:"dateCreated"`
		LastUpdated string  `json:"lastUpdated"`
	} `json:"items"`
	Stats struct {
		Price    int64  `json:"price"`
		Currency string `json:"currency"`
		Unit     string `json:"unit"`
	} `json:"stats"`
}

type InventoryItem struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	Quantity      int64  `json:"quantity"`
	Status        string `json:"status"`
	StatusMessage string `json:"statusMessage"`
	RefType       string `json:"refType"`
	Execution     struct {
		ID        int64  `json:"id"`
		JobID     int64  `json:"jobId"`
		Status    string `json:"status"`
		StartDate string `json:"startDate"`
		EndDate   string `json:"endDate"`
		Duration  int64  `json:"duration"`
	} `json:"execution"`
	App struct {
		ID        int64  `json:"id"`
		Name      string `json:"name"`
		Status    string `json:"status"`
		Instances []struct {
			ID              int64    `json:"id"`
			Name            string   `json:"name"`
			Status          string   `json:"status"`
			Locations       []string `json:"locations"`
			Virtualmachines int64    `json:"virtualMachines"`
			Version         string   `json:"version"`
		} `json:"instances"`
	} `json:"app"`
	Instance struct {
		ID              int64    `json:"id"`
		Name            string   `json:"name"`
		Status          string   `json:"status"`
		Locations       []string `json:"locations"`
		Virtualmachines int64    `json:"virtualMachines"`
		Version         string   `json:"version"`
	} `json:"instance"`
	OrderDate   string `json:"orderDate"`
	DateCreated string `json:"dateCreated"`
	LastUpdated string `json:"lastUpdated"`
}

// GetCatalogCartResult structure parses the get catalog response payload
type GetCatalogCartResult struct {
	Cart *Cart `json:"cart"`
}

// ClearCatalogResult structure parses the clear catalog response payload
type ClearCatalogResult struct {
	DeleteResult
}

// AddCatalogItemCartResult structure parses the add catalog item response payload
type AddCatalogItemCartResult struct {
	Success bool      `json:"success"`
	Item    *CartItem `json:"item"`
}

// RemoveCatalogItemCartResult structure parses the remove catalog item response payload
type RemoveCatalogItemCartResult struct {
	DeleteResult
}

type ListCatalogInventoryItemsResult struct {
	CatalogInventoryItems *[]InventoryItem `json:"items"`
	Meta                  *MetaResult      `json:"meta"`
}

type GetCatalogInventoryItemResult struct {
	CatalogInventoryItem *InventoryItem `json:"item"`
	Meta                 *MetaResult    `json:"meta"`
}

// DeleteServiceCatalogInventoryItemResult structure parses the delete catalog inventory item response payload
type DeleteCatalogInventoryItemResult struct {
	DeleteResult
}

// PlaceCatalogOrderResult structure parses the place catalog order response payload
type PlaceCatalogOrderResult struct {
	Success    bool                   `json:"success"`
	Msg        string                 `json:"msg"`
	Errors     map[string]interface{} `json:"errors"`
	ItemErrors interface{}            `json:"itemErrors"`
	Order      CatalogOrder           `json:"order"`
}

// ListCatalogTypesResult structure parses the list catalog item types response payload
type ListCatalogItemTypesResult struct {
	CatalogItemTypes *[]CatalogItemType `json:"catalogItemTypes"`
	Meta             *MetaResult        `json:"meta"`
}

// GetCatalogItemTypeResult structure parses the get catalog item type response payload
type GetCatalogItemTypeResult struct {
	CatalogItemType *CatalogItemType `json:"catalogItemType"`
}

// GetCatalogCart gets the contents of the cart
func (client *Client) GetCatalogCart(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        ServiceCatalogPath,
		QueryParams: req.QueryParams,
		Result:      &GetCatalogCartResult{},
	})
}

// ClearCatalogCart removes the contents of the cart
func (client *Client) ClearCatalogCart(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        ServiceCatalogPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &ClearCatalogResult{},
	})
}

// AddCatalogItemCart adds an item to the cart
func (client *Client) AddCatalogItemCart(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        ServiceCatalogItemsPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &AddCatalogItemCartResult{},
	})
}

// RemoveCatalogItemCart removes an existing item from the cart
func (client *Client) RemoveCatalogItemCart(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", ServiceCatalogItemsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &RemoveCatalogItemCartResult{},
	})
}

// ListCatalogInventoryItems list existing inventory items
func (client *Client) ListCatalogInventoryItems(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        "/api/catalog/items",
		QueryParams: req.QueryParams,
		Result:      &ListCatalogInventoryItemsResult{},
	})
}

// GetCatalogInventoryItem gets an existing inventory item
func (client *Client) GetCatalogInventoryItem(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", ServiceCatalogTypesPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetCatalogInventoryItemResult{},
	})
}

// DeleteCatalogInventoryItem deletes an existing inventory item
func (client *Client) DeleteCatalogInventoryItem(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", ServiceCatalogItemsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteCatalogInventoryItemResult{},
	})
}

// PlaceCatalogOrder places a catalog order
func (client *Client) PlaceCatalogOrder(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        "/api/catalog/orders",
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &PlaceCatalogOrderResult{},
	})
}

// ListCatalogItemTypes lists all catalog item types
func (client *Client) ListCatalogItemTypes(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        ServiceCatalogTypesPath,
		QueryParams: req.QueryParams,
		Result:      &ListCatalogItemTypesResult{},
	})
}

// GetCatalogItemType gets an existing catalog item type
func (client *Client) GetCatalogItemType(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", ServiceCatalogTypesPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetCatalogItemTypeResult{},
	})
}
