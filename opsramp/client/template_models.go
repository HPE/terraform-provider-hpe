// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package client

// Template models for the G2 monitoring template API (metricGroups based).
// The shapes below mirror the payload the New UI sends on create/update, which is
// richer than the abbreviated examples on the public API docs pages.

// TemplateStatus represents the activity status of a template.
type TemplateStatus struct {
	State   string `json:"state"` // ACTIVE, INACTIVE, EOL
	EolDate string `json:"eolDate,omitempty"`
}

// TemplateThreshold represents a single threshold on a linked template metric.
type TemplateThreshold struct {
	CriticalThreshold   string `json:"criticalThreshold,omitempty"`
	CriticalRepeatCount int    `json:"criticalRepeatCount,omitempty"`
	WarningThreshold    string `json:"warningThreshold,omitempty"`
	WarningRepeatCount  int    `json:"warningRepeatCount,omitempty"`
	WarningOperator     string `json:"warningOperator,omitempty"`
	CriticalOperator    string `json:"criticalOperator,omitempty"`
	RaiseAlert          bool   `json:"raiseAlert"`
	Availability        bool   `json:"availability"`
	LevelType           string `json:"levelType,omitempty"`
	CompName            string `json:"compName,omitempty"`
}

// TemplateMetricValue maps one OID/formula input feeding a SNMP/JMX/VMware metric's expression.
type TemplateMetricValue struct {
	Name          string `json:"name"`
	Oid           string `json:"oid,omitempty"`
	OidType       string `json:"oidType,omitempty"`
	DataType      string `json:"dataType,omitempty"`
	Formula       string `json:"formula,omitempty"`
	PollFrequency string `json:"pollFrequency,omitempty"`
}

// TemplateEnumeratedMapEntry maps one raw collected value to a display label.
type TemplateEnumeratedMapEntry struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// TemplateFormatMetricValue controls how a metric's collected value is interpreted/displayed.
type TemplateFormatMetricValue struct {
	UsedBy        string                       `json:"usedBy,omitempty"`
	EnumeratedMap []TemplateEnumeratedMapEntry `json:"enumeratedMap"`
	Alerts        bool                         `json:"alerts"`
	Graph         bool                         `json:"graph"`
}

// TemplateFormatMetricValueDetails wraps the display/formatting configuration of a metric.
type TemplateFormatMetricValueDetails struct {
	FormatMetricValue TemplateFormatMetricValue `json:"formatMetricValue"`
	Units             string                    `json:"units,omitempty"`
	Graph             bool                      `json:"graph"`
}

// TemplateMetric links an existing metric into a metric group, or (for SNMP/JMX/VMware
// metric groups) defines a new metric inline via MetricValues/Expression.
type TemplateMetric struct {
	Name                     string                            `json:"name"`
	DisplayName              string                            `json:"displayName,omitempty"`
	Description              string                            `json:"description,omitempty"`
	Category                 string                            `json:"category,omitempty"`
	ThresholdType            string                            `json:"thresholdType"`
	Thresholds               []TemplateThreshold               `json:"thresholds,omitempty"`
	Availability             bool                              `json:"availability"`
	AlertOnMissingComponents bool                              `json:"alertOnMissingComponents"`
	AlertSubject             string                            `json:"alertSubject,omitempty"`
	AlertContent             string                            `json:"alertContent,omitempty"`
	MetricValues             []TemplateMetricValue             `json:"metricValues"`
	Expression               string                            `json:"expression,omitempty"`
	FormatMetricValueDetails *TemplateFormatMetricValueDetails `json:"formatMetricValueDetails,omitempty"`
	ComponentFilters         []interface{}                     `json:"componentFilters"`
	// MetricCategory is nullable on the API and not yet exposed for configuration; always sent as null.
	MetricCategory *string `json:"metricCategory"`
}

// TemplateConfigurationParameter represents a configurable input of a metric group.
// Id is assigned by the API once the parameter exists on a saved template; VistaraCustomParamId
// references the (global, cross-template) catalog entry for well-known parameters.
type TemplateConfigurationParameter struct {
	Id                   int    `json:"id,omitempty"`
	VistaraCustomParamId int    `json:"vistaraCustomParamId,omitempty"`
	Name                 string `json:"name"`
	DisplayName          string `json:"displayName,omitempty"`
	CollectorType        string `json:"collectorType,omitempty"`
	DefaultValue         string `json:"defaultValue,omitempty"`
	Description          string `json:"description,omitempty"`
	Mandatory            bool   `json:"mandatory"`
	Configurable         bool   `json:"configurable"`
	InputType            string `json:"inputType,omitempty"`
	InputOptionValues    string `json:"inputOptionValues,omitempty"`
}

// TemplateComponentProperties describes how component instances are named/identified.
type TemplateComponentProperties struct {
	Id        string `json:"id"`
	Name      string `json:"name,omitempty"`
	Mode      string `json:"mode,omitempty"`
	Oid       string `json:"oid,omitempty"`
	LookupOid string `json:"lookupOid,omitempty"`
	Format    string `json:"format,omitempty"`
}

// TemplateDiscoveryFilter is a condition used to automatically match component instances,
// when a TemplateInstanceFilter's InstanceType is AUTOMATIC_INSTANCE.
type TemplateDiscoveryFilter struct {
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

// TemplateInstanceFilter scopes which discovered component instances a metric group applies to.
type TemplateInstanceFilter struct {
	InstanceType     string                    `json:"instanceType,omitempty"`
	DiscoveryFilters []TemplateDiscoveryFilter `json:"discoveryFilters"`
	InstanceNames    []string                  `json:"instanceNames"`
}

// TemplateAdditionalProperty is an extra SNMP OID collected alongside a metric group's metrics
// (e.g. for use in component naming or reference), not itself alertable.
// FormmatValue mirrors the API's own (misspelled) "formmatValue" JSON key.
type TemplateAdditionalProperty struct {
	Id            string                       `json:"id"`
	Name          string                       `json:"name"`
	OidType       string                       `json:"oidType,omitempty"`
	Oid           string                       `json:"oid,omitempty"`
	Operator      string                       `json:"operator"`
	Operand       string                       `json:"operand"`
	FormatValue   string                       `json:"formmatValue,omitempty"`
	EnumeratedMap []TemplateEnumeratedMapEntry `json:"enumeratedMap"`
	Units         string                       `json:"units"`
	DefaultValue  string                       `json:"defaultValue"`
	Description   string                       `json:"description"`
}

// TemplateOtelMetricsConfig configures OpenTelemetry-based metric collection.
type TemplateOtelMetricsConfig struct {
	ConfigType                    string `json:"configType,omitempty"`
	ReceiverConfigYaml            string `json:"receiverConfigYaml,omitempty"`
	AlertDefinitionsYaml          string `json:"alertDefinitionsYaml,omitempty"`
	ExternalExporter              bool   `json:"externalExporter"`
	ExternalExporterName          string `json:"externalExporterName,omitempty"`
	ExternalExporterCmdLineConfig string `json:"externalExporterCmdLineConfig,omitempty"`
}

// TemplateMetricGroup represents a G2 metric group (replaces the legacy "monitors" array).
type TemplateMetricGroup struct {
	Id                      string                           `json:"id,omitempty"`
	Name                    string                           `json:"name"`
	Frequency               int                              `json:"frequency"`
	AdapterType             string                           `json:"adapterType,omitempty"`
	ApplicationType         string                           `json:"applicationType,omitempty"`
	Script                  string                           `json:"script,omitempty"`
	Metrics                 []TemplateMetric                 `json:"metrics,omitempty"`
	ComponentProperties     *TemplateComponentProperties     `json:"componentProperties,omitempty"`
	IsMultiInstance         bool                             `json:"isMultiInstance"`
	AdditionalProperties    []TemplateAdditionalProperty     `json:"additionalProperties"`
	InputScopeComponent     bool                             `json:"inputScopeComponent"`
	DiscoveryMetric         string                           `json:"discoveryMetric,omitempty"`
	ConfigurationParameters []TemplateConfigurationParameter `json:"configurationParameters,omitempty"`
	OtelMetricsConfig       *TemplateOtelMetricsConfig       `json:"otelMetricsConfig,omitempty"`
	InstanceFilter          *TemplateInstanceFilter          `json:"instanceFilter,omitempty"`
	IsG2MonitorType         bool                             `json:"isG2MonitorType"`
}

// TemplateRequest is the request body for creating/updating a G2 template.
type TemplateRequest struct {
	Name          string `json:"name"`
	Description   string `json:"description,omitempty"`
	CollectorType string `json:"collectorType"`
	ResourceType  string `json:"resourceType"` // DEVICE, SYNTHETIC
	// CheckType: HTTPS, HTTP, DNS, IMAP, PING, RTT, POP3, SCRIPT, SMTP, SSL, TCP, UDP, FTP, SIP
	CheckType       string                `json:"checkType,omitempty"`
	Status          TemplateStatus        `json:"status"`
	Tags            string                `json:"tags,omitempty"`
	Generation      int                   `json:"generation"` // 2
	TemplateFamily  string                `json:"templateFamily,omitempty"`
	Notes           string                `json:"notes,omitempty"`
	Prerequisites   string                `json:"prerequisites,omitempty"`
	DeploymentType  string                `json:"deploymentType,omitempty"`
	ApplicationType []string              `json:"applicationType"`
	TemplateType    []string              `json:"templateType,omitempty"`
	MetricGroups    []TemplateMetricGroup `json:"metricGroups,omitempty"`
	// MonitorType is always "METRICS": this resource only manages metric-based G2 templates.
	MonitorType     string        `json:"monitorType"`
	Metrics         []interface{} `json:"metrics"`
	EventLogMetrics []interface{} `json:"eventLogMetrics"`
	OtherMetrics    []interface{} `json:"otherMetrics"`
	ClonedTemplate  bool          `json:"clonedTemplate"`
}

// TemplateCreateResponse is the minimal 201 response returned by the create endpoint.
type TemplateCreateResponse struct {
	TemplateId string `json:"templateId"`
}

// Template is the full template representation returned by GET.
type Template struct {
	Id                  string                `json:"id"`
	Name                string                `json:"name"`
	Description         string                `json:"description,omitempty"`
	CollectorType       string                `json:"collectorType"`
	ResourceType        string                `json:"resourceType"`
	CheckType           string                `json:"checkType,omitempty"`
	Status              TemplateStatus        `json:"status"`
	Tags                string                `json:"tags,omitempty"`
	Generation          int                   `json:"generation"`
	CreatedDate         string                `json:"createdDate,omitempty"`
	UpdatedDate         string                `json:"updatedDate,omitempty"`
	Scope               string                `json:"scope,omitempty"`
	TemplateFamily      string                `json:"templateFamily,omitempty"`
	Notes               string                `json:"notes,omitempty"`
	Prerequisites       string                `json:"prerequisites,omitempty"`
	DeploymentType      string                `json:"deploymentType,omitempty"`
	AppName             string                `json:"appName,omitempty"`
	NativeType          string                `json:"nativeType,omitempty"`
	AssignedToResources bool                  `json:"assignedToResources"`
	ApplicationType     []string              `json:"applicationType,omitempty"`
	TemplateType        []string              `json:"templateType,omitempty"`
	MetricGroups        []TemplateMetricGroup `json:"metricGroups,omitempty"`
	SdkTag              string                `json:"sdkTag,omitempty"`
	MonTemplateVer      int                   `json:"monTemplateVer,omitempty"`
}
