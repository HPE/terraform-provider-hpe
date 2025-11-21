// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"

	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/helpers"
)

type codeRepositories struct {
	Success bool `json:"success"`
	Data    []struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	} `json:"data"`
}

func DataSourceServiceNowWorkFlow() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus ServiceNow workflow data source.",
		ReadContext: dataSourceServiceNowWorkFlowRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeInt,
				Description:   "The ID of the ServiceNow integration",
				Optional:      true,
				ConflictsWith: []string{"name"},
				Computed:      true,
			},
			"name": {
				Type:          schema.TypeString,
				Description:   "The name of the ServiceNow workflow",
				Optional:      true,
				ConflictsWith: []string{"id"},
			},
			"integration_id": {
				Type:        schema.TypeInt,
				Description: "The ID of the ServiceNow integration",
				Required:    true,
			},
		},
	}
}

func dataSourceServiceNowWorkFlowRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	var diags diag.Diagnostics

	var name string
	if v, ok := d.Get("name").(string); ok {
		name = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
	}

	var integrationID int
	if v, ok := d.Get("integration_id").(int); ok {
		integrationID = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("integration_id", d.Get("integration_id")))
	}

	resp, err := client.Execute(&morpheus.Request{
		Method: "GET",
		Path: fmt.Sprintf(
			"/api/options/deleteApprovalServiceNowWorkflows?config.accountIntegrationId=%d",
			integrationID,
		),
		QueryParams: map[string]string{},
	})
	if err != nil {
		log.Println("API ERROR: ", err)

		return diag.FromErr(err)
	}
	log.Println("API RESPONSE:", resp)

	var itemResponsePayload codeRepositories
	if err := json.Unmarshal(resp.Body, &itemResponsePayload); err != nil {
		return diag.FromErr(err)
	}

	if itemResponsePayload.Data == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Data"))
	}

	foundWorkflow := false
	for _, v := range itemResponsePayload.Data {
		if v.Name == name {
			foundWorkflow = true
			d.SetId(convert.IntToString(v.Value))
		}
	}
	if !foundWorkflow {
		return diag.Errorf("Workflow named %s not found", name)
	}

	return diags
}
