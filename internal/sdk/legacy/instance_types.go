package morpheus

import (
	"fmt"
)

var (
	// InstanceTypesPath is the API endpoint for instance types
	InstanceTypesPath = "/api/library/instance-types"
)

// InstanceType structures for use in request and response payloads
type InstanceType struct {
	ID      int64 `json:"id"`
	Account struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"account"`
	Name                string      `json:"name"`
	Labels              []string    `json:"labels"`
	Code                string      `json:"code"`
	Description         string      `json:"description"`
	ProvisionTypeCode   string      `json:"provisionTypeCode"`
	Category            string      `json:"category"`
	Active              bool        `json:"active"`
	HasProvisioningStep bool        `json:"hasProvisioningStep"`
	HasDeployment       bool        `json:"hasDeployment"`
	HasConfig           bool        `json:"hasConfig"`
	HasSettings         bool        `json:"hasSettings"`
	HasAutoscale        bool        `json:"hasAutoScale"`
	ProxyType           interface{} `json:"proxyType"`
	ProxyPort           interface{} `json:"proxyPort"`
	ProxyProtocol       interface{} `json:"proxyProtocol"`
	EnvironmentPrefix   string      `json:"environmentPrefix"`
	BackupType          interface{} `json:"backupType"`
	Config              struct {
	} `json:"config"`
	Visibility          string   `json:"visibility"`
	Featured            bool     `json:"featured"`
	Versions            []string `json:"versions"`
	InstanceTypeLayouts []struct {
		ID           int64 `json:"id"`
		InstanceType struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
			Code string `json:"code"`
		} `json:"instanceType"`
		Account struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
		} `json:"account"`
		Code                     string      `json:"code"`
		Name                     string      `json:"name"`
		InstanceVersion          string      `json:"instanceVersion"`
		Description              interface{} `json:"description"`
		Creatable                bool        `json:"creatable"`
		MemoryRequirement        interface{} `json:"memoryRequirement"`
		SortOrder                int64       `json:"sortOrder"`
		SupportsConvertToManaged bool        `json:"supportsConvertToManaged"`
	} `json:"instanceTypeLayouts"`
	OptionTypes []struct {
		ID                 int64       `json:"id"`
		Name               string      `json:"name"`
		Description        interface{} `json:"description"`
		Code               string      `json:"code"`
		FieldName          string      `json:"fieldName"`
		FieldLabel         string      `json:"fieldLabel"`
		FieldCode          interface{} `json:"fieldCode"`
		FieldContext       string      `json:"fieldContext"`
		FieldGroup         interface{} `json:"fieldGroup"`
		FieldClass         interface{} `json:"fieldClass"`
		FieldAddon         interface{} `json:"fieldAddOn"`
		FieldComponent     interface{} `json:"fieldComponent"`
		FieldInput         interface{} `json:"fieldInput"`
		Placeholder        interface{} `json:"placeHolder"`
		VerifyPattern      interface{} `json:"verifyPattern"`
		HelpBlock          interface{} `json:"helpBlock"`
		HelpBlockFieldCode interface{} `json:"helpBlockFieldCode"`
		DefaultValue       interface{} `json:"defaultValue"`
		OptionSource       string      `json:"optionSource"`
		OptionSourceType   interface{} `json:"optionSourceType"`
		OptionList         struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
		} `json:"optionList"`
		Type       string `json:"type"`
		Advanced   bool   `json:"advanced"`
		Required   bool   `json:"required"`
		ExportMeta bool   `json:"exportMeta"`
		Editable   bool   `json:"editable"`
		Creatable  bool   `json:"creatable"`
		Config     struct {
		} `json:"config"`
		DisplayOrder          int64       `json:"displayOrder"`
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
	EnvironmentVariables []struct {
		EvarName         string `json:"evarName"`
		Name             string `json:"name"`
		DefaultValue     string `json:"defaultValue"`
		DefaultValueHash string `json:"defaultValueHash"`
		ValueType        string `json:"valueType"`
		Export           bool   `json:"export"`
		Masked           bool   `json:"masked"`
	} `json:"environmentVariables"`
	ImagePath     string `json:"imagePath"`
	DarkImagePath string `json:"darkImagePath"`
	PriceSets     []struct {
		ID        int64  `json:"id"`
		Name      string `json:"name"`
		Code      string `json:"code"`
		PriceUnit string `json:"priceUnit"`
	} `json:"priceSets"`
}

// ListInstanceTypesResult structure parses the list instance types response payload
type ListInstanceTypesResult struct {
	InstanceTypes *[]InstanceType `json:"instanceTypes"`
	Meta          *MetaResult     `json:"meta"`
}

type GetInstanceTypeResult struct {
	InstanceType *InstanceType `json:"instanceType"`
}

type CreateInstanceTypeResult struct {
	Success      bool              `json:"success"`
	Message      string            `json:"msg"`
	Errors       map[string]string `json:"errors"`
	InstanceType *InstanceType     `json:"instanceType"`
}

type UpdateInstanceTypeResult struct {
	CreateInstanceTypeResult
}

type DeleteInstanceTypeResult struct {
	DeleteResult
}

// Client request methods

func (client *Client) ListInstanceTypes(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        InstanceTypesPath,
		QueryParams: req.QueryParams,
		Result:      &ListInstanceTypesResult{},
	})
}

func (client *Client) GetInstanceType(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", InstanceTypesPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetInstanceTypeResult{},
	})
}

func (client *Client) CreateInstanceType(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        InstanceTypesPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateInstanceTypeResult{},
	})
}

func (client *Client) UpdateInstanceType(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", InstanceTypesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateInstanceTypeResult{},
	})
}

func (client *Client) DeleteInstanceType(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", InstanceTypesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteInstanceTypeResult{},
	})
}

func (client *Client) ToggleFeaturedInstanceType(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d/toggle-featured", InstanceTypesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateInstanceTypeResult{},
	})
}

func (client *Client) UpdateInstanceTypeLogo(id int64, filePayload []*FilePayload, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:         "POST",
		Path:           fmt.Sprintf("/api/library/instance-types/%d/update-logo", id),
		IsMultiPart:    true,
		MultiPartFiles: filePayload,
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
		Result: &UpdateInstanceTypeResult{},
	})
}

// helper functions
func (client *Client) FindInstanceTypeByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListInstanceTypes(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListInstanceTypesResult)
	instanceTypesCount := len(*listResult.InstanceTypes)
	if instanceTypesCount != 1 {
		return resp, fmt.Errorf("found %d InstanceTypes for %v", instanceTypesCount, name)
	}
	firstRecord := (*listResult.InstanceTypes)[0]
	instanceTypeId := firstRecord.ID
	return client.GetInstanceType(instanceTypeId, &Request{})
}
