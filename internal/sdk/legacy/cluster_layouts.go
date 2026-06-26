package morpheus

import (
	"fmt"
	"time"
)

var (
	// ClusterLayoutsPath is the API endpoint for cluster layouts
	ClusterLayoutsPath = "/api/library/cluster-layouts"
)

// ClusterLayout structures for use in request and response payloads
type ClusterLayout struct {
	ID                int64     `json:"id"`
	ServerCount       int       `json:"serverCount"`
	DateCreated       time.Time `json:"dateCreated"`
	Code              string    `json:"code"`
	LastUpdated       time.Time `json:"lastUpdated"`
	HasAutoScale      bool      `json:"hasAutoScale"`
	MemoryRequirement int       `json:"memoryRequirement"`
	ComputeVersion    string    `json:"computeVersion"`
	ProvisionType     struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Code string `json:"code"`
	} `json:"provisionType"`
	Config      string `json:"config"`
	HasSettings bool   `json:"hasSettings"`
	SortOrder   int    `json:"sortOrder"`
	HasConfig   bool   `json:"hasConfig"`
	GroupType   struct {
		ID   int    `json:"id"`
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"groupType"`
	Name   string   `json:"name"`
	Labels []string `json:"labels"`
	Type   struct {
		ID   int    `json:"id"`
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"type"`
	Account struct {
		ID int `json:"id"`
	} `json:"account"`
	Creatable            bool          `json:"creatable"`
	Enabled              bool          `json:"enabled"`
	Description          string        `json:"description"`
	EnvironmentVariables []interface{} `json:"environmentVariables"`
	OptionTypes          []struct {
		ID                 int         `json:"id"`
		Name               string      `json:"name"`
		Description        interface{} `json:"description"`
		Code               string      `json:"code"`
		FieldName          string      `json:"fieldName"`
		FieldLabel         string      `json:"fieldLabel"`
		FieldCode          string      `json:"fieldCode"`
		FieldContext       string      `json:"fieldContext"`
		FieldGroup         string      `json:"fieldGroup"`
		FieldClass         interface{} `json:"fieldClass"`
		FieldAddon         interface{} `json:"fieldAddOn"`
		FieldComponent     interface{} `json:"fieldComponent"`
		FieldInput         interface{} `json:"fieldInput"`
		PlaceHolder        interface{} `json:"placeHolder"`
		VerifyPattern      interface{} `json:"verifyPattern"`
		HelpBlock          string      `json:"helpBlock"`
		HelpBlockFieldCode interface{} `json:"helpBlockFieldCode"`
		DefaultValue       string      `json:"defaultValue"`
		OptionSource       interface{} `json:"optionSource"`
		OptionSourceType   interface{} `json:"optionSourceType"`
		OptionList         interface{} `json:"optionList"`
		Type               string      `json:"type"`
		Advanced           bool        `json:"advanced"`
		Required           bool        `json:"required"`
		ExportMeta         bool        `json:"exportMeta"`
		Editable           bool        `json:"editable"`
		Creatable          bool        `json:"creatable"`
		Config             struct {
		} `json:"config"`
		DisplayOrder          int         `json:"displayOrder"`
		WrapperClass          interface{} `json:"wrapperClass"`
		Enabled               bool        `json:"enabled"`
		NoBlank               bool        `json:"noBlank"`
		DependsOnCode         interface{} `json:"dependsOnCode"`
		VisibleOnCode         interface{} `json:"visibleOnCode"`
		RequireOnCode         interface{} `json:"requireOnCode"`
		ContextualDefault     bool        `json:"contextualDefault"`
		DisplayValueOnDetails bool        `json:"displayValueOnDetails"`
		ShowOnCreate          bool        `json:"showOnCreate"`
		ShowOnEdit            bool        `json:"showOnEdit"`
		LocalCredential       interface{} `json:"localCredential"`
	} `json:"optionTypes"`
	Actions        []interface{} `json:"actions"`
	ComputeServers []struct {
		ID                      int         `json:"id"`
		PriorityOrder           int         `json:"priorityOrder"`
		NodeCount               int         `json:"nodeCount"`
		NodeType                string      `json:"nodeType"`
		MinNodeCount            int         `json:"minNodeCount"`
		MaxNodeCount            interface{} `json:"maxNodeCount"`
		DynamicCount            bool        `json:"dynamicCount"`
		InstallContainerRuntime bool        `json:"installContainerRuntime"`
		InstallStorageRuntime   bool        `json:"installStorageRuntime"`
		Name                    string      `json:"name"`
		Code                    string      `json:"code"`
		Category                interface{} `json:"category"`
		Config                  interface{} `json:"config"`
		ContainerType           struct {
			ID               int         `json:"id"`
			Account          interface{} `json:"account"`
			Name             string      `json:"name"`
			Shortname        string      `json:"shortName"`
			Code             string      `json:"code"`
			ContainerVersion string      `json:"containerVersion"`
			ProvisionType    struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
				Code string `json:"code"`
			} `json:"provisionType"`
			VirtualImage struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"virtualImage"`
			Category string `json:"category"`
			Config   struct {
			} `json:"config"`
			Containerports []struct {
				ID                  int         `json:"id"`
				Name                string      `json:"name"`
				Port                int         `json:"port"`
				LoadBalanceProtocol interface{} `json:"loadBalanceProtocol"`
				ExportName          string      `json:"exportName"`
			} `json:"containerPorts"`
			ContainerScripts []struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"containerScripts"`
			ContainerTemplates   []interface{} `json:"containerTemplates"`
			EnvironmentVariables []interface{} `json:"environmentVariables"`
		} `json:"containerType"`
		Computeservertype struct {
			ID             int    `json:"id"`
			Code           string `json:"code"`
			Name           string `json:"name"`
			Managed        bool   `json:"managed"`
			ExternalDelete bool   `json:"externalDelete"`
		} `json:"computeServerType"`
		ProvisionService interface{} `json:"provisionService"`
		PlanCategory     interface{} `json:"planCategory"`
		NamePrefix       interface{} `json:"namePrefix"`
		NameSuffix       string      `json:"nameSuffix"`
		ForceNameIndex   bool        `json:"forceNameIndex"`
		LoadBalance      bool        `json:"loadBalance"`
	} `json:"computeServers"`
	InstallContainerRuntime bool `json:"installContainerRuntime"`
	SpecTemplates           []struct {
		ID      int `json:"id"`
		Account struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"account"`
		Name string      `json:"name"`
		Code interface{} `json:"code"`
		Type struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			Code string `json:"code"`
		} `json:"type"`
		ExternalID   interface{} `json:"externalId"`
		ExternalType interface{} `json:"externalType"`
		DeploymentID interface{} `json:"deploymentId"`
		Status       interface{} `json:"status"`
		File         struct {
			ID          int         `json:"id"`
			SourceType  string      `json:"sourceType"`
			ContentRef  interface{} `json:"contentRef"`
			ContentPath interface{} `json:"contentPath"`
			Repository  interface{} `json:"repository"`
			Content     string      `json:"content"`
		} `json:"file"`
		Config struct {
			CloudFormation struct {
				Iam                  string `json:"IAM"`
				CapabilityAutoExpand string `json:"CAPABILITY_AUTO_EXPAND"`
				CapabilityNamedIam   string `json:"CAPABILITY_NAMED_IAM"`
			} `json:"cloudformation"`
		} `json:"config"`
		CreatedBy   string      `json:"createdBy"`
		UpdatedBy   interface{} `json:"updatedBy"`
		DateCreated time.Time   `json:"dateCreated"`
		LastUpdated time.Time   `json:"lastUpdated"`
	} `json:"specTemplates"`
	TaskSets []struct {
		ID   int         `json:"id"`
		Code interface{} `json:"code"`
		Name string      `json:"name"`
	} `json:"taskSets"`
}

type ListClusterLayoutsResult struct {
	ClusterLayouts *[]ClusterLayout `json:"layouts"`
	Meta           *MetaResult      `json:"meta"`
}

type GetClusterLayoutResult struct {
	ClusterLayout *ClusterLayout `json:"layout"`
}

type CreateClusterLayoutResult struct {
	Success       bool              `json:"success"`
	Message       string            `json:"msg"`
	Errors        map[string]string `json:"errors"`
	ClusterLayout *ClusterLayout    `json:"layout"`
}

type UpdateClusterLayoutResult struct {
	CreateClusterLayoutResult
}

type DeleteClusterLayoutResult struct {
	DeleteResult
}

// Client request methods
func (client *Client) ListClusterLayouts(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        ClusterLayoutsPath,
		QueryParams: req.QueryParams,
		Result:      &ListClusterLayoutsResult{},
	})
}

func (client *Client) GetClusterLayout(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", ClusterLayoutsPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetClusterLayoutResult{},
	})
}

func (client *Client) CreateClusterLayout(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        ClusterLayoutsPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateClusterLayoutResult{},
	})
}

func (client *Client) UpdateClusterLayout(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", ClusterLayoutsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateClusterLayoutResult{},
	})
}

func (client *Client) DeleteClusterLayout(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", ClusterLayoutsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteClusterLayoutResult{},
	})
}

// FindClusterLayoutByName gets an existing cluster layout by name
func (client *Client) FindClusterLayoutByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListClusterLayouts(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListClusterLayoutsResult)
	clusterLayoutCount := len(*listResult.ClusterLayouts)
	if clusterLayoutCount != 1 {
		return resp, fmt.Errorf("found %d Cluster Layouts for %v", clusterLayoutCount, name)
	}
	firstRecord := (*listResult.ClusterLayouts)[0]
	clusterLayoutID := firstRecord.ID
	return client.GetClusterLayout(clusterLayoutID, &Request{})
}
