package morpheus

import (
	"fmt"
)

var (
	// SecurityGroupsPath is the API endpoint for security groups
	SecurityGroupsPath = "/api/security-groups"
)

// SecurityGroup structures for use in request and response payloads
type SecurityGroup struct {
	ID          int64       `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	AccountId   int64       `json:"accountId"`
	GroupSource interface{} `json:"groupSource"`
	ExternalId  string      `json:"externalId"`
	Enabled     bool        `json:"enabled"`
	SyncSource  string      `json:"syncSource"`
	Visibility  string      `json:"visibility"`
	Active      bool        `json:"active"`
	Zone        struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"zone"`
	Locations []SecurityGroupLocation `json:"locations"`
	Rules     []SecurityGroupRule     `json:"rules"`
	Tenants   []struct {
		ID        int64  `json:"id"`
		Name      string `json:"name"`
		CanManage bool   `json:"canManage"`
	} `json:"tenants"`
	ResourcePermission struct {
		DefaultStore  bool `json:"defaultStore"`
		AllPlans      bool `json:"allPlans"`
		DefaultTarget bool `json:"defaultTarget"`
		CanManage     bool `json:"canManage"`
		All           bool `json:"all"`
		Account       struct {
			ID int64 `json:"id"`
		} `json:"account"`
		Sites []struct {
			ID      int64  `json:"id"`
			Name    string `json:"name"`
			Default bool   `json:"default"`
		} `json:"sites"`
		Plans []struct {
			ID      int64  `json:"id"`
			Name    string `json:"name"`
			Default bool   `json:"default"`
		} `json:"plans"`
	} `json:"resourcePermission"`
}

type SecurityGroupLocation struct {
	ID          int64       `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Externalid  string      `json:"externalId"`
	IacId       interface{} `json:"iacId"`
	Zone        struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"zone"`
	ZonePool struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"zonePool"`
	Status     string      `json:"status"`
	Priority   interface{} `json:"priority"`
	GroupLayer interface{} `json:"groupLayer"`
}

type SecurityGroupRule struct {
	ID               int64  `json:"id"`
	Name             string `json:"name"`
	RuleType         string `json:"ruleType"`
	CustomRule       bool   `json:"customRule"`
	InstanceTypeId   int64  `json:"instanceTypeId"`
	Direction        string `json:"direction"`
	Policy           string `json:"policy"`
	SourceType       string `json:"sourceType"`
	Source           string `json:"source"`
	SourceGroup      string `json:"sourceGroup"`
	SourceTier       string `json:"sourceTier"`
	PortRange        string `json:"portRange"`
	Protocol         string `json:"protocol"`
	DestinationType  string `json:"destinationType"`
	Destination      string `json:"destination"`
	DestinationGroup string `json:"destinationGroup"`
	DestinationTier  string `json:"destinationTier"`
	ExternalId       string `json:"externalId"`
	Enabled          bool   `json:"enabled"`
}

type ListSecurityGroupsResult struct {
	SecurityGroups *[]SecurityGroup `json:"securityGroups"`
	Meta           *MetaResult      `json:"meta"`
}

type ListSecurityGroupRulesResult struct {
	SecurityGroupRules *[]SecurityGroupRule `json:"rules"`
	Meta               *MetaResult          `json:"meta"`
}

type GetSecurityGroupResult struct {
	SecurityGroup *SecurityGroup `json:"securityGroup"`
}

type GetSecurityGroupRuleResult struct {
	SecurityGroupRule *SecurityGroupRule `json:"rule"`
}

type CreateSecurityGroupResult struct {
	Success       bool              `json:"success"`
	Message       string            `json:"msg"`
	Errors        map[string]string `json:"errors"`
	SecurityGroup *SecurityGroup    `json:"securityGroup"`
}

type CreateSecurityGroupRuleResult struct {
	Success           bool               `json:"success"`
	Message           string             `json:"msg"`
	Errors            map[string]string  `json:"errors"`
	SecurityGroupRule *SecurityGroupRule `json:"rule"`
}

type CreateSecurityGroupLocationResult struct {
	Success               bool                   `json:"success"`
	Message               string                 `json:"msg"`
	Errors                map[string]string      `json:"errors"`
	SecurityGroupLocation *SecurityGroupLocation `json:"location"`
}

type UpdateSecurityGroupResult struct {
	CreateSecurityGroupResult
}

type UpdateSecurityGroupRuleResult struct {
	CreateSecurityGroupRuleResult
}

type DeleteSecurityGroupResult struct {
	DeleteResult
}

type DeleteSecurityGroupRuleResult struct {
	DeleteResult
}

type DeleteSecurityGroupLocationResult struct {
	DeleteResult
}

// Client request methods
func (client *Client) ListSecurityGroups(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        SecurityGroupsPath,
		QueryParams: req.QueryParams,
		Result:      &ListSecurityGroupsResult{},
	})
}

func (client *Client) ListSecurityGroupRules(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d/rules", SecurityGroupsPath, id),
		QueryParams: req.QueryParams,
		Result:      &ListSecurityGroupRulesResult{},
	})
}

func (client *Client) GetSecurityGroup(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", SecurityGroupsPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetSecurityGroupResult{},
	})
}

func (client *Client) GetSecurityGroupRule(id int64, ruleId int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d/rules/%d", SecurityGroupsPath, id, ruleId),
		QueryParams: req.QueryParams,
		Result:      &GetSecurityGroupRuleResult{},
	})
}

func (client *Client) CreateSecurityGroup(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        SecurityGroupsPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateSecurityGroupResult{},
	})
}

func (client *Client) CreateSecurityGroupRule(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        fmt.Sprintf("%s/%d/rules", SecurityGroupsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateSecurityGroupRuleResult{},
	})
}

func (client *Client) UpdateSecurityGroup(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", SecurityGroupsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateSecurityGroupResult{},
	})
}

func (client *Client) UpdateSecurityGroupRule(id int64, ruleId int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d/rules/%d", SecurityGroupsPath, id, ruleId),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateSecurityGroupRuleResult{},
	})
}

func (client *Client) DeleteSecurityGroup(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", SecurityGroupsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteSecurityGroupResult{},
	})
}

func (client *Client) DeleteSecurityGroupRule(id int64, ruleId int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d/rules/%d", SecurityGroupsPath, id, ruleId),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteSecurityGroupRuleResult{},
	})
}

// FindSecurityGroupByName gets an existing security group by name
func (client *Client) FindSecurityGroupByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListSecurityGroups(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListSecurityGroupsResult)
	securityGroupCount := len(*listResult.SecurityGroups)
	if securityGroupCount != 1 {
		return resp, fmt.Errorf("found %d security groups for %v", securityGroupCount, name)
	}
	firstRecord := (*listResult.SecurityGroups)[0]
	securityGroupID := firstRecord.ID
	return client.GetSecurityGroup(securityGroupID, &Request{})
}

// FindSecurityGroupByName gets an existing security group by name
func (client *Client) FindSecurityGroupRuleByName(id int64, name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListSecurityGroupRules(id, &Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListSecurityGroupRulesResult)
	securityGroupRuleCount := len(*listResult.SecurityGroupRules)
	if securityGroupRuleCount != 1 {
		return resp, fmt.Errorf("found %d security group rules for %v", securityGroupRuleCount, name)
	}
	firstRecord := (*listResult.SecurityGroupRules)[0]
	securityGroupRuleID := firstRecord.ID
	return client.GetSecurityGroupRule(id, securityGroupRuleID, &Request{})
}

func (client *Client) CreateSecurityGroupLocation(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        fmt.Sprintf("%s/%d/locations", SecurityGroupsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateSecurityGroupLocationResult{},
	})
}

func (client *Client) DeleteSecurityGroupLocation(id int64, locationId int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d/locations/%d", SecurityGroupsPath, id, locationId),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteSecurityGroupLocationResult{},
	})
}
