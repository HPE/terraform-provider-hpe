package morpheus

import (
	"fmt"
)

var (
	// GroupsPath is the API endpoint for groups
	GroupsPath = "/api/groups"
)

// Group structures for use in request and response payloads
type Group struct {
	ID        int64    `json:"id"`
	UUID      string   `json:"uuid"`
	Name      string   `json:"name"`
	Code      string   `json:"code"`
	Labels    []string `json:"labels"`
	Location  string   `json:"location"`
	AccountID int64    `json:"accountId"`
	Config    struct {
		DNSIntegrationID    string `json:"dnsIntegrationId"`
		ConfigCMDBID        string `json:"configCmdbId"`
		ConfigCMID          string `json:"configCmId"`
		ServiceRegistryID   string `json:"serviceRegistryId"`
		ConfigManagementID  string `json:"configManagementId"`
		ConfigCMDBDiscovery bool   `json:"configCmdbDiscovery"`
	} `json:"config"`
	Clouds      []Zone `json:"zones"`
	DateCreated string `json:"dateCreated"`
	LastUpdated string `json:"lastUpdated"`
	Stats       struct {
		InstanceCounts struct {
			All int64 `json:"all"`
		} `json:"instanceCounts"`
		ServerCounts struct {
			All           int64 `json:"all"`
			Host          int64 `json:"host"`
			Hypervisor    int64 `json:"hypervisor"`
			ContainerHost int64 `json:"containerHost"`
			VM            int64 `json:"vm"`
			Baremetal     int64 `json:"baremetal"`
			Unmanaged     int64 `json:"unmanaged"`
		} `json:"serverCounts"`
	} `json:"stats"`
	ServerCount int64 `json:"serverCount"`
}

type Zone struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// ListGroupsResult structure parses the list groups response payload
type ListGroupsResult struct {
	Groups *[]Group    `json:"groups"`
	Meta   *MetaResult `json:"meta"`
}

type GetGroupResult struct {
	Group *Group `json:"group"`
}

type CreateGroupResult struct {
	Success bool              `json:"success"`
	Message string            `json:"msg"`
	Errors  map[string]string `json:"errors"`
	Group   *Group            `json:"group"`
}

type UpdateGroupResult struct {
	CreateGroupResult
}

type DeleteGroupResult struct {
	DeleteResult
}

// API endpoints

func (client *Client) ListGroups(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        GroupsPath,
		QueryParams: req.QueryParams,
		Result:      &ListGroupsResult{},
	})
}

func (client *Client) GetGroup(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", GroupsPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetGroupResult{},
	})
}

func (client *Client) CreateGroup(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        GroupsPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateGroupResult{},
	})
}

func (client *Client) UpdateGroup(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", GroupsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateGroupResult{},
	})
}

func (client *Client) UpdateGroupClouds(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d/update-zones", GroupsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &StandardResult{},
	})
}

func (client *Client) UpdateGroupZones(id int64, req *Request) (*Response, error) {
	return client.UpdateGroupClouds(id, req)
}

// DeleteGroup deletes an existing group
func (client *Client) DeleteGroup(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", GroupsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteGroupResult{},
	})
}

// FindGroupByName gets an existing group by name
func (client *Client) FindGroupByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListGroups(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListGroupsResult)
	groupsCount := len(*listResult.Groups)
	if groupsCount != 1 {
		return resp, fmt.Errorf("found %d Groups for %v", groupsCount, name)
	}
	firstRecord := (*listResult.Groups)[0]
	groupId := firstRecord.ID
	return client.GetGroup(groupId, &Request{})
}
