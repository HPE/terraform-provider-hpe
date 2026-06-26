package morpheus

import (
	"fmt"
)

var (
	// IntegrationsPath is the API endpoint for integrations
	IntegrationsPath = "/api/integrations"
)

// Integration structures for use in request and response payloads
type Integration struct {
	ID              int64  `json:"id"`
	Name            string `json:"name"`
	Enabled         bool   `json:"enabled"`
	Type            string `json:"type"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	PasswordHash    string `json:"passwordHash"`
	Port            string `json:"port"`
	Version         string `json:"version"`
	IntegrationType struct {
		ID   int64  `json:"id"`
		Code string `json:"code"`
		Name string `json:"name"`
	}
	URL             string `json:"url"`
	WindowsVersion  string `json:"windowsVersion"`
	RefType         string `json:"refType"`
	RefId           string `json:"refId"`
	ServiceUrl      string `json:"serviceUrl"`
	ServiceUsername string `json:"serviceUsername"`
	ServicePassword string `json:"servicePassword"`
	Token           string `json:"token"`
	TokenHash       string `json:"tokenHash"`
	AuthType        string `json:"authType"`
	AuthId          string `json:"authId"`
	Config          struct {
		Inventory string `json:"inventory"`
		// Git Integration
		DefaultBranch string      `json:"defaultBranch"`
		CacheEnabled  interface{} `json:"cacheEnabled"`
		// Ansible Integration
		AnsiblePlaybooks                string                            `json:"ansiblePlaybooks"`
		AnsibleRoles                    string                            `json:"ansibleRoles"`
		AnsibleGroupVars                string                            `json:"ansibleGroupVars"`
		AnsibleHostVars                 string                            `json:"ansibleHostVars"`
		AnsibleCommandBus               interface{}                       `json:"ansibleCommandBus"`
		AnsibleVerbose                  interface{}                       `json:"ansibleVerbose"`
		AnsibleGalaxyEnabled            interface{}                       `json:"ansibleGalaxyEnabled"`
		AnsibleDefaultBranch            string                            `json:"ansibleDefaultBranch"`
		Plugin                          interface{}                       `json:"plugin"`
		IncidentAccess                  bool                              `json:"incidentAccess"`
		RequestAccess                   bool                              `json:"requestAccess"`
		ServiceNowCMDBMode              string                            `json:"cmdbMode"`
		ServiceNowIgnoreCertErrors      bool                              `json:"ignoreCertErrors"`
		ServiceNowCMDBBusinessObject    string                            `json:"serviceNowCMDBBusinessObject"`
		ServiceNowCustomCmdbMapping     string                            `json:"serviceNowCustomCmdbMapping"`
		ServiceNowCmdbClassMapping      []serviceNowCmdbClassMappingEntry `json:"serviceNowCmdbClassMapping"`
		ServiceNowCmdbClassMappingInput []string                          `json:"serviceNowCmdbClassMapping.input"`
		PreparedForSync                 bool                              `json:"preparedForSync"`
		PreparedForSyncv2               bool                              `json:"preparedForSync-v2"`
		PreparedForSyncv3               bool                              `json:"preparedForSync-v3"`
		PreparedForSyncv4               bool                              `json:"preparedForSync-v4"`
		WebServiceImportUrl             string                            `json:"webServiceImportUrl"`
		WebServiceImportSysId           string                            `json:"webServiceImportSysId"`
		WebServiceOperationUrl          string                            `json:"webServiceOperationUrl"`
		ApprovalUser                    string                            `json:"approvalUser"`
		Company                         string                            `json:"company"`
		AppID                           string                            `json:"appId"`
		InventoryExisting               string                            `json:"inventoryExisting"`
		ExtraAttributes                 string                            `json:"extraAttributes"`
		EngineMount                     string                            `json:"engineMount"`
		SecretPath                      string                            `json:"secretPath"`
		SecretEngine                    string                            `json:"secretEngine"`
		SecretPathHash                  string                            `json:"secretPathHash"`
		SecretEngineHash                string                            `json:"secretEngineHash"`
		// Chef Integration
		Databags              []chefDatabagEntry `json:"databag"`
		Endpoint              string             `json:"endpoint"`
		Org                   string             `json:"org"`
		ChefUser              string             `json:"chefUser"`
		UseFqdn               bool               `json:"useFqdn"`
		ChefUseFQDN           bool               `json:"chefUseFqdn"`
		UserKey               string             `json:"userKey"`
		OrgKey                string             `json:"orgKey"`
		Version               string             `json:"version"`
		WindowsVersion        string             `json:"windowsVersion"`
		WindowsInstallURL     string             `json:"windowsInstallUrl"`
		UserKeyHash           string             `json:"userKeyHash"`
		OrgKeyHash            string             `json:"orgKeyHash"`
		PuppetMaster          string             `json:"puppetMaster"`
		PuppetFireNow         string             `json:"puppetFireNow"`
		PuppetSshUser         string             `json:"puppetSshUser"`
		PuppetSshPassword     string             `json:"puppetSshPassword"`
		PuppetSshPasswordHash string             `json:"puppetSshPasswordHash"`
		// Cherwell integration
		CherwellCustomCmdbMapping string `json:"cherwellCustomCmdbMapping"`
		CherwellClientKey         string `json:"cherwellClientKey"`
		CherwellCreatedBy         string `json:"cherwellCreatedBy"`
		CherwellStartDate         string `json:"cherwellStartDate"`
		CherwellEndDate           string `json:"cherwellEndDate"`
		CherwellIgnoreSSLErrors   string `json:"cherwellIgnoreSSLErrors"`
		CherwellBusinessObject    string `json:"cherwellBusinessObject"`
		ZoneFilter                string `json:"zoneFilter"`
	}
	Status           string `json:"status"`
	StatusDate       string `json:"statusDate"`
	StatusMessage    string `json:"statusMessage"`
	LastSync         string `json:"lastSync"`
	LastSyncDuration string `json:"lastSyncDuration"`
	IsPlugin         bool   `json:"isPlugin"`
	ServiceKey       struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	}
	ServiceMode string `json:"serviceMode"`
	ServiceFlag bool   `json:"serviceFlag"`
	Credential  struct {
		ID    int64    `json:"id"`
		Name  string   `json:"name"`
		Type  string   `json:"type"`
		Types []string `json:"types"`
	}
}

type serviceNowCmdbClassMappingEntry struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Code     string `json:"code"`
	NowClass string `json:"nowClass"`
}

type chefDatabagEntry struct {
	Path string `json:"path"`
	Key  string `json:"key"`
}

// ListIntegrationsResult structure parses the list integrations response payload
type ListIntegrationsResult struct {
	Integrations *[]Integration `json:"integrations"`
	Meta         *MetaResult    `json:"meta"`
}

type ListIntegrationObjectsResult struct {
	Objects []struct {
		ID              int64  `json:"id"`
		Name            string `json:"name"`
		Type            string `json:"type"`
		RefType         string `json:"refType"`
		RefID           int64  `json:"refId"`
		CatalogItemType struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
		}
	} `json:"objects"`
	Meta *MetaResult `json:"meta"`
}

type GetIntegrationResult struct {
	Integration *Integration `json:"integration"`
}

type CreateIntegrationResult struct {
	Success     bool              `json:"success"`
	Message     string            `json:"msg"`
	Errors      map[string]string `json:"errors"`
	Integration *Integration      `json:"integration"`
}

type CreateIntegrationObjectResult struct {
	Success bool              `json:"success"`
	Message string            `json:"msg"`
	Errors  map[string]string `json:"errors"`
	Object  struct {
		ID int64 `json:"id"`
	} `json:"object"`
}

type UpdateIntegrationResult struct {
	CreateIntegrationResult
}

type DeleteIntegrationResult struct {
	DeleteResult
}

// Client request methods
func (client *Client) ListIntegrations(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        IntegrationsPath,
		QueryParams: req.QueryParams,
		Result:      &ListIntegrationsResult{},
	})
}

func (client *Client) ListIntegrationObjects(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d/objects", IntegrationsPath, id),
		QueryParams: req.QueryParams,
		Result:      &ListIntegrationsResult{},
	})
}

func (client *Client) GetIntegration(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", IntegrationsPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetIntegrationResult{},
	})
}

func (client *Client) CreateIntegration(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        IntegrationsPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateIntegrationResult{},
	})
}

func (client *Client) RefreshIntegration(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        fmt.Sprintf("%s/%d/refresh", IntegrationsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateIntegrationResult{},
	})
}

func (client *Client) CreateIntegrationObject(id int64, objectId int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        fmt.Sprintf("%s/%d/objects/%d", IntegrationsPath, id, objectId),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateIntegrationResult{},
	})
}

func (client *Client) UpdateIntegration(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", IntegrationsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateIntegrationResult{},
	})
}

func (client *Client) DeleteIntegration(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", IntegrationsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteIntegrationResult{},
	})
}

func (client *Client) DeleteIntegrationObject(id int64, objectId int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d/objects/%d", IntegrationsPath, id, objectId),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteIntegrationResult{},
	})
}

// helper functions

func (client *Client) FindIntegrationByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListIntegrations(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListIntegrationsResult)
	integrationCount := len(*listResult.Integrations)
	if integrationCount != 1 {
		return resp, fmt.Errorf("found %d integrations named %v", integrationCount, name)
	}
	firstRecord := (*listResult.Integrations)[0]
	integrationID := firstRecord.ID
	return client.GetIntegration(integrationID, &Request{})
}
