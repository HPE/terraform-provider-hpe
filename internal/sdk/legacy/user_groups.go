package morpheus

import (
	"fmt"
	"time"
)

var (
	// UserGroupsPath is the API endpoint for user groups
	UserGroupsPath = "/api/user-groups"
)

// UserGroup structures for use in request and response payloads
type UserGroup struct {
	ID          int64  `json:"id"`
	AccountID   int64  `json:"accountId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
	SudoUser    bool   `json:"sudoUser"`
	ServerGroup string `json:"serverGroup"`
	Users       []struct {
		ID          int64  `json:"id"`
		Username    string `json:"username"`
		DisplayName string `json:"displayName"`
	} `json:"users"`
	Account struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"account"`
	DateCreated time.Time `json:"dateCreated"`
	LastUpdated time.Time `json:"lastUpdated"`
}

type ListUserGroupsResult struct {
	UserGroups *[]UserGroup `json:"userGroups"`
	Meta       *MetaResult  `json:"meta"`
}

type GetUserGroupResult struct {
	UserGroup *UserGroup `json:"userGroup"`
}

type CreateUserGroupResult struct {
	Success   bool              `json:"success"`
	Message   string            `json:"msg"`
	Errors    map[string]string `json:"errors"`
	UserGroup *UserGroup        `json:"userGroup"`
}

type UpdateUserGroupResult struct {
	CreateUserGroupResult
}

type DeleteUserGroupResult struct {
	DeleteResult
}

// Client request methods
func (client *Client) ListUserGroups(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        UserGroupsPath,
		QueryParams: req.QueryParams,
		Result:      &ListUserGroupsResult{},
	})
}

func (client *Client) GetUserGroup(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", UserGroupsPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetUserGroupResult{},
	})
}

func (client *Client) CreateUserGroup(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        UserGroupsPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateUserGroupResult{},
	})
}

func (client *Client) UpdateUserGroup(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", UserGroupsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateUserGroupResult{},
	})
}

func (client *Client) DeleteUserGroup(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", UserGroupsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteUserGroupResult{},
	})
}

// FindUserGroupByName gets an existing user group by name
func (client *Client) FindUserGroupByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListUserGroups(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListUserGroupsResult)
	userGroupCount := len(*listResult.UserGroups)
	if userGroupCount != 1 {
		return resp, fmt.Errorf("found %d User Groups for %v", userGroupCount, name)
	}
	firstRecord := (*listResult.UserGroups)[0]
	userGroupID := firstRecord.ID
	return client.GetUserGroup(userGroupID, &Request{})
}
