package morpheus

import (
	"fmt"
)

var (
	// BlueprintsPath is the API endpoint for blueprints
	BlueprintsPath = "/api/blueprints"
)

// Blueprint structures for use in request and response payloads
type Blueprint struct {
	ID          int64    `json:"id"`
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Labels      []string `json:"labels"`
	Category    string   `json:"category"`
	Visibility  string   `json:"visibility"`
	Config      struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Arm         struct {
			ConfigType       string `json:"configType"`
			OsType           string `json:"osType"`
			CloudInitEnabled bool   `json:"cloudInitEnabled"`
			InstallAgent     bool   `json:"installAgent"`
			JSON             string `json:"json"`
			Git              struct {
				Path          string `json:"path"`
				RepoId        int64  `json:"repoId"`
				IntegrationId int64  `json:"integrationId"`
				Branch        string `json:"branch"`
			} `json:"git"`
		} `json:"arm"`
		CloudFormation struct {
			ConfigType       string `json:"configType"`
			CloudInitEnabled bool   `json:"cloudInitEnabled"`
			InstallAgent     bool   `json:"installAgent"`
			JSON             string `json:"json"`
			YAML             string `json:"yaml"`
			IAM              bool   `json:"IAM"`
			IAMNamed         bool   `json:"CAPABILITY_NAMED_IAM"`
			AutoExpand       bool   `json:"CAPABILITY_AUTO_EXPAND"`
			Git              struct {
				Path          string `json:"path"`
				RepoId        int64  `json:"repoId"`
				IntegrationId int64  `json:"integrationId"`
				Branch        string `json:"branch"`
			} `json:"git"`
		} `json:"cloudformation"`
		Helm struct {
			ConfigType string `json:"configType"`
			Git        struct {
				Path          string `json:"path"`
				RepoId        int    `json:"repoId"`
				IntegrationId int    `json:"integrationId"`
				Branch        string `json:"branch"`
			} `json:"git"`
		} `json:"helm"`
		Kubernetes struct {
			ConfigType string `json:"configType"`
			Git        struct {
				Path          string `json:"path"`
				RepoId        int    `json:"repoId"`
				IntegrationId int    `json:"integrationId"`
				Branch        string `json:"branch"`
			} `json:"git"`
		} `json:"kubernetes"`
		Terraform struct {
			TfVersion      string `json:"tfVersion"`
			Tf             string `json:"tf"`
			TfVarSecret    string `json:"tfvarSecret"`
			CommandOptions string `json:"commandOptions"`
			ConfigType     string `json:"configType"`
			JSON           string `json:"json"`
			Git            struct {
				Path          string `json:"path"`
				RepoId        int64  `json:"repoId"`
				IntegrationId int64  `json:"integrationId"`
				Branch        string `json:"branch"`
			} `json:"git"`
		} `json:"terraform"`
		Config struct {
			Specs []struct {
				ID    int64  `json:"id"`
				Value string `json:"value"`
				Name  string `json:"name"`
			} `json:"specs"`
		} `json:"config"`
		Type     string `json:"type"`
		Category string `json:"category"`
		Image    string `json:"image"`
	} `json:"config"`
	ResourcePermission struct {
		All      bool          `json:"all"`
		Sites    []interface{} `json:"sites"`
		AllPlans bool          `json:"allPlans"`
		Plans    []interface{} `json:"plans"`
	} `json:"resourcePermission"`
	Owner struct {
		ID       int64  `json:"id"`
		Username string `json:"username"`
	} `json:"owner"`
	Tenant struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"tenant"`
}

// ListBlueprintsResult structure parses the list blueprints response payload
type ListBlueprintsResult struct {
	Blueprints *[]Blueprint `json:"blueprints"`
	Meta       *MetaResult  `json:"meta"`
}

type GetBlueprintResult struct {
	Blueprint *Blueprint `json:"blueprint"`
}

type CreateBlueprintResult struct {
	Success   bool              `json:"success"`
	Message   string            `json:"msg"`
	Errors    map[string]string `json:"errors"`
	Blueprint *Blueprint        `json:"blueprint"`
}

type UpdateBlueprintResult struct {
	CreateBlueprintResult
}

type DeleteBlueprintResult struct {
	DeleteResult
}

// Client request methods

func (client *Client) ListBlueprints(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        BlueprintsPath,
		QueryParams: req.QueryParams,
		Result:      &ListBlueprintsResult{},
	})
}

func (client *Client) GetBlueprint(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", BlueprintsPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetBlueprintResult{},
	})
}

func (client *Client) CreateBlueprint(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        BlueprintsPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateBlueprintResult{},
	})
}

// UpdateBlueprint updates an existing blueprint
func (client *Client) UpdateBlueprint(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", BlueprintsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateBlueprintResult{},
	})
}

// UpdateBlueprintLogo updates an existing blueprint logo
func (client *Client) UpdateBlueprintLogo(id int64, filePayload []*FilePayload, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:         "POST",
		Path:           fmt.Sprintf("/api/blueprints/%d/image", id),
		IsMultiPart:    true,
		MultiPartFiles: filePayload,
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
		Result: &UpdateBlueprintResult{},
	})
}

// DeleteBlueprint deletes an existing blueprint
func (client *Client) DeleteBlueprint(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", BlueprintsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteBlueprintResult{},
	})
}

func (client *Client) FindBlueprintByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListBlueprints(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListBlueprintsResult)
	blueprintsCount := len(*listResult.Blueprints)
	if blueprintsCount != 1 {
		return resp, fmt.Errorf("found %d Blueprints for %v", blueprintsCount, name)
	}
	firstRecord := (*listResult.Blueprints)[0]
	blueprintID := firstRecord.ID
	return client.GetBlueprint(blueprintID, &Request{})
}
