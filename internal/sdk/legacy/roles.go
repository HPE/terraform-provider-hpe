package morpheus

import (
	"fmt"
	"time"
)

var (
	// RolesPath is the API endpoint for roles
	RolesPath = "/api/roles"
	// RolesPath is the API endpoint for tenant roles
	TenantRolesPath = "/api/accounts/available-roles"
)

// Role structures for use in request and response payloads
type GetRoleResult struct {
	Role               Role `json:"role"`
	FeaturePermissions []struct {
		ID          int64  `json:"id"`
		Code        string `json:"code"`
		Name        string `json:"name"`
		Access      string `json:"access"`
		SubCategory string `json:"subCategory"`
	} `json:"featurePermissions"`
	GlobalSiteAccess string `json:"globalSiteAccess"`
	Sites            []struct {
		ID     int64  `json:"id"`
		Name   string `json:"name"`
		Access string `json:"access"`
	} `json:"sites"`
	GlobalZoneAccess string `json:"globalZoneAccess"`
	Zones            []struct {
		ID     int64  `json:"id"`
		Code   string `json:"code"`
		Name   string `json:"name"`
		Access string `json:"access"`
	} `json:"zones"`
	GlobalInstanceTypeAccess string `json:"globalInstanceTypeAccess"`
	InstanceTypePermissions  []struct {
		ID     int64  `json:"id"`
		Code   string `json:"code"`
		Name   string `json:"name"`
		Access string `json:"access"`
	} `json:"instanceTypePermissions"`
	GlobalAppTemplateAccess string `json:"globalAppTemplateAccess"`
	AppTemplatePermissions  []struct {
		ID     int64  `json:"id"`
		Code   string `json:"code"`
		Name   string `json:"name"`
		Access string `json:"access"`
	} `json:"appTemplatePermissions"`
	GlobalCatalogItemTypeAccess string `json:"globalCatalogItemTypeAccess"`
	CatalogItemTypePermissions  []struct {
		ID     int64  `json:"id"`
		Name   string `json:"name"`
		Access string `json:"access"`
	} `json:"catalogItemTypePermissions"`
	GlobalPersonaAccess string `json:"globalPersonaAccess"`
	PersonaPermissions []struct {
		ID     int64  `json:"id"`
		Code   string `json:"code"`
		Name   string `json:"name"`
		Access string `json:"access"`
	} `json:"personaPermissions"`
	GlobalVDIPoolAccess string `json:"globalVdiPoolAccess"`
	VDIPoolPermissions  []struct {
		ID     int64  `json:"id"`
		Code   string `json:"code"`
		Name   string `json:"name"`
		Access string `json:"access"`
	} `json:"vdiPoolPermissions"`
	GlobalReportTypeAccess string `json:"globalReportTypeAccess"`
	ReportTypePermissions  []struct {
		ID     int64  `json:"id"`
		Code   string `json:"code"`
		Name   string `json:"name"`
		Access string `json:"access"`
	} `json:"reportTypePermissions"`
	GlobalTaskAccess string `json:"globalTaskAccess"`
	TaskPermissions  []struct {
		ID     int64  `json:"id"`
		Code   string `json:"code"`
		Name   string `json:"name"`
		Access string `json:"access"`
	} `json:"taskPermissions"`
	GlobalTaskSetAccess string `json:"globalTaskSetAccess"`
	TaskSetPermissions  []struct {
		ID     int64  `json:"id"`
		Code   string `json:"code"`
		Name   string `json:"name"`
		Access string `json:"access"`
	} `json:"taskSetPermissions"`
	GlobalClusterTypeAccess string `json:"globalClusterTypeAccess"`
	ClusterTypePermissions  []struct {
		ID     int64  `json:"id"`
		Code   string `json:"code"`
		Name   string `json:"name"`
		Access string `json:"access"`
	} `json:"clusterTypePermissions"`
}

type Role struct {
	ID                int64  `json:"id"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	LandingUrl        string `json:"landingUrl"`
	Scope             string `json:"scope"`
	RoleType          string `json:"roleType"`
	MultiTenant       bool   `json:"multitenant"`
	MultiTenantLocked bool   `json:"multitenantLocked"`
	Diverged          bool   `json:"diverged"`
	OwnerId           int64  `json:"ownerId"`
	Authority         string `json:"authority"`
	Owner             struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"owner"`
	DefaultPersona struct {
		ID   int    `json:"id"`
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"defaultPersona"`
	DateCreated time.Time `json:"dateCreated"`
	LastUpdated time.Time `json:"lastUpdated"`
}

type ListRolesResult struct {
	Roles *[]Role     `json:"roles"`
	Meta  *MetaResult `json:"meta"`
}

type CreateRoleResult struct {
	Success bool              `json:"success"`
	Message string            `json:"msg"`
	Errors  map[string]string `json:"errors"`
	Role    *Role             `json:"role"`
}

type UpdateRoleResult struct {
	CreateRoleResult
}

type UpdateRolePermissionResult struct {
	Success bool              `json:"success"`
	Message string            `json:"msg"`
	Errors  map[string]string `json:"errors"`
	Access  string            `json:"access"`
}

type DeleteRoleResult struct {
	DeleteResult
}

// Client request methods
func (client *Client) ListRoles(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        RolesPath,
		QueryParams: req.QueryParams,
		Result:      &ListRolesResult{},
	})
}

func (client *Client) ListTenantRoles(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        TenantRolesPath,
		QueryParams: req.QueryParams,
		Result:      &ListRolesResult{},
	})
}

func (client *Client) GetRole(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", RolesPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetRoleResult{},
	})
}

func (client *Client) CreateRole(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        RolesPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateRoleResult{},
	})
}

func (client *Client) UpdateRole(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", RolesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateRoleResult{},
	})
}

func (client *Client) UpdateRoleBlueprintAccess(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d/update-blueprint", RolesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateRolePermissionResult{},
	})
}

func (client *Client) UpdateRoleCatalogItemTypeAccess(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d/update-catalog-item-type", RolesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateRolePermissionResult{},
	})
}

func (client *Client) UpdateRoleCloudAccess(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d/update-cloud", RolesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateRolePermissionResult{},
	})
}

func (client *Client) UpdateRoleGroupAccess(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d/update-group", RolesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateRolePermissionResult{},
	})
}

func (client *Client) UpdateRoleInstanceTypeAccess(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d/update-instance-type", RolesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateRolePermissionResult{},
	})
}

func (client *Client) UpdateRoleFeaturePermission(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d/update-permission", RolesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateRolePermissionResult{},
	})
}

func (client *Client) UpdateRolePersonaAccess(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d/update-persona", RolesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateRolePermissionResult{},
	})
}

func (client *Client) UpdateRoleReportTypeAccess(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d/update-report-type", RolesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateRolePermissionResult{},
	})
}

func (client *Client) UpdateRoleTaskAccess(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d/update-task", RolesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateRolePermissionResult{},
	})
}

func (client *Client) UpdateRoleWorkflowAccess(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d/update-task-set", RolesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateRolePermissionResult{},
	})
}

func (client *Client) UpdateRoleVDIPoolAccess(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d/update-vdi-pool", RolesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateRolePermissionResult{},
	})
}

func (client *Client) DeleteRole(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", RolesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteRoleResult{},
	})
}

// FindRoleByName gets an existing role by name
func (client *Client) FindRoleByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListRoles(&Request{
		QueryParams: map[string]string{
			"authority": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListRolesResult)
	rolesCount := len(*listResult.Roles)
	if rolesCount != 1 {
		return resp, fmt.Errorf("found %d Roles for %v", rolesCount, name)
	}
	firstRecord := (*listResult.Roles)[0]
	roleID := firstRecord.ID
	return client.GetRole(roleID, &Request{})
}

// FindTenantRoleByName gets an existing tenant role by name
func (client *Client) FindTenantRoleByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListTenantRoles(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListRolesResult)
	for _, role := range *listResult.Roles {
		if role.Authority == name {
			return client.GetRole(role.ID, &Request{})
		}
	}
	return resp, fmt.Errorf("not matching role found for %v", name)
}
