// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package client

// AlertCorrelationPolicy represents an alert correlation policy
type AlertCorrelationPolicy struct {
	Id                       string                `json:"id,omitempty"`
	Name                     string                `json:"name"`
	Enabled                  bool                  `json:"enabled,omitempty"`
	EnabledMode              string                `json:"enabledMode,omitempty"`
	OrganizationMatchingType string                `json:"organizationMatchingType,omitempty"`
	IncludedClients          []string              `json:"includedClients,omitempty"`
	Precedence               int                   `json:"precedence,omitempty"`
	Type                     string                `json:"type,omitempty"`
	FilterQuery              string                `json:"filterQuery"`
	InferenceQuery           string                `json:"inferenceQuery"`
	Review                   bool                  `json:"review"`
	InferenceSubject         string                `json:"inferenceSubject,omitempty"`
	AlgorithmCorrelation     *AlgorithmCorrelation `json:"algorithmCorrelation,omitempty"`
	MachineLearning          *MachineLearning      `json:"machineLearning,omitempty"`
}

// AlgorithmCorrelation defines algorithm-based correlation settings
type AlgorithmCorrelation struct {
	MatchingConditions []MatchingCondition `json:"matchingConditions"`
	AlertsTimeWindow   string              `json:"alertsTimeWindow,omitempty"`
	AlertTrigger       *AlertTrigger       `json:"alertTrigger,omitempty"`
}

// AlertTrigger defines alert trigger conditions
type AlertTrigger struct {
	Rules    []AlertTriggerRule `json:"rules"`
	Duration int                `json:"duration,omitempty"`
}

// AlertTriggerRule defines a single trigger rule
type AlertTriggerRule struct {
	EntityName  string `json:"entityName"`
	Operator    string `json:"operator"`
	EntityValue string `json:"entityValue"`
}

// MachineLearning defines the machine learning settings for alert correlation
type MachineLearning struct {
	ContinuousLearning bool                `json:"continuousLearning,omitempty"`
	Topology           bool                `json:"topology,omitempty"`
	TopologyDepth      int                 `json:"topologyDepth,omitempty"`
	MatchingConditions []MatchingCondition `json:"matchingConditions"`
}

// MatchingCondition defines a matching condition for co-occurrence correlation
type MatchingCondition struct {
	Property  string `json:"property"`
	MatchType string `json:"matchType"`
}
