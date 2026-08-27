package morpheus

import (
	"fmt"
	"time"
)

var (
	// AppsPath is the API endpoint for apps
	AppsPath = "/api/apps"
)

// App structures for use in request and response payloads
type App struct {
	ID          int64    `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Labels      []string `json:"labels"`
	Environment string   `json:"environment"`
	AccountId   int64    `json:"accountId"`
	Account     struct {
		Id   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"account"`
	Owner struct {
		Id       int64  `json:"id"`
		Username string `json:"username"`
	} `json:"owner"`
	SiteId int64 `json:"siteId"`
	Group  struct {
		Id   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"group"`
	Blueprint struct {
		Id   int64  `json:"id"`
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"blueprint"`
	Type           string    `json:"type"`
	DateCreated    time.Time `json:"dateCreated"`
	LastUpdated    time.Time `json:"lastUpdated"`
	AppContext     string    `json:"appContext"`
	Status         string    `json:"status"`
	AppStatus      string    `json:"appStatus"`
	InstanceCount  int64     `json:"instanceCount"`
	ContainerCount int64     `json:"containerCount"`
	AppTiers       []struct {
		Tier struct {
			Id   int64  `json:"id"`
			Name string `json:"name"`
		} `json:"tier"`
		AppInstances []struct {
			Config   interface{} `json:"config"`
			Instance Instance    `json:"instance"`
		} `json:"appInstances"`
		BootSequence int64 `json:"bootSequence"`
	} `json:"appTiers"`
	Instances []struct {
		Id   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"instances"`
	Stats struct {
		UsedMemory            int64   `json:"usedMemory"`
		MaxMemory             int64   `json:"maxMemory"`
		UsedStorage           int64   `json:"usedStorage"`
		MaxStorage            int64   `json:"maxStorage"`
		Running               int64   `json:"running"`
		Total                 int64   `json:"total"`
		CpuUsage              float64 `json:"cpuUsage"`
		InstanceCount         int64   `json:"instanceCount"`
		InstanceDayCount      []int64 `json:"instanceDayCount"`
		InstanceDayCountTotal int64   `json:"instanceDayCountTotal"`
	} `json:"stats"`
}

// ListAppsResult structure parses the list apps response payload
type ListAppsResult struct {
	Apps *[]App      `json:"apps"`
	Meta *MetaResult `json:"meta"`
}

type GetAppResult struct {
	App *App `json:"app"`
}

type GetAppStateResult struct {
	Success   bool       `json:"success"`
	Workloads []Workload `json:"workloads"`
	IacDrift  bool       `json:"iacDrift"`
	Specs     []Spec     `json:"specs"`
	PlanData  string     `json:"planData"`
	Input     struct {
		Variables []struct {
			Name      string      `json:"name"`
			Sensitive bool        `json:"sensitive"`
			Value     interface{} `json:"value"`
			Type      interface{} `json:"type"`
		} `json:"variables"`
		Providers []struct {
			Name string `json:"name"`
		} `json:"providers"`
	} `json:"input"`
	Output struct {
		Outputs []struct {
			Name  string `json:"name"`
			Value struct {
				Sensitive bool        `json:"sensitive"`
				Value     interface{} `json:"value"`
				Type      interface{} `json:"type"`
			} `json:"value"`
		} `json:"outputs"`
	} `json:"output"`
	StateData string `json:"stateData"`
}

type Output struct {
}

type Workload struct {
	RefType    string `json:"refType"`
	RefId      int64  `json:"refId"`
	RefName    string `json:"refName"`
	SubRefName string `json:"subRefName"`
	StateDate  string `json:"stateDate"`
	Status     string `json:"status"`
	IacDrift   bool   `json:"iacDrift"`
}

type Spec struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Template struct {
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"template"`
	Isolated bool `json:"isolated"`
}

type CreateAppResult struct {
	Success bool              `json:"success"`
	Message string            `json:"msg"`
	Errors  map[string]string `json:"errors"`
	App     *App              `json:"app"`
}

type UpdateAppResult struct {
	CreateAppResult
}

type DeleteAppResult struct {
	DeleteResult
}

type ApplyStateForAppResult struct {
	Success     bool              `json:"success"`
	Message     string            `json:"msg"`
	Errors      map[string]string `json:"errors"`
	ExecutionId string            `json:"executionId"`
}

type ValidateApplyStateForAppResult struct {
	Success     bool              `json:"success"`
	Message     string            `json:"msg"`
	Errors      map[string]string `json:"errors"`
	ExecutionId string            `json:"executionId"`
}

type PrepareToApplyAppResult struct {
	Success bool              `json:"success"`
	Message string            `json:"msg"`
	Errors  map[string]string `json:"errors"`
	Data    struct {
		Image        string `json:"image"`
		Name         string `json:"name"`
		AutoValidate bool   `json:"autoValidate"`
		Terraform    struct {
			RefreshMode string `json:"refreshMode"`
			BackendType string `json:"backendType"`
			TimeoutMode string `json:"timeoutMode"`
			ConfigType  string `json:"configType"`
		} `json:"terraform"`
		Type   string `json:"type"`
		Config struct {
			Specs []Spec `json:"specs"`
		} `json:"config"`
		BlueprintName string `json:"blueprintName"`
		Description   string `json:"description"`
		TemplateId    int64  `json:"templateId"`
		BlueprintId   int64  `json:"blueprintId"`
		Group         struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
		} `json:"group"`
	} `json:"data"`
}

// Client request methods

func (client *Client) ListApps(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        AppsPath,
		QueryParams: req.QueryParams,
		Result:      &ListAppsResult{},
	})
}

func (client *Client) GetApp(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", AppsPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetAppResult{},
	})
}

func (client *Client) GetAppState(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d/state", AppsPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetAppStateResult{},
	})
}

func (client *Client) CreateApp(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        AppsPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateAppResult{},
	})
}

func (client *Client) UpdateApp(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", AppsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateAppResult{},
	})
}

func (client *Client) AddInstanceToApp(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        fmt.Sprintf("%s/%d/add-instance", AppsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateAppResult{},
	})
}

func (client *Client) RemoveInstanceFromApp(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        fmt.Sprintf("%s/%d/remove-instance", AppsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateAppResult{},
	})
}

func (client *Client) PrepareToApplyAppState(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d/prepare-apply", AppsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &PrepareToApplyAppResult{},
	})
}

func (client *Client) ApplyAppState(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        fmt.Sprintf("%s/%d/add-instance", AppsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &ApplyStateForAppResult{},
	})
}

func (client *Client) UndoAppDelete(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d/cancel-removal", AppsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateAppResult{},
	})
}

func (client *Client) RefreshAppState(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        fmt.Sprintf("%s/%d/refresh", AppsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &StandardResult{},
	})
}

func (client *Client) DeleteApp(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", AppsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteAppResult{},
	})
}

func (client *Client) ValidateApplyStateForApp(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        fmt.Sprintf("%s/%d/validate-apply", AppsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &ValidateApplyStateForAppResult{},
	})
}

// helper functions
func (client *Client) FindAppByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListApps(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListAppsResult)
	appsCount := len(*listResult.Apps)
	if appsCount != 1 {
		return resp, fmt.Errorf("found %d Apps for %v", appsCount, name)
	}
	firstRecord := (*listResult.Apps)[0]
	appID := firstRecord.ID
	return client.GetApp(appID, &Request{})
}
