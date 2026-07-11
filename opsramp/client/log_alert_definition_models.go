// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package client

// --- Log Alert Definition ---

// LogAlertDefinitionRequest is the top-level request for creating log alerts
type LogAlertDefinitionRequest struct {
	Alerts []LogAlertWrapper `json:"alerts"`
}

// LogAlertWrapper wraps a single log alert in the request/response
type LogAlertWrapper struct {
	Alert LogAlertDefinition `json:"alert"`
}

// LogAlertDefinition represents a log alert definition
type LogAlertDefinition struct {
	AlertID     string              `json:"alertId,omitempty"`
	Name        string              `json:"name"`
	Type        string              `json:"type"`
	TenantID    string              `json:"tenantId,omitempty"`
	AlertNoData string              `json:"alertNoData"`
	Conditions  []LogAlertCondition `json:"conditions"`
	Schedule    *LogAlertSchedule   `json:"schedule,omitempty"`
	Notification       LogAlertNotification `json:"notification"`
	Query              string               `json:"query,omitempty"`
	AdvancedQuery      string               `json:"advancedQuery,omitempty"`
	Status             string               `json:"status"`
	HealQuery          string               `json:"healQuery,omitempty"`
	EntityType         string               `json:"entityType"`
	Component          string               `json:"component,omitempty"`
	ResourceAttributes map[string]string    `json:"resourceAttributes,omitempty"`
	Labels             map[string]string    `json:"labels,omitempty"`
}

// LogAlertCondition defines a severity threshold for log alerts
type LogAlertCondition struct {
	Severity string `json:"severity"`
	Operator string `json:"operator"`
	Value    int    `json:"value"`
}

// LogAlertSchedule defines when the log alert is evaluated
type LogAlertSchedule struct {
	Pattern   LogAlertPattern `json:"pattern"`
	StartTime string          `json:"startTime,omitempty"`
	EndTime   string          `json:"endTime,omitempty"`
	Timezone  string          `json:"timezone,omitempty"`
}

// LogAlertPattern defines the schedule frequency
type LogAlertPattern struct {
	Type            string `json:"type"`
	RepeatFrequency int    `json:"repeatFrequency,omitempty"`
	WeekDays        string `json:"weekDays,omitempty"`
	DayOfMonth      string `json:"dayOfMonth,omitempty"`
	WeekIndex       string `json:"weekIndex,omitempty"`
	DayOfWeek       string `json:"dayOfWeek,omitempty"`
}

// LogAlertNotification defines the alert notification content
type LogAlertNotification struct {
	Subject     string `json:"subject"`
	Description string `json:"description"`
}

// LogAlertDefinitionResponse is the API response when creating log alerts
type LogAlertDefinitionResponse struct {
	Alerts []LogAlertWrapper `json:"alerts"`
	Errors []any             `json:"errors"`
}

// --- Shared ---

// NameValuePair is a generic label or attribute entry
type NameValuePair struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
