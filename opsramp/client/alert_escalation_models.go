// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package client

// AlertEscalationPolicy represents an alert escalation policy
type AlertEscalationPolicy struct {
	Id              string                    `json:"id,omitempty"`
	Name            string                    `json:"name"`
	Description     string                    `json:"description"`
	Precedence      *int                      `json:"precedence"`
	EscalationType  string                    `json:"escalationType"`
	PolicyType      string                    `json:"policyType"`
	EnabledMode     string                    `json:"enabledMode"`
	Active          bool                      `json:"active,omitempty"`
	AllClients      bool                      `json:"allClients"`
	IncludedClients []ClientRef               `json:"includedClients,omitempty"`
	Scope           *EscalationScope          `json:"scope"`
	Resources       []EscalationResource      `json:"resources"`
	Escalations     []EscalationLevel         `json:"escalations"`
	FilterCriteria  *EscalationFilterCriteria `json:"filterCriteria,omitempty"`
}

type ClientRef struct {
	UniqueId string `json:"uniqueId"`
}

// EscalationScope defines the scope of the policy
type EscalationScope struct {
	Uuid string `json:"uuid,omitempty"`
}

// EscalationResource defines a resource in scope
type EscalationResource struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}

// EscalationFilterCriteria defines filter criteria for the policy
type EscalationFilterCriteria struct {
	MatchingType        string `json:"matchingType,omitempty"`
	SearchQuery         string `json:"searchQuery,omitempty"`
	ResourceSearchQuery string `json:"resourceSearchQuery,omitempty"`
}

// EscalationLevel defines an escalation level/step
type EscalationLevel struct {
	WaitMins int    `json:"waitMins"`
	Action   string `json:"action"` // NOTIFICATION, INCIDENT

	Priority               string                `json:"priority,omitempty"` // LOW, MEDIUM, HIGH, CRITICAL
	RepeatFrequency        int                   `json:"repeatFrequency,omitempty"`
	NotifyLimitCount       int                   `json:"notifyLimitCount,omitempty"`
	NotificationType       string                `json:"notificationType,omitempty"` // basic, advanced
	NotificationTemplateId string                `json:"notificationTemplateId,omitempty"`
	Recipients             []EscalationRecipient `json:"recipients,omitempty"`

	Incident       *EscalationIncident       `json:"incident,omitempty"`
	UpdateIncident *EscalationUpdateIncident `json:"updateIncident,omitempty"`
}

// EscalationRecipient defines a notification recipient
type EscalationRecipient struct {
	Name string `json:"name,omitempty"`
	Id   string `json:"id"`
	Type string `json:"type"` // USER, USER_GROUP, ROSTER, USER_GROUP_DL
}

// EscalationIncident defines an incident to create
type EscalationIncident struct {
	Priority    string `json:"priority"`
	Subject     string `json:"subject"`
	Description string `json:"description"`

	AssigneeGroup  *EscalationUniqueRef `json:"assigneeGroup,omitempty"`
	AssignedUser   *EscalationUserRef   `json:"assignedUser,omitempty"`
	Category       *EscalationUniqueRef `json:"category,omitempty"`
	SubCategory    *EscalationUniqueRef `json:"subCategory,omitempty"`
	BusinessImpact *EscalationUniqueRef `json:"businessImpact,omitempty"`
	Urgency        *EscalationUniqueRef `json:"urgency,omitempty"`
	NotifyRoster   *EscalationRosterRef `json:"notifyRoster,omitempty"`

	AttachedArticles    []EscalationArticleRef `json:"attachedArticles,omitempty"`
	KnowledgeArticleIds []string               `json:"knowledgeArticleIds,omitempty"`

	Cc string `json:"cc,omitempty"`
}

// EscalationUniqueRef is a reference by uniqueId
type EscalationUniqueRef struct {
	UniqueId string `json:"uniqueId"`
}

// EscalationUserRef is a reference to a user
type EscalationUserRef struct {
	Id        string `json:"id"`
	LoginName string `json:"loginName,omitempty"`
}

// EscalationRosterRef is a reference to a roster
type EscalationRosterRef struct {
	Id string `json:"id"`
}

// EscalationArticleRef is a reference to a KB article
type EscalationArticleRef struct {
	Id      string `json:"id"`
	Subject string `json:"subject,omitempty"`
}

// EscalationUpdateIncident defines incident update settings
type EscalationUpdateIncident struct {
	// UpdateIncidentMode
	UpdateWhenAlertStateChange bool `json:"updateWhenAlertStateChange,omitempty"`
	UpdateForEveryRepeatAlert  bool `json:"updateForEveryRepeatAlert,omitempty"`

	UpdateWithRuleWhenAlertStateChange bool `json:"updateWithRuleWhenAlertStateChange,omitempty"`
	UpdateWithRuleForEveryRepeatAlert  bool `json:"updateWithRuleForEveryRepeatAlert,omitempty"`

	// UpdateIncidentSubjectMode
	UpdateIncidentSubject         bool `json:"updateIncidentSubject,omitempty"`
	UpdateIncidentSubjectWithRule bool `json:"updateIncidentSubjectWithRule,omitempty"`

	// AutoResolveIncidentMode
	AutoResolveIncident           bool `json:"autoResolveIncident,omitempty"`
	AutoResolveUnassignedIncident bool `json:"autoResolveUnassignedIncident,omitempty"`
	AutoHealWaitTime              int  `json:"autoHealWaitTime,omitempty"`

	UpdatePriorityByMLConfiguration bool                     `json:"updatePriorityByMLConfiguration,omitempty"`
	PriorityRules                   []EscalationPriorityRule `json:"priorityRules,omitempty"`
}

// EscalationPriorityRule defines a priority rule for incident updates
type EscalationPriorityRule struct {
	Key            string               `json:"key"`
	Operator       string               `json:"operator"`
	Value          string               `json:"value"`
	BusinessImpact *EscalationUniqueRef `json:"businessImpact,omitempty"`
	Urgency        *EscalationUniqueRef `json:"urgency,omitempty"`
	Priority       string               `json:"priority"` // Very Low, Low, Normal, High, Urgent
}
