// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package client

// DeviceManagementPolicy represents a Resource Management (Device Management) policy.
// A policy filters resources using rules and applies actions (monitoring templates,
// availability rules, credential sets, custom attributes, jobs, knowledge articles) to
// the matching resources once the policy is run.
type DeviceManagementPolicy struct {
	Id                  int                            `json:"id,omitempty"`
	UID                 string                         `json:"uid,omitempty"`
	Name                string                         `json:"name"`
	ResourceType        string                         `json:"resourceType,omitempty"`
	MatchType           string                         `json:"matchType,omitempty"`
	Rules               []DeviceManagementPolicyRule   `json:"rules,omitempty"`
	Actions             []DeviceManagementPolicyAction `json:"actions,omitempty"`
	ExecutionOrder      int                            `json:"executionOrder,omitempty"`
	ClientMatchingType  string                         `json:"clientMatchingType,omitempty"`
	Clients             []string                       `json:"clients,omitempty"`
	SearchQuery         string                         `json:"searchQuery,omitempty"`
	ParentPolicy        int                            `json:"parentPolicy,omitempty"`
	SystemQuery         string                         `json:"systemQuery,omitempty"`
	CompleteSearchQuery string                         `json:"completeSearchQuery,omitempty"`
	ParentPolicyUID     string                         `json:"parentPolicyUID,omitempty"`
	Customized          bool                           `json:"customized,omitempty"`
	GlobalPolicy        bool                           `json:"globalPolicy,omitempty"`
	AvailableAppId      string                         `json:"availableAppId,omitempty"`
	NativeType          string                         `json:"nativeType,omitempty"`
}

// DeviceManagementPolicyRule filters resources by a native/custom attribute, or by a
// system filterType (e.g. "ANY_CLOUD_RESOURCE").
type DeviceManagementPolicyRule struct {
	Id                int    `json:"id,omitempty"`
	Key               string `json:"key,omitempty"`
	Operator          string `json:"operator,omitempty"`
	Value             string `json:"value,omitempty"`
	IsNativeAttribute *bool  `json:"isNativeAttribute,omitempty"`
	FilterType        string `json:"filterType,omitempty"`
}

// DeviceManagementPolicyAction assigns or unassigns an item (monitoring template,
// availability rule, credential set, custom attribute, job, knowledge article) to the
// resources matched by the policy rules.
type DeviceManagementPolicyAction struct {
	Action                string                             `json:"action"`
	Items                 []DeviceManagementPolicyActionItem `json:"items,omitempty"`
	ForceAssignOrUnassign *bool                              `json:"forceAssignOrUnassign,omitempty"`
}

// DeviceManagementPolicyActionItem is a single item assigned by an action.
type DeviceManagementPolicyActionItem struct {
	Item       string            `json:"item"`
	ItemType   string            `json:"itemType,omitempty"`
	ItemUID    string            `json:"itemUid,omitempty"`
	Parameters map[string]string `json:"parameters,omitempty"`
}

type DeviceManagementPolicySearchResult struct {
	Results         []DeviceManagementPolicy `json:"results"`
	TotalResults    int                      `json:"totalResults"`
	OrderBy         string                   `json:"orderBy,omitempty"`
	PageNo          int                      `json:"pageNo,omitempty"`
	PageSize        int                      `json:"pageSize,omitempty"`
	TotalPages      int                      `json:"totalPages,omitempty"`
	NextPage        bool                     `json:"nextPage,omitempty"`
	PreviousPageNo  int                      `json:"previousPageNo,omitempty"`
	DescendingOrder bool                     `json:"descendingOrder,omitempty"`
}
