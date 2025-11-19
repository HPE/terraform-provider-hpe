// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/helpers"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func DataSourceIntegrationGit() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus git integration data source.",
		ReadContext: dataSourceIntegrationGitRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeInt,
				Description:   "The ID of the Morpheus git integration",
				Optional:      true,
				ConflictsWith: []string{"name"},
				Computed:      true,
			},
			"name": {
				Type:          schema.TypeString,
				Description:   "The name of the Morpheus git integration.",
				Optional:      true,
				ConflictsWith: []string{"id"},
			},
			"repository_ids": {
				Computed:    true,
				Type:        schema.TypeMap,
				Description: "A map of git repository ids for use with integrations that reference a git repository",
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
		},
	}
}

func dataSourceIntegrationGitRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	var name string
	if v, ok := d.Get("name").(string); ok {
		name = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
	}

	var id int
	if v, ok := d.Get("id").(int); ok {
		id = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("id", d.Get("id")))
	}

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == 0 && name != "" {
		resp, err = client.FindIntegrationByName(name)
	} else if id != 0 {
		resp, err = client.GetIntegration(int64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Integration cannot be read without name or id")
	}
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %v", resp, err)

			return nil
		} else {
			log.Printf("API FAILURE: %s - %v", resp, err)

			return diag.FromErr(err)
		}
	}
	log.Printf("API RESPONSE: %s", resp)

	// store resource data
	var result *morpheus.GetIntegrationResult
	if v, ok := resp.Result.(*morpheus.GetIntegrationResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("GetIntegrationResult", resp.Result))
	}

	if result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("GetIntegrationResult"))
	}

	integration := result.Integration
	if integration == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Integration"))
	}

	d.SetId(convert.Int64ToString(integration.ID))
	d.Set("name", integration.Name)
	resp, err = client.Execute(&morpheus.Request{
		Method:      "GET",
		Path:        fmt.Sprintf("/api/options/codeRepositories?integrationId=%d", integration.ID),
		QueryParams: map[string]string{},
	})
	if err != nil {
		log.Println("API ERROR: ", err)
	}
	log.Println("API RESPONSE:", resp)
	repoIDs := make(map[string]int)

	var itemResponsePayload CodeRepositories
	json.Unmarshal(resp.Body, &itemResponsePayload)
	if itemResponsePayload.Data != nil {
		for _, v := range itemResponsePayload.Data {
			repoIDs[v.Name] = v.Value
		}
	}
	d.Set("repository_ids", repoIDs)

	return diags
}

type CodeRepositories struct {
	Success bool `json:"success"`
	Data    []struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	} `json:"data"`
}
