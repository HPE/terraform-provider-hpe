package morpheus

import (
	"fmt"
)

var (
	// WikisPath is the API endpoint for wikis
	WikisPath = "/api/wiki/pages"
	// WikiCategoriesPath is the API endpoint for wiki categories
	WikiCategoriesPath = "/api/wiki/categories"
)

// Wiki structures for use in request and response payloads
type Wiki struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	UrlName   string `json:"urlName"`
	Category  string `json:"category"`
	RefID     int64  `json:"refId"`
	RefType   string `json:"refType"`
	Format    string `json:"format"`
	Content   string `json:"content"`
	CreatedBy struct {
		ID       int64  `json:"id"`
		Username string `json:"username"`
	}
	UpdatedBy struct {
		ID       int64  `json:"id"`
		Username string `json:"username"`
	}
	DateCreated string `json:"dateCreated"`
	LastUpdated string `json:"lastUpdated"`
}

// CreateWiki structure defines the create wiki payload
type CreateWiki struct {
	Name     string `json:"name"`
	Category string `json:"category"`
	Content  string `json:"content"`
}

// WikiCategory structure defines the wiki category payload
type WikiCategory struct {
	Name      string `json:"name"`
	PageCount int64  `json:"pageCount"`
}

type ListWikisResult struct {
	Wikis *[]Wiki     `json:"pages"`
	Meta  *MetaResult `json:"meta"`
}

type ListWikiCategoriesResult struct {
	WikiCategories *[]WikiCategory `json:"categories"`
	Meta           *MetaResult     `json:"meta"`
}

type GetWikiResult struct {
	Wiki *Wiki `json:"page"`
}

type CreateWikiResult struct {
	Success bool              `json:"success"`
	Message string            `json:"msg"`
	Errors  map[string]string `json:"errors"`
	Wiki    *Wiki             `json:"page"`
}

type UpdateWikiResult struct {
	CreateWikiResult
}

type DeleteWikiResult struct {
	DeleteResult
}

// Client request methods

// ListWikis lists all the wikis
// https://apidocs.morpheusdata.com/#get-all-wiki-pages
func (client *Client) ListWikis(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        WikisPath,
		QueryParams: req.QueryParams,
		Result:      &ListWikisResult{},
	})
}

// ListWikiCategories lists all wiki categories
// https://apidocs.morpheusdata.com/#get-all-wiki-categories
func (client *Client) ListWikiCategories(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        WikiCategoriesPath,
		QueryParams: req.QueryParams,
		Result:      &ListWikiCategoriesResult{},
	})
}

// GetAppWiki gets a wiki
// https://apidocs.morpheusdata.com/#get-a-specific-wiki-page
func (client *Client) GetWiki(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", WikisPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetWikiResult{},
	})
}

// GetInstanceWiki gets an instance wiki
// https://apidocs.morpheusdata.com/#get-a-wiki-page-for-instance
func (client *Client) GetInstanceWiki(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("/api/instances/%d/wiki", id),
		QueryParams: req.QueryParams,
		Result:      &GetWikiResult{},
	})
}

// GetAppWiki gets an app wiki
// https://apidocs.morpheusdata.com/#get-a-wiki-page-for-app
func (client *Client) GetAppWiki(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("/api/apps/%d/wiki", id),
		QueryParams: req.QueryParams,
		Result:      &GetWikiResult{},
	})
}

// GetClusterWiki gets a cluster wiki
// https://apidocs.morpheusdata.com/#get-a-wiki-page-for-cluster
func (client *Client) GetClusterWiki(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("/api/clusters/%d/wiki", id),
		QueryParams: req.QueryParams,
		Result:      &GetWikiResult{},
	})
}

// GetServerWiki gets a server wiki
// https://apidocs.morpheusdata.com/#get-a-wiki-page-for-server
func (client *Client) GetServerWiki(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("/api/servers/%d/wiki", id),
		QueryParams: req.QueryParams,
		Result:      &GetWikiResult{},
	})
}

// GetCloudWiki gets a cloud wiki
// https://apidocs.morpheusdata.com/#get-a-wiki-page-for-cloud
func (client *Client) GetCloudWiki(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("/api/zones/%d/wiki", id),
		QueryParams: req.QueryParams,
		Result:      &GetWikiResult{},
	})
}

// GetGroupWiki gets a wiki
// https://apidocs.morpheusdata.com/#get-a-wiki-page-for-group
func (client *Client) GetGroupWiki(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("/api/groups/%d/wiki", id),
		QueryParams: req.QueryParams,
		Result:      &GetWikiResult{},
	})
}

// CreateWiki creates a new wiki
// https://apidocs.morpheusdata.com/#create-a-wiki-page
func (client *Client) CreateWiki(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        WikisPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateWikiResult{},
	})
}

// UpdateWiki updates an existing wiki
// https://apidocs.morpheusdata.com/#update-a-wiki-page
func (client *Client) UpdateWiki(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", WikisPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateWikiResult{},
	})
}

// UpdateInstanceWiki updates an existing instance wiki
// https://apidocs.morpheusdata.com/#update-a-wiki-page-for-instance
func (client *Client) UpdateInstanceWiki(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("/api/instances/%d/wiki", id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateWikiResult{},
	})
}

// UpdateAppWiki updates an existing app wiki
// https://apidocs.morpheusdata.com/#update-a-wiki-page-for-app
func (client *Client) UpdateAppWiki(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("/api/apps/%d/wiki", id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateWikiResult{},
	})
}

// https://apidocs.morpheusdata.com/#update-a-wiki-page-for-cluster
func (client *Client) UpdateClusterWiki(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("/api/clusters/%d/wiki", id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateWikiResult{},
	})
}

// UpdateServerWiki updates an existing server wiki
// https://apidocs.morpheusdata.com/#update-a-wiki-page-for-server
func (client *Client) UpdateServerWiki(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("/api/servers/%d/wiki", id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateWikiResult{},
	})
}

// UpdateCloudWiki updates an existing cloud wiki
// https://apidocs.morpheusdata.com/#update-a-wiki-page-for-cloud
func (client *Client) UpdateCloudWiki(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("/api/zones/%d/wiki", id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateWikiResult{},
	})
}

// UpdateGroupWiki updates an existing group wiki
// https://apidocs.morpheusdata.com/#update-a-wiki-page-for-group
func (client *Client) UpdateGroupWiki(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("/api/groups/%d/wiki", id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateWikiResult{},
	})
}

// DeleteWiki deletes an existing wiki
// https://apidocs.morpheusdata.com/#delete-a-wiki-page
func (client *Client) DeleteWiki(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", WikisPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteWikiResult{},
	})
}

// FindWikiByName gets an existing wiki by name
func (client *Client) FindWikiByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListWikis(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListWikisResult)
	wikiCount := len(*listResult.Wikis)
	if wikiCount != 1 {
		return resp, fmt.Errorf("found %d Wikis for %v", wikiCount, name)
	}
	firstRecord := (*listResult.Wikis)[0]
	wikiID := firstRecord.ID
	return client.GetWiki(wikiID, &Request{})
}
