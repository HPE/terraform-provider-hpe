// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package client

// Alert Definition models for V3 API (metric-alerts and log-alerts)

// --- Metric Alert Definition ---

// MetricAlertDefinitionRequest represents the request body for creating/updating a metric alert definition
type MetricAlertDefinitionRequest struct {
	Name                 string                   `json:"name"`
	Query                string                   `json:"query"`
	AlertType            string                   `json:"alertType"`
	AlertThresholdType   string                   `json:"alertThresholdType"`
	AlertThresholdData   MetricAlertThresholdData `json:"alertThresholdData"`
	AlertTriggerDuration string                   `json:"alertTriggerDuration"`
	NoDataCondition      string                   `json:"noDataCondition,omitempty"`
	Labels               []NameValuePair          `json:"labels,omitempty"`
	Attributes           []NameValuePair          `json:"attributes"`
	Subject              string                   `json:"subject,omitempty"`
	Description          string                   `json:"description,omitempty"`
	EntityType           []string                 `json:"entityType"`
	Component            []string                 `json:"component,omitempty"`
	Status               bool                     `json:"status"`
	IsObsolete           *bool                    `json:"isObsolete,omitempty"`
}

// MetricAlertThresholdData holds warning and/or critical threshold conditions
type MetricAlertThresholdData struct {
	WarningCondition  string `json:"warningCondition,omitempty"`
	CriticalCondition string `json:"criticalCondition,omitempty"`
	Limit             int64  `json:"limit,omitempty"`
	Direction         string `json:"direction,omitempty"`
	LearningPeriod    string `json:"learningPeriod,omitempty"`
	StandardDeviation int64  `json:"standardDeviation,omitempty"`
}

// MetricAlertDefinitionResponse is the API response for create/update metric alert
type MetricAlertDefinitionResponse struct {
	Message                  string `json:"Message"`
	AlertDefinitionUniqueId  string `json:"Alert Definition UniqueId"`
	AlertDefinitionUniqueId2 string `json:"Alert Definition UniqueId "`
}

// GetID returns the alert definition ID from whichever response field is populated
func (r *MetricAlertDefinitionResponse) GetID() string {
	if r.AlertDefinitionUniqueId != "" {
		return r.AlertDefinitionUniqueId
	}

	return r.AlertDefinitionUniqueId2
}
