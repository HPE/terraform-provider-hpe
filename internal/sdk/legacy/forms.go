package morpheus

import (
	"fmt"
)

var (
	// FormsPath is the API endpoint for forms
	FormsPath = "/api/library/option-type-forms"
)

// Form structures for use in request and response payloads
type Form struct {
	ID          int64        `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Code        string       `json:"code"`
	Context     string       `json:"context"`
	Locked      bool         `json:"locked"`
	Labels      []string     `json:"labels"`
	Options     []Option     `json:"options"`
	FieldGroups []FieldGroup `json:"fieldGroups"`
}

type FieldGroup struct {
	Name             string   `json:"name"`
	Code             string   `json:"code"`
	Description      string   `json:"description"`
	LocalizedName    string   `json:"localizedName"`
	Collapsible      bool     `json:"collapsible"`
	DefaultCollapsed bool     `json:"defaultCollapsed"`
	VisibleOnCode    string   `json:"visibleOnCode"`
	Options          []Option `json:"options"`
}

type Option struct {
	ID                 int64       `json:"id"`
	Name               string      `json:"name"`
	Description        string      `json:"description"`
	Labels             []string    `json:"labels"`
	Code               string      `json:"code"`
	FieldName          string      `json:"fieldName"`
	FieldLabel         string      `json:"fieldLabel"`
	FieldCode          string      `json:"fieldCode"`
	FieldContext       string      `json:"fieldContext"`
	FieldGroup         interface{} `json:"fieldGroup"`
	FieldClass         interface{} `json:"fieldClass"`
	FieldAddOn         interface{} `json:"fieldAddOn"`
	FieldComponent     interface{} `json:"fieldComponent"`
	FieldInput         interface{} `json:"fieldInput"`
	PlaceHolder        string      `json:"placeHolder"`
	VerifyPattern      string      `json:"verifyPattern"`
	HelpBlock          string      `json:"helpBlock"`
	HelpBlockFieldCode string      `json:"helpBlockFieldCode"`
	DefaultValue       string      `json:"defaultValue"`
	OptionSource       string      `json:"optionSource"`
	OptionSourceType   string      `json:"optionSourceType"`
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
		AddOn                      string `json:"addon"`
		AddOnPosition              string `json:"addonPosition"`
		AllowDuplicates            bool   `json:"allowDuplicates"`
		AllowReadonly              bool   `json:"allowReadonly"`
		AsObject                   bool   `json:"asObject"`
		CanPeek                    bool   `json:"canPeek"`
		CloudCodeField             string `json:"cloudCodeField"`
		CloudField                 string `json:"cloudField"`
		CloudFieldType             string `json:"cloudFieldType"`
		CloudType                  string `json:"cloudType"`
		CloudId                    string `json:"cloudId"`
		CustomData                 string `json:"customData"`
		DefaultValue               string `json:"defaultValue"`
		DiskField                  string `json:"diskField"`
		Display                    string `json:"display"`
		EnableDatastoreSelection   bool   `json:"enableDatastoreSelection"`
		EnableDiskTypeSelection    bool   `json:"enableDiskTypeSelection"`
		EnableIPModeSelection      bool   `json:"enableIPModeSelection"`
		EnableStorageTypeSelection bool   `json:"enableStorageTypeSelection"`
		FilterResource             bool   `json:"filterResource"`
		Group                      string `json:"group"`
		GroupId                    string `json:"groupId"`
		GroupField                 string `json:"groupField"`
		GroupFieldType             string `json:"groupFieldType"`
		ImageField                 string `json:"imageField"`
		ImageId                    string `json:"imageId"`
		InstanceTypeCode           string `json:"instanceTypeCode"`
		InstanceTypeFieldCode      string `json:"instanceTypeFieldCode"`
		InstanceTypeFieldType      string `json:"instanceTypeFieldType"`
		KeyPlaceholder             string `json:"keyPlaceholder"`
		Lang                       string `json:"lang"`
		LayoutId                   string `json:"layoutId"`
		LayoutField                string `json:"layoutField"`
		LayoutFieldType            string `json:"layoutFieldType"`
		LockDisplay                bool   `json:"lockDisplay"`
		Sortable                   bool   `json:"sortable"`
		MultiSelect                bool   `json:"multiSelect"`
		PlanFieldType              string `json:"planFieldType"`
		PlanField                  string `json:"planField"`
		PlanId                     string `json:"planId"`
		PoolId                     string `json:"poolId"`
		PoolField                  string `json:"poolField"`
		PoolFieldType              string `json:"poolFieldType"`
		ResourcePoolField          string `json:"resourcePoolField"`
		Separator                  string `json:"separator"`
		ShowLineNumbers            bool   `json:"showLineNumbers"`
		ShowNetworkTypeSelection   bool   `json:"showNetworkTypeSelection"`
		ShowPricing                bool   `json:"showPricing"`
		Step                       int64  `json:"step"`
		Rows                       int64  `json:"rows"`
		ValuePlaceholder           string `json:"valuePlaceholder"`
		VirtualImageFieldType      string `json:"virtualImageFieldType"`
	} `json:"config"`
	DisplayOrder          int64       `json:"displayOrder"`
	WrapperClass          interface{} `json:"wrapperClass"`
	Enabled               bool        `json:"enabled"`
	NoBlank               bool        `json:"noBlank"`
	DependsOnCode         string      `json:"dependsOnCode"`
	VisibleOnCode         string      `json:"visibleOnCode"`
	RequireOnCode         string      `json:"requireOnCode"`
	ContextualDefault     bool        `json:"contextualDefault"`
	DisplayValueOnDetails bool        `json:"displayValueOnDetails"`
	ShowOnCreate          bool        `json:"showOnCreate"`
	ShowOnEdit            bool        `json:"showOnEdit"`
	LocalCredential       bool        `json:"localCredential"`
	FormField             bool        `json:"formField"`
	ExcludeFromSearch     bool        `json:"excludeFromSearch"`
	IsHidden              bool        `json:"isHidden"`
	IsLocked              bool        `json:"isLocked"`
	MinVal                int64       `json:"minVal"`
	MaxVal                int64       `json:"maxVal"`
}

// ListFormsResult structure parses the list forms response payload
type ListFormsResult struct {
	Forms *[]Form     `json:"optionTypeForms"`
	Meta  *MetaResult `json:"meta"`
}

// GetFormResult structure parses the get forms response payload
type GetFormResult struct {
	Form *Form `json:"optionTypeForm"`
}

type CreateFormResult struct {
	Success bool              `json:"success"`
	Message string            `json:"msg"`
	Errors  map[string]string `json:"errors"`
	Form    *Form             `json:"optionTypeForm"`
}

type UpdateFormResult struct {
	CreateFormResult
}

type DeleteFormResult struct {
	DeleteResult
}

// API endpoints
func (client *Client) ListForms(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        FormsPath,
		QueryParams: req.QueryParams,
		Result:      &ListFormsResult{},
	})
}

func (client *Client) GetForm(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", FormsPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetFormResult{},
	})
}

func (client *Client) CreateForm(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        FormsPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateFormResult{},
	})
}

func (client *Client) UpdateForm(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", FormsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateFormResult{},
	})
}

// DeleteForm deletes an existing form
func (client *Client) DeleteForm(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", FormsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteFormResult{},
	})
}

// FindFormByName gets an existing form by name
func (client *Client) FindFormByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListForms(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListFormsResult)
	formsCount := len(*listResult.Forms)
	if formsCount != 1 {
		return resp, fmt.Errorf("found %d forms for %v", formsCount, name)
	}
	firstRecord := (*listResult.Forms)[0]
	formId := firstRecord.ID
	return client.GetForm(formId, &Request{})
}
