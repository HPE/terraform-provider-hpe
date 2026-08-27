package morpheus

import (
	"fmt"
	"time"
)

var (
	// UsersPath is the API endpoint for user
	UsersPath = "/api/users"
)

// User structures for use in request and response payloads
type User struct {
	ID                   int64     `json:"id"`
	AccountID            int64     `json:"accountId"`
	Username             string    `json:"username"`
	DisplayName          string    `json:"displayName"`
	Email                string    `json:"email"`
	FirstName            string    `json:"firstName"`
	LastName             string    `json:"lastName"`
	Enabled              bool      `json:"enabled"`
	ReceiveNotifications bool      `json:"receiveNotifications"`
	Isusing2FA           bool      `json:"isUsing2FA"`
	AccountExpired       bool      `json:"accountExpired"`
	AccountLocked        bool      `json:"accountLocked"`
	PasswordExpired      bool      `json:"passwordExpired"`
	LoginCount           int64     `json:"loginCount"`
	LoginAttempts        int64     `json:"loginAttempts"`
	LastLoginDate        time.Time `json:"lastLoginDate"`
	Roles                []struct {
		ID          int64  `json:"id"`
		Name        string `json:"name"`
		Authority   string `json:"authority"`
		Description string `json:"description"`
	} `json:"roles"`
	Account struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"account"`
	LinuxUsername   string      `json:"linuxUsername"`
	LinuxPassword   string      `json:"linuxPassword"`
	LinuxKeyPairID  int64       `json:"linuxKeyPairId"`
	WindowsUsername interface{} `json:"windowsUsername"`
	WindowsPassword interface{} `json:"windowsPassword"`
	DefaultPersona  interface{} `json:"defaultPersona"`
	DateCreated     time.Time   `json:"dateCreated"`
	LastUpdated     time.Time   `json:"lastUpdated"`
}

type ListUsersResult struct {
	Users *[]User     `json:"users"`
	Meta  *MetaResult `json:"meta"`
}

type GetUserResult struct {
	User *User `json:"user"`
}

type CreateUserResult struct {
	Success bool              `json:"success"`
	Message string            `json:"msg"`
	Errors  map[string]string `json:"errors"`
	User    *User             `json:"user"`
}

type UpdateUserResult struct {
	CreateUserResult
}

type DeleteUserResult struct {
	DeleteResult
}

// Client request methods
func (client *Client) ListUsers(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        UsersPath,
		QueryParams: req.QueryParams,
		Result:      &ListUsersResult{},
	})
}

func (client *Client) GetUser(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", UsersPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetUserResult{},
	})
}

func (client *Client) CreateUser(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        UsersPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateUserResult{},
	})
}

func (client *Client) UpdateUser(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", UsersPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateUserResult{},
	})
}

func (client *Client) DeleteUserResult(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", UsersPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteUserResult{},
	})
}

// FindUserByName gets an existing user by name
func (client *Client) FindUserByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListUsers(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListUsersResult)
	userCount := len(*listResult.Users)
	if userCount != 1 {
		return resp, fmt.Errorf("found %d Users for %v", userCount, name)
	}
	firstRecord := (*listResult.Users)[0]
	userID := firstRecord.ID
	return client.GetUser(userID, &Request{})
}

func (client *Client) FindUserByExactName(name string) (*Response, error) {
	// Find by name, then get by ID

	resp, err := client.ListUsers(&Request{
		QueryParams: map[string]string{
			"name": name,
			"max":  "1000", //this is needed for tenants with a lot of users that could be returned
		},
	})
	if err != nil {
		return resp, err
	}

	listResult := resp.Result.(*ListUsersResult)
	var userID int64 = 0
	for _, user := range *listResult.Users {
		if user.Username == name {
			userID = user.ID
		}
	}
	return client.GetUser(userID, &Request{})
}
