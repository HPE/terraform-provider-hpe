// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package setting

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"

	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/helpers"
)

func ResourceSettingMonitoring() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus monitoring setting resource.",
		CreateContext: resourceSettingMonitoringCreate,
		ReadContext:   resourceSettingMonitoringRead,
		UpdateContext: resourceSettingMonitoringUpdate,
		DeleteContext: resourceSettingMonitoringDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the monitoring setting",
				Computed:    true,
			},
			"morpheus_auto_create_checks": {
				Type:        schema.TypeBool,
				Description: "When enabled a Monitoring Check will automatically be create for Instances and Apps",
				Optional:    true,
				Computed:    true,
			},
			"morpheus_availability_time_frame": {
				Type:        schema.TypeInt,
				Description: "The number of days availability should be calculated for",
				Optional:    true,
				Computed:    true,
			},
			"morpheus_availability_precision": {
				Type:         schema.TypeInt,
				Description:  "The number of decimal places availability should be displayed in, can be anywhere between 0 and 5",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntInSlice([]int{0, 1, 2, 3, 4, 5}),
			},
			"morpheus_default_check_interval": {
				Type: schema.TypeInt,
				Description: "The default interval in minutes to use when creating new checks " +
					"(1, 2, 3, 4, 5, 10, 15, 20, 25, 30, 45, 60, 120, 180)",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntInSlice([]int{1, 2, 3, 4, 5, 10, 15, 20, 25, 30, 45, 60, 120, 180}),
			},
			"servicenow_monitoring_enabled": {
				Type:         schema.TypeBool,
				Description:  "Whether the ServiceNow monitoring integration is enabled",
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"servicenow_integration_id"},
			},
			"servicenow_integration_id": {
				Type:         schema.TypeInt,
				Description:  "The id of the ServiceNow monitoring integration",
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"servicenow_monitoring_enabled"},
			},
			"servicenow_new_incident_action": {
				Type:         schema.TypeString,
				Description:  "The Service Now action to take when a Morpheus incident is created (create, none)",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"create", "none"}, false),
				RequiredWith: []string{"servicenow_monitoring_enabled"},
			},
			"servicenow_close_incident_action": {
				Type:         schema.TypeString,
				Description:  "The Service Now action to take when a Morpheus incident is closed (activity, close, none)",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"activity", "close", "none"}, false),
				RequiredWith: []string{"servicenow_monitoring_enabled"},
			},
			"servicenow_severity_info_impact": {
				Type:         schema.TypeString,
				Description:  "The ServiceNow impact level to map to the Morpheus info severity (high, medium, low)",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"high", "medium", "low"}, false),
				RequiredWith: []string{"servicenow_monitoring_enabled"},
			},
			"servicenow_severity_warning_impact": {
				Type:         schema.TypeString,
				Description:  "The ServiceNow impact level to map to the Morpheus warning severity (high, medium, low)",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"high", "medium", "low"}, false),
				RequiredWith: []string{"servicenow_monitoring_enabled"},
			},
			"servicenow_severity_critical_impact": {
				Type:         schema.TypeString,
				Description:  "The ServiceNow impact level to map to the Morpheus critical severity (high, medium, low)",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"high", "medium", "low"}, false),
				RequiredWith: []string{"servicenow_monitoring_enabled"},
			},
			"new_relic_monitoring_enabled": {
				Type:        schema.TypeBool,
				Description: "Whether the New Relic monitoring integration is enabled",
				Optional:    true,
				Computed:    true,
			},
			"new_relic_license_key": {
				Type:        schema.TypeString,
				Description: "The New Relic license key",
				Optional:    true,
				Computed:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceSettingMonitoringCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	var diags diag.Diagnostics

	monitoringSettings := make(map[string]any)

	var morpheusAutoCreateChecks bool
	if v, ok := d.Get("morpheus_auto_create_checks").(bool); ok {
		morpheusAutoCreateChecks = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("morpheus_auto_create_checks", d.Get("morpheus_auto_create_checks")))
	}
	monitoringSettings["autoManageChecks"] = morpheusAutoCreateChecks

	availabilityTimeFrame, availabilityTimeFrameok := d.GetOk("morpheus_availability_time_frame")
	if availabilityTimeFrameok {
		var availabilityTimeFrameInt int
		if v, ok := availabilityTimeFrame.(int); ok {
			availabilityTimeFrameInt = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("morpheus_availability_time_frame", availabilityTimeFrame))
		}
		monitoringSettings["availabilityTimeFrame"] = availabilityTimeFrameInt
	}

	availabilityPrecision, availabilityPrecisionok := d.GetOk("morpheus_availability_precision")
	if availabilityPrecisionok {
		var availabilityPrecisionInt int
		if v, ok := availabilityPrecision.(int); ok {
			availabilityPrecisionInt = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("morpheus_availability_precision", availabilityPrecision))
		}
		monitoringSettings["availabilityPrecision"] = availabilityPrecisionInt
	}

	defaultCheckInterval, defaultCheckIntervalok := d.GetOk("morpheus_default_check_interval")
	if defaultCheckIntervalok {
		var defaultCheckIntervalInt int
		if v, ok := defaultCheckInterval.(int); ok {
			defaultCheckIntervalInt = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("morpheus_default_check_interval", defaultCheckInterval))
		}
		monitoringSettings["defaultCheckInterval"] = defaultCheckIntervalInt
	}

	serviceNowSettings := make(map[string]any)

	serviceNowNewIntegrationId, serviceNowNewIntegrationIdok := d.GetOk("servicenow_integration_id")
	if serviceNowNewIntegrationIdok {
		var serviceNowNewIntegrationIdInt int
		if v, ok := serviceNowNewIntegrationId.(int); ok {
			serviceNowNewIntegrationIdInt = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("servicenow_integration_id", serviceNowNewIntegrationId))
		}
		serviceNowIntegration := make(map[string]any)
		serviceNowIntegration["id"] = serviceNowNewIntegrationIdInt

		var serviceNowMonitoringEnabled bool
		if v, ok := d.Get("servicenow_monitoring_enabled").(bool); ok {
			serviceNowMonitoringEnabled = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError(
				"servicenow_monitoring_enabled",
				d.Get("servicenow_monitoring_enabled"),
			))
		}
		serviceNowSettings["integration"] = serviceNowIntegration
		serviceNowSettings["enabled"] = serviceNowMonitoringEnabled
		monitoringSettings["serviceNow"] = serviceNowSettings
	}

	serviceNowNewIncidentAction, serviceNowNewIncidentActionok := d.GetOk("servicenow_new_incident_action")
	if serviceNowNewIncidentActionok {
		var serviceNowNewIncidentActionStr string
		if v, ok := serviceNowNewIncidentAction.(string); ok {
			serviceNowNewIncidentActionStr = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("servicenow_new_incident_action", serviceNowNewIncidentAction))
		}
		serviceNowSettings["newIncidentAction"] = serviceNowNewIncidentActionStr
		monitoringSettings["serviceNow"] = serviceNowSettings
	}

	serviceNowCloseIncidentAction, serviceNowCloseIncidentActionok := d.GetOk("servicenow_close_incident_action")
	if serviceNowCloseIncidentActionok {
		var serviceNowCloseIncidentActionStr string
		if v, ok := serviceNowCloseIncidentAction.(string); ok {
			serviceNowCloseIncidentActionStr = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("servicenow_close_incident_action", serviceNowCloseIncidentAction))
		}
		serviceNowSettings["closeIncidentAction"] = serviceNowCloseIncidentActionStr
		monitoringSettings["serviceNow"] = serviceNowSettings
	}

	serviceNowInfoImpact, serviceNowInfoImpactok := d.GetOk("servicenow_severity_info_impact")
	if serviceNowInfoImpactok {
		var serviceNowInfoImpactStr string
		if v, ok := serviceNowInfoImpact.(string); ok {
			serviceNowInfoImpactStr = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("servicenow_severity_info_impact", serviceNowInfoImpact))
		}
		serviceNowSettings["infoMapping"] = serviceNowInfoImpactStr
		monitoringSettings["serviceNow"] = serviceNowSettings
	}

	serviceNowWarningImpact, serviceNowWarningImpactok := d.GetOk("servicenow_severity_warning_impact")
	if serviceNowWarningImpactok {
		var serviceNowWarningImpactStr string
		if v, ok := serviceNowWarningImpact.(string); ok {
			serviceNowWarningImpactStr = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("servicenow_severity_warning_impact", serviceNowWarningImpact))
		}
		serviceNowSettings["warningMapping"] = serviceNowWarningImpactStr
		monitoringSettings["serviceNow"] = serviceNowSettings
	}

	serviceNowCriticalImpact, serviceNowCriticalImpactok := d.GetOk("servicenow_severity_critical_impact")
	if serviceNowCriticalImpactok {
		var serviceNowCriticalImpactStr string
		if v, ok := serviceNowCriticalImpact.(string); ok {
			serviceNowCriticalImpactStr = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("servicenow_severity_critical_impact", serviceNowCriticalImpact))
		}
		serviceNowSettings["criticalMapping"] = serviceNowCriticalImpactStr
		monitoringSettings["serviceNow"] = serviceNowSettings
	}

	newRelicSettings := make(map[string]any)

	var newRelicMonitoringEnabled bool
	if v, ok := d.Get("new_relic_monitoring_enabled").(bool); ok {
		newRelicMonitoringEnabled = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError(
			"new_relic_monitoring_enabled",
			d.Get("new_relic_monitoring_enabled"),
		))
	}
	newRelicSettings["enabled"] = newRelicMonitoringEnabled
	monitoringSettings["newRelic"] = newRelicSettings

	newRelicLicenseKey, newRelicLicenseKeyok := d.GetOk("new_relic_license_key")
	if newRelicLicenseKeyok {
		var newRelicLicenseKeyStr string
		if v, ok := newRelicLicenseKey.(string); ok {
			newRelicLicenseKeyStr = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("new_relic_license_key", newRelicLicenseKey))
		}
		newRelicSettings["licenseKey"] = newRelicLicenseKeyStr
		monitoringSettings["newRelic"] = newRelicSettings
	}

	req := &morpheus.Request{
		Body: map[string]any{
			"monitoringSettings": monitoringSettings,
		},
	}

	jsonRequest, _ := json.Marshal(req.Body)
	log.Printf("API JSON REQUEST: %s", string(jsonRequest))

	resp, err := client.UpdateMonitoringSettings(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	var result *morpheus.UpdateMonitoringSettingsResult
	if v, ok := resp.Result.(*morpheus.UpdateMonitoringSettingsResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.MonitoringSettings == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("MonitoringSettings"))
	}

	d.SetId(convert.Int64ToString(1))

	diags = append(diags, resourceSettingMonitoringRead(ctx, d, meta)...)

	return diags
}

func resourceSettingMonitoringRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	var diags diag.Diagnostics

	var resp *morpheus.Response
	var err error

	resp, err = client.GetMonitoringSettings(&morpheus.Request{})
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %s", resp, err)
			log.Printf("Forcing recreation of resource")
			d.SetId("")

			return diags
		}

		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	var result *morpheus.GetMonitoringSettingsResult
	if v, ok := resp.Result.(*morpheus.GetMonitoringSettingsResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.MonitoringSettings == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("MonitoringSettings"))
	}

	monitoringSetting := result.MonitoringSettings
	d.SetId(convert.Int64ToString(1))
	d.Set("morpheus_auto_create_checks", monitoringSetting.AutoManageChecks)
	d.Set("morpheus_availability_time_frame", monitoringSetting.AvailabilityTimeFrame)
	d.Set("morpheus_availability_precision", monitoringSetting.AvailabilityPrecision)
	d.Set("morpheus_default_check_interval", monitoringSetting.DefaultCheckInterval)
	d.Set("servicenow_monitoring_enabled", monitoringSetting.ServiceNow.Enabled)
	d.Set("servicenow_integration_id", monitoringSetting.ServiceNow.Integration.ID)
	d.Set("servicenow_new_incident_action", monitoringSetting.ServiceNow.NewIncidentAction)
	d.Set("servicenow_close_incident_action", monitoringSetting.ServiceNow.CloseIncidentAction)
	d.Set("servicenow_severity_info_impact", monitoringSetting.ServiceNow.InfoMapping)
	d.Set("servicenow_severity_warning_impact", monitoringSetting.ServiceNow.WarningMapping)
	d.Set("servicenow_severity_critical_impact", monitoringSetting.ServiceNow.CriticalMapping)
	d.Set("new_relic_monitoring_enabled", monitoringSetting.NewRelic.Enabled)
	d.Set("new_relic_license_key", monitoringSetting.NewRelic.LicenseKey)

	return diags
}

func resourceSettingMonitoringUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	monitoringSettings := make(map[string]any)

	if d.HasChange("morpheus_auto_create_checks") {
		var morpheusAutoCreateChecks bool
		if v, ok := d.Get("morpheus_auto_create_checks").(bool); ok {
			morpheusAutoCreateChecks = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("morpheus_auto_create_checks", d.Get("morpheus_auto_create_checks")))
		}
		monitoringSettings["autoManageChecks"] = morpheusAutoCreateChecks
	}

	if d.HasChange("morpheus_availability_time_frame") {
		var morpheusAvailabilityTimeFrame int
		if v, ok := d.Get("morpheus_availability_time_frame").(int); ok {
			morpheusAvailabilityTimeFrame = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError(
				"morpheus_availability_time_frame",
				d.Get("morpheus_availability_time_frame")))
		}
		monitoringSettings["availabilityTimeFrame"] = morpheusAvailabilityTimeFrame
	}

	if d.HasChange("morpheus_availability_precision") {
		var morpheusAvailabilityPrecision int
		if v, ok := d.Get("morpheus_availability_precision").(int); ok {
			morpheusAvailabilityPrecision = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError(
				"morpheus_availability_precision",
				d.Get("morpheus_availability_precision")))
		}
		monitoringSettings["availabilityPrecision"] = morpheusAvailabilityPrecision
	}

	if d.HasChange("morpheus_default_check_interval") {
		var morpheusDefaultCheckInterval int
		if v, ok := d.Get("morpheus_default_check_interval").(int); ok {
			morpheusDefaultCheckInterval = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError(
				"morpheus_default_check_interval",
				d.Get("morpheus_default_check_interval")))
		}
		monitoringSettings["defaultCheckInterval"] = morpheusDefaultCheckInterval
	}

	serviceNowSettings := make(map[string]any)

	if d.HasChange("servicenow_monitoring_enabled") {
		var serviceNowMonitoringEnabled bool
		if v, ok := d.Get("servicenow_monitoring_enabled").(bool); ok {
			serviceNowMonitoringEnabled = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError(
				"servicenow_monitoring_enabled",
				d.Get("servicenow_monitoring_enabled"),
			))
		}
		serviceNowSettings["enabled"] = serviceNowMonitoringEnabled
		monitoringSettings["serviceNow"] = serviceNowSettings
	}

	if d.HasChange("servicenow_integration_id") {
		var serviceNowIntegrationId int
		if v, ok := d.Get("servicenow_integration_id").(int); ok {
			serviceNowIntegrationId = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("servicenow_integration_id", d.Get("servicenow_integration_id")))
		}
		serviceNowIntegration := make(map[string]any)
		serviceNowIntegration["id"] = serviceNowIntegrationId
		serviceNowSettings["integration"] = serviceNowIntegration
		monitoringSettings["serviceNow"] = serviceNowSettings
	}

	if d.HasChange("servicenow_new_incident_action") {
		var serviceNowNewIncidentAction string
		if v, ok := d.Get("servicenow_new_incident_action").(string); ok {
			serviceNowNewIncidentAction = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError(
				"servicenow_new_incident_action",
				d.Get("servicenow_new_incident_action")))
		}
		serviceNowSettings["newIncidentAction"] = serviceNowNewIncidentAction
		monitoringSettings["serviceNow"] = serviceNowSettings
	}

	if d.HasChange("servicenow_close_incident_action") {
		var serviceNowCloseIncidentAction string
		if v, ok := d.Get("servicenow_close_incident_action").(string); ok {
			serviceNowCloseIncidentAction = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError(
				"servicenow_close_incident_action",
				d.Get("servicenow_close_incident_action")))
		}
		serviceNowSettings["closeIncidentAction"] = serviceNowCloseIncidentAction
		monitoringSettings["serviceNow"] = serviceNowSettings
	}

	if d.HasChange("servicenow_severity_info_impact") {
		var serviceNowSeverityInfoImpact string
		if v, ok := d.Get("servicenow_severity_info_impact").(string); ok {
			serviceNowSeverityInfoImpact = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError(
				"servicenow_severity_info_impact",
				d.Get("servicenow_severity_info_impact")))
		}
		serviceNowSettings["infoMapping"] = serviceNowSeverityInfoImpact
		monitoringSettings["serviceNow"] = serviceNowSettings
	}

	if d.HasChange("servicenow_severity_warning_impact") {
		var serviceNowSeverityWarningImpact string
		if v, ok := d.Get("servicenow_severity_warning_impact").(string); ok {
			serviceNowSeverityWarningImpact = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError(
				"servicenow_severity_warning_impact",
				d.Get("servicenow_severity_warning_impact")))
		}
		serviceNowSettings["warningMapping"] = serviceNowSeverityWarningImpact
		monitoringSettings["serviceNow"] = serviceNowSettings
	}

	if d.HasChange("servicenow_severity_critical_impact") {
		var serviceNowSeverityCriticalImpact string
		if v, ok := d.Get("servicenow_severity_critical_impact").(string); ok {
			serviceNowSeverityCriticalImpact = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError(
				"servicenow_severity_critical_impact",
				d.Get("servicenow_severity_critical_impact")))
		}
		serviceNowSettings["criticalMapping"] = serviceNowSeverityCriticalImpact
		monitoringSettings["serviceNow"] = serviceNowSettings
	}

	newRelicSettings := make(map[string]any)

	if d.HasChange("new_relic_monitoring_enabled") {
		var newRelicMonitoringEnabled bool
		if v, ok := d.Get("new_relic_monitoring_enabled").(bool); ok {
			newRelicMonitoringEnabled = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError(
				"new_relic_monitoring_enabled",
				d.Get("new_relic_monitoring_enabled"),
			))
		}
		newRelicSettings["enabled"] = newRelicMonitoringEnabled
		monitoringSettings["newRelic"] = newRelicSettings
	}

	if d.HasChange("new_relic_license_key") {
		var newRelicLicenseKey string
		if v, ok := d.Get("new_relic_license_key").(string); ok {
			newRelicLicenseKey = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("new_relic_license_key", d.Get("new_relic_license_key")))
		}
		newRelicSettings["licenseKey"] = newRelicLicenseKey
		monitoringSettings["newRelic"] = newRelicSettings
	}

	req := &morpheus.Request{
		Body: map[string]any{
			"monitoringSettings": monitoringSettings,
		},
	}

	jsonRequest, _ := json.Marshal(req.Body)
	log.Printf("API JSON REQUEST: %s", string(jsonRequest))

	resp, err := client.UpdateMonitoringSettings(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	var result *morpheus.UpdateMonitoringSettingsResult
	if v, ok := resp.Result.(*morpheus.UpdateMonitoringSettingsResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.MonitoringSettings == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("MonitoringSettings"))
	}

	d.SetId(convert.Int64ToString(1))

	return resourceSettingMonitoringRead(ctx, d, meta)
}

func resourceSettingMonitoringDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics

	d.SetId("")

	return diags
}
