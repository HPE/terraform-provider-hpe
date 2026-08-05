// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package client

// FirstResponsePolicy represents a first response policy
type FirstResponsePolicy struct {
	Id                       string                       `json:"id,omitempty"`
	Name                     string                       `json:"name"`
	EnabledMode              string                       `json:"enabledMode,omitempty"`
	OrganizationMatchingType string                       `json:"organizationMatchingType,omitempty"`
	IncludedClients          []string                     `json:"includedClients,omitempty"`
	FilterQuery              string                       `json:"filterQuery"`
	AttributeActions         *FirstResponseAttrActions    `json:"attributeActions,omitempty"`
	PatternActions           *FirstResponsePatternActions `json:"patternActions,omitempty"`
}

// FirstResponseAttrActions represents attribute-based actions
type FirstResponseAttrActions struct {
	ContinuousLearning bool                       `json:"continuousLearning"`
	Suppress           *FirstResponseAttrSuppress `json:"suppress"`
	Insights           *FirstResponseInsights     `json:"insights"`
	RunProcess         *FirstResponseRunProcess   `json:"runProcess,omitempty"`
}

// FirstResponseAttrSuppress represents suppress settings for attribute actions
type FirstResponseAttrSuppress struct {
	LearnedConfiguration bool `json:"learnedConfiguration"`
	SuppressDuration     int  `json:"suppressDuration,omitempty"`
}

// FirstResponseInsights represents insights settings
type FirstResponseInsights struct {
	CreatePrcInsights bool `json:"createPrcInsights"`
}

// FirstResponseRunProcess represents run process settings
type FirstResponseRunProcess struct {
	LearnedConfiguration bool     `json:"learnedConfiguration"`
	RunImmediately       bool     `json:"runImmediately"`
	ProcessIds           []string `json:"processIds,omitempty"`
}

// FirstResponsePatternActions represents pattern-based actions
type FirstResponsePatternActions struct {
	SeasonalityTimeFrame string                        `json:"seasonalityTimeFrame"`
	Suppress             *FirstResponsePatternSuppress `json:"suppress"`
}

// FirstResponsePatternSuppress represents suppress settings for pattern actions
type FirstResponsePatternSuppress struct {
	SeasonalAlerts bool `json:"seasonalAlerts"`
}
