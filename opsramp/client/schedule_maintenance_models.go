// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package client

// Scheduled Maintenance Window models

// ScheduleMaintenanceRequest represents the request body to create/update a maintenance window
type ScheduleMaintenanceRequest struct {
	Name              string                   `json:"name"`
	Description       string                   `json:"description"`
	RunRBA            bool                     `json:"runRBA"`
	InstallPatch      bool                     `json:"installPatch"`
	RunEscalateAction bool                     `json:"runEscalateAction,omitempty"`
	CorrelateAlerts   bool                     `json:"correlateAlerts"`
	Schedule          ScheduleMaintenanceTime  `json:"schedule"`
	Devices           []ScheduleDevice         `json:"devices,omitempty"`
	DeviceGroups      []ScheduleDeviceGroup    `json:"deviceGroups,omitempty"`
	Locations         []ScheduleLocation       `json:"locations,omitempty"`
	AlertConditions   *ScheduleAlertConditions `json:"alertConditions,omitempty"`
}

// ScheduleMaintenanceTime represents the schedule timing configuration
type ScheduleMaintenanceTime struct {
	Type      string             `json:"type"`              // "one-time" or "recurring"
	StartTime string             `json:"startTime"`         // ISO 8601 format
	EndTime   string             `json:"endTime"`           // ISO 8601 format
	EndBy     string             `json:"endBy,omitempty"`   // "Never" or specific end date
	Timezone  string             `json:"timezone"`          // e.g. "America/New_York", "GMT"
	Pattern   *SMSchedulePattern `json:"pattern,omitempty"` // Required for Recurring type
}

// SMSchedulePattern defines the recurrence pattern for scheduled maintenance
type SMSchedulePattern struct {
	Type            string `json:"type"`                      // "daily", "weekly", "monthly"
	WeekDays        string `json:"weekDays,omitempty"`        // Comma-separated days for weekly: "Monday,Wednesday"
	DayOfWeek       string `json:"dayOfWeek,omitempty"`       // Day name for monthly pattern
	WeekIndex       string `json:"weekIndex,omitempty"`       // "First", "Second", etc. for monthly pattern
	RepeatFrequency int    `json:"repeatFrequency,omitempty"` // Number of units between occurrences
	Months          string `json:"months,omitempty"`          // Comma-separated month names for yearly pattern
	DayFrequency    string `json:"dayFrequency,omitempty"`    // Number of days between occurrences for daily pattern
	DayOfMonth      string `json:"dayOfMonth,omitempty"`      // Day of month for monthly pattern
}

// ScheduleDevice identifies a device for the maintenance window
type ScheduleDevice struct {
	HostName string `json:"hostName,omitempty"`
	UniqueId string `json:"uniqueId,omitempty"`
}

// ScheduleDeviceGroup identifies a device group for the maintenance window
type ScheduleDeviceGroup struct {
	Name string `json:"name,omitempty"`
	Id   string `json:"id,omitempty"`
}

// ScheduleLocation identifies a location for the maintenance window
type ScheduleLocation struct {
	Name string `json:"name,omitempty"`
	Id   string `json:"id,omitempty"`
}

// ScheduleAlertConditions defines alert matching conditions during maintenance
type ScheduleAlertConditions struct {
	MatchingType string              `json:"matchingType"` // "ANY" or "ALL"
	Rules        []ScheduleAlertRule `json:"rules"`
}

// ScheduleAlertRule defines a single alert condition rule
type ScheduleAlertRule struct {
	Key      string `json:"key"`      // e.g. "subject", "description", "serviceName", "resourceName"
	Operator string `json:"operator"` // e.g. "CONTAINS", "EQUALS"
	Value    string `json:"value"`
}

// ScheduleMaintenanceResponse represents the full response for a maintenance window
type ScheduleMaintenanceResponse struct {
	Id                string                   `json:"id,omitempty"`
	UniqueId          string                   `json:"uniqueId"`
	Name              string                   `json:"name"`
	Description       string                   `json:"description"`
	RunRBA            bool                     `json:"runRBA"`
	InstallPatch      bool                     `json:"installPatch"`
	DontRunRBA        string                   `json:"dontRunRBA,omitempty"`
	DontInstallPatch  string                   `json:"dontInstallPatch,omitempty"`
	RunEscalateAction bool                     `json:"runEscalateAction"`
	CorrelateAlerts   bool                     `json:"correlateAlerts"`
	Schedule          ScheduleMaintenanceTime  `json:"schedule"`
	Devices           []ScheduleDeviceResponse `json:"devices,omitempty"`
	DeviceGroups      []ScheduleDeviceGroup    `json:"deviceGroups,omitempty"`
	Locations         []ScheduleLocation       `json:"locations,omitempty"`
	AlertConditions   *ScheduleAlertConditions `json:"alertConditions,omitempty"`
	Status            string                   `json:"status,omitempty"`
	CreatedTime       string                   `json:"createdTime,omitempty"`
	UpdatedTime       string                   `json:"updatedTime,omitempty"`
}

// ScheduleDeviceResponse represents a device in the GET response (richer than request)
type ScheduleDeviceResponse struct {
	Id             string                 `json:"id,omitempty"`
	GeneralInfo    map[string]interface{} `json:"generalInfo,omitempty"`
	ClientUniqueId string                 `json:"clientUniqueId,omitempty"`
	Type           string                 `json:"type,omitempty"`
}

// ScheduleMaintenanceResourcesRequest is used for POST/DELETE .../resources
type ScheduleMaintenanceResourcesRequest struct {
	Devices      []ScheduleDevice      `json:"devices,omitempty"`
	DeviceGroups []ScheduleDeviceGroup `json:"deviceGroups,omitempty"`
	Locations    []ScheduleLocation    `json:"locations,omitempty"`
}

// ScheduleMaintenanceResourceTypeResponse is the response from GET .../resources/{type}
type ScheduleMaintenanceResourceTypeResponse struct {
	Results      []map[string]interface{} `json:"results,omitempty"`
	TotalResults int                      `json:"totalResults,omitempty"`
	PageNo       int                      `json:"pageNo,omitempty"`
	PageSize     int                      `json:"pageSize,omitempty"`
	TotalPages   int                      `json:"totalPages,omitempty"`
	NextPage     bool                     `json:"nextPage,omitempty"`
}

// ScheduleMaintenanceCreateResponse is the minimal response from create/update
type ScheduleMaintenanceCreateResponse struct {
	UniqueId string `json:"uniqueId"`
}
