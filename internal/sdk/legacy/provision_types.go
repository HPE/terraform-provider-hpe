package morpheus

import (
	"fmt"
)

var (
	ProvisionTypesPath = "/api/provision-types"
)

// Provision Type structures for use in request and response payloads
type ProvisionType struct {
	ID                   int64       `json:"id"`
	Name                 string      `json:"name"`
	Description          interface{} `json:"description"`
	Code                 string      `json:"code"`
	Aclenabled           bool        `json:"aclEnabled"`
	Multitenant          bool        `json:"multiTenant"`
	Managed              bool        `json:"managed"`
	Hostnetwork          bool        `json:"hostNetwork"`
	Customsupported      bool        `json:"customSupported"`
	Mapports             bool        `json:"mapPorts"`
	Exportserver         interface{} `json:"exportServer"`
	Viewset              string      `json:"viewSet"`
	Servertype           string      `json:"serverType"`
	Hosttype             string      `json:"hostType"`
	Addvolumes           bool        `json:"addVolumes"`
	Hasdatastore         bool        `json:"hasDatastore"`
	Hasnetworks          interface{} `json:"hasNetworks"`
	Maxnetworks          interface{} `json:"maxNetworks"`
	Customizevolume      bool        `json:"customizeVolume"`
	Rootdiskcustomizable bool        `json:"rootDiskCustomizable"`
	Lvmsupported         bool        `json:"lvmSupported"`
	Hostdiskmode         string      `json:"hostDiskMode"`
	Mindisk              int64       `json:"minDisk"`
	Maxdisk              interface{} `json:"maxDisk"`
	Resizecopiesvolumes  bool        `json:"resizeCopiesVolumes"`
	Optiontypes          []struct {
		Name         string      `json:"name"`
		Description  interface{} `json:"description"`
		Fieldname    string      `json:"fieldName"`
		Fieldlabel   string      `json:"fieldLabel"`
		Fieldcontext string      `json:"fieldContext"`
		Fieldaddon   interface{} `json:"fieldAddOn"`
		Placeholder  interface{} `json:"placeHolder"`
		Helpblock    string      `json:"helpBlock"`
		Defaultvalue interface{} `json:"defaultValue"`
		Optionsource string      `json:"optionSource"`
		Type         string      `json:"type"`
		Advanced     bool        `json:"advanced"`
		Required     bool        `json:"required"`
		Editable     bool        `json:"editable"`
		Displayorder int64       `json:"displayOrder"`
	} `json:"optionTypes"`
	Customoptiontypes []interface{} `json:"customOptionTypes"`
	Networktypes      []interface{} `json:"networkTypes"`
	Storagetypes      []struct {
		ID                int64       `json:"id"`
		Code              string      `json:"code"`
		Name              string      `json:"name"`
		Displayorder      int64       `json:"displayOrder"`
		Defaulttype       bool        `json:"defaultType"`
		Customlabel       bool        `json:"customLabel"`
		Customsize        bool        `json:"customSize"`
		Customsizeoptions interface{} `json:"customSizeOptions"`
	} `json:"storageTypes"`
	Rootstoragetypes []struct {
		ID                int64       `json:"id"`
		Code              string      `json:"code"`
		Name              string      `json:"name"`
		Displayorder      int64       `json:"displayOrder"`
		Defaulttype       bool        `json:"defaultType"`
		Customlabel       bool        `json:"customLabel"`
		Customsize        bool        `json:"customSize"`
		Customsizeoptions interface{} `json:"customSizeOptions"`
	} `json:"rootStorageTypes"`
	Controllertypes []interface{} `json:"controllerTypes"`
}

// ListProvisionTypeResult structure parses the list provision types response payload
type ListProvisionTypesResult struct {
	ProvisionTypes *[]ProvisionType `json:"provisionTypes"`
	Meta           *MetaResult      `json:"meta"`
}

// GetProvisionTypeResult structure parses the get provision type response payload
type GetProvisionTypeResult struct {
	ProvisionType *ProvisionType `json:"provisionType"`
}

// API endpoints
// ListProvisionTypes lists all provision types
func (client *Client) ListProvisionTypes(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        ProvisionTypesPath,
		QueryParams: req.QueryParams,
		Result:      &ListProvisionTypesResult{},
	})
}

// GetProvisionType gets a provision type by ID
func (client *Client) GetProvisionType(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", ProvisionTypesPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetProvisionTypeResult{},
	})
}

// FindProvisionTypeByName gets an existing provision type by name
func (client *Client) FindProvisionTypeByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListProvisionTypes(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListProvisionTypesResult)
	provisionTypeCount := len(*listResult.ProvisionTypes)
	if provisionTypeCount != 1 {
		return resp, fmt.Errorf("found %d Provision Types for %v", provisionTypeCount, name)
	}
	firstRecord := (*listResult.ProvisionTypes)[0]
	provisionTypeID := firstRecord.ID
	return client.GetProvisionType(provisionTypeID, &Request{})
}
